package main

import (
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
		name: "SelectPostgresTestBAD_JoinBadQuantifierColumnType",
		err:  errors.BadExpressionTypeForQuantifierError,
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
		name: "SelectPostgresTestBAD_WhereBadFieldTypeForQuantifier",
		err:  errors.IllegalFieldTypeForQuantifierError,
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
		name: "SelectPostgresTestBAD_WhereColumnBadTypeForQuantifier",
		err:  errors.BadExpressionTypeForQuantifierError,
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
		err: typeError{
			errorCode:    errNoRelationColumn,
			pkgPath:      "path/to/test",
			targetName:   "SelectPostgresTestBAD_OrderByColumnNotFound",
			fieldType:    "github.com/frk/gosql.OrderBy",
			fieldName:    "_",
			tagValue:     "",
			dbName:       "gosql_test_db",
			relQualifier: "public",
			relName:      "column_tests_1",
			colQualifier: "c",
			colName:      "col_xyz",
			fileName:     "../../testdata/postgres_bad.go",
			fileLine:     346,
		},
	}, {
		name: "SelectPostgresTestBAD_OrderByRelationNotFound",
		err: typeError{
			errorCode:    errNoDatabaseRelation,
			pkgPath:      "path/to/test",
			targetName:   "SelectPostgresTestBAD_OrderByRelationNotFound",
			fieldType:    "github.com/frk/gosql.OrderBy",
			fieldName:    "_",
			tagValue:     "",
			dbName:       "gosql_test_db",
			relQualifier: "",
			relName:      "",
			colQualifier: "x",
			colName:      "col_a",
			fileName:     "../../testdata/postgres_bad.go",
			fileLine:     352,
		},
	}, {
		name: "InsertPostgresTestBAD_DefaultBadRelationAlias",
		err: typeError{
			errorCode:    errBadColumnQualifier,
			pkgPath:      "path/to/test",
			targetName:   "InsertPostgresTestBAD_DefaultBadRelationAlias",
			fieldType:    "github.com/frk/gosql.Default",
			fieldName:    "_",
			tagValue:     "",
			dbName:       "gosql_test_db",
			relQualifier: "",
			relName:      "",
			colQualifier: "x",
			colName:      "col_b",
			fileName:     "../../testdata/postgres_bad.go",
			fileLine:     358,
		},
	}, {
		name: "InsertPostgresTestBAD_DefaultColumnNotFound",
		err: typeError{
			errorCode:    errNoRelationColumn,
			pkgPath:      "path/to/test",
			targetName:   "InsertPostgresTestBAD_DefaultColumnNotFound",
			fieldType:    "github.com/frk/gosql.Default",
			fieldName:    "_",
			tagValue:     "",
			dbName:       "gosql_test_db",
			relQualifier: "public",
			relName:      "column_tests_1",
			colQualifier: "c",
			colName:      "col_xyz",
			fileName:     "../../testdata/postgres_bad.go",
			fileLine:     364,
		},
	}, {
		name: "InsertPostgresTestBAD_DefaultNotSet",
		err: typeError{
			errorCode:    errNoColumnDefault,
			pkgPath:      "path/to/test",
			targetName:   "InsertPostgresTestBAD_DefaultNotSet",
			fieldType:    "github.com/frk/gosql.Default",
			fieldName:    "_",
			tagValue:     "",
			dbName:       "gosql_test_db",
			relQualifier: "public",
			relName:      "column_tests_1",
			colQualifier: "c",
			colName:      "col_b",
			colType:      "text",
			fileName:     "../../testdata/postgres_bad.go",
			fileLine:     370,
		},
	}, {
		name: "InsertPostgresTestBAD_ForceColumnNotFound",
		err: typeError{
			errorCode:    errNoRelationColumn,
			pkgPath:      "path/to/test",
			targetName:   "InsertPostgresTestBAD_ForceColumnNotFound",
			fieldType:    "github.com/frk/gosql.Force",
			fieldName:    "_",
			tagValue:     "",
			dbName:       "gosql_test_db",
			relQualifier: "public",
			relName:      "column_tests_1",
			colQualifier: "c",
			colName:      "col_xyz",
			fileName:     "../../testdata/postgres_bad.go",
			fileLine:     376,
		},
	}, {
		name: "InsertPostgresTestBAD_ForceRelationNotFound",
		err: typeError{
			errorCode:    errNoDatabaseRelation,
			pkgPath:      "path/to/test",
			targetName:   "InsertPostgresTestBAD_ForceRelationNotFound",
			fieldType:    "github.com/frk/gosql.Force",
			fieldName:    "_",
			tagValue:     "",
			dbName:       "gosql_test_db",
			relQualifier: "",
			relName:      "",
			colQualifier: "x",
			colName:      "col_a",
			fileName:     "../../testdata/postgres_bad.go",
			fileLine:     382,
		},
	}, {
		name: "UpdatePostgresTestBAD_ReturnColumnNotFound",
		err: typeError{
			errorCode:    errNoRelationColumn,
			pkgPath:      "path/to/test",
			targetName:   "UpdatePostgresTestBAD_ReturnColumnNotFound",
			fieldType:    "github.com/frk/gosql.Return",
			fieldName:    "_",
			tagValue:     "",
			dbName:       "gosql_test_db",
			relQualifier: "public",
			relName:      "column_tests_1",
			colQualifier: "c",
			colName:      "col_xyz",
			fileName:     "../../testdata/postgres_bad.go",
			fileLine:     388,
		},
	}, {
		name: "UpdatePostgresTestBAD_ReturnRelationNotFound",
		err: typeError{
			errorCode:    errNoDatabaseRelation,
			pkgPath:      "path/to/test",
			targetName:   "UpdatePostgresTestBAD_ReturnRelationNotFound",
			fieldType:    "github.com/frk/gosql.Return",
			fieldName:    "_",
			tagValue:     "",
			dbName:       "gosql_test_db",
			relQualifier: "",
			relName:      "",
			colQualifier: "x",
			colName:      "col_a",
			fileName:     "../../testdata/postgres_bad.go",
			fileLine:     394,
		},
	}, {
		name: "UpdatePostgresTestBAD_ReturnFieldNotFound",
		err: typeError{
			errorCode:    errNoColumnField,
			pkgPath:      "path/to/test",
			targetName:   "UpdatePostgresTestBAD_ReturnFieldNotFound",
			fieldType:    "github.com/frk/gosql.Return",
			fieldName:    "_",
			tagValue:     "",
			dbName:       "gosql_test_db",
			relQualifier: "public",
			relName:      "column_tests_1",
			colQualifier: "c",
			colName:      "col_d",
			fileName:     "../../testdata/postgres_bad.go",
			fileLine:     400,
		},
		//TODO }, {
		// TODO 	name: "FilterPostgresTestBAD_TextSearchColumnNotFound",
		// TODO 	err:  errors.NoDBColumnError,
		// TODO }, {
		// TODO 	name: "FilterPostgresTestBAD_TextSearchRelationNotFound",
		// TODO 	err:  errors.NoDBRelationError,
		// TODO }, {
		// TODO 	name: "FilterPostgresTestBAD_TextSearchBadColumnType",
		// TODO 	err:  errors.BadDBColumnTypeError,
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
		name: "InsertPostgresTestBAD_BadFieldToColumnType",
		err:  errors.BadFieldToColumnTypeError,
	}, {
		name: "InsertPostgresTestBAD_ResultColumnNotFound",
		err:  errors.NoDBColumnError,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ti, err := runAnalysis(tt.name, t)
			if err != nil {
				t.Fatal(err)
			}

			dbc := new(pgTypeCheck)
			dbc.fset = tdata.Fset
			dbc.pg = testdb.pg
			dbc.ti = ti
			if dbc.ti.query != nil {
				dbc.ti.dataField = dbc.ti.query.dataField
			} else if dbc.ti.filter != nil {
				dbc.ti.dataField = dbc.ti.filter.dataField
			}

			err = dbc.run()
			if e := compare.Compare(err, tt.err); e != nil {
				t.Errorf("%v - %#v %v", e, err, err)
			}
		})
	}
}
func Test_pgchecker_loadrelation(t *testing.T) {
	tests := []struct {
		relId relId
		want  *pgRelationInfo
		err   error
	}{{
		relId: relId{name: "relation_test", qual: "public"},
		want:  &pgRelationInfo{name: "relation_test", namespace: "public", relkind: "r"},
		err:   nil,
	}, {
		relId: relId{name: "column_tests_1", qual: "public"},
		want:  &pgRelationInfo{name: "column_tests_1", namespace: "public", relkind: "r"},
		err:   nil,
	}, {
		relId: relId{name: "view_test"},
		want:  &pgRelationInfo{name: "view_test", namespace: "public", relkind: "v"},
		err:   nil,
	}, {
		relId: relId{name: "no_relation", qual: "public"},
		err:   errors.NoDBRelationError,
	}, {
		relId: relId{name: "view_test", qual: "no_namespace"},
		err:   errors.NoDBRelationError,
	}}

	for i, tt := range tests {
		dbc := new(pgTypeCheck)
		dbc.fset = tdata.Fset
		dbc.pg = testdb.pg
		dbc.ti = &targetInfo{query: &queryStruct{dataField: &dataField{relId: tt.relId}}}
		dbc.ti.dataField = dbc.ti.query.dataField

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
			t.Error(i, e)
		}
	}
}

func Test_pgchecker_loadcolumns(t *testing.T) {
	tests := []struct {
		relId relId
		want  []*pgcolumn
		err   error
	}{{
		relId: relId{name: "relation_test", qual: "public"},
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
		relId: relId{name: "column_tests_1"},
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
		dbc := new(pgTypeCheck)
		dbc.fset = tdata.Fset
		dbc.pg = testdb.pg
		dbc.ti = &targetInfo{query: &queryStruct{dataField: &dataField{relId: tt.relId}}}
		dbc.ti.dataField = dbc.ti.query.dataField

		err := dbc.run()
		if err == nil {
			if e := compare.Compare(dbc.rel.columns, tt.want); e != nil {
				t.Error(i, e)
			}
		}

		if e := compare.Compare(err, tt.err); e != nil {
			t.Error(i, e)
		}
	}
}

func Test_pgchecker_loadconstraints(t *testing.T) {
	tests := []struct {
		relId relId
		want  []*pgconstraint
		err   error
	}{{
		relId: relId{name: "column_tests_1"},
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
		dbc := new(pgTypeCheck)
		dbc.fset = tdata.Fset
		dbc.pg = testdb.pg
		dbc.ti = &targetInfo{query: &queryStruct{dataField: &dataField{relId: tt.relId}}}
		dbc.ti.dataField = dbc.ti.query.dataField

		err := dbc.run()
		if err == nil {
			if e := compare.Compare(dbc.rel.constraints, tt.want); e != nil {
				t.Error(i, e)
			}
		}

		if e := compare.Compare(err, tt.err); e != nil {
			t.Error(i, e)
		}
	}
}

func Test_pgchecker_loadindexes(t *testing.T) {
	tests := []struct {
		relId relId
		want  []*pgindex
		err   error
	}{{
		relId: relId{name: "column_tests_1"},
		want: []*pgindex{{
			name:        "column_tests_1_pkey",
			natts:       1,
			isunique:    true,
			isprimary:   true,
			isimmediate: true,
			isready:     true,
			key:         []int16{1},
			indexdef:    "CREATE UNIQUE INDEX column_tests_1_pkey ON public.column_tests_1 USING btree (col_a)",
			indexpr:     "col_a",
		}, {
			name:        "column_tests_1_col_b_key",
			natts:       1,
			isunique:    true,
			isimmediate: true,
			isready:     true,
			key:         []int16{2},
			indexdef:    "CREATE UNIQUE INDEX column_tests_1_col_b_key ON public.column_tests_1 USING btree (col_b)",
			indexpr:     "col_b",
		}},
		err: nil,
	}, {
		relId: relId{name: "test_onconflict"},

		want: []*pgindex{{
			name:        "test_onconflict_pkey",
			natts:       1,
			isunique:    true,
			isprimary:   true,
			isimmediate: true,
			isready:     true,
			key:         []int16{1},
			indexdef:    "CREATE UNIQUE INDEX test_onconflict_pkey ON public.test_onconflict USING btree (id)",
			indexpr:     "id",
		}, {
			name:        "test_onconflict_key_idx",
			natts:       1,
			isunique:    true,
			isimmediate: true,
			isready:     true,
			key:         []int16{2},
			indexdef:    "CREATE UNIQUE INDEX test_onconflict_key_idx ON public.test_onconflict USING btree (key)",
			indexpr:     "key",
		}, {
			name:        "test_onconflict_key_name_idx",
			natts:       2,
			isunique:    true,
			isimmediate: true,
			isready:     true,
			key:         []int16{2, 3},
			indexdef:    "CREATE UNIQUE INDEX test_onconflict_key_name_idx ON public.test_onconflict USING btree (key, name)",
			indexpr:     "key, name",
		}, {
			name:        "test_onconflict_name_fruit_idx",
			natts:       2,
			isunique:    true,
			isimmediate: true,
			isready:     true,
			key:         []int16{0, 0},
			indexdef:    `CREATE UNIQUE INDEX test_onconflict_name_fruit_idx ON public.test_onconflict USING btree (lower(name), upper(fruit) COLLATE "C" text_pattern_ops)`,
			indexpr:     `lower(name), upper(fruit) COLLATE "C" text_pattern_ops`,
		}, {
			name:        "test_onconflict_fruit_key_name_idx",
			natts:       3,
			isunique:    true,
			isimmediate: true,
			isready:     true,
			key:         []int16{0, 2, 0},
			indexdef:    "CREATE UNIQUE INDEX test_onconflict_fruit_key_name_idx ON public.test_onconflict USING btree (lower(fruit), key, upper(name)) WHERE (key < 5)",
			indexpr:     "lower(fruit), key, upper(name)",
			indpred:     "key < 5",
		}, {
			name:        "test_onconflict_key_value_key",
			natts:       2,
			isunique:    true,
			isimmediate: true,
			isready:     true,
			key:         []int16{2, 5},
			indexdef:    "CREATE UNIQUE INDEX test_onconflict_key_value_key ON public.test_onconflict USING btree (key, value)",
			indexpr:     "key, value",
		}},
		err: nil,
	}}

	for i, tt := range tests {
		dbc := new(pgTypeCheck)
		dbc.fset = tdata.Fset
		dbc.pg = testdb.pg
		dbc.ti = &targetInfo{query: &queryStruct{dataField: &dataField{relId: tt.relId}}}
		dbc.ti.dataField = dbc.ti.query.dataField

		err := dbc.run()
		if err == nil {
			if e := compare.Compare(dbc.rel.indexes, tt.want); e != nil {
				t.Error(i, e)
			}
		}

		if e := compare.Compare(err, tt.err); e != nil {
			t.Error(i, e)
		}
	}
}
