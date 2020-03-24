package convert

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/frk/compare"
	_ "github.com/lib/pq"
)

var testdb *convtest

func TestMain(m *testing.M) {
	var exitcode int

	func() { // use a func wrapper so we can rely on defer
		testdb = new(convtest)
		defer testdb.close()

		if err := testdb.init(); err != nil {
			panic(err)
		}

		exitcode = m.Run()
	}()

	os.Exit(exitcode)
}

type convtest struct {
	root   *sql.DB
	db     *sql.DB
	dbname string
	dburl  string
}

func (t *convtest) init() (err error) {
	// open the default db
	if t.root, err = sql.Open("postgres", "postgres:///?sslmode=disable"); err != nil {
		return err
	} else if err = t.root.Ping(); err != nil {
		return err
	}

	// create a new database so that the default one isn't polluted with the test data.
	t.dbname = "convert_test_db"
	t.dburl = "postgres:///" + t.dbname + "?sslmode=disable"
	if _, err = t.root.Exec("DROP DATABASE IF EXISTS " + t.dbname); err != nil {
		return err
	}
	if _, err = t.root.Exec("CREATE DATABASE " + t.dbname); err != nil {
		return err
	}

	// open the new database
	if t.db, err = sql.Open("postgres", t.dburl); err != nil {
		return err
	} else if err = t.db.Ping(); err != nil {
		return err
	}

	// populate test db
	const populatedbquery = `
CREATE EXTENSION hstore;
CREATE TABLE coltype_test (
	id serial PRIMARY KEY

	-- scalar types
	, col_bit         bit          DEFAULT '1'
	, col_bool        boolean      DEFAULT true
	, col_bpchar      bpchar       DEFAULT 'a'
	, col_bpchar3     bpchar(3)    DEFAULT 'abc'
	, col_char        char         DEFAULT 'a'
	, col_char3       char(3)      DEFAULT 'abc'
	, col_cidr        cidr         DEFAULT '192.168.100.128/25'
	, col_date        date         DEFAULT '1999-01-08'
	, col_float4      float4       DEFAULT 0.1
	, col_float8      float8       DEFAULT 0.1
	, col_inet        inet         DEFAULT '192.168.0.1/24'
	, col_int2        int2         DEFAULT 1
	, col_int4        int4         DEFAULT 1
	, col_int8        int8         DEFAULT 1
	, col_interval    interval     DEFAULT '1 day'
	, col_macaddr     macaddr      DEFAULT '08:00:2b:01:02:03'
	, col_macaddr8    macaddr8     DEFAULT '08:00:2b:01:02:03:04:05'
	, col_money       money        DEFAULT '$1.20'
	, col_numeric     numeric      DEFAULT 1.2
	, col_text        text         DEFAULT 'foo'
	, col_time        time         DEFAULT '04:05:06.789'
	, col_timestamp   timestamp    DEFAULT '1999-01-08 04:05:06'
	, col_timestamptz timestamptz  DEFAULT 'January 8 04:05:06 1999 PST'
	, col_timetz      timetz       DEFAULT '04:05:06.789-8'
	, col_tsquery     tsquery      DEFAULT 'fat & rat'
	, col_tsvector    tsvector     DEFAULT 'a fat cat sat on a mat and ate a fat rat'
	, col_uuid        uuid         DEFAULT 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'
	, col_varbit      varbit       DEFAULT '101'
	, col_varbit1     varbit(1)    DEFAULT '1'
	, col_varchar     varchar      DEFAULT 'foo'
	, col_varchar1    varchar(1)   DEFAULT 'a'

	-- aggregate types
	, col_bitarr         bit[]          DEFAULT '{0,1,0,1}'::bit[]
	, col_boolarr        boolean[]      DEFAULT '{t,f}'::boolean[]
	, col_boxarr         box[]          DEFAULT '{(0,1),(0.5,1.5);(0,1),(0.5,1.5)}'::box[]
	, col_bpchararr      bpchar[]       DEFAULT ARRAY['a','b']::bpchar[]
	, col_bpchar3arr     bpchar(3)[]    DEFAULT ARRAY['abc','def']::bpchar(3)[]
	, col_bytea          bytea          DEFAULT '\xDEADBEEF'
	, col_byteaarr       bytea[]        DEFAULT ARRAY['\xDEADBEEF', '\xDEADBEEF']::bytea[]
	, col_chararr        char[]         DEFAULT ARRAY['a','b']::char[]
	, col_char3arr       char(3)[]      DEFAULT ARRAY['abc','def']::char(3)[]
	, col_cidrarr        cidr[]         DEFAULT ARRAY['192.168.100.128/25','2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128']::cidr[]
	, col_circlearr      circle[]       DEFAULT ARRAY['<(0,0), 3.5>','((0.5,1), 5)']::circle[]
	, col_datearr        date[]         DEFAULT ARRAY['1999-01-08', 'May 5, 2001']::date[]
	, col_daterangearr   daterange[]    DEFAULT ARRAY['[1999-01-08, 2001-01-08)', '(1999-01-08, 2001-01-08]']::daterange[]
	, col_float4arr      float4[]       DEFAULT ARRAY[0.1, 0.2]::float4[]
	, col_float8arr      float8[]       DEFAULT ARRAY[0.1, 0.2]::float8[]
	, col_hstorearr      hstore[]       DEFAULT ARRAY['a=>1,b=>2', 'c=>3,d=>4']::hstore[]
	, col_inetarr        inet[]         DEFAULT ARRAY['192.168.0.1/24', '128.0.0.0/16']::inet[]
	, col_int2arr        int2[]         DEFAULT ARRAY[1,2]::int2[]
	, col_int2vector     int2vector     DEFAULT '1 2'
	, col_int2vectorarr  int2vector[]   DEFAULT ARRAY['1 2', '3 4']::int2vector[]
	, col_int4arr        int4[]         DEFAULT ARRAY[1,2]::int4[]
	, col_int4rangearr   int4range[]    DEFAULT ARRAY['[0,9)', '(0,9]']::int4range[]
	, col_int8arr        int8[]         DEFAULT ARRAY[1,2]::int8[]
	, col_int8rangearr   int8range[]    DEFAULT ARRAY['[0,9)', '(0,9]']::int8range[]
	, col_intervalarr    interval[]     DEFAULT ARRAY['1 day','5 years 4 months 34 minutes ago']::interval[]
	, col_jsonarr        json[]         DEFAULT ARRAY['{"foo":["bar", "baz", 123]}', '["foo", 123]']::json[]
	, col_jsonbarr       jsonb[]        DEFAULT ARRAY['{"foo":["bar", "baz", 123]}', '["foo", 123]']::jsonb[]
	, col_linearr        line[]         DEFAULT ARRAY['{1,2,3}', '{4,5,6}']::line[]
	, col_lsegarr        lseg[]         DEFAULT ARRAY['[(1,2), (3,4)]', '[(1.5,2.5), (3.5,4.5)]']::lseg[]
	, col_macaddrarr     macaddr[]      DEFAULT ARRAY['08:00:2b:01:02:03', '08002b010203']::macaddr[]
	, col_macaddr8arr    macaddr8[]     DEFAULT ARRAY['08:00:2b:01:02:03:04:05', '08002b0102030405']::macaddr8[]
	, col_moneyarr       money[]        DEFAULT ARRAY['$1.20', '$0.99']::money[]
	, col_numericarr     numeric[]      DEFAULT ARRAY[1.2,3.4]::numeric[]
	, col_numrangearr    numrange[]     DEFAULT ARRAY['[1.2,3.4)', '(1.2,3.4]']::numrange[]
	, col_patharr        path[]         DEFAULT ARRAY['[(1,1),(2,2),(3,3)]', '[(1.5,1.5),(2.5,2.5),(3.5,3.5)]']::path[]
	, col_pointarr       point[]        DEFAULT ARRAY['(1,1)', '(2,2)']::point[]
	, col_polygonarr     polygon[]      DEFAULT ARRAY['((1,1),(2,2),(3,3))', '((1.5,1.5),(2.5,2.5),(3.5,3.5))']::polygon[]
	, col_textarr        text[]         DEFAULT ARRAY['foo', 'bar']::text[]
	, col_timearr        time[]         DEFAULT ARRAY['04:05:06.789', '040506']::time[]
	, col_timestamparr   timestamp[]    DEFAULT ARRAY['1999-01-08 04:05:06', '2004-10-19 10:23:54']::timestamp[]
	, col_timestamptzarr timestamptz[]  DEFAULT ARRAY['January 8 04:05:06 1999 PST','2004-10-19 10:23:54+02']::timestamptz[]
	, col_timetzarr      timetz[]       DEFAULT ARRAY['04:05:06.789-8','2003-04-12 04:05:06 America/New_York']::timetz[]
	, col_tsqueryarr     tsquery[]      DEFAULT ARRAY['fat & rat', 'fat & rat & ! cat']::tsquery[]
	, col_tsrangearr     tsrange[]      DEFAULT ARRAY['[1999-01-08 04:05:06, 2004-10-19 10:23:54)', '(1999-01-08 04:05:06, 2004-10-19 10:23:54]']::tsrange[]
	, col_tstzrangearr   tstzrange[]    DEFAULT ARRAY['[January 8 04:05:06 1999 PST, 2004-10-19 10:23:54+02)','(January 8 04:05:06 1999 PST, 2004-10-19 10:23:54+02]']::tstzrange[]
	, col_tsvectorarr    tsvector[]     DEFAULT ARRAY['a fat cat sat on a mat','and ate a fat rat']::tsvector[]
	, col_uuidarr        uuid[]         DEFAULT ARRAY['a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11','a0eebc999c0b4ef8bb6d6bb9bd380a11']::uuid[]
	, col_varbitarr      varbit[]       DEFAULT '{101, 00}'::varbit[]
	, col_varbit1arr     varbit(1)[]    DEFAULT '{1, 0}'::varbit(1)[]
	, col_varchararr     varchar[]      DEFAULT ARRAY['foo', 'bar']::varchar[]
	, col_varchar1arr    varchar(1)[]   DEFAULT ARRAY['a', 'b']::varchar(1)[]
	, col_xmlarr         xml[]          DEFAULT ARRAY['<foo>bar</foo>','<bar>foo</bar>']::xml[]

	-- composite types
	, col_box     box     DEFAULT '(0.5,100.5),(0.5,10.5)'::box
	, col_circle  circle  DEFAULT '<(0,0), 3.5>'
	, col_hstore  hstore  DEFAULT 'a=>1,b=>2'
	, col_line    line    DEFAULT '{1,2,3}'
	, col_lseg    lseg    DEFAULT '[(1,2), (3,4)]'
	, col_path    path    DEFAULT '[(1,1),(2,2),(3,3)]'
	, col_point   point   DEFAULT '(1,1)'
	, col_polygon polygon DEFAULT '((1,1),(2,2),(3,3))'

	-- range types
	, col_daterange daterange  DEFAULT '[1999-01-08, 2001-01-08)'
	, col_int4range int4range  DEFAULT '[0,9)'
	, col_int8range int8range  DEFAULT '[0,9)'
	, col_numrange  numrange   DEFAULT '[1.2,3.4)'
	, col_tsrange   tsrange    DEFAULT '[1999-01-08 04:05:06, 2004-10-19 10:23:54)'
	, col_tstzrange tstzrange  DEFAULT '[January 8 04:05:06 1999 PST, 2004-10-19 10:23:54+02)'

	-- "format" types
	, col_json  json  DEFAULT '{"foo":["bar", "baz", 123]}'
	, col_jsonb jsonb DEFAULT '{"foo":["bar", "baz", 123]}'
	, col_xml   xml   DEFAULT '<foo>bar</foo>'
);
` //`

	if _, err = t.db.Exec(populatedbquery); err != nil {
		return err
	}

	return nil
}

func (t *convtest) close() {
	if t.db != nil {
		if err := t.db.Close(); err != nil {
			log.Println("error closing test db handle:", err)
		}
	}
	if t.root != nil {
		if _, err := t.root.Exec("DROP DATABASE " + t.dbname); err != nil {
			log.Println("error dropping test db:", err)
		}
		if err := t.root.Close(); err != nil {
			log.Println("error closing root db handle:", err)
		}
	}
}

type test_scanner []struct {
	// scnr is called per each row in the rows slice and returns dest and result:
	// - dest - is used by the test runner as the destination to scan the
	//   selected row's column.
	// - result - is the the pointer that points to the value that the scanner
	//   populates with the data of the column, it is then used to compare
	//   against the value specified into a row's want field.
	scanner func() (scanner, result interface{})
	rows    []test_scanner_row
}

type test_scanner_row struct {
	typ  string
	in   interface{}
	want interface{}
}

func (table test_scanner) execute(t *testing.T) {
	for _, tt := range table {
		for _, r := range tt.rows {
			scanner, got := tt.scanner()
			if scanner == nil {
				scanner = got
			}

			name := fmt.Sprintf("%T::<%s>_to_<%T>_using_%v", scanner, r.typ, r.want, r.in)
			t.Run(name, func(t *testing.T) {
				var id int
				var col = "col_" + r.typ
				row := testdb.db.QueryRow(`insert into coltype_test (`+col+`) values ($1) returning id`, r.in)
				if err := row.Scan(&id); err != nil {
					t.Fatal(err)
				}

				row = testdb.db.QueryRow(`select `+col+` from coltype_test where id = $1`, id)
				if err := row.Scan(scanner); err != nil {
					t.Error(err)
				} else if e := compare.Compare(got, r.want); e != nil {
					t.Error(e)
				}
			})
		}
	}
}

type test_valuer []struct {
	valuer func() interface{}
	rows   []test_valuer_row
}

type test_valuer_row struct {
	typ  string
	in   interface{}
	want *string
}

func (table test_valuer) execute(t *testing.T) {
	for _, tt := range table {
		for _, r := range tt.rows {
			valuer := tt.valuer()
			if valuer != nil {
				rv := reflect.ValueOf(valuer).Elem()
				if fv := reflect.ValueOf(r.in); fv.IsValid() {
					rv.Field(0).Set(fv)
				}
			} else {
				valuer = r.in
			}

			name := fmt.Sprintf("%T::<%s>_from_<%T>_using_%v", valuer, r.typ, r.in, r.in)
			t.Run(name, func(t *testing.T) {
				var id int
				var col = "col_" + r.typ
				row := testdb.db.QueryRow(`insert into coltype_test (`+col+`) values ($1) returning id`, valuer)
				if err := row.Scan(&id); err != nil {
					t.Fatal(err)
				}

				var dest *string
				row = testdb.db.QueryRow(`select `+col+`::text from coltype_test where id = $1`, id)
				if err := row.Scan(&dest); err != nil {
					t.Error(err)
				} else if e := compare.Compare(dest, r.want); e != nil {
					t.Error(e)
				}
			})
		}
	}
}

// helper
func strptr(v string) *string   { return &v }
func bytesptr(v string) *[]byte { vv := []byte(v); return &vv }
func byteptr(v byte) *byte      { return &v }
func runeptr(v rune) *rune      { return &v }
func boolptr(v bool) *bool      { return &v }
func uptr(v uint) *uint         { return &v }
func u8ptr(v uint8) *uint8      { return &v }
