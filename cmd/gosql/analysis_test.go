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

	ti := new(targetInfo)
	if err := analyze(named, ti); err != nil {
		return nil, err
	}

	return ti, nil
}

func TestAnalysis_queryStruct(t *testing.T) {
	// for reuse, analyzed common.User typeInfo
	commonUserTypeinfo := typeInfo{
		name:       "User",
		kind:       kindStruct,
		pkgPath:    "github.com/frk/gosql/testdata/common",
		pkgName:    "common",
		pkgLocal:   "common",
		isImported: true,
	}

	commonUserFields := []*fieldInfo{{
		name:       "Id",
		typ:        typeInfo{kind: kindInt},
		isExported: true,
		colId:      colId{name: "id"},
		tag:        tagutil.Tag{"sql": {"id"}},
	}, {
		name:       "Email",
		typ:        typeInfo{kind: kindString},
		isExported: true,
		colId:      colId{name: "email"},
		tag:        tagutil.Tag{"sql": {"email"}},
	}, {
		name:       "FullName",
		typ:        typeInfo{kind: kindString},
		isExported: true,
		colId:      colId{name: "full_name"},
		tag:        tagutil.Tag{"sql": {"full_name"}},
	}, {
		name: "CreatedAt",
		typ: typeInfo{
			name:              "Time",
			kind:              kindStruct,
			pkgPath:           "time",
			pkgName:           "time",
			pkgLocal:          "time",
			isImported:        true,
			isTime:            true,
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
				kind:     kindStruct,
				pkgPath:  "path/to/test",
				pkgName:  "testdata",
				pkgLocal: "testdata",
			},
			isSlice: true,
			fields: []*fieldInfo{{
				typ:        typeInfo{kind: kindString},
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
			kind:     kindStruct,
			pkgPath:  "path/to/test",
			pkgName:  "testdata",
			pkgLocal: "testdata",
		},
		fields: []*fieldInfo{{
			typ:        typeInfo{kind: kindString},
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
						kind: kindStruct,
					},
					fields: []*fieldInfo{{
						name:       "Name3",
						typ:        typeInfo{kind: kindString},
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
						kind: kindStruct,
					},
					fields: []*fieldInfo{{
						name:   "a",
						typ:    typeInfo{kind: kindInt},
						colId:  colId{name: "a"},
						tag:    tagutil.Tag{"sql": {"a", "pk"}},
						isPKey: true,
					}, {
						name:      "b",
						typ:       typeInfo{kind: kindInt},
						colId:     colId{name: "b"},
						tag:       tagutil.Tag{"sql": {"b", "nullempty"}},
						nullEmpty: true,
					}, {
						name:     "c",
						typ:      typeInfo{kind: kindInt},
						colId:    colId{name: "c"},
						tag:      tagutil.Tag{"sql": {"c", "ro", "json"}},
						readOnly: true,
						useJSON:  true,
					}, {
						name:      "d",
						typ:       typeInfo{kind: kindInt},
						colId:     colId{name: "d"},
						tag:       tagutil.Tag{"sql": {"d", "wo"}},
						writeOnly: true,
					}, {
						name:   "e",
						typ:    typeInfo{kind: kindInt},
						colId:  colId{name: "e"},
						tag:    tagutil.Tag{"sql": {"e", "add"}},
						useAdd: true,
					}, {
						name:        "f",
						typ:         typeInfo{kind: kindInt},
						colId:       colId{name: "f"},
						tag:         tagutil.Tag{"sql": {"f", "coalesce"}},
						useCoalesce: true,
					}, {
						name:          "g",
						typ:           typeInfo{kind: kindInt},
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
						kind: kindStruct,
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
						typ:        typeInfo{kind: kindString},
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
						typ:        typeInfo{kind: kindString},
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
				data:  dataType{typeInfo: typeInfo{kind: kindStruct}},
			},
			whereBlock: &whereBlock{name: "Where", items: []*predicateItem{{
				pred: &predicateField{
					name:  "ID",
					typ:   typeInfo{kind: kindInt},
					colId: colId{name: "id"},
					kind:  isEQ,
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
				data:  dataType{typeInfo: typeInfo{kind: kindStruct}},
			},
			whereBlock: &whereBlock{name: "Where", items: []*predicateItem{
				{pred: &predicateColumn{colId: colId{name: "column_a"}, kind: notNull}},
				{bool: boolAnd, pred: &predicateColumn{colId: colId{name: "column_b"}, kind: isNull}},
				{bool: boolOr, pred: &predicateColumn{colId: colId{name: "column_c"}, kind: notTrue}},
				{bool: boolAnd, pred: &predicateColumn{colId: colId{name: "column_d"}, kind: isTrue}},
				{bool: boolOr, pred: &predicateColumn{colId: colId{name: "column_e"}, kind: notFalse}},
				{bool: boolOr, pred: &predicateColumn{colId: colId{name: "column_f"}, kind: isFalse}},
				{bool: boolAnd, pred: &predicateColumn{colId: colId{name: "column_g"}, kind: notUnknown}},
				{bool: boolAnd, pred: &predicateColumn{colId: colId{name: "column_h"}, kind: isUnknown}},
				{bool: boolAnd, pred: &predicateColumn{colId: colId{name: "column_i"}, kind: isTrue}},
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
				data:  dataType{typeInfo: typeInfo{kind: kindStruct}},
			},
			whereBlock: &whereBlock{name: "Where", items: []*predicateItem{
				{pred: &predicateNested{name: "x", items: []*predicateItem{
					{pred: &predicateField{
						name:  "foo",
						typ:   typeInfo{kind: kindInt},
						colId: colId{name: "column_foo"},
						kind:  isEQ,
					}},
					{bool: boolAnd, pred: &predicateColumn{colId: colId{name: "column_a"}, kind: isNull}},
				}}},
				{bool: boolOr, pred: &predicateNested{name: "y", items: []*predicateItem{
					{pred: &predicateColumn{colId: colId{name: "column_b"}, kind: notTrue}},
					{bool: boolOr, pred: &predicateField{
						name:  "bar",
						typ:   typeInfo{kind: kindString},
						colId: colId{name: "column_bar"},
						kind:  isEQ,
					}},
					{bool: boolAnd, pred: &predicateNested{name: "z", items: []*predicateItem{
						{pred: &predicateField{
							name:  "baz",
							typ:   typeInfo{kind: kindBool},
							colId: colId{name: "column_baz"},
							kind:  isEQ,
						}},
						{bool: boolAnd, pred: &predicateField{
							name:  "quux",
							typ:   typeInfo{kind: kindString},
							colId: colId{name: "column_quux"},
							kind:  isEQ,
						}},
						{bool: boolOr, pred: &predicateColumn{colId: colId{name: "column_c"}, kind: isTrue}},
					}}},
				}}},
				{bool: boolOr, pred: &predicateColumn{colId: colId{name: "column_d"}, kind: notFalse}},
				{bool: boolAnd, pred: &predicateColumn{colId: colId{name: "column_e"}, kind: isFalse}},
				{bool: boolAnd, pred: &predicateField{
					name:  "foo",
					typ:   typeInfo{kind: kindInt},
					colId: colId{name: "column_foo"},
					kind:  isEQ,
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
				data:  dataType{typeInfo: typeInfo{kind: kindStruct}},
			},
			whereBlock: &whereBlock{name: "Where", items: []*predicateItem{
				{pred: &predicateField{name: "a", typ: typeInfo{kind: kindInt}, colId: colId{name: "column_a"}, kind: isLT}},
				{bool: boolAnd, pred: &predicateField{name: "b", typ: typeInfo{kind: kindInt}, colId: colId{name: "column_b"}, kind: isGT}},
				{bool: boolAnd, pred: &predicateField{name: "c", typ: typeInfo{kind: kindInt}, colId: colId{name: "column_c"}, kind: isLTE}},
				{bool: boolAnd, pred: &predicateField{name: "d", typ: typeInfo{kind: kindInt}, colId: colId{name: "column_d"}, kind: isGTE}},
				{bool: boolAnd, pred: &predicateField{name: "e", typ: typeInfo{kind: kindInt}, colId: colId{name: "column_e"}, kind: isEQ}},
				{bool: boolAnd, pred: &predicateField{name: "f", typ: typeInfo{kind: kindInt}, colId: colId{name: "column_f"}, kind: notEQ}},
				{bool: boolAnd, pred: &predicateField{name: "g", typ: typeInfo{kind: kindInt}, colId: colId{name: "column_g"}, kind: isEQ}},
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
				data:  dataType{typeInfo: typeInfo{kind: kindStruct}},
			},
			whereBlock: &whereBlock{name: "Where", items: []*predicateItem{
				{pred: &predicateColumn{colId: colId{name: "column_a"}, kind: notEQ, colId2: colId{name: "column_b"}}},
				{bool: boolAnd, pred: &predicateColumn{colId: colId{qual: "t", name: "column_c"}, kind: isEQ, colId2: colId{qual: "u", name: "column_d"}}},
				{bool: boolAnd, pred: &predicateColumn{colId: colId{qual: "t", name: "column_e"}, kind: isGT, literal: "123"}},
				{bool: boolAnd, pred: &predicateColumn{colId: colId{qual: "t", name: "column_f"}, kind: isEQ, literal: "'active'"}},
				{bool: boolAnd, pred: &predicateColumn{colId: colId{qual: "t", name: "column_g"}, kind: notEQ, literal: "true"}},
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
				data:  dataType{typeInfo: typeInfo{kind: kindStruct}},
			},
			whereBlock: &whereBlock{name: "Where", items: []*predicateItem{
				{pred: &predicateBetween{
					name:  "a",
					colId: colId{name: "column_a"},
					kind:  isBetween,
					x:     &fieldDatum{name: "x", typ: typeInfo{kind: kindInt}},
					y:     &fieldDatum{name: "y", typ: typeInfo{kind: kindInt}},
				}},
				{bool: boolAnd, pred: &predicateBetween{
					name:  "b",
					colId: colId{name: "column_b"},
					kind:  isBetweenSym,
					x:     colId{name: "column_x"},
					y:     colId{name: "column_y"},
				}},
				{bool: boolAnd, pred: &predicateBetween{
					name:  "c",
					colId: colId{name: "column_c"},
					kind:  notBetweenSym,
					x:     colId{name: "column_z"},
					y:     &fieldDatum{name: "z", typ: typeInfo{kind: kindInt}},
				}},
				{bool: boolAnd, pred: &predicateBetween{
					name:  "d",
					colId: colId{name: "column_d"},
					kind:  notBetween,
					x:     &fieldDatum{name: "z", typ: typeInfo{kind: kindInt}},
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
				data:  dataType{typeInfo: typeInfo{kind: kindStruct}},
			},
			whereBlock: &whereBlock{name: "Where", items: []*predicateItem{
				{pred: &predicateField{
					name:  "a",
					typ:   typeInfo{kind: kindInt},
					colId: colId{name: "column_a"},
					kind:  isDistinct,
				}},
				{bool: boolAnd, pred: &predicateField{
					name:  "b",
					typ:   typeInfo{kind: kindInt},
					colId: colId{name: "column_b"},
					kind:  notDistinct,
				}},
				{bool: boolAnd, pred: &predicateColumn{colId: colId{name: "column_c"}, kind: isDistinct, colId2: colId{name: "column_x"}}},
				{bool: boolAnd, pred: &predicateColumn{colId: colId{name: "column_d"}, kind: notDistinct, colId2: colId{name: "column_y"}}},
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
				data:  dataType{typeInfo: typeInfo{kind: kindStruct}},
			},
			whereBlock: &whereBlock{name: "Where", items: []*predicateItem{
				{pred: &predicateField{
					name: "a",
					typ: typeInfo{
						kind: kindSlice,
						elem: &typeInfo{
							kind: kindInt,
						},
					},
					colId: colId{name: "column_a"},
					kind:  isIn,
				}},
				{bool: boolAnd, pred: &predicateField{
					name: "b",
					typ: typeInfo{
						kind: kindArray,
						elem: &typeInfo{
							kind: kindInt,
						},
						arrayLen: 5,
					},
					colId: colId{name: "column_b"},
					kind:  notIn,
				}},
				{bool: boolAnd, pred: &predicateField{
					name: "c",
					typ: typeInfo{
						kind: kindSlice,
						elem: &typeInfo{
							kind: kindInt,
						},
					},
					colId: colId{name: "column_c"},
					kind:  isEQ,
					qua:   quantAny,
				}},
				{bool: boolAnd, pred: &predicateField{
					name: "d",
					typ: typeInfo{
						kind: kindArray,
						elem: &typeInfo{
							kind: kindInt,
						},
						arrayLen: 10,
					},
					colId: colId{name: "column_d"},
					kind:  isGT,
					qua:   quantSome,
				}},
				{bool: boolAnd, pred: &predicateField{
					name: "e",
					typ: typeInfo{
						kind: kindSlice,
						elem: &typeInfo{
							kind: kindInt,
						},
					},
					colId: colId{name: "column_e"},
					kind:  isLTE,
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
				data:  dataType{typeInfo: typeInfo{kind: kindStruct}},
			},
			whereBlock: &whereBlock{name: "Where", items: []*predicateItem{
				{pred: &predicateField{
					name:  "a",
					typ:   typeInfo{kind: kindString},
					colId: colId{name: "column_a"},
					kind:  isLike,
				}},
				{bool: boolAnd, pred: &predicateField{
					name:  "b",
					typ:   typeInfo{kind: kindString},
					colId: colId{name: "column_b"},
					kind:  notLike,
				}},
				{bool: boolAnd, pred: &predicateField{
					name:  "c",
					typ:   typeInfo{kind: kindString},
					colId: colId{name: "column_c"},
					kind:  isSimilar,
				}},
				{bool: boolAnd, pred: &predicateField{
					name:  "d",
					typ:   typeInfo{kind: kindString},
					colId: colId{name: "column_d"},
					kind:  notSimilar,
				}},
				{bool: boolAnd, pred: &predicateField{
					name:  "e",
					typ:   typeInfo{kind: kindString},
					colId: colId{name: "column_e"},
					kind:  isMatch,
				}},
				{bool: boolAnd, pred: &predicateField{
					name:  "f",
					typ:   typeInfo{kind: kindString},
					colId: colId{name: "column_f"},
					kind:  isMatchi,
				}},
				{bool: boolAnd, pred: &predicateField{
					name:  "g",
					typ:   typeInfo{kind: kindString},
					colId: colId{name: "column_g"},
					kind:  notMatch,
				}},
				{bool: boolAnd, pred: &predicateField{
					name:  "h",
					typ:   typeInfo{kind: kindString},
					colId: colId{name: "column_h"},
					kind:  notMatchi,
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
				data:  dataType{typeInfo: typeInfo{kind: kindStruct}},
			},
			joinBlock: &joinBlock{relId: relId{name: "relation_b", alias: "b"}, items: []*joinItem{
				{joinType: joinLeft, relId: relId{name: "relation_c", alias: "c"}, predicates: []*predicateItem{{
					pred: &predicateColumn{
						colId:  colId{qual: "c", name: "b_id"},
						colId2: colId{qual: "b", name: "id"},
						kind:   isEQ,
					}}}},
				{joinType: joinRight, relId: relId{name: "relation_d", alias: "d"}, predicates: []*predicateItem{{
					pred: &predicateColumn{
						colId:  colId{qual: "d", name: "c_id"},
						colId2: colId{qual: "c", name: "id"},
						kind:   isEQ,
					},
				}, {
					bool: boolOr,
					pred: &predicateColumn{
						colId:  colId{qual: "d", name: "num"},
						colId2: colId{qual: "b", name: "num"},
						kind:   isGT,
					},
				}}},
				{joinType: joinFull, relId: relId{name: "relation_e", alias: "e"}, predicates: []*predicateItem{{
					pred: &predicateColumn{
						colId:  colId{qual: "e", name: "d_id"},
						colId2: colId{qual: "d", name: "id"},
						kind:   isEQ,
					},
				}, {
					bool: boolAnd,
					pred: &predicateColumn{
						colId: colId{qual: "e", name: "is_foo"},
						kind:  isFalse,
					},
				}}},
				{joinType: joinCross, relId: relId{name: "relation_f", alias: "f"}},
			}},
			whereBlock: &whereBlock{name: "Where", items: []*predicateItem{
				{pred: &predicateColumn{
					colId:  colId{qual: "a", name: "id"},
					kind:   isEQ,
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
				data:  dataType{typeInfo: typeInfo{kind: kindStruct}},
			},
			joinBlock: &joinBlock{relId: relId{name: "relation_b", alias: "b"}, items: []*joinItem{
				{joinType: joinLeft, relId: relId{name: "relation_c", alias: "c"}, predicates: []*predicateItem{{
					pred: &predicateColumn{
						colId:  colId{qual: "c", name: "b_id"},
						colId2: colId{qual: "b", name: "id"},
						kind:   isEQ,
					},
				}}},
				{joinType: joinRight, relId: relId{name: "relation_d", alias: "d"}, predicates: []*predicateItem{{
					pred: &predicateColumn{
						colId:  colId{qual: "d", name: "c_id"},
						colId2: colId{qual: "c", name: "id"},
						kind:   isEQ,
					},
				}, {
					bool: boolOr,
					pred: &predicateColumn{
						colId:  colId{qual: "d", name: "num"},
						colId2: colId{qual: "b", name: "num"},
						kind:   isGT,
					},
				}}},
				{joinType: joinFull, relId: relId{name: "relation_e", alias: "e"}, predicates: []*predicateItem{{
					pred: &predicateColumn{
						colId:  colId{qual: "e", name: "d_id"},
						colId2: colId{qual: "d", name: "id"},
						kind:   isEQ,
					},
				}, {
					bool: boolAnd,
					pred: &predicateColumn{
						colId: colId{qual: "e", name: "is_foo"},
						kind:  isFalse,
					},
				}}},
				{joinType: joinCross, relId: relId{name: "relation_f", alias: "f"}},
			}},
			whereBlock: &whereBlock{name: "Where", items: []*predicateItem{
				{pred: &predicateColumn{
					colId:  colId{qual: "a", name: "id"},
					colId2: colId{qual: "d", name: "a_id"},
					kind:   isEQ,
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
				data:  dataType{typeInfo: typeInfo{kind: kindStruct}},
			},
			joinBlock: &joinBlock{items: []*joinItem{
				{joinType: joinLeft, relId: relId{name: "relation_b", alias: "b"}, predicates: []*predicateItem{{
					pred: &predicateColumn{
						colId:  colId{qual: "b", name: "a_id"},
						colId2: colId{qual: "a", name: "id"},
						kind:   isEQ,
					},
				}}},
				{joinType: joinLeft, relId: relId{name: "relation_c", alias: "c"}, predicates: []*predicateItem{{
					pred: &predicateColumn{
						colId:  colId{qual: "c", name: "b_id"},
						colId2: colId{qual: "b", name: "id"},
						kind:   isEQ,
					},
				}}},
				{joinType: joinRight, relId: relId{name: "relation_d", alias: "d"}, predicates: []*predicateItem{{
					pred: &predicateColumn{
						colId:  colId{qual: "d", name: "c_id"},
						colId2: colId{qual: "c", name: "id"},
						kind:   isEQ,
					},
				}, {
					bool: boolOr,
					pred: &predicateColumn{
						colId:  colId{qual: "d", name: "num"},
						colId2: colId{qual: "b", name: "num"},
						kind:   isGT,
					},
				}}},
				{joinType: joinFull, relId: relId{name: "relation_e", alias: "e"}, predicates: []*predicateItem{{
					pred: &predicateColumn{
						colId:  colId{qual: "e", name: "d_id"},
						colId2: colId{qual: "d", name: "id"},
						kind:   isEQ,
					},
				}, {
					bool: boolAnd,
					pred: &predicateColumn{
						colId: colId{qual: "e", name: "is_foo"},
						kind:  isFalse,
					},
				}}},
				{joinType: joinCross, relId: relId{name: "relation_f", alias: "f"}},
			}},
			whereBlock: &whereBlock{name: "Where", items: []*predicateItem{
				{pred: &predicateColumn{
					colId:  colId{qual: "a", name: "id"},
					kind:   isEQ,
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
				data:  dataType{typeInfo: typeInfo{kind: kindStruct}},
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
				data:  dataType{typeInfo: typeInfo{kind: kindStruct}},
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
				data:  dataType{typeInfo: typeInfo{kind: kindStruct}},
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
				data:  dataType{typeInfo: typeInfo{kind: kindStruct}},
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
				data:  dataType{typeInfo: typeInfo{kind: kindStruct}},
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
				data:  dataType{typeInfo: typeInfo{kind: kindStruct}},
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
				data:  dataType{typeInfo: typeInfo{kind: kindStruct}},
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
				data:  dataType{typeInfo: typeInfo{kind: kindStruct}},
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
				data:  dataType{typeInfo: typeInfo{kind: kindStruct}},
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
				data:  dataType{typeInfo: typeInfo{kind: kindStruct}},
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
			whereBlock: &whereBlock{name: "Where", items: []*predicateItem{
				{pred: &predicateColumn{colId: colId{qual: "a", name: "is_inactive"}, kind: isTrue}},
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
			whereBlock: &whereBlock{name: "Where", items: []*predicateItem{
				{pred: &predicateColumn{colId: colId{qual: "a", name: "is_inactive"}, kind: isTrue}},
			}},
			rowsAffectedField: &rowsAffectedField{
				name: "RowsAffected",
				kind: kindInt,
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
					typeInfo: typeInfo{kind: kindStruct},
					fields: []*fieldInfo{{
						name: "f1", typ: typeInfo{kind: kindBool},
						colId: colId{name: "c1"},
						tag:   tagutil.Tag{"sql": {"c1"}},
					}, {
						name: "f2", typ: typeInfo{kind: kindUint8, isByte: true},
						colId: colId{name: "c2"},
						tag:   tagutil.Tag{"sql": {"c2"}},
					}, {
						name: "f3", typ: typeInfo{kind: kindInt32, isRune: true},
						colId: colId{name: "c3"},
						tag:   tagutil.Tag{"sql": {"c3"}},
					}, {
						name: "f4", typ: typeInfo{kind: kindInt8},
						colId: colId{name: "c4"},
						tag:   tagutil.Tag{"sql": {"c4"}},
					}, {
						name: "f5", typ: typeInfo{kind: kindInt16},
						colId: colId{name: "c5"},
						tag:   tagutil.Tag{"sql": {"c5"}},
					}, {
						name: "f6", typ: typeInfo{kind: kindInt32},
						colId: colId{name: "c6"},
						tag:   tagutil.Tag{"sql": {"c6"}},
					}, {
						name: "f7", typ: typeInfo{kind: kindInt64},
						colId: colId{name: "c7"},
						tag:   tagutil.Tag{"sql": {"c7"}},
					}, {
						name: "f8", typ: typeInfo{kind: kindInt},
						colId: colId{name: "c8"},
						tag:   tagutil.Tag{"sql": {"c8"}},
					}, {
						name: "f9", typ: typeInfo{kind: kindUint8},
						colId: colId{name: "c9"},
						tag:   tagutil.Tag{"sql": {"c9"}},
					}, {
						name: "f10", typ: typeInfo{kind: kindUint16},
						colId: colId{name: "c10"},
						tag:   tagutil.Tag{"sql": {"c10"}},
					}, {
						name: "f11", typ: typeInfo{kind: kindUint32},
						colId: colId{name: "c11"},
						tag:   tagutil.Tag{"sql": {"c11"}},
					}, {
						name: "f12", typ: typeInfo{kind: kindUint64},
						colId: colId{name: "c12"},
						tag:   tagutil.Tag{"sql": {"c12"}},
					}, {
						name: "f13", typ: typeInfo{kind: kindUint},
						colId: colId{name: "c13"},
						tag:   tagutil.Tag{"sql": {"c13"}},
					}, {
						name: "f14", typ: typeInfo{kind: kindUintptr},
						colId: colId{name: "c14"},
						tag:   tagutil.Tag{"sql": {"c14"}},
					}, {
						name: "f15", typ: typeInfo{kind: kindFloat32},
						colId: colId{name: "c15"},
						tag:   tagutil.Tag{"sql": {"c15"}},
					}, {
						name: "f16", typ: typeInfo{kind: kindFloat64},
						colId: colId{name: "c16"},
						tag:   tagutil.Tag{"sql": {"c16"}},
					}, {
						name: "f17", typ: typeInfo{kind: kindComplex64},
						colId: colId{name: "c17"},
						tag:   tagutil.Tag{"sql": {"c17"}},
					}, {
						name: "f18", typ: typeInfo{kind: kindComplex128},
						colId: colId{name: "c18"},
						tag:   tagutil.Tag{"sql": {"c18"}},
					}, {
						name: "f19", typ: typeInfo{kind: kindString},
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
					typeInfo: typeInfo{kind: kindStruct},
					fields: []*fieldInfo{{
						name: "f1", typ: typeInfo{
							kind: kindSlice,
							elem: &typeInfo{kind: kindBool},
						},
						colId: colId{name: "c1"},
						tag:   tagutil.Tag{"sql": {"c1"}},
					}, {
						name: "f2", typ: typeInfo{
							kind: kindSlice,
							elem: &typeInfo{kind: kindUint8, isByte: true},
						},
						colId: colId{name: "c2"},
						tag:   tagutil.Tag{"sql": {"c2"}},
					}, {
						name: "f3", typ: typeInfo{
							kind: kindSlice,
							elem: &typeInfo{kind: kindInt32, isRune: true},
						},
						colId: colId{name: "c3"},
						tag:   tagutil.Tag{"sql": {"c3"}},
					}, {
						name: "f4", typ: typeInfo{
							name:       "HardwareAddr",
							kind:       kindSlice,
							pkgPath:    "net",
							pkgName:    "net",
							pkgLocal:   "net",
							isImported: true,
							elem:       &typeInfo{kind: kindUint8, isByte: true},
						},
						colId: colId{name: "c4"},
						tag:   tagutil.Tag{"sql": {"c4"}},
					}, {
						name: "f5", typ: typeInfo{
							name:              "RawMessage",
							kind:              kindSlice,
							pkgPath:           "encoding/json",
							pkgName:           "json",
							pkgLocal:          "json",
							isImported:        true,
							isJSONMarshaler:   true,
							isJSONUnmarshaler: true,
							elem:              &typeInfo{kind: kindUint8, isByte: true},
						},
						colId: colId{name: "c5"},
						tag:   tagutil.Tag{"sql": {"c5"}},
					}, {
						name: "f6", typ: typeInfo{
							kind: kindSlice,
							elem: &typeInfo{
								name:            "Marshaler",
								kind:            kindInterface,
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
							kind: kindSlice,
							elem: &typeInfo{
								name:              "RawMessage",
								kind:              kindSlice,
								pkgPath:           "encoding/json",
								pkgName:           "json",
								pkgLocal:          "json",
								isImported:        true,
								isJSONMarshaler:   true,
								isJSONUnmarshaler: true,
								elem:              &typeInfo{kind: kindUint8, isByte: true},
							},
						},
						colId: colId{name: "c7"},
						tag:   tagutil.Tag{"sql": {"c7"}},
					}, {
						name: "f8", typ: typeInfo{
							kind: kindSlice,
							elem: &typeInfo{
								kind: kindSlice,
								elem: &typeInfo{kind: kindUint8, isByte: true},
							},
						},
						colId: colId{name: "c8"},
						tag:   tagutil.Tag{"sql": {"c8"}},
					}, {
						name: "f9", typ: typeInfo{
							kind: kindSlice,
							elem: &typeInfo{
								kind:     kindArray,
								arrayLen: 2,
								elem: &typeInfo{
									kind:     kindArray,
									arrayLen: 2,
									elem:     &typeInfo{kind: kindFloat64},
								},
							},
						},
						colId: colId{name: "c9"},
						tag:   tagutil.Tag{"sql": {"c9"}},
					}, {
						name: "f10", typ: typeInfo{
							kind: kindSlice,
							elem: &typeInfo{
								kind: kindSlice,
								elem: &typeInfo{
									kind:     kindArray,
									arrayLen: 2,
									elem:     &typeInfo{kind: kindFloat64},
								},
							},
						},
						colId: colId{name: "c10"},
						tag:   tagutil.Tag{"sql": {"c10"}},
					}, {
						name: "f11", typ: typeInfo{
							kind: kindMap,
							key:  &typeInfo{kind: kindString},
							elem: &typeInfo{
								name:       "NullString",
								kind:       kindStruct,
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
							kind: kindSlice,
							elem: &typeInfo{
								kind: kindMap,
								key:  &typeInfo{kind: kindString},
								elem: &typeInfo{
									kind: kindPtr,
									elem: &typeInfo{kind: kindString},
								},
							},
						},
						colId: colId{name: "c12"},
						tag:   tagutil.Tag{"sql": {"c12"}},
					}, {
						name: "f13", typ: typeInfo{
							kind: kindSlice,
							elem: &typeInfo{
								kind:     kindArray,
								arrayLen: 2,
								elem: &typeInfo{
									kind: kindPtr,
									elem: &typeInfo{
										name:              "Int",
										kind:              kindStruct,
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
					typeInfo: typeInfo{kind: kindStruct},
					fields: []*fieldInfo{{
						name: "f1", typ: typeInfo{
							name:            "Marshaler",
							kind:            kindInterface,
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
							kind:              kindInterface,
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
							kind:              kindInterface,
							isJSONMarshaler:   true,
							isJSONUnmarshaler: true,
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
			kind:     kindStruct,
			pkgPath:  "path/to/test",
			pkgName:  "testdata",
			pkgLocal: "testdata",
		},
		fields: []*fieldInfo{{
			typ:        typeInfo{kind: kindString},
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
		want string
	}{
		{"f01", gotypbool},
		{"f02", gotypbool},
		{"f03", gotypbools},
		{"f04", gotypstring},
		{"f05", gotypstring},
		{"f06", gotypstrings},
		{"f07", gotypstringss},
		{"f08", gotypstringm},
		{"f09", gotypstringm},
		{"f10", gotypstringms},
		{"f11", gotypstringms},
		{"f12", gotypbyte},
		{"f13", gotypbyte},
		{"f14", gotypbytes},
		{"f15", gotypbytess},
		{"f16", gotypbytea16},
		{"f17", gotypbytea16s},
		{"f18", gotyprune},
		{"f19", gotyprune},
		{"f20", gotyprunes},
		{"f21", gotypruness},
		{"f22", gotypint8},
		{"f23", gotypint8},
		{"f24", gotypint8s},
		{"f25", gotypint8ss},
		{"f26", gotypint16},
		{"f27", gotypint16},
		{"f28", gotypint16s},
		{"f29", gotypint16ss},
		{"f30", gotypint32},
		{"f31", gotypint32},
		{"f32", gotypint32s},
		{"f33", gotypint32a2},
		{"f34", gotypint32a2s},
		{"f35", gotypint64},
		{"f36", gotypint64},
		{"f37", gotypint64s},
		{"f38", gotypint64a2},
		{"f39", gotypint64a2s},
		{"f40", gotypfloat32},
		{"f41", gotypfloat32},
		{"f42", gotypfloat32s},
		{"f43", gotypfloat64},
		{"f44", gotypfloat64},
		{"f45", gotypfloat64s},
		{"f46", gotypfloat64a2},
		{"f47", gotypfloat64a2s},
		{"f48", gotypfloat64a2ss},
		{"f49", gotypfloat64a2a2},
		{"f50", gotypfloat64a2a2s},
		{"f51", gotypfloat64a3},
		{"f52", gotypfloat64a3s},
		{"f53", gotypipnet},
		{"f54", gotypipnets},
		{"f55", gotyptime},
		{"f56", gotyptime},
		{"f57", gotyptimes},
		{"f58", gotyptimes},
		{"f59", gotyptimea2},
		{"f60", gotyptimea2s},
		{"f61", gotypbytes},
		{"f62", gotypbytess},
		{"f63", gotypbigint},
		{"f64", gotypbigint},
		{"f65", gotypbigints},
		{"f66", gotypbigints},
		{"f67", gotypbiginta2},
		{"f68", gotypbiginta2},
		{"f69", gotypbiginta2s},
		{"f70", gotypnullstringm},
		{"f71", gotypnullstringms},
		{"f72", gotypbytes},
		{"f73", gotypbytess},
	}

	ti, err := runAnalysis("SelectAnalysisTestOK_typeinfo_string", t)
	if err != nil {
		t.Error(err)
	}
	fields := ti.query.dataField.data.fields
	for i := 0; i < len(fields); i++ {
		ff := fields[i]
		tt := tests[i]

		got := ff.typ.string(true)
		if ff.name != tt.name || got != tt.want {
			t.Errorf("got %s::%s, want %s::%s", ff.name, got, tt.name, tt.want)
		}
	}
}
