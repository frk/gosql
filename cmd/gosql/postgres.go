package main

// TODO the postgres types circle(arr) and interval(arr) need a corresponding go type
// TODO "cancast" per-field tag option, as well as a global option, as well as a list-of-castable-types option

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/frk/gosql/internal/errors"

	"github.com/lib/pq"
)

const (
	pgselectdbname = `SELECT current_database()` //`

	pgselectdbrelation = `SELECT
		c.oid
		, c.relkind
	FROM pg_class c
	WHERE c.relname = $1
	AND c.relnamespace = to_regnamespace($2)`  //`

	pgselectdbcolumns = `SELECT
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
	ORDER BY a.attnum`  //`

	pgselectdbconstraints = `SELECT
		c.conname
		, c.contype
		, c.condeferrable
		, c.condeferred
		, c.conkey
		, c.confkey
	FROM pg_constraint c
	LEFT JOIN pg_index i ON i.indexrelid = c.conindid
	WHERE c.conrelid = $1
	ORDER BY c.oid`  //`

	pgselectdbindexes = `SELECT
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
	ORDER BY i.indexrelid`  //`

	pgselecttypes = `SELECT
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
	AND t.typcategory <> 'P'`  //`

	pgselectoperators = `SELECT
		o.oid
		, o.oprname
		, o.oprkind
		, o.oprleft
		, o.oprright
		, o.oprresult
	FROM pg_operator o `  //`

	pgselectcasts = `SELECT
		c.oid
		, c.castsource
		, c.casttarget
		, c.castcontext
	FROM pg_cast c `  //`

	pgselectprocs_11plus = `SELECT
		p.oid
		, p.proname
		, p.proargtypes[0]
		, p.prorettype
		, p.prokind = 'a'
	FROM pg_proc p
	WHERE p.pronargs = 1
	AND p.proname NOT LIKE 'pg_%'
	AND p.proname NOT LIKE '_pg_%'
	`  //`

	pgselectprocs_pre11 = `SELECT
		p.oid
		, p.proname
		, p.proargtypes[0]
		, p.prorettype
		, p.proisagg
	FROM pg_proc p
	WHERE p.pronargs = 1
	AND p.proname NOT LIKE 'pg_%'
	AND p.proname NOT LIKE '_pg_%'
	`  //`

	pgshowversionnum = `SHOW server_version_num` //`

	pgselectexprtype = `SELECT id::oid FROM pg_typeof(%s) AS id` //`

)

type postgres struct {
	db   *sql.DB
	cat  *pgcatalogue
	url  string
	name string // name of the current database, used mainly for error reporting
}

func (pg *postgres) init() (err error) {
	if pg.db, err = sql.Open("postgres", pg.url); err != nil {
		return err
	} else if err := pg.db.Ping(); err != nil {
		return err
	}

	// the name of the current database
	if err := pg.db.QueryRow(pgselectdbname).Scan(&pg.name); err != nil {
		return err
	}

	pg.cat = new(pgcatalogue)
	if err = pg.cat.load(pg.db, pg.url); err != nil {
		return err
	}
	return nil
}

func (pg *postgres) close() error {
	return pg.db.Close()
}

type pgchecker struct {
	pg *postgres
	ti *targetInfo

	rel      *pgrelation   // the target relation
	joinlist []*pgrelation // joined relations
	relmap   map[string]*pgrelation
}

func (c *pgchecker) run() (err error) {
	c.relmap = make(map[string]*pgrelation)
	c.ti.searchConditionFieldColumns = make(map[*searchConditionField]*pgcolumn)
	if c.rel, err = c.loadRelation(c.ti.dataField.relId); err != nil {
		return err
	}

	// Map the target relation to the "" (empty string) key, this will allow
	// columns, constraints, and indexes that were specified without a qualifier
	// to be associated with this target relation.
	c.relmap[""] = c.rel

	if c.ti.query != nil {
		return c.checkQueryStruct()
	}
	if c.ti.filter != nil {
		return c.checkFilterStruct()
	}

	panic("nothing to db-check")
	return nil
}

func (c *pgchecker) checkQueryStruct() (err error) {
	if join := c.ti.query.joinBlock; join != nil {
		if err := c.checkJoinBlock(join); err != nil {
			return err
		}
	}
	if onconf := c.ti.query.onConflictBlock; onconf != nil {
		if err := c.checkOnConflictBlock(onconf); err != nil {
			return err
		}
	}
	if where := c.ti.query.whereBlock; where != nil {
		if err := c.checkWhereBlock(where); err != nil {
			return err
		}
	}

	// If an OrderBy directive was used, make sure that the specified
	// columns are present in the loaded relations.
	if c.ti.query.orderByList != nil {
		for _, item := range c.ti.query.orderByList.items {
			if _, err := c.findColumnByColId(item.colId); err != nil {
				return err
			}
		}
	}

	// If a Default directive was provided, make sure that the specified
	// columns are present in the target relation.
	if c.ti.query.defaultList != nil {
		for _, item := range c.ti.query.defaultList.items {
			// Qualifier, if present, must match the target table's alias.
			if len(item.qual) > 0 && item.qual != c.ti.query.dataField.relId.alias {
				return errors.BadTargetTableForDefaultError
			}

			if col := c.rel.findColumn(item.name); col == nil {
				return errors.NoDBColumnError
			} else if !col.hasdefault {
				return errors.NoColumnDefaultSetError
			}
		}
	}

	// If a Force directive was provided, make sure that the specified
	// columns are present in the loaded relations.
	if c.ti.query.forceList != nil {
		for _, item := range c.ti.query.forceList.items {
			if _, err := c.findColumnByColId(item); err != nil {
				return err
			}
		}
	}

	// If a Return directive was provided, make sure that the specified
	// columns are present in the loaded relations.
	var strict bool
	var outputfields []*fieldInfo
	if res := c.ti.query.resultField; res != nil {
		strict = true
		outputfields = res.data.fields
	} else if c.ti.query.returnList != nil {
		if c.ti.query.returnList.all {
			strict = false
			outputfields = c.ti.query.dataField.data.fields
		} else {
			// If an explicit list of columns was provided, make sure that
			// they are present in one of the associated relations, and that
			// each one of them has a corresponding field, if not return an error.
			strict = true

			for _, colId := range c.ti.query.returnList.items {
				// NOTE(mkopriva): currently the findFieldByColId method returns
				// fields matched by just the column's name, i.e. the qualifiers
				// are ignored, this means that one could pass in two different
				// colIds with the same name and the method would return the same field.
				field, err := c.findFieldByColId(colId)
				if err != nil {
					return err
				}
				outputfields = append(outputfields, field)
			}
		}

	}
	if len(outputfields) > 0 {
		if err := c.checkFields(outputfields, dataRead, strict); err != nil {
			return err
		}
	}

	if dataField := c.ti.query.dataField; dataField != nil && !dataField.isDirective {
		var dataOp dataOperation
		if c.ti.query.kind == queryKindSelect {
			dataOp = dataRead
		} else if c.ti.query.kind == queryKindInsert || c.ti.query.kind == queryKindUpdate {
			dataOp = dataWrite
		}

		if err := c.checkFields(dataField.data.fields, dataOp, true); err != nil {
			return err
		}
	}

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

	//func (q *UpdatePKeyReturningAllSingleQuery) Exec(c gosql.Conn) error {
	//   	const queryString = `UPDATE "test_user" SET (
	//   		"email"
	//   		, "password"
	//   		, "created_at"
	//   		, "updated_at"
	//   	) = (
	//   		$1
	//   		, $2
	//   		, $3
	//   		, $4
	//   	)
	//   	WHERE "id" = $5 AND "id" = $6
	//   	RETURNING
	//   	"id"
	//   	, "email"
	//   	, "created_at"
	//   	, "updated_at"` // `

	//   	row := c.QueryRow(queryString,
	//   		q.User.Email,
	//   		q.User.Password,
	//   		q.User.CreatedAt,
	//   		q.User.UpdatedAt,
	//   		q.User.Id,
	//   		q.User.Id,
	//   	)
	//   	return row.Scan(
	//   		&q.User.Id,
	//   		&q.User.Email,
	//   		&q.User.CreatedAt,
	//   		&q.User.UpdatedAt,
	//   	)
	//   }
}

func (c *pgchecker) checkFilterStruct() (err error) {
	// If a TextSearch directive was provided, make sure that the
	// specified column is present in one of the loaded relations
	// and that it has the correct type.
	if c.ti.filter.textSearchColId != nil {
		col, err := c.findColumnByColId(*c.ti.filter.textSearchColId)
		if err != nil {
			return err
		} else if col.typ.oid != pgtyp_tsvector {
			return errors.BadDBColumnTypeError
		}
	}

	if dataField := c.ti.filter.dataField; dataField != nil {
		if err := c.checkFields(dataField.data.fields, dataTest, false); err != nil {
			return err
		}
	}
	return nil
}

// checkFields checks column existence, type compatibility, operation validity
func (c *pgchecker) checkFields(fields []*fieldInfo, dataOp dataOperation, strict bool) (err error) {
	if dataOp == dataNop {
		return nil
	}

	for _, fld := range fields {
		var col *pgcolumn
		// TODO(mkopriva): currenlty checkFields requires that every
		// field has a corresponding column in the target relation,
		// which represents an ideal where each relation in the db
		// has a matching type in the app, this however may not always
		// be practical and there may be cases where a single struct
		// type is used to represent multiple different relations
		// having fields that have corresponding columns in one
		// relation but not in another...
		if dataOp == dataRead {
			// TODO(mkopriva): currently columns specified in
			// the fields of the struct representing the record
			// aren't really meant to include the relation alias
			// which makes this a bit of a non-issue, however in
			// the future it would be good to provide a way to do
			// that, like when selecting columns from multiple
			// joined tables.. therefore this stays here, at least for now...

			// If this is a SELECT, or the target type is
			// from the "Result" field, lookup the column
			// in all of the associated relations since its
			// ok to select columns from joined relations.
			if col, err = c.findColumnByColId(fld.colId); err != nil {
				if strict {
					return err
				}
				continue
			}
		} else {
			// If this is a dataWrite or dataTest operation the column
			// must be present directly in the target relation, columns
			// from associated relations (joins) are not allowed.
			if col = c.rel.findColumn(fld.colId.name); col == nil {
				if strict {
					return errors.NoDBColumnError
				}
				continue
			}
		}

		if dataOp == dataWrite {
			if fld.useDefault && !col.hasdefault {
				// TODO error
			}
		}

		cid := colId{name: fld.colId.name, qual: c.ti.dataField.relId.alias}
		info := &fieldColumnInfo{field: fld, column: col, colId: cid}

		// Make sure that a value of the given field's type
		// can be assigned to given column, and vice versa.
		if !c.canAssign(info, col, fld, dataOp) {
			return errors.BadFieldToColumnTypeError
		}

		if dataOp == dataRead {
			c.ti.output = append(c.ti.output, info)
		} else if dataOp == dataWrite || dataOp == dataTest {
			c.ti.input = append(c.ti.input, info)
		}
		// aggrgate primary keys for writes only
		if col.isprimary && dataOp == dataWrite {
			c.ti.primaryKeys = append(c.ti.primaryKeys, info)
		}
	}
	return nil
}

func (c *pgchecker) checkJoinBlock(jb *joinBlock) error {
	if len(jb.relId.name) > 0 {
		rel, err := c.loadRelation(jb.relId)
		if err != nil {
			return err
		}
		c.joinlist = append(c.joinlist, rel)
	}
	for _, join := range jb.items {
		rel, err := c.loadRelation(join.relId)
		if err != nil {
			return err
		}
		c.joinlist = append(c.joinlist, rel)

		for _, item := range join.conds {
			switch cond := item.cond.(type) {
			case *searchConditionColumn:
				// A join condition's left-hand-side column MUST always
				// reference a column of the relation being joined, so to
				// avoid confusion make sure that node.colId has either no
				// qualifier or, if it has one, it matches the alias of
				// the joined table.
				if cond.colId.qual != "" && cond.colId.qual != join.relId.alias {
					return errors.NoDBRelationError
				}

				// Make sure that colId is present in relation being joined.
				col := rel.findColumn(cond.colId.name)
				if col == nil {
					return errors.NoDBColumnError
				}

				if cond.pred.isUnary() {
					// Column type must be bool if the predicate type produces
					// the "IS [NOT] { FALSE | TRUE | UNKNOWN }" SQL predicate.
					if cond.pred.isBoolean() && col.typ.oid != pgtyp_bool {
						return errors.BadColumnTypeForUnaryOpError
					}
					// Column must be NULLable if the predicate type
					// produces the "IS [NOT] NULL" SQL predicate.
					if col.hasnotnull && (cond.pred == isNull || cond.pred == notNull) {
						return errors.BadColumnNULLSettingForNULLOpError
					}
				} else {
					var typ *pgtype
					// Get the type of the right hand side, which is
					// either a column or a literal expression.
					if len(cond.colId2.name) > 0 {
						colId2, err := c.findColumnByColId(cond.colId2)
						if err != nil {
							return err
						}
						typ = colId2.typ
					} else if len(cond.literal) > 0 {
						var oid pgoid
						row := c.pg.db.QueryRow(fmt.Sprintf(pgselectexprtype, cond.literal))
						if err := row.Scan(&oid); err != nil {
							return errors.BadLiteralExpressionError
						}
						typ = c.pg.cat.types[oid]
					}

					if cond.pred.isQuantified() || cond.qua > 0 {
						// Check that the quantifier can be used
						// with the type of the RHS expression.
						if typ.category != pgtypcategory_array {
							return errors.BadExpressionTypeForQuantifierError
						}
						typ = c.pg.cat.types[typ.elem]
					}

					rhsoids := []pgoid{typ.oid}
					if !c.canCompare(col, rhsoids, cond.pred) {
						return errors.BadColumnToLiteralComparisonError
					}
				}
			default:
				// NOTE(mkopriva): currently predicates other than
				// searchConditionColumn are not supported as a join condition.
			}
		}
	}
	return nil
}

func (c *pgchecker) checkOnConflictBlock(onconf *onConflictBlock) error {
	rel := c.rel

	// If a Column directive was provided in an OnConflict block, make
	// sure that the listed columns are present in the target table.
	// Make also msure that the list of columns matches the full list
	// of columns of a unique index that's present on the target table.
	if len(onconf.column) > 0 {
		var attnums []int16
		for _, id := range onconf.column {
			col := rel.findColumn(id.name)
			if col == nil {
				return errors.NoDBColumnError
			}
			attnums = append(attnums, col.num)
		}

		var match bool
		for _, ind := range rel.indexes {
			if !ind.isunique && !ind.isprimary {
				continue
			}
			if !matchNumbers(ind.key, attnums) {
				continue
			}

			match = true
			c.ti.onConflictIndex = ind
			break
		}
		if !match {
			return errors.NoDBIndexForColumnListError
		}
	}

	// If an Index directive was provided check that the specified
	// index is actually present on the target table and also make
	// sure that it is a unique index.
	if len(onconf.index) > 0 {
		ind := rel.findIndex(onconf.index)
		if ind == nil {
			return errors.NoDBIndexError
		}
		if !ind.isunique && !ind.isprimary {
			return errors.NoDBIndexError
		}

		c.ti.onConflictIndex = ind
	}

	// If a Constraint directive was provided make sure that the
	// specified constraint is present on the target table and that
	// it is a unique constraint.
	if len(onconf.constraint) > 0 {
		con := rel.findConstraint(onconf.constraint)
		if con == nil {
			return errors.NoDBConstraintError
		}
		if con.typ != pgconstraint_pkey && con.typ != pgconstraint_unique {
			return errors.NoDBConstraintError
		}
	}

	// If an Update directive was provided, make sure that each
	// listed column is present in the target table.
	if onconf.update != nil {
		for _, item := range onconf.update.items {
			if col := rel.findColumn(item.name); col == nil {
				return errors.NoDBColumnError
			}
		}
	}
	return nil
}

func (c *pgchecker) checkWhereBlock(where *whereBlock) error {
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
			case *searchConditionColumn:
				// Check that the referenced Column is present
				// in one of the associated relations.
				col, err := c.findColumnByColId(cond.colId)
				if err != nil {
					return err
				}

				if cond.pred.isUnary() {
					// Column type must be bool if the predicate type produces
					// the "IS [NOT] { FALSE | TRUE | UNKNOWN }" SQL predicate.
					if cond.pred.isBoolean() && col.typ.oid != pgtyp_bool {
						return errors.BadColumnTypeForUnaryOpError
					}
					// Column must be NULLable if the predicate type
					// produces the "IS [NOT] NULL" SQL predicate.
					if col.hasnotnull && (cond.pred == isNull || cond.pred == notNull) {
						return errors.BadColumnNULLSettingForNULLOpError
					}
				} else {
					var typ *pgtype

					// Get the type of the right hand side, which is
					// either a column or a literal expression.
					if len(cond.colId2.name) > 0 {
						col2, err := c.findColumnByColId(cond.colId2)
						if err != nil {
							return err
						}
						typ = col2.typ
					} else if len(cond.literal) > 0 {
						var oid pgoid
						row := c.pg.db.QueryRow(fmt.Sprintf(pgselectexprtype, cond.literal))
						if err := row.Scan(&oid); err != nil {
							return errors.BadLiteralExpressionError
						}
						typ = c.pg.cat.types[oid]
					} else {
						panic("shouldn't happen")
					}

					if cond.pred.isQuantified() || cond.qua > 0 {
						// Check that the quantifier can be used
						// with the type of the RHS expression.
						if typ.category != pgtypcategory_array {
							return errors.BadExpressionTypeForQuantifierError
						}
						typ = c.pg.cat.types[typ.elem]
					}

					rhsoids := []pgoid{typ.oid}
					if !c.canCompare(col, rhsoids, cond.pred) {
						return errors.BadColumnToLiteralComparisonError
					}
				}
			case *searchConditionBetween:
				// Check that the referenced Column is present
				// in one of the associated relations.
				col, err := c.findColumnByColId(cond.colId)
				if err != nil {
					return err
				}

				// Check that both predicands, x and y, can be compared to the column.
				for _, arg := range []interface{}{cond.x, cond.y} {
					var argoids []pgoid
					switch a := arg.(type) {
					case colId:
						col2, err := c.findColumnByColId(a)
						if err != nil {
							return err
						}
						argoids = []pgoid{col2.typ.oid}
					case *fieldDatum:
						argoids = c.typeoids(a.typ)
					}

					if !c.canCompare(col, argoids, isGT) {
						return errors.BadColumnToColumnTypeComparisonError
					}
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

// typeoids returns a list of OIDs of those PostgreSQL types that can be
// used to hold a value of a Go type represented by the given typeInfo.
func (c *pgchecker) typeoids(typ typeInfo) []pgoid {
	switch typstr := typ.goTypeId(false, false, true); typstr {
	case goTypeStringMap, goTypeNullStringMap, goTypeStringPtrMap:
		if t := c.pg.cat.typebyname("hstore"); t != nil {
			return []pgoid{t.oid}
		}
	case goTypeStringMapSlice, goTypeNullStringMapSlice, goTypeStringPtrMapSlice:
		if t := c.pg.cat.typebyname("_hstore"); t != nil {
			return []pgoid{t.oid}
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
func (c *pgchecker) canCompare(col *pgcolumn, rhstypes []pgoid, pred predicate) bool {
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

// canAssign reports whether a value
func (c *pgchecker) canAssign(info *fieldColumnInfo, col *pgcolumn, field *fieldInfo, dataOp dataOperation) bool {
	// TODO(mkopriva): this returns on success, so write tests that test
	// successful scenarios so that every code branch is covered.

	// If the column is gonna be written to and the field's type knows
	// how to encode itself to a database value, accept.
	if dataOp == dataWrite && field.typ.isValuer {
		return true
	}

	// If the column is gonna be read from and the field's type knows how
	// to decode itself from a database value, accept.
	if dataOp == dataRead && field.typ.isScanner {
		return true
	}

	oid := col.typ.oid
	gotyp := field.typ.goTypeId(false, false, true)
	typkey := pgsqlTypeKey{oid: oid}

	// Because the pgsql.JSON and pgsql.XML transformers both accept the empty
	// interface as their argument, if the field's type is not considered as
	// raw data (string or []byte types) it can be substituted with the interface{}
	// type for the purpose of resolving the pgsql transformer.
	if oid == pgtyp_json || oid == pgtyp_jsonb || oid == pgtyp_xml {
		if gotyp != goTypeString && gotyp != goTypeByteSlice {
			gotyp = goTypeEmptyInterface
		}

		info.pgsql = pgsqlTypeTable[typkey][gotyp]
		return true
	}

	// Columns with a type in the bit or char family and a typmod of 1 have
	// a distinct Go representation then those with a typmod != 1.
	if col.typ.isbase(pgtyp_bit, pgtyp_varbit, pgtyp_char, pgtyp_varchar, pgtyp_bpchar) {
		typkey.typmod1 = (col.typmod == 1)
	}

	// Columns with type numeric that have a precision but no scale, have
	// a distinct Go representation then those numeric types that have a
	// different precision and scale.
	if col.typ.is(pgtyp_numeric, pgtyp_numericarr) {
		precision := ((col.typmod - 4) >> 16) & 65535
		scale := (col.typmod - 4) & 65535
		if precision > 0 && scale == 0 {
			typkey.noscale = true
		}
	}

	if entry, ok := pgsqlTypeTable[typkey][gotyp]; ok {
		info.pgsql = entry
		return true
	}

	// If the target type is of the string category, accept.
	if col.typ.category == pgtypcategory_string {
		return true
	}

	// If the field's type is []string or [][]byte and the column's type is an
	// array type whose element type belongs to the string category then the oid
	// can be substituted with the text[] type for the purpose of resolving a
	// compatible pgsql transformer.
	if (gotyp == goTypeStringSlice || gotyp == goTypeByteSliceSlice) && col.typ.category == pgtypcategory_array {
		elemtyp := c.pg.cat.types[col.typ.elem]
		if elemtyp != nil && elemtyp.category == pgtypcategory_string {
			return true
		}
	}

	// TODO(mkopriva): implement canAssign for domain types
	if col.typ.typ == pgtyptype_domain {
		return false
	}
	// TODO(mkopriva): implement canAssign for composite types
	if col.typ.typ == pgtyptype_composite {
		return false
	}
	return false
}

func (c *pgchecker) canCoerceOID(targetid, inputid pgoid) bool {
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
func (c *pgchecker) checkModifierFunction(fn funcName, col *pgcolumn, fieldoids []pgoid) error {
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

func (c *pgchecker) loadRelation(id relId) (*pgrelation, error) {
	rel := new(pgrelation)
	rel.name = id.name
	rel.namespace = id.qual
	if len(rel.namespace) == 0 {
		rel.namespace = "public"
	}

	// retrieve relation info
	row := c.pg.db.QueryRow(pgselectdbrelation, rel.name, rel.namespace)
	if err := row.Scan(&rel.oid, &rel.relkind); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NoDBRelationError
		}
		return nil, err
	}

	// retrieve column info
	rows, err := c.pg.db.Query(pgselectdbcolumns, rel.oid)
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

	// load constraints info
	rows, err = c.pg.db.Query(pgselectdbconstraints, rel.oid)
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

	// load index info
	rows, err = c.pg.db.Query(pgselectdbindexes, rel.oid)
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

func (c *pgchecker) findColumnByColId(id colId) (*pgcolumn, error) {
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

func (c *pgchecker) findFieldByColId(id colId) (*fieldInfo, error) {
	// make sure the column actually exists
	if _, err := c.findColumnByColId(id); err != nil {
		return nil, err
	}

	for _, field := range c.ti.query.dataField.data.fields {
		if field.colId.name == id.name {
			return field, nil
		}
	}
	return nil, errors.NoColumnFieldError
}

type pgrelation struct {
	oid         pgoid  // The object identifier of the relation.
	name        string // The name of the relation.
	namespace   string // The name of the namespace to which the relation belongs.
	relkind     string // The relation's kind, we're only interested in r, v, and m.
	columns     []*pgcolumn
	constraints []*pgconstraint
	indexes     []*pgindex
}

func (rel *pgrelation) findColumn(name string) *pgcolumn {
	for _, col := range rel.columns {
		if col.name == name {
			return col
		}
	}
	return nil
}

func (rel *pgrelation) findConstraint(name string) *pgconstraint {
	for _, con := range rel.constraints {
		if con.name == name {
			return con
		}
	}
	return nil
}

func (rel *pgrelation) findIndex(name string) *pgindex {
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
	typoid pgoid
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
	oid pgoid
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
	elem pgoid
}

func (t pgtype) is(oids ...pgoid) bool {
	for _, oid := range oids {
		if t.oid == oid {
			return true
		}
	}
	return false
}

// isbase returns true if t's oid matches one of the given oids, or if t is an
// array type isbase returns true if t's elem matches one of the given oids.
func (t pgtype) isbase(oids ...pgoid) bool {
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
	oid    pgoid
	name   string
	kind   string
	left   pgoid
	right  pgoid
	result pgoid
}

type pgopkey struct {
	name  string
	left  pgoid
	right pgoid
}

type pgcast struct {
	oid     pgoid
	source  pgoid
	target  pgoid
	context string
}

type pgproc struct {
	oid     pgoid
	name    string
	argtype pgoid
	rettype pgoid
	isagg   bool
}

type pgcastkey struct {
	target pgoid
	source pgoid
}

type pgcatalogue struct {
	types     map[pgoid]*pgtype
	operators map[pgopkey]*pgoperator
	casts     map[pgcastkey]*pgcast
	procs     map[funcName][]*pgproc
}

var pgcataloguecache = struct {
	sync.RWMutex
	m map[string]*pgcatalogue
}{m: make(map[string]*pgcatalogue)}

func (c *pgcatalogue) typebyname(name string) *pgtype {
	for _, t := range c.types {
		if t.name == name {
			return t
		}
	}
	return nil
}

func (c *pgcatalogue) typebyoid(oid pgoid) *pgtype {
	return c.types[oid]
}

// cancasti reports whether s can be cast to t *implicitly* or in assignment.
func (c *pgcatalogue) cancasti(t, s pgoid) bool {
	key := pgcastkey{target: t, source: s}
	if cast := c.casts[key]; cast != nil {
		return cast.context == pgcast_implicit || cast.context == pgcast_assignment
	}
	return false
}

func (c *pgcatalogue) load(db *sql.DB, key string) error {
	pgcataloguecache.Lock()
	defer pgcataloguecache.Unlock()

	cat := pgcataloguecache.m[key]
	if cat != nil {
		*c = *cat
		return nil
	}

	c.types = make(map[pgoid]*pgtype)

	// load types
	rows, err := db.Query(pgselecttypes)
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

	c.operators = make(map[pgopkey]*pgoperator)

	// load operators
	rows, err = db.Query(pgselectoperators)
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

	c.casts = make(map[pgcastkey]*pgcast)

	// load casts
	rows, err = db.Query(pgselectcasts)
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

	c.procs = make(map[funcName][]*pgproc)

	var version int
	var pgselectprocs string
	if err := db.QueryRow(pgshowversionnum).Scan(&version); err != nil {
		return err
	} else if version >= 110000 {
		pgselectprocs = pgselectprocs_11plus
	} else {
		pgselectprocs = pgselectprocs_pre11
	}

	// load procs
	rows, err = db.Query(pgselectprocs)
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

	pgcataloguecache.m[key] = c
	return nil
}

type dataOperation uint8

const (
	dataNop dataOperation = iota
	dataRead
	dataWrite
	dataTest
)

type pgoid uint32

func (oid pgoid) getArrayOID() pgoid {
	return pgoidToArrayOID[oid]
}

// postgres types
const (
	pgtyp_any            pgoid = 2276
	pgtyp_bit            pgoid = 1560
	pgtyp_bitarr         pgoid = 1561
	pgtyp_bool           pgoid = 16
	pgtyp_boolarr        pgoid = 1000
	pgtyp_box            pgoid = 603
	pgtyp_boxarr         pgoid = 1020
	pgtyp_bpchar         pgoid = 1042
	pgtyp_bpchararr      pgoid = 1014
	pgtyp_bytea          pgoid = 17
	pgtyp_byteaarr       pgoid = 1001
	pgtyp_char           pgoid = 18
	pgtyp_chararr        pgoid = 1002
	pgtyp_cidr           pgoid = 650
	pgtyp_cidrarr        pgoid = 651
	pgtyp_circle         pgoid = 718
	pgtyp_circlearr      pgoid = 719
	pgtyp_date           pgoid = 1082
	pgtyp_datearr        pgoid = 1182
	pgtyp_daterange      pgoid = 3912
	pgtyp_daterangearr   pgoid = 3913
	pgtyp_float4         pgoid = 700
	pgtyp_float4arr      pgoid = 1021
	pgtyp_float8         pgoid = 701
	pgtyp_float8arr      pgoid = 1022
	pgtyp_inet           pgoid = 869
	pgtyp_inetarr        pgoid = 1041
	pgtyp_int2           pgoid = 21
	pgtyp_int2arr        pgoid = 1005
	pgtyp_int2vector     pgoid = 22
	pgtyp_int2vectorarr  pgoid = 1006
	pgtyp_int4           pgoid = 23
	pgtyp_int4arr        pgoid = 1007
	pgtyp_int4range      pgoid = 3904
	pgtyp_int4rangearr   pgoid = 3905
	pgtyp_int8           pgoid = 20
	pgtyp_int8arr        pgoid = 1016
	pgtyp_int8range      pgoid = 3926
	pgtyp_int8rangearr   pgoid = 3927
	pgtyp_interval       pgoid = 1186
	pgtyp_intervalarr    pgoid = 1187
	pgtyp_json           pgoid = 114
	pgtyp_jsonarr        pgoid = 199
	pgtyp_jsonb          pgoid = 3802
	pgtyp_jsonbarr       pgoid = 3807
	pgtyp_line           pgoid = 628
	pgtyp_linearr        pgoid = 629
	pgtyp_lseg           pgoid = 601
	pgtyp_lsegarr        pgoid = 1018
	pgtyp_macaddr        pgoid = 829
	pgtyp_macaddrarr     pgoid = 1040
	pgtyp_macaddr8       pgoid = 774
	pgtyp_macaddr8arr    pgoid = 775
	pgtyp_money          pgoid = 790
	pgtyp_moneyarr       pgoid = 791
	pgtyp_numeric        pgoid = 1700
	pgtyp_numericarr     pgoid = 1231
	pgtyp_numrange       pgoid = 3906
	pgtyp_numrangearr    pgoid = 3907
	pgtyp_oidvector      pgoid = 30
	pgtyp_path           pgoid = 602
	pgtyp_patharr        pgoid = 1019
	pgtyp_point          pgoid = 600
	pgtyp_pointarr       pgoid = 1017
	pgtyp_polygon        pgoid = 604
	pgtyp_polygonarr     pgoid = 1027
	pgtyp_text           pgoid = 25
	pgtyp_textarr        pgoid = 1009
	pgtyp_time           pgoid = 1083
	pgtyp_timearr        pgoid = 1183
	pgtyp_timestamp      pgoid = 1114
	pgtyp_timestamparr   pgoid = 1115
	pgtyp_timestamptz    pgoid = 1184
	pgtyp_timestamptzarr pgoid = 1185
	pgtyp_timetz         pgoid = 1266
	pgtyp_timetzarr      pgoid = 1270
	pgtyp_tsquery        pgoid = 3615
	pgtyp_tsqueryarr     pgoid = 3645
	pgtyp_tsrange        pgoid = 3908
	pgtyp_tsrangearr     pgoid = 3909
	pgtyp_tstzrange      pgoid = 3910
	pgtyp_tstzrangearr   pgoid = 3911
	pgtyp_tsvector       pgoid = 3614
	pgtyp_tsvectorarr    pgoid = 3643
	pgtyp_uuid           pgoid = 2950
	pgtyp_uuidarr        pgoid = 2951
	pgtyp_unknown        pgoid = 705
	pgtyp_varbit         pgoid = 1562
	pgtyp_varbitarr      pgoid = 1563
	pgtyp_varchar        pgoid = 1043
	pgtyp_varchararr     pgoid = 1015
	pgtyp_xml            pgoid = 142
	pgtyp_xmlarr         pgoid = 143

	pgtyp_hstore    pgoid = 9999
	pgtyp_hstorearr pgoid = 9998
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

var pgoidToArrayOID = map[pgoid]pgoid{
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

var go2pgoids = map[goTypeId][]pgoid{
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
