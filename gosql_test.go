package gosql

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testdb *testdbtype

func TestMain(m *testing.M) {
	var exitcode int

	func() { // use a func wrapper so we can rely on defer
		testdb = new(testdbtype)
		defer testdb.close()

		if err := testdb.init(); err != nil {
			panic(err)
		}

		exitcode = m.Run()
	}()

	os.Exit(exitcode)
}

type testdbtype struct {
	root *sql.DB
	pg   *postgres
}

func (t *testdbtype) init() (err error) {
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

	t.pg = new(postgres)
	t.pg.url = "postgres:///" + dbname + "?sslmode=disable"
	if err := t.pg.init(); err != nil {
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
	, created_at timestamptz not null
	, is_active boolean not null default true
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

CREATE UNIQUE INDEX column_tests_2_unique_index ON column_tests_2 (col_indkey2, col_indkey1);
CREATE INDEX column_tests_2_nonunique_index ON column_tests_2 (col_indkey2, col_indkey1);

ALTER TABLE column_tests_2 ADD CONSTRAINT column_tests_2_unique_constraint
UNIQUE (col_conkey1, col_conkey2);

ALTER TABLE column_tests_2 ADD CONSTRAINT column_tests_2_nonunique_constraint
FOREIGN KEY (col_conkey1) REFERENCES column_tests_1 (col_b);

CREATE FUNCTION increment(i integer) RETURNS integer AS $$
BEGIN
	RETURN i + 1;
END;
$$ LANGUAGE plpgsql;

` //`

	if _, err = t.pg.db.Exec(populatedbquery); err != nil {
		return err
	}

	return nil
}

func (t *testdbtype) close() {
	if t.pg.db != nil {
		if err := t.pg.close(); err != nil {
			log.Println("error closing test db handle:", err)
		}
	}
	if t.root != nil {
		if _, err := t.root.Exec("DROP DATABASE " + t.pg.name); err != nil {
			log.Println("error dropping test db:", err)
		}
		if err := t.root.Close(); err != nil {
			log.Println("error closing root db handle:", err)
		}
	}
}
