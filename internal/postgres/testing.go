package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type TestDB struct {
	root *sql.DB
	DB   *DB
}

func (t *TestDB) Init() (err error) {
	// open the default db
	if t.root, err = sql.Open("postgres", "postgres:///?sslmode=disable"); err != nil {
		return err
	} else if err = t.root.Ping(); err != nil {
		return err
	}

	// create a new database so that the default one isn't polluted with the test data.
	dbname := "gosql_test_db"
	if _, err = t.root.Exec("DROP DATABASE IF EXISTS " + dbname); err != nil {
		return err
	}
	if _, err = t.root.Exec("CREATE DATABASE " + dbname); err != nil {
		return err
	}

	url := "postgres:///" + dbname + "?sslmode=disable"
	if t.DB, err = Open(url); err != nil {
		t.root.Close()
		return err
	}

	// populate test db
	const populatedbquery = `
CREATE TABLE relation_test (
	col_stub char
);

CREATE TABLE column_tests_1 (
	col_a serial primary key
	, col_b text not null unique
	, col_c boolean
	, col_d float8 default 0.0
	, col_e timestamp not null default now()
);

CREATE TABLE column_tests_2 (
	col_text_search_ok tsvector
	, col_text_search_bad text -- text search column must be tsvector
	, col_orderby_a text
	, col_orderby_b integer
	, col_foo integer default 0
	, col_bar text default ''
	, col_baz boolean not null default false
	, col_indkey1 text
	, col_indkey2 integer
	, col_conkey1 text
	, col_conkey2 integer
);

CREATE UNIQUE INDEX column_tests_2_unique_index ON column_tests_2 (col_indkey2, col_indkey1);
CREATE INDEX column_tests_2_nonunique_index ON column_tests_2 (col_indkey2, col_indkey1);
ALTER TABLE column_tests_2 ADD CONSTRAINT column_tests_2_unique_constraint UNIQUE (col_conkey1, col_conkey2);
ALTER TABLE column_tests_2 ADD CONSTRAINT column_tests_2_nonunique_constraint FOREIGN KEY (col_conkey1) REFERENCES column_tests_1 (col_b);

CREATE TABLE column_type_tests (
	col_bool bool
	, col_boola boolean[]
	, col_float4 float4
	, col_float4a float4[]
	, col_float8 float8
	, col_float8a float8[]
	, col_int2 int2
	, col_int2a int2[]
	, col_int4 int4
	, col_int4a int4[]
	, col_int8 int8
	, col_int8a int8[]
	, col_text text
	, col_texta text[]
);

CREATE TABLE test_user (
	id serial primary key
	, email text not null
	, full_name text not null
	, password bytea not null default ''
	, is_active boolean not null default true
	, metadata1 json not null default '{}'
	, metadata2 jsonb
	, _search_document tsvector
	, created_at timestamptz not null
	, updated_at timestamptz not null default now()
);

CREATE TABLE test_user_with_defaults (
	id serial primary key
	, email text not null default 'joe@example.com'
	, full_name text not null default 'john doe'
	, is_active boolean not null default true
	, created_at timestamptz not null default now()
	, updated_at timestamptz not null default now()
);

CREATE TABLE test_post (
	id serial primary key
	, user_id integer not null REFERENCES test_user (id)
	, content text not null
	, is_spam boolean not null
	, created_at timestamptz not null
);

CREATE TABLE test_nested (
	foo_bar_baz_val text not null default ''
	, foo_baz_val text not null default ''
	, foo2_bar_baz_val text not null default ''
	, foo2_baz_val text not null default ''
);

CREATE TABLE test_join1 (
	id serial primary key
	, post_id integer not null REFERENCES test_post (id)
);

CREATE TABLE test_join2 (
	id serial primary key
	, join1_id integer not null REFERENCES test_join1 (id)
);

CREATE TABLE test_join3 (
	id serial primary key
	, join2_id integer not null REFERENCES test_join2 (id)
);

CREATE TABLE test_join4 (
	id serial primary key
	, join3_id integer not null REFERENCES test_join3 (id)
);

CREATE VIEW view_test AS SELECT
	col_a
	, col_b
	, col_c
	, col_d
	, col_e
	, (length(col_b) > 0) AS col_z
FROM column_tests_1;

CREATE TABLE test_onconflict (
	id serial primary key
	, key int4
	, name text
	, fruit text
	, value float8
);

CREATE UNIQUE INDEX test_onconflict_key_idx ON test_onconflict (key);
CREATE UNIQUE INDEX test_onconflict_key_name_idx ON test_onconflict (key, name);
CREATE UNIQUE INDEX test_onconflict_name_fruit_idx ON test_onconflict (lower(name), upper(fruit) collate "C" text_pattern_ops);
CREATE UNIQUE INDEX test_onconflict_fruit_key_name_idx ON test_onconflict (lower(fruit), key, upper(name)) where key < 5;

ALTER TABLE test_onconflict ADD CONSTRAINT test_onconflict_key_value_key UNIQUE (key, value);

CREATE TABLE test_composite_pkey (
	id serial
	, key int4
	, name text
	, fruit text
	, value float8
	, PRIMARY KEY (id, key, name)
);

CREATE FUNCTION increment(i integer) RETURNS integer AS $$
BEGIN
	RETURN i + 1;
END;
$$ LANGUAGE plpgsql;

CREATE EXTENSION hstore;
CREATE TABLE pgsql_test (
	id serial PRIMARY KEY

	-- scalar types
	, col_bit         bit
	, col_bool        boolean
	, col_bpchar      bpchar(1)
	, col_bpchar3     bpchar(3)
	, col_char        char(1)
	, col_char3       char(3)
	, col_cidr        cidr
	, col_date        date
	, col_inet        inet
	, col_interval    interval
	, col_macaddr     macaddr
	, col_macaddr8    macaddr8
	, col_money       money
	, col_numeric     numeric
	, col_text        text
	, col_time        time
	, col_timestamp   timestamp
	, col_timestamptz timestamptz
	, col_timetz      timetz
	, col_tsquery     tsquery
	, col_tsvector    tsvector
	, col_uuid        uuid
	, col_varbit      varbit
	, col_varbit1     varbit(1)
	, col_varchar     varchar
	, col_varchar1    varchar(1)

	-- aggregate types
	, col_bitarr         bit[]          DEFAULT '{0,1,0,1}'::bit[]
	, col_boolarr        boolean[]      DEFAULT '{t,f}'::boolean[]
	, col_boxarr         box[]          DEFAULT '{(0,1),(0.5,1.5);(0,1),(0.5,1.5)}'::box[]
	, col_bpchararr      bpchar(1)[]    DEFAULT ARRAY['a','b']::bpchar[]
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
	--, col_hstorearr      hstore[]       DEFAULT ARRAY['a=>1,b=>2', 'c=>3,d=>4']::hstore[]
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
	--, col_hstore  hstore  DEFAULT 'a=>1,b=>2'
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

	if _, err = t.DB.Exec(populatedbquery); err != nil {
		return err
	}
	return nil
}

func (t *TestDB) Close() {
	if t.DB.DB != nil {
		if err := t.DB.Close(); err != nil {
			log.Println("error closing test db handle:", err)
		}
	}
	if t.root != nil {
		if _, err := t.root.Exec("DROP DATABASE " + t.DB.name); err != nil {
			log.Println("error dropping test db:", err)
		}
		if err := t.root.Close(); err != nil {
			log.Println("error closing root db handle:", err)
		}
	}
}
