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
					if err := c.checktype(col, node.typ); err != nil {
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

func (c *dbchecker) checktype(col *dbcolumn, typ typeinfo) error {
	var ok bool

	// TODO

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
	pgtyp_bool           = "bool"         // bool
	pgtyp_bytea          = "bytea"        // []byte
	pgtyp_char           = "char"         // byte, rune, string, uint8 - if column.typmod=1; else - string
	pgtyp_int8           = "int8"         // int64
	pgtyp_int2           = "int2"         // int16, int8
	pgtyp_int2vector     = "int2vector"   // []int16, int8
	pgtyp_int4           = "int4"         // int32
	pgtyp_text           = "text"         // string
	pgtyp_json           = "json"         // []byte, json.RawMessage, json.Marshaler, usejson=true
	pgtyp_xml            = "xml"          // []byte, xml.Marshaler, usexml=true
	pgtyp_xmlarr         = "_xml"         // [][]byte, []xml.Marshaler, usexml=true
	pgtyp_jsonarr        = "_json"        // [][]byte, []json.RawMessage, []json.Marshaler, usejson=true
	pgtyp_point          = "point"        // [2]float64
	pgtyp_lseg           = "lseg"         // [2][2]float64
	pgtyp_path           = "path"         // [][2]float64
	pgtyp_box            = "box"          // [2][2]float64
	pgtyp_polygon        = "polygon"      // [][2]float64
	pgtyp_line           = "line"         // [3]float64
	pgtyp_linearr        = "_line"        // [][3]float64
	pgtyp_cidr           = "cidr"         // *net.IPNet, string
	pgtyp_cidrarr        = "_cidr"        // []*net.IPNet, string
	pgtyp_float4         = "float4"       // float32
	pgtyp_float8         = "float8"       // float64
	pgtyp_circle         = "circle"       // { centerpoint [2]float64, radius float64 }
	pgtyp_circlearr      = "_circle"      // []{ centerpoint [2]float64, radius float64 }
	pgtyp_macaddr8       = "macaddr8"     // net.HardwareAddr
	pgtyp_macaddr8arr    = "_macaddr8"    // []net.HardwareAddr
	pgtyp_money          = "money"        // int64 (cents)
	pgtyp_moneyarr       = "_money"       // []int64 (cents)
	pgtyp_macaddr        = "macaddr"      // net.HardwareAddr
	pgtyp_inet           = "inet"         // *net.IPNet
	pgtyp_boolarr        = "_bool"        // []bool
	pgtyp_byteaarr       = "_bytea"       // [][]byte
	pgtyp_chararr        = "_char"        // []byte, []rune, []string, []uint8
	pgtyp_int2arr        = "_int2"        // []int16, []int8
	pgtyp_int2vectorarr  = "_int2vector"  // [][]int16, []int8
	pgtyp_int4arr        = "_int4"        // []int32
	pgtyp_textarr        = "_text"        // []string
	pgtyp_bpchararr      = "_bpchar"      // []string
	pgtyp_varchararr     = "_varchar"     // []string
	pgtyp_int8arr        = "_int8"        // []int64
	pgtyp_pointarr       = "_point"       // [][2]float64
	pgtyp_lsegarr        = "_lseg"        // [][2][2]float64
	pgtyp_patharr        = "_path"        // [][][2]float64
	pgtyp_boxarr         = "_box"         // [][2][2]float64
	pgtyp_float4arr      = "_float4"      // []float32
	pgtyp_float8arr      = "_float8"      // []float64
	pgtyp_polygonarr     = "_polygon"     // [][][2]float64
	pgtyp_macaddrarr     = "_macaddr"     // []net.HardwareAddr
	pgtyp_inetarr        = "_inet"        // []*net.IPNet
	pgtyp_bpchar         = "bpchar"       // string (blank-padded char [the internal name of the character data type])
	pgtyp_varchar        = "varchar"      // string
	pgtyp_date           = "date"         // time.Time
	pgtyp_time           = "time"         // time.Time
	pgtyp_timestamp      = "timestamp"    // time.Time
	pgtyp_timestamparr   = "_timestamp"   // []time.Time
	pgtyp_datearr        = "_date"        // []time.Time
	pgtyp_timearr        = "_time"        // []time.Time
	pgtyp_timestamptz    = "timestamptz"  // time.Time
	pgtyp_timestamptzarr = "_timestamptz" // []time.Time
	pgtyp_interval       = "interval"     // {microsecs int64, days int32, months int32}
	pgtyp_intervalarr    = "_interval"    // []{microsecs int64, days int32, months int32}
	pgtyp_numericarr     = "_numeric"     // ? []*big.Int
	pgtyp_timetz         = "timetz"       // time.Time
	pgtyp_timetzarr      = "_timetz"      // []time.Time
	pgtyp_bit            = "bit"          // string, []byte
	pgtyp_bitarr         = "_bit"         // []string, [][]byte
	pgtyp_varbit         = "varbit"       // string, []byte
	pgtyp_varbitarr      = "_varbit"      // []string, [][]byte
	pgtyp_numeric        = "numeric"      // ? *big.Int
	pgtyp_uuid           = "uuid"         // string, [16]byte, []byte
	pgtyp_uuidarr        = "_uuid"        // []string, [][16]byte, [][]byte
	pgtyp_tsvector       = "tsvector"     // []string
	pgtyp_tsquery        = "tsquery"      // string
	pgtyp_tsvectorarr    = "_tsvector"    // [][]string
	pgtyp_tsqueryarr     = "_tsquery"     // []string
	pgtyp_jsonb          = "jsonb"        // []byte, json.RawMessage, json.Marshaler, usejson=true
	pgtyp_jsonbarr       = "_jsonb"       // [][]byte, []json.RawMessage, []json.Marshaler, usejson=true
	pgtyp_int4range      = "int4range"    // [2]int32 (bound types expressed in tags?)
	pgtyp_int4rangearr   = "_int4range"   // [][2]int32 (bound types expressed in tags?)
	pgtyp_numrange       = "numrange"     // [2]*big.Int (bound types expressed in tags?)
	pgtyp_numrangearr    = "_numrange"    // [][2]*big.Int (bound types expressed in tags?)
	pgtyp_tsrange        = "tsrange"      // [2]time.Time (bound types expressed in tags?)
	pgtyp_tsrangearr     = "_tsrange"     // [][2]time.Time (bound types expressed in tags?)
	pgtyp_tstzrange      = "tstzrange"    // [2]time.Time (bound types expressed in tags?)
	pgtyp_tstzrangearr   = "_tstzrange"   // [][2]time.Time (bound types expressed in tags?)
	pgtyp_daterange      = "daterange"    // [2]time.Time (bound types expressed in tags?)
	pgtyp_daterangearr   = "_daterange"   // [][2]time.Time (bound types expressed in tags?)
	pgtyp_int8range      = "int8range"    // [2]int64 (bound types expressed in tags?)
	pgtyp_int8rangearr   = "_int8range"   // [][2]int64 (bound types expressed in tags?)
	pgtyp_hstore         = "hstore"       // map[string]sql.NullString, map[string]*string, map[string]string
	pgtyp_hstorearr      = "_hstore"      // []map[string]sql.NullString, []map[string]*string, map[string]string
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
