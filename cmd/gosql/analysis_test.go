package main

import (
	"testing"

	"github.com/frk/compare"
	"github.com/frk/gosql/internal/testutil"
	"github.com/frk/tagutil"
)

var tdata = testutil.ParseTestdata("../../testdata")

func runAnalysis(name string, t *testing.T) (*analyzer, error) {
	named := testutil.FindNamedType(name, tdata)
	if named == nil {
		// Stop the test if no type with the given name was found.
		t.Fatal(name, " not found")
		return nil, nil
	}

	a := &analyzer{fset: tdata.Fset, named: named}
	if err := a.run(); err != nil {
		return nil, err
	}
	return a, nil
}

func TestAnalysis_queryStruct(t *testing.T) {
	// for reuse, analyzed common.User typeInfo
	commonUserTypeinfo := typeInfo{
		name:       "User",
		kind:       typeKindStruct,
		pkgPath:    "github.com/frk/gosql/testdata/common",
		pkgName:    "common",
		pkgLocal:   "common",
		isImported: true,
	}

	commonUserFields := []*fieldInfo{{
		name:       "Id",
		typ:        typeInfo{kind: typeKindInt},
		isExported: true,
		colId:      colId{name: "id"},
		tag:        tagutil.Tag{"sql": {"id"}},
	}, {
		name:       "Email",
		typ:        typeInfo{kind: typeKindString},
		isExported: true,
		colId:      colId{name: "email"},
		tag:        tagutil.Tag{"sql": {"email"}},
	}, {
		name:       "FullName",
		typ:        typeInfo{kind: typeKindString},
		isExported: true,
		colId:      colId{name: "full_name"},
		tag:        tagutil.Tag{"sql": {"full_name"}},
	}, {
		name: "CreatedAt",
		typ: typeInfo{
			name:              "Time",
			kind:              typeKindStruct,
			pkgPath:           "time",
			pkgName:           "time",
			pkgLocal:          "time",
			isImported:        true,
			isJSONMarshaler:   true,
			isJSONUnmarshaler: true,
		},
		isExported: true,
		colId:      colId{name: "created_at"},
		tag:        tagutil.Tag{"sql": {"created_at"}},
	}}

	reldummyslice := &dataField{
		name:  "Rel",
		relId: relId{name: "relation_a", alias: "a"},
		data: dataType{
			typeInfo: typeInfo{
				name:     "T",
				kind:     typeKindStruct,
				pkgPath:  "path/to/test",
				pkgName:  "testdata",
				pkgLocal: "testdata",
			},
			isSlice: true,
			fields: []*fieldInfo{{
				typ:        typeInfo{kind: typeKindString},
				name:       "F",
				isExported: true,
				tag:        tagutil.Tag{"sql": {"f"}},
				colId:      colId{name: "f"},
			}},
		},
	}

	dummyrecord := dataType{
		typeInfo: typeInfo{
			name:     "T",
			kind:     typeKindStruct,
			pkgPath:  "path/to/test",
			pkgName:  "testdata",
			pkgLocal: "testdata",
		},
		fields: []*fieldInfo{{
			typ:        typeInfo{kind: typeKindString},
			name:       "F",
			isExported: true,
			tag:        tagutil.Tag{"sql": {"f"}},
			colId:      colId{name: "f"},
		}},
	}

	tests := []struct {
		name string
		want *queryStruct
		err  error
	}{{
		name: "InsertAnalysisTestBAD_NoDataField",
		err: analysisError{
			errorCode:  errNoTargetField,
			pkgPath:    "path/to/test",
			targetName: "InsertAnalysisTestBAD_NoDataField",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   11,
		},
	}, {
		name: "InsertAnalysisTestBAD3",
		err: analysisError{
			errorCode:  errDataType,
			pkgPath:    "path/to/test",
			targetName: "InsertAnalysisTestBAD3",
			blockName:  "",
			fieldType:  "string",
			fieldName:  "User",
			tagValue:   "",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   17,
		},
	}, {
		name: "DeleteAnalysisTestBAD_BadRelId",
		err: analysisError{
			errorCode:  errBadRelIdTagValue,
			pkgPath:    "path/to/test",
			targetName: "DeleteAnalysisTestBAD_BadRelId",
			blockName:  "",
			fieldType:  "path/to/test.T",
			fieldName:  "Rel",
			tagValue:   "foo.123:bar",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   22,
		},
	}, {
		name: "SelectAnalysisTestBAD_MultipleRelTags",
		err: analysisError{
			errorCode:  errRelTagConflict,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_MultipleRelTags",
			blockName:  "",
			fieldType:  "path/to/test.T",
			fieldName:  "Rel2",
			tagValue:   "",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   28,
		},
	}, {
		name: "DeleteAnalysisTestBAD_IllegalCountField",
		err: analysisError{
			errorCode:  errIllegalField,
			blockName:  "",
			fieldType:  "int",
			fieldName:  "Count",
			targetName: "DeleteAnalysisTestBAD_IllegalCountField",
			pkgPath:    "path/to/test",
			fileLine:   33,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "UpdateAnalysisTestBAD_IllegalExistsField",
		err: analysisError{
			errorCode:  errIllegalField,
			blockName:  "",
			fieldType:  "bool",
			fieldName:  "Exists",
			targetName: "UpdateAnalysisTestBAD_IllegalExistsField",
			pkgPath:    "path/to/test",
			fileLine:   38,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "InsertAnalysisTestBAD_IllegalNotExistsField",
		err: analysisError{
			errorCode:  errIllegalField,
			blockName:  "",
			fieldType:  "bool",
			fieldName:  "NotExists",
			targetName: "InsertAnalysisTestBAD_IllegalNotExistsField",
			pkgPath:    "path/to/test",
			fileLine:   43,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "SelectAnalysisTestBAD_IllegalRelationDirective",
		err: analysisError{
			errorCode:  errIllegalField,
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Relation",
			fieldName:  "_",
			targetName: "SelectAnalysisTestBAD_IllegalRelationDirective",
			pkgPath:    "path/to/test",
			fileLine:   48,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "SelectAnalysisTestBAD_UnnamedBaseStructType",
		err: analysisError{
			errorCode:  errDataType,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_UnnamedBaseStructType",
			fieldType:  "[]*struct{}",
			fieldName:  "Rel",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   53,
		},
	}, {
		name: "SelectAnalysisTestBAD_IllegalAllDirective",
		err: analysisError{
			errorCode:  errIllegalField,
			blockName:  "",
			fieldType:  "github.com/frk/gosql.All",
			fieldName:  "_",
			targetName: "SelectAnalysisTestBAD_IllegalAllDirective",
			pkgPath:    "path/to/test",
			fileLine:   59,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "InsertAnalysisTestBAD_IllegalAllDirective",
		err: analysisError{
			errorCode:  errIllegalField,
			blockName:  "",
			fieldType:  "github.com/frk/gosql.All",
			fieldName:  "_",
			targetName: "InsertAnalysisTestBAD_IllegalAllDirective",
			pkgPath:    "path/to/test",
			fileLine:   65,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "UpdateAnalysisTestBAD_ConflictWhereProducer",
		err: analysisError{
			errorCode:  errFieldConflict,
			pkgPath:    "path/to/test",
			targetName: "UpdateAnalysisTestBAD_ConflictWhereProducer",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.All",
			fieldName:  "_",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   74,
		},
	}, {
		name: "DeleteAnalysisTestBAD_IllegalDefaultDirective",
		err: analysisError{
			errorCode:  errIllegalField,
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Default",
			fieldName:  "_",
			targetName: "DeleteAnalysisTestBAD_IllegalDefaultDirective",
			pkgPath:    "path/to/test",
			fileLine:   80,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "UpdateAnalysisTestBAD_EmptyDefaultDirectiveCollist",
		err: analysisError{
			errorCode:  errNoTagValue,
			fieldType:  "github.com/frk/gosql.Default",
			fieldName:  "_",
			targetName: "UpdateAnalysisTestBAD_EmptyDefaultDirectiveCollist",
			pkgPath:    "path/to/test",
			fileLine:   86,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "SelectAnalysisTestBAD_IllegalForceDirective",
		err: analysisError{
			errorCode:  errIllegalField,
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Force",
			fieldName:  "_",
			targetName: "SelectAnalysisTestBAD_IllegalForceDirective",
			pkgPath:    "path/to/test",
			fileLine:   92,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "UpdateAnalysisTestBAD_BadForceDirectiveColId",
		err: analysisError{
			errorCode:  errBadColIdTagValue,
			pkgPath:    "path/to/test",
			targetName: "UpdateAnalysisTestBAD_BadForceDirectiveColId",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Force",
			fieldName:  "_",
			tagValue:   "1234",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   98,
		},
	}, {
		name: "DeleteAnalysisTestBAD_ConflictResultProducer",
		err: analysisError{
			errorCode:  errFieldConflict,
			pkgPath:    "path/to/test",
			targetName: "DeleteAnalysisTestBAD_ConflictResultProducer",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Return",
			fieldName:  "_",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   111,
		},
	}, {
		name: "UpdateAnalysisTestBAD_EmptyReturnDirectiveCollist",
		err: analysisError{
			errorCode:  errNoTagValue,
			fieldType:  "github.com/frk/gosql.Return",
			fieldName:  "_",
			targetName: "UpdateAnalysisTestBAD_EmptyReturnDirectiveCollist",
			pkgPath:    "path/to/test",
			fileLine:   117,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "InsertAnalysisTestBAD_IllegalLimitField",
		err: analysisError{
			errorCode:  errIllegalField,
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Limit",
			fieldName:  "_",
			targetName: "InsertAnalysisTestBAD_IllegalLimitField",
			pkgPath:    "path/to/test",
			fileLine:   123,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "UpdateAnalysisTestBAD_IllegalOffsetField",
		err: analysisError{
			errorCode:  errIllegalField,
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Offset",
			fieldName:  "_",
			targetName: "UpdateAnalysisTestBAD_IllegalOffsetField",
			pkgPath:    "path/to/test",
			fileLine:   129,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "InsertAnalysisTestBAD_IllegalOrderByDirective",
		err: analysisError{
			errorCode:  errIllegalField,
			blockName:  "",
			fieldType:  "github.com/frk/gosql.OrderBy",
			fieldName:  "_",
			targetName: "InsertAnalysisTestBAD_IllegalOrderByDirective",
			pkgPath:    "path/to/test",
			fileLine:   135,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "DeleteAnalysisTestBAD_IllegalOverrideDirective",
		err: analysisError{
			errorCode:  errIllegalField,
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Override",
			fieldName:  "_",
			targetName: "DeleteAnalysisTestBAD_IllegalOverrideDirective",
			pkgPath:    "path/to/test",
			fileLine:   141,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "SelectAnalysisTestBAD_IllegalTextSearchDirective",
		err: analysisError{
			errorCode:  errIllegalField,
			blockName:  "",
			fieldType:  "github.com/frk/gosql.TextSearch",
			fieldName:  "_",
			targetName: "SelectAnalysisTestBAD_IllegalTextSearchDirective",
			pkgPath:    "path/to/test",
			fileLine:   147,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "SelectAnalysisTestBAD_IllegalColumnDirective",
		err: analysisError{
			errorCode:  errIllegalField,
			fieldType:  "github.com/frk/gosql.Column",
			fieldName:  "_",
			targetName: "SelectAnalysisTestBAD_IllegalColumnDirective",
			pkgPath:    "path/to/test",
			fileLine:   153,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "InsertAnalysisTestBAD_IllegalWhereBlock",
		err: analysisError{
			errorCode:  errIllegalField,
			blockName:  "",
			fieldType:  "struct{Id int \"sql:\\\"id\\\"\"}",
			fieldName:  "Where",
			targetName: "InsertAnalysisTestBAD_IllegalWhereBlock",
			pkgPath:    "path/to/test",
			fileLine:   159,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "UpdateAnalysisTestBAD_IllegalJoinBlock",
		err: analysisError{
			errorCode:  errIllegalField,
			blockName:  "",
			fieldType:  "struct{_ github.com/frk/gosql.Relation}",
			fieldName:  "Join",
			targetName: "UpdateAnalysisTestBAD_IllegalJoinBlock",
			pkgPath:    "path/to/test",
			fileLine:   167,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "DeleteAnalysisTestBAD_IllegalFromBlock",
		err: analysisError{
			errorCode:  errIllegalField,
			blockName:  "",
			fieldType:  "struct{_ github.com/frk/gosql.Relation}",
			fieldName:  "From",
			targetName: "DeleteAnalysisTestBAD_IllegalFromBlock",
			pkgPath:    "path/to/test",
			fileLine:   175,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "SelectAnalysisTestBAD_IllegalUsingBlock",
		err: analysisError{
			errorCode:  errIllegalField,
			blockName:  "",
			fieldType:  "struct{_ github.com/frk/gosql.Relation}",
			fieldName:  "Using",
			targetName: "SelectAnalysisTestBAD_IllegalUsingBlock",
			pkgPath:    "path/to/test",
			fileLine:   183,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "UpdateAnalysisTestBAD_IllegalOnConflictBlock",
		err: analysisError{
			errorCode:  errIllegalField,
			blockName:  "",
			fieldType:  "struct{}",
			fieldName:  "OnConflict",
			targetName: "UpdateAnalysisTestBAD_IllegalOnConflictBlock",
			pkgPath:    "path/to/test",
			fileLine:   191,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "SelectAnalysisTestBAD_IllegalResultField",
		err: analysisError{
			errorCode:  errIllegalField,
			blockName:  "",
			fieldType:  "path/to/test.T",
			fieldName:  "Result",
			targetName: "SelectAnalysisTestBAD_IllegalResultField",
			pkgPath:    "path/to/test",
			fileLine:   199,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "SelectAnalysisTestBAD_ConflictLimitProducer",
		err: analysisError{
			errorCode:  errFieldConflict,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_ConflictLimitProducer",
			blockName:  "",
			fieldType:  "int",
			fieldName:  "Limit",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   206,
		},
	}, {
		name: "SelectAnalysisTestBAD_ConflictOffsetProducer",
		err: analysisError{
			errorCode:  errFieldConflict,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_ConflictOffsetProducer",
			blockName:  "",
			fieldType:  "int",
			fieldName:  "Offset",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   213,
		},
	}, {
		name: "SelectAnalysisTestBAD_IllegalRowsAffectedField",
		err: analysisError{
			errorCode:  errIllegalField,
			blockName:  "",
			fieldType:  "int",
			fieldName:  "RowsAffected",
			targetName: "SelectAnalysisTestBAD_IllegalRowsAffectedField",
			pkgPath:    "path/to/test",
			fileLine:   219,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "InsertAnalysisTestBAD_IllegalFilterField",
		err: analysisError{
			errorCode:  errIllegalField,
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Filter",
			fieldName:  "F",
			targetName: "InsertAnalysisTestBAD_IllegalFilterField",
			pkgPath:    "path/to/test",
			fileLine:   225,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "SelectAnalysisTestBAD_ConflictWhereProducer",
		err: analysisError{
			errorCode:  errFieldConflict,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_ConflictWhereProducer",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Filter",
			fieldName:  "F",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   234,
		},
	}, {
		name: "DeleteAnalysisTestBAD_ConflictWhereProducer",
		err: analysisError{
			errorCode:  errFieldConflict,
			pkgPath:    "path/to/test",
			targetName: "DeleteAnalysisTestBAD_ConflictWhereProducer",
			blockName:  "",
			fieldType:  "path/to/test.myerrorhandler",
			fieldName:  "erh",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   241,
		},
	}, {
		name: "SelectAnalysisTestBAD_IteratorWithTooManyMethods",
		err: analysisError{
			errorCode:  errIterType,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_IteratorWithTooManyMethods",
			fieldType:  "path/to/test.badIterator",
			fieldName:  "Rel",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   255,
		},
	}, {
		name: "SelectAnalysisTestBAD_IteratorWithBadSignature",
		err: analysisError{
			errorCode:  errIterType,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_IteratorWithBadSignature",
			fieldType:  "func(*github.com/frk/gosql/testdata/common.User) int",
			fieldName:  "Rel",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   260,
		},
	}, {
		name: "SelectAnalysisTestBAD_IteratorWithBadSignatureIface",
		err: analysisError{
			errorCode:  errIterType,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_IteratorWithBadSignatureIface",
			fieldType:  "path/to/test.badIterator2",
			fieldName:  "Rel",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   265,
		},
	}, {
		name: "SelectAnalysisTestBAD_IteratorWithUnexportedMethod",
		err: analysisError{
			errorCode:  errIterType,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_IteratorWithUnexportedMethod",
			fieldType:  "github.com/frk/gosql/testdata/common.BadIterator",
			fieldName:  "Rel",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   270,
		},
	}, {
		name: "SelectAnalysisTestBAD_IteratorWithUnnamedArgument",
		err: analysisError{
			errorCode:  errIterType,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_IteratorWithUnnamedArgument",
			fieldType:  "func(*struct{}) error",
			fieldName:  "Rel",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   275,
		},
	}, {
		name: "SelectAnalysisTestBAD_IteratorWithNonStructArgument",
		err: analysisError{
			errorCode:  errIterType,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_IteratorWithNonStructArgument",
			fieldType:  "func(*path/to/test.notstruct) error",
			fieldName:  "Rel",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   280,
		},
	}, {
		name: "InsertAnalysisTestBAD_BadRelfiedlStructBaseType",
		err: analysisError{
			errorCode:  errDataType,
			pkgPath:    "path/to/test",
			targetName: "InsertAnalysisTestBAD_BadRelfiedlStructBaseType",
			fieldType:  "[]*path/to/test.notstruct",
			fieldName:  "Rel",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   287,
		},
	}, {
		name: "UpdateAnalysisTestBAD_BadRelTypeFieldColId",
		err: analysisError{
			errorCode:  errBadColIdTagValue,
			pkgPath:    "path/to/test",
			targetName: "UpdateAnalysisTestBAD_BadRelTypeFieldColId",
			blockName:  "",
			fieldType:  "string",
			fieldName:  "Foo",
			tagValue:   "1234",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   293,
		},
	}, {
		name: "UpdateAnalysisTestBAD_ConflictWhereProducer2",
		err: analysisError{
			errorCode:  errFieldConflict,
			pkgPath:    "path/to/test",
			targetName: "UpdateAnalysisTestBAD_ConflictWhereProducer2",
			blockName:  "",
			fieldType:  "struct{Id int \"sql:\\\"id\\\"\"}",
			fieldName:  "Where",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   301,
		},
	}, {
		name: "DeleteAnalysisTestBAD_BadWhereBlockType",
		err: analysisError{
			errorCode:  errFieldBlock,
			pkgPath:    "path/to/test",
			targetName: "DeleteAnalysisTestBAD_BadWhereBlockType",
			blockName:  "",
			fieldType:  "[]string",
			fieldName:  "Where",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   309,
		},
	}, {
		name: "SelectAnalysisTestBAD_BadBoolTagValue",
		err: analysisError{
			errorCode:  errBadBoolTagValue,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_BadBoolTagValue",
			blockName:  "Where",
			fieldType:  "string",
			fieldName:  "Name",
			tagValue:   "abc",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   317,
		},
	}, {
		name: "SelectAnalysisTestBAD_BadNestedWhereBlockType",
		err: analysisError{
			errorCode:  errFieldBlock,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_BadNestedWhereBlockType",
			blockName:  "Where",
			fieldType:  "path/to/test.notstruct",
			fieldName:  "X",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   326,
		},
	}, {
		name: "SelectAnalysisTestBAD_BadColumnExpressionLHS",
		err: analysisError{
			errorCode:  errBadColIdTagValue,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_BadColumnExpressionLHS",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Column",
			fieldName:  "_",
			tagValue:   "123",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   334,
		},
	}, {
		name: "SelectAnalysisTestBAD_BadColumnPredicateCombo",
		err: analysisError{
			errorCode:  errIllegalPredicateQuantifier,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_BadColumnPredicateCombo",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Column",
			fieldName:  "_",
			tagValue:   "x isin any y",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   342,
		},
	}, {
		name: "DeleteAnalysisTestBAD_BadColumnExpressionLHS",
		err: analysisError{
			errorCode:  errBadColIdTagValue,
			pkgPath:    "path/to/test",
			targetName: "DeleteAnalysisTestBAD_BadColumnExpressionLHS",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Column",
			fieldName:  "_",
			tagValue:   "123",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   350,
		},
	}, {
		name: "UpdateAnalysisTestBAD_BadUnaryOp",
		err: analysisError{
			errorCode:  errBadUnaryPredicate,
			pkgPath:    "path/to/test",
			targetName: "UpdateAnalysisTestBAD_BadUnaryOp",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Column",
			fieldName:  "_",
			tagValue:   "x <=",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   358,
		},
	}, {
		name: "UpdateAnalysisTestBAD_ExtraQuantifier",
		err: analysisError{
			errorCode:  errIllegalPredicateQuantifier,
			pkgPath:    "path/to/test",
			targetName: "UpdateAnalysisTestBAD_ExtraQuantifier",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Column",
			fieldName:  "_",
			tagValue:   "x isnull any",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   366,
		},
	}, {
		name: "SelectAnalysisTestBAD_BadBetweenFieldType",
		err: analysisError{
			errorCode:  errBadBetweenPredicate,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_BadBetweenFieldType",
			fieldType:  "path/to/test.notstruct",
			fieldName:  "between",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   374,
		},
	}, {
		name: "SelectAnalysisTestBAD_BadBetweenFieldType2",
		err: analysisError{
			errorCode:  errBadBetweenPredicate,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_BadBetweenFieldType2",
			fieldType:  "struct{x int; y int; z int}",
			fieldName:  "between",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   382,
		},
	}, {
		name: "SelectAnalysisTestBAD_BadBetweenArgColId",
		err: analysisError{
			errorCode:  errBadColIdTagValue,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_BadBetweenArgColId",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Column",
			fieldName:  "_",
			tagValue:   "123",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   394,
		},
	}, {
		name: "SelectAnalysisTestBAD_NoBetweenXYArg",
		err: analysisError{
			errorCode:  errBadBetweenPredicate,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_NoBetweenXYArg",
			fieldType:  "struct{_ github.com/frk/gosql.Column \"sql:\\\"a.bar\\\"\"; _ github.com/frk/gosql.Column \"sql:\\\"a.baz,y\\\"\"}",
			fieldName:  "between",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   403,
		},
	}, {
		name: "SelectAnalysisTestBAD_BadBetweenColId",
		err: analysisError{
			errorCode:  errBadColIdTagValue,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_BadBetweenColId",
			blockName:  "",
			fieldType:  "struct{_ github.com/frk/gosql.Column \"sql:\\\"a.bar,x\\\"\"; _ github.com/frk/gosql.Column \"sql:\\\"a.baz,y\\\"\"}",
			fieldName:  "between",
			tagValue:   "123",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   414,
		},
	}, {
		name: "DeleteAnalysisTestBAD_BadWhereFieldColId",
		err: analysisError{
			errorCode:  errBadColIdTagValue,
			pkgPath:    "path/to/test",
			targetName: "DeleteAnalysisTestBAD_BadWhereFieldColId",
			blockName:  "",
			fieldType:  "int",
			fieldName:  "Id",
			tagValue:   "123",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   425,
		},
	}, {
		name: "DeleteAnalysisTestBAD_BadWhereFieldPredicateCombo",
		err: analysisError{
			errorCode:  errIllegalPredicateQuantifier,
			pkgPath:    "path/to/test",
			targetName: "DeleteAnalysisTestBAD_BadWhereFieldPredicateCombo",
			fieldType:  "int",
			fieldName:  "Id",
			tagValue:   "a.id notin any",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   433,
		},
	}, {
		name: "DeleteAnalysisTestBAD_IllegalWhereFieldUnaryPredicate",
		err: analysisError{
			errorCode:  errIllegalUnaryPredicate,
			pkgPath:    "path/to/test",
			targetName: "DeleteAnalysisTestBAD_IllegalWhereFieldUnaryPredicate",
			blockName:  "",
			fieldType:  "int",
			fieldName:  "Id",
			tagValue:   "a.id istrue",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   441,
		},
	}, {
		name: "UpdateAnalysisTestBAD_BadWhereFieldTypeForQuantifier",
		err: analysisError{
			errorCode:  errIllegalPredicateQuantifier,
			pkgPath:    "path/to/test",
			targetName: "UpdateAnalysisTestBAD_BadWhereFieldTypeForQuantifier",
			fieldType:  "int",
			fieldName:  "Id",
			tagValue:   "a.id = any",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   449,
		},
	}, {
		name: "SelectAnalysisTestBAD_BadJoinBlockType",
		err: analysisError{
			errorCode:  errFieldBlock,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_BadJoinBlockType",
			blockName:  "",
			fieldType:  "path/to/test.notstruct",
			fieldName:  "Join",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   456,
		},
	}, {
		name: "SelectAnalysisTestBAD_IllegalJoinBlockRelationDirective",
		err: analysisError{
			errorCode:  errIllegalField,
			blockName:  "Join",
			fieldType:  "github.com/frk/gosql.Relation",
			fieldName:  "_",
			targetName: "SelectAnalysisTestBAD_IllegalJoinBlockRelationDirective",
			pkgPath:    "path/to/test",
			fileLine:   463,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "DeleteAnalysisTestBAD_ConflictRelationDirective",
		err: analysisError{
			errorCode:  errFieldConflict,
			pkgPath:    "path/to/test",
			targetName: "DeleteAnalysisTestBAD_ConflictRelationDirective",
			fieldType:  "github.com/frk/gosql.Relation",
			fieldName:  "_",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   472,
		},
	}, {
		name: "UpdateAnalysisTestBAD_BadFromRelationRelId",
		err: analysisError{
			errorCode:  errBadRelIdTagValue,
			pkgPath:    "path/to/test",
			targetName: "UpdateAnalysisTestBAD_BadFromRelationRelId",
			fieldType:  "github.com/frk/gosql.Relation",
			fieldName:  "_",
			tagValue:   "123",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   480,
		},
	}, {
		name: "SelectAnalysisTestBAD_BadJoinDirectiveRelId",
		err: analysisError{
			errorCode:  errBadRelIdTagValue,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_BadJoinDirectiveRelId",
			fieldType:  "github.com/frk/gosql.LeftJoin",
			fieldName:  "_",
			tagValue:   "123",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   488,
		},
	}, {
		name: "SelectAnalysisTestBAD_BadJoinDirectiveExpressionColId",
		err: analysisError{
			errorCode:  errBadColIdTagValue,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_BadJoinDirectiveExpressionColId",
			fieldType:  "github.com/frk/gosql.LeftJoin",
			fieldName:  "_",
			tagValue:   "123",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   496,
		},
	}, {
		name: "SelectAnalysisTestBAD_BadJoinDirectiveExpressionPredicate",
		err: analysisError{
			errorCode:  errBadUnaryPredicate,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_BadJoinDirectiveExpressionPredicate",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.LeftJoin",
			fieldName:  "_",
			tagValue:   "b.foo =",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   504,
		},
	}, {
		name: "SelectAnalysisTestBAD_BadJoinDirectiveExpressionExtraQuantifier",
		err: analysisError{
			errorCode:  errIllegalPredicateQuantifier,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_BadJoinDirectiveExpressionExtraQuantifier",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.LeftJoin",
			fieldName:  "_",
			tagValue:   "b.foo isnull any",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   512,
		},
	}, {
		name: "SelectAnalysisTestBAD_BadJoinDirectiveExpressionPredicateCombo",
		err: analysisError{
			errorCode:  errIllegalPredicateQuantifier,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_BadJoinDirectiveExpressionPredicateCombo",
			fieldType:  "github.com/frk/gosql.LeftJoin",
			fieldName:  "_",
			tagValue:   "b.foo isin any a.bar",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   520,
		},
	}, {
		name: "DeleteAnalysisTestBAD_IllegalJoinBlockDirective",
		err: analysisError{
			errorCode:  errIllegalField,
			blockName:  "Using",
			fieldType:  "github.com/frk/gosql.Column",
			fieldName:  "_",
			targetName: "DeleteAnalysisTestBAD_IllegalJoinBlockDirective",
			pkgPath:    "path/to/test",
			fileLine:   528,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "InsertAnalysisTestBAD_BadOnConflictBlockType",
		err: analysisError{
			errorCode:  errFieldBlock,
			pkgPath:    "path/to/test",
			targetName: "InsertAnalysisTestBAD_BadOnConflictBlockType",
			blockName:  "",
			fieldType:  "path/to/test.notstruct",
			fieldName:  "OnConflict",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   535,
		},
	}, {
		name: "InsertAnalysisTestBAD_ConflictOnConflictBlockTargetProducer",
		err: analysisError{
			errorCode:  errFieldConflict,
			pkgPath:    "path/to/test",
			targetName: "InsertAnalysisTestBAD_ConflictOnConflictBlockTargetProducer",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Column",
			fieldName:  "_",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   543,
		},
	}, {
		name: "InsertAnalysisTestBAD_ConflictOnConflictBlockTargetProducer2",
		err: analysisError{
			errorCode:  errFieldConflict,
			pkgPath:    "path/to/test",
			targetName: "InsertAnalysisTestBAD_ConflictOnConflictBlockTargetProducer2",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Index",
			fieldName:  "_",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   552,
		},
	}, {
		name: "InsertAnalysisTestBAD_ConflictOnConflictBlockTargetProducer3",
		err: analysisError{
			errorCode:  errFieldConflict,
			pkgPath:    "path/to/test",
			targetName: "InsertAnalysisTestBAD_ConflictOnConflictBlockTargetProducer3",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Constraint",
			fieldName:  "_",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   561,
		},
	}, {
		name: "InsertAnalysisTestBAD_ConflictOnConflictBlockActionProducer",
		err: analysisError{
			errorCode:  errFieldConflict,
			pkgPath:    "path/to/test",
			targetName: "InsertAnalysisTestBAD_ConflictOnConflictBlockActionProducer",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Ignore",
			fieldName:  "_",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   571,
		},
	}, {
		name: "InsertAnalysisTestBAD_ConflictOnConflictBlockActionProducer2",
		err: analysisError{
			errorCode:  errFieldConflict,
			pkgPath:    "path/to/test",
			targetName: "InsertAnalysisTestBAD_ConflictOnConflictBlockActionProducer2",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Update",
			fieldName:  "_",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   581,
		},
	}, {
		name: "InsertAnalysisTestBAD_BadOnConflictColumnTargetValue",
		err: analysisError{
			errorCode:  errBadColIdTagValue,
			pkgPath:    "path/to/test",
			targetName: "InsertAnalysisTestBAD_BadOnConflictColumnTargetValue",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Column",
			fieldName:  "_",
			tagValue:   "a.1234",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   589,
		},
	}, {
		name: "InsertAnalysisTestBAD_BadOnConflictIndexTargetIdent",
		err: analysisError{
			errorCode:  errBadTagValue,
			pkgPath:    "path/to/test",
			targetName: "InsertAnalysisTestBAD_BadOnConflictIndexTargetIdent",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Index",
			fieldName:  "_",
			tagValue:   "1234",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   597,
		},
	}, {
		name: "InsertAnalysisTestBAD_BadOnConflictConstraintTargetIdent",
		err: analysisError{
			errorCode:  errBadTagValue,
			pkgPath:    "path/to/test",
			targetName: "InsertAnalysisTestBAD_BadOnConflictConstraintTargetIdent",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Constraint",
			fieldName:  "_",
			tagValue:   "1234",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   605,
		},
	}, {
		name: "InsertAnalysisTestBAD_BadOnConflictUpdateActionCollist",
		err: analysisError{
			errorCode:  errBadColIdTagValue,
			pkgPath:    "path/to/test",
			targetName: "InsertAnalysisTestBAD_BadOnConflictUpdateActionCollist",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Update",
			fieldName:  "_",
			tagValue:   "a.1234",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   614,
		},
	}, {
		name: "InsertAnalysisTestBAD_IllegalOnConflictDirective",
		err: analysisError{
			errorCode:  errIllegalField,
			blockName:  "OnConflict",
			fieldType:  "github.com/frk/gosql.LeftJoin",
			fieldName:  "_",
			targetName: "InsertAnalysisTestBAD_IllegalOnConflictDirective",
			pkgPath:    "path/to/test",
			fileLine:   622,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "InsertAnalysisTestBAD_NoOnConflictTarget",
		err: analysisError{
			errorCode:  errNoTargetField,
			pkgPath:    "path/to/test",
			targetName: "InsertAnalysisTestBAD_NoOnConflictTarget",
			blockName:  "",
			fieldType:  "struct{_ github.com/frk/gosql.Update \"sql:\\\"a.foo,a.bar\\\"\"}",
			fieldName:  "OnConflict",
			tagValue:   "",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   629,
		},
	}, {
		name: "SelectAnalysisTestBAD_BadLimitFieldType",
		err: analysisError{
			errorCode:  errFieldType,
			fieldType:  "string",
			fieldName:  "Limit",
			targetName: "SelectAnalysisTestBAD_BadLimitFieldType",
			pkgPath:    "path/to/test",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   637,
		},
	}, {
		name: "SelectAnalysisTestBAD_NoLimitDirectiveValue",
		err: analysisError{
			errorCode:  errNoTagValue,
			fieldType:  "github.com/frk/gosql.Limit",
			fieldName:  "_",
			targetName: "SelectAnalysisTestBAD_NoLimitDirectiveValue",
			pkgPath:    "path/to/test",
			fileLine:   643,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "SelectAnalysisTestBAD_BadLimitDirectiveValue",
		err: analysisError{
			errorCode:  errBadTagValue,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_BadLimitDirectiveValue",
			fieldType:  "github.com/frk/gosql.Limit",
			fieldName:  "_",
			tagValue:   "abc",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   649,
		},
	}, {
		name: "SelectAnalysisTestBAD_BadOffsetFieldType",
		err: analysisError{
			errorCode:  errFieldType,
			fieldType:  "string",
			fieldName:  "Offset",
			targetName: "SelectAnalysisTestBAD_BadOffsetFieldType",
			pkgPath:    "path/to/test",
			fileLine:   655,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "SelectAnalysisTestBAD_NoOffsetDirectiveValue",
		err: analysisError{
			errorCode:  errNoTagValue,
			fieldType:  "github.com/frk/gosql.Offset",
			fieldName:  "_",
			targetName: "SelectAnalysisTestBAD_NoOffsetDirectiveValue",
			pkgPath:    "path/to/test",
			fileLine:   661,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "SelectAnalysisTestBAD_BadOffsetDirectiveValue",
		err: analysisError{
			errorCode:  errBadTagValue,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_BadOffsetDirectiveValue",
			fieldType:  "github.com/frk/gosql.Offset",
			fieldName:  "_",
			tagValue:   "abc",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   667,
		},
	}, {
		name: "SelectAnalysisTestBAD_EmptyOrderByDirectiveCollist",
		err: analysisError{
			errorCode:  errNoTagValue,
			fieldType:  "github.com/frk/gosql.OrderBy",
			fieldName:  "_",
			targetName: "SelectAnalysisTestBAD_EmptyOrderByDirectiveCollist",
			pkgPath:    "path/to/test",
			fileLine:   673,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "SelectAnalysisTestBAD_BadOrderByDirectiveNullsOrderValue",
		err: analysisError{
			errorCode:  errBadTagValue,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_BadOrderByDirectiveNullsOrderValue",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.OrderBy",
			fieldName:  "_",
			tagValue:   "a.id:nullsthird",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   679,
		},
	}, {
		name: "SelectAnalysisTestBAD_BadOrderByDirectiveCollist",
		err: analysisError{
			errorCode:  errBadColIdTagValue,
			pkgPath:    "path/to/test",
			targetName: "SelectAnalysisTestBAD_BadOrderByDirectiveCollist",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.OrderBy",
			fieldName:  "_",
			tagValue:   "a.1234",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   685,
		},
	}, {
		name: "InsertAnalysisTestBAD_BadOverrideDirectiveKindValue",
		err: analysisError{
			errorCode:  errBadTagValue,
			pkgPath:    "path/to/test",
			targetName: "InsertAnalysisTestBAD_BadOverrideDirectiveKindValue",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.Override",
			fieldName:  "_",
			tagValue:   "foo",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   691,
		},
	}, {
		name: "UpdateAnalysisTestBAD_ConflictResultProducer",
		err: analysisError{
			errorCode:  errFieldConflict,
			pkgPath:    "path/to/test",
			targetName: "UpdateAnalysisTestBAD_ConflictResultProducer",
			blockName:  "",
			fieldType:  "[]path/to/test.T",
			fieldName:  "Result",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   698,
		},
	}, {
		name: "UpdateAnalysisTestBAD_BadResultFieldType",
		err: analysisError{
			errorCode:  errDataType,
			pkgPath:    "path/to/test",
			targetName: "UpdateAnalysisTestBAD_BadResultFieldType",
			fieldType:  "[]path/to/test.notstruct",
			fieldName:  "Result",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   704,
		},
	}, {
		name: "DeleteAnalysisTestBAD_ConflictResultProducer2",
		err: analysisError{
			errorCode:  errFieldConflict,
			pkgPath:    "path/to/test",
			targetName: "DeleteAnalysisTestBAD_ConflictResultProducer2",
			blockName:  "",
			fieldType:  "int",
			fieldName:  "RowsAffected",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   711,
		},
	}, {
		name: "DeleteAnalysisTestBAD_BadRowsAffecteFieldType",
		err: analysisError{
			errorCode:  errFieldType,
			fieldType:  "string",
			fieldName:  "RowsAffected",
			targetName: "DeleteAnalysisTestBAD_BadRowsAffecteFieldType",
			pkgPath:    "path/to/test",
			fileLine:   717,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "InsertAnalysisTestOK1",
		want: &queryStruct{
			name: "InsertAnalysisTestOK1",
			kind: queryKindInsert,
			dataField: &dataField{
				name:  "UserRec",
				relId: relId{name: "users_table"},
				data: dataType{
					typeInfo:  commonUserTypeinfo,
					fields:    commonUserFields,
					isPointer: true,
				},
			},
		},
	}, {
		name: "InsertAnalysisTestOK2",
		want: &queryStruct{
			name: "InsertAnalysisTestOK2",
			kind: queryKindInsert,
			dataField: &dataField{
				name:  "UserRec",
				relId: relId{name: "users_table"},
				data: dataType{
					typeInfo: typeInfo{
						kind: typeKindStruct,
					},
					fields: []*fieldInfo{{
						name:       "Name3",
						typ:        typeInfo{kind: typeKindString},
						isExported: true,
						colId:      colId{name: "name"},
						tag:        tagutil.Tag{"sql": {"name"}},
					}},
				},
			},
		},
	}, {
		name: "SelectAnalysisTestOK3",
		want: &queryStruct{
			name: "SelectAnalysisTestOK3",
			kind: queryKindSelect,
			dataField: &dataField{
				name:  "User",
				relId: relId{name: "users_table"},
				data: dataType{
					typeInfo:  commonUserTypeinfo,
					fields:    commonUserFields,
					isPointer: true,
					isIter:    true,
				},
			},
		},
	}, {
		name: "SelectAnalysisTestOK4",
		want: &queryStruct{
			name: "SelectAnalysisTestOK4",
			kind: queryKindSelect,
			dataField: &dataField{
				name:  "User",
				relId: relId{name: "users_table"},
				data: dataType{
					typeInfo:  commonUserTypeinfo,
					fields:    commonUserFields,
					isPointer: true,
					isIter:    true,
				},
			},
		},
	}, {
		name: "SelectAnalysisTestOK5",
		want: &queryStruct{
			name: "SelectAnalysisTestOK5",
			kind: queryKindSelect,
			dataField: &dataField{
				name:  "User",
				relId: relId{name: "users_table"},
				data: dataType{
					typeInfo:   commonUserTypeinfo,
					fields:     commonUserFields,
					isPointer:  true,
					isIter:     true,
					iterMethod: "Fn",
				},
			},
		},
	}, {
		name: "SelectAnalysisTestOK6",
		want: &queryStruct{
			name: "SelectAnalysisTestOK6",
			kind: queryKindSelect,
			dataField: &dataField{
				name:  "User",
				relId: relId{name: "users_table"},
				data: dataType{
					typeInfo:   commonUserTypeinfo,
					fields:     commonUserFields,
					isPointer:  true,
					isIter:     true,
					iterMethod: "Fn",
				},
			},
		},
	}, {
		name: "SelectAnalysisTestOK7",
		want: &queryStruct{
			name: "SelectAnalysisTestOK7",
			kind: queryKindSelect,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a"},
				data: dataType{
					typeInfo: typeInfo{
						kind: typeKindStruct,
					},
					fields: []*fieldInfo{{
						name:  "a",
						typ:   typeInfo{kind: typeKindInt},
						colId: colId{name: "a"},
						tag:   tagutil.Tag{"sql": {"a", "pk"}},
					}, {
						name:      "b",
						typ:       typeInfo{kind: typeKindInt},
						colId:     colId{name: "b"},
						tag:       tagutil.Tag{"sql": {"b", "nullempty"}},
						nullEmpty: true,
					}, {
						name:     "c",
						typ:      typeInfo{kind: typeKindInt},
						colId:    colId{name: "c"},
						tag:      tagutil.Tag{"sql": {"c", "ro", "json"}},
						readOnly: true,
					}, {
						name:      "d",
						typ:       typeInfo{kind: typeKindInt},
						colId:     colId{name: "d"},
						tag:       tagutil.Tag{"sql": {"d", "wo"}},
						writeOnly: true,
					}, {
						name:   "e",
						typ:    typeInfo{kind: typeKindInt},
						colId:  colId{name: "e"},
						tag:    tagutil.Tag{"sql": {"e", "add"}},
						useAdd: true,
					}, {
						name:        "f",
						typ:         typeInfo{kind: typeKindInt},
						colId:       colId{name: "f"},
						tag:         tagutil.Tag{"sql": {"f", "coalesce"}},
						useCoalesce: true,
					}, {
						name:          "g",
						typ:           typeInfo{kind: typeKindInt},
						colId:         colId{name: "g"},
						tag:           tagutil.Tag{"sql": {"g", "coalesce(-1)"}},
						useCoalesce:   true,
						coalesceValue: "-1",
					}},
				},
			},
		},
	}, {
		name: "InsertAnalysisTestOK8",
		want: &queryStruct{
			name: "InsertAnalysisTestOK8",
			kind: queryKindInsert,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a"},
				data: dataType{
					typeInfo: typeInfo{
						kind: typeKindStruct,
					},
					fields: []*fieldInfo{{
						name: "Val",
						path: []*fieldNode{
							{
								name:         "Foobar",
								tag:          tagutil.Tag{"sql": {">foo_"}},
								typeName:     "Foo",
								typePkgPath:  "github.com/frk/gosql/testdata/common",
								typePkgName:  "common",
								typePkgLocal: "common",
								isExported:   true,
								isImported:   true,
							},
							{
								name:         "Bar",
								tag:          tagutil.Tag{"sql": {">bar_"}},
								typeName:     "Bar",
								typePkgPath:  "github.com/frk/gosql/testdata/common",
								typePkgName:  "common",
								typePkgLocal: "common",
								isImported:   true,
								isExported:   true,
							},
							{
								name:         "Baz",
								tag:          tagutil.Tag{"sql": {">baz_"}},
								typeName:     "Baz",
								typePkgPath:  "github.com/frk/gosql/testdata/common",
								typePkgName:  "common",
								typePkgLocal: "common",
								isExported:   true,
								isEmbedded:   true,
								isImported:   true,
							},
						},
						isExported: true,
						typ:        typeInfo{kind: typeKindString},
						colId:      colId{name: "foo_bar_baz_val"},
						tag:        tagutil.Tag{"sql": {"val"}},
					}, {
						name: "Val",
						path: []*fieldNode{{
							name:         "Foobar",
							tag:          tagutil.Tag{"sql": {">foo_"}},
							typeName:     "Foo",
							typePkgPath:  "github.com/frk/gosql/testdata/common",
							typePkgName:  "common",
							typePkgLocal: "common",
							isExported:   true,
							isImported:   true,
						}, {
							name:         "Baz",
							tag:          tagutil.Tag{"sql": {">baz_"}},
							typeName:     "Baz",
							typePkgPath:  "github.com/frk/gosql/testdata/common",
							typePkgName:  "common",
							typePkgLocal: "common",
							isImported:   true,
							isExported:   true,
							isEmbedded:   false,
							isPointer:    true,
						}},
						isExported: true,
						typ:        typeInfo{kind: typeKindString},
						colId:      colId{name: "foo_baz_val"},
						tag:        tagutil.Tag{"sql": {"val"}},
					}},
				},
			},
		},
	}, {
		name: "DeleteAnalysisTestOK9",
		want: &queryStruct{
			name: "DeleteAnalysisTestOK9",
			kind: queryKindDelete,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a"},
				data:  dataType{typeInfo: typeInfo{kind: typeKindStruct}},
			},
			whereBlock: &whereBlock{name: "Where", conds: []*searchCondition{{
				cond: &searchConditionField{
					name:  "ID",
					typ:   typeInfo{kind: typeKindInt},
					colId: colId{name: "id"},
					pred:  isEQ,
				},
			}}},
		},
	}, {
		name: "DeleteAnalysisTestOK10",
		want: &queryStruct{
			name: "DeleteAnalysisTestOK10",
			kind: queryKindDelete,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a"},
				data:  dataType{typeInfo: typeInfo{kind: typeKindStruct}},
			},
			whereBlock: &whereBlock{name: "Where", conds: []*searchCondition{
				{cond: &searchConditionColumn{colId: colId{name: "column_a"}, pred: notNull}},
				{bool: boolAnd, cond: &searchConditionColumn{colId: colId{name: "column_b"}, pred: isNull}},
				{bool: boolOr, cond: &searchConditionColumn{colId: colId{name: "column_c"}, pred: notTrue}},
				{bool: boolAnd, cond: &searchConditionColumn{colId: colId{name: "column_d"}, pred: isTrue}},
				{bool: boolOr, cond: &searchConditionColumn{colId: colId{name: "column_e"}, pred: notFalse}},
				{bool: boolOr, cond: &searchConditionColumn{colId: colId{name: "column_f"}, pred: isFalse}},
				{bool: boolAnd, cond: &searchConditionColumn{colId: colId{name: "column_g"}, pred: notUnknown}},
				{bool: boolAnd, cond: &searchConditionColumn{colId: colId{name: "column_h"}, pred: isUnknown}},
				{bool: boolAnd, cond: &searchConditionColumn{colId: colId{name: "column_i"}, pred: isTrue}},
			}},
		},
	}, {
		name: "DeleteAnalysisTestOK11",
		want: &queryStruct{
			name: "DeleteAnalysisTestOK11",
			kind: queryKindDelete,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a"},
				data:  dataType{typeInfo: typeInfo{kind: typeKindStruct}},
			},
			whereBlock: &whereBlock{name: "Where", conds: []*searchCondition{
				{cond: &searchConditionNested{name: "x", conds: []*searchCondition{
					{cond: &searchConditionField{
						name:  "foo",
						typ:   typeInfo{kind: typeKindInt},
						colId: colId{name: "column_foo"},
						pred:  isEQ,
					}},
					{bool: boolAnd, cond: &searchConditionColumn{colId: colId{name: "column_a"}, pred: isNull}},
				}}},
				{bool: boolOr, cond: &searchConditionNested{name: "y", conds: []*searchCondition{
					{cond: &searchConditionColumn{colId: colId{name: "column_b"}, pred: notTrue}},
					{bool: boolOr, cond: &searchConditionField{
						name:  "bar",
						typ:   typeInfo{kind: typeKindString},
						colId: colId{name: "column_bar"},
						pred:  isEQ,
					}},
					{bool: boolAnd, cond: &searchConditionNested{name: "z", conds: []*searchCondition{
						{cond: &searchConditionField{
							name:  "baz",
							typ:   typeInfo{kind: typeKindBool},
							colId: colId{name: "column_baz"},
							pred:  isEQ,
						}},
						{bool: boolAnd, cond: &searchConditionField{
							name:  "quux",
							typ:   typeInfo{kind: typeKindString},
							colId: colId{name: "column_quux"},
							pred:  isEQ,
						}},
						{bool: boolOr, cond: &searchConditionColumn{colId: colId{name: "column_c"}, pred: isTrue}},
					}}},
				}}},
				{bool: boolOr, cond: &searchConditionColumn{colId: colId{name: "column_d"}, pred: notFalse}},
				{bool: boolAnd, cond: &searchConditionColumn{colId: colId{name: "column_e"}, pred: isFalse}},
				{bool: boolAnd, cond: &searchConditionField{
					name:  "foo",
					typ:   typeInfo{kind: typeKindInt},
					colId: colId{name: "column_foo"},
					pred:  isEQ,
				}},
			}},
		},
	}, {
		name: "DeleteAnalysisTestOK12",
		want: &queryStruct{
			name: "DeleteAnalysisTestOK12",
			kind: queryKindDelete,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a"},
				data:  dataType{typeInfo: typeInfo{kind: typeKindStruct}},
			},
			whereBlock: &whereBlock{name: "Where", conds: []*searchCondition{
				{cond: &searchConditionField{name: "a", typ: typeInfo{kind: typeKindInt}, colId: colId{name: "column_a"}, pred: isLT}},
				{bool: boolAnd, cond: &searchConditionField{name: "b", typ: typeInfo{kind: typeKindInt}, colId: colId{name: "column_b"}, pred: isGT}},
				{bool: boolAnd, cond: &searchConditionField{name: "c", typ: typeInfo{kind: typeKindInt}, colId: colId{name: "column_c"}, pred: isLTE}},
				{bool: boolAnd, cond: &searchConditionField{name: "d", typ: typeInfo{kind: typeKindInt}, colId: colId{name: "column_d"}, pred: isGTE}},
				{bool: boolAnd, cond: &searchConditionField{name: "e", typ: typeInfo{kind: typeKindInt}, colId: colId{name: "column_e"}, pred: isEQ}},
				{bool: boolAnd, cond: &searchConditionField{name: "f", typ: typeInfo{kind: typeKindInt}, colId: colId{name: "column_f"}, pred: notEQ}},
				{bool: boolAnd, cond: &searchConditionField{name: "g", typ: typeInfo{kind: typeKindInt}, colId: colId{name: "column_g"}, pred: isEQ}},
			}},
		},
	}, {
		name: "DeleteAnalysisTestOK13",
		want: &queryStruct{
			name: "DeleteAnalysisTestOK13",
			kind: queryKindDelete,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a"},
				data:  dataType{typeInfo: typeInfo{kind: typeKindStruct}},
			},
			whereBlock: &whereBlock{name: "Where", conds: []*searchCondition{
				{cond: &searchConditionColumn{colId: colId{name: "column_a"}, pred: notEQ, colId2: colId{name: "column_b"}}},
				{bool: boolAnd, cond: &searchConditionColumn{colId: colId{qual: "t", name: "column_c"}, pred: isEQ, colId2: colId{qual: "u", name: "column_d"}}},
				{bool: boolAnd, cond: &searchConditionColumn{colId: colId{qual: "t", name: "column_e"}, pred: isGT, literal: "123"}},
				{bool: boolAnd, cond: &searchConditionColumn{colId: colId{qual: "t", name: "column_f"}, pred: isEQ, literal: "'active'"}},
				{bool: boolAnd, cond: &searchConditionColumn{colId: colId{qual: "t", name: "column_g"}, pred: notEQ, literal: "true"}},
			}},
		},
	}, {
		name: "DeleteAnalysisTestOK14",
		want: &queryStruct{
			name: "DeleteAnalysisTestOK14",
			kind: queryKindDelete,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a"},
				data:  dataType{typeInfo: typeInfo{kind: typeKindStruct}},
			},
			whereBlock: &whereBlock{name: "Where", conds: []*searchCondition{
				{cond: &searchConditionBetween{
					name:  "a",
					colId: colId{name: "column_a"},
					pred:  isBetween,
					x:     &fieldDatum{name: "x", typ: typeInfo{kind: typeKindInt}},
					y:     &fieldDatum{name: "y", typ: typeInfo{kind: typeKindInt}},
				}},
				{bool: boolAnd, cond: &searchConditionBetween{
					name:  "b",
					colId: colId{name: "column_b"},
					pred:  isBetweenSym,
					x:     colId{name: "column_x"},
					y:     colId{name: "column_y"},
				}},
				{bool: boolAnd, cond: &searchConditionBetween{
					name:  "c",
					colId: colId{name: "column_c"},
					pred:  notBetweenSym,
					x:     colId{name: "column_z"},
					y:     &fieldDatum{name: "z", typ: typeInfo{kind: typeKindInt}},
				}},
				{bool: boolAnd, cond: &searchConditionBetween{
					name:  "d",
					colId: colId{name: "column_d"},
					pred:  notBetween,
					x:     &fieldDatum{name: "z", typ: typeInfo{kind: typeKindInt}},
					y:     colId{name: "column_z"},
				}},
			}},
		},
	}, {
		name: "DeleteAnalysisTestOK_DistinctFrom",
		want: &queryStruct{
			name: "DeleteAnalysisTestOK_DistinctFrom",
			kind: queryKindDelete,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a"},
				data:  dataType{typeInfo: typeInfo{kind: typeKindStruct}},
			},
			whereBlock: &whereBlock{name: "Where", conds: []*searchCondition{
				{cond: &searchConditionField{
					name:  "a",
					typ:   typeInfo{kind: typeKindInt},
					colId: colId{name: "column_a"},
					pred:  isDistinct,
				}},
				{bool: boolAnd, cond: &searchConditionField{
					name:  "b",
					typ:   typeInfo{kind: typeKindInt},
					colId: colId{name: "column_b"},
					pred:  notDistinct,
				}},
				{bool: boolAnd, cond: &searchConditionColumn{colId: colId{name: "column_c"}, pred: isDistinct, colId2: colId{name: "column_x"}}},
				{bool: boolAnd, cond: &searchConditionColumn{colId: colId{name: "column_d"}, pred: notDistinct, colId2: colId{name: "column_y"}}},
			}},
		},
	}, {
		name: "DeleteAnalysisTestOK_ArrayPredicate",
		want: &queryStruct{
			name: "DeleteAnalysisTestOK_ArrayPredicate",
			kind: queryKindDelete,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a"},
				data:  dataType{typeInfo: typeInfo{kind: typeKindStruct}},
			},
			whereBlock: &whereBlock{name: "Where", conds: []*searchCondition{
				{cond: &searchConditionField{
					name: "a",
					typ: typeInfo{
						kind: typeKindSlice,
						elem: &typeInfo{
							kind: typeKindInt,
						},
					},
					colId: colId{name: "column_a"},
					pred:  isIn,
				}},
				{bool: boolAnd, cond: &searchConditionField{
					name: "b",
					typ: typeInfo{
						kind: typeKindArray,
						elem: &typeInfo{
							kind: typeKindInt,
						},
						arrayLen: 5,
					},
					colId: colId{name: "column_b"},
					pred:  notIn,
				}},
				{bool: boolAnd, cond: &searchConditionField{
					name: "c",
					typ: typeInfo{
						kind: typeKindSlice,
						elem: &typeInfo{
							kind: typeKindInt,
						},
					},
					colId: colId{name: "column_c"},
					pred:  isEQ,
					qua:   quantAny,
				}},
				{bool: boolAnd, cond: &searchConditionField{
					name: "d",
					typ: typeInfo{
						kind: typeKindArray,
						elem: &typeInfo{
							kind: typeKindInt,
						},
						arrayLen: 10,
					},
					colId: colId{name: "column_d"},
					pred:  isGT,
					qua:   quantSome,
				}},
				{bool: boolAnd, cond: &searchConditionField{
					name: "e",
					typ: typeInfo{
						kind: typeKindSlice,
						elem: &typeInfo{
							kind: typeKindInt,
						},
					},
					colId: colId{name: "column_e"},
					pred:  isLTE,
					qua:   quantAll,
				}},
			}},
		},
	}, {
		name: "DeleteAnalysisTestOK_PatternMatching",
		want: &queryStruct{
			name: "DeleteAnalysisTestOK_PatternMatching",
			kind: queryKindDelete,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a"},
				data:  dataType{typeInfo: typeInfo{kind: typeKindStruct}},
			},
			whereBlock: &whereBlock{name: "Where", conds: []*searchCondition{
				{cond: &searchConditionField{
					name:  "a",
					typ:   typeInfo{kind: typeKindString},
					colId: colId{name: "column_a"},
					pred:  isLike,
				}},
				{bool: boolAnd, cond: &searchConditionField{
					name:  "b",
					typ:   typeInfo{kind: typeKindString},
					colId: colId{name: "column_b"},
					pred:  notLike,
				}},
				{bool: boolAnd, cond: &searchConditionField{
					name:  "c",
					typ:   typeInfo{kind: typeKindString},
					colId: colId{name: "column_c"},
					pred:  isSimilar,
				}},
				{bool: boolAnd, cond: &searchConditionField{
					name:  "d",
					typ:   typeInfo{kind: typeKindString},
					colId: colId{name: "column_d"},
					pred:  notSimilar,
				}},
				{bool: boolAnd, cond: &searchConditionField{
					name:  "e",
					typ:   typeInfo{kind: typeKindString},
					colId: colId{name: "column_e"},
					pred:  isMatch,
				}},
				{bool: boolAnd, cond: &searchConditionField{
					name:  "f",
					typ:   typeInfo{kind: typeKindString},
					colId: colId{name: "column_f"},
					pred:  isMatchi,
				}},
				{bool: boolAnd, cond: &searchConditionField{
					name:  "g",
					typ:   typeInfo{kind: typeKindString},
					colId: colId{name: "column_g"},
					pred:  notMatch,
				}},
				{bool: boolAnd, cond: &searchConditionField{
					name:  "h",
					typ:   typeInfo{kind: typeKindString},
					colId: colId{name: "column_h"},
					pred:  notMatchi,
				}},
			}},
		},
	}, {
		name: "DeleteAnalysisTestOK_Using",
		want: &queryStruct{
			name: "DeleteAnalysisTestOK_Using",
			kind: queryKindDelete,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a", alias: "a"},
				data:  dataType{typeInfo: typeInfo{kind: typeKindStruct}},
			},
			joinBlock: &joinBlock{name: "Using", relId: relId{name: "relation_b", alias: "b"}, items: []*joinItem{
				{joinType: joinLeft, relId: relId{name: "relation_c", alias: "c"}, conds: []*searchCondition{{
					cond: &searchConditionColumn{
						colId:  colId{qual: "c", name: "b_id"},
						colId2: colId{qual: "b", name: "id"},
						pred:   isEQ,
					}}}},
				{joinType: joinRight, relId: relId{name: "relation_d", alias: "d"}, conds: []*searchCondition{{
					cond: &searchConditionColumn{
						colId:  colId{qual: "d", name: "c_id"},
						colId2: colId{qual: "c", name: "id"},
						pred:   isEQ,
					},
				}, {
					bool: boolOr,
					cond: &searchConditionColumn{
						colId:  colId{qual: "d", name: "num"},
						colId2: colId{qual: "b", name: "num"},
						pred:   isGT,
					},
				}}},
				{joinType: joinFull, relId: relId{name: "relation_e", alias: "e"}, conds: []*searchCondition{{
					cond: &searchConditionColumn{
						colId:  colId{qual: "e", name: "d_id"},
						colId2: colId{qual: "d", name: "id"},
						pred:   isEQ,
					},
				}, {
					bool: boolAnd,
					cond: &searchConditionColumn{
						colId: colId{qual: "e", name: "is_foo"},
						pred:  isFalse,
					},
				}}},
				{joinType: joinCross, relId: relId{name: "relation_f", alias: "f"}},
			}},
			whereBlock: &whereBlock{name: "Where", conds: []*searchCondition{
				{cond: &searchConditionColumn{
					colId:  colId{qual: "a", name: "id"},
					pred:   isEQ,
					colId2: colId{qual: "d", name: "a_id"},
				}},
			}},
		},
	}, {
		name: "UpdateAnalysisTestOK_From",
		want: &queryStruct{
			name: "UpdateAnalysisTestOK_From",
			kind: queryKindUpdate,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a", alias: "a"},
				data:  dataType{typeInfo: typeInfo{kind: typeKindStruct}},
			},
			joinBlock: &joinBlock{name: "From", relId: relId{name: "relation_b", alias: "b"}, items: []*joinItem{
				{joinType: joinLeft, relId: relId{name: "relation_c", alias: "c"}, conds: []*searchCondition{{
					cond: &searchConditionColumn{
						colId:  colId{qual: "c", name: "b_id"},
						colId2: colId{qual: "b", name: "id"},
						pred:   isEQ,
					},
				}}},
				{joinType: joinRight, relId: relId{name: "relation_d", alias: "d"}, conds: []*searchCondition{{
					cond: &searchConditionColumn{
						colId:  colId{qual: "d", name: "c_id"},
						colId2: colId{qual: "c", name: "id"},
						pred:   isEQ,
					},
				}, {
					bool: boolOr,
					cond: &searchConditionColumn{
						colId:  colId{qual: "d", name: "num"},
						colId2: colId{qual: "b", name: "num"},
						pred:   isGT,
					},
				}}},
				{joinType: joinFull, relId: relId{name: "relation_e", alias: "e"}, conds: []*searchCondition{{
					cond: &searchConditionColumn{
						colId:  colId{qual: "e", name: "d_id"},
						colId2: colId{qual: "d", name: "id"},
						pred:   isEQ,
					},
				}, {
					bool: boolAnd,
					cond: &searchConditionColumn{
						colId: colId{qual: "e", name: "is_foo"},
						pred:  isFalse,
					},
				}}},
				{joinType: joinCross, relId: relId{name: "relation_f", alias: "f"}},
			}},
			whereBlock: &whereBlock{name: "Where", conds: []*searchCondition{
				{cond: &searchConditionColumn{
					colId:  colId{qual: "a", name: "id"},
					colId2: colId{qual: "d", name: "a_id"},
					pred:   isEQ,
				}},
			}},
		},
	}, {
		name: "SelectAnalysisTestOK_Join",
		want: &queryStruct{
			name: "SelectAnalysisTestOK_Join",
			kind: queryKindSelect,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a", alias: "a"},
				data:  dataType{typeInfo: typeInfo{kind: typeKindStruct}},
			},
			joinBlock: &joinBlock{name: "Join", items: []*joinItem{
				{joinType: joinLeft, relId: relId{name: "relation_b", alias: "b"}, conds: []*searchCondition{{
					cond: &searchConditionColumn{
						colId:  colId{qual: "b", name: "a_id"},
						colId2: colId{qual: "a", name: "id"},
						pred:   isEQ,
					},
				}}},
				{joinType: joinLeft, relId: relId{name: "relation_c", alias: "c"}, conds: []*searchCondition{{
					cond: &searchConditionColumn{
						colId:  colId{qual: "c", name: "b_id"},
						colId2: colId{qual: "b", name: "id"},
						pred:   isEQ,
					},
				}}},
				{joinType: joinRight, relId: relId{name: "relation_d", alias: "d"}, conds: []*searchCondition{{
					cond: &searchConditionColumn{
						colId:  colId{qual: "d", name: "c_id"},
						colId2: colId{qual: "c", name: "id"},
						pred:   isEQ,
					},
				}, {
					bool: boolOr,
					cond: &searchConditionColumn{
						colId:  colId{qual: "d", name: "num"},
						colId2: colId{qual: "b", name: "num"},
						pred:   isGT,
					},
				}}},
				{joinType: joinFull, relId: relId{name: "relation_e", alias: "e"}, conds: []*searchCondition{{
					cond: &searchConditionColumn{
						colId:  colId{qual: "e", name: "d_id"},
						colId2: colId{qual: "d", name: "id"},
						pred:   isEQ,
					},
				}, {
					bool: boolAnd,
					cond: &searchConditionColumn{
						colId: colId{qual: "e", name: "is_foo"},
						pred:  isFalse,
					},
				}}},
				{joinType: joinCross, relId: relId{name: "relation_f", alias: "f"}},
			}},
			whereBlock: &whereBlock{name: "Where", conds: []*searchCondition{
				{cond: &searchConditionColumn{
					colId:  colId{qual: "a", name: "id"},
					pred:   isEQ,
					colId2: colId{qual: "d", name: "a_id"},
				}},
			}},
		},
	}, {
		name: "UpdateAnalysisTestOK_All",
		want: &queryStruct{
			name: "UpdateAnalysisTestOK_All",
			kind: queryKindUpdate,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a", alias: "a"},
				data:  dataType{typeInfo: typeInfo{kind: typeKindStruct}},
			},
			all: true,
		},
	}, {
		name: "DeleteAnalysisTestOK_All",
		want: &queryStruct{
			name: "DeleteAnalysisTestOK_All",
			kind: queryKindDelete,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a", alias: "a"},
				data:  dataType{typeInfo: typeInfo{kind: typeKindStruct}},
			},
			all: true,
		},
	}, {
		name: "DeleteAnalysisTestOK_Return",
		want: &queryStruct{
			name: "DeleteAnalysisTestOK_Return",
			kind: queryKindDelete,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a", alias: "a"},
				data:  dummyrecord,
			},
			returnList: &colIdList{all: true},
		},
	}, {
		name: "InsertAnalysisTestOK_Return",
		want: &queryStruct{
			name: "InsertAnalysisTestOK_Return",
			kind: queryKindInsert,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a", alias: "a"},
				data:  dummyrecord,
			},
			returnList: &colIdList{items: []colId{
				{qual: "a", name: "foo"},
				{qual: "a", name: "bar"},
				{qual: "a", name: "baz"},
			}},
		},
	}, {
		name: "UpdateAnalysisTestOK_Return",
		want: &queryStruct{
			name: "UpdateAnalysisTestOK_Return",
			kind: queryKindUpdate,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a", alias: "a"},
				data:  dummyrecord,
			},
			returnList: &colIdList{items: []colId{
				{qual: "a", name: "foo"},
				{qual: "a", name: "bar"},
				{qual: "a", name: "baz"},
			}},
		},
	}, {
		name: "InsertAnalysisTestOK_Default",
		want: &queryStruct{
			name: "InsertAnalysisTestOK_Default",
			kind: queryKindInsert,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a", alias: "a"},
				data:  dataType{typeInfo: typeInfo{kind: typeKindStruct}},
			},
			defaultList: &colIdList{all: true},
		},
	}, {
		name: "UpdateAnalysisTestOK_Default",
		want: &queryStruct{
			name: "UpdateAnalysisTestOK_Default",
			kind: queryKindUpdate,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a", alias: "a"},
				data:  dataType{typeInfo: typeInfo{kind: typeKindStruct}},
			},
			defaultList: &colIdList{items: []colId{
				{qual: "a", name: "foo"},
				{qual: "a", name: "bar"},
				{qual: "a", name: "baz"},
			}},
		},
	}, {
		name: "InsertAnalysisTestOK_Force",
		want: &queryStruct{
			name: "InsertAnalysisTestOK_Force",
			kind: queryKindInsert,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a", alias: "a"},
				data:  dataType{typeInfo: typeInfo{kind: typeKindStruct}},
			},
			forceList: &colIdList{all: true},
		},
	}, {
		name: "UpdateAnalysisTestOK_Force",
		want: &queryStruct{
			name: "UpdateAnalysisTestOK_Force",
			kind: queryKindUpdate,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a", alias: "a"},
				data:  dataType{typeInfo: typeInfo{kind: typeKindStruct}},
			},
			forceList: &colIdList{items: []colId{
				{qual: "a", name: "foo"},
				{qual: "a", name: "bar"},
				{qual: "a", name: "baz"},
			}},
		},
	}, {
		name: "SelectAnalysisTestOK_ErrorHandler",
		want: &queryStruct{
			name: "SelectAnalysisTestOK_ErrorHandler",
			kind: queryKindSelect,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a", alias: "a"},
				data:  dataType{typeInfo: typeInfo{kind: typeKindStruct}},
			},
			errorHandlerField: &errorHandlerField{name: "eh"},
		},
	}, {
		name: "InsertAnalysisTestOK_ErrorHandler",
		want: &queryStruct{
			name: "InsertAnalysisTestOK_ErrorHandler",
			kind: queryKindInsert,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a", alias: "a"},
				data:  dataType{typeInfo: typeInfo{kind: typeKindStruct}},
			},
			errorHandlerField: &errorHandlerField{name: "myerrorhandler"},
		},
	}, {
		name: "SelectAnalysisTestOK_ErrorInfoHandler",
		want: &queryStruct{
			name: "SelectAnalysisTestOK_ErrorInfoHandler",
			kind: queryKindSelect,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a", alias: "a"},
				data:  dataType{typeInfo: typeInfo{kind: typeKindStruct}},
			},
			errorHandlerField: &errorHandlerField{name: "eh", isInfo: true},
		},
	}, {
		name: "InsertAnalysisTestOK_ErrorInfoHandler",
		want: &queryStruct{
			name: "InsertAnalysisTestOK_ErrorInfoHandler",
			kind: queryKindInsert,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a", alias: "a"},
				data:  dataType{typeInfo: typeInfo{kind: typeKindStruct}},
			},
			errorHandlerField: &errorHandlerField{name: "myerrorinfohandler", isInfo: true},
		},
	}, {
		name: "SelectAnalysisTestOK_Count",
		want: &queryStruct{
			name: "SelectAnalysisTestOK_Count",
			kind: queryKindSelectCount,
			dataField: &dataField{
				name:  "Count",
				relId: relId{name: "relation_a", alias: "a"},
			},
		},
	}, {
		name: "SelectAnalysisTestOK_Exists",
		want: &queryStruct{
			name: "SelectAnalysisTestOK_Exists",
			kind: queryKindSelectExists,
			dataField: &dataField{
				name:  "Exists",
				relId: relId{name: "relation_a", alias: "a"},
			},
		},
	}, {
		name: "SelectAnalysisTestOK_NotExists",
		want: &queryStruct{
			name: "SelectAnalysisTestOK_NotExists",
			kind: queryKindSelectNotExists,
			dataField: &dataField{
				name:  "NotExists",
				relId: relId{name: "relation_a", alias: "a"},
			},
		},
	}, {
		name: "DeleteAnalysisTestOK_Relation",
		want: &queryStruct{
			name: "DeleteAnalysisTestOK_Relation",
			kind: queryKindDelete,
			dataField: &dataField{
				name:        "_",
				relId:       relId{name: "relation_a", alias: "a"},
				isDirective: true,
			},
		},
	}, {
		name: "SelectAnalysisTestOK_LimitDirective",
		want: &queryStruct{
			name:       "SelectAnalysisTestOK_LimitDirective",
			kind:       queryKindSelect,
			dataField:  reldummyslice,
			limitField: &limitField{value: 25},
		},
	}, {
		name: "SelectAnalysisTestOK_LimitField",
		want: &queryStruct{
			name:       "SelectAnalysisTestOK_LimitField",
			kind:       queryKindSelect,
			dataField:  reldummyslice,
			limitField: &limitField{name: "Limit", value: 10},
		},
	}, {
		name: "SelectAnalysisTestOK_OffsetDirective",
		want: &queryStruct{
			name:        "SelectAnalysisTestOK_OffsetDirective",
			kind:        queryKindSelect,
			dataField:   reldummyslice,
			offsetField: &offsetField{value: 25},
		},
	}, {
		name: "SelectAnalysisTestOK_OffsetField",
		want: &queryStruct{
			name:        "SelectAnalysisTestOK_OffsetField",
			kind:        queryKindSelect,
			dataField:   reldummyslice,
			offsetField: &offsetField{name: "Offset", value: 10},
		},
	}, {
		name: "SelectAnalysisTestOK_OrderByDirective",
		want: &queryStruct{
			name:      "SelectAnalysisTestOK_OrderByDirective",
			kind:      queryKindSelect,
			dataField: reldummyslice,
			orderByList: &orderByList{items: []*orderByItem{
				{colId: colId{qual: "a", name: "foo"}, direction: orderAsc, nulls: nullsFirst},
				{colId: colId{qual: "a", name: "bar"}, direction: orderDesc, nulls: nullsFirst},
				{colId: colId{qual: "a", name: "baz"}, direction: orderDesc, nulls: 0},
				{colId: colId{qual: "a", name: "quux"}, direction: orderAsc, nulls: nullsLast},
			}},
		},
	}, {
		name: "InsertAnalysisTestOK_OverrideDirective",
		want: &queryStruct{
			name:           "InsertAnalysisTestOK_OverrideDirective",
			kind:           queryKindInsert,
			dataField:      reldummyslice,
			overridingKind: overridingSystem,
		},
	}, {
		name: "InsertAnalysisTestOK_OnConflict",
		want: &queryStruct{
			name:            "InsertAnalysisTestOK_OnConflict",
			kind:            queryKindInsert,
			dataField:       reldummyslice,
			onConflictBlock: &onConflictBlock{ignore: true},
		},
	}, {
		name: "InsertAnalysisTestOK_OnConflictColumn",
		want: &queryStruct{
			name:      "InsertAnalysisTestOK_OnConflictColumn",
			kind:      queryKindInsert,
			dataField: reldummyslice,
			onConflictBlock: &onConflictBlock{
				column: []colId{{qual: "a", name: "id"}},
				ignore: true,
			},
		},
	}, {
		name: "InsertAnalysisTestOK_OnConflictConstraint",
		want: &queryStruct{
			name:      "InsertAnalysisTestOK_OnConflictConstraint",
			kind:      queryKindInsert,
			dataField: reldummyslice,
			onConflictBlock: &onConflictBlock{
				constraint: "relation_constraint_xyz",
				update: &colIdList{items: []colId{
					{qual: "a", name: "foo"},
					{qual: "a", name: "bar"},
					{qual: "a", name: "baz"},
				}},
			},
		},
	}, {
		name: "InsertAnalysisTestOK_OnConflictIndex",
		want: &queryStruct{
			name:      "InsertAnalysisTestOK_OnConflictIndex",
			kind:      queryKindInsert,
			dataField: reldummyslice,
			onConflictBlock: &onConflictBlock{
				index:  "relation_index_xyz",
				update: &colIdList{all: true},
			},
		},
	}, {
		name: "DeleteAnalysisTestOK_ResultField",
		want: &queryStruct{
			name: "DeleteAnalysisTestOK_ResultField",
			kind: queryKindDelete,
			dataField: &dataField{
				name:        "_",
				relId:       relId{name: "relation_a", alias: "a"},
				isDirective: true,
			},
			whereBlock: &whereBlock{name: "Where", conds: []*searchCondition{
				{cond: &searchConditionColumn{colId: colId{qual: "a", name: "is_inactive"}, pred: isTrue}},
			}},
			resultField: &resultField{
				name: "Result",
				data: reldummyslice.data,
			},
		},
	}, {
		name: "DeleteAnalysisTestOK_RowsAffected",
		want: &queryStruct{
			name: "DeleteAnalysisTestOK_RowsAffected",
			kind: queryKindDelete,
			dataField: &dataField{
				name:        "_",
				relId:       relId{name: "relation_a", alias: "a"},
				isDirective: true,
			},
			whereBlock: &whereBlock{name: "Where", conds: []*searchCondition{
				{cond: &searchConditionColumn{colId: colId{qual: "a", name: "is_inactive"}, pred: isTrue}},
			}},
			rowsAffectedField: &rowsAffectedField{
				name: "RowsAffected",
				kind: typeKindInt,
			},
		},
	}, {
		name: "SelectAnalysisTestOK_FilterField",
		want: &queryStruct{
			name:        "SelectAnalysisTestOK_FilterField",
			kind:        queryKindSelect,
			dataField:   reldummyslice,
			filterField: "Filter",
		},
	}, {
		name: "SelectAnalysisTestOK_FieldTypesBasic",
		want: &queryStruct{
			name: "SelectAnalysisTestOK_FieldTypesBasic",
			kind: queryKindSelect,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a", alias: "a"},
				data: dataType{
					typeInfo: typeInfo{kind: typeKindStruct},
					fields: []*fieldInfo{{
						name: "f1", typ: typeInfo{kind: typeKindBool},
						colId: colId{name: "c1"},
						tag:   tagutil.Tag{"sql": {"c1"}},
					}, {
						name: "f2", typ: typeInfo{kind: typeKindUint8, isByte: true},
						colId: colId{name: "c2"},
						tag:   tagutil.Tag{"sql": {"c2"}},
					}, {
						name: "f3", typ: typeInfo{kind: typeKindInt32, isRune: true},
						colId: colId{name: "c3"},
						tag:   tagutil.Tag{"sql": {"c3"}},
					}, {
						name: "f4", typ: typeInfo{kind: typeKindInt8},
						colId: colId{name: "c4"},
						tag:   tagutil.Tag{"sql": {"c4"}},
					}, {
						name: "f5", typ: typeInfo{kind: typeKindInt16},
						colId: colId{name: "c5"},
						tag:   tagutil.Tag{"sql": {"c5"}},
					}, {
						name: "f6", typ: typeInfo{kind: typeKindInt32},
						colId: colId{name: "c6"},
						tag:   tagutil.Tag{"sql": {"c6"}},
					}, {
						name: "f7", typ: typeInfo{kind: typeKindInt64},
						colId: colId{name: "c7"},
						tag:   tagutil.Tag{"sql": {"c7"}},
					}, {
						name: "f8", typ: typeInfo{kind: typeKindInt},
						colId: colId{name: "c8"},
						tag:   tagutil.Tag{"sql": {"c8"}},
					}, {
						name: "f9", typ: typeInfo{kind: typeKindUint8},
						colId: colId{name: "c9"},
						tag:   tagutil.Tag{"sql": {"c9"}},
					}, {
						name: "f10", typ: typeInfo{kind: typeKindUint16},
						colId: colId{name: "c10"},
						tag:   tagutil.Tag{"sql": {"c10"}},
					}, {
						name: "f11", typ: typeInfo{kind: typeKindUint32},
						colId: colId{name: "c11"},
						tag:   tagutil.Tag{"sql": {"c11"}},
					}, {
						name: "f12", typ: typeInfo{kind: typeKindUint64},
						colId: colId{name: "c12"},
						tag:   tagutil.Tag{"sql": {"c12"}},
					}, {
						name: "f13", typ: typeInfo{kind: typeKindUint},
						colId: colId{name: "c13"},
						tag:   tagutil.Tag{"sql": {"c13"}},
					}, {
						name: "f14", typ: typeInfo{kind: typeKindUintptr},
						colId: colId{name: "c14"},
						tag:   tagutil.Tag{"sql": {"c14"}},
					}, {
						name: "f15", typ: typeInfo{kind: typeKindFloat32},
						colId: colId{name: "c15"},
						tag:   tagutil.Tag{"sql": {"c15"}},
					}, {
						name: "f16", typ: typeInfo{kind: typeKindFloat64},
						colId: colId{name: "c16"},
						tag:   tagutil.Tag{"sql": {"c16"}},
					}, {
						name: "f17", typ: typeInfo{kind: typeKindComplex64},
						colId: colId{name: "c17"},
						tag:   tagutil.Tag{"sql": {"c17"}},
					}, {
						name: "f18", typ: typeInfo{kind: typeKindComplex128},
						colId: colId{name: "c18"},
						tag:   tagutil.Tag{"sql": {"c18"}},
					}, {
						name: "f19", typ: typeInfo{kind: typeKindString},
						colId: colId{name: "c19"},
						tag:   tagutil.Tag{"sql": {"c19"}},
					}},
				},
			},
		},
	}, {
		name: "SelectAnalysisTestOK_FieldTypesSlices",
		want: &queryStruct{
			name: "SelectAnalysisTestOK_FieldTypesSlices",
			kind: queryKindSelect,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a", alias: "a"},
				data: dataType{
					typeInfo: typeInfo{kind: typeKindStruct},
					fields: []*fieldInfo{{
						name: "f1", typ: typeInfo{
							kind: typeKindSlice,
							elem: &typeInfo{kind: typeKindBool},
						},
						colId: colId{name: "c1"},
						tag:   tagutil.Tag{"sql": {"c1"}},
					}, {
						name: "f2", typ: typeInfo{
							kind: typeKindSlice,
							elem: &typeInfo{kind: typeKindUint8, isByte: true},
						},
						colId: colId{name: "c2"},
						tag:   tagutil.Tag{"sql": {"c2"}},
					}, {
						name: "f3", typ: typeInfo{
							kind: typeKindSlice,
							elem: &typeInfo{kind: typeKindInt32, isRune: true},
						},
						colId: colId{name: "c3"},
						tag:   tagutil.Tag{"sql": {"c3"}},
					}, {
						name: "f4", typ: typeInfo{
							name:       "HardwareAddr",
							kind:       typeKindSlice,
							pkgPath:    "net",
							pkgName:    "net",
							pkgLocal:   "net",
							isImported: true,
							elem:       &typeInfo{kind: typeKindUint8, isByte: true},
						},
						colId: colId{name: "c4"},
						tag:   tagutil.Tag{"sql": {"c4"}},
					}, {
						name: "f5", typ: typeInfo{
							name:              "RawMessage",
							kind:              typeKindSlice,
							pkgPath:           "encoding/json",
							pkgName:           "json",
							pkgLocal:          "json",
							isImported:        true,
							isJSONMarshaler:   true,
							isJSONUnmarshaler: true,
							elem:              &typeInfo{kind: typeKindUint8, isByte: true},
						},
						colId: colId{name: "c5"},
						tag:   tagutil.Tag{"sql": {"c5"}},
					}, {
						name: "f6", typ: typeInfo{
							kind: typeKindSlice,
							elem: &typeInfo{
								name:            "Marshaler",
								kind:            typeKindInterface,
								pkgPath:         "encoding/json",
								pkgName:         "json",
								pkgLocal:        "json",
								isImported:      true,
								isJSONMarshaler: true,
							},
						},
						colId: colId{name: "c6"},
						tag:   tagutil.Tag{"sql": {"c6"}},
					}, {
						name: "f7", typ: typeInfo{
							kind: typeKindSlice,
							elem: &typeInfo{
								name:              "RawMessage",
								kind:              typeKindSlice,
								pkgPath:           "encoding/json",
								pkgName:           "json",
								pkgLocal:          "json",
								isImported:        true,
								isJSONMarshaler:   true,
								isJSONUnmarshaler: true,
								elem:              &typeInfo{kind: typeKindUint8, isByte: true},
							},
						},
						colId: colId{name: "c7"},
						tag:   tagutil.Tag{"sql": {"c7"}},
					}, {
						name: "f8", typ: typeInfo{
							kind: typeKindSlice,
							elem: &typeInfo{
								kind: typeKindSlice,
								elem: &typeInfo{kind: typeKindUint8, isByte: true},
							},
						},
						colId: colId{name: "c8"},
						tag:   tagutil.Tag{"sql": {"c8"}},
					}, {
						name: "f9", typ: typeInfo{
							kind: typeKindSlice,
							elem: &typeInfo{
								kind:     typeKindArray,
								arrayLen: 2,
								elem: &typeInfo{
									kind:     typeKindArray,
									arrayLen: 2,
									elem:     &typeInfo{kind: typeKindFloat64},
								},
							},
						},
						colId: colId{name: "c9"},
						tag:   tagutil.Tag{"sql": {"c9"}},
					}, {
						name: "f10", typ: typeInfo{
							kind: typeKindSlice,
							elem: &typeInfo{
								kind: typeKindSlice,
								elem: &typeInfo{
									kind:     typeKindArray,
									arrayLen: 2,
									elem:     &typeInfo{kind: typeKindFloat64},
								},
							},
						},
						colId: colId{name: "c10"},
						tag:   tagutil.Tag{"sql": {"c10"}},
					}, {
						name: "f11", typ: typeInfo{
							kind: typeKindMap,
							key:  &typeInfo{kind: typeKindString},
							elem: &typeInfo{
								name:       "NullString",
								kind:       typeKindStruct,
								pkgPath:    "database/sql",
								pkgName:    "sql",
								pkgLocal:   "sql",
								isImported: true,
								isScanner:  true,
								isValuer:   true,
							},
						},
						colId: colId{name: "c11"},
						tag:   tagutil.Tag{"sql": {"c11"}},
					}, {
						name: "f12", typ: typeInfo{
							kind: typeKindSlice,
							elem: &typeInfo{
								kind: typeKindMap,
								key:  &typeInfo{kind: typeKindString},
								elem: &typeInfo{
									kind: typeKindPtr,
									elem: &typeInfo{kind: typeKindString},
								},
							},
						},
						colId: colId{name: "c12"},
						tag:   tagutil.Tag{"sql": {"c12"}},
					}, {
						name: "f13", typ: typeInfo{
							kind: typeKindSlice,
							elem: &typeInfo{
								kind:     typeKindArray,
								arrayLen: 2,
								elem: &typeInfo{
									kind: typeKindPtr,
									elem: &typeInfo{
										name:              "Int",
										kind:              typeKindStruct,
										pkgPath:           "math/big",
										pkgName:           "big",
										pkgLocal:          "big",
										isImported:        true,
										isJSONMarshaler:   true,
										isJSONUnmarshaler: true,
									},
								},
							},
						},
						colId: colId{name: "c13"},
						tag:   tagutil.Tag{"sql": {"c13"}},
					}},
				},
			},
		},
	}, {
		name: "SelectAnalysisTestOK_FieldTypesInterfaces",
		want: &queryStruct{
			name: "SelectAnalysisTestOK_FieldTypesInterfaces",
			kind: queryKindSelect,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a", alias: "a"},
				data: dataType{
					typeInfo: typeInfo{kind: typeKindStruct},
					fields: []*fieldInfo{{
						name: "f1", typ: typeInfo{
							name:            "Marshaler",
							kind:            typeKindInterface,
							pkgPath:         "encoding/json",
							pkgName:         "json",
							pkgLocal:        "json",
							isImported:      true,
							isJSONMarshaler: true,
						},
						colId: colId{name: "c1"},
						tag:   tagutil.Tag{"sql": {"c1"}},
					}, {
						name: "f2", typ: typeInfo{
							name:              "Unmarshaler",
							kind:              typeKindInterface,
							pkgPath:           "encoding/json",
							pkgName:           "json",
							pkgLocal:          "json",
							isImported:        true,
							isJSONUnmarshaler: true,
						},
						colId: colId{name: "c2"},
						tag:   tagutil.Tag{"sql": {"c2"}},
					}, {
						name: "f3", typ: typeInfo{
							kind: typeKindInterface,
						},
						colId: colId{name: "c3"},
						tag:   tagutil.Tag{"sql": {"c3"}},
					}},
				},
			},
		},
	}, {
		name: "SelectAnalysisTestOK_FieldTypesEmptyInterfaces",
		want: &queryStruct{
			name: "SelectAnalysisTestOK_FieldTypesEmptyInterfaces",
			kind: queryKindSelect,
			dataField: &dataField{
				name:  "Rel",
				relId: relId{name: "relation_a", alias: "a"},
				data: dataType{
					typeInfo: typeInfo{kind: typeKindStruct},
					fields: []*fieldInfo{{
						name: "f1", typ: typeInfo{
							kind:             typeKindInterface,
							isEmptyInterface: true,
						},
						colId: colId{name: "c1"},
						tag:   tagutil.Tag{"sql": {"c1"}},
					}, {
						name: "f2", typ: typeInfo{
							kind: typeKindPtr,
							elem: &typeInfo{
								kind:             typeKindInterface,
								isEmptyInterface: true,
							},
						},
						colId: colId{name: "c2"},
						tag:   tagutil.Tag{"sql": {"c2"}},
					}, {
						name: "f3", typ: typeInfo{
							name:             "donothing",
							kind:             typeKindInterface,
							pkgPath:          "path/to/test",
							pkgName:          "testdata",
							pkgLocal:         "testdata",
							isEmptyInterface: true,
						},
						colId: colId{name: "c3"},
						tag:   tagutil.Tag{"sql": {"c3"}},
					}, {
						name: "f4", typ: typeInfo{
							kind: typeKindPtr,
							elem: &typeInfo{
								name:             "donothing",
								kind:             typeKindInterface,
								pkgPath:          "path/to/test",
								pkgName:          "testdata",
								pkgLocal:         "testdata",
								isEmptyInterface: true,
							},
						},
						colId: colId{name: "c4"},
						tag:   tagutil.Tag{"sql": {"c4"}},
					}},
				},
			},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *queryStruct
			a, err := runAnalysis(tt.name, t)
			if a != nil && a.info != nil {
				got = a.info.query
			}

			if e := compare.Compare(err, tt.err); e != nil {
				t.Errorf("%v - %#v %v", e, err, err)
			}
			if e := compare.Compare(got, tt.want); e != nil {
				t.Error(e)
			}
		})
	}
}

func TestAnalysis_filterStruct(t *testing.T) {
	dummyrecord := dataType{
		typeInfo: typeInfo{
			name:     "T",
			kind:     typeKindStruct,
			pkgPath:  "path/to/test",
			pkgName:  "testdata",
			pkgLocal: "testdata",
		},
		fields: []*fieldInfo{{
			typ:        typeInfo{kind: typeKindString},
			name:       "F",
			isExported: true,
			tag:        tagutil.Tag{"sql": {"f"}},
			colId:      colId{name: "f"},
		}},
	}

	tests := []struct {
		name string
		want *filterStruct
		err  error
	}{{
		name: "FilterAnalysisTestBAD_IllegalReturnDirective",
		err: analysisError{
			errorCode:  errIllegalField,
			fieldType:  "github.com/frk/gosql.Return",
			fieldName:  "_",
			targetName: "FilterAnalysisTestBAD_IllegalReturnDirective",
			pkgPath:    "path/to/test",
			fileLine:   104,
			fileName:   "../../testdata/analysis_bad.go",
		},
	}, {
		name: "FilterAnalysisTestBAD_BadTextSearchDirectiveColId",
		err: analysisError{
			errorCode:  errBadColIdTagValue,
			pkgPath:    "path/to/test",
			targetName: "FilterAnalysisTestBAD_BadTextSearchDirectiveColId",
			blockName:  "",
			fieldType:  "github.com/frk/gosql.TextSearch",
			fieldName:  "_",
			tagValue:   "123",
			fileName:   "../../testdata/analysis_bad.go",
			fileLine:   723,
		},
	}, {
		name: "FilterAnalysisTestOK_TextSearchDirective",
		want: &filterStruct{
			name: "FilterAnalysisTestOK_TextSearchDirective",
			dataField: &dataField{
				name:  "_",
				relId: relId{name: "relation_a", alias: "a"},
				data:  dummyrecord,
			},
			textSearchColId: &colId{qual: "a", name: "ts_document"},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *filterStruct
			a, err := runAnalysis(tt.name, t)
			if a != nil && a.info != nil {
				got = a.info.filter
			}

			if e := compare.Compare(err, tt.err); e != nil {
				t.Errorf("%v - %#v %v", e, err, err)
			}
			if e := compare.Compare(got, tt.want); e != nil {
				t.Error(e)
			}
		})
	}
}

func TestTypeinfo_string(t *testing.T) {
	tests := []struct {
		name string
		want goTypeId
	}{
		{"f01", goTypeBool},
		{"f02", goTypeBool},
		{"f03", goTypeBoolSlice},
		{"f04", goTypeString},
		{"f05", goTypeString},
		{"f06", goTypeStringSlice},
		{"f07", goTypeStringSliceSlice},
		{"f08", goTypeStringMap},
		{"f09", goTypeStringPtrMap},
		{"f10", goTypeStringMapSlice},
		{"f11", goTypeStringPtrMapSlice},
		{"f12", goTypeByte},
		{"f13", goTypeByte},
		{"f14", goTypeByteSlice},
		{"f15", goTypeByteSliceSlice},
		{"f16", goTypeByteArray16},
		{"f17", goTypeByteArray16Slice},
		{"f18", goTypeRune},
		{"f19", goTypeRune},
		{"f20", goTypeRuneSlice},
		{"f21", goTypeRuneSliceSlice},
		{"f22", goTypeInt8},
		{"f23", goTypeInt8},
		{"f24", goTypeInt8Slice},
		{"f25", goTypeInt8SliceSlice},
		{"f26", goTypeInt16},
		{"f27", goTypeInt16},
		{"f28", goTypeInt16Slice},
		{"f29", goTypeInt16SliceSlice},
		{"f30", goTypeInt32},
		{"f31", goTypeInt32},
		{"f32", goTypeInt32Slice},
		{"f33", goTypeInt32Array2},
		{"f34", goTypeInt32Array2Slice},
		{"f35", goTypeInt64},
		{"f36", goTypeInt64},
		{"f37", goTypeInt64Slice},
		{"f38", goTypeInt64Array2},
		{"f39", goTypeInt64Array2Slice},
		{"f40", goTypeFloat32},
		{"f41", goTypeFloat32},
		{"f42", goTypeFloat32Slice},
		{"f43", goTypeFloat64},
		{"f44", goTypeFloat64},
		{"f45", goTypeFloat64Slice},
		{"f46", goTypeFloat64Array2},
		{"f47", goTypeFloat64Array2Slice},
		{"f48", goTypeFloat64Array2SliceSlice},
		{"f49", goTypeFloat64Array2Array2},
		{"f50", goTypeFloat64Array2Array2Slice},
		{"f51", goTypeFloat64Array3},
		{"f52", goTypeFloat64Array3Slice},
		{"f53", goTypeIPNet},
		{"f54", "[]*net.IPNet"},
		{"f55", goTypeTime},
		{"f56", goTypeTime},
		{"f57", goTypeTimeSlice},
		{"f58", "[]*time.Time"},
		{"f59", goTypeTimeArray2},
		{"f60", goTypeTimeArray2Slice},
		{"f61", goTypeHardwareAddr},
		{"f62", goTypeHardwareAddrSlice},
		{"f63", goTypeBigInt},
		{"f64", goTypeBigInt},
		{"f65", goTypeBigIntSlice},
		{"f66", "[]*big.Int"},
		{"f67", goTypeBigIntArray2},
		{"f68", "[2]*big.Int"},
		{"f69", "[][2]*big.Int"},
		{"f70", goTypeNullStringMap},
		{"f71", goTypeNullStringMapSlice},
		{"f72", "json.RawMessage"},
		{"f73", "[]json.RawMessage"},
	}

	a, err := runAnalysis("SelectAnalysisTestOK_typeinfo_string", t)
	if err != nil {
		t.Error(err)
	}
	fields := a.info.query.dataField.data.fields
	for i := 0; i < len(fields); i++ {
		ff := fields[i]
		tt := tests[i]

		//got := goTypeId(ff.typ.string(true))
		got := ff.typ.goTypeId(false, false, true)
		if ff.name != tt.name || got != tt.want {
			t.Errorf("got %s::%s, want %s::%s", ff.name, got, tt.name, tt.want)
		}
	}
}
