package gosql

import (
	"log"
	"testing"

	"github.com/frk/compare"
	"github.com/frk/gosql/internal/errors"
)

func Test_dbchecker_loadrelation(t *testing.T) {
	tests := []struct {
		relid relid
		want  *dbrelation
		err   error
	}{{
		relid: relid{name: "relation_test", qual: "public"},
		want:  &dbrelation{name: "relation_test", namespace: "public", relkind: "r"},
		err:   nil,
	}, {
		relid: relid{name: "column_tests_1", qual: "public"},
		want:  &dbrelation{name: "column_tests_1", namespace: "public", relkind: "r"},
		err:   nil,
	}, {
		relid: relid{name: "view_test"},
		want:  &dbrelation{name: "view_test", namespace: "public", relkind: "v"},
		err:   nil,
	}, {
		relid: relid{name: "no_relation", qual: "public"},
		err:   errors.NoDBRelationError,
	}, {
		relid: relid{name: "view_test", qual: "no_namespace"},
		err:   errors.NoDBRelationError,
	}}

	for i, tt := range tests {
		dbc := new(dbchecker)
		dbc.db = testdb.db
		dbc.cmd = &command{rel: &relfield{relid: tt.relid}}

		err := dbc.load()
		rel := dbc.rel
		if err == nil {
			if rel.oid == 0 {
				t.Error(i, "expected rel.oid to be not 0")
			}

			// we don't care about these in this test
			rel.columns = nil
			rel.constraints = nil
			rel.indexes = nil

			// non-deterministic value, all we care about is that
			// it's not 0, after checking that we can move on.
			tt.want.oid = rel.oid

			if e := compare.Compare(rel, tt.want); e != nil {
				t.Error(i, e)
			}
		}

		if e := compare.Compare(err, tt.err); e != nil {
			log.Printf("%#v\n", err)
			t.Error(i, e)
		}
	}
}

func Test_dbchecker_loadcolumns(t *testing.T) {
	tests := []struct {
		relid relid
		want  []*dbcolumn
		err   error
	}{{
		relid: relid{name: "relation_test", qual: "public"},
		want: []*dbcolumn{{
			num: 1, name: "col_stub", typmod: 5,
			typ: dbtype{
				name:     pgtyp_bpchar,
				size:     -1,
				typ:      pgtyptype_base,
				category: pgtypcategory_string,
			},
		}},
		err: nil,
	}, {
		relid: relid{name: "column_tests_1"},
		want: []*dbcolumn{{
			num:        1,
			name:       "col_a",
			typmod:     -1,
			hasnotnull: true,
			hasdefault: true,
			isprimary:  true,
			typ: dbtype{
				name:     pgtyp_int4,
				size:     4,
				typ:      pgtyptype_base,
				category: pgtypcategory_numeric,
			},
		}, {
			num:        2,
			name:       "col_b",
			typmod:     -1,
			hasnotnull: true,
			typ: dbtype{
				name:        pgtyp_text,
				size:        -1,
				typ:         pgtyptype_base,
				category:    pgtypcategory_string,
				ispreferred: true,
			},
		}, {
			num:    3,
			name:   "col_c",
			typmod: -1,
			typ: dbtype{
				name:        pgtyp_bool,
				size:        1,
				typ:         pgtyptype_base,
				category:    pgtypcategory_boolean,
				ispreferred: true,
			},
		}, {
			num:        4,
			name:       "col_d",
			typmod:     -1,
			hasdefault: true,
			typ: dbtype{
				name:        pgtyp_float8,
				size:        8,
				typ:         pgtyptype_base,
				category:    pgtypcategory_numeric,
				ispreferred: true,
			},
		}, {
			num:        5,
			name:       "col_e",
			typmod:     -1,
			hasnotnull: true,
			hasdefault: true,
			typ: dbtype{
				name:     pgtyp_timestamp,
				size:     8,
				typ:      pgtyptype_base,
				category: pgtypcategory_datetime,
			},
		}},
		err: nil,
	}}

	for i, tt := range tests {
		dbc := new(dbchecker)
		dbc.db = testdb.db
		dbc.cmd = &command{rel: &relfield{relid: tt.relid}}

		err := dbc.load()
		if err == nil {
			if e := compare.Compare(dbc.rel.columns, tt.want); e != nil {
				t.Error(i, e)
			}
		}

		if e := compare.Compare(err, tt.err); e != nil {
			log.Printf("%#v\n", err)
			t.Error(i, e)
		}
	}
}

func Test_dbchecker_loadconstraints(t *testing.T) {
	tests := []struct {
		relid relid
		want  []*dbconstraint
		err   error
	}{{
		relid: relid{name: "column_tests_1"},
		want: []*dbconstraint{{
			name: "column_tests_1_pkey",
			typ:  pgconstraint_pkey,
			key:  []int64{1},
		}, {
			name: "column_tests_1_col_b_key",
			typ:  pgconstraint_unique,
			key:  []int64{2},
		}},
		err: nil,
	}}

	for i, tt := range tests {
		dbc := new(dbchecker)
		dbc.db = testdb.db
		dbc.cmd = &command{rel: &relfield{relid: tt.relid}}

		err := dbc.load()
		if err == nil {
			if e := compare.Compare(dbc.rel.constraints, tt.want); e != nil {
				t.Error(i, e)
			}
		}

		if e := compare.Compare(err, tt.err); e != nil {
			log.Printf("%#v\n", err)
			t.Error(i, e)
		}
	}
}

func Test_dbchecker_loadindexes(t *testing.T) {
	tests := []struct {
		relid relid
		want  []*dbindex
		err   error
	}{{
		relid: relid{name: "column_tests_1"},
		want: []*dbindex{{
			name:        "column_tests_1_pkey",
			natts:       1,
			isunique:    true,
			isprimary:   true,
			isimmediate: true,
			isready:     true,
			key:         []int16{1},
			indexdef:    "CREATE UNIQUE INDEX column_tests_1_pkey ON public.column_tests_1 USING btree (col_a)",
		}, {
			name:        "column_tests_1_col_b_key",
			natts:       1,
			isunique:    true,
			isimmediate: true,
			isready:     true,
			key:         []int16{2},
			indexdef:    "CREATE UNIQUE INDEX column_tests_1_col_b_key ON public.column_tests_1 USING btree (col_b)",
		}},
		err: nil,
	}}

	for i, tt := range tests {
		dbc := new(dbchecker)
		dbc.db = testdb.db
		dbc.cmd = &command{rel: &relfield{relid: tt.relid}}

		err := dbc.load()
		if err == nil {
			if e := compare.Compare(dbc.rel.indexes, tt.want); e != nil {
				t.Error(i, e)
			}
		}

		if e := compare.Compare(err, tt.err); e != nil {
			log.Printf("%#v\n", err)
			t.Error(i, e)
		}
	}
}

func Test_dbchecker_check_textsearch(t *testing.T) {
	tests := []struct {
		cmd *command
		err error
	}{{
		cmd: &command{
			rel:        &relfield{relid: relid{name: "column_tests_2"}},
			textsearch: &colid{qual: "", name: "col_text_search_ok"},
		},
		err: nil,
	}, {
		cmd: &command{
			rel:        &relfield{relid: relid{name: "column_tests_2", alias: "c"}},
			textsearch: &colid{qual: "c", name: "col_text_search_ok"},
		},
		err: nil,
	}, {
		cmd: &command{
			rel:        &relfield{relid: relid{name: "column_tests_2", alias: "c"}},
			textsearch: &colid{qual: "d", name: "col_text_search_ok"},
		},
		err: errors.NoDBRelationError,
	}, {
		cmd: &command{
			rel:        &relfield{relid: relid{name: "column_tests_2", alias: "c"}},
			textsearch: &colid{qual: "c", name: "col_none"},
		},
		err: errors.NoDBColumnError,
	}, {
		cmd: &command{
			rel:        &relfield{relid: relid{name: "column_tests_2", alias: "c"}},
			textsearch: &colid{qual: "c", name: "col_text_search_bad"},
		},
		err: errors.BadDBColumnTypeError,
	}}

	for i, tt := range tests {
		dbc := new(dbchecker)
		dbc.db = testdb.db
		dbc.cmd = tt.cmd

		if err := dbc.load(); err != nil {
			log.Printf("%#v\n", err)
			t.Error(err)
			return
		}

		err := dbc.check()
		if e := compare.Compare(err, tt.err); e != nil {
			log.Printf("%#v\n", err)
			t.Error(i, e)
		}
	}
}

func Test_dbchecker_check_orderby(t *testing.T) {
	tests := []struct {
		cmd *command
		err error
	}{{
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2"}},
			orderby: &orderbylist{items: []*orderbyitem{
				{col: colid{name: "col_orderby_a"}},
				{col: colid{name: "col_orderby_b"}},
			}},
		},
		err: nil,
	}, {
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2", alias: "c"}},
			orderby: &orderbylist{items: []*orderbyitem{
				{col: colid{qual: "c", name: "col_orderby_a"}},
				{col: colid{qual: "c", name: "col_orderby_b"}},
			}},
		},
		err: nil,
	}, {
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2", alias: "c"}},
			orderby: &orderbylist{items: []*orderbyitem{
				{col: colid{qual: "d", name: "col_orderby_a"}},
				{col: colid{qual: "d", name: "col_orderby_b"}},
			}},
		},
		err: errors.NoDBRelationError,
	}, {
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2", alias: "c"}},
			orderby: &orderbylist{items: []*orderbyitem{
				{col: colid{qual: "c", name: "col_none"}},
			}},
		},
		err: errors.NoDBColumnError,
	}}

	for i, tt := range tests {
		dbc := new(dbchecker)
		dbc.db = testdb.db
		dbc.cmd = tt.cmd

		if err := dbc.load(); err != nil {
			log.Printf("%#v\n", err)
			t.Error(err)
			return
		}

		err := dbc.check()
		if e := compare.Compare(err, tt.err); e != nil {
			log.Printf("%#v\n", err)
			t.Error(i, e)
		}
	}
}

func Test_dbchecker_check_defaults(t *testing.T) {
	tests := []struct {
		cmd *command
		err error
	}{{
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2"}},
			defaults: &collist{items: []colid{
				{name: "col_foo"},
				{name: "col_bar"},
				{name: "col_baz"},
			}},
		},
		err: nil,
	}, {
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2", alias: "c"}},
			defaults: &collist{items: []colid{
				{qual: "c", name: "col_foo"},
				{qual: "c", name: "col_baz"},
			}},
		},
		err: nil,
	}, {
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2", alias: "c"}},
			defaults: &collist{items: []colid{
				{qual: "c", name: "col_foo"},
				{qual: "d", name: "col_bar"},
			}},
		},
		err: errors.NoDBRelationError,
	}, {
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2", alias: "c"}},
			defaults: &collist{items: []colid{
				{qual: "c", name: "col_foo"},
				{qual: "c", name: "col_none"},
			}},
		},
		err: errors.NoDBColumnError,
	}}

	for i, tt := range tests {
		dbc := new(dbchecker)
		dbc.db = testdb.db
		dbc.cmd = tt.cmd

		if err := dbc.load(); err != nil {
			log.Printf("%#v\n", err)
			t.Error(err)
			return
		}

		err := dbc.check()
		if e := compare.Compare(err, tt.err); e != nil {
			log.Printf("%#v\n", err)
			t.Error(i, e)
		}
	}
}

func Test_dbchecker_check_force(t *testing.T) {
	tests := []struct {
		cmd *command
		err error
	}{{
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2"}},
			force: &collist{items: []colid{
				{name: "col_foo"},
				{name: "col_bar"},
				{name: "col_baz"},
			}},
		},
		err: nil,
	}, {
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2", alias: "c"}},
			force: &collist{items: []colid{
				{qual: "c", name: "col_foo"},
				{qual: "c", name: "col_baz"},
			}},
		},
		err: nil,
	}, {
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2", alias: "c"}},
			force: &collist{items: []colid{
				{qual: "c", name: "col_foo"},
				{qual: "d", name: "col_bar"},
			}},
		},
		err: errors.NoDBRelationError,
	}, {
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2", alias: "c"}},
			force: &collist{items: []colid{
				{qual: "c", name: "col_foo"},
				{qual: "c", name: "col_none"},
			}},
		},
		err: errors.NoDBColumnError,
	}}

	for i, tt := range tests {
		dbc := new(dbchecker)
		dbc.db = testdb.db
		dbc.cmd = tt.cmd

		if err := dbc.load(); err != nil {
			log.Printf("%#v\n", err)
			t.Error(err)
			return
		}

		err := dbc.check()
		if e := compare.Compare(err, tt.err); e != nil {
			log.Printf("%#v\n", err)
			t.Error(i, e)
		}
	}
}

func Test_dbchecker_check_returning(t *testing.T) {
	tests := []struct {
		cmd *command
		err error
	}{{
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2"}},
			returning: &collist{items: []colid{
				{name: "col_foo"},
				{name: "col_bar"},
				{name: "col_baz"},
			}},
		},
		err: nil,
	}, {
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2", alias: "c"}},
			returning: &collist{items: []colid{
				{qual: "c", name: "col_foo"},
				{qual: "c", name: "col_baz"},
			}},
		},
		err: nil,
	}, {
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2", alias: "c"}},
			returning: &collist{items: []colid{
				{qual: "c", name: "col_foo"},
				{qual: "d", name: "col_bar"},
			}},
		},
		err: errors.NoDBRelationError,
	}, {
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2", alias: "c"}},
			returning: &collist{items: []colid{
				{qual: "c", name: "col_foo"},
				{qual: "c", name: "col_none"},
			}},
		},
		err: errors.NoDBColumnError,
	}}

	for i, tt := range tests {
		dbc := new(dbchecker)
		dbc.db = testdb.db
		dbc.cmd = tt.cmd

		if err := dbc.load(); err != nil {
			log.Printf("%#v\n", err)
			t.Error(err)
			return
		}

		err := dbc.check()
		if e := compare.Compare(err, tt.err); e != nil {
			log.Printf("%#v\n", err)
			t.Error(i, e)
		}
	}
}

func Test_dbchecker_check_onconflict(t *testing.T) {
	tests := []struct {
		cmd *command
		err error
	}{{
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2"}},
			onconflict: &onconflictblock{
				column: []colid{
					{name: "col_indkey1"},
					{name: "col_indkey2"},
				},
			},
		},
		err: nil,
	}, {
		cmd: &command{
			rel:        &relfield{relid: relid{name: "column_tests_2"}},
			onconflict: &onconflictblock{column: []colid{{name: "col_none"}}},
		},
		err: errors.NoDBColumnError,
	}, {
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2"}},
			onconflict: &onconflictblock{
				column: []colid{{name: "col_indkey2"}},
			},
		},
		err: errors.NoDBIndexForColumnListError,
	}, {
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2"}},
			onconflict: &onconflictblock{
				column: []colid{{name: "col_indkey1"}},
			},
		},
		err: errors.NoDBIndexForColumnListError,
	}, {
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2"}},
			onconflict: &onconflictblock{
				column: []colid{
					{name: "col_indkey1"},
					{name: "col_indkey2"},
					{name: "col_foo"},
				},
			},
		},
		err: errors.NoDBIndexForColumnListError,
	}, {
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2"}},
			onconflict: &onconflictblock{
				index: "column_tests_2_unique_index",
			},
		},
		err: nil,
	}, {
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2"}},
			onconflict: &onconflictblock{
				index: "column_tests_2_index_none",
			},
		},
		err: errors.NoDBIndexError,
	}, {
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2"}},
			onconflict: &onconflictblock{
				index: "column_tests_2_nonunique_index",
			},
		},
		err: errors.NoDBIndexError,
	}, {
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2"}},
			onconflict: &onconflictblock{
				constraint: "column_tests_2_unique_constraint",
			},
		},
		err: nil,
	}, {
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2"}},
			onconflict: &onconflictblock{
				constraint: "column_tests_2_unique_constraint_none",
			},
		},
		err: errors.NoDBConstraintError,
	}, {
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2"}},
			onconflict: &onconflictblock{
				constraint: "column_tests_2_nonunique_constraint",
			},
		},
		err: errors.NoDBConstraintError,
	}, {
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2"}},
			onconflict: &onconflictblock{update: &collist{items: []colid{
				{name: "col_foo"},
				{name: "col_bar"},
				{name: "col_baz"},
			}}},
		},
		err: nil,
	}, {
		cmd: &command{
			rel: &relfield{relid: relid{name: "column_tests_2"}},
			onconflict: &onconflictblock{update: &collist{items: []colid{
				{name: "col_foo"},
				{name: "col_bar"},
				{name: "col_none"},
			}}},
		},
		err: errors.NoDBColumnError,
	}}

	for i, tt := range tests {
		dbc := new(dbchecker)
		dbc.db = testdb.db
		dbc.cmd = tt.cmd

		if err := dbc.load(); err != nil {
			log.Printf("%#v\n", err)
			t.Error(err)
			return
		}

		err := dbc.check()
		if e := compare.Compare(err, tt.err); e != nil {
			log.Printf("%#v\n", err)
			t.Error(i, e)
		}
	}
}
