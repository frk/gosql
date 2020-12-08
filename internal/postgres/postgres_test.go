package postgres

import (
	"fmt"
	"log"
	"testing"

	"github.com/frk/compare"
	"github.com/frk/gosql/internal/analysis"
	"github.com/frk/gosql/internal/postgres/oid"
	"github.com/frk/gosql/internal/testutil"

	"github.com/lib/pq"
)

func init() {
	compare.DefaultConfig.ObserveFieldTag = "cmp"
}

var tdata = testutil.ParseTestdata("../testdata")

func testCheck(name string, t *testing.T) (*TargetInfo, error) {
	named, pos := testutil.FindNamedType(name, tdata)
	if named == nil {
		// Stop the test if no type with the given name was found.
		t.Fatal(name, " not found")
		return nil, nil
	}

	info := new(analysis.Info)
	ts, err := analysis.Run(tdata.Fset, named, pos, info)
	if err != nil {
		return nil, err
	}

	return Check(testdb.DB, ts, info)
}

func TestOpen(t *testing.T) {
	tests := []struct {
		dsn      string
		err      error
		printerr bool
	}{
		{
			dsn: "postgres:///gosql_test_db?sslmode=disable",
		}, {
			dsn: "postgres://foo@bar/baz?sslmode=disable",
			err: &dbError{
				Code: errDatabaseOpen,
				DB:   dbInfo{DSN: "postgres://foo@bar/baz?sslmode=disable"},
				Err:  &pq.Error{},
			},
		},
	}

	for _, tt := range tests {
		db, err := Open(tt.dsn)
		if err == nil {
			db.Close()
		}
		if e := compare.Compare(err, tt.err); e != nil {
			t.Errorf("%v - %#v %v", e, err, err)
		}
		if tt.printerr && err != nil {
			fmt.Println(err)
		}
	}
}

func TestCheck(t *testing.T) {
	test_dbinfo := dbInfo{DSN: testdb.DB.dsn, Name: testdb.DB.name, User: testdb.DB.user, SearchPath: testdb.DB.searchpath}
	column_tests_1, err := loadRelation(&checker{db: testdb.DB}, testdb.DB, analysis.RelIdent{"column_tests_1", "", ""}, 0)
	if err != nil {
		log.Fatalf("relation not found: %v\n", err)
	}
	column_tests_2, err := loadRelation(&checker{db: testdb.DB}, testdb.DB, analysis.RelIdent{"column_tests_2", "", ""}, 0)
	if err != nil {
		log.Fatalf("relation not found: %v\n", err)
	}
	test_user, err := loadRelation(&checker{db: testdb.DB}, testdb.DB, analysis.RelIdent{"test_user", "", ""}, 0)
	if err != nil {
		log.Fatalf("relation not found: %v\n", err)
	}
	pgsql_test, err := loadRelation(&checker{db: testdb.DB}, testdb.DB, analysis.RelIdent{"pgsql_test", "", ""}, 0)
	if err != nil {
		log.Fatalf("relation not found: %v\n", err)
	}

	tests := []struct {
		name     string
		err      error
		printerr bool
	}{{
		name: "SelectPostgresTestOK_Simple",
		err:  nil,
	}, {
		name:     "SelectPostgresTestOK_Enums",
		printerr: true,
		err:      nil,
	}, {
		name:     "InsertPostgresTestOK_Enums",
		printerr: true,
		err:      nil,
	}, {
		name:     "SelectPostgresTestOK_CustomTypePointer",
		printerr: true,
		err:      nil,
	}, {
		name:     "InsertPostgresTestOK_CustomTypePointer",
		printerr: true,
		err:      nil,
	}, {
		name: "SelectPostgresTestBAD_NoRelation",
		err: &dbError{
			Code: errRelationUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_NoRelation",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 11,
				},
			},
			Field: fieldInfo{
				Name: "Columns",
				Type: "path/to/test.CT1",
				Tag:  `rel:"norel"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 12,
				},
			},
			Rel: relInfo{
				Id: analysis.RelIdent{"norel", "", ""},
			},
		},
	}, {
		name: "DeletePostgresTestBAD_JoinNoRelation",
		err: &dbError{
			Code: errRelationUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "DeletePostgresTestBAD_JoinNoRelation",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 16,
				},
			},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.Relation",
				Tag:  `sql:"norel:b"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 19,
				},
			},
			Rel: relInfo{
				Id: analysis.RelIdent{"norel", "b", ""},
			},
		},
	}, {
		name: "DeletePostgresTestBAD_JoinNoRelation2",
		err: &dbError{
			Code: errRelationUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "DeletePostgresTestBAD_JoinNoRelation2",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 27,
				},
			},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.LeftJoin",
				Tag:  `sql:"norel:c,c.b_id = b.id"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 31,
				},
			},
			Rel: relInfo{
				Id: analysis.RelIdent{"norel", "c", ""},
			},
		},
	}, {
		name: "SelectPostgresTestBAD_JoinNoColumn",
		err: &dbError{
			Code: errColumnUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_JoinNoColumn",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 39,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_2", "b", ""}, Relation: column_tests_2},
			Col: colInfo{Id: analysis.ColIdent{"nocol", "b"}},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.LeftJoin",
				Tag:  `sql:"column_tests_2:b,b.nocol = a.nocol"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 42,
				},
			},
		},
	}, {
		name: "SelectPostgresTestBAD_JoinNoColumn2",
		err: &dbError{
			Code: errColumnUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_JoinNoColumn2",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 47,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col: colInfo{Id: analysis.ColIdent{"nocol", "a"}},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.LeftJoin",
				Tag:  `sql:"column_tests_2:b,b.col_foo = a.nocol"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 50,
				},
			},
		},
	}, {
		name: "SelectPostgresTestBAD_JoinBadUnaryBoolColumn",
		err: &dbError{
			Code: errPredicateOperandBool,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_JoinBadUnaryBoolColumn",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 55,
				},
			},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.LeftJoin",
				Tag:  `sql:"column_tests_2:b,b.col_foo istrue"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 58,
				},
			},
			Rel:  relInfo{Id: analysis.RelIdent{"column_tests_2", "", "public"}, Relation: column_tests_2},
			Col:  colInfo{Id: analysis.ColIdent{"col_foo", "b"}, Column: findRelColumn(column_tests_2, "col_foo")},
			Pred: analysis.IsTrue,
		},
	}, {
		name: "SelectPostgresTestBAD_JoinBadUnaryNullColumn",
		err: &dbError{
			Code: errPredicateOperandNull,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_JoinBadUnaryNullColumn",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 63,
				},
			},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.LeftJoin",
				Tag:  `sql:"column_tests_2:b,b.col_baz isnull"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 66,
				},
			},
			Rel:  relInfo{Id: analysis.RelIdent{"column_tests_2", "", "public"}, Relation: column_tests_2},
			Col:  colInfo{Id: analysis.ColIdent{"col_baz", "b"}, Column: findRelColumn(column_tests_2, "col_baz")},
			Pred: analysis.IsNull,
		},
	}, {
		name: "SelectPostgresTestBAD_JoinBadLiteralExpression",
		err: &dbError{
			Code: errPredicateLiteralExpr,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_JoinBadLiteralExpression",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 71,
				},
			},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.LeftJoin",
				Tag:  `sql:"column_tests_2:b,b.col_baz = 'foo'bar "`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 74,
				},
			},
			Rel:    relInfo{Id: analysis.RelIdent{"column_tests_2", "", "public"}, Relation: column_tests_2},
			Col:    colInfo{Id: analysis.ColIdent{"col_baz", "b"}, Column: findRelColumn(column_tests_2, "col_baz")},
			RHSLit: exprInfo{Expr: "'foo'bar"},
			Pred:   analysis.IsEQ,
		},
	}, {
		name: "SelectPostgresTestBAD_JoinBadQuantifierColumnType",
		err: &dbError{
			Code: errPredicateOperandQuantifier,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_JoinBadQuantifierColumnType",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 79,
				},
			},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.LeftJoin",
				Tag:  `sql:"column_tests_2:b,b.col_foo >any a.col_a"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 82,
				},
			},
			Rel:    relInfo{Id: analysis.RelIdent{"column_tests_2", "", "public"}, Relation: column_tests_2},
			Col:    colInfo{Id: analysis.ColIdent{"col_foo", "b"}, Column: findRelColumn(column_tests_2, "col_foo")},
			RHSCol: colInfo{Id: analysis.ColIdent{"col_a", "a"}, Column: findRelColumn(column_tests_1, "col_a")},
			Pred:   analysis.IsGT,
			Quant:  analysis.QuantAny,
		},
	}, {
		name: "SelectPostgresTestBAD_JoinBadComparisonOperandType",
		err: &dbError{
			Code: errColumnComparison,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_JoinBadComparisonOperandType",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 87,
				},
			},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.LeftJoin",
				Tag:  `sql:"column_tests_2:b,b.col_baz < 'baz'"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 90,
				},
			},
			Rel:    relInfo{Id: analysis.RelIdent{"column_tests_2", "", "public"}, Relation: column_tests_2},
			Col:    colInfo{Id: analysis.ColIdent{"col_baz", "b"}, Column: findRelColumn(column_tests_2, "col_baz")},
			RHSLit: exprInfo{Expr: "'baz'", Type: testdb.DB.catalog.Types[oid.Unknown]},
			Pred:   analysis.IsLT,
		},
	}, {
		name: "InsertPostgresTestBAD_OnConflictNoColumn",
		err: &dbError{
			Code: errColumnUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "InsertPostgresTestBAD_OnConflictNoColumn",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 95,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col: colInfo{Id: analysis.ColIdent{"col_xyz", "c"}},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.Column",
				Tag:  "sql:\"c.col_xyz\"",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 98,
				},
			},
		},
	}, {
		name: "InsertPostgresTestBAD_OnConflictColumnNoIndexMatch",
		err: &dbError{
			Code: errOnConflictIndexColumnsUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "InsertPostgresTestBAD_OnConflictColumnNoIndexMatch",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 103,
				},
			},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.Column",
				Tag:  `sql:"c.col_a,c.col_b"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 106,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
		},
	}, {
		name: "InsertPostgresTestBAD_OnConflictNoIndex",
		err: &dbError{
			Code: errOnConflictIndexUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "InsertPostgresTestBAD_OnConflictNoIndex",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 111,
				},
			},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.Index",
				Tag:  "sql:\"some_index\"",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 114,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
		},
	}, {
		name: "InsertPostgresTestBAD_OnConflictNoUniqueIndex",
		err: &dbError{
			Code: errOnConflictIndexNotUnique,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "InsertPostgresTestBAD_OnConflictNoUniqueIndex",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 119,
				},
			},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.Index",
				Tag:  "sql:\"column_tests_2_nonunique_index\"",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 122,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_2", "", "public"}, Relation: column_tests_2},
		},
	}, {
		name: "InsertPostgresTestBAD_OnConflictNoConstraint",
		err: &dbError{
			Code: errOnConflictConstraintUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "InsertPostgresTestBAD_OnConflictNoConstraint",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 127,
				},
			},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.Constraint",
				Tag:  "sql:\"some_constraint\"",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 130,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
		},
	}, {
		name: "InsertPostgresTestBAD_OnConflictNoUniqueConstraint",
		err: &dbError{
			Code: errOnConflictConstraintNotUnique,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "InsertPostgresTestBAD_OnConflictNoUniqueConstraint",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 135,
				},
			},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.Constraint",
				Tag:  "sql:\"column_tests_2_nonunique_constraint\"",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 138,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_2", "", "public"}, Relation: column_tests_2},
		},
	}, {
		name: "InsertPostgresTestBAD_OnConflictUpdateColumnNotFound",
		err: &dbError{
			Code: errColumnUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "InsertPostgresTestBAD_OnConflictUpdateColumnNotFound",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 143,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_2", "", "public"}, Relation: column_tests_2},
			Col: colInfo{Id: analysis.ColIdent{"col_a", "c"}},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.Update",
				Tag:  "sql:\"c.col_a,c.col_b,c.col_xyz\"",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 147,
				},
			},
		},
	}, {
		name: "SelectPostgresTestBAD_WhereFieldColumnNotFound",
		err: &dbError{
			Code: errColumnUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_WhereFieldColumnNotFound",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 152,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col: colInfo{Id: analysis.ColIdent{"id", "c"}},
			Field: fieldInfo{
				Name: "Id",
				Type: "int",
				Tag:  "sql:\"c.id\"",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 155,
				},
			},
		},
	}, {
		name: "SelectPostgresTestBAD_WhereCannotCompareTypes",
		err: &dbError{
			Code: errColumnFieldComparison,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_WhereCannotCompareTypes",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 160,
				},
			},
			Field: fieldInfo{
				Name: "D",
				Type: "float64",
				Tag:  "sql:\"c.col_e ~\"",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 163,
				},
			},
			Rel:  relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col:  colInfo{Id: analysis.ColIdent{"col_e", "c"}, Column: findRelColumn(column_tests_1, "col_e")},
			Pred: analysis.IsMatch,
		},
	}, {
		name: "SelectPostgresTestBAD_WhereColumnTypeForFuncname",
		err: &dbError{
			Code: errProcedureUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_WhereColumnTypeForFuncname",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 168,
				},
			},
			Field: fieldInfo{
				Name: "D",
				Type: "float64",
				Tag:  "sql:\"c.col_d,@lower\"",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 171,
				},
			},
			Rel:  relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col:  colInfo{Id: analysis.ColIdent{"col_d", "c"}, Column: findRelColumn(column_tests_1, "col_d")},
			Pred: analysis.IsEQ,
			Func: "lower",
		},
	}, {
		name: "SelectPostgresTestBAD_WhereColumnNotFound",
		err: &dbError{
			Code: errColumnUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_WhereColumnNotFound",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 176,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col: colInfo{Id: analysis.ColIdent{"col_xyz", "c"}},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.Column",
				Tag:  "sql:\"c.col_xyz istrue\"",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 179,
				},
			},
		},
	}, {
		name: "SelectPostgresTestBAD_WhereColumnBadBoolOp",
		err: &dbError{
			Code: errPredicateOperandBool,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_WhereColumnBadBoolOp",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 184,
				},
			},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.Column",
				Tag:  "sql:\"c.col_a istrue\"",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 187,
				},
			},
			Rel:  relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col:  colInfo{Id: analysis.ColIdent{"col_a", "c"}, Column: findRelColumn(column_tests_1, "col_a")},
			Pred: analysis.IsTrue,
		},
	}, {
		name: "SelectPostgresTestBAD_WhereColumnNotFoundRHS",
		err: &dbError{
			Code: errColumnUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_WhereColumnNotFoundRHS",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 200,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col: colInfo{Id: analysis.ColIdent{"col_xyz", "c"}},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.Column",
				Tag:  "sql:\"c.col_a = c.col_xyz\"",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 203,
				},
			},
		},
	}, {
		name: "SelectPostgresTestBAD_WhereColumnBadLiteralExpression",
		err: &dbError{
			Code: errPredicateLiteralExpr,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_WhereColumnBadLiteralExpression",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 208,
				},
			},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.Column",
				Tag:  `sql:"c.col_a = 123abc"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 211,
				},
			},
			Rel:    relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col:    colInfo{Id: analysis.ColIdent{"col_a", "c"}, Column: findRelColumn(column_tests_1, "col_a")},
			RHSLit: exprInfo{Expr: "123abc"},
			Pred:   analysis.IsEQ,
		},
	}, {
		name: "SelectPostgresTestBAD_WhereColumnBadTypeForQuantifier",
		err: &dbError{
			Code: errPredicateOperandArray,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_WhereColumnBadTypeForQuantifier",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 216,
				},
			},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.Column",
				Tag:  `sql:"c.col_a isin c.col_b"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 219,
				},
			},
			Rel:    relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col:    colInfo{Id: analysis.ColIdent{"col_a", "c"}, Column: findRelColumn(column_tests_1, "col_a")},
			RHSCol: colInfo{Id: analysis.ColIdent{"col_b", "c"}, Column: findRelColumn(column_tests_1, "col_b")},
			Pred:   analysis.IsIn,
		},
	}, {
		name: "SelectPostgresTestBAD_WhereColumnBadTypeComparison",
		err: &dbError{
			Code: errColumnComparison,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_WhereColumnBadTypeComparison",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 224,
				},
			},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.Column",
				Tag:  `sql:"c.col_a = c.col_b"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 227,
				},
			},
			Rel:    relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col:    colInfo{Id: analysis.ColIdent{"col_a", "c"}, Column: findRelColumn(column_tests_1, "col_a")},
			RHSCol: colInfo{Id: analysis.ColIdent{"col_b", "c"}, Column: findRelColumn(column_tests_1, "col_b")},
			Pred:   analysis.IsEQ,
		},
	}, {
		name: "SelectPostgresTestBAD_WhereBetweenColumnNotFound",
		err: &dbError{
			Code: errColumnUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_WhereBetweenColumnNotFound",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 232,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col: colInfo{Id: analysis.ColIdent{"col_xyz", "c"}},
			Field: fieldInfo{
				Name: "a",
				Type: `struct{_ github.com/frk/gosql.Column "sql:\"c.col_a,x\""; _ github.com/frk/gosql.Column "sql:\"c.col_c,y\""}`,
				Tag:  "sql:\"c.col_xyz isbetween\"",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 235,
				},
			},
			WBField: fieldInfo{
				Name: "a",
				Type: `struct{_ github.com/frk/gosql.Column "sql:\"c.col_a,x\""; _ github.com/frk/gosql.Column "sql:\"c.col_c,y\""}`,
				Tag:  `sql:"c.col_xyz isbetween"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 235,
				},
			},
		},
	}, {
		name: "SelectPostgresTestBAD_WhereBetweenArgColumnNotFound",
		err: &dbError{
			Code: errColumnUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_WhereBetweenArgColumnNotFound",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 243,
				},
			},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.Column",
				Tag:  "sql:\"c.col_xyz,x\"",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 247,
				},
			},
			WBField: fieldInfo{
				Name: "a",
				Type: `struct{_ github.com/frk/gosql.Column "sql:\"c.col_xyz,x\""; _ github.com/frk/gosql.Column "sql:\"c.col_c,y\""}`,
				Tag:  `sql:"c.col_a isbetween"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 246,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col: colInfo{Id: analysis.ColIdent{"col_xyz", "c"}},
		},
	}, {
		name: "SelectPostgresTestBAD_WhereBetweenComparisonBadArgType",
		err: &dbError{
			Code: errBetweenFieldComparison,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_WhereBetweenComparisonBadArgType",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 254,
				},
			},
			Field: fieldInfo{
				Name: "y",
				Type: "bool",
				Tag:  "sql:\"y\"",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 259,
				},
			},
			WBField: fieldInfo{
				Name: "a",
				Type: `struct{x int "sql:\"x\""; y bool "sql:\"y\""}`,
				Tag:  `sql:"c.col_a isbetween"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 257,
				},
			},
			Rel:  relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col:  colInfo{Id: analysis.ColIdent{"col_a", "c"}, Column: findRelColumn(column_tests_1, "col_a")},
			Pred: analysis.IsLTE,
		},
	}, {
		name: "SelectPostgresTestBAD_OrderByColumnNotFound",
		err: &dbError{
			Code: errColumnUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_OrderByColumnNotFound",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 265,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col: colInfo{Id: analysis.ColIdent{"col_xyz", "c"}},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.OrderBy",
				Tag:  "sql:\"c.col_a,c.col_xyz\"",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 267,
				},
			},
		},
	}, {
		name: "InsertPostgresTestBAD_DefaultColumnNotFound",
		err: &dbError{
			Code: errColumnUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "InsertPostgresTestBAD_DefaultColumnNotFound",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 271,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col: colInfo{Id: analysis.ColIdent{"col_xyz", "c"}},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.Default",
				Tag:  "sql:\"c.col_xyz\"",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 273,
				},
			},
		},
	}, {
		name: "InsertPostgresTestBAD_DefaultNotSet",
		err: &dbError{
			Code: errColumnDefaultUnset,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "InsertPostgresTestBAD_DefaultNotSet",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 277,
				},
			},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.Default",
				Tag:  `sql:"c.col_b"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 279,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col: colInfo{Id: analysis.ColIdent{"col_b", "c"}, Column: findRelColumn(column_tests_1, "col_b")},
		},
	}, {
		name: "UpdatePostgresTestBAD_ReturnColumnNotFound",
		err: &dbError{
			Code: errColumnUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "UpdatePostgresTestBAD_ReturnColumnNotFound",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 283,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col: colInfo{Id: analysis.ColIdent{"col_xyz", "c"}},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.Return",
				Tag:  "sql:\"c.col_xyz\"",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 285,
				},
			},
		},
	}, {
		name: "FilterPostgresTestBAD_TextSearchColumnNotFound",
		err: &dbError{
			Code: errColumnUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "FilterPostgresTestBAD_TextSearchColumnNotFound",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 289,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col: colInfo{Id: analysis.ColIdent{"col_xyz", "c"}},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.TextSearch",
				Tag:  "sql:\"c.col_xyz\"",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 291,
				},
			},
		},
	}, {
		name: "FilterPostgresTestBAD_TextSearchBadColumnType",
		err: &dbError{
			Code: errColumnTextSearchType,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "FilterPostgresTestBAD_TextSearchBadColumnType",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 296,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col: colInfo{Id: analysis.ColIdent{"col_b", "c"}, Column: findRelColumn(column_tests_1, "col_b")},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.TextSearch",
				Tag:  "sql:\"c.col_b\"",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 298,
				},
			},
		},
	}, {
		name: "SelectPostgresTestBAD_RelationColumnNotFound",
		err: &dbError{
			Code: errColumnUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_RelationColumnNotFound",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 303,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col: colInfo{Id: analysis.ColIdent{"col_xyz", ""}},
			Field: fieldInfo{
				Name: "Xyz",
				Type: "string",
				Tag:  "sql:\"col_xyz\"",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 305,
				},
			},
		},
	}, {
		name: "InsertPostgresTestBAD_RelationColumnNotFound",
		err: &dbError{
			Code: errColumnUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "InsertPostgresTestBAD_RelationColumnNotFound",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 310,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col: colInfo{Id: analysis.ColIdent{"col_xyz", ""}},
			Field: fieldInfo{
				Name: "XYZ",
				Type: "string",
				Tag:  "sql:\"col_xyz\"",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 312,
				},
			},
		},
	}, {
		name: "InsertPostgresTestBAD_BadFieldToColumnType",
		err: &dbError{
			Code: errColumnFieldTypeWrite,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "InsertPostgresTestBAD_BadFieldToColumnType",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 331,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col: colInfo{Id: analysis.ColIdent{"col_c", ""}, Column: findRelColumn(column_tests_1, "col_c")},
			Field: fieldInfo{
				Name: "B",
				Type: "int",
				Tag:  `sql:"col_c"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 333,
				},
			},
		},
	}, {
		name: "InsertPostgresTestBAD_ResultColumnNotFound",
		err: &dbError{
			Code: errColumnUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "InsertPostgresTestBAD_ResultColumnNotFound",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 338,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col: colInfo{Id: analysis.ColIdent{"col_xyz", ""}},
			Field: fieldInfo{
				Name: "A",
				Type: "int",
				Tag:  "sql:\"col_xyz\"",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 341,
				},
			},
		},
	}, {
		name: "SelectPostgresTestBAD_NoSchema",
		err: &dbError{
			Code: errRelationUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_NoSchema",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 346,
				},
			},
			Field: fieldInfo{
				Name: "Columns",
				Type: "path/to/test.CT1",
				Tag:  `rel:"noschema.column_tests_1:c"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 347,
				},
			},
			Rel: relInfo{
				Id: analysis.RelIdent{"column_tests_1", "c", "noschema"},
			},
		},
	}, {
		name: "SelectPostgresTestBAD_RelationColumnNotFound2",
		err: &dbError{
			Code: errColumnUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_RelationColumnNotFound2",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 351,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col: colInfo{Id: analysis.ColIdent{"col_xyz", ""}},
			Field: fieldInfo{
				Name: "XYZ",
				Type: "string",
				Tag:  "sql:\"col_xyz\"",
				File: fileInfo{
					Name: "../testdata/types.go",
					Line: 31,
				},
			},
		},
	}, {
		name: "FilterPostgresTestBAD_BadFieldWriteType",
		err: &dbError{
			Code: errColumnFieldTypeWrite,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "FilterPostgresTestBAD_BadFieldWriteType",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 356,
				},
			},
			Field: fieldInfo{
				Name: "Metadata",
				Type: "func()",
				Tag:  `sql:"metadata2"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 358,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"test_user", "", "public"}, Relation: test_user},
			Col: colInfo{Column: findRelColumn(test_user, "metadata2")},
		},
	}, {
		name: "FilterPostgresTestBAD_BadFieldWriteType2",
		err: &dbError{
			Code: errColumnFieldTypeWrite,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "FilterPostgresTestBAD_BadFieldWriteType2",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 364,
				},
			},
			Field: fieldInfo{
				Name: "Envelope",
				Type: "chan struct{}",
				Tag:  `sql:"envelope"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 366,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"test_user", "", "public"}, Relation: test_user},
			Col: colInfo{Column: findRelColumn(test_user, "envelope")},
		},
	}, {
		name: "FilterPostgresTestBAD_BadFieldWriteType3",
		err: &dbError{
			Code: errColumnFieldTypeWrite,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "FilterPostgresTestBAD_BadFieldWriteType3",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 372,
				},
			},
			Field: fieldInfo{
				Name: "Lines",
				Type: "float64",
				Tag:  `sql:"col_linearr"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 374,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"pgsql_test", "", "public"}, Relation: pgsql_test},
			Col: colInfo{Column: findRelColumn(pgsql_test, "col_linearr")},
		},
	}, {
		name: "SelectPostgresTestBAD_WhereLiteralBadTypeForQuantifier",
		err: &dbError{
			Code: errPredicateOperandArray,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_WhereLiteralBadTypeForQuantifier",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 380,
				},
			},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.Column",
				Tag:  `sql:"c.col_a notin 'foo bar'"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 383,
				},
			},
			Rel:    relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col:    colInfo{Id: analysis.ColIdent{"col_a", "c"}, Column: findRelColumn(column_tests_1, "col_a")},
			RHSLit: exprInfo{Expr: "'foo bar'", Type: testdb.DB.catalog.Types[oid.Unknown]},
			Pred:   analysis.NotIn,
		},
	}, {
		name: "SelectPostgresTestBAD_WhereUnknownFunc",
		err: &dbError{
			Code: errProcedureUnknown,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_WhereUnknownFunc",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 388,
				},
			},
			Field: fieldInfo{
				Name: "D",
				Type: "float64",
				Tag:  `sql:"c.col_d,@unknown_func"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 391,
				},
			},
			Rel:  relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col:  colInfo{Id: analysis.ColIdent{"col_d", "c"}, Column: findRelColumn(column_tests_1, "col_d")},
			Pred: analysis.IsEQ,
			Func: "unknown_func",
		},
	}, {
		name: "InsertPostgresTestBAD_DefaultOption",
		err: &dbError{
			Code: errColumnDefaultUnset,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "InsertPostgresTestBAD_DefaultOption",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 396,
				},
			},
			Field: fieldInfo{
				Name: "B",
				Type: "string",
				Tag:  `sql:"col_b,default"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 398,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col: colInfo{Id: analysis.ColIdent{"col_b", ""}, Column: findRelColumn(column_tests_1, "col_b")},
		},
	}, {
		name: "SelectPostgresTestBAD_ColumnTypeToBadField",
		err: &dbError{
			Code: errColumnFieldTypeRead,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_ColumnTypeToBadField",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 403,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col: colInfo{Id: analysis.ColIdent{"col_c", ""}, Column: findRelColumn(column_tests_1, "col_c")},
			Field: fieldInfo{
				Name: "B",
				Type: "int",
				Tag:  `sql:"col_c"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 405,
				},
			},
		},
	}, {
		name: "SelectPostgresTestBAD_ColumnTypeToBadField2",
		err: &dbError{
			Code: errColumnFieldTypeRead,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_ColumnTypeToBadField2",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 410,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"test_user", "", "public"}, Relation: test_user},
			Col: colInfo{Id: analysis.ColIdent{"envelope", ""}, Column: findRelColumn(test_user, "envelope")},
			Field: fieldInfo{
				Name: "Envelope",
				Type: "chan struct{}",
				Tag:  `sql:"envelope"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 412,
				},
			},
		},
	}, {
		name: "SelectPostgresTestBAD_ColumnTypeToBadField3",
		err: &dbError{
			Code: errColumnFieldTypeRead,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_ColumnTypeToBadField3",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 417,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"pgsql_test", "", "public"}, Relation: pgsql_test},
			Col: colInfo{Id: analysis.ColIdent{"col_linearr", ""}, Column: findRelColumn(pgsql_test, "col_linearr")},
			Field: fieldInfo{
				Name: "Lines",
				Type: "float64",
				Tag:  `sql:"col_linearr"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 419,
				},
			},
		},
	}, {
		name: "SelectPostgresTestBAD_BadBetweenColumnComparison",
		err: &dbError{
			Code: errBetweenColumnComparison,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "SelectPostgresTestBAD_BadBetweenColumnComparison",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 424,
				},
			},
			Rel:    relInfo{Id: analysis.RelIdent{"column_tests_1", "", "public"}, Relation: column_tests_1},
			Col:    colInfo{Id: analysis.ColIdent{"col_a", "c"}, Column: findRelColumn(column_tests_1, "col_a")},
			RHSCol: colInfo{Id: analysis.ColIdent{"col_c", "c"}, Column: findRelColumn(column_tests_1, "col_c")},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.Column",
				Tag:  `sql:"c.col_c,y"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 429,
				},
			},
			WBField: fieldInfo{
				Name: "a",
				Type: `struct{_ github.com/frk/gosql.Column "sql:\"c.col_d,x\""; _ github.com/frk/gosql.Column "sql:\"c.col_c,y\""}`,
				Tag:  `sql:"c.col_a isbetween"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 427,
				},
			},
			Pred: analysis.IsLTE,
		},
	}, {
		name: "InsertPostgresTestBAD_OnConflictIndexColumnsNotUnique",
		err: &dbError{
			Code: errOnConflictIndexColumnsNotUnique,
			DB:   test_dbinfo,
			Target: targetInfo{
				Pkg:  "path/to/test",
				Name: "InsertPostgresTestBAD_OnConflictIndexColumnsNotUnique",
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 435,
				},
			},
			Field: fieldInfo{
				Name: "_",
				Type: "github.com/frk/gosql.Column",
				Tag:  `sql:"c.col_indkey1,c.col_indkey2,c.col_indkey3"`,
				File: fileInfo{
					Name: "../testdata/postgres_bad.go",
					Line: 438,
				},
			},
			Rel: relInfo{Id: analysis.RelIdent{"column_tests_2", "", "public"}, Relation: column_tests_2},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := testCheck(tt.name, t)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Errorf("%v - %#v %v", e, err, err)
			}
			_ = info

			if tt.printerr && err != nil {
				fmt.Println(err)
			}
		})
	}
}

type anTargetStruct struct {
	ts analysis.TargetStruct
}

// helper func to retrieve analysis info on specific test type
func _analyzeTargetStruct(name string) anTargetStruct {
	named, pos := testutil.FindNamedType(name, tdata)
	if named == nil {
		panic(name + " not found")
	}

	info := new(analysis.Info)
	ts, err := analysis.Run(tdata.Fset, named, pos, info)
	if err != nil {
		panic(err)
	}

	return anTargetStruct{ts}
}

func (a anTargetStruct) relField() *analysis.RelField {
	switch v := a.ts.(type) {
	case *analysis.QueryStruct:
		return v.Rel
	case *analysis.FilterStruct:
		return v.Rel
	}
	return nil
}

func (a anTargetStruct) resultField() *analysis.ResultField {
	switch v := a.ts.(type) {
	case *analysis.QueryStruct:
		return v.Result
	}
	return nil
}
