// TODO(mkopriva):
// - given the gosql.Index directive inside an OnConflict block use pg_get_indexdef(index_oid)
//   to retrieve the index's definition, parse that to extract the index expression and then
//   use that expression when generating the ON CONFLICT clause.
//   (https://www.postgresql.org/message-id/204ADCAA-853B-4B5A-A080-4DFA0470B790%40justatheory.com)
// TODO the postgres types circle(arr) and interval(arr) need a corresponding go type

package gosql

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
	selectdbname = `SELECT current_database()` //`

	selectdbrelation = `SELECT
		c.oid
		, c.relkind
	FROM pg_class c
	WHERE c.relname = $1
	AND c.relnamespace = to_regnamespace($2)`  //`

	selectdbcolumns = `SELECT
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

	selectdbconstraints = `SELECT
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

	selectdbindexes = `SELECT
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

	selecttypes = `SELECT
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

	selectoperators = `SELECT
		o.oid
		, o.oprname
		, o.oprkind
		, o.oprleft
		, o.oprright
		, o.oprresult
	FROM pg_operator o `  //`

	selectcasts = `SELECT
		c.oid
		, c.castsource
		, c.casttarget
		, c.castcontext
	FROM pg_cast c `  //`

	selectexprtype = `SELECT pg_typeof(%s)` //`

)

func dbcheck(url string, cmds []*command) error {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return err
	} else if err := db.Ping(); err != nil {
		return err
	}
	defer db.Close()

	// the name of the current database
	var dbname string
	if err := db.QueryRow(selectdbname).Scan(&dbname); err != nil {
		return err
	}

	pgcat := new(pgcatalogue)
	if err := pgcat.load(db, url); err != nil {
		return err
	}

	for _, cmd := range cmds {
		c := new(dbchecker)
		c.db = db
		c.dbname = dbname
		c.cmd = cmd
		c.pgcat = pgcat
		if err := c.run(); err != nil {
			return err
		}
	}
	return nil
}

type dbchecker struct {
	db       *sql.DB
	dbname   string // name of the current database, used mainly for error reporting
	cmd      *command
	rel      *pgrelation   // the target relation
	joinlist []*pgrelation // joined relations
	relmap   map[string]*pgrelation

	pgcat *pgcatalogue
}

func (c *dbchecker) run() (err error) {
	c.relmap = make(map[string]*pgrelation)
	if c.rel, err = c.loadrelation(c.cmd.rel.relid); err != nil {
		return err
	}

	// Map the target relation to the "" (empty string) key, this will allow
	// columns, constraints, and indexes that were specified without a qualifier
	// to be associated with this target relation.
	c.relmap[""] = c.rel

	if join := c.cmd.join; join != nil {
		if err := c.checkjoin(join); err != nil {
			return err
		}
	}
	if onconf := c.cmd.onconflict; onconf != nil {
		if err := c.checkonconflict(onconf); err != nil {
			return err
		}
	}
	if where := c.cmd.where; where != nil {
		if err := c.checkwhere(where); err != nil {
			return err
		}
	}

	// If an OrderBy directive was used, make sure that the specified
	// columns are present in the loaded relations.
	if c.cmd.orderby != nil {
		for _, item := range c.cmd.orderby.items {
			if _, err := c.column(item.col); err != nil {
				return err
			}
		}
	}

	// If a Default directive was provided, make sure that the specified
	// columns are present in the target relation.
	if c.cmd.defaults != nil {
		for _, item := range c.cmd.defaults.items {
			// Qualifier, if present, must match the target table's alias.
			if len(item.qual) > 0 && item.qual != c.cmd.rel.relid.alias {
				return errors.BadTargetTableForDefaultError
			}
			if col := c.rel.column(item.name); col == nil {
				return errors.NoDBColumnError
			}
		}
	}

	// If a Force directive was provided, make sure that the specified
	// columns are present in the loaded relations.
	if c.cmd.force != nil {
		for _, item := range c.cmd.force.items {
			if _, err := c.column(item); err != nil {
				return err
			}
		}
	}

	// If a Return directive was provided, make sure that the specified
	// columns are present in the loaded relations.
	if c.cmd.returning != nil {
		for _, item := range c.cmd.returning.items {
			if _, err := c.column(item); err != nil {
				return err
			}
		}
	}

	// If a TextSearch directive was provided, make sure that the
	// specified column is present in one of the loaded relations
	// and that it has the correct type.
	if c.cmd.textsearch != nil {
		col, err := c.column(*c.cmd.textsearch)
		if err != nil {
			return err
		} else if col.typ.oid != pgtyp_tsvector {
			return errors.BadDBColumnTypeError
		}
	}

	if rel := c.cmd.rel; rel != nil && !rel.isdir {
		if err := c.checkfields(rel.datatype.base, false); err != nil {
			return err
		}
	}

	if res := c.cmd.result; res != nil {
		if err := c.checkfields(res.datatype.base, true); err != nil {
			return err
		}
	}

	return nil
}

func (c *dbchecker) checkfields(typ typeinfo, isresult bool) (err error) {
	type loopstate struct {
		typ *typeinfo
		idx int // keeps track of the field index
	}

	stack := []*loopstate{{typ: &typ}} // lifo stack

stackloop:
	for len(stack) > 0 {
		loop := stack[len(stack)-1]
		for loop.idx < len(loop.typ.fields) {
			fld := loop.typ.fields[loop.idx]

			// Instead of incrementing the index in the for-statement
			// it is done here manually to ensure that it is not skipped
			// when continuing to the outer loop.
			loop.idx++

			if len(fld.typ.fields) > 0 {
				loop2 := new(loopstate)
				loop2.typ = &fld.typ
				stack = append(stack, loop2)
				continue stackloop
			}

			var col *pgcolumn
			if c.cmd.typ == cmdtypeSelect || isresult {
				// If this is a SELECT, or the target type is
				// from the "Result" field, lookup the column
				// in all of the associated relations since its
				// ok to select columns from joined relations.
				if col, err = c.column(fld.colid); err != nil {
					return err
				}
			} else {
				// If this is a non-select command, non-result the
				// columns must be present in the target relation.
				if col = c.rel.column(fld.colid.name); col == nil {
					return errors.NoDBColumnError
				}
			}

			// Make sure that a value of the given field's type
			// can be assigned to given column, and vice versa.
			if !c.cancoerce(col, fld) {
				return nil
			}
		}
		stack = stack[:len(stack)-1]
	}
	return nil
}

func (c *dbchecker) checkjoin(join *joinblock) error {
	if len(join.rel.name) > 0 {
		rel, err := c.loadrelation(join.rel)
		if err != nil {
			return err
		}
		c.joinlist = append(c.joinlist, rel)
	}
	for _, item := range join.items {
		rel, err := c.loadrelation(item.rel)
		if err != nil {
			return err
		}
		c.joinlist = append(c.joinlist, rel)

		for _, cond := range item.conds {
			// Make sure that col1 is present in relation being joined.
			col := rel.column(cond.col1.name)
			if col == nil {
				return errors.NoDBColumnError
			}

			if cond.cmp.isunary() {
				// Column type must be bool if the comparison operator
				// is one of the IS [NOT] {FALSE|TRUE|UNKNOWN} operators.
				if cond.cmp.isbool() && col.typ.oid != pgtyp_bool {
					return errors.BadColumnTypeForUnaryOpError
				}
				// Column must be NULLable if the comparison operator
				// is one of the IS [NOT] NULL operators.
				if col.hasnotnull && (cond.cmp == cmpisnull || cond.cmp == cmpnotnull) {
					return errors.BadColumnNULLSettingForNULLOpError
				}
			} else {
				var typ *pgtype
				// Get the type of the right hand side, which is
				// either a column or a literal expression.
				if len(cond.col2.name) > 0 {
					col2, err := c.column(cond.col2)
					if err != nil {
						return err
					}
					typ = col2.typ
				} else if len(cond.lit) > 0 {
					var oid pgoid
					row := c.db.QueryRow(fmt.Sprintf(selectexprtype, cond.lit))
					if err := row.Scan(&oid); err != nil {
						return err
					}
					typ = c.pgcat.types[oid]
				}

				if cond.cmp.isarr() || cond.sop > 0 {
					// Check that the scalar array operator can
					// be used with the type of the RHS expression.
					if typ.category != pgtypcategory_array {
						return errors.BadExpressionTypeForScalarrOpError
					}
					typ = c.pgcat.types[typ.elem]
				}

				rhsoids := []pgoid{typ.oid}
				if !c.cancompare(col, rhsoids, cond.cmp) {
					return errors.BadColumnToLiteralComparisonError
				}
			}
		}
	}
	return nil
}

func (c *dbchecker) checkonconflict(onconf *onconflictblock) error {
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

func (c *dbchecker) checkwhere(where *whereblock) error {
	type loopstate struct {
		wb  *whereblock
		idx int // keeps track of the field index
	}
	stack := []*loopstate{{wb: where}} // LIFO stack of states.

stackloop:
	for len(stack) > 0 {
		// Loop over the various items of a whereblock, including
		// other nested whereblocks and check them against the db.
		loop := stack[len(stack)-1]
		for loop.idx < len(loop.wb.items) {
			loop.idx++

			switch node := loop.wb.items[loop.idx].node.(type) {
			case *wherefield:
				// Check that the referenced Column is present
				// in one of the associated relations.
				col, err := c.column(node.colid)
				if err != nil {
					return err
				}

				// If the column cannot be set to NULL, then make
				// sure that the field's not a pointer.
				if col.hasnotnull {
					if node.typ.kind == kindptr {
						return errors.IllegalPtrFieldForNotNullColumnError
					}
				}

				// list of types to which the field type can potentially be converted
				var fieldoids = c.pgtypeoids(node.typ)

				// If this is a scalar array comparison then check that
				// the field is a slice or array, and also make sure that
				// the column's type can be compared to the element type
				// of the slice / array.
				if node.sop > 0 || node.cmp.isarr() {
					if node.typ.kind != kindslice && node.typ.kind != kindarray {
						return errors.IllegalFieldTypeForScalarOpError
					}
					fieldoids = c.pgtypeoids(*node.typ.elem)
				}

				// Check that the Field's type can be compared to that of the Column.
				if !c.cancompare(col, fieldoids, node.cmp) {
					return errors.BadFieldToColumnTypeError
				}

				// Check that the modifier function can
				// be used with the given Column's type.
				if node.mod > 0 {
					if err := c.checkmodfunc(col, node.mod); err != nil {
						return err
					}
				}
			case *wherecolumn:
				// Check that the referenced Column is present
				// in one of the associated relations.
				col, err := c.column(node.colid)
				if err != nil {
					return err
				}

				if node.cmp.isunary() {
					// Column type must be bool if the comparison operator
					// is one of the IS [NOT] {FALSE|TRUE|UNKNOWN} operators.
					if node.cmp.isbool() && col.typ.oid != pgtyp_bool {
						return errors.BadColumnTypeForUnaryOpError
					}
					// Column must be NULLable if the comparison operator
					// is one of the IS [NOT] NULL operators.
					if col.hasnotnull && (node.cmp == cmpisnull || node.cmp == cmpnotnull) {
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
						row := c.db.QueryRow(fmt.Sprintf(selectexprtype, node.lit))
						if err := row.Scan(&oid); err != nil {
							return err
						}
						typ = c.pgcat.types[oid]
					} else {
						// bail?
					}

					if node.cmp.isarr() || node.sop > 0 {
						// Check that the scalar array operator can
						// be used with the type of the RHS expression.
						if typ.category != pgtypcategory_array {
							return errors.BadExpressionTypeForScalarrOpError
						}
						typ = c.pgcat.types[typ.elem]
					}

					rhsoids := []pgoid{typ.oid}
					if !c.cancompare(col, rhsoids, node.cmp) {
						return errors.BadColumnToLiteralComparisonError
					}
				}
			case *wherebetween:
				// Check that the referenced Column is present
				// in one of the associated relations.
				col, err := c.column(node.colid)
				if err != nil {
					return err
				}

				// Check that both arguments, x and y, can be compared to the column.
				for _, arg := range []interface{}{node.x, node.y} {
					var argoids []pgoid
					switch a := arg.(type) {
					case colid:
						col2, err := c.column(a)
						if err != nil {
							return err
						}
						argoids = []pgoid{col2.typ.oid}
					case *varinfo:
						argoids = c.pgtypeoids(a.typ)
					}

					if !c.cancompare(col, argoids, cmpgt) {
						return errors.BadColumnToColumnTypeComparisonError
					}
				}
			case *whereblock:
				stack = append(stack, &loopstate{wb: node})
				continue stackloop
			}
		}

		stack = stack[:len(stack)-1]
	}

	return nil
}

func (c *dbchecker) pgtypeoids(typ typeinfo) []pgoid {
	// TODO needs more work

	switch typstr := typ.string(); typstr {
	case gotypstringm, gotypstringpm, gotypnullstringm:
		if t := c.pgcat.typebyname("hstore"); t != nil {
			return []pgoid{t.oid}
		}
	case gotypstringms, gotypstringpms, gotypnullstringms:
		if t := c.pgcat.typebyname("_hstore"); t != nil {
			return []pgoid{t.oid}
		}
	default:
		return gotyp2pgtyp[typstr]
	}
	return nil
}

func (c *dbchecker) cancompare(col *pgcolumn, rhstypes []pgoid, cmp cmpop) bool {
	// TODO needs more work

	name := cmpop2basepgops[cmp]
	left := col.typ.oid
	for _, right := range rhstypes {
		key := pgopkey{name: name, left: left, right: right}
		if _, ok := c.pgcat.operators[key]; ok {
			return true
		}
	}
	return false
}

// cancoerce reports whether or not a value of the given field's type can
// be coerced into a value of the column's type.
func (c *dbchecker) cancoerce(col *pgcolumn, field *fieldinfo) bool {
	// TODO needs more work

	target := col.typ.oid
	sourceoids := c.pgtypeoids(field.typ)
	for _, source := range sourceoids {
		if target == source {
			return true
		}
		if col.typ.category == pgtypcategory_string {
			return true
		}
		if c.pgcat.cancasti(target, source) {
			return true
		}

		if col.typ.category == pgtypcategory_array && (target != pgtyp_int2vector && target != pgtyp_oidvector) {
			if srctyp := c.pgcat.types[source]; srctyp != nil && srctyp.category == pgtypcategory_array {
				if col.typ.elem == srctyp.elem {
					return true
				}
				elemtyp := c.pgcat.types[col.typ.elem]
				if elemtyp != nil && elemtyp.category == pgtypcategory_string {
					return true
				}
				if c.pgcat.cancasti(col.typ.elem, srctyp.elem) {
					return true
				}
			}
		}
	}
	if col.typ.typ == pgtyptype_domain {
		// TODO(mkopriva): implement cancoerce for domain types
		return true
	}
	if col.typ.typ == pgtyptype_composite {
		// TODO(mkopriva): implement cancoerce for composite types
		return true
	}
	return false
}

// check that the column's type matches that of the modfunc's argument.
func (c *dbchecker) checkmodfunc(col *pgcolumn, fn modfunc) error {
	var ok bool
	switch fn {
	case fnlower, fnupper:
		ok = col.typ.oid != pgtyp_text
	}
	if !ok {
		return errors.BadColumnTypeToModFuncError
	}
	return nil
}

func (c *dbchecker) loadrelation(id relid) (*pgrelation, error) {
	rel := new(pgrelation)
	rel.name = id.name
	rel.namespace = id.qual
	if len(rel.namespace) == 0 {
		rel.namespace = "public"
	}

	// retrieve relation info
	row := c.db.QueryRow(selectdbrelation, rel.name, rel.namespace)
	if err := row.Scan(&rel.oid, &rel.relkind); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NoDBRelationError
		}
		return nil, err
	}

	// retrieve column info
	rows, err := c.db.Query(selectdbcolumns, rel.oid)
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
		typ, ok := c.pgcat.types[col.typoid]
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
	rows, err = c.db.Query(selectdbconstraints, rel.oid)
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
	rows, err = c.db.Query(selectdbindexes, rel.oid)
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

func (c *dbchecker) column(id colid) (*pgcolumn, error) {
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

type pgcastkey struct {
	target pgoid
	source pgoid
}

type pgcatalogue struct {
	types     map[pgoid]*pgtype
	operators map[pgopkey]*pgoperator
	casts     map[pgcastkey]*pgcast
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

// cancasti reports whether s can be cast to t implicitly or in assignment.
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
	rows, err := db.Query(selecttypes)
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
	rows, err = db.Query(selectoperators)
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
	rows, err = db.Query(selectcasts)
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

	pgcataloguecache.Lock()
	pgcataloguecache.m[key] = c
	pgcataloguecache.Unlock()
	return nil
}

type pgoid uint32

// postgres types
const (
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
	pgtyp_tinterval      pgoid = 704
	pgtyp_tintervalarr   pgoid = 1025
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

// a mapping of go types to a list of pg types where the list contains those
// pg types that can store a value of the go type.
var gotyp2pgtyp = map[string][]pgoid{
	gotypbool:         {pgtyp_bool},
	gotypboolp:        {pgtyp_bool},
	gotypbools:        {pgtyp_boolarr},
	gotypstring:       {pgtyp_bit, pgtyp_varbit, pgtyp_bpchar, pgtyp_char, pgtyp_varchar, pgtyp_cidr, pgtyp_inet, pgtyp_text, pgtyp_tsquery, pgtyp_uuid},
	gotypstringp:      {pgtyp_bit, pgtyp_varbit, pgtyp_bpchar, pgtyp_char, pgtyp_varchar, pgtyp_cidr, pgtyp_inet, pgtyp_text, pgtyp_tsquery, pgtyp_uuid},
	gotypstrings:      {pgtyp_bitarr, pgtyp_varbitarr, pgtyp_bpchararr, pgtyp_chararr, pgtyp_varchararr, pgtyp_cidrarr, pgtyp_inetarr, pgtyp_textarr, pgtyp_tsqueryarr, pgtyp_tsvector, pgtyp_uuidarr},
	gotypstringss:     {pgtyp_tsvectorarr},
	gotypbyte:         {pgtyp_bit, pgtyp_varbit, pgtyp_bpchar, pgtyp_char, pgtyp_varchar},
	gotypbytep:        {pgtyp_bit, pgtyp_varbit, pgtyp_bpchar, pgtyp_char, pgtyp_varchar},
	gotypbytes:        {pgtyp_bit, pgtyp_varbit, pgtyp_bitarr, pgtyp_varbitarr, pgtyp_bpchar, pgtyp_char, pgtyp_varchar, pgtyp_bpchararr, pgtyp_chararr, pgtyp_varchararr, pgtyp_bytea, pgtyp_json, pgtyp_jsonb, pgtyp_macaddr, pgtyp_macaddr8, pgtyp_uuid, pgtyp_xml},
	gotypbytess:       {pgtyp_bitarr, pgtyp_varbitarr, pgtyp_bpchararr, pgtyp_chararr, pgtyp_varchararr, pgtyp_byteaarr, pgtyp_jsonarr, pgtyp_jsonbarr, pgtyp_macaddrarr, pgtyp_macaddr8arr, pgtyp_uuidarr, pgtyp_xmlarr},
	gotypbytea16:      {pgtyp_uuid},
	gotypbytea16s:     {pgtyp_uuidarr},
	gotyprune:         {pgtyp_bpchar, pgtyp_char, pgtyp_varchar},
	gotyprunep:        {pgtyp_bpchar, pgtyp_char, pgtyp_varchar},
	gotyprunes:        {pgtyp_bpchar, pgtyp_char, pgtyp_varchar, pgtyp_bpchararr, pgtyp_chararr, pgtyp_varchararr},
	gotypruness:       {pgtyp_bpchararr, pgtyp_chararr, pgtyp_varchararr},
	gotypint8:         {pgtyp_int2},
	gotypint8p:        {pgtyp_int2},
	gotypint8s:        {pgtyp_int2arr, pgtyp_int2vector},
	gotypint8ss:       {pgtyp_int2vectorarr},
	gotypint16:        {pgtyp_int2},
	gotypint16p:       {pgtyp_int2},
	gotypint16s:       {pgtyp_int2arr, pgtyp_int2vector},
	gotypint16ss:      {pgtyp_int2vectorarr},
	gotypint32:        {pgtyp_int4},
	gotypint32p:       {pgtyp_int4},
	gotypint32s:       {pgtyp_int4arr},
	gotypint32a2:      {pgtyp_int4range},
	gotypint32a2s:     {pgtyp_int4rangearr},
	gotypint64:        {pgtyp_int8, pgtyp_money},
	gotypint64p:       {pgtyp_int8, pgtyp_money},
	gotypint64s:       {pgtyp_int8arr, pgtyp_moneyarr},
	gotypint64a2:      {pgtyp_int8range},
	gotypint64a2s:     {pgtyp_int8rangearr},
	gotypfloat32:      {pgtyp_float4},
	gotypfloat32p:     {pgtyp_float4},
	gotypfloat32s:     {pgtyp_float4arr},
	gotypfloat64:      {pgtyp_float8},
	gotypfloat64p:     {pgtyp_float8},
	gotypfloat64s:     {pgtyp_float8arr},
	gotypfloat64a2:    {pgtyp_point},
	gotypfloat64a2s:   {pgtyp_path, pgtyp_pointarr, pgtyp_polygon},
	gotypfloat64a2ss:  {pgtyp_patharr, pgtyp_polygonarr},
	gotypfloat64a2a2:  {pgtyp_box, pgtyp_lseg},
	gotypfloat64a2a2s: {pgtyp_boxarr, pgtyp_lsegarr},
	gotypfloat64a3:    {pgtyp_line},
	gotypfloat64a3s:   {pgtyp_linearr},
	gotypipnetp:       {pgtyp_cidr, pgtyp_inet},
	gotypipnetps:      {pgtyp_cidrarr, pgtyp_inetarr},
	gotyptime:         {pgtyp_date, pgtyp_time, pgtyp_timestamp, pgtyp_timestamptz, pgtyp_timetz},
	gotyptimep:        {pgtyp_date, pgtyp_time, pgtyp_timestamp, pgtyp_timestamptz, pgtyp_timetz},
	gotyptimes:        {pgtyp_datearr, pgtyp_timearr, pgtyp_timestamparr, pgtyp_timestamptzarr, pgtyp_timetzarr},
	gotyptimeps:       {pgtyp_datearr, pgtyp_timearr, pgtyp_timestamparr, pgtyp_timestamptzarr, pgtyp_timetzarr},
	gotyptimea2:       {pgtyp_daterange, pgtyp_tsrange, pgtyp_tstzrange},
	gotyptimea2s:      {pgtyp_daterangearr, pgtyp_tsrangearr, pgtyp_tstzrangearr},
	gotypbigint:       {pgtyp_numeric},
	gotypbigintp:      {pgtyp_numeric},
	gotypbigints:      {pgtyp_numericarr},
	gotypbigintps:     {pgtyp_numericarr},
	gotypbiginta2:     {pgtyp_numrange},
	gotypbigintpa2:    {pgtyp_numrange},
	gotypbigintpa2s:   {pgtyp_numrangearr},
}

// Map of supported cmpops to *equivalent* postgres comparison operators. For example
// the constructs IN and NOT IN are essentially the same as comparing the LHS to every
// element of the RHS with the operators "=" and "<>" respectively, and therefore the
// cmpisin maps to "=" and cmpnotin maps to "<>".
var cmpop2basepgops = map[cmpop]string{
	cmpeq:          "=",
	cmpne:          "<>",
	cmpne2:         "<>",
	cmplt:          "<",
	cmpgt:          ">",
	cmple:          "<=",
	cmpge:          ">=",
	cmprexp:        "~",
	cmprexpi:       "~*",
	cmpnotrexp:     "!~",
	cmpnotrexpi:    "!~*",
	cmpisdistinct:  "<>",
	cmpnotdistinct: "=",
	cmpislike:      "~~",
	cmpnotlike:     "!~~",
	cmpisilike:     "~~*",
	cmpnotilike:    "!~~*",
	cmpissimilar:   "~~",
	cmpnotsimilar:  "!~~",
	cmpisin:        "=",
	cmpnotin:       "<>",
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
