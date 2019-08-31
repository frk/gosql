// TODO(mkopriva):
// - given the gosql.Index directive inside an OnConflict block use pg_get_indexdef(index_oid)
//   to retrieve the index's definition, parse that to extract the index expression and then
//   use that expression when generating the ON CONFLICT clause.
//   (https://www.postgresql.org/message-id/204ADCAA-853B-4B5A-A080-4DFA0470B790%40justatheory.com)

package gosql

import (
	"database/sql"
	"strconv"
	"strings"

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
		, t.typname
		, t.typlen
		, t.typtype
		, t.typcategory
		, t.typispreferred
		, t.typelem
	FROM pg_attribute a
	LEFT JOIN pg_type t ON t.oid = a.atttypid
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
		, pg_get_indexdef(i.indexrelid)
	FROM pg_index i
	LEFT JOIN pg_class c ON c.oid = i.indexrelid
	WHERE i.indrelid = $1
	ORDER BY i.indexrelid`  //`
)

func dbcheck(url string, cmds []*command) error {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return err
	} else if err := db.Ping(); err != nil {
		return err
	}
	defer db.Close()

	var dbname string // the name of the current database
	if err := db.QueryRow(selectdbname).Scan(&dbname); err != nil {
		return err
	}

	for _, cmd := range cmds {
		c := new(dbchecker)
		c.db = db
		c.dbname = dbname
		c.cmd = cmd
		if err := c.load(); err != nil {
			return err
		}

		if err := c.check(); err != nil {
			return err
		}
	}
	return nil
}

type dbchecker struct {
	db       *sql.DB
	dbname   string // name of the current database, used mainly for error reporting
	cmd      *command
	rel      *dbrelation   // the target relation
	joinlist []*dbrelation // joined relations

	relmap map[string]*dbrelation
}

func (c *dbchecker) load() (err error) {
	c.relmap = make(map[string]*dbrelation)
	if c.rel, err = c.loadrelation(c.cmd.rel.relid); err != nil {
		return err
	}

	// Map the target relation to the "" (empty string) key, this will allow
	// columns, constraints, and indexes that were specified without a qualifier
	// to be associated with this target relation.
	c.relmap[""] = c.rel

	if join := c.cmd.join; join != nil {
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
		}
	}
	return nil
}

func (c *dbchecker) check() error {
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
	// columns are present in the loaded relations.
	if c.cmd.defaults != nil {
		for _, item := range c.cmd.defaults.items {
			if _, err := c.column(item); err != nil {
				return err
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

	// If a TextSearch directive was provided, make sure that the specified
	// column is present in one of the loaded relations and that it has the
	// correct type.
	if c.cmd.textsearch != nil {
		col, err := c.column(*c.cmd.textsearch)
		if err != nil {
			return err
		} else if col.typ.name != pgtyp_tsvector {
			return errors.BadDBColumnTypeError
		}
	}

	// check onconflict block
	if c.cmd.onconflict != nil {
		oc := c.cmd.onconflict
		rel := c.rel

		// If a Column directive was provided in an OnConflict block make
		// sure that the listed columns are present in the target table.
		// Make also msure that the list of columns matches the full list
		// of columns of a unique index that's present on the target table.
		if len(oc.column) > 0 {
			var attnums []int16
			for _, id := range oc.column {
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
		if len(oc.index) > 0 {
			ind := rel.index(oc.index)
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
		if len(oc.constraint) > 0 {
			con := rel.constraint(oc.constraint)
			if con == nil {
				return errors.NoDBConstraintError
			}
			if con.typ != pgconstraint_pkey && con.typ != pgconstraint_unique {
				return errors.NoDBConstraintError
			}
		}

		// If an Update directive was provided, make sure that each
		// listed column is present in the target table.
		if oc.update != nil {
			for _, item := range oc.update.items {
				if col := rel.column(item.name); col == nil {
					return errors.NoDBColumnError
				}
			}
		}
	}

	// check where block
	if c.cmd.where != nil {
		type loopstate struct {
			wb  *whereblock
			idx int // keeps track of the field index
		}
		stack := []*loopstate{{wb: c.cmd.where}} // LIFO stack of states.

	stackloop:
		for len(stack) > 0 {
			// Loop over the various items of a whereblock, including
			// other nested whereblocks and check them against the db.
			loop := stack[len(stack)-1]
			for loop.idx < len(loop.wb.items) {
				loop.idx++

				switch node := loop.wb.items[loop.idx].node.(type) {
				case *wherefield:
					// In case of a wherefield check that the referenced column
					// is present in one of the associated relations and that
					// the field's type is compatible with that of the column.
					//
					// Also make sure that the comparison operator, the scalar
					// array operator, and the modifier function can actually
					// be used with the given column's type.
					//
					// Also check that if column is NULLable
					// then also the file is NILable
					col, err := c.column(node.colid)
					if err != nil {
						return err
					}
					if err := c.checktype(col, node.typ, typeopts{}); err != nil {
						return err
					}
				case *wherecolumn:
				case *wherebetween:
				case *whereblock:
					stack = append(stack, &loopstate{wb: node})
					continue stackloop
				}
			}

			stack = stack[:len(stack)-1]
		}
	}

	// check join block
	if c.cmd.join != nil {
		// TODO
	}

	// check result field
	if c.cmd.result != nil {
		// TODO
	}

	return nil
}

type typeopts struct {
	usejson bool
	usexml  bool
}

func (c *dbchecker) checktype(col *dbcolumn, typ typeinfo, opts typeopts) error {
	var ok bool

	switch col.typ.name {
	case pgtyp_bit, pgtyp_varbit:
		// string, *string, []byte
		// or byte or *byte if bit(1)
		if col.typmod == 1 {
			ok = typ.is(kindbyte)
		}
		if !ok {
			ok = typ.is(kindstring) || typ.isslice(kindbyte)
		}
	case pgtyp_bitarr, pgtyp_varbitarr:
		// []string, [][]byte
		// or []byte if bit(1)
		if col.typmod == 1 {
			ok = typ.isslice(kindbyte)
		}
		if !ok {
			ok = typ.isslice(kindstring) || typ.isslicen(2, kindbyte)
		}
	case pgtyp_bool:
		// bool, *bool
		ok = typ.is(kindbool)
	case pgtyp_boolarr:
		// []bool
		ok = typ.isslice(kindbool)
	case pgtyp_box:
		// [2][2]float64
		ok = typ.kind == kindarray && typ.arraylen == 2 &&
			typ.elem.kind == kindarray && typ.elem.arraylen == 2 &&
			typ.elem.elem.kind == kindfloat64
	case pgtyp_boxarr:
		// [][2][2]float64
		if t := typ.elem; typ.kind == kindslice {
			ok = t.kind == kindarray && t.arraylen == 2 &&
				t.elem.kind == kindarray && t.elem.arraylen == 2 &&
				t.elem.elem.kind == kindfloat64
		}
	case pgtyp_bpchar, pgtyp_char, pgtyp_varchar:
		// []byte, []rune, string, *string
		// or byte, *byte, rune, *rune if char(1)
		if (col.typmod - 4) == 1 {
			ok = typ.is(kindbyte, kindrune)
		}
		if !ok {
			ok = typ.is(kindstring) || typ.isslice(kindbyte, kindrune)
		}
	case pgtyp_bpchararr, pgtyp_chararr, pgtyp_varchararr:
		// [][]byte, [][]rune, []string
		// or []byte or []rune if char(1)
		if (col.typmod - 4) == 1 {
			ok = typ.isslice(kindbyte, kindrune)
		}
		if !ok {
			ok = typ.isslice(kindstring) || typ.isslicen(2, kindbyte, kindrune)
		}
	case pgtyp_bytea:
		// []byte
		ok = typ.isslice(kindbyte)
	case pgtyp_byteaarr:
		// [][]byte
		ok = typ.isslicen(2, kindbyte)
	case pgtyp_cidr, pgtyp_inet:
		// *net.IPNet, string, *string
		ok = typ.is(kindstring) || typ.isnamed("net", "IPNet")
	case pgtyp_cidrarr, pgtyp_inetarr:
		// []*net.IPNet, []string
		ok = typ.isslice(kindstring) || (typ.kind == kindslice && typ.elem.isnamed("net", "IPNet"))
	case pgtyp_circle:
		// TODO(mkopriva): needs a custom type
		// - struct{ centerpoint [2]float64, radius float64 }
		return errors.UnsupportedColumnTypeError
	case pgtyp_circlearr:
		// TODO(mkopriva): needs a custom type
		// - []struct{ centerpoint [2]float64, radius float64 }
		return errors.UnsupportedColumnTypeError
	case pgtyp_date, pgtyp_time, pgtyp_timestamp, pgtyp_timestamptz, pgtyp_timetz:
		// time.Time, *time.Time
		ok = typ.istime || (typ.kind == kindptr && typ.elem.istime)
	case pgtyp_datearr, pgtyp_timearr, pgtyp_timestamparr, pgtyp_timestamptzarr, pgtyp_timetzarr:
		// []time.Time, []*time.Time
		ok = typ.kind == kindslice && (typ.elem.istime ||
			(typ.kind == kindptr && typ.elem.istime))
	case pgtyp_daterange, pgtyp_tsrange, pgtyp_tstzrange:
		// [2]time.Time (bound types expressed in tags?)
		ok = typ.kind == kindarray && typ.arraylen == 2 && typ.elem.istime
	case pgtyp_daterangearr, pgtyp_tsrangearr, pgtyp_tstzrangearr:
		// [][2]time.Time (bound types expressed in tags?)
		ok = typ.kind == kindslice && typ.elem.kind == kindarray &&
			typ.elem.arraylen == 2 && typ.elem.elem.istime
	case pgtyp_float4:
		// float32, *float32
		ok = typ.is(kindfloat32)
	case pgtyp_float4arr:
		// []float32
		ok = typ.isslice(kindfloat32)
	case pgtyp_float8:
		// float64, *float64
		ok = typ.is(kindfloat64)
	case pgtyp_float8arr:
		// []float64
		ok = typ.isslice(kindfloat64)
	case pgtyp_hstore:
		// map[string]sql.NullString, map[string]*string, map[string]string
		ok = typ.kind == kindmap && typ.key.kind == kindstring &&
			(typ.elem.is(kindstring) || typ.elem.isnamed("database/sql", "NullString"))
	case pgtyp_hstorearr:
		// []map[string]sql.NullString, []map[string]*string, []map[string]string
		ok = typ.kind == kindslice && typ.elem.kind == kindmap && typ.elem.key.kind == kindstring &&
			(typ.elem.elem.is(kindstring) || typ.elem.elem.isnamed("database/sql", "NullString"))
	case pgtyp_int2:
		// int16, int8, uint16, uint8, *int16, *int8, *uint16, *uint8
		ok = typ.is(kindint8, kindint16, kinduint16, kinduint8)
	case pgtyp_int2arr, pgtyp_int2vector:
		// []int16, []int8, []uint16, []uint8
		ok = typ.isslice(kindint8, kindint16, kinduint16, kinduint8)
	case pgtyp_int2vectorarr:
		// [][]int16, [][]int8, [][]uint16, [][]uint8
		ok = typ.isslicen(2, kindint8, kindint16, kinduint16, kinduint8)
	case pgtyp_int4:
		// int32, *int32
		ok = typ.is(kindint32)
	case pgtyp_int4arr:
		// []int32
		ok = typ.isslice(kindint32)
	case pgtyp_int4range:
		// [2]int32 (bound types expressed in tags?)
		ok = typ.kind == kindarray && typ.arraylen == 2 && typ.elem.kind == kindint32
	case pgtyp_int4rangearr:
		// [][2]int32 (bound types expressed in tags?)
		ok = typ.kind == kindslice && typ.elem.kind == kindarray &&
			typ.elem.arraylen == 2 && typ.elem.elem.kind == kindint32
	case pgtyp_int8:
		// int64, *int64
		ok = typ.is(kindint64)
	case pgtyp_int8arr:
		// []int64
		ok = typ.isslice(kindint64)
	case pgtyp_int8range:
		// [2]int64 (bound types expressed in tags?)
		ok = typ.kind == kindarray && typ.arraylen == 2 && typ.elem.kind == kindint64
	case pgtyp_int8rangearr:
		// [][2]int64 (bound types expressed in tags?)
		ok = typ.kind == kindslice && typ.elem.kind == kindarray &&
			typ.elem.arraylen == 2 && typ.elem.elem.kind == kindint64
	case pgtyp_interval:
		// TODO(mkopriva): needs a custom type
		// - struct{ microsecs int64, days int32, months int32 }
		return errors.UnsupportedColumnTypeError
	case pgtyp_intervalarr:
		// TODO(mkopriva): needs a custom type
		// - []struct{ microsecs int64, days int32, months int32 }
		return errors.UnsupportedColumnTypeError
	case pgtyp_json, pgtyp_jsonb:
		// []byte, json.RawMessage, json.Marshaler, usejson=true
		ok = typ.isslice(kindbyte) || (typ.isjsmarshaler && typ.isjsunmarshaler) || opts.usejson
	case pgtyp_jsonarr, pgtyp_jsonbarr:
		// [][]byte, []json.RawMessage, []json.Marshaler
		ok = typ.kind == kindslice && (typ.elem.isslice(kindbyte) ||
			(typ.elem.isjsmarshaler && typ.elem.isjsunmarshaler))
	case pgtyp_line:
		// [3]float64
		ok = typ.kind == kindarray && typ.arraylen == 3 && typ.elem.kind == kindfloat64
	case pgtyp_linearr:
		// [][3]float64
		ok = typ.kind == kindslice && typ.elem.kind == kindarray &&
			typ.elem.arraylen == 3 && typ.elem.elem.kind == kindfloat64
	case pgtyp_lseg:
		// [2][2]float64
		ok = typ.kind == kindarray && typ.arraylen == 2 &&
			typ.elem.kind == kindarray && typ.elem.arraylen == 2 &&
			typ.elem.elem.kind == kindfloat64
	case pgtyp_lsegarr:
		// [][2][2]float64
		if t := typ.elem; typ.kind == kindslice {
			ok = t.kind == kindarray && t.arraylen == 2 &&
				t.elem.kind == kindarray && t.elem.arraylen == 2 &&
				t.elem.elem.kind == kindfloat64
		}
	case pgtyp_macaddr, pgtyp_macaddr8:
		// []byte, net.HardwareAddr
		ok = typ.isslice(kindbyte)
	case pgtyp_macaddrarr, pgtyp_macaddr8arr:
		// []byte, []net.HardwareAddr
		ok = typ.isslicen(2, kindbyte)
	case pgtyp_money:
		// int64 (cents), *int64
		ok = typ.is(kindint64)
	case pgtyp_moneyarr:
		// []int64 (cents)
		ok = typ.isslice(kindint64)
	case pgtyp_numeric:
		// big.Int, *big.Int
		ok = typ.isnamed("math/big", "Int")
	case pgtyp_numericarr:
		// []big.Int, []*big.Int
		ok = typ.kind == kindslice && typ.elem.isnamed("math/big", "Int")
	case pgtyp_numrange:
		// [2]big.Int, [2]*big.Int (bound types expressed in tags?)
		ok = typ.kind == kindarray && typ.arraylen == 2 && typ.elem.isnamed("math/big", "Int")
	case pgtyp_numrangearr:
		// [][2]*big.Int (bound types expressed in tags?)
		ok = typ.kind == kindslice && typ.elem.kind == kindarray &&
			typ.elem.arraylen == 2 && typ.elem.elem.isnamed("math/big", "Int")
	case pgtyp_path, pgtyp_pointarr, pgtyp_polygon:
		// [][2]float64
		ok = typ.kind == kindslice && typ.elem.kind == kindarray &&
			typ.elem.arraylen == 2 && typ.elem.elem.kind == kindfloat64
	case pgtyp_patharr, pgtyp_polygonarr:
		// [][][2]float64
		if t := typ.elem; typ.kind == kindslice {
			ok = t.kind == kindslice && t.elem.kind == kindarray &&
				t.elem.arraylen == 2 && t.elem.elem.kind == kindfloat64
		}
	case pgtyp_point:
		// [2]float64
		ok = typ.kind == kindarray && typ.arraylen == 2 && typ.elem.kind == kindfloat64
	case pgtyp_text, pgtyp_tsquery:
		// string, *string
		ok = typ.is(kindstring)
	case pgtyp_textarr, pgtyp_tsqueryarr, pgtyp_tsvector:
		// []string
		ok = typ.isslice(kindstring)
	case pgtyp_tsvectorarr:
		// [][]string
		ok = typ.isslicen(2, kindstring)
	case pgtyp_uuid:
		// string, *string, []byte, [16]byte
		ok = typ.is(kindstring) || typ.isslice(kindbyte) ||
			(typ.kind == kindarray && typ.arraylen == 16 && typ.elem.kind == kindbyte)
	case pgtyp_uuidarr:
		// []string, [][]byte, [][16]byte
		ok = typ.isslice(kindstring) || typ.isslicen(2, kindbyte) ||
			(typ.kind == kindslice && typ.elem.kind == kindarray &&
				typ.elem.arraylen == 16 && typ.elem.elem.kind == kindbyte)
	case pgtyp_xml:
		// []byte, xml.Marshaler, usexml=true
		ok = typ.isslice(kindbyte) || (typ.isxmlmarshaler && typ.isxmlunmarshaler) || opts.usexml
	case pgtyp_xmlarr:
		// [][]byte, []xml.Marshaler
		ok = typ.kind == kindslice && (typ.elem.isslice(kindbyte) ||
			(typ.elem.isxmlmarshaler && typ.elem.isxmlunmarshaler))
	}

	if !ok {
		return nil // TODO
	}
	return nil
}

func (c *dbchecker) loadrelation(id relid) (*dbrelation, error) {
	rel := new(dbrelation)
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
	c.rel = rel

	// retrieve column info
	rows, err := c.db.Query(selectdbcolumns, rel.oid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		col := new(dbcolumn)
		err := rows.Scan(
			&col.num,
			&col.name,
			&col.typmod,
			&col.ndims,
			&col.hasnotnull,
			&col.hasdefault,
			&col.isprimary,
			&col.typ.name,
			&col.typ.size,
			&col.typ.typ,
			&col.typ.category,
			&col.typ.ispreferred,
			&col.typ.elem,
		)
		if err != nil {
			return nil, err
		}
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
		con := new(dbconstraint)
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
		ind := new(dbindex)
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

func (c *dbchecker) column(id colid) (*dbcolumn, error) {
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

type dbrelation struct {
	oid         int64  // The object identifier of the relation.
	name        string // The name of the relation.
	namespace   string // The name of the namespace to which the relation belongs.
	relkind     string // The relation's kind, we're only interested in r, v, and m.
	columns     []*dbcolumn
	constraints []*dbconstraint
	indexes     []*dbindex
}

func (rel *dbrelation) column(name string) *dbcolumn {
	for _, col := range rel.columns {
		if col.name == name {
			return col
		}
	}
	return nil
}

func (rel *dbrelation) constraint(name string) *dbconstraint {
	for _, con := range rel.constraints {
		if con.name == name {
			return con
		}
	}
	return nil
}

func (rel *dbrelation) index(name string) *dbindex {
	for _, ind := range rel.indexes {
		if ind.name == name {
			return ind
		}
	}
	return nil
}

type dbcolumn struct {
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
	// Info about the column's type.
	typ dbtype
}

type dbtype struct {
	// The name of the type.
	name string
	// The number of bytes for fixed-size types, negative for variable length types.
	size int
	// The type's type.
	typ string
	// An arbitrary classification of data types that is used by the parser
	// to determine which implicit casts should be "preferred".
	category string
	// True if the type is a preferred cast target within its category.
	ispreferred bool
	// If this is an array type then elem identifies the element type
	// of that array type.
	elem int64
}

func (typ *dbtype) is(names ...string) bool {
	for _, name := range names {
		if typ.name == name {
			return true
		}
	}
	return false
}

type dbconstraint struct {
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

type dbindex struct {
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

type dbid struct {
	name      string
	namespace string
}

// postgres types
const (
	pgtyp_bit       = "bit"
	pgtyp_bitarr    = "_bit"
	pgtyp_bool      = "bool"
	pgtyp_boolarr   = "_bool"
	pgtyp_box       = "box"
	pgtyp_boxarr    = "_box"
	pgtyp_bpchar    = "bpchar"
	pgtyp_bpchararr = "_bpchar"
	pgtyp_bytea     = "bytea"
	pgtyp_byteaarr  = "_bytea"

	pgtyp_char      = "char"
	pgtyp_chararr   = "_char"
	pgtyp_cidr      = "cidr"
	pgtyp_cidrarr   = "_cidr"
	pgtyp_circle    = "circle"
	pgtyp_circlearr = "_circle"

	pgtyp_date         = "date"
	pgtyp_datearr      = "_date"
	pgtyp_daterange    = "daterange"
	pgtyp_daterangearr = "_daterange"

	pgtyp_float4    = "float4"
	pgtyp_float4arr = "_float4"
	pgtyp_float8    = "float8"
	pgtyp_float8arr = "_float8"

	pgtyp_hstore    = "hstore"
	pgtyp_hstorearr = "_hstore"

	pgtyp_inet          = "inet"
	pgtyp_inetarr       = "_inet"
	pgtyp_int2          = "int2"
	pgtyp_int2arr       = "_int2"
	pgtyp_int2vector    = "int2vector"
	pgtyp_int2vectorarr = "_int2vector"
	pgtyp_int4          = "int4"
	pgtyp_int4arr       = "_int4"
	pgtyp_int4range     = "int4range"
	pgtyp_int4rangearr  = "_int4range"
	pgtyp_int8          = "int8"
	pgtyp_int8arr       = "_int8"
	pgtyp_int8range     = "int8range"
	pgtyp_int8rangearr  = "_int8range"
	pgtyp_interval      = "interval"
	pgtyp_intervalarr   = "_interval"

	pgtyp_json     = "json"
	pgtyp_jsonarr  = "_json"
	pgtyp_jsonb    = "jsonb"
	pgtyp_jsonbarr = "_jsonb"

	pgtyp_line    = "line"
	pgtyp_linearr = "_line"
	pgtyp_lseg    = "lseg"
	pgtyp_lsegarr = "_lseg"

	pgtyp_macaddr     = "macaddr"
	pgtyp_macaddrarr  = "_macaddr"
	pgtyp_macaddr8    = "macaddr8"
	pgtyp_macaddr8arr = "_macaddr8"
	pgtyp_money       = "money"
	pgtyp_moneyarr    = "_money"

	pgtyp_numeric     = "numeric"
	pgtyp_numericarr  = "_numeric"
	pgtyp_numrange    = "numrange"
	pgtyp_numrangearr = "_numrange"

	pgtyp_path       = "path"
	pgtyp_patharr    = "_path"
	pgtyp_point      = "point"
	pgtyp_pointarr   = "_point"
	pgtyp_polygon    = "polygon"
	pgtyp_polygonarr = "_polygon"

	pgtyp_text           = "text"
	pgtyp_textarr        = "_text"
	pgtyp_time           = "time"
	pgtyp_timestamp      = "timestamp"
	pgtyp_timestamparr   = "_timestamp"
	pgtyp_timearr        = "_time"
	pgtyp_timestamptz    = "timestamptz"
	pgtyp_timestamptzarr = "_timestamptz"
	pgtyp_timetz         = "timetz"
	pgtyp_timetzarr      = "_timetz"
	pgtyp_tsrange        = "tsrange"
	pgtyp_tsrangearr     = "_tsrange"
	pgtyp_tstzrange      = "tstzrange"
	pgtyp_tstzrangearr   = "_tstzrange"

	pgtyp_tsquery     = "tsquery"
	pgtyp_tsqueryarr  = "_tsquery"
	pgtyp_tsvector    = "tsvector"
	pgtyp_tsvectorarr = "_tsvector"

	pgtyp_uuid    = "uuid"
	pgtyp_uuidarr = "_uuid"

	pgtyp_varbit     = "varbit"
	pgtyp_varbitarr  = "_varbit"
	pgtyp_varchar    = "varchar"
	pgtyp_varchararr = "_varchar"

	pgtyp_xml    = "xml"
	pgtyp_xmlarr = "_xml"
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
