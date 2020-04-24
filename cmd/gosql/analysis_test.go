package main

import (
	"testing"

	"github.com/frk/compare"
	"github.com/frk/gosql/internal/errors"
	"github.com/frk/gosql/internal/testutil"
	"github.com/frk/tagutil"
)

var tdata = testutil.ParseTestdata("../../testdata")

func runAnalysis(name string, t *testing.T) (*targetInfo, error) {
	named := testutil.FindNamedType(name, tdata)
	if named == nil {
		// Stop the test if no type with the given name was found.
		t.Fatal(name, " not found")
		return nil, nil
	}

	a := &analyzer{named: named}
	if err := a.run(); err != nil {
		return nil, err
	}

	return a.targetInfo(), nil
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
			name:       "Time",
			kind:       typeKindStruct,
			pkgPath:    "time",
			pkgName:    "time",
			pkgLocal:   "time",
			isImported: true,
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
		err:  errors.NoDataFieldError,
	}, {
		name: "InsertAnalysisTestBAD3",
		err:  errors.BadDataFieldTypeError,
	}, {
		name: "DeleteAnalysisTestBAD_BadRelId",
		err:  errors.BadRelIdError,
	}, {
		name: "SelectAnalysisTestBAD_MultipleRelTags",
		err:  errors.MultipleDataFieldsError,
	}, {
		name: "DeleteAnalysisTestBAD_IllegalCountField",
		err:  errors.IllegalCountFieldError,
	}, {
		name: "UpdateAnalysisTestBAD_IllegalExistsField",
		err:  errors.IllegalExistsFieldError,
	}, {
		name: "InsertAnalysisTestBAD_IllegalNotExistsField",
		err:  errors.IllegalNotExistsFieldError,
	}, {
		name: "SelectAnalysisTestBAD_IllegalRelationDirective",
		err:  errors.IllegalRelationDirectiveError,
	}, {
		name: "SelectAnalysisTestBAD_UnnamedBaseStructType",
		err:  errors.BadDataFieldTypeError,
	}, {
		name: "SelectAnalysisTestBAD_IllegalAllDirective",
		err:  errors.IllegalAllDirectiveError,
	}, {
		name: "InsertAnalysisTestBAD_IllegalAllDirective",
		err:  errors.IllegalAllDirectiveError,
	}, {
		name: "UpdateAnalysisTestBAD_ConflictWhereProducer",
		err:  errors.ConflictWhereProducerError,
	}, {
		name: "DeleteAnalysisTestBAD_IllegalDefaultDirective",
		err:  errors.IllegalDefaultDirectiveError,
	}, {
		name: "UpdateAnalysisTestBAD_EmptyDefaultDirectiveCollist",
		err:  errors.EmptyColListError,
	}, {
		name: "SelectAnalysisTestBAD_IllegalForceDirective",
		err:  errors.IllegalForceDirectiveError,
	}, {
		name: "UpdateAnalysisTestBAD_BadForceDirectiveColId",
		err:  errors.BadColIdError,
	}, {
		name: "DeleteAnalysisTestBAD_ConflictResultProducer",
		err:  errors.ConflictResultProducerError,
	}, {
		name: "UpdateAnalysisTestBAD_EmptyReturnDirectiveCollist",
		err:  errors.EmptyColListError,
	}, {
		name: "InsertAnalysisTestBAD_IllegalLimitField",
		err:  errors.IllegalLimitFieldOrDirectiveError,
	}, {
		name: "UpdateAnalysisTestBAD_IllegalOffsetField",
		err:  errors.IllegalOffsetFieldOrDirectiveError,
	}, {
		name: "InsertAnalysisTestBAD_IllegalOrderByDirective",
		err:  errors.IllegalOrderByDirectiveError,
	}, {
		name: "DeleteAnalysisTestBAD_IllegalOverrideDirective",
		err:  errors.IllegalOverrideDirectiveError,
	}, {
		name: "SelectAnalysisTestBAD_IllegalTextSearchDirective",
		err:  errors.IllegalTextSearchDirectiveError,
	}, {
		name: "SelectAnalysisTestBAD_IllegalColumnDirective",
		err:  errors.IllegalCommandDirectiveError,
	}, {
		name: "InsertAnalysisTestBAD_IllegalWhereBlock",
		err:  errors.IllegalWhereBlockError,
	}, {
		name: "UpdateAnalysisTestBAD_IllegalJoinBlock",
		err:  errors.IllegalJoinBlockError,
	}, {
		name: "DeleteAnalysisTestBAD_IllegalFromBlock",
		err:  errors.IllegalFromBlockError,
	}, {
		name: "SelectAnalysisTestBAD_IllegalUsingBlock",
		err:  errors.IllegalUsingBlockError,
	}, {
		name: "UpdateAnalysisTestBAD_IllegalOnConflictBlock",
		err:  errors.IllegalOnConflictBlockError,
	}, {
		name: "SelectAnalysisTestBAD_IllegalResultField",
		err:  errors.IllegalResultFieldError,
	}, {
		name: "SelectAnalysisTestBAD_ConflictLimitProducer",
		err:  errors.ConflictLimitProducerError,
	}, {
		name: "SelectAnalysisTestBAD_ConflictOffsetProducer",
		err:  errors.ConflictOffsetProducerError,
	}, {
		name: "SelectAnalysisTestBAD_IllegalRowsAffectedField",
		err:  errors.IllegalRowsAffectedFieldError,
	}, {
		name: "InsertAnalysisTestBAD_IllegalFilterField",
		err:  errors.IllegalFilterFieldError,
	}, {
		name: "SelectAnalysisTestBAD_ConflictWhereProducer",
		err:  errors.ConflictWhereProducerError,
	}, {
		name: "DeleteAnalysisTestBAD_ConflictWhereProducer",
		err:  errors.ConflictErrorHandlerFieldError,
	}, {
		name: "SelectAnalysisTestBAD_IteratorWithTooManyMethods",
		err:  errors.BadIteratorTypeError,
	}, {
		name: "SelectAnalysisTestBAD_IteratorWithBadSignature",
		err:  errors.BadIteratorTypeError,
	}, {
		name: "SelectAnalysisTestBAD_IteratorWithBadSignatureIface",
		err:  errors.BadIteratorTypeError,
	}, {
		name: "SelectAnalysisTestBAD_IteratorWithUnexportedMethod",
		err:  errors.BadIteratorTypeError,
	}, {
		name: "SelectAnalysisTestBAD_IteratorWithUnnamedArgument",
		err:  errors.BadIteratorTypeError,
	}, {
		name: "SelectAnalysisTestBAD_IteratorWithNonStructArgument",
		err:  errors.BadIteratorTypeError,
	}, {
		name: "InsertAnalysisTestBAD_BadRelfiedlStructBaseType",
		err:  errors.BadDataFieldTypeError,
	}, {
		name: "UpdateAnalysisTestBAD_BadRelTypeFieldColId",
		err:  errors.BadColIdError,
	}, {
		name: "UpdateAnalysisTestBAD_ConflictWhereProducer2",
		err:  errors.ConflictWhereProducerError,
	}, {
		name: "DeleteAnalysisTestBAD_BadWhereBlockType",
		err:  errors.BadWhereBlockTypeError,
	}, {
		name: "SelectAnalysisTestBAD_BadBoolTagValue",
		err:  errors.BadBoolTagValueError,
	}, {
		name: "SelectAnalysisTestBAD_BadNestedWhereBlockType",
		err:  errors.BadWhereBlockTypeError,
	}, {
		name: "SelectAnalysisTestBAD_BadColumnExpressionLHS",
		err:  errors.BadColIdError,
	}, {
		name: "SelectAnalysisTestBAD_BadColumnPredicateCombo",
		err:  errors.BadPredicateComboError,
	}, {
		name: "DeleteAnalysisTestBAD_BadColumnExpressionLHS",
		err:  errors.BadColIdError,
	}, {
		name: "UpdateAnalysisTestBAD_BadUnaryOp",
		err:  errors.BadUnaryPredicateError,
	}, {
		name: "UpdateAnalysisTestBAD_ExtraQuantifier",
		err:  errors.ExtraQuantifierError,
	}, {
		name: "SelectAnalysisTestBAD_BadBetweenFieldType",
		err:  errors.BadBetweenTypeError,
	}, {
		name: "SelectAnalysisTestBAD_BadBetweenFieldType2",
		err:  errors.BadBetweenTypeError,
	}, {
		name: "SelectAnalysisTestBAD_BadBetweenArgColId",
		err:  errors.BadColIdError,
	}, {
		name: "SelectAnalysisTestBAD_NoBetweenXYArg",
		err:  errors.NoBetweenXYArgsError,
	}, {
		name: "SelectAnalysisTestBAD_BadBetweenColId",
		err:  errors.BadColIdError,
	}, {
		name: "DeleteAnalysisTestBAD_BadWhereFieldColId",
		err:  errors.BadColIdError,
	}, {
		name: "DeleteAnalysisTestBAD_BadWhereFieldPredicateCombo",
		err:  errors.BadPredicateComboError,
	}, {
		name: "DeleteAnalysisTestBAD_IllegalWhereFieldUnaryPredicate",
		err:  errors.IllegalUnaryPredicateError,
	}, {
		name: "UpdateAnalysisTestBAD_BadWhereFieldTypeForQuantifier",
		err:  errors.BadQuantifierFieldTypeError,
	}, {
		name: "SelectAnalysisTestBAD_BadJoinBlockType",
		err:  errors.BadJoinBlockTypeError,
	}, {
		name: "SelectAnalysisTestBAD_IllegalJoinBlockRelationDirective",
		err:  errors.IllegalJoinBlockRelationDirectiveError,
	}, {
		name: "DeleteAnalysisTestBAD_ConflictRelationDirective",
		err:  errors.ConflictJoinBlockRelationDirectiveError,
	}, {
		name: "UpdateAnalysisTestBAD_BadFromRelationRelId",
		err:  errors.BadRelIdError,
	}, {
		name: "SelectAnalysisTestBAD_BadJoinDirectiveRelId",
		err:  errors.BadRelIdError,
	}, {
		name: "SelectAnalysisTestBAD_BadJoinDirectiveExpressionColId",
		err:  errors.BadColIdError,
	}, {
		name: "SelectAnalysisTestBAD_BadJoinDirectiveExpressionPredicate",
		err:  errors.BadUnaryPredicateError,
	}, {
		name: "SelectAnalysisTestBAD_BadJoinDirectiveExpressionExtraQuantifier",
		err:  errors.ExtraQuantifierError,
	}, {
		name: "SelectAnalysisTestBAD_BadJoinDirectiveExpressionPredicateCombo",
		err:  errors.BadPredicateComboError,
	}, {
		name: "DeleteAnalysisTestBAD_IllegalJoinBlockDirective",
		err:  errors.IllegalJoinBlockDirectiveError,
	}, {
		name: "InsertAnalysisTestBAD_BadOnConflictBlockType",
		err:  errors.BadOnConflictBlockTypeError,
	}, {
		name: "InsertAnalysisTestBAD_ConflictOnConflictBlockTargetProducer",
		err:  errors.ConflictOnConflictBlockTargetProducerError,
	}, {
		name: "InsertAnalysisTestBAD_ConflictOnConflictBlockTargetProducer2",
		err:  errors.ConflictOnConflictBlockTargetProducerError,
	}, {
		name: "InsertAnalysisTestBAD_ConflictOnConflictBlockTargetProducer3",
		err:  errors.ConflictOnConflictBlockTargetProducerError,
	}, {
		name: "InsertAnalysisTestBAD_ConflictOnConflictBlockActionProducer",
		err:  errors.ConflictOnConflictBlockActionProducerError,
	}, {
		name: "InsertAnalysisTestBAD_ConflictOnConflictBlockActionProducer2",
		err:  errors.ConflictOnConflictBlockActionProducerError,
	}, {
		name: "InsertAnalysisTestBAD_BadOnConflictColumnTargetValue",
		err:  errors.BadColIdError,
	}, {
		name: "InsertAnalysisTestBAD_BadOnConflictIndexTargetIdent",
		err:  errors.BadIndexIdentifierValueError,
	}, {
		name: "InsertAnalysisTestBAD_BadOnConflictConstraintTargetIdent",
		err:  errors.BadConstraintIdentifierValueError,
	}, {
		name: "InsertAnalysisTestBAD_BadOnConflictUpdateActionCollist",
		err:  errors.BadColIdError,
	}, {
		name: "InsertAnalysisTestBAD_IllegalOnConflictDirective",
		err:  errors.IllegalOnConflictBlockDirectiveError,
	}, {
		name: "InsertAnalysisTestBAD_NoOnConflictTarget",
		err:  errors.NoOnConflictTargetError,
	}, {
		name: "SelectAnalysisTestBAD_BadLimitFieldType",
		err:  errors.BadLimitTypeError,
	}, {
		name: "SelectAnalysisTestBAD_NoLimitDirectiveValue",
		err:  errors.NoLimitDirectiveValueError,
	}, {
		name: "SelectAnalysisTestBAD_BadLimitDirectiveValue",
		err:  errors.BadLimitValueError,
	}, {
		name: "SelectAnalysisTestBAD_BadOffsetFieldType",
		err:  errors.BadOffsetTypeError,
	}, {
		name: "SelectAnalysisTestBAD_NoOffsetDirectiveValue",
		err:  errors.NoOffsetDirectiveValueError,
	}, {
		name: "SelectAnalysisTestBAD_BadOffsetDirectiveValue",
		err:  errors.BadOffsetValueError,
	}, {
		name: "SelectAnalysisTestBAD_EmptyOrderByDirectiveCollist",
		err:  errors.EmptyOrderByListError,
	}, {
		name: "SelectAnalysisTestBAD_BadOrderByDirectiveNullsOrderValue",
		err:  errors.BadNullsOrderOptionValueError,
	}, {
		name: "SelectAnalysisTestBAD_BadOrderByDirectiveCollist",
		err:  errors.BadColIdError,
	}, {
		name: "InsertAnalysisTestBAD_BadOverrideDirectiveKindValue",
		err:  errors.BadOverrideKindValueError,
	}, {
		name: "UpdateAnalysisTestBAD_ConflictResultProducer",
		err:  errors.ConflictResultProducerError,
	}, {
		name: "UpdateAnalysisTestBAD_BadResultFieldType",
		err:  errors.BadDataFieldTypeError,
	}, {
		name: "DeleteAnalysisTestBAD_ConflictResultProducer2",
		err:  errors.ConflictResultProducerError,
	}, {
		name: "DeleteAnalysisTestBAD_BadRowsAffecteFieldType",
		err:  errors.BadRowsAffectedTypeError,
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
			joinBlock: &joinBlock{relId: relId{name: "relation_b", alias: "b"}, items: []*joinItem{
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
			joinBlock: &joinBlock{relId: relId{name: "relation_b", alias: "b"}, items: []*joinItem{
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
			joinBlock: &joinBlock{items: []*joinItem{
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
				{qual: "a", name: "baz"}}},
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
				{qual: "a", name: "baz"}}},
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
				{qual: "a", name: "baz"}}},
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
				{qual: "a", name: "baz"}}},
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
							name:       "RawMessage",
							kind:       typeKindSlice,
							pkgPath:    "encoding/json",
							pkgName:    "json",
							pkgLocal:   "json",
							isImported: true,
							elem:       &typeInfo{kind: typeKindUint8, isByte: true},
						},
						colId: colId{name: "c5"},
						tag:   tagutil.Tag{"sql": {"c5"}},
					}, {
						name: "f6", typ: typeInfo{
							kind: typeKindSlice,
							elem: &typeInfo{
								name:       "Marshaler",
								kind:       typeKindInterface,
								pkgPath:    "encoding/json",
								pkgName:    "json",
								pkgLocal:   "json",
								isImported: true,
							},
						},
						colId: colId{name: "c6"},
						tag:   tagutil.Tag{"sql": {"c6"}},
					}, {
						name: "f7", typ: typeInfo{
							kind: typeKindSlice,
							elem: &typeInfo{
								name:       "RawMessage",
								kind:       typeKindSlice,
								pkgPath:    "encoding/json",
								pkgName:    "json",
								pkgLocal:   "json",
								isImported: true,
								elem:       &typeInfo{kind: typeKindUint8, isByte: true},
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
										name:       "Int",
										kind:       typeKindStruct,
										pkgPath:    "math/big",
										pkgName:    "big",
										pkgLocal:   "big",
										isImported: true,
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
							name:       "Marshaler",
							kind:       typeKindInterface,
							pkgPath:    "encoding/json",
							pkgName:    "json",
							pkgLocal:   "json",
							isImported: true,
						},
						colId: colId{name: "c1"},
						tag:   tagutil.Tag{"sql": {"c1"}},
					}, {
						name: "f2", typ: typeInfo{
							name:       "Unmarshaler",
							kind:       typeKindInterface,
							pkgPath:    "encoding/json",
							pkgName:    "json",
							pkgLocal:   "json",
							isImported: true,
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
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *queryStruct
			ti, err := runAnalysis(tt.name, t)
			if ti != nil {
				got = ti.query
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
		err:  errors.IllegalCommandDirectiveError,
	}, {
		name: "FilterAnalysisTestBAD_BadTextSearchDirectiveColId",
		err:  errors.BadColIdError,
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
			ti, err := runAnalysis(tt.name, t)
			if ti != nil {
				got = ti.filter
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

	ti, err := runAnalysis("SelectAnalysisTestOK_typeinfo_string", t)
	if err != nil {
		t.Error(err)
	}
	fields := ti.query.dataField.data.fields
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
