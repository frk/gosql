package typetests

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var libpq *libpqtest

type libpqtest struct {
	root   *sql.DB
	db     *sql.DB
	dbname string
	dburl  string
}

func (t *libpqtest) init() (err error) {
	// open the default db
	if t.root, err = sql.Open("postgres", "postgres:///?sslmode=disable"); err != nil {
		return err
	} else if err = t.root.Ping(); err != nil {
		return err
	}

	// create a new database so that the default one isn't polluted with the test data.
	t.dbname = "libpq_typetest_db"
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
CREATE TABLE typetest_table (
	id serial PRIMARY KEY
	, col_bit bit NOT NULL DEFAULT '1'
	, col_bitarr bit[] NOT NULL DEFAULT ARRAY['1','0']::bit[]
	, col_bit3 bit(3) NOT NULL DEFAULT '101'
	, col_bit3arr bit(3)[] NOT NULL DEFAULT ARRAY['101','010']::bit(3)[]
	, col_bool boolean NOT NULL DEFAULT true
	, col_boolarr boolean[] NOT NULL DEFAULT ARRAY[true,false]
	, col_box box NOT NULL DEFAULT '((0,1), (0.5,1.5))'
	, col_boxarr box[] NOT NULL DEFAULT ARRAY['((0,1), (0.5,1.5))','((0,1), (0.5,1.5))']::box[]
	, col_bpchar bpchar NOT NULL DEFAULT 'a'
	, col_bpchararr bpchar[] NOT NULL DEFAULT ARRAY['a','b']::bpchar[]
	, col_bpchar3 bpchar(3) NOT NULL DEFAULT 'abc'
	, col_bpchar3arr bpchar(3)[] NOT NULL DEFAULT ARRAY['abc','def']::bpchar(3)[]
	, col_bytea bytea NOT NULL DEFAULT '\xDEADBEEF'
	, col_byteaarr bytea[] NOT NULL DEFAULT ARRAY['\xDEADBEEF', '\xDEADBEEF']::bytea[]
	, col_char char NOT NULL DEFAULT 'a'
	, col_chararr char[] NOT NULL DEFAULT ARRAY['a','b']::char[]
	, col_char3 char(3) NOT NULL DEFAULT 'abc'
	, col_char3arr char(3)[] NOT NULL DEFAULT ARRAY['abc','def']::char(3)[]
	, col_cidr cidr NOT NULL DEFAULT '192.168.100.128/25'
	, col_cidrarr cidr[] NOT NULL DEFAULT ARRAY['192.168.100.128/25','2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128']::cidr[]
	, col_circle circle NOT NULL DEFAULT '<(0,0), 3.5>'
	, col_circlearr circle[] NOT NULL DEFAULT ARRAY['<(0,0), 3.5>','((0.5,1), 5)']::circle[]
	, col_date date NOT NULL DEFAULT '1999-01-08'
	, col_datearr date[] NOT NULL DEFAULT ARRAY['1999-01-08', 'May 5, 2001']::date[]
	, col_daterange daterange NOT NULL DEFAULT '[1999-01-08, 2001-01-08)'
	, col_daterangearr daterange[] NOT NULL DEFAULT ARRAY['[1999-01-08, 2001-01-08)', '(1999-01-08, 2001-01-08]']::daterange[]
	, col_float4 float4 NOT NULL DEFAULT 0.1
	, col_float4arr float4[] NOT NULL DEFAULT ARRAY[0.1, 0.2]::float4[]
	, col_float8 float8 NOT NULL DEFAULT 0.1
	, col_float8arr float8[] NOT NULL DEFAULT ARRAY[0.1, 0.2]::float8[]
	, col_hstore hstore NOT NULL DEFAULT 'a=>1,b=>2'
	, col_hstorearr hstore[] NOT NULL DEFAULT ARRAY['a=>1,b=>2', 'c=>3,d=>4']::hstore[]
	, col_inet inet NOT NULL DEFAULT '192.168.0.1/24'
	, col_inetarr inet[] NOT NULL DEFAULT ARRAY['192.168.0.1/24', '128.0.0.0/16']::inet[]
	, col_int2 int2 NOT NULL DEFAULT 1
	, col_int2arr int2[] NOT NULL DEFAULT ARRAY[1,2]::int2[]
	, col_int2vector int2vector NOT NULL DEFAULT '1 2'
	, col_int2vectorarr int2vector[] NOT NULL DEFAULT ARRAY['1 2', '3 4']::int2vector[]
	, col_int4 int4 NOT NULL DEFAULT 1
	, col_int4arr int4[] NOT NULL DEFAULT ARRAY[1,2]::int4[]
	, col_int4range int4range NOT NULL DEFAULT '[0,9)'
	, col_int4rangearr int4range[] NOT NULL DEFAULT ARRAY['[0,9)', '(0,9]']::int4range[]
	, col_int8 int8 NOT NULL DEFAULT 1
	, col_int8arr int8[] NOT NULL DEFAULT ARRAY[1,2]::int8[]
	, col_int8range int8range NOT NULL DEFAULT '[0,9)'
	, col_int8rangearr int8range[] NOT NULL DEFAULT ARRAY['[0,9)', '(0,9]']::int8range[]
	, col_interval interval NOT NULL DEFAULT '1 day'
	, col_intervalarr interval[] NOT NULL DEFAULT ARRAY['1 day','5 years 4 months 34 minutes ago']::interval[]
	, col_json json NOT NULL DEFAULT '{"foo":["bar", "baz", 123]}'
	, col_jsonarr json[] NOT NULL DEFAULT ARRAY['{"foo":["bar", "baz", 123]}', '["foo", 123]']::json[]
	, col_jsonb jsonb NOT NULL DEFAULT '{"foo":["bar", "baz", 123]}'
	, col_jsonbarr jsonb[] NOT NULL DEFAULT ARRAY['{"foo":["bar", "baz", 123]}', '["foo", 123]']::jsonb[]
	, col_line line NOT NULL DEFAULT '{1,2,3}'
	, col_linearr line[] NOT NULL DEFAULT ARRAY['{1,2,3}', '{4,5,6}']::line[]
	, col_lseg lseg NOT NULL DEFAULT '[(1,2), (3,4)]'
	, col_lsegarr lseg[] NOT NULL DEFAULT ARRAY['[(1,2), (3,4)]', '[(1.5,2.5), (3.5,4.5)]']::lseg[]
	, col_macaddr macaddr NOT NULL DEFAULT '08:00:2b:01:02:03'
	, col_macaddrarr macaddr[] NOT NULL DEFAULT ARRAY['08:00:2b:01:02:03', '08002b010203']::macaddr[]
	, col_macaddr8 macaddr8 NOT NULL DEFAULT '08:00:2b:01:02:03:04:05'
	, col_macaddr8arr macaddr8[] NOT NULL DEFAULT ARRAY['08:00:2b:01:02:03:04:05', '08002b0102030405']::macaddr8[]
	, col_money money NOT NULL DEFAULT '$1.20'
	, col_moneyarr money[] NOT NULL DEFAULT ARRAY['$1.20', '$0.99']::money[]
	, col_numeric numeric NOT NULL DEFAULT 1.2
	, col_numericarr numeric[] NOT NULL DEFAULT ARRAY[1.2,3.4]::numeric[]
	, col_numrange numrange NOT NULL DEFAULT '[1.2,3.4)'
	, col_numrangearr numrange[] NOT NULL DEFAULT ARRAY['[1.2,3.4)', '(1.2,3.4]']::numrange[]
	, col_path path NOT NULL DEFAULT '[(1,1),(2,2),(3,3)]'
	, col_patharr path[] NOT NULL DEFAULT ARRAY['[(1,1),(2,2),(3,3)]', '[(1.5,1.5),(2.5,2.5),(3.5,3.5)]']::path[]
	, col_point point NOT NULL DEFAULT '(1,1)'
	, col_pointarr point[] NOT NULL DEFAULT ARRAY['(1,1)', '(2,2)']::point[]
	, col_polygon polygon NOT NULL DEFAULT '((1,1),(2,2),(3,3))'
	, col_polygonarr polygon[] NOT NULL DEFAULT ARRAY['((1,1),(2,2),(3,3))', '((1.5,1.5),(2.5,2.5),(3.5,3.5))']::polygon[]
	, col_text text NOT NULL DEFAULT 'foo'
	, col_textarr text[] NOT NULL DEFAULT ARRAY['foo', 'bar']::text[]
	, col_time time NOT NULL DEFAULT '04:05:06.789'
	, col_timearr time[] NOT NULL DEFAULT ARRAY['04:05:06.789', '040506']::time[]
	, col_timestamp timestamp NOT NULL DEFAULT '1999-01-08 04:05:06'
	, col_timestamparr timestamp[] NOT NULL DEFAULT ARRAY['1999-01-08 04:05:06', '2004-10-19 10:23:54']::timestamp[]
	, col_timestamptz timestamptz NOT NULL DEFAULT 'January 8 04:05:06 1999 PST'
	, col_timestamptzarr timestamptz[] NOT NULL DEFAULT ARRAY['January 8 04:05:06 1999 PST','2004-10-19 10:23:54+02']::timestamptz[]
	, col_timetz timetz NOT NULL DEFAULT '04:05:06.789-8'
	, col_timetzarr timetz[] NOT NULL DEFAULT ARRAY['04:05:06.789-8','2003-04-12 04:05:06 America/New_York']::timetz[]
	, col_tsquery tsquery NOT NULL DEFAULT 'fat & rat'
	, col_tsqueryarr tsquery[] NOT NULL DEFAULT ARRAY['fat & rat', 'fat & rat & ! cat']::tsquery[]
	, col_tsrange tsrange NOT NULL DEFAULT '[1999-01-08 04:05:06, 2004-10-19 10:23:54)'
	, col_tsrangearr tsrange[] NOT NULL DEFAULT ARRAY['[1999-01-08 04:05:06, 2004-10-19 10:23:54)', '(1999-01-08 04:05:06, 2004-10-19 10:23:54]']::tsrange[]
	, col_tstzrange tstzrange NOT NULL DEFAULT '[January 8 04:05:06 1999 PST, 2004-10-19 10:23:54+02)'
	, col_tstzrangearr tstzrange[] NOT NULL DEFAULT ARRAY['[January 8 04:05:06 1999 PST, 2004-10-19 10:23:54+02)','(January 8 04:05:06 1999 PST, 2004-10-19 10:23:54+02]']::tstzrange[]
	, col_tsvector tsvector NOT NULL DEFAULT 'a fat cat sat on a mat and ate a fat rat'
	, col_tsvectorarr tsvector[] NOT NULL DEFAULT ARRAY['a fat cat sat on a mat','and ate a fat rat']::tsvector[]
	, col_uuid uuid NOT NULL DEFAULT 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'
	, col_uuidarr uuid[] NOT NULL DEFAULT ARRAY['a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11','a0eebc999c0b4ef8bb6d6bb9bd380a11']::uuid[]
	, col_varbit varbit NOT NULL DEFAULT '101'
	, col_varbitarr varbit[] NOT NULL DEFAULT ARRAY['101', '00']::varbit[]
	, col_varbit1 varbit(1) NOT NULL DEFAULT '1'
	, col_varbit1arr varbit(1)[] NOT NULL DEFAULT ARRAY['1', '0']::varbit(1)[]
	, col_varchar varchar NOT NULL DEFAULT 'foo'
	, col_varchararr varchar[] NOT NULL DEFAULT ARRAY['foo', 'bar']::varchar[]
	, col_varchar1 varchar(1) NOT NULL DEFAULT 'a'
	, col_varchar1arr varchar(1)[] NOT NULL DEFAULT ARRAY['a', 'b']::varchar(1)[]
	, col_xml xml NOT NULL DEFAULT '<foo>bar</foo>'
	, col_xmlarr xml[] NOT NULL DEFAULT ARRAY['<foo>bar</foo>','<bar>foo</bar>']::xml[]
);

INSERT INTO typetest_table DEFAULT VALUES;` //`

	if _, err = t.db.Exec(populatedbquery); err != nil {
		return err
	}

	return nil
}

func (t *libpqtest) close() {
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