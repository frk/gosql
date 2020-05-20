package main

// TODO: mysqldbInfo & mysqlTypeCheck

import (
	"database/sql"
	"fmt"
	"go/token"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/frk/gosql/internal/errors"

	"github.com/lib/pq"
)

// pgdbInfo handles the connection pool to the target postgres database
// and holds additional information about that database.
//
// pgdbInfo is NOT safe for concurrent use, an instance of pgdbInfo is intended
// to be reused by separate runs of the type checker although not concurrently.
type pgdbInfo struct {
	// The connection pool to the target database.
	db *sql.DB
	// The name of the current database. (intended mainly for error reporting)
	name string
	// The url used to open connections to the database.
	// Used also as the key for caching the catalogue.
	url string
	// The catalog for the target database.
	cat *pgCatalog
}

// init creates a new connection pool to the url specified database and loads
// the catalog information.
func (pg *pgdbInfo) init() (err error) {
	if pg.db, err = sql.Open("postgres", pg.url); err != nil {
		return err
	} else if err := pg.db.Ping(); err != nil {
		return err
	}

	const selectDBName = `SELECT current_database()` //`
	if err := pg.db.QueryRow(selectDBName).Scan(&pg.name); err != nil {
		return err
	}

	pg.cat = new(pgCatalog)
	if err = pg.cat.load(pg.db, pg.url); err != nil {
		return err
	}
	return nil
}

// close closes the underlying connection pool.
func (pg *pgdbInfo) close() error {
	return pg.db.Close()
}

// pgTypeCheck holds the state of the postgres specific type checker.
type pgTypeCheck struct {
	fset *token.FileSet
	pg   *pgdbInfo
	ti   *targetInfo

	// The target relation.
	rel *pgRelationInfo
	// A map of column qualifier (alias or relname or "") to the
	// relations denoted by said qualifier.
	relmap map[string]*pgRelationInfo
	// A map that holds the set of joined relations.
	joinrels map[relId]*pgRelationInfo
}

func (c *pgTypeCheck) init() (err error) {
	c.ti.searchConditionFieldColumns = make(map[*searchConditionField]*pgcolumn)

	c.joinrels = make(map[relId]*pgRelationInfo)
	c.relmap = make(map[string]*pgRelationInfo)
	if c.rel, err = c.loadRelation(c.ti.dataField.relId); err != nil {
		return err
	}

	// Map the "" (empty string) key to the target relation, this will allow
	// columns, constraints, and indexes that were specified without a qualifier
	// to be associated with this target relation.
	c.relmap[""] = c.rel

	// Load all the joined relations if a joinStruct is present.
	if c.ti.query != nil && c.ti.query.joinStruct != nil {
		jb := c.ti.query.joinStruct
		if len(jb.relId.name) > 0 {
			rel, err := c.loadRelation(jb.relId)
			if err != nil {
				return err
			}
			c.joinrels[jb.relId] = rel
		}
		for _, join := range jb.items {
			rel, err := c.loadRelation(join.relId)
			if err != nil {
				return err
			}
			c.joinrels[join.relId] = rel
		}
	}
	return nil
}

// run initializes and executes the type checker.
func (c *pgTypeCheck) run() (err error) {
	if err := c.init(); err != nil {
		return err
	}

	if c.ti.query != nil {
		return c.checkQueryStruct()
	} else if c.ti.filter != nil {
		return c.checkFilterStruct()
	}

	panic("nothing to db-check")
	return nil
}

// checkFilterStruct type checks the target filter struct.
func (c *pgTypeCheck) checkFilterStruct() (err error) {
	checks := []func() error{
		c.checkFilterDataField, // TODO
		c.checkTextSearchColId, // TODO
	}
	for i := 0; i < len(checks); i++ {
		if err := checks[i](); err != nil {
			return err
		}
	}
	return nil
}

// checkQueryStruct type checks the target query struct.
func (c *pgTypeCheck) checkQueryStruct() (err error) {
	checks := []func() error{
		c.checkForceList,
		c.checkDefaultList,
		c.checkOrderByList,

		c.checkReturnList,
		c.checkResultField,
		c.checkQueryDataField,

		c.checkJoinBlock,
		c.checkWhereBlock, // TODO
		c.checkOnConflictBlock,
	}
	for i := 0; i < len(checks); i++ {
		if err := checks[i](); err != nil {
			return err
		}
	}

	////////////////////////////////////////////////////////////////////////

	if c.ti.query.kind == queryKindInsert {
		// TODO(mkopriva): if this is an insert make sure that all columns that
		// do not have a DEFAULT set but do have a NOT NULL set, have a corresponding
		// field in the relation's dataType. (keep in mind that such a column could
		// be also set by a BEFORE TRIGGER, so maybe instead of erroring only warning
		// would be thrown, or make this check optional, something that can be turned
		// on/off...?)
	}

	// TODO(mkopriva): if this is an insert or update (i.e. write) the column
	// generated from the field tags cannot contain duplicate columns. Conversely
	// if this is a read query (select, or I,U,D with returning) the columns do
	// not have to be unique.

	// TODO(mkopriva): if this is an UPDATE but none of the columns associated
	// with the available fields constitute a complete primary key of the table,
	// an error should be returned as then there's no way to properly match the
	// data instances with specific rows in the table.
	return nil
}

// checkForceList checks the columns listed in the gosql.Force directive's tag.
//
// CHECKLIST:
//  ✅ Each column MUST be present in one of the loaded relations.
func (c *pgTypeCheck) checkForceList() error {
	list := c.ti.query.forceList
	if list == nil {
		return nil
	}

	for _, colId := range list.items {
		if errcode := c.checkColumnExists(colId); errcode > 0 {
			return c.newError(errcode, colId, nil, list)
		}
	}
	return nil
}

// checkOrderByList checks the columns listed in the gosql.OrderBy directive's tag.
//
// CHECKLIST:
//  ✅ Each column MUST be present in one of the loaded relations.
func (c *pgTypeCheck) checkOrderByList() error {
	list := c.ti.query.orderByList
	if list == nil {
		return nil
	}

	for _, item := range list.items {
		if errcode := c.checkColumnExists(item.colId); errcode > 0 {
			return c.newError(errcode, item.colId, nil, list)
		}
	}
	return nil
}

// checkDefaultList checks the columns listed in the gosql.Default directive's tag.
//
// CHECKLIST:
//  ✅ Each column MUST be present in the TARGET relation.
//  ✅ Each column MUST have a DEFAULT set.
//  ✅ If a column has a qualifier it MUST match the alias,
//     or name, of the target relation.
func (c *pgTypeCheck) checkDefaultList() error {
	list := c.ti.query.defaultList
	if list == nil {
		return nil
	}

	for _, item := range list.items {
		if len(item.qual) > 0 {
			relId := c.ti.query.dataField.relId
			if item.qual != relId.alias && item.qual != relId.name {
				return c.newError(errBadColumnQualifier, item, nil, list)
			}
		}

		if col := c.rel.findColumn(item.name); col == nil {
			return c.newError(errNoRelationColumn, item, nil, list)
		} else if !col.hasdefault {
			return c.newError(errNoColumnDefault, item, col, list)
		}
	}
	return nil
}

// checkReturnList checks the columns listed in the gosql.Return directive's tag.
//
// CHECKLIST:
// - If "*" tag was used:
//  ✅ Each field of the target data type MAY have a corresponding column
//     in one of the loaded relations (denoted by the field's tag).
//  ✅ Each field of the target data type that has a corresponding column MUST
//     be of a type that is readable from a value of that column.
//
// - If "<column_list>" tag was used:
//  ✅ Each listed column MUST be present in the TARGET relation.
//  ✅ Each listed column MUST have a corresponding field in the target data type.
//  ✅ Each listed column's qualifier, if it has one, MUST match the alias, or name,
//     of the TARGET relation.
//  ✅ Each listed column's corresponding field MUST be of a type that IS READABLE
//     from a value of that column's type.
func (c *pgTypeCheck) checkReturnList() error {
	list := c.ti.query.returnList
	if list == nil {
		return nil
	}

	var strict bool
	var fields []*fieldInfo

	if list.all {
		strict = false
		fields = c.ti.query.dataField.data.fields
	} else {
		strict = true
		for _, colId := range list.items {
			fi, errcode := c.findColumnField(colId, strict)
			if errcode > 0 {
				return c.newError(errcode, colId, nil, list)
			} else if fi != nil {
				fields = append(fields, fi)
			}
		}
	}

	for _, field := range fields {
		if err := c.checkColumnRead(field, strict); err != nil {
			return err
		}
	}
	return nil
}

// checkResultField checks the fields of the resultField struct data type.
//
// CHECKLIST:
//  ✅ Each field of the result data type MUST have a corresponding column in
//     one of the loaded relations.
//  ✅ Each field of the result data type MUST be of a type that is readable
//     from a value of the corresponding column.
func (c *pgTypeCheck) checkResultField() error {
	result := c.ti.query.resultField
	if result == nil {
		return nil
	}

	for _, field := range result.data.fields {
		if err := c.checkColumnRead(field, true); err != nil {
			return err
		}
	}
	return nil
}

// checkQueryDataField checks the fields of the query's struct data type.
//
// CHECKLIST:
// - If the query represents a SELECT:
//  ✅ Each field of the query's struct data type MUST be of a type that
//     is readable from a value of the corresponding column.
// - If the query represents an INSERT or UPDATE:
//  ✅ Each field of the query's struct data type MUST be of a type that
//     is writeable to the corresponding column.
func (c *pgTypeCheck) checkQueryDataField() error {
	dataField := c.ti.query.dataField
	if dataField == nil || dataField.isDirective {
		return nil
	}

	if c.ti.query.kind == queryKindSelect {
		for _, field := range dataField.data.fields {
			if err := c.checkColumnRead(field, true); err != nil {
				return err
			}
		}
	} else if c.ti.query.kind == queryKindInsert || c.ti.query.kind == queryKindUpdate {
		for _, field := range dataField.data.fields {
			if err := c.checkColumnWrite(field, true); err != nil {
				return err
			}
		}
	}
	return nil
}

// tableJoin [ TODO documentation ]
type tableJoin struct {
	joinType joinType
	relId    relId
	tests    []*columnTest
}

// checkJoinBlock [ TODO documentation ]
func (c *pgTypeCheck) checkJoinBlock() error {
	jb := c.ti.query.joinStruct
	if jb == nil {
		return nil
	}

	for _, join := range jb.items {
		rel := c.joinrels[join.relId]

		tj := new(tableJoin)
		tj.joinType = join.joinType
		tj.relId = join.relId

		for _, item := range join.conds {
			switch cond := item.cond.(type) {
			case *searchConditionColumn:
				// A join condition's left-hand-side column MUST always
				// reference a column of the relation being joined, so to
				// avoid confusion make sure that node.colId has either no
				// qualifier or, if it has one, it matches the alias of
				// the joined table.
				if len(cond.lhsColId.qual) > 0 && cond.lhsColId.qual != join.relId.alias {
					return c.newError(errNoDatabaseRelation, cond.lhsColId, nil, cond)
				}

				test, err := c.checkColumnTest(cond, rel)
				if err != nil {
					return err
				}
				tj.tests = append(tj.tests, test)
			default:
				// NOTE(mkopriva): currently predicates other than
				// searchConditionColumn are not supported as a join condition.
			}
		}
		c.ti.joins = append(c.ti.joins, tj)
	}
	return nil
}

type onConflictType struct {
	constraint *pgconstraint
	index      *pgindex
	columns    []*pgcolumn
	update     []*pgcolumn
}

// checkOnConflictBlock
//
// CHECKLIST:
//  ✅ If a gosql.Column directive was used in the OnConflict block, the columns
//     listed in the directive's tag MUST be present in the target table, they
//     also MUST constitute a unique index of the target table.
//  ✅ If a gosql.Index directive was used in the OnConflict block, the index
//     specified in the directive's tag MUST be present on the target table
//     and it MUST be a unique index.
//  ✅ If a gosql.Constraint directive was used in the OnConflict block, the
//     constraint specified in the directive's tag MUST be present on the target
//     table and it MUST be a unique constraint.
//  ✅ If a gosql.Update directive was used in the OnConflict block, the columns
//     listed in the directive's tag MUST be present in the target table.
func (c *pgTypeCheck) checkOnConflictBlock() error {
	onconf := c.ti.query.onConflict
	if onconf == nil {
		return nil
	}

	conflict := new(onConflictType)

	// check the column list and ensure a unique index exists
	if onconf.column != nil {
		var attnums []int16
		for _, id := range onconf.column.sqlTagColIds {
			col := c.rel.findColumn(id.name)
			if col == nil {
				cid := colId{qual: id.qualifier, name: id.name} // TEMP
				return c.newError(errNoRelationColumn, cid, nil, onconf.column)
			}
			attnums = append(attnums, col.num)
			conflict.columns = append(conflict.columns, col)
		}

		var match bool
		for _, ind := range c.rel.indexes {
			if !ind.isunique && !ind.isprimary {
				continue
			}
			if !matchNumbers(ind.key, attnums) {
				continue
			}

			conflict.index = ind
			match = true
			break
		}
		if !match {
			return c.newError(errNoUniqueIndex, colId{}, nil, onconf.column)
		}
	}

	// check that the index exists and is unique.
	if index := onconf.index; index != nil {
		ind := c.rel.findIndex(index.sqlTagValue)
		if ind == nil || (!ind.isunique && !ind.isprimary) {
			return c.newError(errNoUniqueIndex, colId{}, nil, onconf.index)
		}
		conflict.index = ind
	}

	// check that the constraint exists and is unique.
	if constraint := onconf.constraint; constraint != nil {
		con := c.rel.findConstraint(constraint.sqlTagValue)
		if con == nil || (con.typ != pgconstraint_pkey && con.typ != pgconstraint_unique) {
			return c.newError(errNoUniqueConstraint, colId{}, nil, constraint)
		}
		conflict.constraint = con
	}

	// check that each specified column is present in the target table
	if onconf.update != nil {
		if onconf.update.sqlTagColIdList.all {
			// XXX this probably shouldn't be here, it's not doing any
			// type checking... maybe this could be handled before or
			// after the type checker has done its job...
			for _, write := range c.ti.writes {
				if c.ti.skipwrite2(write) {
					continue
				}
				conflict.update = append(conflict.update, write.column)
			}
		} else {
			for _, item := range onconf.update.sqlTagColIdList.items {
				col := c.rel.findColumn(item.name)
				if col == nil {
					return c.newError(errNoRelationColumn, colId{name: item.name}, nil, onconf.update)
				}
				conflict.update = append(conflict.update, col)
			}
		}
	}

	c.ti.conflict = conflict
	return nil
}

// columnRead holds the information needed by the generator to produce the
// expression nodes that constitute a column-to-field read operation.
type columnRead struct {
	// The column identifier.
	colId colId
	// The column from which the data will be read.
	column *pgcolumn
	// Info on the field into which the column will be read.
	field *fieldInfo
	// The name of the scanner to be used for reading the column.
	scanner string
}

// checkColumnRead checks if a value from the column that is associated with
// the given field can be read into that field. If strict=false and there
// is no column associated with the given field the check will be skipped.
// If the check is successful a new instance of *columnRead will be appended
// to the reads field of the *tagetInfo instance.
//
// CHECKLIST:
//  ✅ The field's type MUST NOT be a non-empty interface type.
//  ✅ The field's type MAY be a non-interface type that implements sql.Scanner.
//
// - If the column's type is json or jsonb:
//  ✅ The field's type MUST NOT be a chan, func, or a complex type IF it does
//     not implement the json.Unmarshaler interface,
//  ✅ otherwise the field's type MAY be any other type.
//
//  ✅ The field's type MAY be the empty-interface type.
//
// - If the column's type is xml:
//  ✅ The field's type MUST NOT be a func, chan, or map type IF it does not
//     implement the xml.Unmarshaler interface,
//  ✅ otherwise the field's type MAY be any other type.
//
//  ✅ The field's type MAY be any type that, together with the column's type,
//     has an entry in the relevant type table.
func (c *pgTypeCheck) checkColumnRead(field *fieldInfo, strict bool) error {
	col, err := c.findColumnByColId(field.colId)
	if err != nil && strict {
		return err
	} else if col == nil && !strict {
		return nil
	}

	check := func(fld *fieldInfo, col *pgcolumn) (scanner string, errcode typeErrorCode) {
		// non-empty interface, reject
		if fld.typ.kind == typeKindInterface && !fld.typ.isEmptyInterface {
			return "", errBadColumnReadIfaceType
		}

		// implements sql.Scanner & non-interface, accept as is
		if fld.typ.isScanner && fld.typ.kind != typeKindInterface {
			return "", 0
		}

		if col.typ.is(pgtyp_json, pgtyp_jsonb) {
			if !fld.typ.canJSONUnmarshal() {
				// chan, func, or complex, reject
				if fld.typ.is(typeKindChan, typeKindFunc, typeKindComplex64, typeKindComplex128) {
					return "", errBadColumnReadType
				}

				// []byte type, accept as is
				if fld.typ.goTypeId(false, false, true) == goTypeByteSlice {
					return "", 0
				}

				// string kind, accept as is
				if fld.typ.is(typeKindString) {
					return "", 0
				}
			}

			// everything else, accept with JSON
			return "JSON", 0
		}

		// empty interface, accept with AnyToEmptyInterface
		if fld.typ.isEmptyInterface {
			return "AnyToEmptyInterface", 0
		}

		if col.typ.is(pgtyp_xml) {
			if !fld.typ.canXMLUnmarshal() {
				if fld.typ.is(typeKindChan, typeKindFunc, typeKindMap) {
					return "", errBadColumnReadType
				}

				// []byte type, accept as is
				if fld.typ.goTypeId(false, false, true) == goTypeByteSlice {
					return "", 0
				}

				// string kind, accept as is
				if fld.typ.is(typeKindString) {
					return "", 0
				}
			}

			// everything else, accept with XML
			return "XML", 0
		}

		if c.isLength1Type(col) {
			// length-1 type table entry exists, accept
			gotyp := fld.typ.goTypeId(false, false, true)
			if e, ok := pgLength1TypeTable[col.typ.oid][gotyp]; ok {
				return e.scanner, 0
			}
			return "", errBadColumnReadType
		} else {
			// type table entry exists, accept
			gotyp := fld.typ.goTypeId(false, false, true)
			if e, ok := pgTypeTable[pgTypeKey{oid: col.typ.oid}][gotyp]; ok {
				return e.scanner, 0
			}

			// try to salvage this
			oid := col.typ.oid
			if col.typ.category == pgtypcategory_string {
				oid = pgtyp_text
			} else if col.typ.category == pgtypcategory_array {
				if et := c.pg.cat.types[col.typ.elem]; et != nil && et.category == pgtypcategory_string {
					oid = pgtyp_textarr
				}
			}
			if e, ok := pgTypeTable[pgTypeKey{oid: oid}][gotyp]; ok {
				return e.scanner, 0
			}
		}

		if false { // TODO(mkopriva): [ ... ]
			if col.typ.is(pgtyp_circle, pgtyp_circlearr) {
				// ...
			} else if col.typ.is(pgtyp_interval, pgtyp_intervalarr) {
				// ...
			} else if col.typ.typ == pgtyptype_domain {
				// ...
			} else if col.typ.typ == pgtyptype_composite {
				// ...
			}
		}

		return "", errBadColumnReadType
	}

	scanner, errcode := check(field, col)
	if errcode > 0 {
		return c.newError(errcode, field.colId, col, field)
	}

	read := new(columnRead)
	read.colId = field.colId
	read.column = col
	read.field = field
	read.scanner = scanner

	// NOTE(mkopriva): currently reading is allowed ONLY from the target
	// relation therefore here the alias of the target relation is used.
	// Once it is allowed to read from other, joined relations this will
	// need to be updated to properly handle that scenario.
	read.colId.qual = c.ti.dataField.relId.alias

	c.ti.reads = append(c.ti.reads, read)
	return nil
}

// columnWrite holds the information needed by the generator to produce the
// expression nodes that constitute a field-to-column write operation.
type columnWrite struct {
	// The column identifier.
	colId colId
	// The column to which the data will be written.
	column *pgcolumn
	// Info on the field from which the column will be written.
	field *fieldInfo
	// The name of the valuer to be used for writing the column.
	valuer string
}

// checkColumnWrite checks if a value of the given field can be written to the
// column that is associated with that field. If strict=false and there
// is no column associated with the given field the check will be skipped.
// If the check is successful a new instance of *columnWrite will be appended
// to the writes field of the *tagetInfo instance.
//
// CHECKLIST:
//  ✅ If strict=true; the column denoted by the given field MUST be present
//     in the target relation.
//  ✅ If the given field has the "default" tag, the target column MUST have
//     a default data value assigned.
//  ✅ The field's type MAY implement the driver.Valuer interface.
//
// - If the column's type is json or jsonb:
//  ✅ The field's type MUST NOT be a chan, func, or a complex type IF it does
//     not implement the json.Marshaler interface,
//  ✅ otherwise the field's type MAY be any other type.
//
// - If the column's type is xml:
//  ✅ The field's type MUST NOT be a func, chan, or map type IF it does not
//     implement the xml.Marshaler interface,
//  ✅ otherwise the field's type MAY be any other type.
//
//  ✅ The field's type MAY be any type that, together with the column's type,
//     has an entry in the relevant type table.
func (c *pgTypeCheck) checkColumnWrite(field *fieldInfo, strict bool) error {
	col := c.rel.findColumn(field.colId.name)
	if col == nil && strict {
		// colId with no qualifier to resolve to the target relation
		cid := colId{name: field.colId.name}
		return c.newError(errNoRelationColumn, cid, nil, field)
	} else if col == nil && !strict {
		return nil
	}

	check := func(fld *fieldInfo, col *pgcolumn) (valuer string, errcode typeErrorCode) {
		// default requested but non available, reject
		if fld.useDefault && !col.hasdefault {
			return "", errNoColumnDefault
		}

		// implements driver.Valuer, accept as is
		if fld.typ.isValuer {
			return "", 0
		}

		if col.typ.is(pgtyp_json, pgtyp_jsonb) {
			if !fld.typ.canJSONMarshal() {
				// chan, func, or complex, reject
				if fld.typ.is(typeKindChan, typeKindFunc, typeKindComplex64, typeKindComplex128) {
					return "", errBadColumnWriteType
				}

				// []byte type, accept as is
				if fld.typ.goTypeId(false, false, true) == goTypeByteSlice {
					return "", 0
				}

				// string kind, accept as is
				if fld.typ.is(typeKindString) {
					return "", 0
				}
			}

			// everything else, accept with JSON
			return "JSON", 0
		} else if col.typ.is(pgtyp_xml) {
			if !fld.typ.canXMLMarshal() {
				if fld.typ.is(typeKindChan, typeKindFunc, typeKindMap) {
					return "", errBadColumnWriteType
				}

				// []byte type, accept as is
				if fld.typ.goTypeId(false, false, true) == goTypeByteSlice {
					return "", 0
				}

				// string kind, accept as is
				if fld.typ.is(typeKindString) {
					return "", 0
				}
			}

			// everything else, accept with XML
			return "XML", 0
		}

		if c.isLength1Type(col) {
			// length-1 type table entry exists, accept
			gotyp := fld.typ.goTypeId(false, false, true)
			if e, ok := pgLength1TypeTable[col.typ.oid][gotyp]; ok {
				return e.valuer, 0
			}
			return "", errBadColumnWriteType
		} else {
			// type table entry exists, accept
			gotyp := fld.typ.goTypeId(false, false, true)
			if e, ok := pgTypeTable[pgTypeKey{oid: col.typ.oid}][gotyp]; ok {
				return e.valuer, 0
			}

			// try to salvage this
			oid := col.typ.oid
			if col.typ.category == pgtypcategory_string {
				oid = pgtyp_text
			} else if col.typ.category == pgtypcategory_array {
				if et := c.pg.cat.types[col.typ.elem]; et != nil && et.category == pgtypcategory_string {
					oid = pgtyp_textarr
				}
			}
			if e, ok := pgTypeTable[pgTypeKey{oid: oid}][gotyp]; ok {
				return e.valuer, 0
			}
		}

		if false { // TODO(mkopriva): [ ... ]
			if col.typ.is(pgtyp_circle, pgtyp_circlearr) {
				// ...
			} else if col.typ.is(pgtyp_interval, pgtyp_intervalarr) {
				// ...
			} else if col.typ.typ == pgtyptype_domain {
				// ...
			} else if col.typ.typ == pgtyptype_composite {
				// ...
			}
		}

		return "", errBadColumnWriteType
	}

	valuer, errcode := check(field, col)
	if errcode > 0 {
		return c.newError(errcode, field.colId, col, field)
	}

	write := new(columnWrite)
	write.colId = field.colId
	write.colId.qual = c.ti.dataField.relId.alias
	write.column = col
	write.field = field
	write.valuer = valuer

	c.ti.writes = append(c.ti.writes, write)
	if col.isprimary {
		c.ti.pkeys = append(c.ti.pkeys, write)
	}

	if true { // TODO TO BE REMOVED
		fc := &fieldColumnInfo{
			field:  write.field,
			column: write.column,
			colId:  write.colId,
			pgsql:  pgTypeEntry{valuer: write.valuer},
		}

		c.ti.input = append(c.ti.input, fc)
		if col.isprimary {
			c.ti.primaryKeys = append(c.ti.primaryKeys, fc)
		}
	}
	return nil
}

// columnFilter [ TODO documentation ]
type columnFilter struct {
	// The column identifier.
	colId colId
	// The column which to filter by.
	column *pgcolumn
	// Info on the field that holds the value to be used as the filter parameter.
	field *fieldInfo
	// The name of the valuer to be used for converting the field value.
	valuer string
}

// checkColumnFilter [ TODO documentation }
func (c *pgTypeCheck) checkColumnFilter(field *fieldInfo, strict bool) error {
	col := c.rel.findColumn(field.colId.name)
	if col == nil && strict {
		// colId with no qualifier to resolve to the target relation
		cid := colId{name: field.colId.name}
		return c.newError(errNoRelationColumn, cid, nil, field)
	} else if col == nil && !strict {
		return nil
	}

	check := func(fld *fieldInfo, col *pgcolumn) (valuer string, errcode typeErrorCode) {
		// implements driver.Valuer, accept as is
		if fld.typ.isValuer {
			return "", 0
		}

		if col.typ.is(pgtyp_json, pgtyp_jsonb) {
			if !fld.typ.canJSONMarshal() {
				// chan, func, or complex, reject
				if fld.typ.is(typeKindChan, typeKindFunc, typeKindComplex64, typeKindComplex128) {
					return "", errBadColumnWriteType
				}

				// []byte type, accept as is
				if fld.typ.goTypeId(false, false, true) == goTypeByteSlice {
					return "", 0
				}

				// string kind, accept as is
				if fld.typ.is(typeKindString) {
					return "", 0
				}
			}

			// everything else, accept with JSON
			return "JSON", 0
		} else if col.typ.is(pgtyp_xml) {
			if !fld.typ.canXMLMarshal() {
				if fld.typ.is(typeKindChan, typeKindFunc, typeKindMap) {
					return "", errBadColumnWriteType
				}

				// []byte type, accept as is
				if fld.typ.goTypeId(false, false, true) == goTypeByteSlice {
					return "", 0
				}

				// string kind, accept as is
				if fld.typ.is(typeKindString) {
					return "", 0
				}
			}

			// everything else, accept with XML
			return "XML", 0
		}

		if c.isLength1Type(col) {
			// length-1 type table entry exists, accept
			gotyp := fld.typ.goTypeId(false, false, true)
			if e, ok := pgLength1TypeTable[col.typ.oid][gotyp]; ok {
				return e.valuer, 0
			}
			return "", errBadColumnWriteType
		} else {
			// type table entry exists, accept
			gotyp := fld.typ.goTypeId(false, false, true)
			if e, ok := pgTypeTable[pgTypeKey{oid: col.typ.oid}][gotyp]; ok {
				return e.valuer, 0
			}

			// try to salvage this
			oid := col.typ.oid
			if col.typ.category == pgtypcategory_string {
				oid = pgtyp_text
			} else if col.typ.category == pgtypcategory_array {
				if et := c.pg.cat.types[col.typ.elem]; et != nil && et.category == pgtypcategory_string {
					oid = pgtyp_textarr
				}
			}
			if e, ok := pgTypeTable[pgTypeKey{oid: oid}][gotyp]; ok {
				return e.valuer, 0
			}
		}

		if false { // TODO(mkopriva): [ ... ]
			if col.typ.is(pgtyp_circle, pgtyp_circlearr) {
				// ...
			} else if col.typ.is(pgtyp_interval, pgtyp_intervalarr) {
				// ...
			} else if col.typ.typ == pgtyptype_domain {
				// ...
			} else if col.typ.typ == pgtyptype_composite {
				// ...
			}
		}

		return "", errBadColumnWriteType
	}

	valuer, errcode := check(field, col)
	if errcode > 0 {
		return c.newError(errcode, field.colId, col, field)
	}

	filter := new(columnFilter)
	filter.colId = field.colId
	filter.colId.qual = c.ti.dataField.relId.alias
	filter.column = col
	filter.field = field
	filter.valuer = valuer

	c.ti.filters = append(c.ti.filters, filter)
	return nil
}

// columnTest holds the information needed by the generator
// to produce an SQL predicate boolean expression.
type columnTest struct {
	// Left hand side column id.
	lhsColId colId
	// Left hand side column.
	lhsColumn *pgcolumn
	// Right hand side column id, or empty.
	rhsColId colId
	// Right hand side column, or nil.
	rhsColumn *pgcolumn
	// Right hand side literal expression, or empty.
	rhsLiteral string
	// Type of the right hand side column or literal expression.
	rhsType *pgtype
	// The type of the predicate, or 0.
	pred predicate
	// The predicate quantifier, or 0.
	qua quantifier
}

// checkColumnTest checks if the information in the given searchConditionColumn
// can be used to produce a valid, type-safe boolean expression. If the check is
// successful a new instance of *columnTest holding information relevant to the
// generator will be returned.
//
// CHECKLIST:
// - If the *pgRelationInfo argument is not nil:
//  ✅ The LHS column MUST be present in given relation.
//
// - If the *pgRelationInfo argument is nil:
//  ✅ The LHS column MUST be present in one of the loaded relations.
//
// - If the predicate of the condition argument is unary:
//  ✅ If the unary predicate is one of the "IS [NOT] { FALSE | TRUE | UNKNOWN }"
//     predicates then the type of the LHS column/expression MUST be boolean.
//  ✅ If the unary predicate is one of the "IS [NOT] NULL" predicates then
//     the LHS column MUST NOT have the "NOT NULL" constraint.
//
// - If a predicate quantifier was provided:
//  ✅ The RHS column or literal expression MUST be quantifiable.
//
//  ✅ The LHS and RHS types MUST be comparable with the given predicate and quantifier.
func (c *pgTypeCheck) checkColumnTest(cond *searchConditionColumn, rel *pgRelationInfo) (*columnTest, error) {
	var lhsColumn *pgcolumn

	if rel != nil {
		if lhsColumn = rel.findColumn(cond.lhsColId.name); lhsColumn == nil {
			return nil, c.newError(errNoRelationColumn, cond.lhsColId, nil, cond)
		}
	} else {
		var errcode typeErrorCode
		if lhsColumn, errcode = c.findColumn(cond.lhsColId); errcode > 0 {
			return nil, c.newError(errcode, cond.lhsColId, nil, cond)
		}
	}
	test := new(columnTest)
	test.lhsColId = cond.lhsColId
	test.lhsColumn = lhsColumn

	if cond.pred.isUnary() {
		if cond.pred.isBoolean() && lhsColumn.typ.oid != pgtyp_bool {
			return nil, c.newPredicateError(errBadUnaryPredicateType, lhsColumn, nil, cond)
		}
		if lhsColumn.hasnotnull && (cond.pred == isNull || cond.pred == notNull) {
			return nil, c.newPredicateError(errBadUnaryPredicateType, lhsColumn, nil, cond)
		}
	} else {
		var rhsType *pgtype     // the type of the RHS column/literal
		var rhsColumn *pgcolumn // for error reporting

		if len(cond.rhsColId.name) > 0 {
			col, err := c.findColumnByColId(cond.rhsColId)
			if err != nil {
				return nil, err
			}
			rhsType = col.typ
			rhsColumn = col
		} else if len(cond.literal) > 0 {
			lityp, err := c.typeOfLiteral(cond.literal)
			if err != nil {
				return nil, err
			}
			rhsType = lityp
		}

		if cond.pred.isQuantified() || cond.qua > 0 {
			if rhsType.category != pgtypcategory_array {
				return nil, c.newPredicateError(errBadPredicateQuantifierType, lhsColumn, rhsColumn, cond)
			}
		}

		if !c.checkComparison(lhsColumn.typ, rhsType, cond.pred, cond.qua) {
			return nil, c.newPredicateError(errBadComparisonOperation, lhsColumn, rhsColumn, cond)
		}

		test.rhsColId = cond.rhsColId
		test.rhsColumn = rhsColumn
		test.rhsLiteral = cond.literal
		test.rhsType = rhsType
		test.pred = cond.pred
		test.qua = cond.qua
	}

	return test, nil
}

type whereInfo struct {
	conds []interface{}
}

// checkWhereBlock
func (c *pgTypeCheck) checkWhereBlock() error {
	where := c.ti.query.whereStruct
	if where == nil {
		return nil
	}

	c.ti.where = new(whereInfo)

	type loopstate struct {
		conds []*searchCondition // the current iteration search conditions
		idx   int                // keeps track of the field index
	}
	stack := []*loopstate{{conds: where.conds}} // LIFO stack of states.

stackloop:
	for len(stack) > 0 {
		// Loop over the various conds of a whereblock, including
		// other nested whereblocks and check them against the db.
		loop := stack[len(stack)-1]
		for loop.idx < len(loop.conds) {
			item := loop.conds[loop.idx]

			// Instead of incrementing the index in the for-statement
			// it is done here manually to ensure that it is not skipped
			// when continuing to the outer loop.
			loop.idx++

			switch cond := item.cond.(type) {
			case *searchConditionField:
				if err := c.checkFieldTest(cond); err != nil {
					return err
				}
			case *searchConditionColumn:
				test, err := c.checkColumnTest(cond, nil)
				if err != nil {
					return err
				}
				c.ti.where.conds = append(c.ti.where.conds, test)
			case *searchConditionBetween:
				if err := c.checkBetweenTest(cond); err != nil {
					return err
				}
			case *searchConditionNested:
				stack = append(stack, &loopstate{conds: cond.conds})
				continue stackloop
			}
		}

		stack = stack[:len(stack)-1]
	}

	return nil
}

func (c *pgTypeCheck) checkFieldTest(cond *searchConditionField) error {
	// CHECKLIST: [column, operator, field]
	// ☑️  1. The column must exist in one of the associated relations.
	// ☑️  2. A NOT NULL column requires non-pointer field.
	//	*** Not sure this is necessary, how can this help? Are there cases
	//          where this is more of a burden than help?
	// ☑️  3. column must be, using the given operator, comparable to one of the
	//    types to which the field type can be converted.
	// ☑️  4. If a quantifier (ANY, ALL, etc.) was provided the field must be a slice/array.
	// ☑️  5. If a quantifier (ANY, ALL, etc.) was provided the field's element type must be
	//       used in check (3).
	// ☑️  6. If a modifier function (lower, upper, etc.) was provided the column's type must
	//       match the function's argument type.
	// ☑️  7. If a modifier function (lower, upper, etc.) was provided the field's type must
	//       must be convertible to a type accepted by the function as its argument.
	// ☑️  8. If a modifier function (lower, upper, etc.) was provided the result types of
	//       the LHS and RHS expressions must be comparable using the given operator.

	// TODO(mkopriva): COLUMN TO FIELD COMPARISON:
	// 0. construct a map of: [LHS (postgres_oid) -> CMP (comparison_operator) -> RHS (postgres_oid[])]
	// 1. resolve the RHS list using the known LHS and CMP
	// 2. using the RHS list, the typetable, and the field's type
	//    we can look for a compatible match, but, the first
	//    oid in the RHS list should be, if present, the LHS oid
	//    so that that is chosen if the field can be converted to
	//    that.

	// Check that the referenced Column is present
	// in one of the associated relations.
	col, err := c.findColumnByColId(cond.colId)
	if err != nil {
		return err
	}
	// If the column cannot be set to NULL, then make
	// sure that the field's not a pointer.
	if col.hasnotnull {
		if cond.typ.kind == typeKindPtr {
			return errors.IllegalPtrFieldForNotNullColumnError
		}
	}

	// list of types to which the field type can potentially be converted
	var fieldoids = c.typeoids(cond.typ)

	// If this is a quantified predicate then check that
	// the field is a slice or array, and also make sure that
	// the column's type can be compared to the element type
	// of the slice / array.
	if cond.qua > 0 || cond.pred.isQuantified() {
		if cond.typ.kind != typeKindSlice && cond.typ.kind != typeKindArray {
			return errors.IllegalFieldTypeForQuantifierError
		}
		fieldoids = c.typeoids(*cond.typ.elem)
	}

	if len(cond.modFunc) > 0 {
		// Check that the modifier function can
		// be used with the given Column's type.
		if err := c.checkModifierFunction(cond.modFunc, col, fieldoids); err != nil {
			return err
		}
	} else {
		// Check that the Field's type can be
		// compared to that of the Column.
		if !c.canCompare(col, fieldoids, cond.pred) {
			return errors.BadFieldToColumnTypeError
		}
	}

	c.ti.searchConditionFieldColumns[cond] = col

	return nil
}

// Check that the referenced Column is present
// in one of the associated relations.
func (c *pgTypeCheck) checkBetweenTest(cond *searchConditionBetween) error {
	col, errcode := c.findColumn(cond.colId)
	if errcode > 0 {
		return c.newError(errcode, cond.colId, nil, cond)
	}

	// Check that both predicands, x and y, can be compared to the column.
	for _, arg := range []interface{}{cond.x, cond.y} {
		var argoids []pgOID
		switch a := arg.(type) {
		case *colId:
			col2, err := c.findColumnByColId(*a)
			if err != nil {
				return err
			}
			argoids = []pgOID{col2.typ.oid}
		case *fieldDatum:
			argoids = c.typeoids(a.typ)
		}

		if !c.canCompare(col, argoids, isGT) {
			return errors.BadColumnToColumnTypeComparisonError
		}
	}
	return nil
}

// If a TextSearch directive was provided, make sure that the
// specified column is present in one of the loaded relations
// and that it has the correct type.
func (c *pgTypeCheck) checkTextSearchColId() error {
	colid := c.ti.filter.textSearchColId
	if colid == nil {
		return nil
	}

	col, err := c.findColumnByColId(*colid)
	if err != nil {
		return err
	} else if col.typ.oid != pgtyp_tsvector {
		return errors.BadDBColumnTypeError
	}

	return nil
}

func (c *pgTypeCheck) checkFilterDataField() error {
	dataField := c.ti.filter.dataField
	if dataField == nil {
		return nil
	}

	for _, field := range dataField.data.fields {
		if err := c.checkColumnFilter(field, false); err != nil {
			return err
		}
	}
	return nil
}

// checkComparison checks whether or not a valid comparison operation
// expression can be generated from the provided arguments.
//
// CHECKLIST:
//  ✅ If the predicate together with the quantifier constitute an array comparison
//     the check MUST use the RHS's element type instead of the RHS type.
//  ✅ ACCEPT if the LHS type belongs to the string category and the RHS type is unknown.
//  ✅ ACCEPT if the combination of LHS type, RHS type, and the predicate has
//     an entry in the pg_operator table.
func (c *pgTypeCheck) checkComparison(ltyp, rtyp *pgtype, pred predicate, qua quantifier) bool {
	if pred.isQuantified() || qua > 0 {
		rtyp = c.pg.cat.types[rtyp.elem]
	}
	if ltyp.category == pgtypcategory_string && rtyp.oid == pgtyp_unknown {
		return true
	}

	key := pgopkey{left: ltyp.oid, right: rtyp.oid}
	key.name = predicateToBasePGOps[pred]
	if _, ok := c.pg.cat.operators[key]; ok {
		return true
	}
	return false
}

// typeoids returns a list of OIDs of those PostgreSQL types that can be
// used to hold a value of a Go type represented by the given typeInfo.
func (c *pgTypeCheck) typeoids(typ typeInfo) []pgOID {
	switch typstr := typ.goTypeId(false, false, true); typstr {
	case goTypeStringMap, goTypeNullStringMap, goTypeStringPtrMap:
		if t := c.pg.cat.typeByName("hstore"); t != nil {
			return []pgOID{t.oid}
		}
	case goTypeStringMapSlice, goTypeNullStringMapSlice, goTypeStringPtrMapSlice:
		if t := c.pg.cat.typeByName("_hstore"); t != nil {
			return []pgOID{t.oid}
		}
	default:
		if oids, ok := go2pgoids[typstr]; ok {
			return oids
		}
	}
	return nil
}

// canCompare reports whether a value of the given col's type can be compared to,
// using the predicate, a value of one of the types specified by the given oids.
func (c *pgTypeCheck) canCompare(col *pgcolumn, rhstypes []pgOID, pred predicate) bool {
	name := predicateToBasePGOps[pred]
	left := col.typ.oid
	for _, right := range rhstypes {
		if col.typ.category == pgtypcategory_string && right == pgtyp_unknown {
			return true
		}

		key := pgopkey{name: name, left: left, right: right}
		if _, ok := c.pg.cat.operators[key]; ok {
			return true
		}
	}
	return false
}

func (c *pgTypeCheck) canCoerceOID(targetid, inputid pgOID) bool {
	// no problem if same type
	if targetid == inputid {
		return true
	}

	// accept if target is ANY
	if targetid == pgtyp_any {
		return true
	}

	// if input is an untyped string constant assume it can be converted to anything
	if inputid == pgtyp_unknown {
		return true
	}

	// if pg_cast says ok, accept.
	if c.pg.cat.cancasti(targetid, inputid) {
		return true
	}
	return false
}

// check that the column's and field's type can be used as the argument to the specified function.
func (c *pgTypeCheck) checkModifierFunction(fn funcName, col *pgcolumn, fieldoids []pgOID) error {
	var (
		ok1 bool // column's type can be coerced to the functions argument's type
		ok2 bool // field's type can be assigned to the functions argument's type
	)

	if proclist, ok := c.pg.cat.procs[fn]; ok {
		for _, proc := range proclist {
			// ok if the column's type can be coerced to the function argument's type
			if c.canCoerceOID(proc.argtype, col.typ.oid) {
				ok1 = true
			}

			// ok if one of the fieldoids types can be assigned to the function argument's type
			for _, foid := range fieldoids {
				if c.canCoerceOID(proc.argtype, foid) {
					ok2 = true
				}
			}
		}
	}
	if !ok1 || !ok2 {
		return errors.BadColumnTypeForDBFuncError
	}
	return nil
}

func (c *pgTypeCheck) loadRelation(id relId) (*pgRelationInfo, error) {
	rel := new(pgRelationInfo)
	rel.name = id.name
	rel.namespace = id.qual
	if len(rel.namespace) == 0 {
		rel.namespace = "public"
	}

	const selectRelationInfo = `SELECT
		c.oid
		, c.relkind
	FROM pg_class c
	WHERE c.relname = $1
	AND c.relnamespace = to_regnamespace($2)` //`
	row := c.pg.db.QueryRow(selectRelationInfo, rel.name, rel.namespace)
	if err := row.Scan(&rel.oid, &rel.relkind); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NoDBRelationError
		}
		return nil, err
	}

	const selectRelationColumns = `SELECT
		a.attnum
		, a.attname
		, a.atttypmod
		, a.attndims
		, a.attnotnull
		, a.atthasdef
		, COALESCE(i.indisprimary, false)
		, a.atttypid
	FROM pg_attribute a
	LEFT JOIN pg_index i ON (
		i.indisprimary
		AND i.indrelid = a.attrelid
		AND a.attnum = ANY(i.indkey)
	)
	WHERE a.attrelid = $1
	AND a.attnum > 0
	AND NOT a.attisdropped
	ORDER BY a.attnum` //`
	rows, err := c.pg.db.Query(selectRelationColumns, rel.oid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		col := new(pgcolumn)
		err := rows.Scan(
			&col.num,
			&col.name,
			&col.typmod,
			&col.ndims,
			&col.hasnotnull,
			&col.hasdefault,
			&col.isprimary,
			&col.typoid,
		)
		if err != nil {
			return nil, err
		}
		typ, ok := c.pg.cat.types[col.typoid]
		if !ok {
			return nil, errors.UnknownPostgresTypeError
		}

		col.typ = typ
		rel.columns = append(rel.columns, col)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	const selectRelationConstraints = `SELECT
		c.conname
		, c.contype
		, c.condeferrable
		, c.condeferred
		, c.conkey
		, c.confkey
	FROM pg_constraint c
	LEFT JOIN pg_index i ON i.indexrelid = c.conindid
	WHERE c.conrelid = $1
	ORDER BY c.oid` //`
	rows, err = c.pg.db.Query(selectRelationConstraints, rel.oid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		con := new(pgconstraint)
		err := rows.Scan(
			&con.name,
			&con.typ,
			&con.isdeferrable,
			&con.isdeferred,
			(*pq.Int64Array)(&con.key),
			(*pq.Int64Array)(&con.fkey),
		)
		if err != nil {
			return nil, err
		}
		rel.constraints = append(rel.constraints, con)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	const selectRelationIndexes = `SELECT
		c.relname
		, i.indnatts
		, i.indisunique
		, i.indisprimary
		, i.indisexclusion
		, i.indimmediate
		, i.indisready
		, i.indkey
		, pg_catalog.pg_get_indexdef(i.indexrelid)
		, COALESCE(pg_catalog.pg_get_expr(i.indpred, i.indrelid, true), '')
	FROM pg_index i
	LEFT JOIN pg_class c ON c.oid = i.indexrelid
	WHERE i.indrelid = $1
	ORDER BY i.indexrelid` //`
	rows, err = c.pg.db.Query(selectRelationIndexes, rel.oid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		ind := new(pgindex)
		err := rows.Scan(
			&ind.name,
			&ind.natts,
			&ind.isunique,
			&ind.isprimary,
			&ind.isexclusion,
			&ind.isimmediate,
			&ind.isready,
			(*int2vec)(&ind.key),
			&ind.indexdef,
			&ind.indpred,
		)
		if err != nil {
			return nil, err
		}

		ind.indexpr = parseIndexExpr(ind.indexdef)
		rel.indexes = append(rel.indexes, ind)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Map the relation to its alias or name if no alias was specified.
	// This will help with matching columns, constraints, and indexes with
	// the relation with which their identifiers were qualified.
	if len(id.alias) > 0 {
		c.relmap[id.alias] = rel
	} else {
		c.relmap[id.name] = rel
	}
	return rel, nil
}

// DEPRECATED in favor of findColumnByColId2
func (c *pgTypeCheck) findColumnByColId(id colId) (*pgcolumn, error) {
	rel, ok := c.relmap[id.qual]
	if !ok {
		return nil, errors.NoDBRelationError
	}
	col := rel.findColumn(id.name)
	if col == nil {
		return nil, errors.NoDBColumnError
	}
	return col, nil
}

// findColumnField finds and returns the *fieldInfo of the target data type's
// field that is tagged with the column identified by the given colId. If strict
// is true, findColumnField will also check whether or not the column actually exists.
//
// NOTE(mkopriva): currently the findColumnField method returns fields matched
// by just the column's name, i.e. the qualifiers are ignored, this means that
// one could pass in two different colIds with the same name and the method
// would return the same field.
func (c *pgTypeCheck) findColumnField(id colId, strict bool) (*fieldInfo, typeErrorCode) {
	for _, field := range c.ti.query.dataField.data.fields {
		if field.colId.name == id.name {
			if strict {
				if errcode := c.checkColumnExists(id); errcode > 0 {
					return nil, errcode
				}
			}
			return field, 0
		}
	}
	return nil, errNoColumnField
}

// checkColumnExists checks whether or not the column, or its relation, denoted by
// the given id is present in the database. If the column exists 0 will be returned,
// if the column doesn't exist errNoRelationColumn will be returned, and if the
// column's relation doesn't exist then errNoDatabaseRelation will be returned.
func (c *pgTypeCheck) checkColumnExists(id colId) typeErrorCode {
	if rel, ok := c.relmap[id.qual]; ok {
		for _, col := range rel.columns {
			if col.name == id.name {
				return 0 // found
			}
		}
		return errNoRelationColumn
	}
	return errNoDatabaseRelation
}

// findColumn finds and returns the *pgcolumn identified by the given colId.
// If the relation denoted by the column's qualifier doesn't exist, or if the
// relation exists but the column itself is not present in that relation an
// error code will be returned.
func (c *pgTypeCheck) findColumn(id colId) (*pgcolumn, typeErrorCode) {
	if rel, ok := c.relmap[id.qual]; ok {
		for _, col := range rel.columns {
			if col.name == id.name {
				return col, 0 // found
			}
		}
		return nil, errNoRelationColumn
	}
	return nil, errNoDatabaseRelation
}

// newError constructs and returns a new typeError value.
func (c *pgTypeCheck) newError(errcode typeErrorCode, cid colId, col *pgcolumn, fptr fieldPtr) error {
	e := typeError{errorCode: errcode}
	e.pkgPath = c.ti.pkgPath
	e.targetName = c.ti.typeName
	e.dbName = c.pg.name
	e.colQualifier = cid.qual
	e.colName = cid.name

	if rel, ok := c.relmap[cid.qual]; ok {
		e.relQualifier = rel.namespace
		e.relName = rel.name
	}
	if f, ok := c.ti.fieldmap[fptr]; ok {
		p := c.fset.Position(f.fvar.Pos())
		e.fieldType = f.fvar.Type().String()
		e.fieldName = f.fvar.Name()
		e.fileName = p.Filename
		e.fileLine = p.Line
	}
	if col != nil {
		e.colType = col.typ.namefmt
	}
	return e
}

// newPredicateError constructs and returns a new predicateError value.
func (c *pgTypeCheck) newPredicateError(errcode predicateErrorCode, lhsColumn *pgcolumn, rhsColumn *pgcolumn, fptr fieldPtr) error {
	e := predicateError{errorCode: errcode}
	e.pkgPath = c.ti.pkgPath
	e.targetName = c.ti.typeName
	e.dbName = c.pg.name

	if f, ok := c.ti.fieldmap[fptr]; ok {
		p := c.fset.Position(f.fvar.Pos())
		e.fieldType = f.fvar.Type().String()
		e.fieldName = f.fvar.Name()
		e.fileName = p.Filename
		e.fileLine = p.Line
	}

	switch cond := fptr.(type) {
	case *searchConditionColumn:
		e.lhsExpr = cond.lhsColId.quoted()
		if rel, ok := c.relmap[cond.lhsColId.qual]; ok {
			e.lhsRel = fmt.Sprintf("%q.%q", rel.namespace, rel.name)
		}
		if lhsColumn != nil {
			e.lhsType = lhsColumn.typ.namefmt
		}

		if !cond.rhsColId.isEmpty() {
			e.rhsExpr = cond.rhsColId.quoted()
			if rel, ok := c.relmap[cond.rhsColId.qual]; ok {
				e.rhsRel = fmt.Sprintf("%q.%q", rel.namespace, rel.name)
			}
			if rhsColumn != nil {
				e.rhsType = rhsColumn.typ.namefmt
			}
		} else if len(cond.literal) > 0 {
			e.rhsExpr = cond.literal
			if typ, err := c.typeOfLiteral(cond.literal); err == nil {
				e.rhsType = typ.namefmt
			}
		}

		// TODO
		e.predicate = ""
		e.quantifier = ""
	}
	return e
}

// isLength1Type reports whether or not the given column's type
// is a "length 1" type, i.e. char(1), varchar(1), or bit(1)[], etc.
func (c *pgTypeCheck) isLength1Type(col *pgcolumn) bool {
	typ := col.typ
	if typ.category == pgtypcategory_array {
		typ = c.pg.cat.types[typ.elem]
	}

	if typ.category == pgtypcategory_bitstring {
		return (col.typmod == 1)
	} else if typ.category == pgtypcategory_string {
		return ((col.typmod - 4) == 1)
	}
	return false
}

// isNumericWithoutScale reports whether or not the given column's type is
// a numeric type with no scale. Numeric types without scale are stored by
// postgres as integers which might be useful in deciding how to decode the
// column's value in Go.
func (c *pgTypeCheck) isNumericWithoutScale(col *pgcolumn) bool {
	if col.typ.is(pgtyp_numeric, pgtyp_numericarr) {
		precision := ((col.typmod - 4) >> 16) & 65535
		scale := (col.typmod - 4) & 65535
		return (precision > 0 && scale == 0)
	}
	return false
}

// typeOfLiteral returns the type of the given literal expression.
func (c *pgTypeCheck) typeOfLiteral(expr string) (*pgtype, error) {
	const pgselectexprtype = `SELECT id::oid FROM pg_typeof(%s) AS id` //`

	var oid pgOID
	row := c.pg.db.QueryRow(fmt.Sprintf(pgselectexprtype, expr))
	if err := row.Scan(&oid); err != nil {
		return nil, errors.BadLiteralExpressionError
	}
	return c.pg.cat.types[oid], nil
}

type pgRelationInfo struct {
	oid         pgOID  // The object identifier of the relation.
	name        string // The name of the relation.
	namespace   string // The name of the namespace to which the relation belongs.
	relkind     string // The relation's kind, we're only interested in r, v, and m.
	columns     []*pgcolumn
	constraints []*pgconstraint
	indexes     []*pgindex
}

func (rel *pgRelationInfo) findColumn(name string) *pgcolumn {
	for _, col := range rel.columns {
		if col.name == name {
			return col
		}
	}
	return nil
}

func (rel *pgRelationInfo) findConstraint(name string) *pgconstraint {
	for _, con := range rel.constraints {
		if con.name == name {
			return con
		}
	}
	return nil
}

func (rel *pgRelationInfo) findIndex(name string) *pgindex {
	for _, ind := range rel.indexes {
		if ind.name == name {
			return ind
		}
	}
	return nil
}

type pgcolumn struct {
	num  int16  // The number of the column. Ordinary columns are numbered from 1 up.
	name string // The name of the member's column.
	// The number of dimensions if the column is an array type, otherwise 0.
	ndims int
	// Records type-specific data supplied at table creation time (for example,
	// the maximum length of a varchar column). It is passed to type-specific
	// input functions and length coercion functions. The value will generally
	// be -1 for types that do not need.
	//
	// NOTE(mkopriva): to get the actual value subtract 4.
	// NOTE(mkopriva): in the case of NUMERIC(precision, scale) types, to
	// calculate the precision use ((typmod - 4) >> 16) & 65535 and to
	// calculate the scale use (typmod - 4) && 65535
	typmod int
	// Indicates whether or not the column has a NOT NULL constraint.
	hasnotnull bool
	// Indicates whether or not the column has a DEFAULT value.
	hasdefault bool
	// Reports whether or not the column is a primary key.
	isprimary bool
	// The id of the column's type.
	typoid pgOID
	// Info about the column's type.
	typ *pgtype
}

type pgconstraint struct {
	// Constraint name (not necessarily unique!)
	name string
	// The type of the constraint
	typ string
	// Indicates whether or not the constraint is deferrable
	isdeferrable bool
	// Indicates whether or not the constraint is deferred by default
	isdeferred bool
	// If a table constraint (including foreign keys, but not constraint triggers),
	// list of the constrained columns
	key []int64
	// If a foreign key, list of the referenced columns
	fkey []int64
}

type pgindex struct {
	// The name of the index.
	name string
	// The total number of columns in the index; this number includes
	// both key and included attributes.
	natts int
	// If true, this is a unique index.
	isunique bool
	// If true, this index represents the primary key of the table.
	isprimary bool
	// If true, this index supports an exclusion constraint.
	isexclusion bool
	// If true, the uniqueness check is enforced immediately on insertion.
	isimmediate bool
	// If true, the index is currently ready for inserts. False means the
	// index must be ignored by INSERT/UPDATE operations.
	isready bool
	// This is an array of values that indicate which table columns this index
	// indexes. For example a value of 1 3 would mean that the first
	// and the third table columns make up the index entries. Key columns come
	// before non-key (included) columns. A zero in this array indicates that
	// the corresponding index attribute is an expression over the table columns,
	// rather than a simple column reference.
	key []int16
	// The index definition.
	indexdef string
	indpred  string
	indexpr  string
}

type pgtype struct {
	oid pgOID
	// The name of the type.
	name string
	// The formatted name of the type.
	namefmt string
	// The number of bytes for fixed-size types, negative for variable length types.
	length int
	// The type's type.
	typ string
	// An arbitrary classification of data types that is used by the parser
	// to determine which implicit casts should be "preferred".
	category string
	// True if the type is a preferred cast target within its category.
	ispreferred bool
	// If this is an array type then elem identifies the element type
	// of that array type.
	elem pgOID
}

func (t pgtype) is(oids ...pgOID) bool {
	for _, oid := range oids {
		if t.oid == oid {
			return true
		}
	}
	return false
}

// isbase returns true if t's oid matches one of the given oids, or if t is an
// array type isbase returns true if t's elem matches one of the given oids.
func (t pgtype) isbase(oids ...pgOID) bool {
	if t.category == pgtypcategory_array {
		for _, oid := range oids {
			if t.elem == oid {
				return true
			}
		}
	}
	return t.is(oids...)
}

type pgoperator struct {
	oid    pgOID
	name   string
	kind   string
	left   pgOID
	right  pgOID
	result pgOID
}

type pgopkey struct {
	name  string
	left  pgOID
	right pgOID
}

type pgcast struct {
	oid     pgOID
	source  pgOID
	target  pgOID
	context string
}

type pgproc struct {
	oid     pgOID
	name    string
	argtype pgOID
	rettype pgOID
	isagg   bool
}

type pgcastkey struct {
	target pgOID
	source pgOID
}

// pgCatalog holds useful information on various objects of the database.
type pgCatalog struct {
	types     map[pgOID]*pgtype
	operators map[pgopkey]*pgoperator
	casts     map[pgcastkey]*pgcast
	procs     map[funcName][]*pgproc
}

// typeByName looks up and returns the pgtype by the given name.
func (c *pgCatalog) typeByName(name string) *pgtype {
	for _, t := range c.types {
		if t.name == name {
			return t
		}
	}
	return nil
}

func (c *pgCatalog) typebyoid(oid pgOID) *pgtype {
	return c.types[oid]
}

// cancasti reports whether s can be cast to t *implicitly* or in assignment.
func (c *pgCatalog) cancasti(t, s pgOID) bool {
	key := pgcastkey{target: t, source: s}
	if cast := c.casts[key]; cast != nil {
		return cast.context == pgcast_implicit || cast.context == pgcast_assignment
	}
	return false
}

func (c *pgCatalog) load(db *sql.DB, key string) error {
	pgCatalogCache.Lock()
	defer pgCatalogCache.Unlock()

	cat := pgCatalogCache.m[key]
	if cat != nil {
		*c = *cat
		return nil
	}

	const selectTypes = `SELECT
		t.oid
		, t.typname
		, pg_catalog.format_type(t.oid, NULL)
		, t.typlen
		, t.typtype
		, t.typcategory
		, t.typispreferred
		, t.typelem
	FROM pg_type t
	WHERE t.typrelid = 0
	AND pg_catalog.pg_type_is_visible(t.oid)
	AND t.typcategory <> 'P'` //`
	c.types = make(map[pgOID]*pgtype)
	rows, err := db.Query(selectTypes)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		typ := new(pgtype)
		err := rows.Scan(
			&typ.oid,
			&typ.name,
			&typ.namefmt,
			&typ.length,
			&typ.typ,
			&typ.category,
			&typ.ispreferred,
			&typ.elem,
		)
		if err != nil {
			return err
		}
		c.types[typ.oid] = typ
	}
	if err := rows.Err(); err != nil {
		return err
	}

	const selectOperators = `SELECT
		o.oid
		, o.oprname
		, o.oprkind
		, o.oprleft
		, o.oprright
		, o.oprresult
	FROM pg_operator o ` //`
	c.operators = make(map[pgopkey]*pgoperator)
	rows, err = db.Query(selectOperators)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		op := new(pgoperator)
		err := rows.Scan(
			&op.oid,
			&op.name,
			&op.kind,
			&op.left,
			&op.right,
			&op.result,
		)
		if err != nil {
			return err
		}
		c.operators[pgopkey{name: op.name, left: op.left, right: op.right}] = op
	}
	if err := rows.Err(); err != nil {
		return err
	}

	const selectCasts = `SELECT
		c.oid
		, c.castsource
		, c.casttarget
		, c.castcontext
	FROM pg_cast c ` //`
	c.casts = make(map[pgcastkey]*pgcast)
	rows, err = db.Query(selectCasts)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		cast := new(pgcast)
		err := rows.Scan(
			&cast.oid,
			&cast.source,
			&cast.target,
			&cast.context,
		)
		if err != nil {
			return err
		}
		c.casts[pgcastkey{target: cast.target, source: cast.source}] = cast
	}
	if err := rows.Err(); err != nil {
		return err
	}

	const showVersionNumber = `SHOW server_version_num` //`
	var version int
	if err := db.QueryRow(showVersionNumber).Scan(&version); err != nil {
		return err
	}

	var selectProcs string
	if version >= 110000 {
		// v11+
		selectProcs = `SELECT
			p.oid
			, p.proname
			, p.proargtypes[0]
			, p.prorettype
			, p.prokind = 'a'
		FROM pg_proc p
		WHERE p.pronargs = 1
		AND p.proname NOT LIKE 'pg_%'
		AND p.proname NOT LIKE '_pg_%'
		` //`
	} else {
		// pre v11
		selectProcs = `SELECT
			p.oid
			, p.proname
			, p.proargtypes[0]
			, p.prorettype
			, p.proisagg
		FROM pg_proc p
		WHERE p.pronargs = 1
		AND p.proname NOT LIKE 'pg_%'
		AND p.proname NOT LIKE '_pg_%'
		` //`
	}
	c.procs = make(map[funcName][]*pgproc)
	rows, err = db.Query(selectProcs)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		proc := new(pgproc)
		err := rows.Scan(
			&proc.oid,
			&proc.name,
			&proc.argtype,
			&proc.rettype,
			&proc.isagg,
		)
		if err != nil {
			return err
		}

		c.procs[funcName(proc.name)] = append(c.procs[funcName(proc.name)], proc)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	pgCatalogCache.m[key] = c
	return nil
}

type dataOperation uint8

const (
	dataNop dataOperation = iota
	dataRead
	dataWrite
	dataTest
)

type pgOID uint32

func (oid pgOID) getArrayOID() pgOID {
	return pgoidToArrayOID[oid]
}

// postgres types
const (
	pgtyp_any            pgOID = 2276
	pgtyp_bit            pgOID = 1560
	pgtyp_bitarr         pgOID = 1561
	pgtyp_bool           pgOID = 16
	pgtyp_boolarr        pgOID = 1000
	pgtyp_box            pgOID = 603
	pgtyp_boxarr         pgOID = 1020
	pgtyp_bpchar         pgOID = 1042
	pgtyp_bpchararr      pgOID = 1014
	pgtyp_bytea          pgOID = 17
	pgtyp_byteaarr       pgOID = 1001
	pgtyp_char           pgOID = 18
	pgtyp_chararr        pgOID = 1002
	pgtyp_cidr           pgOID = 650
	pgtyp_cidrarr        pgOID = 651
	pgtyp_circle         pgOID = 718
	pgtyp_circlearr      pgOID = 719
	pgtyp_date           pgOID = 1082
	pgtyp_datearr        pgOID = 1182
	pgtyp_daterange      pgOID = 3912
	pgtyp_daterangearr   pgOID = 3913
	pgtyp_float4         pgOID = 700
	pgtyp_float4arr      pgOID = 1021
	pgtyp_float8         pgOID = 701
	pgtyp_float8arr      pgOID = 1022
	pgtyp_inet           pgOID = 869
	pgtyp_inetarr        pgOID = 1041
	pgtyp_int2           pgOID = 21
	pgtyp_int2arr        pgOID = 1005
	pgtyp_int2vector     pgOID = 22
	pgtyp_int2vectorarr  pgOID = 1006
	pgtyp_int4           pgOID = 23
	pgtyp_int4arr        pgOID = 1007
	pgtyp_int4range      pgOID = 3904
	pgtyp_int4rangearr   pgOID = 3905
	pgtyp_int8           pgOID = 20
	pgtyp_int8arr        pgOID = 1016
	pgtyp_int8range      pgOID = 3926
	pgtyp_int8rangearr   pgOID = 3927
	pgtyp_interval       pgOID = 1186
	pgtyp_intervalarr    pgOID = 1187
	pgtyp_json           pgOID = 114
	pgtyp_jsonarr        pgOID = 199
	pgtyp_jsonb          pgOID = 3802
	pgtyp_jsonbarr       pgOID = 3807
	pgtyp_line           pgOID = 628
	pgtyp_linearr        pgOID = 629
	pgtyp_lseg           pgOID = 601
	pgtyp_lsegarr        pgOID = 1018
	pgtyp_macaddr        pgOID = 829
	pgtyp_macaddrarr     pgOID = 1040
	pgtyp_macaddr8       pgOID = 774
	pgtyp_macaddr8arr    pgOID = 775
	pgtyp_money          pgOID = 790
	pgtyp_moneyarr       pgOID = 791
	pgtyp_numeric        pgOID = 1700
	pgtyp_numericarr     pgOID = 1231
	pgtyp_numrange       pgOID = 3906
	pgtyp_numrangearr    pgOID = 3907
	pgtyp_oidvector      pgOID = 30
	pgtyp_path           pgOID = 602
	pgtyp_patharr        pgOID = 1019
	pgtyp_point          pgOID = 600
	pgtyp_pointarr       pgOID = 1017
	pgtyp_polygon        pgOID = 604
	pgtyp_polygonarr     pgOID = 1027
	pgtyp_text           pgOID = 25
	pgtyp_textarr        pgOID = 1009
	pgtyp_time           pgOID = 1083
	pgtyp_timearr        pgOID = 1183
	pgtyp_timestamp      pgOID = 1114
	pgtyp_timestamparr   pgOID = 1115
	pgtyp_timestamptz    pgOID = 1184
	pgtyp_timestamptzarr pgOID = 1185
	pgtyp_timetz         pgOID = 1266
	pgtyp_timetzarr      pgOID = 1270
	pgtyp_tsquery        pgOID = 3615
	pgtyp_tsqueryarr     pgOID = 3645
	pgtyp_tsrange        pgOID = 3908
	pgtyp_tsrangearr     pgOID = 3909
	pgtyp_tstzrange      pgOID = 3910
	pgtyp_tstzrangearr   pgOID = 3911
	pgtyp_tsvector       pgOID = 3614
	pgtyp_tsvectorarr    pgOID = 3643
	pgtyp_uuid           pgOID = 2950
	pgtyp_uuidarr        pgOID = 2951
	pgtyp_unknown        pgOID = 705
	pgtyp_varbit         pgOID = 1562
	pgtyp_varbitarr      pgOID = 1563
	pgtyp_varchar        pgOID = 1043
	pgtyp_varchararr     pgOID = 1015
	pgtyp_xml            pgOID = 142
	pgtyp_xmlarr         pgOID = 143

	pgtyp_hstore    pgOID = 9999
	pgtyp_hstorearr pgOID = 9998
)

// postgres type types
const (
	pgtyptype_base      = "b"
	pgtyptype_composite = "c"
	pgtyptype_domain    = "d"
	pgtyptype_enum      = "e"
	pgtyptype_pseudo    = "p"
	pgtyptype_range     = "r"
)

// postgres type categories
const (
	pgtypcategory_array       = "A"
	pgtypcategory_boolean     = "B"
	pgtypcategory_composite   = "C"
	pgtypcategory_datetime    = "D"
	pgtypcategory_enum        = "E"
	pgtypcategory_geometric   = "G"
	pgtypcategory_netaddress  = "I"
	pgtypcategory_numeric     = "N"
	pgtypcategory_pseudo      = "P"
	pgtypcategory_range       = "R"
	pgtypcategory_string      = "S"
	pgtypcategory_timespan    = "T"
	pgtypcategory_userdefined = "U"
	pgtypcategory_bitstring   = "V"
	pgtypcategory_unknown     = "X"
)

// postgres constraint types
const (
	pgconstraint_check     = "c"
	pgconstraint_fkey      = "f"
	pgconstraint_pkey      = "p"
	pgconstraint_unique    = "u"
	pgconstraint_trigger   = "t"
	pgconstraint_exclusion = "x"
)

// postgres cast contexts
const (
	pgcast_explicit   = "e"
	pgcast_implicit   = "i"
	pgcast_assignment = "a"
)

var pgoidToArrayOID = map[pgOID]pgOID{
	pgtyp_bit:         pgtyp_bitarr,
	pgtyp_bool:        pgtyp_boolarr,
	pgtyp_box:         pgtyp_boxarr,
	pgtyp_bpchar:      pgtyp_bpchararr,
	pgtyp_bytea:       pgtyp_byteaarr,
	pgtyp_char:        pgtyp_chararr,
	pgtyp_cidr:        pgtyp_cidrarr,
	pgtyp_circle:      pgtyp_circlearr,
	pgtyp_date:        pgtyp_datearr,
	pgtyp_daterange:   pgtyp_daterangearr,
	pgtyp_float4:      pgtyp_float4arr,
	pgtyp_float8:      pgtyp_float8arr,
	pgtyp_inet:        pgtyp_inetarr,
	pgtyp_int2:        pgtyp_int2arr,
	pgtyp_int2vector:  pgtyp_int2vectorarr,
	pgtyp_int4:        pgtyp_int4arr,
	pgtyp_int4range:   pgtyp_int4rangearr,
	pgtyp_int8:        pgtyp_int8arr,
	pgtyp_int8range:   pgtyp_int8rangearr,
	pgtyp_interval:    pgtyp_intervalarr,
	pgtyp_json:        pgtyp_jsonarr,
	pgtyp_jsonb:       pgtyp_jsonbarr,
	pgtyp_line:        pgtyp_linearr,
	pgtyp_lseg:        pgtyp_lsegarr,
	pgtyp_macaddr:     pgtyp_macaddrarr,
	pgtyp_macaddr8:    pgtyp_macaddr8arr,
	pgtyp_money:       pgtyp_moneyarr,
	pgtyp_numeric:     pgtyp_numericarr,
	pgtyp_numrange:    pgtyp_numrangearr,
	pgtyp_path:        pgtyp_patharr,
	pgtyp_point:       pgtyp_pointarr,
	pgtyp_polygon:     pgtyp_polygonarr,
	pgtyp_text:        pgtyp_textarr,
	pgtyp_time:        pgtyp_timearr,
	pgtyp_timestamp:   pgtyp_timestamparr,
	pgtyp_timestamptz: pgtyp_timestamptzarr,
	pgtyp_timetz:      pgtyp_timetzarr,
	pgtyp_tsquery:     pgtyp_tsqueryarr,
	pgtyp_tsrange:     pgtyp_tsrangearr,
	pgtyp_tstzrange:   pgtyp_tstzrangearr,
	pgtyp_tsvector:    pgtyp_tsvectorarr,
	pgtyp_uuid:        pgtyp_uuidarr,
	pgtyp_varbit:      pgtyp_varbitarr,
	pgtyp_varchar:     pgtyp_varchararr,
	pgtyp_xml:         pgtyp_xmlarr,
	pgtyp_hstore:      pgtyp_hstorearr,
}

var go2pgoids = map[goTypeId][]pgOID{
	goTypeBool:                     {pgtyp_bool},
	goTypeBoolSlice:                {pgtyp_boolarr},
	goTypeInt:                      {pgtyp_int4, pgtyp_int2, pgtyp_int8, pgtyp_float4, pgtyp_float8, pgtyp_numeric},
	goTypeIntSlice:                 {pgtyp_int4arr, pgtyp_int2arr, pgtyp_int2vector, pgtyp_int8arr, pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr},
	goTypeIntArray2:                {pgtyp_int4range, pgtyp_int8range, pgtyp_numrange},
	goTypeIntArray2Slice:           {pgtyp_int4rangearr, pgtyp_int8rangearr, pgtyp_numrangearr},
	goTypeInt8:                     {pgtyp_int2, pgtyp_int4, pgtyp_int8, pgtyp_float4, pgtyp_float8, pgtyp_numeric},
	goTypeInt8Slice:                {pgtyp_int2arr, pgtyp_int2vector, pgtyp_int4arr, pgtyp_int8arr, pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr},
	goTypeInt8SliceSlice:           {pgtyp_int2vectorarr},
	goTypeInt16:                    {pgtyp_int2, pgtyp_int4, pgtyp_int8, pgtyp_float4, pgtyp_float8, pgtyp_numeric},
	goTypeInt16Slice:               {pgtyp_int2arr, pgtyp_int2vector, pgtyp_int4arr, pgtyp_int8arr, pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr},
	goTypeInt16SliceSlice:          {pgtyp_int2vectorarr},
	goTypeInt32:                    {pgtyp_int4, pgtyp_int2, pgtyp_int8, pgtyp_float4, pgtyp_float8, pgtyp_numeric},
	goTypeInt32Slice:               {pgtyp_int4arr, pgtyp_int2arr, pgtyp_int8arr, pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr},
	goTypeInt32Array2:              {pgtyp_int4range, pgtyp_int8range, pgtyp_numrange},
	goTypeInt32Array2Slice:         {pgtyp_int4rangearr, pgtyp_int8rangearr, pgtyp_numrangearr},
	goTypeInt64:                    {pgtyp_int8, pgtyp_int4, pgtyp_int2, pgtyp_float4, pgtyp_float8, pgtyp_numeric, pgtyp_money},
	goTypeInt64Slice:               {pgtyp_int8arr, pgtyp_int4arr, pgtyp_int2arr, pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr, pgtyp_moneyarr},
	goTypeInt64Array2:              {pgtyp_int8range, pgtyp_int4range, pgtyp_numrange},
	goTypeInt64Array2Slice:         {pgtyp_int8rangearr, pgtyp_int4rangearr, pgtyp_numrangearr},
	goTypeUint:                     {pgtyp_int4, pgtyp_int2, pgtyp_int8, pgtyp_float4, pgtyp_float8, pgtyp_numeric},
	goTypeUintSlice:                {pgtyp_int4arr, pgtyp_int2arr, pgtyp_int2vector, pgtyp_int8arr, pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr},
	goTypeUint8:                    {pgtyp_int2, pgtyp_int4, pgtyp_int8, pgtyp_float4, pgtyp_float8, pgtyp_numeric},
	goTypeUint8Slice:               {pgtyp_int2arr, pgtyp_int2vector, pgtyp_int4arr, pgtyp_int8arr, pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr},
	goTypeUint16:                   {pgtyp_int2, pgtyp_int4, pgtyp_int8, pgtyp_float4, pgtyp_float8, pgtyp_numeric},
	goTypeUint16Slice:              {pgtyp_int2arr, pgtyp_int2vector, pgtyp_int4arr, pgtyp_int8arr, pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr},
	goTypeUint32:                   {pgtyp_int4, pgtyp_int2, pgtyp_int8, pgtyp_float4, pgtyp_float8, pgtyp_numeric},
	goTypeUint32Slice:              {pgtyp_int4arr, pgtyp_int2arr, pgtyp_int8arr, pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr},
	goTypeUint64:                   {pgtyp_int8, pgtyp_int4, pgtyp_int2, pgtyp_float4, pgtyp_float8, pgtyp_numeric, pgtyp_money},
	goTypeUint64Slice:              {pgtyp_int8arr, pgtyp_int4arr, pgtyp_int2arr, pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr, pgtyp_moneyarr},
	goTypeFloat32:                  {pgtyp_float4, pgtyp_float8, pgtyp_numeric},
	goTypeFloat32Slice:             {pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr},
	goTypeFloat64:                  {pgtyp_float4, pgtyp_float8, pgtyp_numeric},
	goTypeFloat64Slice:             {pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr},
	goTypeFloat64Array2:            {pgtyp_int4rangearr, pgtyp_int8rangearr, pgtyp_numrange, pgtyp_point},
	goTypeFloat64Array2Slice:       {pgtyp_numrangearr, pgtyp_path, pgtyp_pointarr, pgtyp_polygon},
	goTypeFloat64Array2SliceSlice:  {pgtyp_patharr, pgtyp_polygonarr},
	goTypeFloat64Array2Array2:      {pgtyp_box, pgtyp_lseg, pgtyp_path},
	goTypeFloat64Array2Array2Slice: {pgtyp_boxarr, pgtyp_lsegarr, pgtyp_patharr},
	goTypeFloat64Array3:            {pgtyp_line, pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr},
	goTypeFloat64Array3Slice:       {pgtyp_linearr},
	goTypeByte:                     {pgtyp_bpchar, pgtyp_bytea, pgtyp_char, pgtyp_text, pgtyp_varchar},
	goTypeByteSlice:                {pgtyp_bit, pgtyp_bpchar, pgtyp_bytea, pgtyp_char, pgtyp_cidr, pgtyp_inet, pgtyp_json, pgtyp_jsonb, pgtyp_macaddr, pgtyp_macaddr8, pgtyp_text, pgtyp_tsquery, pgtyp_tsvector, pgtyp_uuid, pgtyp_varbit, pgtyp_varchar, pgtyp_xml},
	goTypeByteSliceSlice:           {pgtyp_bitarr, pgtyp_bpchararr, pgtyp_byteaarr, pgtyp_chararr, pgtyp_cidrarr, pgtyp_inetarr, pgtyp_jsonarr, pgtyp_jsonbarr, pgtyp_macaddrarr, pgtyp_macaddr8arr, pgtyp_textarr, pgtyp_tsqueryarr, pgtyp_tsvectorarr, pgtyp_uuidarr, pgtyp_varbitarr, pgtyp_varchararr, pgtyp_xmlarr},
	goTypeString:                   {pgtyp_bit, pgtyp_bpchar, pgtyp_bytea, pgtyp_char, pgtyp_cidr, pgtyp_inet, pgtyp_json, pgtyp_jsonb, pgtyp_macaddr, pgtyp_macaddr8, pgtyp_text, pgtyp_tsquery, pgtyp_tsvector, pgtyp_uuid, pgtyp_varbit, pgtyp_varchar, pgtyp_xml},
	goTypeStringSlice:              {pgtyp_bitarr, pgtyp_bpchararr, pgtyp_byteaarr, pgtyp_chararr, pgtyp_cidrarr, pgtyp_inetarr, pgtyp_jsonarr, pgtyp_jsonbarr, pgtyp_macaddrarr, pgtyp_macaddr8arr, pgtyp_textarr, pgtyp_tsqueryarr, pgtyp_tsvectorarr, pgtyp_uuidarr, pgtyp_varbitarr, pgtyp_varchararr, pgtyp_xmlarr},
	goTypeByteArray16:              {pgtyp_uuid, pgtyp_bytea, pgtyp_text, pgtyp_varchar},
	goTypeByteArray16Slice:         {pgtyp_uuidarr, pgtyp_byteaarr, pgtyp_textarr, pgtyp_varchararr},
	goTypeRune:                     {pgtyp_char, pgtyp_bytea, pgtyp_varchar, pgtyp_bpchar},
	goTypeRuneSlice:                {pgtyp_bpchar, pgtyp_char, pgtyp_text, pgtyp_varchar /* the rest assumes typmod=1 */, pgtyp_bpchararr, pgtyp_chararr, pgtyp_text, pgtyp_varchararr},
	goTypeRuneSliceSlice:           {pgtyp_bpchararr, pgtyp_byteaarr, pgtyp_chararr, pgtyp_textarr, pgtyp_varchararr},
	goTypeIPNet:                    {pgtyp_cidr, pgtyp_inet, pgtyp_text, pgtyp_bpchar, pgtyp_bytea, pgtyp_char, pgtyp_varchar},
	goTypeIPNetSlice:               {pgtyp_cidrarr, pgtyp_inetarr, pgtyp_textarr, pgtyp_bpchararr, pgtyp_byteaarr, pgtyp_chararr, pgtyp_varchararr},
	goTypeTime:                     {pgtyp_date, pgtyp_time, pgtyp_timestamp, pgtyp_timestamptz, pgtyp_timetz},
	goTypeTimeSlice:                {pgtyp_datearr, pgtyp_timearr, pgtyp_timestamparr, pgtyp_timestamptzarr, pgtyp_timetzarr, pgtyp_intervalarr},
	goTypeTimeArray2:               {pgtyp_daterange, pgtyp_tsrange, pgtyp_tstzrange},
	goTypeTimeArray2Slice:          {pgtyp_daterangearr, pgtyp_tsrangearr, pgtyp_tstzrangearr},
	goTypeBigInt:                   {pgtyp_numeric, pgtyp_text},
	goTypeBigIntSlice:              {pgtyp_numericarr, pgtyp_textarr},
	goTypeBigIntArray2:             {pgtyp_numrange},
	goTypeBigIntArray2Slice:        {pgtyp_numrangearr},
	goTypeBigFloat:                 {pgtyp_numeric, pgtyp_text},
	goTypeBigFloatSlice:            {pgtyp_numericarr, pgtyp_textarr},
	goTypeBigFloatArray2:           {pgtyp_numrange},
	goTypeBigFloatArray2Slice:      {pgtyp_numrangearr},

	// NOTE(mkopriva): The hstore pgoids for these 4 go types are returned by
	// the typeoids method, this is because hstore doesn't have a "common" oid.
	//
	// gotypstringm:      {},
	// gotypstringms:     {},
	// gotypnullstringm:  {},
	// gotypnullstringms: {},
}

// Map of supported predicates to *equivalent* postgres comparison operators. For example
// the constructs IN and NOT IN are essentially the same as comparing the LHS to every
// element of the RHS with the operators "=" and "<>" respectively, and therefore the
// isin maps to "=" and notin maps to "<>".
var predicateToBasePGOps = map[predicate]string{
	isEQ:        "=",
	notEQ:       "<>",
	notEQ2:      "<>",
	isLT:        "<",
	isGT:        ">",
	isLTE:       "<=",
	isGTE:       ">=",
	isMatch:     "~",
	isMatchi:    "~*",
	notMatch:    "!~",
	notMatchi:   "!~*",
	isDistinct:  "<>",
	notDistinct: "=",
	isLike:      "~~",
	notLike:     "!~~",
	isILike:     "~~*",
	notILike:    "!~~*",
	isSimilar:   "~~",
	notSimilar:  "!~~",
	isIn:        "=",
	notIn:       "<>",
}

var rxUsingMethod = regexp.MustCompile(`(?i:\susing\s)`)

// parseIndexExpr extracts the index_expression from a CREATE INDEX command.
// e.g. CREATE INDEX index_name ON table_name USING method_name (index_expression);
func parseIndexExpr(s string) (expr string) {
	loc := rxUsingMethod.FindStringIndex(s)
	if len(loc) > 1 {
		s = s[loc[1]:]
	}
	if i := strings.Index(s, "("); i < 0 {
		return
	} else {
		s = s[i+1:]
	}

	var (
		r        rune
		nested   int  // number nested parentheses
		position int  // position in the input
		quoted   bool // quoted text in the input
		escaped  bool // escaped quote `''` in quoted text
	)
	for position, r = range s {
		if !quoted {
			if r == '(' {
				nested += 1 // nest
			} else if r == ')' {
				if nested == 0 {
					break // done
				}
				nested -= 1 // unnest
			} else if r == '\'' {
				quoted = true
			}
		} else {
			if r == '\'' && len(s) > position {
				if s[position+1] == '\'' {
					escaped = true
				} else if escaped {
					escaped = false
				} else {
					quoted = false
				}
			}
		}
	}

	return s[:position]
}

// matchNumbers is a helper func that reports whether a and b both contain
// the same numbers regardless of the order.
func matchNumbers(a, b []int16) bool {
	if len(a) != len(b) {
		return false
	}

toploop:
	for _, x := range a {
		for _, y := range b {
			if x == y {
				continue toploop
			}
		}
		return false // x not found in b
	}
	return true
}

// helper type
type int2vec []int16

func (v *int2vec) Scan(src interface{}) error {
	if b, ok := src.([]byte); ok {
		elems := strings.Split(string(b), " ")
		for _, e := range elems {
			i, err := strconv.ParseInt(e, 10, 16)
			if err != nil {
				return err
			}
			*v = append(*v, int16(i))
		}
	}
	return nil
}

var pgCatalogCache = struct {
	sync.RWMutex
	m map[string]*pgCatalog
}{m: make(map[string]*pgCatalog)}
