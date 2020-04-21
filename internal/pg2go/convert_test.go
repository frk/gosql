package pg2go

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"testing"
	"time"

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

type testlist []struct {
	valuer func() interface{}
	// scanner is called per each row in the rows slice and returns scanner and result:
	// - scanner - is used by the test runner as the destination to scan the
	//   selected row's column.
	// - result - is the the pointer that points to the value that the scanner
	//   populates with the data of the column, it is then used to compare
	//   against the value specified into a row's want field.
	scanner func() (scanner, result interface{})
	data    []testdata
}

type testdata struct {
	input  interface{}
	output interface{}
}

func (list testlist) execute(t *testing.T, coltype string) {
	for _, tt := range list {
		for _, td := range tt.data {
			valuer := tt.valuer()
			if valuer != nil {
				rv := reflect.ValueOf(valuer).Elem()
				if fv := reflect.ValueOf(td.input); fv.IsValid() {
					rv.Field(0).Set(fv)
				}
			} else {
				valuer = td.input
			}

			scanner, dest := tt.scanner()
			if scanner == nil {
				scanner = dest
			}

			name := fmt.Sprintf("%T->%T::<%s>_FROM_<%T>_TO_<%T>_USING_%v", valuer, scanner, coltype, td.input, td.output, td.input)
			t.Run(name, func(t *testing.T) {
				var id int
				var col = "col_" + coltype
				row := testdb.db.QueryRow(`INSERT INTO coltype_test (`+col+`) VALUES ($1) RETURNING id`, valuer)
				if err := row.Scan(&id); err != nil {
					t.Fatal(err)
				}

				row = testdb.db.QueryRow(`SELECT `+col+` FROM coltype_test WHERE id = $1`, id)
				if err := row.Scan(scanner); err != nil {
					t.Error(err)
				} else {
					got := reflect.ValueOf(dest).Elem().Interface()
					_ = got

					//if e := compare.Compare(dest, td.output); e != nil {
					if e := compare.Compare(got, td.output); e != nil {
						t.Error(e)
					}
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

func iptr(v int) *int       { return &v }
func i8ptr(v int8) *int8    { return &v }
func i16ptr(v int16) *int16 { return &v }
func i32ptr(v int32) *int32 { return &v }
func i64ptr(v int64) *int64 { return &v }

func uptr(v uint) *uint       { return &v }
func u8ptr(v uint8) *uint8    { return &v }
func u16ptr(v uint16) *uint16 { return &v }
func u32ptr(v uint32) *uint32 { return &v }
func u64ptr(v uint64) *uint64 { return &v }

func f32ptr(v float32) *float32 { return &v }
func f64ptr(v float64) *float64 { return &v }

func uuid16bytes(v string) [16]byte {
	u, _ := pgParseUUID([]byte(v))
	return u
}

func dateval(y, m, d int) time.Time {
	t := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
	return t
}

func timeval(h, m, s, ms int) time.Time {
	t := time.Date(0, 1, 1, h, m, s, ms*1000*1000, time.UTC)
	return t
}

func timetzval(h, m, s, ms int, loc *time.Location) time.Time {
	t := time.Date(0, 1, 1, h, m, s, ms*1000*1000, loc)
	return t
}

func timestamp(y, m, d, hh, mm, ss, ms int) time.Time {
	return time.Date(y, time.Month(m), d, hh, mm, ss, ms*1000*1000, noZone)
}

func timestamptz(y, m, d, hh, mm, ss, ms int, loc *time.Location) time.Time {
	return time.Date(y, time.Month(m), d, hh, mm, ss, ms*1000*1000, loc)
}

func dateptr(y, m, d int) *time.Time {
	t := dateval(y, m, d)
	return &t
}

func netCIDR(v string) net.IPNet {
	_, n, _ := net.ParseCIDR(v)
	return *n
}

func netCIDRptr(v string) *net.IPNet {
	n := netCIDR(v)
	return &n
}

func netCIDRSlice(vv ...string) []net.IPNet {
	out := make([]net.IPNet, len(vv))
	for i := 0; i < len(vv); i++ {
		out[i] = netCIDR(vv[i])
	}
	return out
}

func netCIDRSliceptr(vv ...string) *[]net.IPNet {
	out := netCIDRSlice(vv...)
	return &out
}

func netIP(v string) net.IP {
	return net.ParseIP(v)
}

func netIPptr(v string) *net.IP {
	ip := netIP(v)
	return &ip
}

func netIPSlice(vv ...string) []net.IP {
	out := make([]net.IP, len(vv))
	for i := 0; i < len(vv); i++ {
		out[i] = netIP(vv[i])
	}
	return out
}

func netIPSliceptr(vv ...string) *[]net.IP {
	out := netIPSlice(vv...)
	return &out
}

func netMAC(v string) net.HardwareAddr {
	mac, err := net.ParseMAC(v)
	if err != nil {
		panic(err)
	}
	return mac
}

func netMACSlice(vv ...string) []net.HardwareAddr {
	out := make([]net.HardwareAddr, len(vv))
	for i := 0; i < len(vv); i++ {
		out[i] = netMAC(vv[i])
	}
	return out
}
