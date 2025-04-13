package postgres

// TODO(mkopriva) an option to specify "insert only" columns (like created_at,
// or anything really that should not change after initially inserted).

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/frk/gosql/internal/analysis"
	"github.com/frk/gosql/internal/postgres/oid"

	"github.com/lib/pq"
)

var _ = log.Println

// DB handles the connection pool to the target postgres database
// and holds additional information about that database.
//
// DB is NOT safe for concurrent use, an instance of DB is intended to
// be reused by separate runs of the type checker, just not concurrently.
type DB struct {
	// The underlying *sql.DB pool handle.
	*sql.DB
	// The DSN used to open connections to the database.
	// Used also as the key for caching the loaded Catalog.
	dsn string
	// The name of the current database. (intended mainly for error reporting)
	name string
	// The name of the current user. (intended mainly for error reporting)
	user string
	// The search_path setting (original).
	searchpath string
	// The version number of the current database.
	version int
	// The catalog for the target database.
	catalog *Catalog
}

// Open opens a new connection pool to the dsn specified postgres
// database and loads the catalog information.
func Open(dsn string) (db *DB, err error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, &dbError{Code: errDatabaseOpen, DB: dbInfo{DSN: dsn}, Err: err}
	}
	defer func() {
		if err != nil {
			conn.Close()
		}
	}()

	if err := conn.Ping(); err != nil {
		return nil, &dbError{Code: errDatabaseOpen, DB: dbInfo{DSN: dsn}, Err: err}
	}

	db = &DB{DB: conn, dsn: dsn}
	if err := db.QueryRow(`SELECT current_database(), current_user`).Scan(&db.name, &db.user); err != nil {
		return nil, &dbError{Code: errDatabaseInit, DB: dbInfo{DSN: dsn}, Err: err}
	}
	if err := db.QueryRow(`SHOW search_path`).Scan(&db.searchpath); err != nil {
		return nil, &dbError{Code: errDatabaseInit, DB: dbInfo{DSN: dsn}, Err: err}
	}
	if err := db.QueryRow(`SHOW server_version_num`).Scan(&db.version); err != nil {
		return nil, &dbError{Code: errDatabaseInit, DB: dbInfo{DSN: dsn}, Err: err}
	}

	if db.catalog, err = loadCatalog(db, dsn); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) Close() error {
	if db.DB == nil {
		return nil
	}
	return db.DB.Close()
}

func (db *DB) dbError(e dbError) *dbError {
	e.DB.DSN = db.dsn
	e.DB.Name = db.name
	e.DB.User = db.user
	e.DB.SearchPath = db.searchpath
	return &e
}

// checker maintains the state of the type checker.
type checker struct {
	db   *DB
	info *analysis.Info
	// The force directive, or nil. Used primarily for resolving which
	// field read / writes to skip and which not to.
	force *analysis.ForceDirective
	// TODO
	optional *analysis.OptionalDirective
	// The identifier of the target relation.
	rid analysis.RelIdent
	// The target relation.
	rel *Relation
	// A map of column qualifier (alias or relname or "") to the
	// relation denoted by said qualifier.
	relMap map[string]*Relation
	// The result info
	res *TargetInfo

	// set only during the type checking of WhereBetweenStruct; unset otherwise
	wb *analysis.WhereBetweenStruct
}

func (c *checker) dbError(e dbError, ptr analysis.FieldPtr) *dbError {
	e.DB = dbInfo{DSN: c.db.dsn, Name: c.db.name, User: c.db.user, SearchPath: c.db.searchpath}

	tpos := c.info.FileSet.Position(c.info.TypeNamePos)
	e.Target.Pkg = c.info.PkgPath
	e.Target.Name = c.info.TypeName
	e.Target.File.Name = tpos.Filename
	e.Target.File.Line = tpos.Line

	// get rel from column
	if e.Rel.Relation == nil && e.Col.Column != nil {
		e.Rel.Relation = e.Col.Column.Relation
	} else if e.Rel.Relation == nil && len(e.Col.Id.Name) > 0 {
		e.Rel.Relation = c.relMap[e.Col.Id.Qualifier]
	}

	// get rel from rel id, or relid from rel
	if e.Rel.Relation == nil && len(e.Rel.Id.Name) > 0 {
		e.Rel.Relation = c.relMap[e.Rel.Id.Name]
	} else if e.Rel.Id.Name == "" && e.Rel.Relation != nil {
		e.Rel.Id = analysis.RelIdent{Name: e.Rel.Relation.Name, Qualifier: e.Rel.Relation.Schema}
	}

	// get field & source info from FieldPtr
	if f, ok := c.info.FieldMap[ptr]; ok {
		fpos := c.info.FileSet.Position(f.Var.Pos())
		e.Field.Name = f.Var.Name()
		e.Field.Type = f.Var.Type().String()
		e.Field.Tag = f.Tag
		e.Field.File.Name = fpos.Filename
		e.Field.File.Line = fpos.Line
	}

	// get field & source info for WhereBetweenStruct, if present
	if c.wb != nil {
		if f, ok := c.info.FieldMap[c.wb]; ok {
			fpos := c.info.FileSet.Position(f.Var.Pos())
			e.WBField.Name = f.Var.Name()
			e.WBField.Type = f.Var.Type().String()
			e.WBField.Tag = f.Tag
			e.WBField.File.Name = fpos.Filename
			e.WBField.File.Line = fpos.Line
		}
	}
	return &e
}

// TargetInfo is the result of type-checking a TargetStruct instance.
// It holds the type, and operation, specific information about the target.
type TargetInfo struct {
	Info   *analysis.Info
	Struct analysis.TargetStruct

	Reads    []*FieldRead
	Writes   []*FieldWrite
	Filters  []*FieldFilter
	PKeys    []*FieldWrite
	Joins    [][]TableJoinConditional
	Where    []WhereConditional
	Conflict *ConflictInfo
}

// Check type-checks the given TargetStruct against the connected-to postgres database.
func Check(db *DB, ts analysis.TargetStruct, info *analysis.Info) (_ *TargetInfo, err error) {
	c := &checker{db: db, info: info}
	c.relMap = make(map[string]*Relation)
	c.res = new(TargetInfo)
	c.res.Struct = ts
	c.res.Info = info

	switch s := ts.(type) {
	case *analysis.FilterStruct:
		if err := loadTargetRelation(c, s.Rel); err != nil {
			return nil, err
		}

		if err := typeCheckFilterStruct(c, s); err != nil {
			return nil, err
		}
	case *analysis.QueryStruct:
		if err := loadTargetRelation(c, s.Rel); err != nil {
			return nil, err
		}

		if err := typeCheckQueryStruct(c, s); err != nil {
			return nil, err
		}
	}
	return c.res, nil
}

// typeCheckFilterStruct type-checks the given FilterStruct.
func typeCheckFilterStruct(c *checker, fs *analysis.FilterStruct) error {
	checks := []func(*checker, *analysis.FilterStruct) error{
		typeCheckFilterRelField,
		typeCheckFilterTextSearchDirective,
	}
	for i := 0; i < len(checks); i++ {
		if err := checks[i](c, fs); err != nil {
			return err
		}
	}
	return nil
}

// typeCheckFilterRelField checks the individual fields of the given FilterStruct's RelField.Type.
func typeCheckFilterRelField(c *checker, fs *analysis.FilterStruct) error {
	for _, f := range fs.Rel.Type.Fields {
		if err := typeCheckFieldFilter(c, f, false); err != nil {
			return err
		}
	}
	return nil
}

// If a TextSearch directive was provided, make sure that the specified column
// is present in one of the loaded relations and that it has the correct type.
func typeCheckFilterTextSearchDirective(c *checker, fs *analysis.FilterStruct) error {
	if fs.TextSearch == nil {
		return nil
	}

	if col := findRelColumn(c.rel, fs.TextSearch.Name); col == nil {
		return c.dbError(dbError{Code: errColumnUnknown, Rel: relInfo{Relation: c.rel},
			Col: colInfo{Id: fs.TextSearch.ColIdent}}, fs.TextSearch)
	} else if col.Type.OID != oid.TSVector {
		return c.dbError(dbError{Code: errColumnTextSearchType, Rel: relInfo{Relation: c.rel},
			Col: colInfo{Id: fs.TextSearch.ColIdent, Column: col}}, fs.TextSearch)
	}
	return nil
}

// typeCheckQueryStruct type-checks the given QueryStruct.
func typeCheckQueryStruct(c *checker, qs *analysis.QueryStruct) error {
	c.optional = qs.Optional

	checks := []func(*checker, *analysis.QueryStruct) error{
		typeCheckQueryJoinStruct,
		typeCheckQueryForceDirective,
		typeCheckQueryDefaultDirective,
		typeCheckQueryOrderByDirective,

		typeCheckQueryReturnDirective,
		typeCheckQueryResultField,
		typeCheckQueryRelField,

		typeCheckQueryWhereStruct,
		typeCheckQueryOnConflictStruct,
	}
	for i := 0; i < len(checks); i++ {
		if err := checks[i](c, qs); err != nil {
			return err
		}
	}

	if qs.Kind == analysis.QueryKindInsert {
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

// typeCheckQueryForceDirective checks the columns listed in the gosql.Force directive's tag.
//
// CHECKLIST:
//
//	✅ Each column MUST be present in one of the loaded relations.
func typeCheckQueryForceDirective(c *checker, qs *analysis.QueryStruct) error {
	if qs.Force == nil {
		return nil
	}
	c.force = qs.Force

	for _, cid := range qs.Force.Items {
		if _, ecode := findColumn(c, cid); ecode > 0 {
			return c.dbError(dbError{Code: ecode, Col: colInfo{Id: cid}}, qs.Force)
		}
	}
	return nil
}

// typeCheckQueryDefaultDirective checks the columns listed in the gosql.Default directive's tag.
//
// CHECKLIST:
//
//	✅ Each column MUST be present in the TARGET relation.
//	✅ Each column MUST have a DEFAULT set.
//	✅ If a column has a qualifier it MUST match the alias,
//	   or name, of the target relation.
func typeCheckQueryDefaultDirective(c *checker, qs *analysis.QueryStruct) error {
	if qs.Default == nil {
		return nil
	}

	for _, cid := range qs.Default.Items {
		if col := findRelColumn(c.rel, cid.Name); col == nil {
			return c.dbError(dbError{Code: errColumnUnknown,
				Rel: relInfo{Relation: c.rel}, Col: colInfo{Id: cid}}, qs.Default)
		} else if !col.HasDefault {
			return c.dbError(dbError{Code: errColumnDefaultUnset,
				Rel: relInfo{Relation: c.rel}, Col: colInfo{Id: cid, Column: col}}, qs.Default)
		}
	}
	return nil
}

// typeCheckQueryOrderByDirective checks the columns listed in the gosql.OrderBy directive's tag.
//
// CHECKLIST:
//
//	✅ Each column MUST be present in one of the loaded relations.
func typeCheckQueryOrderByDirective(c *checker, qs *analysis.QueryStruct) error {
	if qs.OrderBy == nil {
		return nil
	}

	for _, item := range qs.OrderBy.Items {
		if _, ecode := findColumn(c, item.ColIdent); ecode > 0 {
			return c.dbError(dbError{Code: ecode, Col: colInfo{Id: item.ColIdent}}, qs.OrderBy)
		}
	}
	return nil
}

// typeCheckQueryReturnDirective checks the columns listed in the gosql.Return directive's tag.
//
// CHECKLIST:
//
//   - If "*" tag was used:
//     ✅ Each field of the target data type MAY have a corresponding column
//     in one of the loaded relations (denoted by the field's tag).
//     ✅ Each field of the target data type that has a corresponding column MUST
//     be of a type that is readable from a value of that column.
//
//   - If "<column_list>" tag was used:
//     ✅ Each listed column MUST be present in the TARGET relation.
//     ✅ Each listed column MUST have a corresponding field in the target data type.
//     ✅ Each listed column's qualifier, if it has one, MUST match the alias, or name,
//     of the TARGET relation.
//     ✅ Each listed column's corresponding field MUST be of a type that IS READABLE
//     from a value of that column's type.
func typeCheckQueryReturnDirective(c *checker, qs *analysis.QueryStruct) error {
	if qs.Return == nil {
		return nil
	}

	var strict bool
	var fields []*analysis.FieldInfo

	if qs.Return.All {
		strict = false
		fields = qs.Rel.Type.Fields
	} else {
		strict = true
		for _, cid := range qs.Return.Items {
			if f, ecode := findQueryColumnField(c, qs, cid, strict); ecode > 0 {
				return c.dbError(dbError{Code: ecode, Col: colInfo{Id: cid}}, qs.Return)
			} else if f != nil {
				fields = append(fields, f)
			}
		}
	}

	for _, f := range fields {
		if err := typeCheckFieldRead(c, f, strict); err != nil {
			return err
		}
	}
	return nil
}

// typeCheckQueryResultField checks the fields of the resultField struct data type.
//
// CHECKLIST:
//
//	✅ Each field of the result data type MUST have a corresponding column in
//	   one of the loaded relations.
//	✅ Each field of the result data type MUST be of a type that is readable
//	   from a value of the corresponding column.
func typeCheckQueryResultField(c *checker, qs *analysis.QueryStruct) error {
	if qs.Result == nil {
		return nil
	}

	for _, f := range qs.Result.Type.Fields {
		if err := typeCheckFieldRead(c, f, true); err != nil {
			return err
		}
	}
	return nil
}

// typeCheckQueryRelField checks the fields of the query's struct data type.
//
// CHECKLIST:
//   - If the query represents a SELECT:
//     ✅ Each field of the query's struct data type MUST be of a type that
//     is readable from a value of the corresponding column.
//   - If the query represents an INSERT or UPDATE:
//     ✅ Each field of the query's struct data type MUST be of a type that
//     is writeable to the corresponding column.
func typeCheckQueryRelField(c *checker, qs *analysis.QueryStruct) error {
	if qs.Rel == nil || qs.Rel.IsDirective {
		return nil
	}

	if qs.Kind == analysis.QueryKindSelect {
		for _, f := range qs.Rel.Type.Fields {
			strict := (c.optional == nil || !c.optional.Contains(f.ColIdent))
			if err := typeCheckFieldRead(c, f, strict); err != nil {
				return err
			}
		}
	} else if qs.Kind == analysis.QueryKindInsert || qs.Kind == analysis.QueryKindUpdate {
		for _, f := range qs.Rel.Type.Fields {
			w, err := typeCheckFieldWrite(c, f)
			if err != nil {
				return err
			}

			// ✅ If strict=true; the column denoted by the
			// given field MUST be present in the target relation.
			strict := (f.IsReadOnly() == false)
			if w == nil && strict {
				return c.dbError(dbError{Code: errColumnUnknown,
					Col: colInfo{Id: f.ColIdent}, Rel: relInfo{Relation: c.rel}}, f)
			}

			if w != nil {
				if w.Column.IsPrimary {
					c.res.PKeys = append(c.res.PKeys, w)
				}
				switch {
				case qs.Kind == analysis.QueryKindInsert && !skipFieldInsert(c, w):
					c.res.Writes = append(c.res.Writes, w)
				case qs.Kind == analysis.QueryKindUpdate && !skipFieldUpdate(c, w):
					c.res.Writes = append(c.res.Writes, w)
				}
			}
		}
	}
	return nil
}

// typeCheckQueryJoinStruct ...
func typeCheckQueryJoinStruct(c *checker, qs *analysis.QueryStruct) error {
	if qs.Join == nil {
		return nil
	}

	if qs.Join.Relation != nil {
		rid := qs.Join.Relation.RelIdent
		if _, err := loadJoinRelation(c, rid, qs.Join.Relation); err != nil {
			return err
		}
	}

	for _, dir := range qs.Join.Directives {
		rel, err := loadJoinRelation(c, dir.RelIdent, dir)
		if err != nil {
			return err
		}

		conds := make([]TableJoinConditional, len(dir.TagItems))
		for i, v := range dir.TagItems {
			switch item := v.(type) {
			case *analysis.JoinBoolTagItem:
				cond := new(Boolean)
				cond.Value = item.Value
				conds[i] = cond
			case *analysis.JoinConditionTagItem:
				cond := new(ColumnConditional)
				cond.LHSColIdent = item.LHSColIdent
				cond.RHSColIdent = item.RHSColIdent
				cond.RHSLiteral = item.RHSLiteral
				cond.Predicate = item.Predicate
				cond.Quantifier = item.Quantifier
				if err := typeCheckColumnConditional(c, cond, rel, dir.RelIdent, dir); err != nil {
					return err
				}
				conds[i] = cond
			default:
				// NOTE(mkopriva): currently predicates other than
				// analysis.JoinConditionTagItem are not supported
				// as a join condition.
			}
		}

		c.res.Joins = append(c.res.Joins, conds)
	}
	return nil
}

// typeCheckQueryWhereStruct type-checks individual items of the query's WhereStruct.
func typeCheckQueryWhereStruct(c *checker, qs *analysis.QueryStruct) (err error) {
	if qs.Where == nil {
		return nil
	}
	for _, item := range qs.Where.Items {
		cond, err := typeCheckWhereItem(c, item)
		if err != nil {
			return err
		}
		c.res.Where = append(c.res.Where, cond)
	}
	// XXX if qs.TypeName == "DeleteWithUsingJoinBlock1Query" {
	// XXX 	log.Printf("%#v\n", c.res.Where)
	// XXX }
	return nil
}

// typeCheckWhereItem type-checks the given WhereItem and returns the resulting WhereConditional.
func typeCheckWhereItem(c *checker, item analysis.WhereItem) (WhereConditional, error) {
	switch wi := item.(type) {
	case *analysis.WhereBoolTag:
		boolean := new(Boolean)
		boolean.Value = wi.Value
		return boolean, nil
	case *analysis.WhereStruct:
		nested := new(NestedConditional)
		nested.FieldName = wi.FieldName
		for _, item := range wi.Items {
			cond, err := typeCheckWhereItem(c, item)
			if err != nil {
				return nil, err
			}
			nested.Conditionals = append(nested.Conditionals, cond)
		}
		return nested, nil
	case *analysis.WhereColumnDirective:
		column := new(ColumnConditional)
		column.LHSColIdent = wi.LHSColIdent
		column.RHSColIdent = wi.RHSColIdent
		column.RHSLiteral = wi.RHSLiteral
		column.Predicate = wi.Predicate
		column.Quantifier = wi.Quantifier
		if err := typeCheckColumnConditional(c, column, nil, analysis.RelIdent{}, wi); err != nil {
			return nil, err
		}
		return column, nil
	case *analysis.WhereStructField:
		field := new(FieldConditional)
		field.FieldName = wi.Name
		field.FieldType = wi.Type
		field.ColIdent = wi.ColIdent
		field.Predicate = wi.Predicate
		field.Quantifier = wi.Quantifier
		field.FuncName = wi.FuncName

		if field.Column == nil {
			col, ecode := findColumn(c, field.ColIdent)
			if ecode > 0 {
				return nil, c.dbError(dbError{Code: ecode, Col: colInfo{Id: field.ColIdent}}, wi)
			}
			field.Column = col
		}

		if ecode := typeCheckFieldConditional(c, field); ecode > 0 {
			return nil, c.dbError(dbError{Code: ecode,
				Col:  colInfo{Id: wi.ColIdent, Column: field.Column},
				Pred: wi.Predicate, Quant: wi.Quantifier, Func: field.FuncName}, wi)
		}
		return field, nil
	case *analysis.WhereBetweenStruct:
		between, err := typeCheckWhereBetweenStruct(c, wi)
		if err != nil {
			return nil, err
		}
		return between, nil
	default:
		panic("shouldn't reach")
	}
	return nil, nil
}

// typeCheckWhereBetweenStruct
//
// CHECKLIST:
//
//		✅ The column denoting the primary predicand MUST be present in one of the loaded relations.
//		✅ Both range-bound predicands MUST be comparable to the primary column predicand.
//
//	  - If the bound is a *analysis.BetweenColumnDirective:
//	    ✅ The denoted column MUST be present in one of the loaded relations.
//	    ✅ The found column MUST be comparable to the provided acol column.
//
// - If the bound is a *analysis.BetweenStructField:
func typeCheckWhereBetweenStruct(c *checker, wb *analysis.WhereBetweenStruct) (*BetweenConditional, error) {
	c.wb = wb
	defer func() { c.wb = nil }()

	col, ecode := findColumn(c, wb.ColIdent)
	if ecode > 0 {
		return nil, c.dbError(dbError{Code: ecode, Col: colInfo{Id: wb.ColIdent}}, wb)
	}

	rangeInfo := [2]*struct {
		out  RangeBound          // result of type checker
		in   analysis.RangeBound // input from analyzer
		pred analysis.Predicate  // synonymous predicate
	}{{}, {}}
	rangeInfo[0].in = wb.LowerBound
	rangeInfo[1].in = wb.UpperBound

	// Set the predicate here for the purposes of type-checking the comparison.
	if wb.Predicate.IsOneOf(analysis.NotBetween, analysis.NotBetweenSym, analysis.NotBetweenAsym) {
		// "a NOT BETWEEN x AND y" is equivalent to "a < x OR a > y"
		rangeInfo[0].pred = analysis.IsLT
		rangeInfo[1].pred = analysis.IsGT
	} else {
		// "a BETWEEN x AND y" is equivalent to "a >= x AND a <= y"
		rangeInfo[0].pred = analysis.IsGTE
		rangeInfo[1].pred = analysis.IsLTE
	}

	for _, r := range rangeInfo {
		switch b := r.in.(type) {
		default:
			panic("shouldn't reach")
		case *analysis.BetweenColumnDirective:
			column := new(ColumnConditional)
			column.LHSColIdent = wb.ColIdent
			column.LHSColumn = col
			column.RHSColIdent = b.ColIdent
			column.Predicate = r.pred
			if err := typeCheckColumnConditional(c, column, nil, analysis.RelIdent{}, b); err != nil {
				if e, ok := err.(*dbError); ok && e.Code == errColumnComparison {
					return nil, c.dbError(dbError{Code: errBetweenColumnComparison,
						Col:    colInfo{Id: wb.ColIdent, Column: col},
						RHSCol: colInfo{Id: b.ColIdent, Column: column.RHSColumn},
						Pred:   r.pred}, b)
				}

				return nil, err
			}
			r.out = column
		case *analysis.BetweenStructField:
			field := new(FieldConditional)
			field.FieldName = b.Name
			field.FieldType = b.Type
			field.ColIdent = wb.ColIdent
			field.Column = col
			field.Predicate = r.pred
			if ecode := typeCheckFieldConditional(c, field); ecode > 0 {
				if ecode == errColumnFieldComparison {
					return nil, c.dbError(dbError{Code: errBetweenFieldComparison,
						Col:  colInfo{Id: wb.ColIdent, Column: col},
						Pred: r.pred, Func: field.FuncName}, b)
				}
				return nil, c.dbError(dbError{Code: ecode,
					Col:  colInfo{Id: wb.ColIdent, Column: col},
					Pred: r.pred, Func: field.FuncName}, b)
			}
			r.out = field
		}
	}

	cond := new(BetweenConditional)
	cond.FieldName = wb.FieldName
	cond.ColIdent = wb.ColIdent
	cond.Column = col
	cond.Predicate = wb.Predicate
	cond.LowerBound = rangeInfo[0].out
	cond.UpperBound = rangeInfo[1].out
	return cond, nil
}

// typeCheckFieldConditional
//
// CHECKLIST:
//
//	✅ The column denoted by the ColIdent MUST be present in one of the loaded relations.
//	☑️  TODO If the column HasNotNull=true the predicate MUST NOT be a NULL predicate.
//	✅ The pg_operator table MUST contain an entry for the combination of the column's type,
//	   the predicate's operator, and a pg type compatible with the field's type. (Note that,
//	   if a quantifier [ANY, ALL, etc.] was provided, or the predicate is an array predicate,
//	   then the field is expected to be a slice/array and its element type will be used for
//	   the pg_operator lookup).
//	☑️  If a modifier function (lower, upper, etc.) was provided the column's type MUST
//	   match the function's argument type, the field's type MUST be compatible with the
//	   function's argument type, and the result type of both instances of the function
//	   MUST be comparable given the predicate operator, (be mindful of function overloading).
func typeCheckFieldConditional(c *checker, cond *FieldConditional) dbErrorCode {
	if len(cond.FuncName) > 0 {
		return typeCheckFieldConditionalWithFunc(c, cond)
	}

	var (
		ftyp    = cond.FieldType
		ctyp    = cond.Column.Type
		coid    = cond.Column.Type.OID
		oprname = predicateToOprname[cond.Predicate]
	)

	// XXX(mkopriva): not sure this is the correct approach ...
	// The checker can either trust the user to provide a value in the correct
	// format, or it could "request" more information by specifying that some
	// Go types need to implement certain methods that would indicate the pg
	// type with which they are compatible ...
	if !ftyp.Kind.IsBasic() && ftyp.IsValuer {
		ftyp = analysis.TypeInfo{Kind: analysis.TypeKindString}
	}

	if ctyp.Type == TypeTypeEnum { // enum labels are basically text
		coid = oid.Text
	}

	if cond.Predicate.IsArray() || cond.Quantifier > 0 {
		// FIXME(mkopriva): this currently breaks for custom types like
		// enums. Must figure out how to find the oid of a custom type...
		id, ok := oid.TypeToArray[cond.Column.Type.OID]
		if ok {
			ctyp, ok = c.db.catalog.Types[id]
		}
		if !ok && cond.Predicate.IsArray() {
			return errPredicateOperandArray
		} else if !ok && cond.Quantifier > 0 {
			log.Println("EHHEHEHEHE")
			return errPredicateOperandQuantifier
		}
	}

	key := OpKey{Left: coid, Right: coid, Name: oprname}
	if _, ok := c.db.catalog.Operators[key]; ok {
		typmod1 := isLength1Type(c, cond.Column)
		if comp := typeCompatibility(c, ctyp, ftyp, typmod1); comp != nil {
			cond.Valuer = comp.valuer
			return 0
		}
	}

	var comp *compentry
	for key := range c.db.catalog.Operators {
		if key.Name != oprname || key.Left != coid {
			continue
		}

		typ, ok := c.db.catalog.Types[key.Right]
		if !ok {
			continue
		}

		if ce := typeCompatibility(c, typ, ftyp, false); ce != nil {
			if ce.valuer == "" {
				comp = ce
				break
			} else if comp == nil {
				comp = ce
			}
		}
	}
	if comp == nil {
		return errColumnFieldComparison
	}

	cond.Valuer = comp.valuer
	return 0
}

func typeCheckFieldConditionalWithFunc(c *checker, cond *FieldConditional) dbErrorCode {
	procs, ok := c.db.catalog.Procs[string(cond.FuncName)]
	if !ok {
		return errProcedureUnknown
	}

	// - 1. produce a set of LR proc return type & comparison op
	//      combinations that are valid
	// - 2. loop over the set:
	// - 3. lhs column type & L proc arg type must be coercible,castable
	// - 4. rhs field type must be compatible with a pg_type
	//      that is coercible/castable to R proc arg type
	//
	type pair struct{ L, R *Type }

	pairs := []pair{}
	for _, l := range procs {
		lrettyp, ok := c.db.catalog.Types[l.RetType]
		if !ok {
			continue
		}

		for _, r := range procs {
			rrettyp, ok := c.db.catalog.Types[r.RetType]
			if !ok {
				continue
			}

			if typeCheckComparison(c, lrettyp, rrettyp, cond.Predicate, cond.Quantifier) == 0 {
				largtyp, ok := c.db.catalog.Types[l.ArgType]
				if !ok {
					continue
				}
				rargtyp, ok := c.db.catalog.Types[r.ArgType]
				if !ok {
					continue
				}
				pairs = append(pairs, pair{L: largtyp, R: rargtyp})
			}
		}
	}

	for _, p := range pairs {
		if checkTypeCoercion(c, p.L.OID, cond.Column.Type.OID) {
			// XXX(mkopriva): not sure this is the correct approach ...
			// The checker can either trust the user to provide a value in the correct
			// format, or it could "request" more information by specifying that some
			// Go types need to implement certain methods that would indicate the pg
			// type with which they are compatible ...
			ftyp := cond.FieldType
			if !ftyp.Kind.IsBasic() && ftyp.IsValuer {
				ftyp = analysis.TypeInfo{Kind: analysis.TypeKindString}
			}

			if comp := typeCompatibility(c, p.R, ftyp, false); comp != nil {
				cond.Valuer = comp.valuer
				return 0
			}
		}
	}
	return errProcedureUnknown
}

// typeCheckColumnConditional
//
// CHECKLIST:
//
//   - If the *Relation argument is not nil:
//     ✅ The LHS column MUST be present in given relation.
//
//   - If the *Relation argument is nil:
//     ✅ The LHS column MUST be present in one of the loaded relations.
//
//   - If the predicate of the condition argument is unary:
//     ✅ If the unary predicate is one of the "IS [NOT] { FALSE | TRUE | UNKNOWN }"
//     predicates then the type of the LHS column/expression MUST be boolean.
//     ✅ If the unary predicate is one of the "IS [NOT] NULL" predicates then
//     the LHS column MUST NOT have the "NOT NULL" constraint.
//
//   - If a predicate quantifier was provided:
//     ✅ The RHS column or literal expression MUST be quantifiable.
//
//     ✅ The LHS and RHS types MUST be comparable with the given predicate and quantifier.
func typeCheckColumnConditional(c *checker, cond *ColumnConditional, rel *Relation, rid analysis.RelIdent, ptr analysis.FieldPtr) error {
	if cond.LHSColumn == nil {
		if rel != nil {
			col := findRelColumn(rel, cond.LHSColIdent.Name)
			if col == nil {
				return c.dbError(dbError{Code: errColumnUnknown,
					Rel: relInfo{Id: rid, Relation: rel},
					Col: colInfo{Id: cond.LHSColIdent}}, ptr)
			}
			cond.LHSColumn = col
		} else {
			col, ecode := findColumn(c, cond.LHSColIdent)
			if ecode > 0 {
				return c.dbError(dbError{Code: ecode,
					Rel: relInfo{Id: rid, Relation: rel},
					Col: colInfo{Id: cond.LHSColIdent}}, ptr)
			}
			cond.LHSColumn = col
		}
	}

	if cond.Predicate.IsUnary() {
		if cond.Predicate.IsBoolean() && cond.LHSColumn.Type.OID != oid.Bool {
			return c.dbError(dbError{Code: errPredicateOperandBool,
				Col: colInfo{Id: cond.LHSColIdent, Column: cond.LHSColumn}, Pred: cond.Predicate}, ptr)
		}

		if cond.Predicate.IsNull() {
			// if this is a WHERE condition and the referenced column is from a joined table,
			// then allow IS [NOT] NULL even if that column has the NOT NULL constraint since
			// the table's join condition could have been unsuccessful in which case the
			// referenced non-nullable column would be NULL.
			var isWhereJoined bool
			if _, ok := ptr.(*analysis.WhereColumnDirective); ok {
				isWhereJoined = (c.rel != cond.LHSColumn.Relation)
			}
			if cond.LHSColumn.HasNotNull && !isWhereJoined {
				return c.dbError(dbError{Code: errPredicateOperandNull,
					Col: colInfo{Id: cond.LHSColIdent, Column: cond.LHSColumn}, Pred: cond.Predicate}, ptr)
			}
		}
		return nil
	}

	if len(cond.RHSColIdent.Name) > 0 {
		col, ecode := findColumn(c, cond.RHSColIdent)
		if ecode > 0 {
			return c.dbError(dbError{Code: ecode, Col: colInfo{Id: cond.RHSColIdent}}, ptr)
		}
		cond.RHSColumn = col
		cond.RHSType = col.Type
	} else if len(cond.RHSLiteral) > 0 {
		typ, ecode := typeOfLiteral(c, cond.RHSLiteral)
		if ecode > 0 {
			return c.dbError(dbError{Code: ecode,
				Col:    colInfo{Id: cond.LHSColIdent, Column: cond.LHSColumn},
				RHSLit: exprInfo{Expr: cond.RHSLiteral},
				Pred:   cond.Predicate, Quant: cond.Quantifier}, ptr)
		}
		cond.RHSType = typ
	}

	if ecode := typeCheckComparison(c, cond.LHSColumn.Type, cond.RHSType, cond.Predicate, cond.Quantifier); ecode > 0 {
		var rhsCol colInfo
		var rhsLit exprInfo
		if len(cond.RHSColIdent.Name) > 0 {
			rhsCol = colInfo{Id: cond.RHSColIdent, Column: cond.RHSColumn}
		} else {
			rhsLit = exprInfo{Expr: cond.RHSLiteral, Type: cond.RHSType}
		}
		return c.dbError(dbError{Code: ecode,
			Col:    colInfo{Id: cond.LHSColIdent, Column: cond.LHSColumn},
			RHSCol: rhsCol, RHSLit: rhsLit, Pred: cond.Predicate, Quant: cond.Quantifier}, ptr)
	}
	return nil
}

// typeCheckComparison checks whether or not a valid comparison operation
// expression can be generated from the provided arguments.
//
// CHECKLIST:
//
//	✅ If the predicate together with the quantifier constitute an array comparison
//	   the check MUST use the RHS's element type instead of the RHS type.
//	✅ ACCEPT if the LHS type belongs to the string category and the RHS type is unknown.
//	✅ ACCEPT if the combination of LHS type, RHS type, and the predicate has
//	   an entry in the pg_operator table.
func typeCheckComparison(c *checker, ltyp *Type, rtyp *Type, pred analysis.Predicate, qua analysis.Quantifier) dbErrorCode {
	if pred.IsArray() || qua > 0 {
		if rtyp.Category != TypeCategoryArray {
			if pred.IsArray() {
				return errPredicateOperandArray
			}
			return errPredicateOperandQuantifier
		}
		rtyp = c.db.catalog.Types[rtyp.Elem]
	}

	if ltyp.OID == rtyp.OID {
		return 0
	}

	if (ltyp.Category == TypeCategoryString || ltyp.Type == TypeTypeEnum) && rtyp.OID == oid.Unknown {
		return 0
	}

	// if operator exists, accept
	opkey := OpKey{Left: ltyp.OID, Right: rtyp.OID, Name: predicateToOprname[pred]}
	if _, ok := c.db.catalog.Operators[opkey]; ok {
		return 0
	}

	// if L & R are not the same type, try *implicit* casting.
	if opkey.Left != opkey.Right {
		castkey := CastKey{Source: opkey.Left, Target: opkey.Right}
		if cast := c.db.catalog.Casts[castkey]; cast != nil && cast.Context == CastContextImplicit {
			k := OpKey{Left: cast.Target, Right: cast.Target, Name: opkey.Name}
			if _, ok := c.db.catalog.Operators[k]; ok {
				return 0
			}
		}

		// switch target & source
		castkey = CastKey{Source: opkey.Right, Target: opkey.Left}
		if cast := c.db.catalog.Casts[castkey]; cast != nil && cast.Context == CastContextImplicit {
			k := OpKey{Left: cast.Target, Right: cast.Target, Name: opkey.Name}
			if _, ok := c.db.catalog.Operators[k]; ok {
				return 0
			}
		}
	}

	return errColumnComparison
}

// typeCheckQueryOnConflictStruct
//
// CHECKLIST:
//
//	✅ If a gosql.Column directive was used in the OnConflict block, the columns
//	   listed in the directive's tag MUST be present in the target table, they
//	   also MUST constitute a unique index of the target table.
//	✅ If a gosql.Index directive was used in the OnConflict block, the index
//	   specified in the directive's tag MUST be present on the target table
//	   and it MUST be a unique index.
//	✅ If a gosql.Constraint directive was used in the OnConflict block, the
//	   constraint specified in the directive's tag MUST be present on the target
//	   table and it MUST be a unique constraint.
//	✅ If a gosql.Update directive was used in the OnConflict block, the columns
//	   listed in the directive's tag MUST be present in the target table.
func typeCheckQueryOnConflictStruct(c *checker, qs *analysis.QueryStruct) error {
	if qs.OnConflict == nil {
		return nil
	}

	info := new(ConflictInfo)

	// check the column list and ensure a unique index exists
	if qs.OnConflict.Column != nil {
		var attnums []int16
		for _, cid := range qs.OnConflict.Column.ColIdents {
			col := findRelColumn(c.rel, cid.Name)
			if col == nil {
				return c.dbError(dbError{Code: errColumnUnknown,
					Rel: relInfo{Relation: c.rel}, Col: colInfo{Id: cid}},
					qs.OnConflict.Column)
			}
			attnums = append(attnums, col.Num)
		}

		var exists, isunique bool
		for _, ind := range c.rel.Indexes {
			if !matchNumbers(ind.Key, attnums) {
				continue
			}
			exists = true

			if !ind.IsUnique && !ind.IsPrimary {
				continue
			}
			isunique = true

			target := new(ConflictIndex)
			target.Expression = ind.Expression
			target.Predicate = ind.Predicate
			info.Target = target
			break
		}

		if !exists {
			return c.dbError(dbError{Code: errOnConflictIndexColumnsUnknown,
				Rel: relInfo{Relation: c.rel}}, qs.OnConflict.Column)
		} else if !isunique {
			return c.dbError(dbError{Code: errOnConflictIndexColumnsNotUnique,
				Rel: relInfo{Relation: c.rel}}, qs.OnConflict.Column)
		}
	}

	// check that the index exists and is unique.
	if qs.OnConflict.Index != nil {
		ind := findRelIndex(c.rel, qs.OnConflict.Index.Name)
		if ind == nil {
			return c.dbError(dbError{Code: errOnConflictIndexUnknown,
				Rel: relInfo{Relation: c.rel}}, qs.OnConflict.Index)
		} else if !ind.IsUnique && !ind.IsPrimary {
			return c.dbError(dbError{Code: errOnConflictIndexNotUnique,
				Rel: relInfo{Relation: c.rel}}, qs.OnConflict.Index)
		}

		target := new(ConflictIndex)
		target.Expression = ind.Expression
		target.Predicate = ind.Predicate
		info.Target = target
	}

	// check that the constraint exists and is unique.
	if qs.OnConflict.Constraint != nil {
		con := findRelConstraint(c.rel, qs.OnConflict.Constraint.Name)
		if con == nil {
			return c.dbError(dbError{Code: errOnConflictConstraintUnknown,
				Rel: relInfo{Relation: c.rel}}, qs.OnConflict.Constraint)
		} else if con.Type != ConstraintTypePKey && con.Type != ConstraintTypeUnique {
			return c.dbError(dbError{Code: errOnConflictConstraintNotUnique,
				Rel: relInfo{Relation: c.rel}}, qs.OnConflict.Constraint)
		}

		target := new(ConflictConstraint)
		target.Name = con.Name
		info.Target = target
	}

	// check that each specified column is present in the target table
	if qs.OnConflict.Update != nil {
		if qs.OnConflict.Update.All {
			for _, w := range c.res.Writes {
				info.Update = append(info.Update, w.Column)
			}

		} else {

			for _, cid := range qs.OnConflict.Update.Items {
				col := findRelColumn(c.rel, cid.Name)
				if col == nil {
					return c.dbError(dbError{Code: errColumnUnknown,
						Rel: relInfo{Relation: c.rel}, Col: colInfo{Id: cid}},
						qs.OnConflict.Update)
				}
				info.Update = append(info.Update, col)
			}
		}
	}

	c.res.Conflict = info
	return nil
}

// typeCheckFieldRead checks if a value from the column that is associated
// with the given field can be read into that field. If strict=false and there
// is no column associated with the given field the check will be skipped.
// If the check is successful a new instance of *FieldRead will be appended
// to the reads field of the *TargetInfo result.
//
// CHECKLIST:
//
//		✅ The field's type MUST NOT be a non-empty interface type.
//		✅ The field's type MAY be a non-interface type that implements sql.Scanner.
//
//	  - If the column's type is json or jsonb:
//	    ✅ The field's type MUST NOT be a chan, func, or a complex type IF it does
//	    not implement the json.Unmarshaler interface,
//	    ✅ otherwise the field's type MAY be any other type.
//
//	    ✅ The field's type MAY be the empty-interface type.
//
//	  - If the column's type is xml:
//	    ✅ The field's type MUST NOT be a func, chan, or map type IF it does not
//	    implement the xml.Unmarshaler interface,
//	    ✅ otherwise the field's type MAY be any other type.
//
//	    ✅ The field's type MUST be a type that, together with the column's type,
//	    has an entry in the compatibility table.
func typeCheckFieldRead(c *checker, f *analysis.FieldInfo, strict bool) error {
	col := findRelColumn(c.rel, f.ColIdent.Name)
	if col == nil && strict {
		return c.dbError(dbError{Code: errColumnUnknown, Col: colInfo{Id: f.ColIdent},
			Rel: relInfo{Relation: c.rel}}, f)
	} else if col == nil && !strict {
		return nil
	}

	check := func(f *analysis.FieldInfo, col *Column) (scanner string, ecode dbErrorCode) {
		// non-empty interface, reject
		if f.Type.Kind == analysis.TypeKindInterface && !f.Type.IsEmptyInterface {
			return "", errColumnFieldTypeRead
		}

		// implements sql.Scanner & non-interface, accept as is
		if f.Type.ImplementsScanner() {
			return "", 0
		}

		if col.Type.OID == oid.JSON || col.Type.OID == oid.JSONB {
			if !f.Type.ImplementsJSONUnmarshaler() {
				// chan, func, or complex, reject
				if f.Type.IsJSONIllegal() {
					return "", errColumnFieldTypeRead
				}

				// []byte type, accept as is
				if f.Type.IsSliceKind(analysis.TypeKindByte) {
					return "", 0
				}

				// string kind, accept as is
				if f.Type.Is(analysis.TypeKindString) {
					return "", 0
				}
			}

			// everything else, accept with JSON
			return "JSON", 0
		}

		// empty interface, accept with AnyToEmptyInterface
		if f.Type.IsEmptyInterface {
			return "AnyToEmptyInterface", 0
		}

		if col.Type.OID == oid.XML {
			if !f.Type.ImplementsXMLUnmarshaler() {
				if f.Type.IsXMLIllegal() {
					return "", errColumnFieldTypeRead
				}

				// []byte type, accept as is
				if f.Type.IsSliceKind(analysis.TypeKindByte) {
					return "", 0
				}

				// string kind, accept as is
				if f.Type.Is(analysis.TypeKindString) {
					return "", 0
				}
			}

			// everything else, accept with XML
			return "XML", 0
		}

		typmod1 := isLength1Type(c, col)
		if comp := typeCompatibility(c, col.Type, f.Type, typmod1); comp != nil {
			return comp.scanner, 0
		}

		// TODO(mkopriva): ...
		// if col.Type.is(oid.Circle, oid.CircleArr) {
		// 	// ...
		// } else if col.Type.is(oid.Interval, oid.IntervalArr) {
		// 	// ...
		// } else if col.Type.Type == TypeTypeDomain {
		// 	// ...
		// } else if col.Type.Type == TypeTypeComposite {
		// 	// ...
		// }
		return "", errColumnFieldTypeRead
	}

	scanner, ecode := check(f, col)
	if ecode > 0 {
		return c.dbError(dbError{Code: ecode, Rel: relInfo{Relation: c.rel},
			Col: colInfo{Id: f.ColIdent, Column: col}}, f) // TODO test
	}

	read := new(FieldRead)
	read.Field = f
	read.Column = col
	read.ColIdent = f.ColIdent
	read.Scanner = scanner

	// NOTE(mkopriva): currently reading is allowed ONLY from the target
	// relation therefore here the alias of the target relation is used.
	// Once it is allowed to read from other, joined relations this will
	// need to be updated to properly handle that scenario.
	read.ColIdent.Qualifier = c.rid.Alias

	if !skipFieldRead(c, read) {
		c.res.Reads = append(c.res.Reads, read)
	}
	return nil
}

// typeCheckFieldWrite checks if a value of the given field can be written to
// the column that is associated with that field. If strict=false and there
// is no column associated with the given field the check will be skipped.
// If the check is successful a new instance of *columnWrite will be appended
// to the writes field of the *tagetInfo instance.
//
// CHECKLIST:
//
//		✅ If the given field has the "default" tag, the target column MUST have
//		   a default data value assigned.
//		✅ The field's type MAY implement the driver.Valuer interface.
//
//	  - If the column's type is json or jsonb:
//	    ✅ The field's type MUST NOT be a chan, func, or a complex type IF it does
//	    not implement the json.Marshaler interface,
//	    ✅ otherwise the field's type MAY be any other type.
//
//	  - If the column's type is xml:
//	    ✅ The field's type MUST NOT be a func, chan, or map type IF it does not
//	    implement the xml.Marshaler interface,
//	    ✅ otherwise the field's type MAY be any other type.
//
//	    ✅ The field's type MUST be a type that, together with the column's type,
//	    has an entry in the compatibility table.
func typeCheckFieldWrite(c *checker, f *analysis.FieldInfo) (*FieldWrite, error) {
	col := findRelColumn(c.rel, f.ColIdent.Name)
	if col == nil {
		return nil, nil
	}

	check := func(f *analysis.FieldInfo, col *Column) (valuer string, ecode dbErrorCode) {
		// default requested but non available, reject
		if f.UseDefault && !col.HasDefault {
			return "", errColumnDefaultUnset
		}

		// implements driver.Valuer, accept as is
		if f.Type.ImplementsValuer() {
			return "", 0
		}

		if col.Type.OID == oid.JSON || col.Type.OID == oid.JSONB {
			if !f.Type.ImplementsJSONMarshaler() {
				// chan, func, or complex, reject
				if f.Type.IsJSONIllegal() {
					return "", errColumnFieldTypeWrite
				}

				// []byte type, accept as is
				if f.Type.IsSliceKind(analysis.TypeKindByte) {
					return "", 0
				}

				// string kind, accept as is
				if f.Type.Is(analysis.TypeKindString) {
					return "", 0
				}
			}

			// everything else, accept with JSON
			return "JSON", 0
		} else if col.Type.OID == oid.XML {
			if !f.Type.ImplementsXMLMarshaler() {
				if f.Type.IsXMLIllegal() {
					return "", errColumnFieldTypeWrite
				}

				// []byte type, accept as is
				if f.Type.IsSliceKind(analysis.TypeKindByte) {
					return "", 0
				}

				// string kind, accept as is
				if f.Type.Is(analysis.TypeKindString) {
					return "", 0
				}
			}

			// everything else, accept with XML
			return "XML", 0
		}

		typmod1 := isLength1Type(c, col)
		if comp := typeCompatibility(c, col.Type, f.Type, typmod1); comp != nil {
			return comp.valuer, 0
		}

		// TODO(mkopriva): ...
		// if col.Type.is(oid.Circle, oid.CircleArr) {
		// 	// ...
		// } else if col.Type.is(oid.Interval, oid.IntervalArr) {
		// 	// ...
		// } else if col.Type.Type == TypeTypeDomain {
		// 	// ...
		// } else if col.Type.Type == TypeTypeComposite {
		// 	// ...
		// }
		return "", errColumnFieldTypeWrite

	}

	valuer, ecode := check(f, col)
	if ecode > 0 {
		return nil, c.dbError(dbError{Code: ecode, Rel: relInfo{Relation: c.rel},
			Col: colInfo{Id: f.ColIdent, Column: col}}, f) // TODO test
	}

	w := new(FieldWrite)
	w.Field = f
	w.Column = col
	w.ColIdent = f.ColIdent
	w.ColIdent.Qualifier = c.rid.Alias
	w.Valuer = valuer
	return w, nil
}

// typeCheckFieldFilter
//
// CHECKLIST:
//
//		✅ If strict=true; the column denoted by the given field MUST be present
//		   in the target relation.
//		✅ The field's type MAY implement the driver.Valuer interface.
//
//	  - If the column's type is json or jsonb:
//	    ✅ The field's type MUST NOT be a chan, func, or a complex type IF it does
//	    not implement the json.Marshaler interface,
//	    ✅ otherwise the field's type MAY be any other type.
//
//	  - If the column's type is xml:
//	    ✅ The field's type MUST NOT be a func, chan, or map type IF it does not
//	    implement the xml.Marshaler interface,
//	    ✅ otherwise the field's type MAY be any other type.
//
//	    ✅ The field's type MUST be a type that, together with the column's type,
//	    has an entry in the compatibility table.
func typeCheckFieldFilter(c *checker, f *analysis.FieldInfo, strict bool) error {
	if f.ExcludeFromFilter() {
		return nil
	}

	col := findRelColumn(c.rel, f.ColIdent.Name)
	if col == nil && strict {
		// NOTE(mkopriva): at the moment strict is always false so this
		// branch will not be executed, nevertheless we leave it here in
		// case strict=true is used sometime later.
		return c.dbError(dbError{Code: errColumnUnknown,
			Col: colInfo{Id: f.ColIdent}, Rel: relInfo{Relation: c.rel}}, f)
	} else if col == nil && !strict {
		return nil
	}

	check := func(f *analysis.FieldInfo, col *Column) (valuer string, ecode dbErrorCode) {
		// implements driver.Valuer? accept as is
		if f.Type.ImplementsValuer() {
			return "", 0
		}

		if col.Type.OID == oid.JSON || col.Type.OID == oid.JSONB {
			if !f.Type.ImplementsJSONMarshaler() {
				// chan, func, or complex? reject
				if f.Type.IsJSONIllegal() {
					return "", errColumnFieldTypeWrite
				}

				// []byte type? accept as is
				if f.Type.IsSliceKind(analysis.TypeKindByte) {
					return "", 0
				}

				// string kind? accept as is
				if f.Type.Is(analysis.TypeKindString) {
					return "", 0
				}
			}

			// everything else, accept with JSON
			return "JSON", 0
		} else if col.Type.OID == oid.XML {
			if !f.Type.ImplementsXMLMarshaler() {
				// chan, func, or map, reject
				if f.Type.IsXMLIllegal() {
					return "", errColumnFieldTypeWrite
				}

				// []byte type, accept as is
				if f.Type.IsSliceKind(analysis.TypeKindByte) {
					return "", 0
				}

				// string kind, accept as is
				if f.Type.Is(analysis.TypeKindString) {
					return "", 0
				}
			}

			// everything else, accept with XML
			return "XML", 0
		}

		typmod1 := isLength1Type(c, col)
		if comp := typeCompatibility(c, col.Type, f.Type, typmod1); comp != nil {
			return comp.valuer, 0
		}

		// TODO(mkopriva): ...
		// if col.Type.is(oid.Circle, oid.CircleArr) {
		// 	// ...
		// } else if col.Type.is(oid.Interval, oid.IntervalArr) {
		// 	// ...
		// } else if col.Type.Type == TypeTypeDomain {
		// 	// ...
		// } else if col.Type.Type == TypeTypeComposite {
		// 	// ...
		// }
		return "", errColumnFieldTypeWrite
	}

	valuer, ecode := check(f, col)
	if ecode > 0 {
		return c.dbError(dbError{Code: ecode, Col: colInfo{Column: col}, Rel: relInfo{Relation: c.rel}}, f)
	}

	filter := new(FieldFilter)
	filter.ColIdent = f.ColIdent
	filter.ColIdent.Qualifier = c.rid.Alias
	filter.Column = col
	filter.Field = f
	filter.Valuer = valuer

	c.res.Filters = append(c.res.Filters, filter)
	return nil
}

func checkTypeCoercion(c *checker, target oid.OID, source oid.OID) bool {
	// if same type, accept
	if target == source {
		return true
	}

	// if target is ANY, accept
	if target == oid.Any {
		return true
	}

	// if source is an untyped string constant assume it can be converted to anything
	if source == oid.Unknown {
		return true
	}

	// if pg_cast says ok, accept.
	key := CastKey{Target: target, Source: source}
	if cast := c.db.catalog.Casts[key]; cast != nil {
		return cast.Context == CastContextImplicit || cast.Context == CastContextAssignment
	}
	return false
}

////////////////////////////////////////////////////////////////////////////////
// Find Functions
//

// findColumn finds and returns the *Column identified by the given ColIdent.
// If the relation denoted by the column's qualifier doesn't exist, or if the
// relation exists but the column itself is not present in that relation an
// error code will be returned instead.
func findColumn(c *checker, id analysis.ColIdent) (*Column, dbErrorCode) {
	if rel, ok := c.relMap[id.Qualifier]; ok {
		for _, col := range rel.Columns {
			if col.Name == id.Name {
				return col, 0 // found
			}
		}
		return nil, errColumnUnknown
	}
	// NOTE(mkopriva): this should never happend, it ought be caught by analysis.
	return nil, errColumnQualifierUnknown
}

// findRelColumn finds and returns the *Column identified by the given name.
// If no Column is found in the provided Relation, nil will be returned instead.
func findRelColumn(rel *Relation, name string) *Column {
	for _, col := range rel.Columns {
		if col.Name == name {
			return col
		}
	}
	return nil
}

// findRelIndex finds and returns the *Index identified by the given name.
// If no Index is found in the provided Relation, nil will be returned instead.
func findRelIndex(rel *Relation, name string) *Index {
	for _, ind := range rel.Indexes {
		if ind.Name == name {
			return ind
		}
	}
	return nil
}

// findRelConstraint finds and returns the *Constraint identified by the given name.
// If no Constraint is found in the provided Relation, nil will be returned instead.
func findRelConstraint(rel *Relation, name string) *Constraint {
	for _, con := range rel.Constraints {
		if con.Name == name {
			return con
		}
	}
	return nil
}

// findQueryColumnField finds and returns the *analysis.FieldInfo of the target data type's
// field that is tagged with the column identified by the given ColIdent. If strict
// is true, findQueryColumnField will also check whether or not the column actually exists.
//
// NOTE(mkopriva): currently the findQueryColumnField method returns fields matched
// by just the column's name, i.e. the qualifiers are ignored, this means that
// one could pass in two different cids with the same name and the method
// would return the same field.
func findQueryColumnField(c *checker, qs *analysis.QueryStruct, cid analysis.ColIdent, strict bool) (*analysis.FieldInfo, dbErrorCode) {
	for _, f := range qs.Rel.Type.Fields {
		if f.ColIdent.Name == cid.Name {
			if strict {
				if _, ecode := findColumn(c, cid); ecode > 0 {
					return nil, ecode
				}
			}
			return f, 0
		}
	}

	// NOTE(mkopriva): this should be caught by anlysis.
	panic("unreachable")
	return nil, 0
}

// //////////////////////////////////////////////////////////////////////////////
// Load Functions
func loadCatalog(db *DB, key string) (*Catalog, error) {
	catalogCache.Lock()
	defer catalogCache.Unlock()
	if cat, ok := catalogCache.m[key]; ok && cat != nil {
		return cat, nil
	}

	cat := new(Catalog)
	cat.Types = make(map[oid.OID]*Type)
	cat.Operators = make(map[OpKey]*Operator)
	cat.Casts = make(map[CastKey]*Cast)
	cat.Procs = make(map[string][]*Proc)
	cat.Relations = make(map[analysis.RelIdent]*Relation)

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
	WHERE pg_catalog.pg_type_is_visible(t.oid)
	AND t.typcategory <> 'P'` //`
	rows, err := db.Query(selectTypes)
	if err != nil {
		return nil, db.dbError(dbError{Code: errCatalogTypeGet, Err: err})
	}
	defer rows.Close()

	for rows.Next() {
		typ := new(Type)
		err := rows.Scan(
			&typ.OID,
			&typ.Name,
			&typ.NameFmt,
			&typ.Length,
			&typ.Type,
			&typ.Category,
			&typ.IsPreferred,
			&typ.Elem,
		)
		if err != nil {
			return nil, db.dbError(dbError{Code: errCatalogTypeScan, Err: err})
		}
		cat.Types[typ.OID] = typ
	}
	if err := rows.Err(); err != nil {
		return nil, db.dbError(dbError{Code: errCatalogTypeGet, Err: err})
	}

	const selectOperators = `SELECT
		o.oid
		, o.oprname
		, o.oprkind
		, o.oprleft
		, o.oprright
		, o.oprresult
	FROM pg_operator o ` //`
	rows, err = db.Query(selectOperators)
	if err != nil {
		return nil, db.dbError(dbError{Code: errCatalogOperatorGet, Err: err})
	}
	defer rows.Close()

	for rows.Next() {
		op := new(Operator)
		err := rows.Scan(
			&op.OID,
			&op.Name,
			&op.Kind,
			&op.Left,
			&op.Right,
			&op.Result,
		)
		if err != nil {
			return nil, db.dbError(dbError{Code: errCatalogOperatorScan, Err: err})
		}

		key := OpKey{Name: op.Name, Left: op.Left, Right: op.Right}
		cat.Operators[key] = op
	}
	if err := rows.Err(); err != nil {
		return nil, db.dbError(dbError{Code: errCatalogOperatorGet, Err: err})
	}

	const selectCasts = `SELECT
		c.oid
		, c.castsource
		, c.casttarget
		, c.castcontext
	FROM pg_cast c ` //`
	rows, err = db.Query(selectCasts)
	if err != nil {
		return nil, db.dbError(dbError{Code: errCatalogCastGet, Err: err})
	}
	defer rows.Close()

	for rows.Next() {
		cast := new(Cast)
		err := rows.Scan(
			&cast.OID,
			&cast.Source,
			&cast.Target,
			&cast.Context,
		)
		if err != nil {
			return nil, db.dbError(dbError{Code: errCatalogCastScan, Err: err})
		}

		key := CastKey{Target: cast.Target, Source: cast.Source}
		cat.Casts[key] = cast
	}
	if err := rows.Err(); err != nil {
		return nil, db.dbError(dbError{Code: errCatalogCastGet, Err: err})
	}

	var selectProcs string
	if db.version >= 110000 {
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
	rows, err = db.Query(selectProcs)
	if err != nil {
		return nil, db.dbError(dbError{Code: errCatalogProcedureGet, Err: err})
	}
	defer rows.Close()

	for rows.Next() {
		proc := new(Proc)
		err := rows.Scan(
			&proc.OID,
			&proc.Name,
			&proc.ArgType,
			&proc.RetType,
			&proc.IsAgg,
		)
		if err != nil {
			return nil, db.dbError(dbError{Code: errCatalogProcedureScan, Err: err})
		}

		cat.Procs[proc.Name] = append(cat.Procs[proc.Name], proc)
	}
	if err := rows.Err(); err != nil {
		return nil, db.dbError(dbError{Code: errCatalogProcedureGet, Err: err})
	}

	catalogCache.m[key] = cat
	return cat, nil
}

func loadTargetRelation(c *checker, f *analysis.RelField) error {
	rid := f.Id
	rel, err := loadRelation(c, c.db, rid, f)
	if err != nil {
		return err
	}

	// Map the "" (empty string) key to the target relation, this will allow
	// columns, constraints, and indexes that were specified without a qualifier
	// to be associated with this target relation.
	c.relMap[""] = rel
	c.relMap[relIdentKey(rid)] = rel

	c.rel = rel
	c.rid = rid
	return nil
}

func loadJoinRelation(c *checker, rid analysis.RelIdent, ptr analysis.FieldPtr) (rel *Relation, err error) {
	if rel, err = loadRelation(c, c.db, rid, ptr); err != nil {
		return nil, err
	}
	c.relMap[relIdentKey(rid)] = rel
	return rel, nil
}

func loadRelation(c *checker, db *DB, rid analysis.RelIdent, ptr analysis.FieldPtr) (*Relation, error) {
	db.catalog.Lock()
	defer db.catalog.Unlock()
	if rel, ok := db.catalog.Relations[rid]; ok && rel != nil {
		return rel, nil
	}

	rel := new(Relation)
	const selectRelationInfo = `SELECT
		c.oid
		, c.relname
		, c.relkind
		, ns.nspname
	FROM pg_class c
	LEFT JOIN pg_namespace ns ON ns.oid = c.relnamespace
	WHERE c.oid = to_regclass($1)` //`
	row := db.QueryRow(selectRelationInfo, rid.QualifiedName())
	if err := row.Scan(&rel.OID, &rel.Name, &rel.RelKind, &rel.Schema); err != nil {
		if err == sql.ErrNoRows {
			return nil, c.dbError(dbError{Code: errRelationUnknown,
				Rel: relInfo{Id: rid}}, ptr)
		}
		return nil, c.dbError(dbError{Code: errRelationScan,
			Rel: relInfo{Id: rid}, Err: err}, ptr)
	}

	const selectRelationColumns = `SELECT
		a.attnum
		, a.attname
		, a.atttypmod
		, a.attnotnull
		, a.atthasdef
		, COALESCE(i.indisprimary, false)
		, a.attndims
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
	rows, err := db.Query(selectRelationColumns, rel.OID)
	if err != nil {
		return nil, c.dbError(dbError{Code: errRelationColumnGet,
			Rel: relInfo{Id: rid, Relation: rel}, Err: err}, ptr)
	}
	defer rows.Close()

	for rows.Next() {
		col := new(Column)
		err := rows.Scan(
			&col.Num,
			&col.Name,
			&col.TypeMod,
			&col.HasNotNull,
			&col.HasDefault,
			&col.IsPrimary,
			&col.NumDims,
			&col.TypeOID,
		)
		if err != nil {
			return nil, c.dbError(dbError{Code: errRelationColumnScan,
				Rel: relInfo{Id: rid, Relation: rel}, Err: err}, ptr)
		}

		if typ, ok := db.catalog.Types[col.TypeOID]; !ok {
			return nil, c.dbError(dbError{Code: errRelationColumnUnknownType,
				Rel: relInfo{Id: rid, Relation: rel}, Col: colInfo{Column: col}}, ptr)
		} else {
			col.Type = typ
		}

		rel.Columns = append(rel.Columns, col)
		col.Relation = rel
	}
	if err := rows.Err(); err != nil {
		return nil, c.dbError(dbError{Code: errRelationColumnGet,
			Rel: relInfo{Id: rid, Relation: rel}, Err: err}, ptr)
	}

	const selectRelationConstraints = `SELECT
		c.oid
		, c.conname
		, c.contype
		, c.condeferrable
		, c.condeferred
		, c.conkey
		, c.confkey
	FROM pg_constraint c
	LEFT JOIN pg_index i ON i.indexrelid = c.conindid
	WHERE c.conrelid = $1
	ORDER BY c.oid` //`
	rows, err = db.Query(selectRelationConstraints, rel.OID)
	if err != nil {
		return nil, c.dbError(dbError{Code: errRelationConstraintGet,
			Rel: relInfo{Id: rid, Relation: rel}, Err: err}, ptr)
	}
	defer rows.Close()

	for rows.Next() {
		con := new(Constraint)
		err := rows.Scan(
			&con.OID,
			&con.Name,
			&con.Type,
			&con.IsDeferrable,
			&con.IsDeferred,
			(*pq.Int64Array)(&con.Key),
			(*pq.Int64Array)(&con.FKey),
		)
		if err != nil {
			return nil, c.dbError(dbError{Code: errRelationConstraintScan,
				Rel: relInfo{Id: rid, Relation: rel}, Err: err}, ptr)
		}
		rel.Constraints = append(rel.Constraints, con)
	}
	if err := rows.Err(); err != nil {
		return nil, c.dbError(dbError{Code: errRelationConstraintGet,
			Rel: relInfo{Id: rid, Relation: rel}, Err: err}, ptr)
	}

	const selectRelationIndexes = `SELECT
		c.oid
		, c.relname
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
	rows, err = db.Query(selectRelationIndexes, rel.OID)
	if err != nil {
		return nil, c.dbError(dbError{Code: errRelationIndexGet,
			Rel: relInfo{Id: rid, Relation: rel}, Err: err}, ptr)
	}
	defer rows.Close()

	for rows.Next() {
		ind := new(Index)
		err := rows.Scan(
			&ind.OID,
			&ind.Name,
			&ind.NumAtts,
			&ind.IsUnique,
			&ind.IsPrimary,
			&ind.IsExclusion,
			&ind.IsImmediate,
			&ind.IsReady,
			(*int2vec)(&ind.Key),
			&ind.Definition,
			&ind.Predicate,
		)
		if err != nil {
			return nil, c.dbError(dbError{Code: errRelationIndexScan,
				Rel: relInfo{Id: rid, Relation: rel}, Err: err}, ptr)
		}

		ind.Expression = parseIndexExpr(ind.Definition)
		rel.Indexes = append(rel.Indexes, ind)
	}
	if err := rows.Err(); err != nil {
		return nil, c.dbError(dbError{Code: errRelationIndexGet,
			Rel: relInfo{Id: rid, Relation: rel}, Err: err}, ptr)
	}

	db.catalog.Relations[rid] = rel
	return rel, nil
}

////////////////////////////////////////////////////////////////////////////////
// Helper Functions
//

// skipFieldInsert reports whether or not the given FieldWrite should be skipped for INSERTs.
func skipFieldInsert(c *checker, fw *FieldWrite) bool {
	return !fw.Field.Mode.CanInsert() && (c.force == nil || !c.force.Contains(fw.ColIdent))
}

// skipFieldUpdate reports whether or not the given FieldWrite should be skipped for UPDATEs.
func skipFieldUpdate(c *checker, fw *FieldWrite) bool {
	return !fw.Field.Mode.CanUpdate() && (c.force == nil || !c.force.Contains(fw.ColIdent))
}

// skipFieldRead reports whether or not the given FieldRead should be skipped.
func skipFieldRead(c *checker, fr *FieldRead) bool {
	return fr.Field.IsWriteOnly() && (c.force == nil || !c.force.Contains(fr.ColIdent))
}

// typeCompatibility
func typeCompatibility(c *checker, ctyp *Type, ftyp analysis.TypeInfo, typmod1 bool) *compentry {
	lit := ftyp.GenericLiteral()
	key := compkey{oid: ctyp.OID, typmod1: typmod1}
	if ctyp.Type == TypeTypeEnum { // enum labels are basically text
		key.oid = oid.Text
	}

	// type table entry exists, accept
	if comp, ok := compatibility.oid2literal[key][lit]; ok {
		return &comp
	}

	// try to salvage this
	if ctyp.Category == TypeCategoryString || ctyp.Category == TypeCategoryArray {
		if ctyp.Category == TypeCategoryString {
			key.oid = oid.Text
		} else if ctyp.Category == TypeCategoryArray {
			if et := c.db.catalog.Types[ctyp.Elem]; et != nil && et.Category == TypeCategoryString {
				key.oid = oid.TextArr
			}
		}
		if comp, ok := compatibility.oid2literal[key][lit]; ok {
			return &comp
		}
	}
	return nil // not compatible
}

// typeOfLiteral returns the type of the given literal expression.
func typeOfLiteral(c *checker, expr string) (*Type, dbErrorCode) {
	const pgselectexprtype = `SELECT id::oid FROM pg_typeof(%s) AS id` //`

	var typoid oid.OID
	row := c.db.QueryRow(fmt.Sprintf(pgselectexprtype, expr))
	if err := row.Scan(&typoid); err != nil {
		return nil, errPredicateLiteralExpr
	}
	return c.db.catalog.Types[typoid], 0
}

// isLength1Type reports whether or not the given column's type
// is a "length 1" type, i.e. char(1), varchar(1), or bit(1)[], etc.
func isLength1Type(c *checker, col *Column) bool {
	typ := col.Type
	if typ.Category == TypeCategoryArray {
		typ = c.db.catalog.Types[typ.Elem]
	}

	if typ.Category == TypeCategoryBitstring {
		return (col.TypeMod == 1)
	} else if typ.Category == TypeCategoryString {
		return ((col.TypeMod - 4) == 1)
	}
	return false
}

func relIdentKey(rid analysis.RelIdent) string {
	if len(rid.Alias) > 0 {
		return rid.Alias
	}
	return rid.Name
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

type int2vec []int16 // helper type

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

var catalogCache = struct {
	sync.RWMutex
	m map[string]*Catalog
}{m: make(map[string]*Catalog)}
