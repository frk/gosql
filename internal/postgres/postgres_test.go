package postgres

import (
	"testing"

	"github.com/frk/compare"
	"github.com/frk/gosql/internal/analysis"
	"github.com/frk/gosql/internal/x/testutil"
)

func init() {
	compare.DefaultConfig.ObserveFieldTag = "cmp"
}

var tdata = testutil.ParseTestdata("../testdata")

func testCheck(name string, t *testing.T) (*TargetInfo, error) {
	named := testutil.FindNamedType(name, tdata)
	if named == nil {
		// Stop the test if no type with the given name was found.
		t.Fatal(name, " not found")
		return nil, nil
	}

	info := new(analysis.Info)
	ts, err := analysis.Run(tdata.Fset, named, info)
	if err != nil {
		return nil, err
	}

	return Check(testdb.DB, ts, info)
}

func TestCheck(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{{
		name: "SelectPostgresTestOK_Simple",
		err:  nil,
	}, {
		name: "SelectPostgresTestBAD_NoRelation",
		err: Error{
			Code:       ErrNoRelation,
			PkgPath:    "path/to/test",
			DBName:     "gosql_test_db",
			RelName:    "norel",
			TargetName: "SelectPostgresTestBAD_NoRelation",
			FieldType:  "path/to/test.CT1",
			FieldName:  "Columns",
			FieldTag:   `rel:"norel"`,
			FileName:   "../testdata/postgres_bad.go",
			FileLine:   13,
		},
	}, {
		name: "DeletePostgresTestBAD_JoinNoRelation",
		err: Error{
			Code:       ErrNoRelation,
			PkgPath:    "path/to/test",
			DBName:     "gosql_test_db",
			RelName:    "norel",
			TargetName: "DeletePostgresTestBAD_JoinNoRelation",
			FieldType:  "github.com/frk/gosql.Relation",
			FieldName:  "_",
			FieldTag:   `sql:"norel:b"`,
			FileName:   "../testdata/postgres_bad.go",
			FileLine:   20,
		},
	}, {
		name: "DeletePostgresTestBAD_JoinNoRelation2",
		err: Error{
			Code:       ErrNoRelation,
			PkgPath:    "path/to/test",
			DBName:     "gosql_test_db",
			RelName:    "norel",
			TargetName: "DeletePostgresTestBAD_JoinNoRelation2",
			FieldType:  "github.com/frk/gosql.LeftJoin",
			FieldName:  "_",
			FieldTag:   `sql:"norel:c,c.b_id = b.id"`,
			FileName:   "../testdata/postgres_bad.go",
			FileLine:   32,
		},
	}, {
		name: "SelectPostgresTestBAD_JoinNoAliasRelation",
		err: Error{
			Code:         ErrBadAlias,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			RelName:      "column_tests_2",
			ColName:      "col_foo",
			ColQualifier: "x",
			TargetName:   "SelectPostgresTestBAD_JoinNoAliasRelation",
			FieldType:    "github.com/frk/gosql.LeftJoin",
			FieldName:    "_",
			FieldTag:     `sql:"column_tests_2:b,x.col_foo = a.col_a"`,
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     43,
		},
	}, {
		name: "SelectPostgresTestBAD_JoinNoAliasRelation2",
		err: Error{
			Code:         ErrBadAlias,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			ColName:      "col_a",
			ColQualifier: "x",
			TargetName:   "SelectPostgresTestBAD_JoinNoAliasRelation2",
			FieldType:    "github.com/frk/gosql.LeftJoin",
			FieldName:    "_",
			FieldTag:     `sql:"column_tests_2:b,b.col_foo = x.col_a"`,
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     51,
		},
	}, {
		name: "SelectPostgresTestBAD_JoinNoColumn",
		err: Error{
			Code:         ErrNoColumn,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			RelName:      "column_tests_2",
			RelSchema:    "public",
			ColName:      "nocol",
			ColQualifier: "b",
			LitExpr:      "",
			TargetName:   "SelectPostgresTestBAD_JoinNoColumn",
			FieldType:    "github.com/frk/gosql.LeftJoin",
			FieldName:    "_",
			FieldTag:     `sql:"column_tests_2:b,b.nocol = a.nocol"`,
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     59,
		},
	}, {
		name: "SelectPostgresTestBAD_JoinNoColumn2",
		err: Error{
			Code:         ErrNoColumn,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			RelName:      "column_tests_1",
			RelSchema:    "public",
			ColName:      "nocol",
			ColQualifier: "a",
			TargetName:   "SelectPostgresTestBAD_JoinNoColumn2",
			FieldType:    "github.com/frk/gosql.LeftJoin",
			FieldName:    "_",
			FieldTag:     `sql:"column_tests_2:b,b.col_foo = a.nocol"`,
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     67,
		},
	}, {
		name: "SelectPostgresTestBAD_JoinBadUnaryBoolColumn",
		err: Error{
			Code:         ErrBadUnaryPredicateType,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			RelName:      "column_tests_2",
			RelSchema:    "public",
			ColName:      "col_foo",
			ColQualifier: "b",
			ColType:      "integer",
			TargetName:   "SelectPostgresTestBAD_JoinBadUnaryBoolColumn",
			Predicate:    analysis.IsTrue,
			FieldType:    "github.com/frk/gosql.LeftJoin",
			FieldName:    "_",
			FieldTag:     "sql:\"column_tests_2:b,b.col_foo istrue\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     75,
		},
	}, {
		name: "SelectPostgresTestBAD_JoinBadUnaryNullColumn",
		err: Error{
			Code:         ErrBadUnaryPredicateType,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			RelName:      "column_tests_2",
			RelSchema:    "public",
			ColName:      "col_baz",
			ColQualifier: "b",
			ColType:      "boolean",
			TargetName:   "SelectPostgresTestBAD_JoinBadUnaryNullColumn",
			Predicate:    analysis.IsNull,
			FieldType:    "github.com/frk/gosql.LeftJoin",
			FieldName:    "_",
			FieldTag:     "sql:\"column_tests_2:b,b.col_baz isnull\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     83,
		},
	}, {
		name: "SelectPostgresTestBAD_JoinBadLiteralExpression",
		err: Error{
			Code:       ErrBadLiteralExpression,
			PkgPath:    "path/to/test",
			DBName:     "gosql_test_db",
			LitExpr:    "'foo'bar",
			TargetName: "SelectPostgresTestBAD_JoinBadLiteralExpression",
			FieldType:  "github.com/frk/gosql.LeftJoin",
			FieldName:  "_",
			FieldTag:   "sql:\"column_tests_2:b,b.col_baz = 'foo'bar \"",
			FileName:   "../testdata/postgres_bad.go",
			FileLine:   91,
		},
	}, {
		name: "SelectPostgresTestBAD_JoinBadQuantifierColumnType",
		err: Error{
			Code:       ErrBadComparisonOperation,
			PkgPath:    "path/to/test",
			DBName:     "gosql_test_db",
			RelName:    "column_tests_2",
			RelSchema:  "public",
			ColName:    "col_foo",
			ColType:    "integer",
			Col2Name:   "col_a",
			Col2Type:   "integer",
			Predicate:  analysis.IsGT,
			Quantifier: analysis.QuantAny,
			TargetName: "SelectPostgresTestBAD_JoinBadQuantifierColumnType",
			FieldType:  "github.com/frk/gosql.LeftJoin",
			FieldName:  "_",
			FieldTag:   "sql:\"column_tests_2:b,b.col_foo >any a.col_a\"",
			FileName:   "../testdata/postgres_bad.go",
			FileLine:   99,
		},
	}, {
		name: "SelectPostgresTestBAD_JoinBadComparisonOperandType",
		err: Error{
			Code:       ErrBadComparisonOperation,
			PkgPath:    "path/to/test",
			DBName:     "gosql_test_db",
			RelName:    "column_tests_2",
			RelSchema:  "public",
			ColName:    "col_baz",
			ColType:    "boolean",
			Col2Type:   "unknown",
			LitExpr:    "'baz'",
			Predicate:  analysis.IsLT,
			TargetName: "SelectPostgresTestBAD_JoinBadComparisonOperandType",
			FieldType:  "github.com/frk/gosql.LeftJoin",
			FieldName:  "_",
			FieldTag:   "sql:\"column_tests_2:b,b.col_baz < 'baz'\"",
			FileName:   "../testdata/postgres_bad.go",
			FileLine:   107,
		},
	}, {
		name: "InsertPostgresTestBAD_OnConflictNoColumn",
		err: Error{
			Code:         ErrNoColumn,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			RelName:      "column_tests_1",
			RelSchema:    "public",
			ColName:      "col_xyz",
			ColQualifier: "c",
			TargetName:   "InsertPostgresTestBAD_OnConflictNoColumn",
			FieldType:    "github.com/frk/gosql.Column",
			FieldName:    "_",
			FieldTag:     "sql:\"c.col_xyz\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     115,
		},
	}, {
		name: "InsertPostgresTestBAD_OnConflictColumnNoIndexMatch",
		err: Error{
			Code:       ErrNoUniqueIndex,
			PkgPath:    "path/to/test",
			DBName:     "gosql_test_db",
			RelName:    "column_tests_1",
			RelSchema:  "public",
			TargetName: "InsertPostgresTestBAD_OnConflictColumnNoIndexMatch",
			FieldType:  "github.com/frk/gosql.Column",
			FieldName:  "_",
			FieldTag:   "sql:\"c.col_a,c.col_b\"",
			FileName:   "../testdata/postgres_bad.go",
			FileLine:   123,
		},
	}, {
		name: "InsertPostgresTestBAD_OnConflictNoIndex",
		err: Error{
			Code:       ErrNoUniqueIndex,
			PkgPath:    "path/to/test",
			DBName:     "gosql_test_db",
			RelName:    "column_tests_1",
			RelSchema:  "public",
			TargetName: "InsertPostgresTestBAD_OnConflictNoIndex",
			FieldType:  "github.com/frk/gosql.Index",
			FieldName:  "_",
			FieldTag:   "sql:\"some_index\"",
			FileName:   "../testdata/postgres_bad.go",
			FileLine:   131,
		},
	}, {
		name: "InsertPostgresTestBAD_OnConflictNoUniqueIndex",
		err: Error{
			Code:       ErrNoUniqueIndex,
			PkgPath:    "path/to/test",
			DBName:     "gosql_test_db",
			RelName:    "column_tests_2",
			RelSchema:  "public",
			TargetName: "InsertPostgresTestBAD_OnConflictNoUniqueIndex",
			FieldType:  "github.com/frk/gosql.Index",
			FieldName:  "_",
			FieldTag:   "sql:\"column_tests_2_nonunique_index\"",
			FileName:   "../testdata/postgres_bad.go",
			FileLine:   139,
		},
	}, {
		name: "InsertPostgresTestBAD_OnConflictNoConstraint",
		err: Error{
			Code:       ErrNoUniqueConstraint,
			PkgPath:    "path/to/test",
			DBName:     "gosql_test_db",
			RelName:    "column_tests_1",
			RelSchema:  "public",
			TargetName: "InsertPostgresTestBAD_OnConflictNoConstraint",
			FieldType:  "github.com/frk/gosql.Constraint",
			FieldName:  "_",
			FieldTag:   "sql:\"some_constraint\"",
			FileName:   "../testdata/postgres_bad.go",
			FileLine:   147,
		},
	}, {
		name: "InsertPostgresTestBAD_OnConflictNoUniqueConstraint",
		err: Error{
			Code:       ErrNoUniqueConstraint,
			PkgPath:    "path/to/test",
			DBName:     "gosql_test_db",
			RelName:    "column_tests_2",
			RelSchema:  "public",
			TargetName: "InsertPostgresTestBAD_OnConflictNoUniqueConstraint",
			FieldType:  "github.com/frk/gosql.Constraint",
			FieldName:  "_",
			FieldTag:   "sql:\"column_tests_2_nonunique_constraint\"",
			FileName:   "../testdata/postgres_bad.go",
			FileLine:   155,
		},
	}, {
		name: "InsertPostgresTestBAD_OnConflictUpdateColumnNotFound",
		err: Error{
			Code:         ErrNoColumn,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			RelName:      "column_tests_2",
			RelSchema:    "public",
			ColName:      "col_a",
			ColQualifier: "c",
			TargetName:   "InsertPostgresTestBAD_OnConflictUpdateColumnNotFound",
			FieldType:    "github.com/frk/gosql.Update",
			FieldName:    "_",
			FieldTag:     "sql:\"c.col_a,c.col_b,c.col_xyz\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     164,
		},
	}, {
		name: "SelectPostgresTestBAD_WhereFieldColumnNotFound",
		err: Error{
			Code:         ErrNoColumn,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			RelName:      "column_tests_1",
			RelSchema:    "public",
			ColName:      "id",
			ColQualifier: "c",
			Predicate:    analysis.IsEQ,
			Quantifier:   0x0,
			TargetName:   "SelectPostgresTestBAD_WhereFieldColumnNotFound",
			FieldType:    "int",
			FieldName:    "Id",
			FieldTag:     "sql:\"c.id\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     172,
		},
	}, {
		name: "SelectPostgresTestBAD_WhereBadAlias",
		err: Error{
			Code:         ErrBadAlias,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			ColName:      "id",
			ColQualifier: "x",
			Predicate:    analysis.IsEQ,
			TargetName:   "SelectPostgresTestBAD_WhereBadAlias",
			FieldType:    "int",
			FieldName:    "Id",
			FieldTag:     "sql:\"x.id\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     180,
		},
	}, {
		name: "SelectPostgresTestBAD_WhereCannotCompareTypes",
		err: Error{
			Code:         ErrBadComparisonOperation,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			RelName:      "column_tests_1",
			RelSchema:    "public",
			ColName:      "col_e",
			ColQualifier: "c",
			ColType:      "timestamp without time zone",
			Predicate:    analysis.IsEQ,
			TargetName:   "SelectPostgresTestBAD_WhereCannotCompareTypes",
			FieldType:    "float64",
			FieldName:    "D",
			FieldTag:     "sql:\"c.col_e\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     188,
		},
	}, {
		name: "SelectPostgresTestBAD_WhereColumnTypeForFuncname",
		err: Error{
			Code:         ErrBadProcType,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			RelName:      "column_tests_1",
			RelSchema:    "public",
			ColName:      "col_d",
			ColQualifier: "c",
			ColType:      "double precision",
			Predicate:    analysis.IsEQ,
			TargetName:   "SelectPostgresTestBAD_WhereColumnTypeForFuncname",
			FieldType:    "float64",
			FieldName:    "D",
			FieldTag:     "sql:\"c.col_d,@lower\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     196,
		},
	}, {
		name: "SelectPostgresTestBAD_WhereColumnNotFound",
		err: Error{
			Code:         ErrNoColumn,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			RelName:      "column_tests_1",
			RelSchema:    "public",
			ColName:      "col_xyz",
			ColQualifier: "c",
			TargetName:   "SelectPostgresTestBAD_WhereColumnNotFound",
			FieldType:    "github.com/frk/gosql.Column",
			FieldName:    "_",
			FieldTag:     "sql:\"c.col_xyz istrue\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     204,
		},
	}, {
		name: "SelectPostgresTestBAD_WhereColumnNotFoundBadAlias",
		err: Error{
			Code:         ErrBadAlias,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			ColName:      "col_a",
			ColQualifier: "x",
			TargetName:   "SelectPostgresTestBAD_WhereColumnNotFoundBadAlias",
			FieldType:    "github.com/frk/gosql.Column",
			FieldName:    "_",
			FieldTag:     "sql:\"x.col_a = 123\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     212,
		},
	}, {
		name: "SelectPostgresTestBAD_WhereColumnBadBoolOp",
		err: Error{
			Code:         ErrBadUnaryPredicateType,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			RelName:      "column_tests_1",
			RelSchema:    "public",
			ColName:      "col_a",
			ColQualifier: "c",
			ColType:      "integer",
			Predicate:    analysis.IsTrue,
			TargetName:   "SelectPostgresTestBAD_WhereColumnBadBoolOp",
			FieldType:    "github.com/frk/gosql.Column",
			FieldName:    "_",
			FieldTag:     "sql:\"c.col_a istrue\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     220,
		},
	}, {
		name: "SelectPostgresTestBAD_WhereColumnNotFoundRHS",
		err: Error{
			Code:         ErrNoColumn,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			RelName:      "column_tests_1",
			RelSchema:    "public",
			ColName:      "col_xyz",
			ColQualifier: "c",
			TargetName:   "SelectPostgresTestBAD_WhereColumnNotFoundRHS",
			FieldType:    "github.com/frk/gosql.Column",
			FieldName:    "_",
			FieldTag:     "sql:\"c.col_a = c.col_xyz\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     236,
		},
	}, {
		name: "SelectPostgresTestBAD_WhereColumnNotFoundRHSBadAlias",
		err: Error{
			Code:         ErrBadAlias,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			ColName:      "col_a",
			ColQualifier: "x",
			TargetName:   "SelectPostgresTestBAD_WhereColumnNotFoundRHSBadAlias",
			FieldType:    "github.com/frk/gosql.Column",
			FieldName:    "_",
			FieldTag:     "sql:\"c.col_a = x.col_a\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     244,
		},
	}, {
		name: "SelectPostgresTestBAD_WhereColumnBadLiteralExpression",
		err: Error{
			Code:       ErrBadLiteralExpression,
			PkgPath:    "path/to/test",
			DBName:     "gosql_test_db",
			LitExpr:    "123abc",
			TargetName: "SelectPostgresTestBAD_WhereColumnBadLiteralExpression",
			FieldType:  "github.com/frk/gosql.Column",
			FieldName:  "_",
			FieldTag:   "sql:\"c.col_a = 123abc\"",
			FileName:   "../testdata/postgres_bad.go",
			FileLine:   252,
		},
	}, {
		name: "SelectPostgresTestBAD_WhereColumnBadTypeForQuantifier",
		err: Error{
			Code:       ErrBadComparisonOperation,
			PkgPath:    "path/to/test",
			DBName:     "gosql_test_db",
			RelName:    "column_tests_1",
			RelSchema:  "public",
			ColName:    "col_a",
			ColType:    "integer",
			Col2Name:   "col_b",
			Col2Type:   "text",
			Predicate:  analysis.IsIn,
			TargetName: "SelectPostgresTestBAD_WhereColumnBadTypeForQuantifier",
			FieldType:  "github.com/frk/gosql.Column",
			FieldName:  "_",
			FieldTag:   "sql:\"c.col_a isin c.col_b\"",
			FileName:   "../testdata/postgres_bad.go",
			FileLine:   260,
		},
	}, {
		name: "SelectPostgresTestBAD_WhereColumnBadTypeComparison",
		err: Error{
			Code:       ErrBadComparisonOperation,
			PkgPath:    "path/to/test",
			DBName:     "gosql_test_db",
			RelName:    "column_tests_1",
			RelSchema:  "public",
			ColName:    "col_a",
			ColType:    "integer",
			Col2Name:   "col_b",
			Col2Type:   "text",
			Predicate:  analysis.IsEQ,
			TargetName: "SelectPostgresTestBAD_WhereColumnBadTypeComparison",
			FieldType:  "github.com/frk/gosql.Column",
			FieldName:  "_",
			FieldTag:   "sql:\"c.col_a = c.col_b\"",
			FileName:   "../testdata/postgres_bad.go",
			FileLine:   268,
		},
	}, {
		name: "SelectPostgresTestBAD_WhereBetweenColumnNotFound",
		err: Error{
			Code:         ErrNoColumn,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			RelName:      "column_tests_1",
			RelSchema:    "public",
			ColName:      "col_xyz",
			ColQualifier: "c",
			TargetName:   "SelectPostgresTestBAD_WhereBetweenColumnNotFound",
			FieldType:    "struct{_ github.com/frk/gosql.Column \"sql:\\\"c.col_a,x\\\"\"; _ github.com/frk/gosql.Column \"sql:\\\"c.col_c,y\\\"\"}",
			FieldName:    "a",
			FieldTag:     "sql:\"c.col_xyz isbetween\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     276,
		},
	}, {
		name: "SelectPostgresTestBAD_WhereBetweenRelationNotFound",
		err: Error{
			Code:         ErrBadAlias,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			ColName:      "col_a",
			ColQualifier: "x",
			TargetName:   "SelectPostgresTestBAD_WhereBetweenRelationNotFound",
			FieldType:    "struct{_ github.com/frk/gosql.Column \"sql:\\\"c.col_b,x\\\"\"; _ github.com/frk/gosql.Column \"sql:\\\"c.col_c,y\\\"\"}",
			FieldName:    "a",
			FieldTag:     "sql:\"x.col_a isbetween\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     287,
		},
	}, {
		name: "SelectPostgresTestBAD_WhereBetweenArgColumnNotFound",
		err: Error{
			Code:         ErrNoColumn,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			RelName:      "column_tests_1",
			RelSchema:    "public",
			ColName:      "col_xyz",
			ColQualifier: "c",
			TargetName:   "SelectPostgresTestBAD_WhereBetweenArgColumnNotFound",
			FieldType:    "github.com/frk/gosql.Column",
			FieldName:    "_",
			FieldTag:     "sql:\"c.col_xyz,x\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     299,
		},
	}, {
		name: "SelectPostgresTestBAD_WhereBetweenArgRelationNotFound",
		err: Error{
			Code:         ErrBadAlias,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			ColName:      "col_b",
			ColQualifier: "x",
			TargetName:   "SelectPostgresTestBAD_WhereBetweenArgRelationNotFound",
			FieldType:    "github.com/frk/gosql.Column",
			FieldName:    "_",
			FieldTag:     "sql:\"x.col_b,x\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     310,
		},
	}, {
		name: "SelectPostgresTestBAD_WhereBetweenComparisonBadArgType",
		err: Error{
			Code:         ErrBadComparisonOperation,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			RelName:      "column_tests_1",
			RelSchema:    "public",
			ColName:      "col_a",
			ColQualifier: "c",
			ColType:      "integer",
			Predicate:    analysis.IsBetween,
			TargetName:   "SelectPostgresTestBAD_WhereBetweenComparisonBadArgType",
			FieldType:    "bool",
			FieldName:    "y",
			FieldTag:     "sql:\"y\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     322,
		},
	}, {
		name: "SelectPostgresTestBAD_OrderByColumnNotFound",
		err: Error{
			Code:         ErrNoColumn,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			RelName:      "column_tests_1",
			RelSchema:    "public",
			ColName:      "col_xyz",
			ColQualifier: "c",
			TargetName:   "SelectPostgresTestBAD_OrderByColumnNotFound",
			FieldType:    "github.com/frk/gosql.OrderBy",
			FieldName:    "_",
			FieldTag:     "sql:\"c.col_a,c.col_xyz\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     330,
		},
	}, {
		name: "SelectPostgresTestBAD_OrderByRelationNotFound",
		err: Error{
			Code:         ErrNoRelation,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			ColName:      "col_a",
			ColQualifier: "x",
			TargetName:   "SelectPostgresTestBAD_OrderByRelationNotFound",
			FieldType:    "github.com/frk/gosql.OrderBy",
			FieldName:    "_",
			FieldTag:     "sql:\"x.col_a\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     336,
		},
	}, {
		name: "InsertPostgresTestBAD_DefaultBadRelationAlias",
		err: Error{
			Code:         ErrBadColumnQualifier,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			RelName:      "column_tests_1",
			RelSchema:    "public",
			ColName:      "col_b",
			ColQualifier: "x",
			TargetName:   "InsertPostgresTestBAD_DefaultBadRelationAlias",
			FieldType:    "github.com/frk/gosql.Default",
			FieldName:    "_",
			FieldTag:     "sql:\"x.col_b\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     342,
		},
	}, {
		name: "InsertPostgresTestBAD_DefaultColumnNotFound",
		err: Error{
			Code:         ErrNoColumn,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			RelName:      "column_tests_1",
			RelSchema:    "public",
			ColName:      "col_xyz",
			ColQualifier: "c",
			TargetName:   "InsertPostgresTestBAD_DefaultColumnNotFound",
			FieldType:    "github.com/frk/gosql.Default",
			FieldName:    "_",
			FieldTag:     "sql:\"c.col_xyz\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     348,
		},
	}, {
		name: "InsertPostgresTestBAD_DefaultNotSet",
		err: Error{
			Code:         ErrNoColumnDefault,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			RelName:      "column_tests_1",
			RelSchema:    "public",
			ColName:      "col_b",
			ColQualifier: "c",
			ColType:      "text",
			TargetName:   "InsertPostgresTestBAD_DefaultNotSet",
			FieldType:    "github.com/frk/gosql.Default",
			FieldName:    "_",
			FieldTag:     "sql:\"c.col_b\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     354,
		},
	}, {
		name: "InsertPostgresTestBAD_ForceColumnNotFound",
		err: Error{
			Code:         ErrNoColumn,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			RelName:      "column_tests_1",
			RelSchema:    "public",
			ColName:      "col_xyz",
			ColQualifier: "c",
			TargetName:   "InsertPostgresTestBAD_ForceColumnNotFound",
			FieldType:    "github.com/frk/gosql.Force",
			FieldName:    "_",
			FieldTag:     "sql:\"c.col_xyz\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     360,
		},
	}, {
		name: "InsertPostgresTestBAD_ForceRelationNotFound",
		err: Error{
			Code:         ErrNoRelation,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			ColName:      "col_a",
			ColQualifier: "x",
			TargetName:   "InsertPostgresTestBAD_ForceRelationNotFound",
			FieldType:    "github.com/frk/gosql.Force",
			FieldName:    "_",
			FieldTag:     "sql:\"x.col_a\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     366,
		},
	}, {
		name: "UpdatePostgresTestBAD_ReturnColumnNotFound",
		err: Error{
			Code:         ErrNoColumn,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			RelName:      "column_tests_1",
			RelSchema:    "public",
			ColName:      "col_xyz",
			ColQualifier: "c",
			TargetName:   "UpdatePostgresTestBAD_ReturnColumnNotFound",
			FieldType:    "github.com/frk/gosql.Return",
			FieldName:    "_",
			FieldTag:     "sql:\"c.col_xyz\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     372,
		},
	}, {
		name: "UpdatePostgresTestBAD_ReturnRelationNotFound",
		err: Error{
			Code:         ErrNoRelation,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			ColName:      "col_a",
			ColQualifier: "x",
			TargetName:   "UpdatePostgresTestBAD_ReturnRelationNotFound",
			FieldType:    "github.com/frk/gosql.Return",
			FieldName:    "_",
			FieldTag:     "sql:\"x.col_a\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     378,
		},
	}, {
		name: "UpdatePostgresTestBAD_ReturnFieldNotFound",
		err: Error{
			Code:         ErrNoColumnField,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			RelName:      "column_tests_1",
			RelSchema:    "public",
			ColName:      "col_d",
			ColQualifier: "c",
			TargetName:   "UpdatePostgresTestBAD_ReturnFieldNotFound",
			FieldType:    "github.com/frk/gosql.Return",
			FieldName:    "_",
			FieldTag:     "sql:\"c.col_d\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     384,
		},
	}, {
		name: "FilterPostgresTestBAD_TextSearchColumnNotFound",
		err: Error{
			Code:       ErrNoColumn,
			PkgPath:    "path/to/test",
			DBName:     "gosql_test_db",
			RelName:    "column_tests_1",
			RelSchema:  "public",
			TargetName: "FilterPostgresTestBAD_TextSearchColumnNotFound",
			FieldType:  "github.com/frk/gosql.TextSearch",
			FieldName:  "_",
			FieldTag:   "sql:\"c.col_xyz\"",
			FileName:   "../testdata/postgres_bad.go",
			FileLine:   390,
		},
	}, {
		name: "FilterPostgresTestBAD_TextSearchRelationNotFound",
		err: Error{
			Code:         ErrBadColumnQualifier,
			PkgPath:      "path/to/test",
			DBName:       "gosql_test_db",
			RelName:      "column_tests_1",
			RelSchema:    "public",
			ColName:      "col_b",
			ColQualifier: "x",
			TargetName:   "FilterPostgresTestBAD_TextSearchRelationNotFound",
			FieldType:    "github.com/frk/gosql.TextSearch",
			FieldName:    "_",
			FieldTag:     "sql:\"x.col_b\"",
			FileName:     "../testdata/postgres_bad.go",
			FileLine:     396,
		},
	}, {
		name: "FilterPostgresTestBAD_TextSearchBadColumnType",
		err: Error{
			Code:       ErrBadColumnType,
			PkgPath:    "path/to/test",
			DBName:     "gosql_test_db",
			RelName:    "column_tests_1",
			RelSchema:  "public",
			ColName:    "col_b",
			ColType:    "text",
			TargetName: "FilterPostgresTestBAD_TextSearchBadColumnType",
			FieldType:  "github.com/frk/gosql.TextSearch",
			FieldName:  "_",
			FieldTag:   "sql:\"c.col_b\"",
			FileName:   "../testdata/postgres_bad.go",
			FileLine:   402,
		},
	}, {
		name: "SelectPostgresTestBAD_RelationColumnNotFound",
		err: Error{
			Code:       ErrNoColumn,
			PkgPath:    "path/to/test",
			DBName:     "gosql_test_db",
			RelName:    "column_tests_1",
			RelSchema:  "public",
			ColName:    "col_xyz",
			TargetName: "SelectPostgresTestBAD_RelationColumnNotFound",
			FieldType:  "string",
			FieldName:  "Xyz",
			FieldTag:   "sql:\"col_xyz\"",
			FileName:   "../testdata/postgres_bad.go",
			FileLine:   408,
		},
	}, {
		name: "InsertPostgresTestBAD_RelationColumnNotFound",
		err: Error{
			Code:       ErrNoColumn,
			PkgPath:    "path/to/test",
			DBName:     "gosql_test_db",
			RelName:    "column_tests_1",
			RelSchema:  "public",
			ColName:    "col_xyz",
			TargetName: "InsertPostgresTestBAD_RelationColumnNotFound",
			FieldType:  "string",
			FieldName:  "XYZ",
			FieldTag:   "sql:\"col_xyz\"",
			FileName:   "../testdata/postgres_bad.go",
			FileLine:   415,
		},
	}, {
		name: "InsertPostgresTestBAD_BadFieldToColumnType",
		err: Error{
			Code:       ErrBadFieldWriteType,
			PkgPath:    "path/to/test",
			DBName:     "gosql_test_db",
			RelName:    "column_tests_1",
			RelSchema:  "public",
			ColName:    "col_c",
			ColType:    "boolean",
			TargetName: "InsertPostgresTestBAD_BadFieldToColumnType",
			FieldType:  "int",
			FieldName:  "B",
			FieldTag:   "sql:\"col_c\"",
			FileName:   "../testdata/postgres_bad.go",
			FileLine:   436,
		},
	}, {
		name: "InsertPostgresTestBAD_ResultColumnNotFound",
		err: Error{
			Code:       ErrNoColumn,
			PkgPath:    "path/to/test",
			DBName:     "gosql_test_db",
			RelName:    "column_tests_1",
			RelSchema:  "public",
			ColName:    "col_xyz",
			TargetName: "InsertPostgresTestBAD_ResultColumnNotFound",
			FieldType:  "int",
			FieldName:  "A",
			FieldTag:   "sql:\"col_xyz\"",
			FileName:   "../testdata/postgres_bad.go",
			FileLine:   444,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := testCheck(tt.name, t)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Errorf("%v - %#v %v", e, err, err)
			}

			_ = info
		})
	}
}