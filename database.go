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
	if c.cmd.textsearch != nil {
		ts := *c.cmd.textsearch
		rel, ok := c.relmap[ts.qual]
		if !ok {
			return errors.NoDBRelationError
		}
		col := rel.column(ts.name)
		if col == nil {
			return errors.NoDBColumnError
		}
		if col.typ.name != pgtyp_tsvector {
			return errors.BadDBColumnTypeError
		}
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
	num  int    // The number of the column. Ordinary columns are numbered from 1 up.
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
	pgtyp_bool           = "bool"
	pgtyp_bytea          = "bytea"
	pgtyp_char           = "char"
	pgtyp_int8           = "int8"
	pgtyp_int2           = "int2"
	pgtyp_int2vector     = "int2vector"
	pgtyp_int4           = "int4"
	pgtyp_text           = "text"
	pgtyp_json           = "json"
	pgtyp_xml            = "xml"
	pgtyp_xmlarr         = "_xml"
	pgtyp_jsonarr        = "_json"
	pgtyp_point          = "point"
	pgtyp_lseg           = "lseg"
	pgtyp_path           = "path"
	pgtyp_box            = "box"
	pgtyp_polygon        = "polygon"
	pgtyp_line           = "line"
	pgtyp_linearr        = "_line"
	pgtyp_cidr           = "cidr"
	pgtyp_cidrarr        = "_cidr"
	pgtyp_float4         = "float4"
	pgtyp_float8         = "float8"
	pgtyp_abstime        = "abstime"
	pgtyp_reltime        = "reltime"
	pgtyp_tinterval      = "tinterval"
	pgtyp_unknown        = "unknown"
	pgtyp_circle         = "circle"
	pgtyp_circlearr      = "_circle"
	pgtyp_macaddr8       = "macaddr8"
	pgtyp_macaddr8arr    = "_macaddr8"
	pgtyp_money          = "money"
	pgtyp_moneyarr       = "_money"
	pgtyp_macaddr        = "macaddr"
	pgtyp_inet           = "inet"
	pgtyp_boolarr        = "_bool"
	pgtyp_byteaarr       = "_bytea"
	pgtyp_chararr        = "_char"
	pgtyp_int2arr        = "_int2"
	pgtyp_int2vectorarr  = "_int2vector"
	pgtyp_int4arr        = "_int4"
	pgtyp_textarr        = "_text"
	pgtyp_bpchararr      = "_bpchar"
	pgtyp_varchararr     = "_varchar"
	pgtyp_int8arr        = "_int8"
	pgtyp_pointarr       = "_point"
	pgtyp_lsegarr        = "_lseg"
	pgtyp_patharr        = "_path"
	pgtyp_boxarr         = "_box"
	pgtyp_float4arr      = "_float4"
	pgtyp_float8arr      = "_float8"
	pgtyp_abstimearr     = "_abstime"
	pgtyp_reltimearr     = "_reltime"
	pgtyp_tintervalarr   = "_tinterval"
	pgtyp_polygonarr     = "_polygon"
	pgtyp_macaddrarr     = "_macaddr"
	pgtyp_inetarr        = "_inet"
	pgtyp_bpchar         = "bpchar" // blank-padded char (the internal name of the character data type)
	pgtyp_varchar        = "varchar"
	pgtyp_date           = "date"
	pgtyp_time           = "time"
	pgtyp_timestamp      = "timestamp"
	pgtyp_timestamparr   = "_timestamp"
	pgtyp_datearr        = "_date"
	pgtyp_timearr        = "_time"
	pgtyp_timestamptz    = "timestamptz"
	pgtyp_timestamptzarr = "_timestamptz"
	pgtyp_interval       = "interval"
	pgtyp_intervalarr    = "_interval"
	pgtyp_numericarr     = "_numeric"
	pgtyp_cstringarr     = "_cstring"
	pgtyp_timetz         = "timetz"
	pgtyp_timetzarr      = "_timetz"
	pgtyp_bit            = "bit"
	pgtyp_bitarr         = "_bit"
	pgtyp_varbit         = "varbit"
	pgtyp_varbitarr      = "_varbit"
	pgtyp_numeric        = "numeric"
	pgtyp_cstring        = "cstring"
	pgtyp_any            = "any"
	pgtyp_anyarray       = "anyarray"
	pgtyp_void           = "void"
	pgtyp_anyelement     = "anyelement"
	pgtyp_anynonarray    = "anynonarray"
	pgtyp_uuid           = "uuid"
	pgtyp_uuidarr        = "_uuid"
	pgtyp_anyenum        = "anyenum"
	pgtyp_tsvector       = "tsvector"
	pgtyp_tsquery        = "tsquery"
	pgtyp_tsvectorarr    = "_tsvector"
	pgtyp_tsqueryarr     = "_tsquery"
	pgtyp_jsonb          = "jsonb"
	pgtyp_jsonbarr       = "_jsonb"
	pgtyp_anyrange       = "anyrange"
	pgtyp_int4range      = "int4range"
	pgtyp_int4rangearr   = "_int4range"
	pgtyp_numrange       = "numrange"
	pgtyp_numrangearr    = "_numrange"
	pgtyp_tsrange        = "tsrange"
	pgtyp_tsrangearr     = "_tsrange"
	pgtyp_tstzrange      = "tstzrange"
	pgtyp_tstzrangearr   = "_tstzrange"
	pgtyp_daterange      = "daterange"
	pgtyp_daterangearr   = "_daterange"
	pgtyp_int8range      = "int8range"
	pgtyp_int8rangearr   = "_int8range"
	pgtyp_hstore         = "hstore"
	pgtyp_hstorearr      = "_hstore"
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
