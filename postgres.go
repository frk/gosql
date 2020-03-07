package gosql

// TODO(mkopriva):
// - given the gosql.Index directive inside an OnConflict block use pg_get_indexdef(index_oid)
//   to retrieve the index's definition, parse that to extract the index expression and then
//   use that expression when generating the ON CONFLICT clause.
//   (https://www.postgresql.org/message-id/204ADCAA-853B-4B5A-A080-4DFA0470B790%40justatheory.com)
// TODO the postgres types circle(arr) and interval(arr) need a corresponding go type
// TODO "cancast" per-field tag option, as well as a global option, as well as a list-of-castable-types option

import (
	"database/sql"
	"fmt"
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
	pg       *postgres
	spec     *typespec
	rel      *pgrelation   // the target relation
	joinlist []*pgrelation // joined relations
	relmap   map[string]*pgrelation

	// result
	info pginfo
}

// fieldcolumn is used to hold a field and its corresponding column.
type fieldcolumn struct {
	field  *recfield
	column *pgcolumn
	colid  colid
}

type pginfo struct {
	input  []*fieldcolumn
	output []*fieldcolumn
}

func pgcheck(pg *postgres, spec *typespec) (info *pginfo, err error) {
	c := pgchecker{pg: pg, spec: spec}
	if err := c.run(); err != nil {
		return nil, err
	}
	return &c.info, nil
}

func (c *pgchecker) run() (err error) {
	c.relmap = make(map[string]*pgrelation)
	if c.rel, err = c.loadrelation(c.spec.rel.relid); err != nil {
		return err
	}

	// Map the target relation to the "" (empty string) key, this will allow
	// columns, constraints, and indexes that were specified without a qualifier
	// to be associated with this target relation.
	c.relmap[""] = c.rel

	if join := c.spec.join; join != nil {
		if err := c.checkjoin(join); err != nil {
			return err
		}
	}
	if onconf := c.spec.onconflict; onconf != nil {
		if err := c.checkonconflict(onconf); err != nil {
			return err
		}
	}
	if where := c.spec.where; where != nil {
		if err := c.checkwhere(where); err != nil {
			return err
		}
	}

	// If an OrderBy directive was used, make sure that the specified
	// columns are present in the loaded relations.
	if c.spec.orderby != nil {
		for _, item := range c.spec.orderby.items {
			if _, err := c.column(item.col); err != nil {
				return err
			}
		}
	}

	// If a Default directive was provided, make sure that the specified
	// columns are present in the target relation.
	if c.spec.defaults != nil {
		for _, item := range c.spec.defaults.items {
			// Qualifier, if present, must match the target table's alias.
			if len(item.qual) > 0 && item.qual != c.spec.rel.relid.alias {
				return errors.BadTargetTableForDefaultError
			}

			if col := c.rel.column(item.name); col == nil {
				return errors.NoDBColumnError
			} else if !col.hasdefault {
				return errors.NoColumnDefaultSetError
			}
		}
	}

	// If a Force directive was provided, make sure that the specified
	// columns are present in the loaded relations.
	if c.spec.force != nil {
		for _, item := range c.spec.force.items {
			if _, err := c.column(item); err != nil {
				return err
			}
		}
	}

	// If a Return directive was provided, make sure that the specified
	// columns are present in the loaded relations.
	if c.spec.returning != nil {
		if c.spec.returning.all {
			// If all is set to true, collect the to-be-returned list
			// of fieldcolumn pairs by going over the recordtype's fields
			// and matching them up with columns from the target relation.
			// Fields that have no matching column in the target relation
			// will be ignored.
			//
			// NOTE(mkopriva): currently if all is set to true only
			// the columns of the target relation will be considered
			// as candidates for the RETURNING clause, other columns
			// from joined relations will be ignored.
			for _, field := range c.spec.rel.rec.fields {
				if col := c.rel.column(field.colid.name); col != nil {
					cid := colid{name: field.colid.name, qual: c.spec.rel.relid.alias}
					pair := &fieldcolumn{field: field, column: col, colid: cid}
					c.info.output = append(c.info.output, pair)
				}
			}
		} else {
			for _, colid := range c.spec.returning.items {
				// If a list of specific columns was provided,
				// make sure that they are present in one of the
				// associated relations, if not return an error.
				col, err := c.column(colid)
				if err != nil {
					return err
				}

				// The to-be-returned columns must have a
				// corresponding field in the recordtype.
				//
				// NOTE(mkopriva): currently the to-be-returned
				// columns are matched to fields using just the
				// column's name, i.e. the qualifiers are ignored,
				// this means that one could specify two separate
				// columns that have the same name and their values
				// would be scanned into the same field.
				var hasfield bool
				for _, field := range c.spec.rel.rec.fields {
					if field.colid.name == colid.name {
						pair := &fieldcolumn{field: field, column: col, colid: colid}
						c.info.output = append(c.info.output, pair)
						hasfield = true
						break
					}
				}
				if !hasfield {
					return errors.NoFieldColumnError
				}
			}
		}
	}

	// If a TextSearch directive was provided, make sure that the
	// specified column is present in one of the loaded relations
	// and that it has the correct type.
	if c.spec.textsearch != nil {
		col, err := c.column(*c.spec.textsearch)
		if err != nil {
			return err
		} else if col.typ.oid != pgtyp_tsvector {
			return errors.BadDBColumnTypeError
		}
	}

	if rel := c.spec.rel; rel != nil && !rel.isdir {
		if err := c.checkfields(rel.rec, false); err != nil {
			return err
		}
	}
	if res := c.spec.result; res != nil {
		if err := c.checkfields(res.rec, true); err != nil {
			return err
		}
	}

	if c.spec.kind == speckindInsert {
		// TODO(mkopriva): if this is an insert make sure that all columns that
		// do not have a DEFAULT set but do have a NOT NULL set, have a corresponding
		// field in the relation's recordtype. (keep in mind that such a column could
		// be also set by a BEFORE TRIGGER, so maybe instead of erroring only warning
		// would be thrown, or make this check optional, something that can be turned
		// on/off...?)
	}

	// TODO(mkopriva): if this is an insert or update (i.e. write) the column
	// generated from the field tags cannot contain duplicate columns. Conversely
	// if this is a read query (select, or I,U,D with returning) the columns do
	// not have to be unique.

	return nil
}

func (c *pgchecker) checkfields(rec recordtype, isresult bool) (err error) {
	for _, fld := range rec.fields {
		var col *pgcolumn
		var atyp assigntype
		// TODO(mkopriva): currenlty this requires that every field
		// has a corresponding column in the target relation, which
		// represents an ideal where each relation in the db has a
		// matching type in the app, this however may not always be
		// practical and there may be cases where a single struct
		// type is used to represent multiple different relations...
		if c.spec.kind == speckindSelect || isresult {
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
			if col, err = c.column(fld.colid); err != nil {
				return err
			}
			atyp = assignread
		} else {
			// If this is a non-select typespec and non-result
			// the column must be present in the target relation.
			if col = c.rel.column(fld.colid.name); col == nil {
				return errors.NoDBColumnError
			}
			atyp = assignwrite

		}

		if c.spec.kind == speckindInsert || c.spec.kind == speckindUpdate {
			if fld.usedefault && !col.hasdefault {
				// TODO error
			}
		}

		if fld.usejson && !col.typ.is(pgtyp_json, pgtyp_jsonb) {
			return errors.BadUseJSONTargetColumnError
		}
		if fld.usexml && !col.typ.is(pgtyp_xml) {
			return errors.BadUseXMLTargetColumnError
		}

		// Make sure that a value of the given field's type
		// can be assigned to given column, and vice versa.
		if !c.canassign(col, fld, atyp) {
			return errors.BadFieldToColumnTypeError
		}

		if c.spec.kind == speckindSelect || isresult {
			cid := colid{name: fld.colid.name, qual: c.spec.rel.relid.alias}
			pair := &fieldcolumn{field: fld, column: col, colid: cid}
			c.info.output = append(c.info.output, pair)
		}

		if c.spec.kind == speckindInsert || c.spec.kind == speckindUpdate {
			cid := colid{name: fld.colid.name, qual: c.spec.rel.relid.alias}
			pair := &fieldcolumn{field: fld, column: col, colid: cid}
			c.info.input = append(c.info.input, pair)
		}
	}
	return nil
}

func (c *pgchecker) checkjoin(jb *joinblock) error {
	if len(jb.rel.name) > 0 {
		rel, err := c.loadrelation(jb.rel)
		if err != nil {
			return err
		}
		c.joinlist = append(c.joinlist, rel)
	}
	for _, join := range jb.items {
		rel, err := c.loadrelation(join.rel)
		if err != nil {
			return err
		}
		c.joinlist = append(c.joinlist, rel)

		for _, cond := range join.conds {
			switch node := cond.node.(type) {
			case *columnpredicate:
				// A join condition's left-hand-side column MUST always
				// reference a column of the relation being joined, so to
				// avoid confusion make sure that node.colid has either no
				// qualifier or, if it has one, it matches the alias of
				// the joined table.
				if node.colid.qual != "" && node.colid.qual != join.rel.alias {
					return errors.NoDBRelationError
				}

				// Make sure that colid is present in relation being joined.
				col := rel.column(node.colid.name)
				if col == nil {
					return errors.NoDBColumnError
				}

				if node.pred.isunary() {
					// Column type must be bool if the predicate type produces
					// the "IS [NOT] { FALSE | TRUE | UNKNOWN }" SQL predicate.
					if node.pred.istruth() && col.typ.oid != pgtyp_bool {
						return errors.BadColumnTypeForUnaryOpError
					}
					// Column must be NULLable if the predicate type
					// produces the "IS [NOT] NULL" SQL predicate.
					if col.hasnotnull && (node.pred == isnull || node.pred == notnull) {
						return errors.BadColumnNULLSettingForNULLOpError
					}
				} else {
					var typ *pgtype
					// Get the type of the right hand side, which is
					// either a column or a literal expression.
					if len(node.colid2.name) > 0 {
						colid2, err := c.column(node.colid2)
						if err != nil {
							return err
						}
						typ = colid2.typ
					} else if len(node.lit) > 0 {
						var oid pgoid
						row := c.pg.db.QueryRow(fmt.Sprintf(pgselectexprtype, node.lit))
						if err := row.Scan(&oid); err != nil {
							return errors.BadLiteralExpressionError
						}
						typ = c.pg.cat.types[oid]
					}

					if node.pred.isarrcmp() || node.qua > 0 {
						// Check that the quantifier can be used
						// with the type of the RHS expression.
						if typ.category != pgtypcategory_array {
							return errors.BadExpressionTypeForQuantifierError
						}
						typ = c.pg.cat.types[typ.elem]
					}

					rhsoids := []pgoid{typ.oid}
					if !c.cancompare(col, rhsoids, node.pred) {
						return errors.BadColumnToLiteralComparisonError
					}
				}
			default:
				// NOTE(mkopriva): currently predicates other than
				// columnpredicate are not supported as a join condition.
			}
		}
	}
	return nil
}

func (c *pgchecker) checkonconflict(onconf *onconflictblock) error {
	rel := c.rel

	// If a Column directive was provided in an OnConflict block make
	// sure that the listed columns are present in the target table.
	// Make also msure that the list of columns matches the full list
	// of columns of a unique index that's present on the target table.
	if len(onconf.column) > 0 {
		var attnums []int16
		for _, id := range onconf.column {
			col := rel.column(id.name)
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
			if !matchnums(ind.key, attnums) {
				continue
			}

			match = true
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
		ind := rel.index(onconf.index)
		if ind == nil {
			return errors.NoDBIndexError
		}
		if !ind.isunique && !ind.isprimary {
			return errors.NoDBIndexError
		}

		// TODO(mkopriva): retain the index expression so that it can
		// be used later during code generation.
	}

	// If a Constraint directive was provided make sure that the
	// specified constraint is present on the target table and that
	// it is a unique constraint.
	if len(onconf.constraint) > 0 {
		con := rel.constraint(onconf.constraint)
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
			if col := rel.column(item.name); col == nil {
				return errors.NoDBColumnError
			}
		}
	}
	return nil
}

func (c *pgchecker) checkwhere(where *whereblock) error {
	type loopstate struct {
		items []*predicateitem // the current iteration predicate items
		idx   int              // keeps track of the field index
	}
	stack := []*loopstate{{items: where.items}} // LIFO stack of states.

stackloop:
	for len(stack) > 0 {
		// Loop over the various items of a whereblock, including
		// other nested whereblocks and check them against the db.
		loop := stack[len(stack)-1]
		for loop.idx < len(loop.items) {
			item := loop.items[loop.idx]

			// Instead of incrementing the index in the for-statement
			// it is done here manually to ensure that it is not skipped
			// when continuing to the outer loop.
			loop.idx++

			switch node := item.node.(type) {
			case *fieldpredicate:
				// Check that the referenced Column is present
				// in one of the associated relations.
				col, err := c.column(node.colid)
				if err != nil {
					return err
				}

				// If the column cannot be set to NULL, then make
				// sure that the field's not a pointer.
				if col.hasnotnull {
					if node.field.typ.kind == kindptr {
						return errors.IllegalPtrFieldForNotNullColumnError
					}
				}

				// list of types to which the field type can potentially be converted
				var fieldoids = c.typeoids(node.field.typ)

				// If this is a quantified predicate then check that
				// the field is a slice or array, and also make sure that
				// the column's type can be compared to the element type
				// of the slice / array.
				if node.qua > 0 || node.pred.isarrcmp() {
					if node.field.typ.kind != kindslice && node.field.typ.kind != kindarray {
						return errors.IllegalFieldTypeForQuantifierError
					}
					fieldoids = c.typeoids(*node.field.typ.elem)
				}

				if len(node.modfunc) > 0 {
					// Check that the modifier function can
					// be used with the given Column's type.
					if err := c.checkmodfunc(node.modfunc, col, fieldoids); err != nil {
						return err
					}
				} else {
					// Check that the Field's type can be
					// compared to that of the Column.
					if !c.cancompare(col, fieldoids, node.pred) {
						return errors.BadFieldToColumnTypeError
					}
				}

				node.coltype = col.typ.namefmt
			case *columnpredicate:
				// Check that the referenced Column is present
				// in one of the associated relations.
				col, err := c.column(node.colid)
				if err != nil {
					return err
				}

				if node.pred.isunary() {
					// Column type must be bool if the predicate type produces
					// the "IS [NOT] { FALSE | TRUE | UNKNOWN }" SQL predicate.
					if node.pred.istruth() && col.typ.oid != pgtyp_bool {
						return errors.BadColumnTypeForUnaryOpError
					}
					// Column must be NULLable if the predicate type
					// produces the "IS [NOT] NULL" SQL predicate.
					if col.hasnotnull && (node.pred == isnull || node.pred == notnull) {
						return errors.BadColumnNULLSettingForNULLOpError
					}
				} else {
					var typ *pgtype

					// Get the type of the right hand side, which is
					// either a column or a literal expression.
					if len(node.colid2.name) > 0 {
						col2, err := c.column(node.colid2)
						if err != nil {
							return err
						}
						typ = col2.typ
					} else if len(node.lit) > 0 {
						var oid pgoid
						row := c.pg.db.QueryRow(fmt.Sprintf(pgselectexprtype, node.lit))
						if err := row.Scan(&oid); err != nil {
							return errors.BadLiteralExpressionError
						}
						typ = c.pg.cat.types[oid]
					} else {
						panic("shouldn't happen")
					}

					if node.pred.isarrcmp() || node.qua > 0 {
						// Check that the quantifier can be used
						// with the type of the RHS expression.
						if typ.category != pgtypcategory_array {
							return errors.BadExpressionTypeForQuantifierError
						}
						typ = c.pg.cat.types[typ.elem]
					}

					rhsoids := []pgoid{typ.oid}
					if !c.cancompare(col, rhsoids, node.pred) {
						return errors.BadColumnToLiteralComparisonError
					}
				}
			case *betweenpredicate:
				// Check that the referenced Column is present
				// in one of the associated relations.
				col, err := c.column(node.colid)
				if err != nil {
					return err
				}

				// Check that both predicands, x and y, can be compared to the column.
				for _, arg := range []interface{}{node.x, node.y} {
					var argoids []pgoid
					switch a := arg.(type) {
					case colid:
						col2, err := c.column(a)
						if err != nil {
							return err
						}
						argoids = []pgoid{col2.typ.oid}
					case *paramfield:
						argoids = c.typeoids(a.typ)
					}

					if !c.cancompare(col, argoids, isgt) {
						return errors.BadColumnToColumnTypeComparisonError
					}
				}
			case *nestedpredicate:
				stack = append(stack, &loopstate{items: node.items})
				continue stackloop
			}
		}

		stack = stack[:len(stack)-1]
	}

	return nil
}

// typeoids returns a list of OIDs of those PostgreSQL types that can be
// used to hold a value of a Go type represented by the given typeinfo.
func (c *pgchecker) typeoids(typ typeinfo) []pgoid {
	switch typstr := typ.string(true); typstr {
	case gotypstringm, gotypnullstringm:
		if t := c.pg.cat.typebyname("hstore"); t != nil {
			return []pgoid{t.oid}
		}
	case gotypstringms, gotypnullstringms:
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

// cancompare reports whether a value of the given col's type can be compared to,
// using the predicate, a value of one of the types specified by the given oids.
func (c *pgchecker) cancompare(col *pgcolumn, rhstypes []pgoid, pred predicate) bool {
	name := predicate2basepgops[pred]
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

// canassign reports whether a value
func (c *pgchecker) canassign(col *pgcolumn, field *recfield, atyp assigntype) bool {
	// TODO(mkopriva): this returns on success, so write tests that test
	// successful scenarios...

	// If the column is gonna be written to and the field's type knows
	// how to encode itself to a database value, accept.
	if atyp == assignwrite && field.typ.isvaluer {
		return true
	}

	// If the column is gonna be read from and the field's type knows how
	// to decode itself from a database value, accept.
	if atyp == assignread && field.typ.isscanner {
		return true
	}

	// If the column's type is json(b) and the "usejson" directive was used or
	// the field's type implements json.Marshaler/json.Unmarshaler, accept.
	if col.typ.oid == pgtyp_json || col.typ.oid == pgtyp_jsonb {
		if field.usejson || (atyp == assignwrite && field.typ.canjsonmarshal()) ||
			(atyp == assignread && field.typ.canjsonunmarshal()) {
			return true
		}
	}

	// If the column's type is json(b) array and the field's type is a slice
	// whose element type implements json.Marshaler/json.Unmarshaler, accept.
	if col.typ.oid == pgtyp_jsonarr || col.typ.oid == pgtyp_jsonbarr {
		if (atyp == assignwrite && field.typ.kind == kindslice && field.typ.elem.canjsonmarshal()) ||
			(atyp == assignread && field.typ.kind == kindslice && field.typ.elem.canjsonunmarshal()) {
			return true
		}
	}

	// If the column's type is xml and the "usexml" directive was used or
	// the field's type implements xml.Marshaler/xml.Unmarshaler, accept.
	if col.typ.oid == pgtyp_xml {
		if field.usexml || (atyp == assignwrite && field.typ.canxmlmarshal()) ||
			(atyp == assignread && field.typ.canxmlunmarshal()) {
			return true
		}
	}

	// If the column's type is xml array and the field's type is a slice
	// whose element type implements xml.Marshaler/xml.Unmarshaler, accept.
	if col.typ.oid == pgtyp_xmlarr {
		if (atyp == assignwrite && field.typ.kind == kindslice && field.typ.elem.canxmlmarshal()) ||
			(atyp == assignread && field.typ.kind == kindslice && field.typ.elem.canxmlunmarshal()) {
			return true
		}
	}

	conv := pg2goconv{pgtyp: col.typ.oid, gotyp: field.typ.string(true)}

	// Columns with a type in the bit or char family and a typmod of 1 have
	// a distinct Go representation then those with a typmod != 1.
	if col.typ.isbase(pgtyp_bit, pgtyp_varbit, pgtyp_char, pgtyp_varchar, pgtyp_bpchar) {
		conv.typmod1 = (col.typmod == 1)
	}

	// Columns with type numeric that have a precision but no scale, have
	// a distinct Go representation then those numeric types that have a
	// different precision and scale.
	if col.typ.is(pgtyp_numeric, pgtyp_numericarr) {
		precision := ((col.typmod - 4) >> 16) & 65535
		scale := (col.typmod - 4) & 65535
		if precision > 0 && scale == 0 {
			conv.noscale = true
		}
	}

	if _, ok := pg2goconversions[conv]; ok {
		return true
	}

	// If casting is allowed check if the type to which the field will be
	// converted can be coerced to the type of the column.
	if field.cancast {
		return c.cancoerce(col, field)
	}
	return false
}

// cancoerce reports whether or not a value of the given field's type can
// be coerced into a value of the column's type.
func (c *pgchecker) cancoerce(col *pgcolumn, field *recfield) bool {
	// if the target type is of the string category, accept.
	if col.typ.category == pgtypcategory_string {
		return true
	}
	// if the target type is of the array category with an element type of
	// the string category, and the source type is a slice or an array, accept.
	if col.typ.category == pgtypcategory_array && (field.typ.kind == kindslice || field.typ.kind == kindarray) {
		elemtyp := c.pg.cat.types[col.typ.elem]
		if elemtyp != nil && elemtyp.category == pgtypcategory_string {
			return true
		}
	}

	targetid := col.typ.oid
	inputids := c.typeoids(field.typ)
	for _, inputid := range inputids {
		if c.cancoerceoid(targetid, inputid) {
			return true
		}

		if col.typ.category == pgtypcategory_array && (targetid != pgtyp_int2vector && targetid != pgtyp_oidvector) {
			if srctyp := c.pg.cat.types[inputid]; srctyp != nil && srctyp.category == pgtypcategory_array {
				if col.typ.elem == srctyp.elem {
					return true
				}
				if c.pg.cat.cancasti(col.typ.elem, srctyp.elem) {
					return true
				}
			}
		}
	}
	if col.typ.typ == pgtyptype_domain {
		// TODO(mkopriva): implement cancoerce for domain types
		return false
	}
	if col.typ.typ == pgtyptype_composite {
		// TODO(mkopriva): implement cancoerce for composite types
		return false
	}
	return false
}

func (c *pgchecker) cancoerceoid(targetid, inputid pgoid) bool {
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
func (c *pgchecker) checkmodfunc(fn funcname, col *pgcolumn, fieldoids []pgoid) error {
	var (
		ok1 bool // column's type can be coerced to the functions argument's type
		ok2 bool // field's type can be assigned to the functions argument's type
	)

	if proclist, ok := c.pg.cat.procs[fn]; ok {
		for _, proc := range proclist {
			// ok if the column's type can be coerced to the function argument's type
			if c.cancoerceoid(proc.argtype, col.typ.oid) {
				ok1 = true
			}

			// ok if one of the fieldoids types can be assigned to the function argument's type
			for _, foid := range fieldoids {
				if c.cancoerceoid(proc.argtype, foid) {
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

func (c *pgchecker) loadrelation(id relid) (*pgrelation, error) {
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
		)
		if err != nil {
			return nil, err
		}
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

func (c *pgchecker) column(id colid) (*pgcolumn, error) {
	rel, ok := c.relmap[id.qual]
	if !ok {
		return nil, errors.NoDBRelationError
	}
	col := rel.column(id.name)
	if col == nil {
		return nil, errors.NoDBColumnError
	}
	return col, nil
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

func (rel *pgrelation) column(name string) *pgcolumn {
	for _, col := range rel.columns {
		if col.name == name {
			return col
		}
	}
	return nil
}

func (rel *pgrelation) constraint(name string) *pgconstraint {
	for _, con := range rel.constraints {
		if con.name == name {
			return con
		}
	}
	return nil
}

func (rel *pgrelation) index(name string) *pgindex {
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
	procs     map[funcname][]*pgproc
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
	pgcataloguecache.RLock()
	cat := pgcataloguecache.m[key]
	pgcataloguecache.RUnlock()
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

	c.procs = make(map[funcname][]*pgproc)

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

		c.procs[funcname(proc.name)] = append(c.procs[funcname(proc.name)], proc)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	pgcataloguecache.Lock()
	pgcataloguecache.m[key] = c
	pgcataloguecache.Unlock()
	return nil
}

type assigntype uint8

const (
	assignread assigntype = iota + 1
	assignwrite
)

type pgoid uint32

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

// represents a conversion between a postgres type and a go type
type pg2goconv struct {
	pgtyp   pgoid
	gotyp   string
	typmod1 bool // typmod is set to 1
	noscale bool // numeric type with precission but no scale
}

// a map of supported conversions
var pg2goconversions = map[pg2goconv]struct{}{
	// typemod=1
	{pgtyp: pgtyp_char, gotyp: gotyprune, typmod1: true}:        {},
	{pgtyp: pgtyp_chararr, gotyp: gotyprunes, typmod1: true}:    {},
	{pgtyp: pgtyp_varchar, gotyp: gotyprune, typmod1: true}:     {},
	{pgtyp: pgtyp_varchararr, gotyp: gotyprunes, typmod1: true}: {},
	{pgtyp: pgtyp_bpchar, gotyp: gotyprune, typmod1: true}:      {},
	{pgtyp: pgtyp_bpchararr, gotyp: gotyprunes, typmod1: true}:  {},
	// numeric with scale=0
	{pgtyp: pgtyp_numeric, gotyp: gotypbigint, noscale: true}:     {},
	{pgtyp: pgtyp_numericarr, gotyp: gotypbigints, noscale: true}: {},
	// everything else
	{pgtyp: pgtyp_bit, gotyp: gotypstring}:          {},
	{pgtyp: pgtyp_bit, gotyp: gotypbytes}:           {},
	{pgtyp: pgtyp_bitarr, gotyp: gotypstrings}:      {},
	{pgtyp: pgtyp_bitarr, gotyp: gotypbytess}:       {},
	{pgtyp: pgtyp_bool, gotyp: gotypbool}:           {},
	{pgtyp: pgtyp_boolarr, gotyp: gotypbools}:       {},
	{pgtyp: pgtyp_box, gotyp: gotypfloat64a2a2}:     {},
	{pgtyp: pgtyp_boxarr, gotyp: gotypfloat64a2a2s}: {},
	{pgtyp: pgtyp_bpchar, gotyp: gotypstring}:       {},
	{pgtyp: pgtyp_bpchar, gotyp: gotypbytes}:        {},
	{pgtyp: pgtyp_bpchararr, gotyp: gotypstrings}:   {},
	{pgtyp: pgtyp_bpchararr, gotyp: gotypbytess}:    {},
	{pgtyp: pgtyp_bpchararr, gotyp: gotypruness}:    {},
	{pgtyp: pgtyp_bytea, gotyp: gotypstring}:        {},
	{pgtyp: pgtyp_bytea, gotyp: gotypbytes}:         {},
	{pgtyp: pgtyp_byteaarr, gotyp: gotypstrings}:    {},
	{pgtyp: pgtyp_byteaarr, gotyp: gotypbytess}:     {},
	{pgtyp: pgtyp_char, gotyp: gotypstring}:         {},
	{pgtyp: pgtyp_char, gotyp: gotypbytes}:          {},
	{pgtyp: pgtyp_chararr, gotyp: gotypstrings}:     {},
	{pgtyp: pgtyp_chararr, gotyp: gotypbytess}:      {},
	{pgtyp: pgtyp_chararr, gotyp: gotypruness}:      {},
	{pgtyp: pgtyp_cidr, gotyp: gotypstring}:         {},
	{pgtyp: pgtyp_cidr, gotyp: gotypipnet}:          {},
	{pgtyp: pgtyp_cidrarr, gotyp: gotypstrings}:     {},
	{pgtyp: pgtyp_cidrarr, gotyp: gotypipnets}:      {},
	// TODO {pgtyp: pgtyp_circle, gotyp: ""}:        {},
	// TODO {pgtyp: pgtyp_circlearr, gotyp: ""}:        {},
	{pgtyp: pgtyp_date, gotyp: gotyptime}:             {},
	{pgtyp: pgtyp_datearr, gotyp: gotyptimes}:         {},
	{pgtyp: pgtyp_daterange, gotyp: gotyptimea2}:      {},
	{pgtyp: pgtyp_daterangearr, gotyp: gotyptimea2s}:  {},
	{pgtyp: pgtyp_float4, gotyp: gotypfloat32}:        {},
	{pgtyp: pgtyp_float4arr, gotyp: gotypfloat32s}:    {},
	{pgtyp: pgtyp_float8, gotyp: gotypfloat64}:        {},
	{pgtyp: pgtyp_float8arr, gotyp: gotypfloat64s}:    {},
	{pgtyp: pgtyp_inet, gotyp: gotypstring}:           {},
	{pgtyp: pgtyp_inet, gotyp: gotypipnet}:            {},
	{pgtyp: pgtyp_inetarr, gotyp: gotypstrings}:       {},
	{pgtyp: pgtyp_inetarr, gotyp: gotypipnets}:        {},
	{pgtyp: pgtyp_int2, gotyp: gotypint16}:            {},
	{pgtyp: pgtyp_int2arr, gotyp: gotypint16s}:        {},
	{pgtyp: pgtyp_int2vector, gotyp: gotypint16s}:     {},
	{pgtyp: pgtyp_int2vectorarr, gotyp: gotypint16ss}: {},
	{pgtyp: pgtyp_int4, gotyp: gotypint32}:            {},
	{pgtyp: pgtyp_int4, gotyp: gotypint}:              {},
	{pgtyp: pgtyp_int4arr, gotyp: gotypint32s}:        {},
	{pgtyp: pgtyp_int4arr, gotyp: gotypints}:          {},
	{pgtyp: pgtyp_int4range, gotyp: gotypint32a2}:     {},
	{pgtyp: pgtyp_int4range, gotyp: gotypinta2}:       {},
	{pgtyp: pgtyp_int4rangearr, gotyp: gotypint32a2s}: {},
	{pgtyp: pgtyp_int4rangearr, gotyp: gotypinta2s}:   {},
	{pgtyp: pgtyp_int8, gotyp: gotypint64}:            {},
	{pgtyp: pgtyp_int8, gotyp: gotypint}:              {},
	{pgtyp: pgtyp_int8arr, gotyp: gotypint64s}:        {},
	{pgtyp: pgtyp_int8arr, gotyp: gotypints}:          {},
	{pgtyp: pgtyp_int8range, gotyp: gotypint64a2}:     {},
	{pgtyp: pgtyp_int8range, gotyp: gotypinta2}:       {},
	{pgtyp: pgtyp_int8rangearr, gotyp: gotypint64a2s}: {},
	{pgtyp: pgtyp_int8rangearr, gotyp: gotypinta2s}:   {},
	// TODO {pgtyp: pgtyp_interval, gotyp: ""}:   {},
	// TODO {pgtyp: pgtyp_intervalarr, gotyp: ""}:   {},
	{pgtyp: pgtyp_json, gotyp: gotypstring}:             {},
	{pgtyp: pgtyp_json, gotyp: gotypbytes}:              {},
	{pgtyp: pgtyp_jsonarr, gotyp: gotypstrings}:         {},
	{pgtyp: pgtyp_jsonarr, gotyp: gotypbytess}:          {},
	{pgtyp: pgtyp_jsonb, gotyp: gotypstring}:            {},
	{pgtyp: pgtyp_jsonb, gotyp: gotypbytes}:             {},
	{pgtyp: pgtyp_jsonbarr, gotyp: gotypstrings}:        {},
	{pgtyp: pgtyp_jsonbarr, gotyp: gotypbytess}:         {},
	{pgtyp: pgtyp_line, gotyp: gotypfloat64a3}:          {},
	{pgtyp: pgtyp_linearr, gotyp: gotypfloat64a3s}:      {},
	{pgtyp: pgtyp_lseg, gotyp: gotypfloat64a2a2}:        {},
	{pgtyp: pgtyp_lsegarr, gotyp: gotypfloat64a2a2s}:    {},
	{pgtyp: pgtyp_macaddr, gotyp: gotypstring}:          {},
	{pgtyp: pgtyp_macaddr, gotyp: gotypbytes}:           {},
	{pgtyp: pgtyp_macaddrarr, gotyp: gotypstrings}:      {},
	{pgtyp: pgtyp_macaddrarr, gotyp: gotypbytess}:       {},
	{pgtyp: pgtyp_macaddr8, gotyp: gotypstring}:         {},
	{pgtyp: pgtyp_macaddr8, gotyp: gotypbytes}:          {},
	{pgtyp: pgtyp_macaddr8arr, gotyp: gotypstrings}:     {},
	{pgtyp: pgtyp_macaddr8arr, gotyp: gotypbytess}:      {},
	{pgtyp: pgtyp_money, gotyp: gotypint64}:             {},
	{pgtyp: pgtyp_moneyarr, gotyp: gotypint64s}:         {},
	{pgtyp: pgtyp_numeric, gotyp: gotypfloat64}:         {},
	{pgtyp: pgtyp_numeric, gotyp: gotypbigfloat}:        {},
	{pgtyp: pgtyp_numericarr, gotyp: gotypfloat64s}:     {},
	{pgtyp: pgtyp_numericarr, gotyp: gotypbigfloats}:    {},
	{pgtyp: pgtyp_numrange, gotyp: gotypfloat64a2}:      {},
	{pgtyp: pgtyp_numrange, gotyp: gotypbigfloata2}:     {},
	{pgtyp: pgtyp_numrangearr, gotyp: gotypfloat64a2s}:  {},
	{pgtyp: pgtyp_numrangearr, gotyp: gotypbigfloata2s}: {},
	{pgtyp: pgtyp_path, gotyp: gotypfloat64a2s}:         {},
	{pgtyp: pgtyp_patharr, gotyp: gotypfloat64a2ss}:     {},
	{pgtyp: pgtyp_point, gotyp: gotypfloat64a2}:         {},
	{pgtyp: pgtyp_pointarr, gotyp: gotypfloat64a2s}:     {},
	{pgtyp: pgtyp_polygon, gotyp: gotypfloat64a2s}:      {},
	{pgtyp: pgtyp_polygonarr, gotyp: gotypfloat64a2ss}:  {},
	{pgtyp: pgtyp_text, gotyp: gotypstring}:             {},
	{pgtyp: pgtyp_text, gotyp: gotypbytes}:              {},
	{pgtyp: pgtyp_textarr, gotyp: gotypstrings}:         {},
	{pgtyp: pgtyp_textarr, gotyp: gotypbytess}:          {},
	{pgtyp: pgtyp_time, gotyp: gotyptime}:               {},
	{pgtyp: pgtyp_timearr, gotyp: gotyptimes}:           {},
	{pgtyp: pgtyp_timestamp, gotyp: gotyptime}:          {},
	{pgtyp: pgtyp_timestamparr, gotyp: gotyptimes}:      {},
	{pgtyp: pgtyp_timestamptz, gotyp: gotyptime}:        {},
	{pgtyp: pgtyp_timestamptzarr, gotyp: gotyptimes}:    {},
	{pgtyp: pgtyp_timetz, gotyp: gotyptime}:             {},
	{pgtyp: pgtyp_timetzarr, gotyp: gotyptimes}:         {},
	{pgtyp: pgtyp_tsquery, gotyp: gotypstring}:          {},
	{pgtyp: pgtyp_tsquery, gotyp: gotypbytes}:           {},
	{pgtyp: pgtyp_tsqueryarr, gotyp: gotypstring}:       {},
	{pgtyp: pgtyp_tsqueryarr, gotyp: gotypbytess}:       {},
	{pgtyp: pgtyp_tsrange, gotyp: gotyptimea2}:          {},
	{pgtyp: pgtyp_tsrangearr, gotyp: gotyptimea2s}:      {},
	{pgtyp: pgtyp_tstzrange, gotyp: gotyptimea2}:        {},
	{pgtyp: pgtyp_tstzrangearr, gotyp: gotyptimea2s}:    {},
	{pgtyp: pgtyp_tsvector, gotyp: gotypstring}:         {},
	{pgtyp: pgtyp_tsvector, gotyp: gotypbytes}:          {},
	{pgtyp: pgtyp_tsvectorarr, gotyp: gotypstrings}:     {},
	{pgtyp: pgtyp_tsvectorarr, gotyp: gotypbytess}:      {},
	{pgtyp: pgtyp_uuid, gotyp: gotypstring}:             {},
	{pgtyp: pgtyp_uuid, gotyp: gotypbytea16}:            {},
	{pgtyp: pgtyp_uuidarr, gotyp: gotypstrings}:         {},
	{pgtyp: pgtyp_uuidarr, gotyp: gotypbytea16s}:        {},
	{pgtyp: pgtyp_varbit, gotyp: gotypstring}:           {},
	{pgtyp: pgtyp_varbit, gotyp: gotypbytes}:            {},
	{pgtyp: pgtyp_varbitarr, gotyp: gotypstrings}:       {},
	{pgtyp: pgtyp_varbitarr, gotyp: gotypbytess}:        {},
	{pgtyp: pgtyp_varchar, gotyp: gotypstring}:          {},
	{pgtyp: pgtyp_varchar, gotyp: gotypbytes}:           {},
	{pgtyp: pgtyp_varchararr, gotyp: gotypstrings}:      {},
	{pgtyp: pgtyp_varchararr, gotyp: gotypbytess}:       {},
	{pgtyp: pgtyp_varchararr, gotyp: gotypruness}:       {},
	{pgtyp: pgtyp_xml, gotyp: gotypstring}:              {},
	{pgtyp: pgtyp_xml, gotyp: gotypbytes}:               {},
	{pgtyp: pgtyp_xmlarr, gotyp: gotypstrings}:          {},
	{pgtyp: pgtyp_xmlarr, gotyp: gotypbytess}:           {},
}

var go2pgoids = map[string][]pgoid{
	gotypbool:         {pgtyp_bool},
	gotypbools:        {pgtyp_boolarr},
	gotypint:          {pgtyp_int4, pgtyp_int2, pgtyp_int8, pgtyp_float4, pgtyp_float8, pgtyp_numeric},
	gotypints:         {pgtyp_int4arr, pgtyp_int2arr, pgtyp_int2vector, pgtyp_int8arr, pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr},
	gotypinta2:        {pgtyp_int4range, pgtyp_int8range, pgtyp_numrange},
	gotypinta2s:       {pgtyp_int4rangearr, pgtyp_int8rangearr, pgtyp_numrangearr},
	gotypint8:         {pgtyp_int2, pgtyp_int4, pgtyp_int8, pgtyp_float4, pgtyp_float8, pgtyp_numeric},
	gotypint8s:        {pgtyp_int2arr, pgtyp_int2vector, pgtyp_int4arr, pgtyp_int8arr, pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr},
	gotypint8ss:       {pgtyp_int2vectorarr},
	gotypint16:        {pgtyp_int2, pgtyp_int4, pgtyp_int8, pgtyp_float4, pgtyp_float8, pgtyp_numeric},
	gotypint16s:       {pgtyp_int2arr, pgtyp_int2vector, pgtyp_int4arr, pgtyp_int8arr, pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr},
	gotypint16ss:      {pgtyp_int2vectorarr},
	gotypint32:        {pgtyp_int4, pgtyp_int2, pgtyp_int8, pgtyp_float4, pgtyp_float8, pgtyp_numeric},
	gotypint32s:       {pgtyp_int4arr, pgtyp_int2arr, pgtyp_int8arr, pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr},
	gotypint32a2:      {pgtyp_int4range, pgtyp_int8range, pgtyp_numrange},
	gotypint32a2s:     {pgtyp_int4rangearr, pgtyp_int8rangearr, pgtyp_numrangearr},
	gotypint64:        {pgtyp_int8, pgtyp_int4, pgtyp_int2, pgtyp_float4, pgtyp_float8, pgtyp_numeric, pgtyp_money},
	gotypint64s:       {pgtyp_int8arr, pgtyp_int4arr, pgtyp_int2arr, pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr, pgtyp_moneyarr},
	gotypint64a2:      {pgtyp_int8range, pgtyp_int4range, pgtyp_numrange},
	gotypint64a2s:     {pgtyp_int8rangearr, pgtyp_int4rangearr, pgtyp_numrangearr},
	gotypuint:         {pgtyp_int4, pgtyp_int2, pgtyp_int8, pgtyp_float4, pgtyp_float8, pgtyp_numeric},
	gotypuints:        {pgtyp_int4arr, pgtyp_int2arr, pgtyp_int2vector, pgtyp_int8arr, pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr},
	gotypuint8:        {pgtyp_int2, pgtyp_int4, pgtyp_int8, pgtyp_float4, pgtyp_float8, pgtyp_numeric},
	gotypuint8s:       {pgtyp_int2arr, pgtyp_int2vector, pgtyp_int4arr, pgtyp_int8arr, pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr},
	gotypuint16:       {pgtyp_int2, pgtyp_int4, pgtyp_int8, pgtyp_float4, pgtyp_float8, pgtyp_numeric},
	gotypuint16s:      {pgtyp_int2arr, pgtyp_int2vector, pgtyp_int4arr, pgtyp_int8arr, pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr},
	gotypuint32:       {pgtyp_int4, pgtyp_int2, pgtyp_int8, pgtyp_float4, pgtyp_float8, pgtyp_numeric},
	gotypuint32s:      {pgtyp_int4arr, pgtyp_int2arr, pgtyp_int8arr, pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr},
	gotypuint64:       {pgtyp_int8, pgtyp_int4, pgtyp_int2, pgtyp_float4, pgtyp_float8, pgtyp_numeric, pgtyp_money},
	gotypuint64s:      {pgtyp_int8arr, pgtyp_int4arr, pgtyp_int2arr, pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr, pgtyp_moneyarr},
	gotypfloat32:      {pgtyp_float4, pgtyp_float8, pgtyp_numeric},
	gotypfloat32s:     {pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr},
	gotypfloat64:      {pgtyp_float4, pgtyp_float8, pgtyp_numeric},
	gotypfloat64s:     {pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr},
	gotypfloat64a2:    {pgtyp_int4rangearr, pgtyp_int8rangearr, pgtyp_numrange, pgtyp_point},
	gotypfloat64a2s:   {pgtyp_numrangearr, pgtyp_path, pgtyp_pointarr, pgtyp_polygon},
	gotypfloat64a2ss:  {pgtyp_patharr, pgtyp_polygonarr},
	gotypfloat64a2a2:  {pgtyp_box, pgtyp_lseg, pgtyp_path},
	gotypfloat64a2a2s: {pgtyp_boxarr, pgtyp_lsegarr, pgtyp_patharr},
	gotypfloat64a3:    {pgtyp_line, pgtyp_float4arr, pgtyp_float8arr, pgtyp_numericarr},
	gotypfloat64a3s:   {pgtyp_linearr},
	gotypbyte:         {pgtyp_bpchar, pgtyp_bytea, pgtyp_char, pgtyp_text, pgtyp_varchar},
	gotypbytes:        {pgtyp_bit, pgtyp_bpchar, pgtyp_bytea, pgtyp_char, pgtyp_cidr, pgtyp_inet, pgtyp_json, pgtyp_jsonb, pgtyp_macaddr, pgtyp_macaddr8, pgtyp_text, pgtyp_tsquery, pgtyp_tsvector, pgtyp_uuid, pgtyp_varbit, pgtyp_varchar, pgtyp_xml},
	gotypbytess:       {pgtyp_bitarr, pgtyp_bpchararr, pgtyp_byteaarr, pgtyp_chararr, pgtyp_cidrarr, pgtyp_inetarr, pgtyp_jsonarr, pgtyp_jsonbarr, pgtyp_macaddrarr, pgtyp_macaddr8arr, pgtyp_textarr, pgtyp_tsqueryarr, pgtyp_tsvectorarr, pgtyp_uuidarr, pgtyp_varbitarr, pgtyp_varchararr, pgtyp_xmlarr},
	gotypstring:       {pgtyp_bit, pgtyp_bpchar, pgtyp_bytea, pgtyp_char, pgtyp_cidr, pgtyp_inet, pgtyp_json, pgtyp_jsonb, pgtyp_macaddr, pgtyp_macaddr8, pgtyp_text, pgtyp_tsquery, pgtyp_tsvector, pgtyp_uuid, pgtyp_varbit, pgtyp_varchar, pgtyp_xml},
	gotypstrings:      {pgtyp_bitarr, pgtyp_bpchararr, pgtyp_byteaarr, pgtyp_chararr, pgtyp_cidrarr, pgtyp_inetarr, pgtyp_jsonarr, pgtyp_jsonbarr, pgtyp_macaddrarr, pgtyp_macaddr8arr, pgtyp_textarr, pgtyp_tsqueryarr, pgtyp_tsvectorarr, pgtyp_uuidarr, pgtyp_varbitarr, pgtyp_varchararr, pgtyp_xmlarr},
	gotypbytea16:      {pgtyp_uuid, pgtyp_bytea, pgtyp_text, pgtyp_varchar},
	gotypbytea16s:     {pgtyp_uuidarr, pgtyp_byteaarr, pgtyp_textarr, pgtyp_varchararr},
	gotyprune:         {pgtyp_char, pgtyp_bytea, pgtyp_varchar, pgtyp_bpchar},
	gotyprunes:        {pgtyp_bpchar, pgtyp_char, pgtyp_text, pgtyp_varchar /* the rest assumes typmod=1 */, pgtyp_bpchararr, pgtyp_chararr, pgtyp_text, pgtyp_varchararr},
	gotypruness:       {pgtyp_bpchararr, pgtyp_byteaarr, pgtyp_chararr, pgtyp_textarr, pgtyp_varchararr},
	gotypipnet:        {pgtyp_cidr, pgtyp_inet, pgtyp_text, pgtyp_bpchar, pgtyp_bytea, pgtyp_char, pgtyp_varchar},
	gotypipnets:       {pgtyp_cidrarr, pgtyp_inetarr, pgtyp_textarr, pgtyp_bpchararr, pgtyp_byteaarr, pgtyp_chararr, pgtyp_varchararr},
	gotyptime:         {pgtyp_date, pgtyp_time, pgtyp_timestamp, pgtyp_timestamptz, pgtyp_timetz},
	gotyptimes:        {pgtyp_datearr, pgtyp_timearr, pgtyp_timestamparr, pgtyp_timestamptzarr, pgtyp_timetzarr, pgtyp_intervalarr},
	gotyptimea2:       {pgtyp_daterange, pgtyp_tsrange, pgtyp_tstzrange},
	gotyptimea2s:      {pgtyp_daterangearr, pgtyp_tsrangearr, pgtyp_tstzrangearr},
	gotypbigint:       {pgtyp_numeric, pgtyp_text},
	gotypbigints:      {pgtyp_numericarr, pgtyp_textarr},
	gotypbiginta2:     {pgtyp_numrange},
	gotypbiginta2s:    {pgtyp_numrangearr},
	gotypbigfloat:     {pgtyp_numeric, pgtyp_text},
	gotypbigfloats:    {pgtyp_numericarr, pgtyp_textarr},
	gotypbigfloata2:   {pgtyp_numrange},
	gotypbigfloata2s:  {pgtyp_numrangearr},

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
var predicate2basepgops = map[predicate]string{
	iseq:        "=",
	noteq:       "<>",
	noteq2:      "<>",
	islt:        "<",
	isgt:        ">",
	islte:       "<=",
	isgte:       ">=",
	ismatch:     "~",
	ismatchi:    "~*",
	notmatch:    "!~",
	notmatchi:   "!~*",
	isdistinct:  "<>",
	notdistinct: "=",
	islike:      "~~",
	notlike:     "!~~",
	isilike:     "~~*",
	notilike:    "!~~*",
	issimilar:   "~~",
	notsimilar:  "!~~",
	isin:        "=",
	notin:       "<>",
}

////////////////////////////////////////////////////////////////////////////////

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

// helper func
func matchnums(a, b []int16) bool {
	if len(a) != len(b) {
		return false
	}

aloop:
	for _, x := range a {
		for _, y := range b {
			if x == y {
				continue aloop
			}
		}
		return false // x not found in b
	}
	return true
}
