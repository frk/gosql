package analysis

import (
	"fmt"
	"testing"

	"github.com/frk/compare"
	"github.com/frk/gosql/internal/x/testutil"
	"github.com/frk/tagutil"
)

var _ = fmt.Println

func init() {
	compare.DefaultConfig.ObserveFieldTag = "cmp"
}

var tdata = testutil.ParseTestdata("../testdata")

func testRunAnalysis(name string, t *testing.T) (TargetStruct, error) {
	named := testutil.FindNamedType(name, tdata)
	if named == nil {
		// Stop the test if no type with the given name was found.
		t.Fatal(name, " not found")
		return nil, nil
	}

	ts, err := Run(tdata.Fset, named, &Info{})
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func TestAnalysis_queryStruct(t *testing.T) {
	// for reuse, analyzed common.User TypeInfo
	commonUserTypeinfo := TypeInfo{
		Name:       "User",
		Kind:       TypeKindStruct,
		PkgPath:    "github.com/frk/gosql/internal/testdata/common",
		PkgName:    "common",
		PkgLocal:   "common",
		IsImported: true,
	}

	commonUserFields := []*FieldInfo{{
		Name:       "Id",
		Type:       TypeInfo{Kind: TypeKindInt},
		IsExported: true,
		ColIdent:   ColIdent{Name: "id"},
		Tag:        tagutil.Tag{"sql": {"id"}},
	}, {
		Name:       "Email",
		Type:       TypeInfo{Kind: TypeKindString},
		IsExported: true,
		ColIdent:   ColIdent{Name: "email"},
		Tag:        tagutil.Tag{"sql": {"email"}},
	}, {
		Name:       "FullName",
		Type:       TypeInfo{Kind: TypeKindString},
		IsExported: true,
		ColIdent:   ColIdent{Name: "full_name"},
		Tag:        tagutil.Tag{"sql": {"full_name"}},
	}, {
		Name: "CreatedAt",
		Type: TypeInfo{
			Name:              "Time",
			Kind:              TypeKindStruct,
			PkgPath:           "time",
			PkgName:           "time",
			PkgLocal:          "time",
			IsImported:        true,
			IsJSONMarshaler:   true,
			IsJSONUnmarshaler: true,
		},
		IsExported: true,
		ColIdent:   ColIdent{Name: "created_at"},
		Tag:        tagutil.Tag{"sql": {"created_at"}},
	}}

	reldummyslice := &RelField{
		FieldName: "Rel",
		Id:        RelIdent{Name: "relation_a", Alias: "a"},
		Type: RelType{
			Base: TypeInfo{
				Name:     "T",
				Kind:     TypeKindStruct,
				PkgPath:  "path/to/test",
				PkgName:  "testdata",
				PkgLocal: "testdata",
			},
			IsSlice: true,
			Fields: []*FieldInfo{{
				Type:       TypeInfo{Kind: TypeKindString},
				Name:       "F",
				IsExported: true,
				Tag:        tagutil.Tag{"sql": {"f"}},
				ColIdent:   ColIdent{Name: "f"},
			}},
		},
	}

	dummyrecord := RelType{
		Base: TypeInfo{
			Name:     "T",
			Kind:     TypeKindStruct,
			PkgPath:  "path/to/test",
			PkgName:  "testdata",
			PkgLocal: "testdata",
		},
		Fields: []*FieldInfo{{
			Type:       TypeInfo{Kind: TypeKindString},
			Name:       "F",
			IsExported: true,
			Tag:        tagutil.Tag{"sql": {"f"}},
			ColIdent:   ColIdent{Name: "f"},
		}},
	}

	tests := []struct {
		Name     string
		want     *QueryStruct
		err      error
		printerr bool
	}{{
		Name: "InsertAnalysisTestBAD_NoDataField",
		err: &Error{
			Code:       errMissingRelField,
			PkgPath:    "path/to/test",
			TargetName: "InsertAnalysisTestBAD_NoDataField",
			FileName:   "../testdata/analysis_bad.go",
			FileLine:   11,
		},
	}, {
		Name: "InsertAnalysisTestBAD3",
		err: &Error{
			Code:          errBadRelType,
			PkgPath:       "path/to/test",
			TargetName:    "InsertAnalysisTestBAD3",
			BlockName:     "",
			FieldType:     "string",
			FieldTypeKind: "string",
			FieldName:     "User",
			TagString:     "",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      17,
		},
	}, {

		Name: "DeleteAnalysisTestBAD_BadRelId",
		err: &Error{
			Code:          errBadRelIdTagValue,
			PkgPath:       "path/to/test",
			TargetName:    "DeleteAnalysisTestBAD_BadRelId",
			BlockName:     "",
			FieldType:     "path/to/test.T",
			FieldTypeKind: "struct",
			FieldName:     "Rel",
			TagString:     `rel:"foo.123:bar"`,
			TagError:      "foo.123:bar",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      22,
		},
	}, {
		Name: "SelectAnalysisTestBAD_MultipleRelTags",
		err: &Error{
			Code:          errConflictingRelTag,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_MultipleRelTags",
			BlockName:     "",
			FieldType:     "path/to/test.T",
			FieldTypeKind: "struct",
			FieldName:     "Rel2",
			TagString:     `rel:""`,
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      28,
		},
	}, {
		Name: "DeleteAnalysisTestBAD_IllegalCountField",
		err: &Error{
			Code:          errIllegalQueryField,
			BlockName:     "",
			FieldType:     "int",
			FieldTypeKind: "int",
			FieldName:     "Count",
			TagString:     `rel:"relation_a:a"`,
			TargetName:    "DeleteAnalysisTestBAD_IllegalCountField",
			PkgPath:       "path/to/test",
			FileLine:      33,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "UpdateAnalysisTestBAD_IllegalExistsField",
		err: &Error{
			Code:          errIllegalQueryField,
			BlockName:     "",
			FieldType:     "bool",
			FieldTypeKind: "bool",
			FieldName:     "Exists",
			TagString:     `rel:"relation_a:a"`,
			TargetName:    "UpdateAnalysisTestBAD_IllegalExistsField",
			PkgPath:       "path/to/test",
			FileLine:      38,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "InsertAnalysisTestBAD_IllegalNotExistsField",
		err: &Error{
			Code:          errIllegalQueryField,
			BlockName:     "",
			FieldType:     "bool",
			FieldTypeKind: "bool",
			FieldName:     "NotExists",
			TagString:     `rel:"relation_a:a"`,
			TargetName:    "InsertAnalysisTestBAD_IllegalNotExistsField",
			PkgPath:       "path/to/test",
			FileLine:      43,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "SelectAnalysisTestBAD_IllegalRelationDirective",
		err: &Error{
			Code:          errIllegalQueryField,
			BlockName:     "",
			FieldType:     "github.com/frk/gosql.Relation",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `rel:"relation_a:a"`,
			TargetName:    "SelectAnalysisTestBAD_IllegalRelationDirective",
			PkgPath:       "path/to/test",
			FileLine:      48,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "SelectAnalysisTestBAD_UnnamedBaseStructType",
		err: &Error{
			Code:          errBadRelType,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_UnnamedBaseStructType",
			FieldType:     "[]*struct{}",
			FieldTypeKind: "slice",
			FieldName:     "Rel",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      53,
		},
	}, {
		Name: "SelectAnalysisTestBAD_IllegalAllDirective",
		err: &Error{
			Code:          errIllegalQueryField,
			BlockName:     "",
			FieldType:     "github.com/frk/gosql.All",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TargetName:    "SelectAnalysisTestBAD_IllegalAllDirective",
			PkgPath:       "path/to/test",
			FileLine:      59,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "InsertAnalysisTestBAD_IllegalAllDirective",
		err: &Error{
			Code:          errIllegalQueryField,
			BlockName:     "",
			FieldType:     "github.com/frk/gosql.All",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TargetName:    "InsertAnalysisTestBAD_IllegalAllDirective",
			PkgPath:       "path/to/test",
			FileLine:      65,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "UpdateAnalysisTestBAD_ConflictWhereProducer",
		err: &Error{
			Code:          errConflictingWhere,
			PkgPath:       "path/to/test",
			TargetName:    "UpdateAnalysisTestBAD_ConflictWhereProducer",
			BlockName:     "",
			FieldType:     "github.com/frk/gosql.All",
			FieldTypeKind: "struct",
			FieldName:     "_",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      74,
		},
	}, {
		Name: "DeleteAnalysisTestBAD_IllegalDefaultDirective",
		err: &Error{
			Code:          errIllegalQueryField,
			BlockName:     "",
			FieldType:     "github.com/frk/gosql.Default",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"*"`,
			TargetName:    "DeleteAnalysisTestBAD_IllegalDefaultDirective",
			PkgPath:       "path/to/test",
			FileLine:      80,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "UpdateAnalysisTestBAD_EmptyDefaultDirectiveCollist",
		err: &Error{
			Code:          errMissingTagColumnList,
			FieldType:     "github.com/frk/gosql.Default",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TargetName:    "UpdateAnalysisTestBAD_EmptyDefaultDirectiveCollist",
			PkgPath:       "path/to/test",
			FileLine:      86,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "SelectAnalysisTestBAD_IllegalForceDirective",
		err: &Error{
			Code:          errIllegalQueryField,
			BlockName:     "",
			FieldType:     "github.com/frk/gosql.Force",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"*"`,
			TargetName:    "SelectAnalysisTestBAD_IllegalForceDirective",
			PkgPath:       "path/to/test",
			FileLine:      92,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "UpdateAnalysisTestBAD_BadForceDirectiveColId",
		err: &Error{
			Code:          errBadColIdTagValue,
			PkgPath:       "path/to/test",
			TargetName:    "UpdateAnalysisTestBAD_BadForceDirectiveColId",
			BlockName:     "",
			FieldType:     "github.com/frk/gosql.Force",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"a.id,1234"`,
			TagError:      "1234",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      98,
		},
	}, {
		Name: "DeleteAnalysisTestBAD_ConflictResultProducer",
		err: &Error{
			Code:          errConflictingResultTarget,
			PkgPath:       "path/to/test",
			TargetName:    "DeleteAnalysisTestBAD_ConflictResultProducer",
			BlockName:     "",
			FieldType:     "github.com/frk/gosql.Return",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"a.id"`,
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      111,
		},
	}, {
		Name: "UpdateAnalysisTestBAD_EmptyReturnDirectiveCollist",
		err: &Error{
			Code:          errMissingTagColumnList,
			FieldType:     "github.com/frk/gosql.Return",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TargetName:    "UpdateAnalysisTestBAD_EmptyReturnDirectiveCollist",
			PkgPath:       "path/to/test",
			FileLine:      117,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "InsertAnalysisTestBAD_IllegalLimitField",
		err: &Error{
			Code:          errIllegalQueryField,
			BlockName:     "",
			FieldType:     "github.com/frk/gosql.Limit",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"10"`,
			TargetName:    "InsertAnalysisTestBAD_IllegalLimitField",
			PkgPath:       "path/to/test",
			FileLine:      123,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "UpdateAnalysisTestBAD_IllegalOffsetField",
		err: &Error{
			Code:          errIllegalQueryField,
			BlockName:     "",
			FieldType:     "github.com/frk/gosql.Offset",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"2"`,
			TargetName:    "UpdateAnalysisTestBAD_IllegalOffsetField",
			PkgPath:       "path/to/test",
			FileLine:      129,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "InsertAnalysisTestBAD_IllegalOrderByDirective",
		err: &Error{
			Code:          errIllegalQueryField,
			BlockName:     "",
			FieldType:     "github.com/frk/gosql.OrderBy",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"a.id"`,
			TargetName:    "InsertAnalysisTestBAD_IllegalOrderByDirective",
			PkgPath:       "path/to/test",
			FileLine:      135,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "DeleteAnalysisTestBAD_IllegalOverrideDirective",
		err: &Error{
			Code:          errIllegalQueryField,
			BlockName:     "",
			FieldType:     "github.com/frk/gosql.Override",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"user"`,
			TargetName:    "DeleteAnalysisTestBAD_IllegalOverrideDirective",
			PkgPath:       "path/to/test",
			FileLine:      141,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "SelectAnalysisTestBAD_IllegalTextSearchDirective",
		err: &Error{
			Code:          errIllegalQueryField,
			BlockName:     "",
			FieldType:     "github.com/frk/gosql.TextSearch",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TargetName:    "SelectAnalysisTestBAD_IllegalTextSearchDirective",
			PkgPath:       "path/to/test",
			FileLine:      147,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "SelectAnalysisTestBAD_IllegalColumnDirective",
		err: &Error{
			Code:          errIllegalQueryField,
			FieldType:     "github.com/frk/gosql.Column",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TargetName:    "SelectAnalysisTestBAD_IllegalColumnDirective",
			PkgPath:       "path/to/test",
			FileLine:      153,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "InsertAnalysisTestBAD_IllegalWhereBlock",
		err: &Error{
			Code:          errIllegalQueryField,
			BlockName:     "",
			FieldType:     "struct{Id int \"sql:\\\"id\\\"\"}",
			FieldTypeKind: "struct",
			FieldName:     "Where",
			TargetName:    "InsertAnalysisTestBAD_IllegalWhereBlock",
			PkgPath:       "path/to/test",
			FileLine:      159,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "UpdateAnalysisTestBAD_IllegalJoinBlock",
		err: &Error{
			Code:          errIllegalQueryField,
			BlockName:     "",
			FieldType:     "struct{_ github.com/frk/gosql.Relation}",
			FieldTypeKind: "struct",
			FieldName:     "Join",
			TargetName:    "UpdateAnalysisTestBAD_IllegalJoinBlock",
			PkgPath:       "path/to/test",
			FileLine:      167,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "DeleteAnalysisTestBAD_IllegalFromBlock",
		err: &Error{
			Code:          errIllegalQueryField,
			BlockName:     "",
			FieldType:     "struct{_ github.com/frk/gosql.Relation}",
			FieldTypeKind: "struct",
			FieldName:     "From",
			TargetName:    "DeleteAnalysisTestBAD_IllegalFromBlock",
			PkgPath:       "path/to/test",
			FileLine:      175,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "SelectAnalysisTestBAD_IllegalUsingBlock",
		err: &Error{
			Code:          errIllegalQueryField,
			BlockName:     "",
			FieldType:     "struct{_ github.com/frk/gosql.Relation}",
			FieldTypeKind: "struct",
			FieldName:     "Using",
			TargetName:    "SelectAnalysisTestBAD_IllegalUsingBlock",
			PkgPath:       "path/to/test",
			FileLine:      183,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "UpdateAnalysisTestBAD_IllegalOnConflictBlock",
		err: &Error{
			Code:          errIllegalQueryField,
			BlockName:     "",
			FieldType:     "struct{}",
			FieldTypeKind: "struct",
			FieldName:     "OnConflict",
			TargetName:    "UpdateAnalysisTestBAD_IllegalOnConflictBlock",
			PkgPath:       "path/to/test",
			FileLine:      191,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "SelectAnalysisTestBAD_IllegalResultField",
		err: &Error{
			Code:          errIllegalQueryField,
			BlockName:     "",
			FieldType:     "path/to/test.T",
			FieldTypeKind: "struct",
			FieldName:     "Result",
			TargetName:    "SelectAnalysisTestBAD_IllegalResultField",
			PkgPath:       "path/to/test",
			FileLine:      199,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "SelectAnalysisTestBAD_ConflictLimitProducer",
		err: &Error{
			Code:          errConflictingFieldOrDirective,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_ConflictLimitProducer",
			BlockName:     "",
			FieldType:     "int",
			FieldTypeKind: "int",
			FieldName:     "Limit",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      206,
		},
	}, {
		Name: "SelectAnalysisTestBAD_ConflictOffsetProducer",
		err: &Error{
			Code:          errConflictingFieldOrDirective,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_ConflictOffsetProducer",
			BlockName:     "",
			FieldType:     "github.com/frk/gosql.Offset",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"2"`,
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      213,
		},
	}, {
		Name: "SelectAnalysisTestBAD_IllegalRowsAffectedField",
		err: &Error{
			Code:          errIllegalQueryField,
			BlockName:     "",
			FieldType:     "int",
			FieldTypeKind: "int",
			FieldName:     "RowsAffected",
			TargetName:    "SelectAnalysisTestBAD_IllegalRowsAffectedField",
			PkgPath:       "path/to/test",
			FileLine:      219,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "InsertAnalysisTestBAD_IllegalFilterField",
		err: &Error{
			Code:          errIllegalQueryField,
			BlockName:     "",
			FieldType:     "github.com/frk/gosql.Filter",
			FieldTypeKind: "struct",
			FieldName:     "F",
			TargetName:    "InsertAnalysisTestBAD_IllegalFilterField",
			PkgPath:       "path/to/test",
			FileLine:      225,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "SelectAnalysisTestBAD_ConflictWhereProducer",
		err: &Error{
			Code:          errConflictingWhere,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_ConflictWhereProducer",
			BlockName:     "",
			FieldType:     "github.com/frk/gosql.Filter",
			FieldTypeKind: "struct",
			FieldName:     "F",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      234,
		},
	}, {
		Name: "DeleteAnalysisTestBAD_ConflictErrorHandler",
		err: &Error{
			Code:          errConflictingFieldOrDirective,
			PkgPath:       "path/to/test",
			TargetName:    "DeleteAnalysisTestBAD_ConflictErrorHandler",
			BlockName:     "",
			FieldType:     "path/to/test.myerrorhandler",
			FieldTypeKind: "struct",
			FieldName:     "erh",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      241,
		},
	}, {
		Name: "SelectAnalysisTestBAD_IteratorWithTooManyMethods",
		err: &Error{
			Code:          errBadIterTypeInterface,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_IteratorWithTooManyMethods",
			FieldType:     "path/to/test.badIterator",
			FieldTypeKind: "interface",
			FieldName:     "Rel",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      255,
		},
	}, {
		Name: "SelectAnalysisTestBAD_IteratorWithBadSignature",
		err: &Error{
			Code:          errBadIterTypeFunc,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_IteratorWithBadSignature",
			FieldType:     "func(*github.com/frk/gosql/internal/testdata/common.User) int",
			FieldTypeKind: "func",
			FieldName:     "Rel",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      260,
		},
	}, {
		Name: "SelectAnalysisTestBAD_IteratorWithBadSignatureIface",
		err: &Error{
			Code:          errBadIterTypeInterface,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_IteratorWithBadSignatureIface",
			FieldType:     "path/to/test.badIterator2",
			FieldTypeKind: "interface",
			FieldName:     "Rel",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      265,
		},
	}, {
		Name: "SelectAnalysisTestBAD_IteratorWithUnexportedMethod",
		err: &Error{
			Code:          errBadIterTypeInterface,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_IteratorWithUnexportedMethod",
			FieldType:     "github.com/frk/gosql/internal/testdata/common.BadIterator",
			FieldTypeKind: "interface",
			FieldName:     "Rel",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      270,
		},
	}, {
		Name: "SelectAnalysisTestBAD_IteratorWithUnnamedArgument",
		err: &Error{
			Code:          errBadIterTypeFunc,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_IteratorWithUnnamedArgument",
			FieldType:     "func(*struct{}) error",
			FieldTypeKind: "func",
			FieldName:     "Rel",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      275,
		},
	}, {
		Name: "SelectAnalysisTestBAD_IteratorWithNonStructArgument",
		err: &Error{
			Code:          errBadIterTypeFunc,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_IteratorWithNonStructArgument",
			FieldType:     "func(*path/to/test.notstruct) error",
			FieldTypeKind: "func",
			FieldName:     "Rel",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      280,
		},
	}, {
		Name: "InsertAnalysisTestBAD_BadRelfiedlStructBaseType",
		err: &Error{
			Code:          errBadRelType,
			PkgPath:       "path/to/test",
			TargetName:    "InsertAnalysisTestBAD_BadRelfiedlStructBaseType",
			FieldType:     "[]*path/to/test.notstruct",
			FieldTypeKind: "slice",
			FieldName:     "Rel",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      287,
		},
	}, {
		Name: "UpdateAnalysisTestBAD_BadRelTypeFieldColId",
		err: &Error{
			Code:          errBadColIdTagValue,
			PkgPath:       "path/to/test",
			TargetName:    "UpdateAnalysisTestBAD_BadRelTypeFieldColId",
			BlockName:     "",
			FieldType:     "string",
			FieldTypeKind: "string",
			FieldName:     "Foo",
			TagString:     `sql:"1234"`,
			TagError:      "1234",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      293,
		},
	}, {
		Name: "UpdateAnalysisTestBAD_ConflictWhereProducer2",
		err: &Error{
			Code:          errConflictingWhere,
			PkgPath:       "path/to/test",
			TargetName:    "UpdateAnalysisTestBAD_ConflictWhereProducer2",
			BlockName:     "",
			FieldType:     "struct{Id int \"sql:\\\"id\\\"\"}",
			FieldTypeKind: "struct",
			FieldName:     "Where",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      301,
		},
	}, {
		Name: "DeleteAnalysisTestBAD_BadWhereBlockType",
		err: &Error{
			Code:          errBadFieldTypeStruct,
			PkgPath:       "path/to/test",
			TargetName:    "DeleteAnalysisTestBAD_BadWhereBlockType",
			BlockName:     "",
			FieldType:     "[]string",
			FieldTypeKind: "slice",
			FieldName:     "Where",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      309,
		},
	}, {
		Name: "SelectAnalysisTestBAD_BadBoolTagValue",
		err: &Error{
			Code:          errBadBoolTagValue,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_BadBoolTagValue",
			BlockName:     "Where",
			FieldType:     "string",
			FieldTypeKind: "string",
			FieldName:     "Name",
			TagString:     `sql:"name" bool:"abc"`,
			TagError:      "abc",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      317,
		},
	}, {
		Name: "SelectAnalysisTestBAD_BadNestedWhereBlockType",
		err: &Error{
			Code:          errBadFieldTypeStruct,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_BadNestedWhereBlockType",
			BlockName:     "Where",
			FieldType:     "path/to/test.notstruct",
			FieldTypeKind: "string",
			FieldName:     "X",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      326,
		},
	}, {
		Name: "SelectAnalysisTestBAD_BadColumnExpressionLHS",
		err: &Error{
			Code:          errBadColIdTagValue,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_BadColumnExpressionLHS",
			BlockName:     "Where",
			FieldType:     "github.com/frk/gosql.Column",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"123 = x"`,
			TagError:      "123",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      334,
		},
	}, {
		Name:     "SelectAnalysisTestBAD_BadColumnPredicateCombo",
		printerr: true,
		err: &Error{
			Code:          errIllegalPredicateQuantifier,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_BadColumnPredicateCombo",
			BlockName:     "Where",
			FieldType:     "github.com/frk/gosql.Column",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"x isin any y"`,
			TagExpr:       `x isin any y`,
			TagError:      "any",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      342,
		},
	}, {
		Name: "DeleteAnalysisTestBAD_BadColumnExpressionLHS",
		err: &Error{
			Code:          errBadColIdTagValue,
			PkgPath:       "path/to/test",
			TargetName:    "DeleteAnalysisTestBAD_BadColumnExpressionLHS",
			BlockName:     "Where",
			FieldType:     "github.com/frk/gosql.Column",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"123 isnull"`,
			TagError:      "123",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      350,
		},
	}, {
		Name: "UpdateAnalysisTestBAD_BadUnaryOp",
		err: &Error{
			Code:          errBadDirectiveBooleanExpr,
			PkgPath:       "path/to/test",
			TargetName:    "UpdateAnalysisTestBAD_BadUnaryOp",
			BlockName:     "Where",
			FieldType:     "github.com/frk/gosql.Column",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"x <="`,
			TagError:      "x <=",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      358,
		},
	}, {
		Name:     "UpdateAnalysisTestBAD_ExtraQuantifier",
		printerr: true,
		err: &Error{
			Code:          errIllegalPredicateQuantifier,
			PkgPath:       "path/to/test",
			TargetName:    "UpdateAnalysisTestBAD_ExtraQuantifier",
			BlockName:     "Where",
			FieldType:     "github.com/frk/gosql.Column",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"x isnull any"`,
			TagExpr:       `x isnull any`,
			TagError:      "any",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      366,
		},
	}, {
		Name: "SelectAnalysisTestBAD_BadBetweenFieldType",
		err: &Error{
			Code:          errBadBetweenPredicate,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_BadBetweenFieldType",
			BlockName:     "Where",
			FieldType:     "path/to/test.notstruct",
			FieldTypeKind: "string",
			FieldName:     "between",
			TagString:     `sql:"a.foo isbetween"`,
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      374,
		},
	}, {
		Name: "SelectAnalysisTestBAD_BadBetweenFieldType2",
		err: &Error{
			Code:          errBadBetweenPredicate,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_BadBetweenFieldType2",
			BlockName:     "Where",
			FieldType:     "struct{x int; y int; z int}",
			FieldTypeKind: "struct",
			FieldName:     "between",
			TagString:     `sql:"a.foo isbetween"`,
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      382,
		},
	}, {
		Name: "SelectAnalysisTestBAD_BadBetweenArgColId",
		err: &Error{
			Code:          errBadColIdTagValue,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_BadBetweenArgColId",
			BlockName:     "between",
			FieldType:     "github.com/frk/gosql.Column",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"123,y"`,
			TagError:      "123",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      394,
		},
	}, {
		Name: "SelectAnalysisTestBAD_NoBetweenXYArg",
		err: &Error{
			Code:          errBadBetweenPredicate,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_NoBetweenXYArg",
			BlockName:     "Where",
			FieldType:     "struct{_ github.com/frk/gosql.Column \"sql:\\\"a.bar\\\"\"; _ github.com/frk/gosql.Column \"sql:\\\"a.baz,y\\\"\"}",
			FieldTypeKind: "struct",
			FieldName:     "between",
			TagString:     `sql:"a.foo isbetween"`,
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      403,
		},
	}, {
		Name: "SelectAnalysisTestBAD_BadBetweenColId",
		err: &Error{
			Code:          errBadColIdTagValue,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_BadBetweenColId",
			BlockName:     "Where",
			FieldType:     "struct{_ github.com/frk/gosql.Column \"sql:\\\"a.bar,x\\\"\"; _ github.com/frk/gosql.Column \"sql:\\\"a.baz,y\\\"\"}",
			FieldTypeKind: "struct",
			FieldName:     "between",
			TagString:     `sql:"123 isbetween"`,
			TagError:      "123",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      414,
		},
	}, {
		Name: "DeleteAnalysisTestBAD_BadWhereFieldColId",
		err: &Error{
			Code:          errBadColIdTagValue,
			PkgPath:       "path/to/test",
			TargetName:    "DeleteAnalysisTestBAD_BadWhereFieldColId",
			BlockName:     "Where",
			FieldType:     "int",
			FieldTypeKind: "int",
			FieldName:     "Id",
			TagString:     `sql:"123"`,
			TagError:      "123",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      425,
		},
	}, {
		Name:     "DeleteAnalysisTestBAD_BadWhereFieldPredicateCombo",
		printerr: true,
		err: &Error{
			Code:          errIllegalPredicateQuantifier,
			PkgPath:       "path/to/test",
			TargetName:    "DeleteAnalysisTestBAD_BadWhereFieldPredicateCombo",
			BlockName:     "Where",
			FieldType:     "[]int",
			FieldTypeKind: "slice",
			FieldName:     "Id",
			TagString:     `sql:"a.id notin any"`,
			TagExpr:       `a.id notin any`,
			TagError:      "any",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      433,
		},
	}, {
		Name: "DeleteAnalysisTestBAD_IllegalWhereFieldUnaryPredicate",
		err: &Error{
			Code:          errIllegalUnaryPredicate,
			PkgPath:       "path/to/test",
			TargetName:    "DeleteAnalysisTestBAD_IllegalWhereFieldUnaryPredicate",
			BlockName:     "Where",
			FieldType:     "int",
			FieldTypeKind: "int",
			FieldName:     "Id",
			TagString:     `sql:"a.id istrue"`,
			TagExpr:       `a.id istrue`,
			TagError:      `istrue`,
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      441,
		},
	}, {
		Name: "UpdateAnalysisTestBAD_BadWhereFieldTypeForQuantifier",
		err: &Error{
			Code:          errIllegalFieldQuantifier,
			PkgPath:       "path/to/test",
			TargetName:    "UpdateAnalysisTestBAD_BadWhereFieldTypeForQuantifier",
			BlockName:     "Where",
			FieldType:     "int",
			FieldTypeKind: "int",
			FieldName:     "Id",
			TagString:     `sql:"a.id = any"`,
			TagExpr:       `a.id = any`,
			TagError:      "any",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      449,
		},
	}, {
		Name: "SelectAnalysisTestBAD_BadJoinBlockType",
		err: &Error{
			Code:          errBadFieldTypeStruct,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_BadJoinBlockType",
			BlockName:     "",
			FieldType:     "path/to/test.notstruct",
			FieldTypeKind: "string",
			FieldName:     "Join",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      456,
		},
	}, {
		Name: "SelectAnalysisTestBAD_IllegalJoinBlockRelationDirective",
		err: &Error{
			Code:          errIllegalStructDirective,
			TargetName:    "SelectAnalysisTestBAD_IllegalJoinBlockRelationDirective",
			BlockName:     "Join",
			FieldType:     "github.com/frk/gosql.Relation",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"foobar"`,
			PkgPath:       "path/to/test",
			FileLine:      463,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "DeleteAnalysisTestBAD_ConflictRelationDirective",
		err: &Error{
			Code:          errConflictingRelationDirective,
			PkgPath:       "path/to/test",
			TargetName:    "DeleteAnalysisTestBAD_ConflictRelationDirective",
			BlockName:     "Using",
			FieldType:     "github.com/frk/gosql.Relation",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"bar"`,
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      472,
		},
	}, {
		Name: "UpdateAnalysisTestBAD_BadFromRelationRelId",
		err: &Error{
			Code:          errBadRelIdTagValue,
			PkgPath:       "path/to/test",
			TargetName:    "UpdateAnalysisTestBAD_BadFromRelationRelId",
			BlockName:     "From",
			FieldType:     "github.com/frk/gosql.Relation",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"123"`,
			TagError:      "123",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      480,
		},
	}, {
		Name: "SelectAnalysisTestBAD_BadJoinDirectiveRelId",
		err: &Error{
			Code:          errBadRelIdTagValue,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_BadJoinDirectiveRelId",
			BlockName:     "Join",
			FieldType:     "github.com/frk/gosql.LeftJoin",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"123"`,
			TagError:      "123",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      488,
		},
	}, {
		Name: "SelectAnalysisTestBAD_BadJoinDirectiveExpressionColId",
		err: &Error{
			Code:          errBadColIdTagValue,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_BadJoinDirectiveExpressionColId",
			BlockName:     "Join",
			FieldType:     "github.com/frk/gosql.LeftJoin",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"relation_b:b,123 = b.foo"`,
			TagError:      "123",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      496,
		},
	}, {
		Name: "SelectAnalysisTestBAD_BadJoinDirectiveExpressionPredicate",
		err: &Error{
			Code:          errBadDirectiveBooleanExpr,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_BadJoinDirectiveExpressionPredicate",
			BlockName:     "Join",
			FieldType:     "github.com/frk/gosql.LeftJoin",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"relation_b:b,b.foo ="`,
			TagError:      "b.foo =",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      504,
		},
	}, {
		Name:     "SelectAnalysisTestBAD_BadJoinDirectiveExpressionExtraQuantifier",
		printerr: true,
		err: &Error{
			Code:          errIllegalPredicateQuantifier,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_BadJoinDirectiveExpressionExtraQuantifier",
			BlockName:     "Join",
			FieldType:     "github.com/frk/gosql.LeftJoin",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"relation_b:b,b.foo isnull any"`,
			TagExpr:       `b.foo isnull any`,
			TagError:      "any",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      512,
		},
	}, {
		Name:     "SelectAnalysisTestBAD_BadJoinDirectiveExpressionPredicateCombo",
		printerr: true,
		err: &Error{
			Code:          errIllegalPredicateQuantifier,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_BadJoinDirectiveExpressionPredicateCombo",
			BlockName:     "Join",
			FieldType:     "github.com/frk/gosql.LeftJoin",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"relation_b:b,b.foo isin any a.bar"`,
			TagExpr:       `b.foo isin any a.bar`,
			TagError:      "any",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      520,
		},
	}, {
		Name: "DeleteAnalysisTestBAD_IllegalJoinBlockDirective",
		err: &Error{
			Code:          errIllegalStructDirective,
			PkgPath:       "path/to/test",
			TargetName:    "DeleteAnalysisTestBAD_IllegalJoinBlockDirective",
			BlockName:     "Using",
			FieldType:     "github.com/frk/gosql.Column",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"a.foo"`,
			FileLine:      528,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "InsertAnalysisTestBAD_BadOnConflictBlockType",
		err: &Error{
			Code:          errBadFieldTypeStruct,
			PkgPath:       "path/to/test",
			TargetName:    "InsertAnalysisTestBAD_BadOnConflictBlockType",
			BlockName:     "",
			FieldType:     "path/to/test.notstruct",
			FieldTypeKind: "string",
			FieldName:     "OnConflict",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      535,
		},
	}, {
		Name: "InsertAnalysisTestBAD_ConflictOnConflictBlockTargetProducer",
		err: &Error{
			Code:          errConflictingOnConfictTarget,
			PkgPath:       "path/to/test",
			TargetName:    "InsertAnalysisTestBAD_ConflictOnConflictBlockTargetProducer",
			BlockName:     "OnConflict",
			FieldType:     "github.com/frk/gosql.Column",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"a.id"`,
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      543,
		},
	}, {
		Name: "InsertAnalysisTestBAD_ConflictOnConflictBlockTargetProducer2",
		err: &Error{
			Code:          errConflictingOnConfictTarget,
			PkgPath:       "path/to/test",
			TargetName:    "InsertAnalysisTestBAD_ConflictOnConflictBlockTargetProducer2",
			BlockName:     "OnConflict",
			FieldType:     "github.com/frk/gosql.Index",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"some_index"`,
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      552,
		},
	}, {
		Name: "InsertAnalysisTestBAD_ConflictOnConflictBlockTargetProducer3",
		err: &Error{
			Code:          errConflictingOnConfictTarget,
			PkgPath:       "path/to/test",
			TargetName:    "InsertAnalysisTestBAD_ConflictOnConflictBlockTargetProducer3",
			BlockName:     "OnConflict",
			FieldType:     "github.com/frk/gosql.Constraint",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"some_constraint"`,
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      561,
		},
	}, {
		Name: "InsertAnalysisTestBAD_ConflictOnConflictBlockActionProducer",
		err: &Error{
			Code:          errConflictingOnConfictAction,
			PkgPath:       "path/to/test",
			TargetName:    "InsertAnalysisTestBAD_ConflictOnConflictBlockActionProducer",
			BlockName:     "OnConflict",
			FieldType:     "github.com/frk/gosql.Ignore",
			FieldTypeKind: "struct",
			FieldName:     "_",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      571,
		},
	}, {
		Name: "InsertAnalysisTestBAD_ConflictOnConflictBlockActionProducer2",
		err: &Error{
			Code:          errConflictingOnConfictAction,
			PkgPath:       "path/to/test",
			TargetName:    "InsertAnalysisTestBAD_ConflictOnConflictBlockActionProducer2",
			BlockName:     "OnConflict",
			FieldType:     "github.com/frk/gosql.Update",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"a.foo"`,
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      581,
		},
	}, {
		Name: "InsertAnalysisTestBAD_BadOnConflictColumnTargetValue",
		err: &Error{
			Code:          errBadColIdTagValue,
			PkgPath:       "path/to/test",
			TargetName:    "InsertAnalysisTestBAD_BadOnConflictColumnTargetValue",
			BlockName:     "OnConflict",
			FieldType:     "github.com/frk/gosql.Column",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"a.id,a.1234"`,
			TagError:      "a.1234",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      589,
		},
	}, {
		Name: "InsertAnalysisTestBAD_BadOnConflictIndexTargetIdent",
		err: &Error{
			Code:          errBadIdentTagValue,
			PkgPath:       "path/to/test",
			TargetName:    "InsertAnalysisTestBAD_BadOnConflictIndexTargetIdent",
			BlockName:     "OnConflict",
			FieldType:     "github.com/frk/gosql.Index",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"1234"`,
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      597,
		},
	}, {
		Name: "InsertAnalysisTestBAD_BadOnConflictConstraintTargetIdent",
		err: &Error{
			Code:          errBadIdentTagValue,
			PkgPath:       "path/to/test",
			TargetName:    "InsertAnalysisTestBAD_BadOnConflictConstraintTargetIdent",
			BlockName:     "OnConflict",
			FieldType:     "github.com/frk/gosql.Constraint",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"1234"`,
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      605,
		},
	}, {
		Name: "InsertAnalysisTestBAD_BadOnConflictUpdateActionCollist",
		err: &Error{
			Code:          errBadColIdTagValue,
			PkgPath:       "path/to/test",
			TargetName:    "InsertAnalysisTestBAD_BadOnConflictUpdateActionCollist",
			BlockName:     "OnConflict",
			FieldType:     "github.com/frk/gosql.Update",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"a.id,a.1234"`,
			TagError:      "a.1234",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      614,
		},
	}, {
		Name: "InsertAnalysisTestBAD_IllegalOnConflictDirective",
		err: &Error{
			Code:          errIllegalStructDirective,
			TargetName:    "InsertAnalysisTestBAD_IllegalOnConflictDirective",
			BlockName:     "OnConflict",
			FieldType:     "github.com/frk/gosql.LeftJoin",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"relation_b:b,b.id=a.id"`,
			PkgPath:       "path/to/test",
			FileLine:      622,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "InsertAnalysisTestBAD_NoOnConflictTarget",
		err: &Error{
			Code:          errMissingOnConflictTarget,
			PkgPath:       "path/to/test",
			TargetName:    "InsertAnalysisTestBAD_NoOnConflictTarget",
			FieldType:     "struct{_ github.com/frk/gosql.Update \"sql:\\\"a.foo,a.bar\\\"\"}",
			FieldTypeKind: "struct",
			FieldName:     "OnConflict",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      629,
		},
	}, {
		Name: "SelectAnalysisTestBAD_BadLimitFieldType",
		err: &Error{
			Code:          errBadFieldTypeInt,
			TargetName:    "SelectAnalysisTestBAD_BadLimitFieldType",
			FieldType:     "string",
			FieldTypeKind: "string",
			FieldName:     "Limit",
			TagString:     `sql:"123"`,
			PkgPath:       "path/to/test",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      637,
		},
	}, {
		Name: "SelectAnalysisTestBAD_NoLimitDirectiveValue",
		err: &Error{
			Code:          errMissingTagValue,
			TargetName:    "SelectAnalysisTestBAD_NoLimitDirectiveValue",
			FieldType:     "github.com/frk/gosql.Limit",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:""`,
			PkgPath:       "path/to/test",
			FileLine:      643,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "SelectAnalysisTestBAD_BadLimitDirectiveValue",
		err: &Error{
			Code:          errBadUIntegerTagValue,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_BadLimitDirectiveValue",
			FieldType:     "github.com/frk/gosql.Limit",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"abc"`,
			TagError:      `abc`,
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      649,
		},
	}, {
		Name: "SelectAnalysisTestBAD_BadOffsetFieldType",
		err: &Error{
			Code:          errBadFieldTypeInt,
			FieldType:     "string",
			FieldTypeKind: "string",
			FieldName:     "Offset",
			TagString:     `sql:"123"`,
			TargetName:    "SelectAnalysisTestBAD_BadOffsetFieldType",
			PkgPath:       "path/to/test",
			FileLine:      655,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "SelectAnalysisTestBAD_NoOffsetDirectiveValue",
		err: &Error{
			Code:          errMissingTagValue,
			FieldType:     "github.com/frk/gosql.Offset",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:""`,
			TargetName:    "SelectAnalysisTestBAD_NoOffsetDirectiveValue",
			PkgPath:       "path/to/test",
			FileLine:      661,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "SelectAnalysisTestBAD_BadOffsetDirectiveValue",
		err: &Error{
			Code:          errBadUIntegerTagValue,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_BadOffsetDirectiveValue",
			FieldType:     "github.com/frk/gosql.Offset",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"abc"`,
			TagError:      "abc",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      667,
		},
	}, {
		Name: "SelectAnalysisTestBAD_EmptyOrderByDirectiveCollist",
		err: &Error{
			Code:          errMissingTagColumnList,
			FieldType:     "github.com/frk/gosql.OrderBy",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TargetName:    "SelectAnalysisTestBAD_EmptyOrderByDirectiveCollist",
			PkgPath:       "path/to/test",
			FileLine:      673,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "SelectAnalysisTestBAD_BadOrderByDirectiveNullsOrderValue",
		err: &Error{
			Code:          errBadNullsOrderTagValue,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_BadOrderByDirectiveNullsOrderValue",
			BlockName:     "",
			FieldType:     "github.com/frk/gosql.OrderBy",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     "a.id:nullsthird",
			TagError:      "nullsthird",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      679,
		},
	}, {
		Name: "SelectAnalysisTestBAD_BadOrderByDirectiveCollist",
		err: &Error{
			Code:          errBadColIdTagValue,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_BadOrderByDirectiveCollist",
			BlockName:     "",
			FieldType:     "github.com/frk/gosql.OrderBy",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"-a.id:nullsfirst,a.1234"`,
			TagError:      "a.1234",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      685,
		},
	}, {
		Name: "InsertAnalysisTestBAD_BadOverrideDirectiveKindValue",
		err: &Error{
			Code:          errBadOverrideTagValue,
			PkgPath:       "path/to/test",
			TargetName:    "InsertAnalysisTestBAD_BadOverrideDirectiveKindValue",
			BlockName:     "",
			FieldType:     "github.com/frk/gosql.Override",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"foo"`,
			TagError:      "foo",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      691,
		},
	}, {
		Name: "UpdateAnalysisTestBAD_ConflictResultProducer",
		err: &Error{
			Code:          errConflictingResultTarget,
			PkgPath:       "path/to/test",
			TargetName:    "UpdateAnalysisTestBAD_ConflictResultProducer",
			BlockName:     "",
			FieldType:     "[]path/to/test.T",
			FieldTypeKind: "slice",
			FieldName:     "Result",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      698,
		},
	}, {
		Name: "UpdateAnalysisTestBAD_BadResultFieldType",
		err: &Error{
			Code:          errBadRelType,
			PkgPath:       "path/to/test",
			TargetName:    "UpdateAnalysisTestBAD_BadResultFieldType",
			FieldType:     "[]path/to/test.notstruct",
			FieldTypeKind: "slice",
			FieldName:     "Result",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      704,
		},
	}, {
		Name: "DeleteAnalysisTestBAD_ConflictResultProducer2",
		err: &Error{
			Code:          errConflictingResultTarget,
			PkgPath:       "path/to/test",
			TargetName:    "DeleteAnalysisTestBAD_ConflictResultProducer2",
			BlockName:     "",
			FieldType:     "int",
			FieldTypeKind: "int",
			FieldName:     "RowsAffected",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      711,
		},
	}, {
		Name: "DeleteAnalysisTestBAD_BadRowsAffecteFieldType",
		err: &Error{
			Code:          errBadFieldTypeInt,
			FieldType:     "string",
			FieldTypeKind: "string",
			FieldName:     "RowsAffected",
			TargetName:    "DeleteAnalysisTestBAD_BadRowsAffecteFieldType",
			PkgPath:       "path/to/test",
			FileLine:      717,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "UpdateAnalysisTestBAD_IllegalAllDirective",
		err: &Error{
			Code:          errIllegalSliceUpdateModifier,
			TargetName:    "UpdateAnalysisTestBAD_IllegalAllDirective",
			PkgPath:       "path/to/test",
			FieldType:     "github.com/frk/gosql.All",
			FieldTypeKind: "struct",
			FieldName:     "_",
			FileLine:      729,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "UpdateAnalysisTestBAD_IllegalWhereStruct",
		err: &Error{
			Code:          errIllegalSliceUpdateModifier,
			TargetName:    "UpdateAnalysisTestBAD_IllegalWhereStruct",
			PkgPath:       "path/to/test",
			FieldType:     `struct{Name string "sql:\"name\""}`,
			FieldTypeKind: "struct",
			FieldName:     "Where",
			FileLine:      735,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "UpdateAnalysisTestBAD_IllegalFilterField",
		err: &Error{
			Code:          errIllegalSliceUpdateModifier,
			TargetName:    "UpdateAnalysisTestBAD_IllegalFilterField",
			PkgPath:       "path/to/test",
			FieldType:     "github.com/frk/gosql.Filter",
			FieldTypeKind: "struct",
			FieldName:     "F",
			FileLine:      743,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "DeleteAnalysisTestBAD_IllegalUnaryPredicateInExpression",
		err: &Error{
			Code:          errIllegalUnaryPredicate,
			PkgPath:       "path/to/test",
			TargetName:    "DeleteAnalysisTestBAD_IllegalUnaryPredicateInExpression",
			BlockName:     "Where",
			FieldType:     "github.com/frk/gosql.Column",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"a.id isfalse a.foo"`,
			TagExpr:       "a.id isfalse a.foo",
			TagError:      "isfalse",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      750,
		},
	}, {
		Name: "SelectAnalysisTestBAD_IllegalUnaryPredicateInJoinDirectiveExpression",
		err: &Error{
			Code:          errIllegalUnaryPredicate,
			PkgPath:       "path/to/test",
			TargetName:    "SelectAnalysisTestBAD_IllegalUnaryPredicateInJoinDirectiveExpression",
			BlockName:     "Join",
			FieldType:     "github.com/frk/gosql.LeftJoin",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"relation_b:b,b.foo istrue a.bar"`,
			TagExpr:       "b.foo istrue a.bar",
			TagError:      "istrue",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      758,
		},
	}, {
		Name: "DeleteAnalysisTestBAD_ListPredicate",
		err: &Error{
			Code:          errIllegalListPredicate,
			PkgPath:       "path/to/test",
			TargetName:    "DeleteAnalysisTestBAD_ListPredicate",
			BlockName:     "Where",
			FieldType:     "int",
			FieldTypeKind: "int",
			FieldName:     "a",
			TagString:     `sql:"column_a isin"`,
			TagExpr:       "column_a isin",
			TagError:      "isin",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      766,
		},
	}, {
		Name: "InsertAnalysisTestOK1",
		want: &QueryStruct{
			TypeName: "InsertAnalysisTestOK1",
			Kind:     QueryKindInsert,
			Rel: &RelField{
				FieldName: "UserRec",
				Id:        RelIdent{Name: "users_table"},
				Type: RelType{
					Base:      commonUserTypeinfo,
					Fields:    commonUserFields,
					IsPointer: true,
				},
			},
		},
	}, {
		Name: "InsertAnalysisTestOK2",
		want: &QueryStruct{
			TypeName: "InsertAnalysisTestOK2",
			Kind:     QueryKindInsert,
			Rel: &RelField{
				FieldName: "UserRec",
				Id:        RelIdent{Name: "users_table"},
				Type: RelType{
					Base: TypeInfo{
						Kind: TypeKindStruct,
					},
					Fields: []*FieldInfo{{
						Name:       "Name3",
						Type:       TypeInfo{Kind: TypeKindString},
						IsExported: true,
						ColIdent:   ColIdent{Name: "name"},
						Tag:        tagutil.Tag{"sql": {"name"}},
					}},
				},
			},
		},
	}, {
		Name: "SelectAnalysisTestOK3",
		want: &QueryStruct{
			TypeName: "SelectAnalysisTestOK3",
			Kind:     QueryKindSelect,
			Rel: &RelField{
				FieldName: "User",
				Id:        RelIdent{Name: "users_table"},
				Type: RelType{
					Base:      commonUserTypeinfo,
					Fields:    commonUserFields,
					IsPointer: true,
					IsIter:    true,
				},
			},
		},
	}, {
		Name: "SelectAnalysisTestOK4",
		want: &QueryStruct{
			TypeName: "SelectAnalysisTestOK4",
			Kind:     QueryKindSelect,
			Rel: &RelField{
				FieldName: "User",
				Id:        RelIdent{Name: "users_table"},
				Type: RelType{
					Base:      commonUserTypeinfo,
					Fields:    commonUserFields,
					IsPointer: true,
					IsIter:    true,
				},
			},
		},
	}, {
		Name: "SelectAnalysisTestOK5",
		want: &QueryStruct{
			TypeName: "SelectAnalysisTestOK5",
			Kind:     QueryKindSelect,
			Rel: &RelField{
				FieldName: "User",
				Id:        RelIdent{Name: "users_table"},
				Type: RelType{
					Base:       commonUserTypeinfo,
					Fields:     commonUserFields,
					IsPointer:  true,
					IsIter:     true,
					IterMethod: "Fn",
				},
			},
		},
	}, {
		Name: "SelectAnalysisTestOK6",
		want: &QueryStruct{
			TypeName: "SelectAnalysisTestOK6",
			Kind:     QueryKindSelect,
			Rel: &RelField{
				FieldName: "User",
				Id:        RelIdent{Name: "users_table"},
				Type: RelType{
					Base:       commonUserTypeinfo,
					Fields:     commonUserFields,
					IsPointer:  true,
					IsIter:     true,
					IterMethod: "Fn",
				},
			},
		},
	}, {
		Name: "SelectAnalysisTestOK7",
		want: &QueryStruct{
			TypeName: "SelectAnalysisTestOK7",
			Kind:     QueryKindSelect,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a"},
				Type: RelType{
					Base: TypeInfo{
						Kind: TypeKindStruct,
					},
					Fields: []*FieldInfo{{
						Name:     "a",
						Type:     TypeInfo{Kind: TypeKindInt},
						ColIdent: ColIdent{Name: "a"},
						Tag:      tagutil.Tag{"sql": {"a", "pk"}},
					}, {
						Name:      "b",
						Type:      TypeInfo{Kind: TypeKindInt},
						ColIdent:  ColIdent{Name: "b"},
						Tag:       tagutil.Tag{"sql": {"b", "nullempty"}},
						NullEmpty: true,
					}, {
						Name:     "c",
						Type:     TypeInfo{Kind: TypeKindInt},
						ColIdent: ColIdent{Name: "c"},
						Tag:      tagutil.Tag{"sql": {"c", "ro", "json"}},
						ReadOnly: true,
					}, {
						Name:      "d",
						Type:      TypeInfo{Kind: TypeKindInt},
						ColIdent:  ColIdent{Name: "d"},
						Tag:       tagutil.Tag{"sql": {"d", "wo"}},
						WriteOnly: true,
					}, {
						Name:     "e",
						Type:     TypeInfo{Kind: TypeKindInt},
						ColIdent: ColIdent{Name: "e"},
						Tag:      tagutil.Tag{"sql": {"e", "add"}},
						UseAdd:   true,
					}, {
						Name:        "f",
						Type:        TypeInfo{Kind: TypeKindInt},
						ColIdent:    ColIdent{Name: "f"},
						Tag:         tagutil.Tag{"sql": {"f", "coalesce"}},
						UseCoalesce: true,
					}, {
						Name:          "g",
						Type:          TypeInfo{Kind: TypeKindInt},
						ColIdent:      ColIdent{Name: "g"},
						Tag:           tagutil.Tag{"sql": {"g", "coalesce(-1)"}},
						UseCoalesce:   true,
						CoalesceValue: "-1",
					}},
				},
			},
		},
	}, {
		Name: "InsertAnalysisTestOK8",
		want: &QueryStruct{
			TypeName: "InsertAnalysisTestOK8",
			Kind:     QueryKindInsert,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a"},
				Type: RelType{
					Base: TypeInfo{
						Kind: TypeKindStruct,
					},
					Fields: []*FieldInfo{{
						Name: "Val",
						Selector: []*FieldSelectorNode{
							{
								Name:         "Foobar",
								Tag:          tagutil.Tag{"sql": {">foo_"}},
								TypeName:     "Foo",
								TypePkgPath:  "github.com/frk/gosql/internal/testdata/common",
								TypePkgName:  "common",
								TypePkgLocal: "common",
								IsExported:   true,
								IsImported:   true,
							},
							{
								Name:         "Bar",
								Tag:          tagutil.Tag{"sql": {">bar_"}},
								TypeName:     "Bar",
								TypePkgPath:  "github.com/frk/gosql/internal/testdata/common",
								TypePkgName:  "common",
								TypePkgLocal: "common",
								IsImported:   true,
								IsExported:   true,
							},
							{
								Name:         "Baz",
								Tag:          tagutil.Tag{"sql": {">baz_"}},
								TypeName:     "Baz",
								TypePkgPath:  "github.com/frk/gosql/internal/testdata/common",
								TypePkgName:  "common",
								TypePkgLocal: "common",
								IsExported:   true,
								IsEmbedded:   true,
								IsImported:   true,
							},
						},
						IsExported: true,
						Type:       TypeInfo{Kind: TypeKindString},
						ColIdent:   ColIdent{Name: "foo_bar_baz_val"},
						Tag:        tagutil.Tag{"sql": {"val"}},
					}, {
						Name: "Val",
						Selector: []*FieldSelectorNode{{
							Name:         "Foobar",
							Tag:          tagutil.Tag{"sql": {">foo_"}},
							TypeName:     "Foo",
							TypePkgPath:  "github.com/frk/gosql/internal/testdata/common",
							TypePkgName:  "common",
							TypePkgLocal: "common",
							IsExported:   true,
							IsImported:   true,
						}, {
							Name:         "Baz",
							Tag:          tagutil.Tag{"sql": {">baz_"}},
							TypeName:     "Baz",
							TypePkgPath:  "github.com/frk/gosql/internal/testdata/common",
							TypePkgName:  "common",
							TypePkgLocal: "common",
							IsImported:   true,
							IsExported:   true,
							IsEmbedded:   false,
							IsPointer:    true,
						}},
						IsExported: true,
						Type:       TypeInfo{Kind: TypeKindString},
						ColIdent:   ColIdent{Name: "foo_baz_val"},
						Tag:        tagutil.Tag{"sql": {"val"}},
					}},
				},
			},
		},
	}, {
		Name: "DeleteAnalysisTestOK9",
		want: &QueryStruct{
			TypeName: "DeleteAnalysisTestOK9",
			Kind:     QueryKindDelete,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a"},
				Type:      RelType{Base: TypeInfo{Kind: TypeKindStruct}},
			},
			Where: &WhereStruct{FieldName: "Where", Items: []WhereItem{
				&WhereStructField{
					Name:      "ID",
					Type:      TypeInfo{Kind: TypeKindInt},
					ColIdent:  ColIdent{Name: "id"},
					Predicate: IsEQ,
				},
			}},
		},
	}, {
		Name: "DeleteAnalysisTestOK10",
		want: &QueryStruct{
			TypeName: "DeleteAnalysisTestOK10",
			Kind:     QueryKindDelete,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a"},
				Type:      RelType{Base: TypeInfo{Kind: TypeKindStruct}},
			},
			Where: &WhereStruct{FieldName: "Where", Items: []WhereItem{
				&WhereColumnDirective{LHSColIdent: ColIdent{Name: "column_a"}, Predicate: NotNull},
				&WhereBoolTag{BoolAnd},
				&WhereColumnDirective{LHSColIdent: ColIdent{Name: "column_b"}, Predicate: IsNull},
				&WhereBoolTag{BoolOr},
				&WhereColumnDirective{LHSColIdent: ColIdent{Name: "column_c"}, Predicate: NotTrue},
				&WhereBoolTag{BoolAnd},
				&WhereColumnDirective{LHSColIdent: ColIdent{Name: "column_d"}, Predicate: IsTrue},
				&WhereBoolTag{BoolOr},
				&WhereColumnDirective{LHSColIdent: ColIdent{Name: "column_e"}, Predicate: NotFalse},
				&WhereBoolTag{BoolOr},
				&WhereColumnDirective{LHSColIdent: ColIdent{Name: "column_f"}, Predicate: IsFalse},
				&WhereBoolTag{BoolAnd},
				&WhereColumnDirective{LHSColIdent: ColIdent{Name: "column_g"}, Predicate: NotUnknown},
				&WhereBoolTag{BoolAnd},
				&WhereColumnDirective{LHSColIdent: ColIdent{Name: "column_h"}, Predicate: IsUnknown},
				&WhereBoolTag{BoolAnd},
				&WhereColumnDirective{LHSColIdent: ColIdent{Name: "column_i"}, Predicate: IsTrue},
			}},
		},
	}, {
		Name: "DeleteAnalysisTestOK11",
		want: &QueryStruct{
			TypeName: "DeleteAnalysisTestOK11",
			Kind:     QueryKindDelete,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a"},
				Type:      RelType{Base: TypeInfo{Kind: TypeKindStruct}},
			},
			Where: &WhereStruct{FieldName: "Where", Items: []WhereItem{
				&WhereStruct{FieldName: "x", Items: []WhereItem{
					&WhereStructField{
						Name:      "foo",
						Type:      TypeInfo{Kind: TypeKindInt},
						ColIdent:  ColIdent{Name: "column_foo"},
						Predicate: IsEQ,
					},
					&WhereBoolTag{BoolAnd},
					&WhereColumnDirective{LHSColIdent: ColIdent{Name: "column_a"}, Predicate: IsNull},
				}},
				&WhereBoolTag{BoolOr},
				&WhereStruct{FieldName: "y", Items: []WhereItem{
					&WhereColumnDirective{LHSColIdent: ColIdent{Name: "column_b"}, Predicate: NotTrue},
					&WhereBoolTag{BoolOr},
					&WhereStructField{
						Name:      "bar",
						Type:      TypeInfo{Kind: TypeKindString},
						ColIdent:  ColIdent{Name: "column_bar"},
						Predicate: IsEQ,
					},
					&WhereBoolTag{BoolAnd},
					&WhereStruct{FieldName: "z", Items: []WhereItem{
						&WhereStructField{
							Name:      "baz",
							Type:      TypeInfo{Kind: TypeKindBool},
							ColIdent:  ColIdent{Name: "column_baz"},
							Predicate: IsEQ,
						},
						&WhereBoolTag{BoolAnd},
						&WhereStructField{
							Name:      "quux",
							Type:      TypeInfo{Kind: TypeKindString},
							ColIdent:  ColIdent{Name: "column_quux"},
							Predicate: IsEQ,
						},
						&WhereBoolTag{BoolOr},
						&WhereColumnDirective{LHSColIdent: ColIdent{Name: "column_c"}, Predicate: IsTrue},
					}},
				}},
				&WhereBoolTag{BoolOr},
				&WhereColumnDirective{LHSColIdent: ColIdent{Name: "column_d"}, Predicate: NotFalse},
				&WhereBoolTag{BoolAnd},
				&WhereColumnDirective{LHSColIdent: ColIdent{Name: "column_e"}, Predicate: IsFalse},
				&WhereBoolTag{BoolAnd},
				&WhereStructField{
					Name:      "foo",
					Type:      TypeInfo{Kind: TypeKindInt},
					ColIdent:  ColIdent{Name: "column_foo"},
					Predicate: IsEQ,
				},
			}},
		},
	}, {
		Name: "DeleteAnalysisTestOK12",
		want: &QueryStruct{
			TypeName: "DeleteAnalysisTestOK12",
			Kind:     QueryKindDelete,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a"},
				Type:      RelType{Base: TypeInfo{Kind: TypeKindStruct}},
			},
			Where: &WhereStruct{FieldName: "Where", Items: []WhereItem{
				&WhereStructField{Name: "a", Type: TypeInfo{Kind: TypeKindInt}, ColIdent: ColIdent{Name: "column_a"}, Predicate: IsLT},
				&WhereBoolTag{BoolAnd},
				&WhereStructField{Name: "b", Type: TypeInfo{Kind: TypeKindInt}, ColIdent: ColIdent{Name: "column_b"}, Predicate: IsGT},
				&WhereBoolTag{BoolAnd},
				&WhereStructField{Name: "c", Type: TypeInfo{Kind: TypeKindInt}, ColIdent: ColIdent{Name: "column_c"}, Predicate: IsLTE},
				&WhereBoolTag{BoolAnd},
				&WhereStructField{Name: "d", Type: TypeInfo{Kind: TypeKindInt}, ColIdent: ColIdent{Name: "column_d"}, Predicate: IsGTE},
				&WhereBoolTag{BoolAnd},
				&WhereStructField{Name: "e", Type: TypeInfo{Kind: TypeKindInt}, ColIdent: ColIdent{Name: "column_e"}, Predicate: IsEQ},
				&WhereBoolTag{BoolAnd},
				&WhereStructField{Name: "f", Type: TypeInfo{Kind: TypeKindInt}, ColIdent: ColIdent{Name: "column_f"}, Predicate: NotEQ},
				&WhereBoolTag{BoolAnd},
				&WhereStructField{Name: "g", Type: TypeInfo{Kind: TypeKindInt}, ColIdent: ColIdent{Name: "column_g"}, Predicate: IsEQ},
			}},
		},
	}, {
		Name: "DeleteAnalysisTestOK13",
		want: &QueryStruct{
			TypeName: "DeleteAnalysisTestOK13",
			Kind:     QueryKindDelete,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a"},
				Type:      RelType{Base: TypeInfo{Kind: TypeKindStruct}},
			},
			Where: &WhereStruct{FieldName: "Where", Items: []WhereItem{
				&WhereColumnDirective{
					LHSColIdent: ColIdent{Name: "column_a"},
					RHSColIdent: ColIdent{Name: "column_b"},
					Predicate:   NotEQ,
				},
				&WhereBoolTag{BoolAnd},
				&WhereColumnDirective{
					LHSColIdent: ColIdent{Name: "column_c", Qualifier: "t"},
					RHSColIdent: ColIdent{Name: "column_d", Qualifier: "u"},
					Predicate:   IsEQ,
				},
				&WhereBoolTag{BoolAnd},
				&WhereColumnDirective{
					LHSColIdent: ColIdent{Name: "column_e", Qualifier: "t"},
					RHSLiteral:  "123",
					Predicate:   IsGT,
				},
				&WhereBoolTag{BoolAnd},
				&WhereColumnDirective{
					LHSColIdent: ColIdent{Name: "column_f", Qualifier: "t"},
					RHSLiteral:  "'active'",
					Predicate:   IsEQ,
				},
				&WhereBoolTag{BoolAnd},
				&WhereColumnDirective{
					LHSColIdent: ColIdent{Name: "column_g", Qualifier: "t"},
					RHSLiteral:  "true",
					Predicate:   NotEQ,
				},
			}},
		},
	}, {
		Name: "DeleteAnalysisTestOK14",
		want: &QueryStruct{
			TypeName: "DeleteAnalysisTestOK14",
			Kind:     QueryKindDelete,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a"},
				Type:      RelType{Base: TypeInfo{Kind: TypeKindStruct}},
			},
			Where: &WhereStruct{FieldName: "Where", Items: []WhereItem{
				&WhereBetweenStruct{
					FieldName:  "a",
					ColIdent:   ColIdent{Name: "column_a"},
					Predicate:  IsBetween,
					LowerBound: &BetweenStructField{Name: "x", Type: TypeInfo{Kind: TypeKindInt}},
					UpperBound: &BetweenStructField{Name: "y", Type: TypeInfo{Kind: TypeKindInt}},
				},
				&WhereBoolTag{BoolAnd},
				&WhereBetweenStruct{
					FieldName:  "b",
					ColIdent:   ColIdent{Name: "column_b"},
					Predicate:  IsBetweenSym,
					LowerBound: &BetweenColumnDirective{ColIdent{Name: "column_x"}},
					UpperBound: &BetweenColumnDirective{ColIdent{Name: "column_y"}},
				},
				&WhereBoolTag{BoolAnd},
				&WhereBetweenStruct{
					FieldName:  "c",
					ColIdent:   ColIdent{Name: "column_c"},
					Predicate:  NotBetweenSym,
					LowerBound: &BetweenColumnDirective{ColIdent{Name: "column_z"}},
					UpperBound: &BetweenStructField{Name: "z", Type: TypeInfo{Kind: TypeKindInt}},
				},
				&WhereBoolTag{BoolAnd},
				&WhereBetweenStruct{
					FieldName:  "d",
					ColIdent:   ColIdent{Name: "column_d"},
					Predicate:  NotBetween,
					LowerBound: &BetweenStructField{Name: "z", Type: TypeInfo{Kind: TypeKindInt}},
					UpperBound: &BetweenColumnDirective{ColIdent{Name: "column_z"}},
				},
				&WhereBoolTag{BoolAnd},
				&WhereBetweenStruct{
					FieldName:  "d2",
					ColIdent:   ColIdent{Name: "column_d"},
					Predicate:  NotBetween,
					LowerBound: &BetweenStructField{Name: "z", Type: TypeInfo{Kind: TypeKindInt}},
					UpperBound: &BetweenColumnDirective{ColIdent{Name: "column_z"}},
				},
			}},
		},
	}, {
		Name: "DeleteAnalysisTestOK_DistinctFrom",
		want: &QueryStruct{
			TypeName: "DeleteAnalysisTestOK_DistinctFrom",
			Kind:     QueryKindDelete,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a"},
				Type:      RelType{Base: TypeInfo{Kind: TypeKindStruct}},
			},
			Where: &WhereStruct{FieldName: "Where", Items: []WhereItem{
				&WhereStructField{
					Name:      "a",
					Type:      TypeInfo{Kind: TypeKindInt},
					ColIdent:  ColIdent{Name: "column_a"},
					Predicate: IsDistinct,
				},
				&WhereBoolTag{BoolAnd},
				&WhereStructField{
					Name:      "b",
					Type:      TypeInfo{Kind: TypeKindInt},
					ColIdent:  ColIdent{Name: "column_b"},
					Predicate: NotDistinct,
				},
				&WhereBoolTag{BoolAnd},
				&WhereColumnDirective{
					LHSColIdent: ColIdent{Name: "column_c"},
					RHSColIdent: ColIdent{Name: "column_x"},
					Predicate:   IsDistinct,
				},
				&WhereBoolTag{BoolAnd},
				&WhereColumnDirective{
					LHSColIdent: ColIdent{Name: "column_d"},
					RHSColIdent: ColIdent{Name: "column_y"},
					Predicate:   NotDistinct,
				},
			}},
		},
	}, {
		Name: "DeleteAnalysisTestOK_ArrayPredicate",
		want: &QueryStruct{
			TypeName: "DeleteAnalysisTestOK_ArrayPredicate",
			Kind:     QueryKindDelete,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a"},
				Type:      RelType{Base: TypeInfo{Kind: TypeKindStruct}},
			},
			Where: &WhereStruct{FieldName: "Where", Items: []WhereItem{
				&WhereStructField{
					Name: "a",
					Type: TypeInfo{
						Kind: TypeKindSlice,
						Elem: &TypeInfo{
							Kind: TypeKindInt,
						},
					},
					ColIdent:  ColIdent{Name: "column_a"},
					Predicate: IsIn,
				},
				&WhereBoolTag{BoolAnd},
				&WhereStructField{
					Name: "b",
					Type: TypeInfo{
						Kind: TypeKindArray,
						Elem: &TypeInfo{
							Kind: TypeKindInt,
						},
						ArrayLen: 5,
					},
					ColIdent:  ColIdent{Name: "column_b"},
					Predicate: NotIn,
				},
				&WhereBoolTag{BoolAnd},
				&WhereStructField{
					Name: "c",
					Type: TypeInfo{
						Kind: TypeKindSlice,
						Elem: &TypeInfo{
							Kind: TypeKindInt,
						},
					},
					ColIdent:   ColIdent{Name: "column_c"},
					Predicate:  IsEQ,
					Quantifier: QuantAny,
				},
				&WhereBoolTag{BoolAnd},
				&WhereStructField{
					Name: "d",
					Type: TypeInfo{
						Kind: TypeKindArray,
						Elem: &TypeInfo{
							Kind: TypeKindInt,
						},
						ArrayLen: 10,
					},
					ColIdent:   ColIdent{Name: "column_d"},
					Predicate:  IsGT,
					Quantifier: QuantSome,
				},
				&WhereBoolTag{BoolAnd},
				&WhereStructField{
					Name: "e",
					Type: TypeInfo{
						Kind: TypeKindSlice,
						Elem: &TypeInfo{
							Kind: TypeKindInt,
						},
					},
					ColIdent:   ColIdent{Name: "column_e"},
					Predicate:  IsLTE,
					Quantifier: QuantAll,
				},
			}},
		},
	}, {
		Name: "DeleteAnalysisTestOK_PatternMatching",
		want: &QueryStruct{
			TypeName: "DeleteAnalysisTestOK_PatternMatching",
			Kind:     QueryKindDelete,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a"},
				Type:      RelType{Base: TypeInfo{Kind: TypeKindStruct}},
			},
			Where: &WhereStruct{FieldName: "Where", Items: []WhereItem{
				&WhereStructField{
					Name:      "a",
					Type:      TypeInfo{Kind: TypeKindString},
					ColIdent:  ColIdent{Name: "column_a"},
					Predicate: IsLike,
				},
				&WhereBoolTag{BoolAnd},
				&WhereStructField{
					Name:      "b",
					Type:      TypeInfo{Kind: TypeKindString},
					ColIdent:  ColIdent{Name: "column_b"},
					Predicate: NotLike,
				},
				&WhereBoolTag{BoolAnd},
				&WhereStructField{
					Name:      "c",
					Type:      TypeInfo{Kind: TypeKindString},
					ColIdent:  ColIdent{Name: "column_c"},
					Predicate: IsSimilar,
				},
				&WhereBoolTag{BoolAnd},
				&WhereStructField{
					Name:      "d",
					Type:      TypeInfo{Kind: TypeKindString},
					ColIdent:  ColIdent{Name: "column_d"},
					Predicate: NotSimilar,
				},
				&WhereBoolTag{BoolAnd},
				&WhereStructField{
					Name:      "e",
					Type:      TypeInfo{Kind: TypeKindString},
					ColIdent:  ColIdent{Name: "column_e"},
					Predicate: IsMatch,
				},
				&WhereBoolTag{BoolAnd},
				&WhereStructField{
					Name:      "f",
					Type:      TypeInfo{Kind: TypeKindString},
					ColIdent:  ColIdent{Name: "column_f"},
					Predicate: IsMatchi,
				},
				&WhereBoolTag{BoolAnd},
				&WhereStructField{
					Name:      "g",
					Type:      TypeInfo{Kind: TypeKindString},
					ColIdent:  ColIdent{Name: "column_g"},
					Predicate: NotMatch,
				},
				&WhereBoolTag{BoolAnd},
				&WhereStructField{
					Name:      "h",
					Type:      TypeInfo{Kind: TypeKindString},
					ColIdent:  ColIdent{Name: "column_h"},
					Predicate: NotMatchi,
				},
			}},
		},
	}, {
		Name: "DeleteAnalysisTestOK_Using",
		want: &QueryStruct{
			TypeName: "DeleteAnalysisTestOK_Using",
			Kind:     QueryKindDelete,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a", Alias: "a"},
				Type:      RelType{Base: TypeInfo{Kind: TypeKindStruct}},
			},
			Join: &JoinStruct{
				FieldName: "Using",
				Relation:  &RelationDirective{RelIdent{Name: "relation_b", Alias: "b"}},
				Directives: []*JoinDirective{{
					JoinType: JoinTypeLeft,
					RelIdent: RelIdent{Name: "relation_c", Alias: "c"},
					TagItems: []JoinTagItem{
						&JoinConditionTagItem{
							LHSColIdent: ColIdent{Qualifier: "c", Name: "b_id"},
							RHSColIdent: ColIdent{Qualifier: "b", Name: "id"},
							Predicate:   IsEQ,
						},
					},
				}, {
					JoinType: JoinTypeRight,
					RelIdent: RelIdent{Name: "relation_d", Alias: "d"},
					TagItems: []JoinTagItem{
						&JoinConditionTagItem{
							LHSColIdent: ColIdent{Qualifier: "d", Name: "c_id"},
							RHSColIdent: ColIdent{Qualifier: "c", Name: "id"},
							Predicate:   IsEQ,
						},
						&JoinBoolTagItem{BoolOr},
						&JoinConditionTagItem{
							LHSColIdent: ColIdent{Qualifier: "d", Name: "num"},
							RHSColIdent: ColIdent{Qualifier: "b", Name: "num"},
							Predicate:   IsGT,
						},
					},
				}, {
					JoinType: JoinTypeFull,
					RelIdent: RelIdent{Name: "relation_e", Alias: "e"},
					TagItems: []JoinTagItem{
						&JoinConditionTagItem{
							LHSColIdent: ColIdent{Qualifier: "e", Name: "d_id"},
							RHSColIdent: ColIdent{Qualifier: "d", Name: "id"},
							Predicate:   IsEQ,
						},
						&JoinBoolTagItem{BoolAnd},
						&JoinConditionTagItem{
							LHSColIdent: ColIdent{Qualifier: "e", Name: "is_foo"},
							Predicate:   IsFalse,
						},
					},
				}, {
					JoinType: JoinTypeCross,
					RelIdent: RelIdent{Name: "relation_f", Alias: "f"},
				}},
			},
			Where: &WhereStruct{FieldName: "Where", Items: []WhereItem{
				&WhereColumnDirective{
					LHSColIdent: ColIdent{Name: "id", Qualifier: "a"},
					RHSColIdent: ColIdent{Name: "a_id", Qualifier: "d"},
					Predicate:   IsEQ,
				},
			}},
		},
	}, {
		Name: "UpdateAnalysisTestOK_From",
		want: &QueryStruct{
			TypeName: "UpdateAnalysisTestOK_From",
			Kind:     QueryKindUpdate,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a", Alias: "a"},
				Type:      RelType{Base: TypeInfo{Kind: TypeKindStruct}},
			},
			Join: &JoinStruct{
				FieldName: "From",
				Relation:  &RelationDirective{RelIdent{Name: "relation_b", Alias: "b"}},
				Directives: []*JoinDirective{{
					JoinType: JoinTypeLeft,
					RelIdent: RelIdent{Name: "relation_c", Alias: "c"},
					TagItems: []JoinTagItem{
						&JoinConditionTagItem{
							LHSColIdent: ColIdent{Qualifier: "c", Name: "b_id"},
							RHSColIdent: ColIdent{Qualifier: "b", Name: "id"},
							Predicate:   IsEQ,
						},
					},
				}, {
					JoinType: JoinTypeRight,
					RelIdent: RelIdent{Name: "relation_d", Alias: "d"},
					TagItems: []JoinTagItem{
						&JoinConditionTagItem{
							LHSColIdent: ColIdent{Qualifier: "d", Name: "c_id"},
							RHSColIdent: ColIdent{Qualifier: "c", Name: "id"},
							Predicate:   IsEQ,
						},
						&JoinBoolTagItem{BoolOr},
						&JoinConditionTagItem{
							LHSColIdent: ColIdent{Qualifier: "d", Name: "num"},
							RHSColIdent: ColIdent{Qualifier: "b", Name: "num"},
							Predicate:   IsGT,
						},
					},
				}, {
					JoinType: JoinTypeFull,
					RelIdent: RelIdent{Name: "relation_e", Alias: "e"},
					TagItems: []JoinTagItem{
						&JoinConditionTagItem{
							LHSColIdent: ColIdent{Qualifier: "e", Name: "d_id"},
							RHSColIdent: ColIdent{Qualifier: "d", Name: "id"},
							Predicate:   IsEQ,
						},
						&JoinBoolTagItem{BoolAnd},
						&JoinConditionTagItem{
							LHSColIdent: ColIdent{Qualifier: "e", Name: "is_foo"},
							Predicate:   IsFalse,
						},
					},
				}, {
					JoinType: JoinTypeCross,
					RelIdent: RelIdent{Name: "relation_f", Alias: "f"},
				}},
			},
			Where: &WhereStruct{FieldName: "Where", Items: []WhereItem{
				&WhereColumnDirective{
					LHSColIdent: ColIdent{Name: "id", Qualifier: "a"},
					RHSColIdent: ColIdent{Name: "a_id", Qualifier: "d"},
					Predicate:   IsEQ,
				},
			}},
		},
	}, {
		Name: "SelectAnalysisTestOK_Join",
		want: &QueryStruct{
			TypeName: "SelectAnalysisTestOK_Join",
			Kind:     QueryKindSelect,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a", Alias: "a"},
				Type:      RelType{Base: TypeInfo{Kind: TypeKindStruct}},
			},
			Join: &JoinStruct{
				FieldName: "Join",
				Directives: []*JoinDirective{{
					JoinType: JoinTypeLeft, RelIdent: RelIdent{Name: "relation_b", Alias: "b"},
					TagItems: []JoinTagItem{
						&JoinConditionTagItem{
							LHSColIdent: ColIdent{Qualifier: "b", Name: "a_id"},
							RHSColIdent: ColIdent{Qualifier: "a", Name: "id"},
							Predicate:   IsEQ,
						},
					},
				}, {
					JoinType: JoinTypeLeft,
					RelIdent: RelIdent{Name: "relation_c", Alias: "c"},
					TagItems: []JoinTagItem{
						&JoinConditionTagItem{
							LHSColIdent: ColIdent{Qualifier: "c", Name: "b_id"},
							RHSColIdent: ColIdent{Qualifier: "b", Name: "id"},
							Predicate:   IsEQ,
						},
					},
				}, {
					JoinType: JoinTypeRight,
					RelIdent: RelIdent{Name: "relation_d", Alias: "d"},
					TagItems: []JoinTagItem{
						&JoinConditionTagItem{
							LHSColIdent: ColIdent{Qualifier: "d", Name: "c_id"},
							RHSColIdent: ColIdent{Qualifier: "c", Name: "id"},
							Predicate:   IsEQ,
						},
						&JoinBoolTagItem{BoolOr},
						&JoinConditionTagItem{
							LHSColIdent: ColIdent{Qualifier: "d", Name: "num"},
							RHSColIdent: ColIdent{Qualifier: "b", Name: "num"},
							Predicate:   IsGT,
						},
					},
				}, {
					JoinType: JoinTypeFull,
					RelIdent: RelIdent{Name: "relation_e", Alias: "e"},
					TagItems: []JoinTagItem{
						&JoinConditionTagItem{
							LHSColIdent: ColIdent{Qualifier: "e", Name: "d_id"},
							RHSColIdent: ColIdent{Qualifier: "d", Name: "id"},
							Predicate:   IsEQ,
						},
						&JoinBoolTagItem{BoolAnd},
						&JoinConditionTagItem{
							LHSColIdent: ColIdent{Qualifier: "e", Name: "is_foo"},
							Predicate:   IsFalse,
						},
					},
				}, {
					JoinType: JoinTypeCross,
					RelIdent: RelIdent{Name: "relation_f", Alias: "f"},
				}},
			},
			Where: &WhereStruct{FieldName: "Where", Items: []WhereItem{
				&WhereColumnDirective{
					LHSColIdent: ColIdent{Name: "id", Qualifier: "a"},
					RHSColIdent: ColIdent{Name: "a_id", Qualifier: "d"},
					Predicate:   IsEQ,
				},
			}},
		},
	}, {
		Name: "UpdateAnalysisTestOK_All",
		want: &QueryStruct{
			TypeName: "UpdateAnalysisTestOK_All",
			Kind:     QueryKindUpdate,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a", Alias: "a"},
				Type:      RelType{Base: TypeInfo{Kind: TypeKindStruct}},
			},
			All: &AllDirective{},
		},
	}, {
		Name: "DeleteAnalysisTestOK_All",
		want: &QueryStruct{
			TypeName: "DeleteAnalysisTestOK_All",
			Kind:     QueryKindDelete,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a", Alias: "a"},
				Type:      RelType{Base: TypeInfo{Kind: TypeKindStruct}},
			},
			All: &AllDirective{},
		},
	}, {
		Name: "DeleteAnalysisTestOK_Return",
		want: &QueryStruct{
			TypeName: "DeleteAnalysisTestOK_Return",
			Kind:     QueryKindDelete,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a", Alias: "a"},
				Type:      dummyrecord,
			},
			Return: &ReturnDirective{ColIdentList{All: true}},
		},
	}, {
		Name: "InsertAnalysisTestOK_Return",
		want: &QueryStruct{
			TypeName: "InsertAnalysisTestOK_Return",
			Kind:     QueryKindInsert,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a", Alias: "a"},
				Type:      dummyrecord,
			},
			Return: &ReturnDirective{ColIdentList{Items: []ColIdent{
				{Qualifier: "a", Name: "foo"},
				{Qualifier: "a", Name: "bar"},
				{Qualifier: "a", Name: "baz"},
			}}},
		},
	}, {
		Name: "UpdateAnalysisTestOK_Return",
		want: &QueryStruct{
			TypeName: "UpdateAnalysisTestOK_Return",
			Kind:     QueryKindUpdate,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a", Alias: "a"},
				Type:      dummyrecord,
			},
			Return: &ReturnDirective{ColIdentList{Items: []ColIdent{
				{Qualifier: "a", Name: "foo"},
				{Qualifier: "a", Name: "bar"},
				{Qualifier: "a", Name: "baz"},
			}}},
		},
	}, {
		Name: "InsertAnalysisTestOK_Default",
		want: &QueryStruct{
			TypeName: "InsertAnalysisTestOK_Default",
			Kind:     QueryKindInsert,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a", Alias: "a"},
				Type:      RelType{Base: TypeInfo{Kind: TypeKindStruct}},
			},
			Default: &DefaultDirective{ColIdentList{All: true}},
		},
	}, {
		Name: "UpdateAnalysisTestOK_Default",
		want: &QueryStruct{
			TypeName: "UpdateAnalysisTestOK_Default",
			Kind:     QueryKindUpdate,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a", Alias: "a"},
				Type:      RelType{Base: TypeInfo{Kind: TypeKindStruct}},
			},
			Default: &DefaultDirective{ColIdentList{Items: []ColIdent{
				{Qualifier: "a", Name: "foo"},
				{Qualifier: "a", Name: "bar"},
				{Qualifier: "a", Name: "baz"},
			}}},
		},
	}, {
		Name: "InsertAnalysisTestOK_Force",
		want: &QueryStruct{
			TypeName: "InsertAnalysisTestOK_Force",
			Kind:     QueryKindInsert,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a", Alias: "a"},
				Type:      RelType{Base: TypeInfo{Kind: TypeKindStruct}},
			},
			Force: &ForceDirective{ColIdentList{All: true}},
		},
	}, {
		Name: "UpdateAnalysisTestOK_Force",
		want: &QueryStruct{
			TypeName: "UpdateAnalysisTestOK_Force",
			Kind:     QueryKindUpdate,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a", Alias: "a"},
				Type:      RelType{Base: TypeInfo{Kind: TypeKindStruct}},
			},
			Force: &ForceDirective{ColIdentList{Items: []ColIdent{
				{Qualifier: "a", Name: "foo"},
				{Qualifier: "a", Name: "bar"},
				{Qualifier: "a", Name: "baz"},
			}}},
		},
	}, {
		Name: "SelectAnalysisTestOK_ErrorHandler",
		want: &QueryStruct{
			TypeName: "SelectAnalysisTestOK_ErrorHandler",
			Kind:     QueryKindSelect,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a", Alias: "a"},
				Type:      RelType{Base: TypeInfo{Kind: TypeKindStruct}},
			},
			ErrorHandler: &ErrorHandlerField{Name: "eh"},
		},
	}, {
		Name: "InsertAnalysisTestOK_ErrorHandler",
		want: &QueryStruct{
			TypeName: "InsertAnalysisTestOK_ErrorHandler",
			Kind:     QueryKindInsert,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a", Alias: "a"},
				Type:      RelType{Base: TypeInfo{Kind: TypeKindStruct}},
			},
			ErrorHandler: &ErrorHandlerField{Name: "myerrorhandler"},
		},
	}, {
		Name: "SelectAnalysisTestOK_ErrorInfoHandler",
		want: &QueryStruct{
			TypeName: "SelectAnalysisTestOK_ErrorInfoHandler",
			Kind:     QueryKindSelect,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a", Alias: "a"},
				Type:      RelType{Base: TypeInfo{Kind: TypeKindStruct}},
			},
			ErrorHandler: &ErrorHandlerField{Name: "eh", IsInfo: true},
		},
	}, {
		Name: "InsertAnalysisTestOK_ErrorInfoHandler",
		want: &QueryStruct{
			TypeName: "InsertAnalysisTestOK_ErrorInfoHandler",
			Kind:     QueryKindInsert,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a", Alias: "a"},
				Type:      RelType{Base: TypeInfo{Kind: TypeKindStruct}},
			},
			ErrorHandler: &ErrorHandlerField{Name: "myerrorinfohandler", IsInfo: true},
		},
	}, {
		Name: "SelectAnalysisTestOK_Count",
		want: &QueryStruct{
			TypeName: "SelectAnalysisTestOK_Count",
			Kind:     QueryKindSelectCount,
			Rel: &RelField{
				FieldName: "Count",
				Id:        RelIdent{Name: "relation_a", Alias: "a"},
			},
		},
	}, {
		Name: "SelectAnalysisTestOK_Exists",
		want: &QueryStruct{
			TypeName: "SelectAnalysisTestOK_Exists",
			Kind:     QueryKindSelectExists,
			Rel: &RelField{
				FieldName: "Exists",
				Id:        RelIdent{Name: "relation_a", Alias: "a"},
			},
		},
	}, {
		Name: "SelectAnalysisTestOK_NotExists",
		want: &QueryStruct{
			TypeName: "SelectAnalysisTestOK_NotExists",
			Kind:     QueryKindSelectNotExists,
			Rel: &RelField{
				FieldName: "NotExists",
				Id:        RelIdent{Name: "relation_a", Alias: "a"},
			},
		},
	}, {
		Name: "DeleteAnalysisTestOK_Relation",
		want: &QueryStruct{
			TypeName: "DeleteAnalysisTestOK_Relation",
			Kind:     QueryKindDelete,
			Rel: &RelField{
				FieldName:   "_",
				Id:          RelIdent{Name: "relation_a", Alias: "a"},
				IsDirective: true,
			},
		},
	}, {
		Name: "SelectAnalysisTestOK_LimitDirective",
		want: &QueryStruct{
			TypeName: "SelectAnalysisTestOK_LimitDirective",
			Kind:     QueryKindSelect,
			Rel:      reldummyslice,
			Limit:    &LimitField{Value: 25},
		},
	}, {
		Name: "SelectAnalysisTestOK_LimitField",
		want: &QueryStruct{
			TypeName: "SelectAnalysisTestOK_LimitField",
			Kind:     QueryKindSelect,
			Rel:      reldummyslice,
			Limit:    &LimitField{Name: "Limit", Value: 10},
		},
	}, {
		Name: "SelectAnalysisTestOK_OffsetDirective",
		want: &QueryStruct{
			TypeName: "SelectAnalysisTestOK_OffsetDirective",
			Kind:     QueryKindSelect,
			Rel:      reldummyslice,
			Offset:   &OffsetField{Value: 25},
		},
	}, {
		Name: "SelectAnalysisTestOK_OffsetField",
		want: &QueryStruct{
			TypeName: "SelectAnalysisTestOK_OffsetField",
			Kind:     QueryKindSelect,
			Rel:      reldummyslice,
			Offset:   &OffsetField{Name: "Offset", Value: 10},
		},
	}, {
		Name: "SelectAnalysisTestOK_OrderByDirective",
		want: &QueryStruct{
			TypeName: "SelectAnalysisTestOK_OrderByDirective",
			Kind:     QueryKindSelect,
			Rel:      reldummyslice,
			OrderBy: &OrderByDirective{Items: []OrderByTagItem{
				{ColIdent: ColIdent{Qualifier: "a", Name: "foo"}, Direction: OrderAsc, Nulls: NullsFirst},
				{ColIdent: ColIdent{Qualifier: "a", Name: "bar"}, Direction: OrderDesc, Nulls: NullsFirst},
				{ColIdent: ColIdent{Qualifier: "a", Name: "baz"}, Direction: OrderDesc, Nulls: 0},
				{ColIdent: ColIdent{Qualifier: "a", Name: "quux"}, Direction: OrderAsc, Nulls: NullsLast},
			}},
		},
	}, {
		Name: "InsertAnalysisTestOK_OverrideDirective",
		want: &QueryStruct{
			TypeName: "InsertAnalysisTestOK_OverrideDirective",
			Kind:     QueryKindInsert,
			Rel:      reldummyslice,
			Override: &OverrideDirective{OverridingSystem},
		},
	}, {
		Name: "InsertAnalysisTestOK_OnConflict",
		want: &QueryStruct{
			TypeName:   "InsertAnalysisTestOK_OnConflict",
			Kind:       QueryKindInsert,
			Rel:        reldummyslice,
			OnConflict: &OnConflictStruct{FieldName: "OnConflict", Ignore: &IgnoreDirective{}},
		},
	}, {
		Name: "InsertAnalysisTestOK_OnConflictColumn",
		want: &QueryStruct{
			TypeName: "InsertAnalysisTestOK_OnConflictColumn",
			Kind:     QueryKindInsert,
			Rel:      reldummyslice,
			OnConflict: &OnConflictStruct{
				FieldName: "OnConflict",
				Column:    &ColumnDirective{[]ColIdent{{Name: "id", Qualifier: "a"}}},
				Ignore:    &IgnoreDirective{},
			},
		},
	}, {
		Name: "InsertAnalysisTestOK_OnConflictConstraint",
		want: &QueryStruct{
			TypeName: "InsertAnalysisTestOK_OnConflictConstraint",
			Kind:     QueryKindInsert,
			Rel:      reldummyslice,
			OnConflict: &OnConflictStruct{
				FieldName:  "OnConflict",
				Constraint: &ConstraintDirective{"relation_constraint_xyz"},
				Update: &UpdateDirective{ColIdentList{
					Items: []ColIdent{
						{Name: "foo", Qualifier: "a"},
						{Name: "bar", Qualifier: "a"},
						{Name: "baz", Qualifier: "a"},
					},
				}},
			},
		},
	}, {
		Name: "InsertAnalysisTestOK_OnConflictIndex",
		want: &QueryStruct{
			TypeName: "InsertAnalysisTestOK_OnConflictIndex",
			Kind:     QueryKindInsert,
			Rel:      reldummyslice,
			OnConflict: &OnConflictStruct{
				FieldName: "OnConflict",
				Index:     &IndexDirective{"relation_index_xyz"},
				Update:    &UpdateDirective{ColIdentList{All: true}},
			},
		},
	}, {
		Name: "DeleteAnalysisTestOK_ResultField",
		want: &QueryStruct{
			TypeName: "DeleteAnalysisTestOK_ResultField",
			Kind:     QueryKindDelete,
			Rel: &RelField{
				FieldName:   "_",
				Id:          RelIdent{Name: "relation_a", Alias: "a"},
				IsDirective: true,
			},
			Where: &WhereStruct{FieldName: "Where", Items: []WhereItem{
				&WhereColumnDirective{LHSColIdent: ColIdent{Name: "is_inactive", Qualifier: "a"}, Predicate: IsTrue},
			}},
			Result: &ResultField{
				FieldName: "Result",
				Type:      reldummyslice.Type,
			},
		},
	}, {
		Name: "DeleteAnalysisTestOK_RowsAffected",
		want: &QueryStruct{
			TypeName: "DeleteAnalysisTestOK_RowsAffected",
			Kind:     QueryKindDelete,
			Rel: &RelField{
				FieldName:   "_",
				Id:          RelIdent{Name: "relation_a", Alias: "a"},
				IsDirective: true,
			},
			Where: &WhereStruct{FieldName: "Where", Items: []WhereItem{
				&WhereColumnDirective{LHSColIdent: ColIdent{Name: "is_inactive", Qualifier: "a"}, Predicate: IsTrue},
			}},
			RowsAffected: &RowsAffectedField{
				Name:     "RowsAffected",
				TypeKind: TypeKindInt,
			},
		},
	}, {
		Name: "SelectAnalysisTestOK_FilterField",
		want: &QueryStruct{
			TypeName: "SelectAnalysisTestOK_FilterField",
			Kind:     QueryKindSelect,
			Rel:      reldummyslice,
			Filter:   &FilterField{"Filter"},
		},
	}, {
		Name: "SelectAnalysisTestOK_FieldTypesBasic",
		want: &QueryStruct{
			TypeName: "SelectAnalysisTestOK_FieldTypesBasic",
			Kind:     QueryKindSelect,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a", Alias: "a"},
				Type: RelType{
					Base: TypeInfo{Kind: TypeKindStruct},
					Fields: []*FieldInfo{{
						Name: "f1", Type: TypeInfo{Kind: TypeKindBool},
						ColIdent: ColIdent{Name: "c1"},
						Tag:      tagutil.Tag{"sql": {"c1"}},
					}, {
						Name: "f2", Type: TypeInfo{Kind: TypeKindUint8, IsByte: true},
						ColIdent: ColIdent{Name: "c2"},
						Tag:      tagutil.Tag{"sql": {"c2"}},
					}, {
						Name: "f3", Type: TypeInfo{Kind: TypeKindInt32, IsRune: true},
						ColIdent: ColIdent{Name: "c3"},
						Tag:      tagutil.Tag{"sql": {"c3"}},
					}, {
						Name: "f4", Type: TypeInfo{Kind: TypeKindInt8},
						ColIdent: ColIdent{Name: "c4"},
						Tag:      tagutil.Tag{"sql": {"c4"}},
					}, {
						Name: "f5", Type: TypeInfo{Kind: TypeKindInt16},
						ColIdent: ColIdent{Name: "c5"},
						Tag:      tagutil.Tag{"sql": {"c5"}},
					}, {
						Name: "f6", Type: TypeInfo{Kind: TypeKindInt32},
						ColIdent: ColIdent{Name: "c6"},
						Tag:      tagutil.Tag{"sql": {"c6"}},
					}, {
						Name: "f7", Type: TypeInfo{Kind: TypeKindInt64},
						ColIdent: ColIdent{Name: "c7"},
						Tag:      tagutil.Tag{"sql": {"c7"}},
					}, {
						Name: "f8", Type: TypeInfo{Kind: TypeKindInt},
						ColIdent: ColIdent{Name: "c8"},
						Tag:      tagutil.Tag{"sql": {"c8"}},
					}, {
						Name: "f9", Type: TypeInfo{Kind: TypeKindUint8},
						ColIdent: ColIdent{Name: "c9"},
						Tag:      tagutil.Tag{"sql": {"c9"}},
					}, {
						Name: "f10", Type: TypeInfo{Kind: TypeKindUint16},
						ColIdent: ColIdent{Name: "c10"},
						Tag:      tagutil.Tag{"sql": {"c10"}},
					}, {
						Name: "f11", Type: TypeInfo{Kind: TypeKindUint32},
						ColIdent: ColIdent{Name: "c11"},
						Tag:      tagutil.Tag{"sql": {"c11"}},
					}, {
						Name: "f12", Type: TypeInfo{Kind: TypeKindUint64},
						ColIdent: ColIdent{Name: "c12"},
						Tag:      tagutil.Tag{"sql": {"c12"}},
					}, {
						Name: "f13", Type: TypeInfo{Kind: TypeKindUint},
						ColIdent: ColIdent{Name: "c13"},
						Tag:      tagutil.Tag{"sql": {"c13"}},
					}, {
						Name: "f14", Type: TypeInfo{Kind: TypeKindUintptr},
						ColIdent: ColIdent{Name: "c14"},
						Tag:      tagutil.Tag{"sql": {"c14"}},
					}, {
						Name: "f15", Type: TypeInfo{Kind: TypeKindFloat32},
						ColIdent: ColIdent{Name: "c15"},
						Tag:      tagutil.Tag{"sql": {"c15"}},
					}, {
						Name: "f16", Type: TypeInfo{Kind: TypeKindFloat64},
						ColIdent: ColIdent{Name: "c16"},
						Tag:      tagutil.Tag{"sql": {"c16"}},
					}, {
						Name: "f17", Type: TypeInfo{Kind: TypeKindComplex64},
						ColIdent: ColIdent{Name: "c17"},
						Tag:      tagutil.Tag{"sql": {"c17"}},
					}, {
						Name: "f18", Type: TypeInfo{Kind: TypeKindComplex128},
						ColIdent: ColIdent{Name: "c18"},
						Tag:      tagutil.Tag{"sql": {"c18"}},
					}, {
						Name: "f19", Type: TypeInfo{Kind: TypeKindString},
						ColIdent: ColIdent{Name: "c19"},
						Tag:      tagutil.Tag{"sql": {"c19"}},
					}},
				},
			},
		},
	}, {
		Name: "SelectAnalysisTestOK_FieldTypesSlices",
		want: &QueryStruct{
			TypeName: "SelectAnalysisTestOK_FieldTypesSlices",
			Kind:     QueryKindSelect,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a", Alias: "a"},
				Type: RelType{
					Base: TypeInfo{Kind: TypeKindStruct},
					Fields: []*FieldInfo{{
						Name: "f1", Type: TypeInfo{
							Kind: TypeKindSlice,
							Elem: &TypeInfo{Kind: TypeKindBool},
						},
						ColIdent: ColIdent{Name: "c1"},
						Tag:      tagutil.Tag{"sql": {"c1"}},
					}, {
						Name: "f2", Type: TypeInfo{
							Kind: TypeKindSlice,
							Elem: &TypeInfo{Kind: TypeKindUint8, IsByte: true},
						},
						ColIdent: ColIdent{Name: "c2"},
						Tag:      tagutil.Tag{"sql": {"c2"}},
					}, {
						Name: "f3", Type: TypeInfo{
							Kind: TypeKindSlice,
							Elem: &TypeInfo{Kind: TypeKindInt32, IsRune: true},
						},
						ColIdent: ColIdent{Name: "c3"},
						Tag:      tagutil.Tag{"sql": {"c3"}},
					}, {
						Name: "f4", Type: TypeInfo{
							Name:       "HardwareAddr",
							Kind:       TypeKindSlice,
							PkgPath:    "net",
							PkgName:    "net",
							PkgLocal:   "net",
							IsImported: true,
							Elem:       &TypeInfo{Kind: TypeKindUint8, IsByte: true},
						},
						ColIdent: ColIdent{Name: "c4"},
						Tag:      tagutil.Tag{"sql": {"c4"}},
					}, {
						Name: "f5", Type: TypeInfo{
							Name:              "RawMessage",
							Kind:              TypeKindSlice,
							PkgPath:           "encoding/json",
							PkgName:           "json",
							PkgLocal:          "json",
							IsImported:        true,
							IsJSONMarshaler:   true,
							IsJSONUnmarshaler: true,
							Elem:              &TypeInfo{Kind: TypeKindUint8, IsByte: true},
						},
						ColIdent: ColIdent{Name: "c5"},
						Tag:      tagutil.Tag{"sql": {"c5"}},
					}, {
						Name: "f6", Type: TypeInfo{
							Kind: TypeKindSlice,
							Elem: &TypeInfo{
								Name:            "Marshaler",
								Kind:            TypeKindInterface,
								PkgPath:         "encoding/json",
								PkgName:         "json",
								PkgLocal:        "json",
								IsImported:      true,
								IsJSONMarshaler: true,
							},
						},
						ColIdent: ColIdent{Name: "c6"},
						Tag:      tagutil.Tag{"sql": {"c6"}},
					}, {
						Name: "f7", Type: TypeInfo{
							Kind: TypeKindSlice,
							Elem: &TypeInfo{
								Name:              "RawMessage",
								Kind:              TypeKindSlice,
								PkgPath:           "encoding/json",
								PkgName:           "json",
								PkgLocal:          "json",
								IsImported:        true,
								IsJSONMarshaler:   true,
								IsJSONUnmarshaler: true,
								Elem:              &TypeInfo{Kind: TypeKindUint8, IsByte: true},
							},
						},
						ColIdent: ColIdent{Name: "c7"},
						Tag:      tagutil.Tag{"sql": {"c7"}},
					}, {
						Name: "f8", Type: TypeInfo{
							Kind: TypeKindSlice,
							Elem: &TypeInfo{
								Kind: TypeKindSlice,
								Elem: &TypeInfo{Kind: TypeKindUint8, IsByte: true},
							},
						},
						ColIdent: ColIdent{Name: "c8"},
						Tag:      tagutil.Tag{"sql": {"c8"}},
					}, {
						Name: "f9", Type: TypeInfo{
							Kind: TypeKindSlice,
							Elem: &TypeInfo{
								Kind:     TypeKindArray,
								ArrayLen: 2,
								Elem: &TypeInfo{
									Kind:     TypeKindArray,
									ArrayLen: 2,
									Elem:     &TypeInfo{Kind: TypeKindFloat64},
								},
							},
						},
						ColIdent: ColIdent{Name: "c9"},
						Tag:      tagutil.Tag{"sql": {"c9"}},
					}, {
						Name: "f10", Type: TypeInfo{
							Kind: TypeKindSlice,
							Elem: &TypeInfo{
								Kind: TypeKindSlice,
								Elem: &TypeInfo{
									Kind:     TypeKindArray,
									ArrayLen: 2,
									Elem:     &TypeInfo{Kind: TypeKindFloat64},
								},
							},
						},
						ColIdent: ColIdent{Name: "c10"},
						Tag:      tagutil.Tag{"sql": {"c10"}},
					}, {
						Name: "f11", Type: TypeInfo{
							Kind: TypeKindMap,
							Key:  &TypeInfo{Kind: TypeKindString},
							Elem: &TypeInfo{
								Name:       "NullString",
								Kind:       TypeKindStruct,
								PkgPath:    "database/sql",
								PkgName:    "sql",
								PkgLocal:   "sql",
								IsImported: true,
								IsScanner:  true,
								IsValuer:   true,
							},
						},
						ColIdent: ColIdent{Name: "c11"},
						Tag:      tagutil.Tag{"sql": {"c11"}},
					}, {
						Name: "f12", Type: TypeInfo{
							Kind: TypeKindSlice,
							Elem: &TypeInfo{
								Kind: TypeKindMap,
								Key:  &TypeInfo{Kind: TypeKindString},
								Elem: &TypeInfo{
									Kind: TypeKindPtr,
									Elem: &TypeInfo{Kind: TypeKindString},
								},
							},
						},
						ColIdent: ColIdent{Name: "c12"},
						Tag:      tagutil.Tag{"sql": {"c12"}},
					}, {
						Name: "f13", Type: TypeInfo{
							Kind: TypeKindSlice,
							Elem: &TypeInfo{
								Kind:     TypeKindArray,
								ArrayLen: 2,
								Elem: &TypeInfo{
									Kind: TypeKindPtr,
									Elem: &TypeInfo{
										Name:              "Int",
										Kind:              TypeKindStruct,
										PkgPath:           "math/big",
										PkgName:           "big",
										PkgLocal:          "big",
										IsImported:        true,
										IsJSONMarshaler:   true,
										IsJSONUnmarshaler: true,
									},
								},
							},
						},
						ColIdent: ColIdent{Name: "c13"},
						Tag:      tagutil.Tag{"sql": {"c13"}},
					}},
				},
			},
		},
	}, {
		Name: "SelectAnalysisTestOK_FieldTypesInterfaces",
		want: &QueryStruct{
			TypeName: "SelectAnalysisTestOK_FieldTypesInterfaces",
			Kind:     QueryKindSelect,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a", Alias: "a"},
				Type: RelType{
					Base: TypeInfo{Kind: TypeKindStruct},
					Fields: []*FieldInfo{{
						Name: "f1", Type: TypeInfo{
							Name:            "Marshaler",
							Kind:            TypeKindInterface,
							PkgPath:         "encoding/json",
							PkgName:         "json",
							PkgLocal:        "json",
							IsImported:      true,
							IsJSONMarshaler: true,
						},
						ColIdent: ColIdent{Name: "c1"},
						Tag:      tagutil.Tag{"sql": {"c1"}},
					}, {
						Name: "f2", Type: TypeInfo{
							Name:              "Unmarshaler",
							Kind:              TypeKindInterface,
							PkgPath:           "encoding/json",
							PkgName:           "json",
							PkgLocal:          "json",
							IsImported:        true,
							IsJSONUnmarshaler: true,
						},
						ColIdent: ColIdent{Name: "c2"},
						Tag:      tagutil.Tag{"sql": {"c2"}},
					}, {
						Name: "f3", Type: TypeInfo{
							Kind: TypeKindInterface,
						},
						ColIdent: ColIdent{Name: "c3"},
						Tag:      tagutil.Tag{"sql": {"c3"}},
					}},
				},
			},
		},
	}, {
		Name: "SelectAnalysisTestOK_FieldTypesEmptyInterfaces",
		want: &QueryStruct{
			TypeName: "SelectAnalysisTestOK_FieldTypesEmptyInterfaces",
			Kind:     QueryKindSelect,
			Rel: &RelField{
				FieldName: "Rel",
				Id:        RelIdent{Name: "relation_a", Alias: "a"},
				Type: RelType{
					Base: TypeInfo{Kind: TypeKindStruct},
					Fields: []*FieldInfo{{
						Name: "f1", Type: TypeInfo{
							Kind:             TypeKindInterface,
							IsEmptyInterface: true,
						},
						ColIdent: ColIdent{Name: "c1"},
						Tag:      tagutil.Tag{"sql": {"c1"}},
					}, {
						Name: "f2", Type: TypeInfo{
							Kind: TypeKindPtr,
							Elem: &TypeInfo{
								Kind:             TypeKindInterface,
								IsEmptyInterface: true,
							},
						},
						ColIdent: ColIdent{Name: "c2"},
						Tag:      tagutil.Tag{"sql": {"c2"}},
					}, {
						Name: "f3", Type: TypeInfo{
							Name:             "donothing",
							Kind:             TypeKindInterface,
							PkgPath:          "path/to/test",
							PkgName:          "testdata",
							PkgLocal:         "testdata",
							IsEmptyInterface: true,
						},
						ColIdent: ColIdent{Name: "c3"},
						Tag:      tagutil.Tag{"sql": {"c3"}},
					}, {
						Name: "f4", Type: TypeInfo{
							Kind: TypeKindPtr,
							Elem: &TypeInfo{
								Name:             "donothing",
								Kind:             TypeKindInterface,
								PkgPath:          "path/to/test",
								PkgName:          "testdata",
								PkgLocal:         "testdata",
								IsEmptyInterface: true,
							},
						},
						ColIdent: ColIdent{Name: "c4"},
						Tag:      tagutil.Tag{"sql": {"c4"}},
					}},
				},
			},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			var got *QueryStruct
			ts, err := testRunAnalysis(tt.Name, t)
			if qs, ok := ts.(*QueryStruct); ok && qs != nil {
				got = qs
			}

			if e := compare.Compare(err, tt.err); e != nil {
				t.Errorf("%v - %#v %v", e, err, err)
			}
			if e := compare.Compare(got, tt.want); e != nil {
				t.Error(e)
			}

			if tt.printerr && err != nil {
				fmt.Println(err)
			}
		})
	}
}

func TestAnalysis_filterStruct(t *testing.T) {
	dummyrecord := RelType{
		Base: TypeInfo{
			Name:     "T",
			Kind:     TypeKindStruct,
			PkgPath:  "path/to/test",
			PkgName:  "testdata",
			PkgLocal: "testdata",
		},
		Fields: []*FieldInfo{{
			Type:       TypeInfo{Kind: TypeKindString},
			Name:       "F",
			IsExported: true,
			Tag:        tagutil.Tag{"sql": {"f"}},
			ColIdent:   ColIdent{Name: "f"},
		}},
	}

	tests := []struct {
		Name     string
		want     *FilterStruct
		err      error
		printerr bool
	}{{
		Name: "FilterAnalysisTestBAD_IllegalReturnDirective",
		err: &Error{
			Code:          errIllegalQueryField,
			FieldType:     "github.com/frk/gosql.Return",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TargetName:    "FilterAnalysisTestBAD_IllegalReturnDirective",
			PkgPath:       "path/to/test",
			FileLine:      104,
			FileName:      "../testdata/analysis_bad.go",
		},
	}, {
		Name: "FilterAnalysisTestBAD_BadTextSearchDirectiveColId",
		err: &Error{
			Code:          errBadColIdTagValue,
			PkgPath:       "path/to/test",
			TargetName:    "FilterAnalysisTestBAD_BadTextSearchDirectiveColId",
			BlockName:     "",
			FieldType:     "github.com/frk/gosql.TextSearch",
			FieldTypeKind: "struct",
			FieldName:     "_",
			TagString:     `sql:"123"`,
			TagError:      "123",
			FileName:      "../testdata/analysis_bad.go",
			FileLine:      723,
		},
	}, {
		Name: "FilterAnalysisTestOK_TextSearchDirective",
		want: &FilterStruct{
			TypeName: "FilterAnalysisTestOK_TextSearchDirective",
			Rel: &RelField{
				FieldName: "_",
				Id:        RelIdent{Name: "relation_a", Alias: "a"},
				Type:      dummyrecord,
			},
			TextSearch: &TextSearchDirective{ColIdent{Qualifier: "a", Name: "ts_document"}},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			var got *FilterStruct
			ts, err := testRunAnalysis(tt.Name, t)
			if fs, ok := ts.(*FilterStruct); ok && fs != nil {
				got = fs
			}

			if e := compare.Compare(err, tt.err); e != nil {
				t.Errorf("%v - %#v %v", e, err, err)
			}
			if e := compare.Compare(got, tt.want); e != nil {
				t.Error(e)
			}

			if tt.printerr && err != nil {
				fmt.Println(err)
			}
		})
	}
}

func TestTypeinfo_string(t *testing.T) {
	tests := []struct {
		name string
		want LiteralType
	}{
		{"f01", LiteralBool},
		{"f02", LiteralBool},
		{"f03", LiteralBoolSlice},
		{"f04", LiteralString},
		{"f05", LiteralString},
		{"f06", LiteralStringSlice},
		{"f07", LiteralStringSliceSlice},
		{"f08", LiteralStringMap},
		{"f09", LiteralStringPtrMap},
		{"f10", LiteralStringMapSlice},
		{"f11", LiteralStringPtrMapSlice},
		{"f12", LiteralByte},
		{"f13", LiteralByte},
		{"f14", LiteralByteSlice},
		{"f15", LiteralByteSliceSlice},
		{"f16", LiteralByteArray16},
		{"f17", LiteralByteArray16Slice},
		{"f18", LiteralRune},
		{"f19", LiteralRune},
		{"f20", LiteralRuneSlice},
		{"f21", LiteralRuneSliceSlice},
		{"f22", LiteralInt8},
		{"f23", LiteralInt8},
		{"f24", LiteralInt8Slice},
		{"f25", LiteralInt8SliceSlice},
		{"f26", LiteralInt16},
		{"f27", LiteralInt16},
		{"f28", LiteralInt16Slice},
		{"f29", LiteralInt16SliceSlice},
		{"f30", LiteralInt32},
		{"f31", LiteralInt32},
		{"f32", LiteralInt32Slice},
		{"f33", LiteralInt32Array2},
		{"f34", LiteralInt32Array2Slice},
		{"f35", LiteralInt64},
		{"f36", LiteralInt64},
		{"f37", LiteralInt64Slice},
		{"f38", LiteralInt64Array2},
		{"f39", LiteralInt64Array2Slice},
		{"f40", LiteralFloat32},
		{"f41", LiteralFloat32},
		{"f42", LiteralFloat32Slice},
		{"f43", LiteralFloat64},
		{"f44", LiteralFloat64},
		{"f45", LiteralFloat64Slice},
		{"f46", LiteralFloat64Array2},
		{"f47", LiteralFloat64Array2Slice},
		{"f48", LiteralFloat64Array2SliceSlice},
		{"f49", LiteralFloat64Array2Array2},
		{"f50", LiteralFloat64Array2Array2Slice},
		{"f51", LiteralFloat64Array3},
		{"f52", LiteralFloat64Array3Slice},
		{"f53", LiteralIPNet},
		{"f54", "[]*net.IPNet"},
		{"f55", LiteralTime},
		{"f56", LiteralTime},
		{"f57", LiteralTimeSlice},
		{"f58", "[]*time.Time"},
		{"f59", LiteralTimeArray2},
		{"f60", LiteralTimeArray2Slice},
		{"f61", LiteralHardwareAddr},
		{"f62", LiteralHardwareAddrSlice},
		{"f63", LiteralBigInt},
		{"f64", LiteralBigInt},
		{"f65", LiteralBigIntSlice},
		{"f66", "[]*big.Int"},
		{"f67", LiteralBigIntArray2},
		{"f68", "[2]*big.Int"},
		{"f69", "[][2]*big.Int"},
		{"f70", LiteralNullStringMap},
		{"f71", LiteralNullStringMapSlice},
		{"f72", "json.RawMessage"},
		{"f73", "[]json.RawMessage"},
	}

	ts, err := testRunAnalysis("SelectAnalysisTestOK_typeinfo_string", t)
	if err != nil {
		t.Error(err)
	} else if qs, ok := ts.(*QueryStruct); ok && qs != nil {
		fields := qs.Rel.Type.Fields
		for i := 0; i < len(fields); i++ {
			ff := fields[i]
			tt := tests[i]

			got := ff.Type.literal(false, true)
			if ff.Name != tt.name || got != tt.want {
				t.Errorf("got %s::%s, want %s::%s", ff.Name, got, tt.name, tt.want)
			}
		}
	}
}
