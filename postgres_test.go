package gosql

import (
	"log"
	"testing"

	"github.com/frk/compare"
	"github.com/frk/gosql/internal/errors"
)

func Test_pgchecker_run(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{{
		name: "SelectPostgresTestOK_Simple",
		err:  nil,
	}, {
		name: "SelectPostgresTestBAD_NoRelation",
		err:  errors.NoDBRelationError,
	}, {
		name: "DeletePostgresTestBAD_JoinNoRelation",
		err:  errors.NoDBRelationError,
	}, {
		name: "DeletePostgresTestBAD_JoinNoRelation2",
		err:  errors.NoDBRelationError,
	}, {
		name: "SelectPostgresTestBAD_JoinNoAliasRelation",
		err:  errors.NoDBRelationError,
	}, {
		name: "SelectPostgresTestBAD_JoinNoAliasRelation2",
		err:  errors.NoDBRelationError,
	}, {
		name: "SelectPostgresTestBAD_JoinNoColumn",
		err:  errors.NoDBColumnError,
	}, {
		name: "SelectPostgresTestBAD_JoinNoColumn2",
		err:  errors.NoDBColumnError,
	}, {
		name: "SelectPostgresTestBAD_JoinBadUnaryBoolColumn",
		err:  errors.BadColumnTypeForUnaryOpError,
	}, {
		name: "SelectPostgresTestBAD_JoinBadUnaryNullColumn",
		err:  errors.BadColumnNULLSettingForNULLOpError,
	}, {
		name: "SelectPostgresTestBAD_JoinBadLiteralExpression",
		err:  errors.BadLiteralExpressionError,
	}, {
		name: "SelectPostgresTestBAD_JoinBadScalarOpColumnType",
		err:  errors.BadExpressionTypeForScalarrOpError,
	}, {
		name: "SelectPostgresTestBAD_JoinBadComparisonOperandType",
		err:  errors.BadColumnToLiteralComparisonError,
	}, {
		name: "InsertPostgresTestBAD_OnConflictNoColumn",
		err:  errors.NoDBColumnError,
	}, {
		name: "InsertPostgresTestBAD_OnConflictColumnNoIndexMatch",
		err:  errors.NoDBIndexForColumnListError,
	}, {
		name: "InsertPostgresTestBAD_OnConflictNoIndex",
		err:  errors.NoDBIndexError,
	}, {
		name: "InsertPostgresTestBAD_OnConflictNoUniqueIndex",
		err:  errors.NoDBIndexError,
	}, {
		name: "InsertPostgresTestBAD_OnConflictNoConstraint",
		err:  errors.NoDBConstraintError,
	}, {
		name: "InsertPostgresTestBAD_OnConflictNoUniqueConstraint",
		err:  errors.NoDBConstraintError,
	}, {
		name: "InsertPostgresTestBAD_OnConflictUpdateColumnNotFound",
		err:  errors.NoDBColumnError,
	}, {
		name: "SelectPostgresTestBAD_WhereFieldNotFound",
		err:  errors.NoDBColumnError,
	}, {
		name: "SelectPostgresTestBAD_WhereAliasNotFound",
		err:  errors.NoDBRelationError,
	}, {
		name: "SelectPostgresTestBAD_WherePointerFieldForNonNullColumn",
		err:  errors.IllegalPtrFieldForNotNullColumnError,
	}, {
		name: "SelectPostgresTestBAD_WhereBadFieldTypeForScalarOp",
		err:  errors.IllegalFieldTypeForScalarOpError,
	}, {
		name: "SelectPostgresTestBAD_WhereCannotCompareTypes",
		err:  errors.BadFieldToColumnTypeError,
	}, {
		name: "SelectPostgresTestBAD_WhereColumnTypeForFuncname",
		err:  errors.BadColumnTypeForDBFuncError,
	}, {
		name: "SelectPostgresTestBAD_WhereColumnNotFound",
		err:  errors.NoDBColumnError,
	}, {
		name: "SelectPostgresTestBAD_WhereColumnNotFoundBadAlias",
		err:  errors.NoDBRelationError,
	}, {
		name: "SelectPostgresTestBAD_WhereColumnBadBoolOp",
		err:  errors.BadColumnTypeForUnaryOpError,
	}, {
		name: "SelectPostgresTestBAD_WhereColumnNotFoundRHS",
		err:  errors.NoDBColumnError,
	}, {
		name: "SelectPostgresTestBAD_WhereColumnNotFoundRHSBadAlias",
		err:  errors.NoDBRelationError,
	}, {
		name: "SelectPostgresTestBAD_WhereColumnBadLiteralExpression",
		err:  errors.BadLiteralExpressionError,
	}, {
		name: "SelectPostgresTestBAD_WhereColumnBadTypeForScalarOp",
		err:  errors.BadExpressionTypeForScalarrOpError,
	}, {
		name: "SelectPostgresTestBAD_WhereColumnBadTypeComparison",
		err:  errors.BadColumnToLiteralComparisonError,
	}, {
		name: "SelectPostgresTestBAD_WhereBetweenColumnNotFound",
		err:  errors.NoDBColumnError,
	}, {
		name: "SelectPostgresTestBAD_WhereBetweenRelationNotFound",
		err:  errors.NoDBRelationError,
	}, {
		name: "SelectPostgresTestBAD_WhereBetweenArgColumnNotFound",
		err:  errors.NoDBColumnError,
	}, {
		name: "SelectPostgresTestBAD_WhereBetweenArgRelationNotFound",
		err:  errors.NoDBRelationError,
	}, {
		name: "SelectPostgresTestBAD_WhereBetweenComparisonBadArgType",
		err:  errors.BadColumnToColumnTypeComparisonError,
	}, {
		name: "SelectPostgresTestBAD_OrderByColumnNotFound",
		err:  errors.NoDBColumnError,
	}, {
		name: "SelectPostgresTestBAD_OrderByRelationNotFound",
		err:  errors.NoDBRelationError,
	}, {
		name: "InsertPostgresTestBAD_DefaultBadRelationAlias",
		err:  errors.BadTargetTableForDefaultError,
	}, {
		name: "InsertPostgresTestBAD_DefaultColumnNotFound",
		err:  errors.NoDBColumnError,
	}, {
		name: "InsertPostgresTestBAD_DefaultNotSet",
		err:  errors.NoColumnDefaultSetError,
	}, {
		name: "InsertPostgresTestBAD_ForceColumnNotFound",
		err:  errors.NoDBColumnError,
	}, {
		name: "InsertPostgresTestBAD_ForceRelationNotFound",
		err:  errors.NoDBRelationError,
	}, {
		name: "UpdatePostgresTestBAD_ReturnColumnNotFound",
		err:  errors.NoDBColumnError,
	}, {
		name: "UpdatePostgresTestBAD_ReturnRelationNotFound",
		err:  errors.NoDBRelationError,
	}, {
		name: "FilterPostgresTestBAD_TextSearchColumnNotFound",
		err:  errors.NoDBColumnError,
	}, {
		name: "FilterPostgresTestBAD_TextSearchRelationNotFound",
		err:  errors.NoDBRelationError,
	}, {
		name: "FilterPostgresTestBAD_TextSearchBadColumnType",
		err:  errors.BadDBColumnTypeError,
	}, {
		name: "SelectPostgresTestBAD_RelationColumnNotFound",
		err:  errors.NoDBColumnError,
	}, {
		name: "SelectPostgresTestBAD_RelationColumnAliasNotFound",
		err:  errors.NoDBRelationError,
	}, {
		name: "InsertPostgresTestBAD_RelationColumnNotFound",
		err:  errors.NoDBColumnError,
	}, {
		name: "InsertPostgresTestBAD_BadJSONOption",
		err:  errors.BadUseJSONTargetColumnError,
	}, {
		name: "InsertPostgresTestBAD_BadXMLOption",
		err:  errors.BadUseXMLTargetColumnError,
	}, {
		name: "InsertPostgresTestBAD_BadFieldToColumnType",
		err:  errors.FieldToColumnTypeError,
	}, {
		name: "InsertPostgresTestBAD_ResultColumnNotFound",
		err:  errors.NoDBColumnError,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := runAnalysis(tt.name, t)
			if err != nil {
				t.Fatal(err)
			}

			dbc := new(pgchecker)
			dbc.db = testdb.db
			dbc.pgcat = testdb.pgcat
			dbc.cmd = cmd

			err = dbc.run()
			if e := compare.Compare(err, tt.err); e != nil {
				t.Errorf("%v - %#v %v", e, err, err)
			}
		})
	}
}
func Test_pgchecker_loadrelation(t *testing.T) {
	tests := []struct {
		relid relid
		want  *pgrelation
		err   error
	}{{
		relid: relid{name: "relation_test", qual: "public"},
		want:  &pgrelation{name: "relation_test", namespace: "public", relkind: "r"},
		err:   nil,
	}, {
		relid: relid{name: "column_tests_1", qual: "public"},
		want:  &pgrelation{name: "column_tests_1", namespace: "public", relkind: "r"},
		err:   nil,
	}, {
		relid: relid{name: "view_test"},
		want:  &pgrelation{name: "view_test", namespace: "public", relkind: "v"},
		err:   nil,
	}, {
		relid: relid{name: "no_relation", qual: "public"},
		err:   errors.NoDBRelationError,
	}, {
		relid: relid{name: "view_test", qual: "no_namespace"},
		err:   errors.NoDBRelationError,
	}}

	for i, tt := range tests {
		dbc := new(pgchecker)
		dbc.db = testdb.db
		dbc.pgcat = testdb.pgcat
		dbc.cmd = &command{rel: &relfield{relid: tt.relid}}

		err := dbc.run()
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

func Test_pgchecker_loadcolumns(t *testing.T) {
	tests := []struct {
		relid relid
		want  []*pgcolumn
		err   error
	}{{
		relid: relid{name: "relation_test", qual: "public"},
		want: []*pgcolumn{{
			num: 1, name: "col_stub", typmod: 5,
			typoid: pgtyp_bpchar,
			typ: &pgtype{
				oid:      pgtyp_bpchar,
				name:     "bpchar",
				namefmt:  "character",
				length:   -1,
				typ:      pgtyptype_base,
				category: pgtypcategory_string,
			},
		}},
		err: nil,
	}, {
		relid: relid{name: "column_tests_1"},
		want: []*pgcolumn{{
			num:        1,
			name:       "col_a",
			typmod:     -1,
			hasnotnull: true,
			hasdefault: true,
			isprimary:  true,
			typoid:     pgtyp_int4,
			typ: &pgtype{
				oid:      pgtyp_int4,
				name:     "int4",
				namefmt:  "integer",
				length:   4,
				typ:      pgtyptype_base,
				category: pgtypcategory_numeric,
			},
		}, {
			num:        2,
			name:       "col_b",
			typmod:     -1,
			hasnotnull: true,
			typoid:     pgtyp_text,
			typ: &pgtype{
				oid:         pgtyp_text,
				name:        "text",
				namefmt:     "text",
				length:      -1,
				typ:         pgtyptype_base,
				category:    pgtypcategory_string,
				ispreferred: true,
			},
		}, {
			num:    3,
			name:   "col_c",
			typmod: -1,
			typoid: pgtyp_bool,
			typ: &pgtype{
				oid:         pgtyp_bool,
				name:        "bool",
				namefmt:     "boolean",
				length:      1,
				typ:         pgtyptype_base,
				category:    pgtypcategory_boolean,
				ispreferred: true,
			},
		}, {
			num:        4,
			name:       "col_d",
			typmod:     -1,
			hasdefault: true,
			typoid:     pgtyp_float8,
			typ: &pgtype{
				oid:         pgtyp_float8,
				name:        "float8",
				namefmt:     "double precision",
				length:      8,
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
			typoid:     pgtyp_timestamp,
			typ: &pgtype{
				oid:      pgtyp_timestamp,
				name:     "timestamp",
				namefmt:  "timestamp without time zone",
				length:   8,
				typ:      pgtyptype_base,
				category: pgtypcategory_datetime,
			},
		}},
		err: nil,
	}}

	for i, tt := range tests {
		dbc := new(pgchecker)
		dbc.db = testdb.db
		dbc.pgcat = testdb.pgcat
		dbc.cmd = &command{rel: &relfield{relid: tt.relid}}

		err := dbc.run()
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

func Test_pgchecker_loadconstraints(t *testing.T) {
	tests := []struct {
		relid relid
		want  []*pgconstraint
		err   error
	}{{
		relid: relid{name: "column_tests_1"},
		want: []*pgconstraint{{
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
		dbc := new(pgchecker)
		dbc.db = testdb.db
		dbc.pgcat = testdb.pgcat
		dbc.cmd = &command{rel: &relfield{relid: tt.relid}}

		err := dbc.run()
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

func Test_pgchecker_loadindexes(t *testing.T) {
	tests := []struct {
		relid relid
		want  []*pgindex
		err   error
	}{{
		relid: relid{name: "column_tests_1"},
		want: []*pgindex{{
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
		dbc := new(pgchecker)
		dbc.db = testdb.db
		dbc.pgcat = testdb.pgcat
		dbc.cmd = &command{rel: &relfield{relid: tt.relid}}

		err := dbc.run()
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

func Test_pgchecker_check_textsearch(t *testing.T) {
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
		dbc := new(pgchecker)
		dbc.db = testdb.db
		dbc.pgcat = testdb.pgcat
		dbc.cmd = tt.cmd

		err := dbc.run()
		if e := compare.Compare(err, tt.err); e != nil {
			log.Printf("%#v\n", err)
			t.Error(i, e)
		}
	}
}

func Test_pgchecker_check_orderby(t *testing.T) {
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
		dbc := new(pgchecker)
		dbc.db = testdb.db
		dbc.pgcat = testdb.pgcat
		dbc.cmd = tt.cmd

		err := dbc.run()
		if e := compare.Compare(err, tt.err); e != nil {
			log.Printf("%#v\n", err)
			t.Error(i, e)
		}
	}
}

func Test_pgchecker_check_defaults(t *testing.T) {
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
		err: errors.BadTargetTableForDefaultError,
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
		dbc := new(pgchecker)
		dbc.db = testdb.db
		dbc.pgcat = testdb.pgcat
		dbc.cmd = tt.cmd

		err := dbc.run()
		if e := compare.Compare(err, tt.err); e != nil {
			log.Printf("%#v\n", err)
			t.Error(i, e)
		}
	}
}

func Test_pgchecker_check_force(t *testing.T) {
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
		dbc := new(pgchecker)
		dbc.db = testdb.db
		dbc.pgcat = testdb.pgcat
		dbc.cmd = tt.cmd

		err := dbc.run()
		if e := compare.Compare(err, tt.err); e != nil {
			log.Printf("%#v\n", err)
			t.Error(i, e)
		}
	}
}

func Test_pgchecker_check_returning(t *testing.T) {
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
		dbc := new(pgchecker)
		dbc.db = testdb.db
		dbc.pgcat = testdb.pgcat
		dbc.cmd = tt.cmd

		err := dbc.run()
		if e := compare.Compare(err, tt.err); e != nil {
			log.Printf("%#v\n", err)
			t.Error(i, e)
		}
	}
}

func Test_pgchecker_check_onconflict(t *testing.T) {
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
		dbc := new(pgchecker)
		dbc.db = testdb.db
		dbc.pgcat = testdb.pgcat
		dbc.cmd = tt.cmd

		err := dbc.run()
		if e := compare.Compare(err, tt.err); e != nil {
			log.Printf("%#v\n", err)
			t.Error(i, e)
		}
	}
}
