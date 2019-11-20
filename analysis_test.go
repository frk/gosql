package gosql

import (
	"testing"

	"github.com/frk/compare"
	"github.com/frk/gosql/internal/errors"
	"github.com/frk/gosql/internal/testutil"
	"github.com/frk/tagutil"
)

var tdata = testutil.ParseTestdata("testdata")

func runAnalysis(name string, t *testing.T) (*typespec, error) {
	named := testutil.FindNamedType(name, tdata)
	if named == nil {
		// Stop the test if no type with the given name was found.
		t.Fatal(name, " not found")
		return nil, nil
	}
	return analyze(named)
}

func TestAnalysis(t *testing.T) {

	// for reuse, analyzed common.User typeinfo
	commonUserTypeinfo := typeinfo{
		name:       "User",
		kind:       kindstruct,
		pkgpath:    "github.com/frk/gosql/testdata/common",
		pkgname:    "common",
		pkglocal:   "common",
		isimported: true,
	}

	commonUserFields := []*fieldinfo{{
		name:       "Id",
		typ:        typeinfo{kind: kindint},
		isexported: true,
		colid:      colid{name: "id"},
		tag:        tagutil.Tag{"sql": {"id"}},
	}, {
		name:       "Email",
		typ:        typeinfo{kind: kindstring},
		isexported: true,
		colid:      colid{name: "email"},
		tag:        tagutil.Tag{"sql": {"email"}},
	}, {
		name:       "FullName",
		typ:        typeinfo{kind: kindstring},
		isexported: true,
		colid:      colid{name: "full_name"},
		tag:        tagutil.Tag{"sql": {"full_name"}},
	}, {
		name: "CreatedAt",
		typ: typeinfo{
			name:            "Time",
			kind:            kindstruct,
			pkgpath:         "time",
			pkgname:         "time",
			pkglocal:        "time",
			isimported:      true,
			istime:          true,
			isjsmarshaler:   true,
			isjsunmarshaler: true,
		},
		isexported: true,
		colid:      colid{name: "created_at"},
		tag:        tagutil.Tag{"sql": {"created_at"}},
	}}

	reldummyslice := &relfield{
		name:  "Rel",
		relid: relid{name: "relation_a", alias: "a"},
		rec: recordtype{
			base: typeinfo{
				name:     "T",
				kind:     kindstruct,
				pkgpath:  "path/to/test",
				pkgname:  "testdata",
				pkglocal: "testdata",
			},
			isslice: true,
			fields: []*fieldinfo{{
				typ:        typeinfo{kind: kindstring},
				name:       "F",
				isexported: true,
				tag:        tagutil.Tag{"sql": {"f"}},
				colid:      colid{name: "f"},
			}},
		},
	}

	dummyrecord := recordtype{
		base: typeinfo{
			name:     "T",
			kind:     kindstruct,
			pkgpath:  "path/to/test",
			pkgname:  "testdata",
			pkglocal: "testdata",
		},
		fields: []*fieldinfo{{
			typ:        typeinfo{kind: kindstring},
			name:       "F",
			isexported: true,
			tag:        tagutil.Tag{"sql": {"f"}},
			colid:      colid{name: "f"},
		}},
	}

	tests := []struct {
		name string
		want *typespec
		err  error
	}{{
		name: "InsertAnalysisTestBAD_NoRelfield",
		err:  errors.NoRelfieldError,
	}, {
		name: "InsertAnalysisTestBAD3",
		err:  errors.BadRelfieldTypeError,
	}, {
		name: "DeleteAnalysisTestBAD_BadRelId",
		err:  errors.BadRelIdError,
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
		err:  errors.BadRelfieldTypeError,
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
		name: "FilterAnalysisTestBAD_IllegalReturnDirective",
		err:  errors.IllegalReturnDirectiveError,
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
		err:  errors.BadRelfieldTypeError,
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
		name: "SelectAnalysisTestBAD_BadColumnCmpopCombo",
		err:  errors.BadCmpopComboError,
	}, {
		name: "DeleteAnalysisTestBAD_BadColumnExpressionLHS",
		err:  errors.BadColIdError,
	}, {
		name: "UpdateAnalysisTestBAD_BadUnaryOp",
		err:  errors.BadUnaryCmpopError,
	}, {
		name: "UpdateAnalysisTestBAD_ExtraScalarrop",
		err:  errors.ExtraScalarropError,
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
		name: "DeleteAnalysisTestBAD_BadWhereFieldCmpopCombo",
		err:  errors.BadCmpopComboError,
	}, {
		name: "DeleteAnalysisTestBAD_IllegalWhereFieldUnaryCmp",
		err:  errors.IllegalUnaryComparisonOperatorError,
	}, {
		name: "UpdateAnalysisTestBAD_BadWhereFieldTypeForScalarrop",
		err:  errors.BadScalarFieldTypeError,
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
		name: "SelectAnalysisTestBAD_BadJoinDirectiveExpressionCmpop",
		err:  errors.BadUnaryCmpopError,
	}, {
		name: "SelectAnalysisTestBAD_BadJoinDirectiveExpressionExtraScalarrop",
		err:  errors.ExtraScalarropError,
	}, {
		name: "SelectAnalysisTestBAD_BadJoinDirectiveExpressionCmpopCombo",
		err:  errors.BadCmpopComboError,
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
		err:  errors.BadRelfieldTypeError,
	}, {
		name: "DeleteAnalysisTestBAD_ConflictResultProducer2",
		err:  errors.ConflictResultProducerError,
	}, {
		name: "DeleteAnalysisTestBAD_BadRowsAffecteFieldType",
		err:  errors.BadRowsAffectedTypeError,
	}, {
		name: "FilterAnalysisTestBAD_BadTextSearchDirectiveColId",
		err:  errors.BadColIdError,
	}, {
		name: "InsertAnalysisTestOK1",
		want: &typespec{name: "InsertAnalysisTestOK1", kind: speckindInsert, rel: &relfield{
			name:  "UserRec",
			relid: relid{name: "users_table"},
			rec: recordtype{
				base:      commonUserTypeinfo,
				fields:    commonUserFields,
				ispointer: true,
			},
		}},
	}, {
		name: "InsertAnalysisTestOK2",
		want: &typespec{name: "InsertAnalysisTestOK2", kind: speckindInsert, rel: &relfield{
			name:  "UserRec",
			relid: relid{name: "users_table"},
			rec: recordtype{
				base: typeinfo{
					kind: kindstruct,
				},
				fields: []*fieldinfo{{
					name:       "Name3",
					typ:        typeinfo{kind: kindstring},
					isexported: true,
					colid:      colid{name: "name"},
					tag:        tagutil.Tag{"sql": {"name"}},
				}},
			},
		}},
	}, {
		name: "SelectAnalysisTestOK3",
		want: &typespec{name: "SelectAnalysisTestOK3", kind: speckindSelect, rel: &relfield{
			name:  "User",
			relid: relid{name: "users_table"},
			rec: recordtype{
				base:      commonUserTypeinfo,
				fields:    commonUserFields,
				ispointer: true,
				isiter:    true,
			},
		}},
	}, {
		name: "SelectAnalysisTestOK4",
		want: &typespec{name: "SelectAnalysisTestOK4", kind: speckindSelect, rel: &relfield{
			name:  "User",
			relid: relid{name: "users_table"},
			rec: recordtype{
				base:      commonUserTypeinfo,
				fields:    commonUserFields,
				ispointer: true,
				isiter:    true,
			},
		}},
	}, {
		name: "SelectAnalysisTestOK5",
		want: &typespec{name: "SelectAnalysisTestOK5", kind: speckindSelect, rel: &relfield{
			name:  "User",
			relid: relid{name: "users_table"},
			rec: recordtype{
				base:       commonUserTypeinfo,
				fields:     commonUserFields,
				ispointer:  true,
				isiter:     true,
				itermethod: "Fn",
			},
		}},
	}, {
		name: "SelectAnalysisTestOK6",
		want: &typespec{name: "SelectAnalysisTestOK6", kind: speckindSelect, rel: &relfield{
			name:  "User",
			relid: relid{name: "users_table"},
			rec: recordtype{
				base:       commonUserTypeinfo,
				fields:     commonUserFields,
				ispointer:  true,
				isiter:     true,
				itermethod: "Fn",
			},
		}},
	}, {
		name: "SelectAnalysisTestOK7",
		want: &typespec{name: "SelectAnalysisTestOK7", kind: speckindSelect, rel: &relfield{
			name:  "Rel",
			relid: relid{name: "relation_a"},
			rec: recordtype{
				base: typeinfo{
					kind: kindstruct,
				},
				fields: []*fieldinfo{{
					name:   "a",
					typ:    typeinfo{kind: kindint},
					colid:  colid{name: "a"},
					tag:    tagutil.Tag{"sql": {"a", "pk", "auto"}},
					ispkey: true,
					auto:   true,
				}, {
					name:      "b",
					typ:       typeinfo{kind: kindint},
					colid:     colid{name: "b"},
					tag:       tagutil.Tag{"sql": {"b", "nullempty"}},
					nullempty: true,
				}, {
					name:     "c",
					typ:      typeinfo{kind: kindint},
					colid:    colid{name: "c"},
					tag:      tagutil.Tag{"sql": {"c", "ro", "json"}},
					readonly: true,
					usejson:  true,
				}, {
					name:      "d",
					typ:       typeinfo{kind: kindint},
					colid:     colid{name: "d"},
					tag:       tagutil.Tag{"sql": {"d", "wo"}},
					writeonly: true,
				}, {
					name:   "e",
					typ:    typeinfo{kind: kindint},
					colid:  colid{name: "e"},
					tag:    tagutil.Tag{"sql": {"e", "+"}},
					binadd: true,
				}, {
					name:        "f",
					typ:         typeinfo{kind: kindint},
					colid:       colid{name: "f"},
					tag:         tagutil.Tag{"sql": {"f", "coalesce"}},
					usecoalesce: true,
				}, {
					name:        "g",
					typ:         typeinfo{kind: kindint},
					colid:       colid{name: "g"},
					tag:         tagutil.Tag{"sql": {"g", "coalesce(-1)"}},
					usecoalesce: true,
					coalesceval: "-1",
				}},
			},
		}},
	}, {
		name: "InsertAnalysisTestOK8",
		want: &typespec{name: "InsertAnalysisTestOK8", kind: speckindInsert, rel: &relfield{
			name:  "Rel",
			relid: relid{name: "relation_a"},
			rec: recordtype{
				base: typeinfo{
					kind: kindstruct,
				},
				fields: []*fieldinfo{{
					name: "Val",
					path: []*fieldelem{
						{
							name:         "Foobar",
							tag:          tagutil.Tag{"sql": {">foo_"}},
							typename:     "Foo",
							typepkgpath:  "github.com/frk/gosql/testdata/common",
							typepkgname:  "common",
							typepkglocal: "common",
							isexported:   true,
							isimported:   true,
						},
						{
							name:         "Bar",
							tag:          tagutil.Tag{"sql": {">bar_"}},
							typename:     "Bar",
							typepkgpath:  "github.com/frk/gosql/testdata/common",
							typepkgname:  "common",
							typepkglocal: "common",
							isimported:   true,
							isexported:   true,
						},
						{
							name:         "Baz",
							tag:          tagutil.Tag{"sql": {">baz_"}},
							typename:     "Baz",
							typepkgpath:  "github.com/frk/gosql/testdata/common",
							typepkgname:  "common",
							typepkglocal: "common",
							isexported:   true,
							isembedded:   true,
							isimported:   true,
						},
					},
					isexported: true,
					typ:        typeinfo{kind: kindstring},
					colid:      colid{name: "foo_bar_baz_val"},
					tag:        tagutil.Tag{"sql": {"val"}},
				}, {
					name: "Val",
					path: []*fieldelem{{
						name:         "Foobar",
						tag:          tagutil.Tag{"sql": {">foo_"}},
						typename:     "Foo",
						typepkgpath:  "github.com/frk/gosql/testdata/common",
						typepkgname:  "common",
						typepkglocal: "common",
						isexported:   true,
						isimported:   true,
					}, {
						name:         "Baz",
						tag:          tagutil.Tag{"sql": {">baz_"}},
						typename:     "Baz",
						typepkgpath:  "github.com/frk/gosql/testdata/common",
						typepkgname:  "common",
						typepkglocal: "common",
						isimported:   true,
						isexported:   true,
						isembedded:   false,
						ispointer:    true,
					}},
					isexported: true,
					typ:        typeinfo{kind: kindstring},
					colid:      colid{name: "foo_baz_val"},
					tag:        tagutil.Tag{"sql": {"val"}},
				}},
			},
		}},
	}, {
		name: "DeleteAnalysisTestOK9",
		want: &typespec{
			name: "DeleteAnalysisTestOK9",
			kind: speckindDelete,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a"},
				rec:   recordtype{base: typeinfo{kind: kindstruct}},
			},
			where: &whereblock{name: "Where", items: []*whereitem{{
				node: &wherefield{
					name:  "ID",
					colid: colid{name: "id"},
					typ:   typeinfo{kind: kindint},
					cmp:   cmpeq,
				},
			}}},
		},
	}, {
		name: "DeleteAnalysisTestOK10",
		want: &typespec{
			name: "DeleteAnalysisTestOK10",
			kind: speckindDelete,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a"},
				rec:   recordtype{base: typeinfo{kind: kindstruct}},
			},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &wherecolumn{colid: colid{name: "column_a"}, cmp: cmpnotnull}},
				{op: booland, node: &wherecolumn{colid: colid{name: "column_b"}, cmp: cmpisnull}},
				{op: boolor, node: &wherecolumn{colid: colid{name: "column_c"}, cmp: cmpnottrue}},
				{op: booland, node: &wherecolumn{colid: colid{name: "column_d"}, cmp: cmpistrue}},
				{op: boolor, node: &wherecolumn{colid: colid{name: "column_e"}, cmp: cmpnotfalse}},
				{op: boolor, node: &wherecolumn{colid: colid{name: "column_f"}, cmp: cmpisfalse}},
				{op: booland, node: &wherecolumn{colid: colid{name: "column_g"}, cmp: cmpnotunknown}},
				{op: booland, node: &wherecolumn{colid: colid{name: "column_h"}, cmp: cmpisunknown}},
				{op: booland, node: &wherecolumn{colid: colid{name: "column_i"}, cmp: cmpistrue}},
			}},
		},
	}, {
		name: "DeleteAnalysisTestOK11",
		want: &typespec{
			name: "DeleteAnalysisTestOK11",
			kind: speckindDelete,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a"},
				rec:   recordtype{base: typeinfo{kind: kindstruct}},
			},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &whereblock{name: "x", items: []*whereitem{
					{node: &wherefield{
						name:  "foo",
						typ:   typeinfo{kind: kindint},
						colid: colid{name: "column_foo"},
						cmp:   cmpeq,
					}},
					{op: booland, node: &wherecolumn{colid: colid{name: "column_a"}, cmp: cmpisnull}},
				}}},
				{op: boolor, node: &whereblock{name: "y", items: []*whereitem{
					{node: &wherecolumn{colid: colid{name: "column_b"}, cmp: cmpnottrue}},
					{op: boolor, node: &wherefield{
						name:  "bar",
						typ:   typeinfo{kind: kindstring},
						colid: colid{name: "column_bar"},
						cmp:   cmpeq,
					}},
					{op: booland, node: &whereblock{name: "z", items: []*whereitem{
						{node: &wherefield{
							name:  "baz",
							typ:   typeinfo{kind: kindbool},
							colid: colid{name: "column_baz"},
							cmp:   cmpeq,
						}},
						{op: booland, node: &wherefield{
							name:  "quux",
							typ:   typeinfo{kind: kindstring},
							colid: colid{name: "column_quux"},
							cmp:   cmpeq,
						}},
						{op: boolor, node: &wherecolumn{colid: colid{name: "column_c"}, cmp: cmpistrue}},
					}}},
				}}},
				{op: boolor, node: &wherecolumn{colid: colid{name: "column_d"}, cmp: cmpnotfalse}},
				{op: booland, node: &wherecolumn{colid: colid{name: "column_e"}, cmp: cmpisfalse}},
				{op: booland, node: &wherefield{
					name:  "foo",
					typ:   typeinfo{kind: kindint},
					colid: colid{name: "column_foo"},
					cmp:   cmpeq,
				}},
			}},
		},
	}, {
		name: "DeleteAnalysisTestOK12",
		want: &typespec{
			name: "DeleteAnalysisTestOK12",
			kind: speckindDelete,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a"},
				rec:   recordtype{base: typeinfo{kind: kindstruct}},
			},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &wherefield{name: "a", typ: typeinfo{kind: kindint}, colid: colid{name: "column_a"}, cmp: cmplt}},
				{op: booland, node: &wherefield{name: "b", typ: typeinfo{kind: kindint}, colid: colid{name: "column_b"}, cmp: cmpgt}},
				{op: booland, node: &wherefield{name: "c", typ: typeinfo{kind: kindint}, colid: colid{name: "column_c"}, cmp: cmple}},
				{op: booland, node: &wherefield{name: "d", typ: typeinfo{kind: kindint}, colid: colid{name: "column_d"}, cmp: cmpge}},
				{op: booland, node: &wherefield{name: "e", typ: typeinfo{kind: kindint}, colid: colid{name: "column_e"}, cmp: cmpeq}},
				{op: booland, node: &wherefield{name: "f", typ: typeinfo{kind: kindint}, colid: colid{name: "column_f"}, cmp: cmpne}},
				{op: booland, node: &wherefield{name: "g", typ: typeinfo{kind: kindint}, colid: colid{name: "column_g"}, cmp: cmpeq}},
			}},
		},
	}, {
		name: "DeleteAnalysisTestOK13",
		want: &typespec{
			name: "DeleteAnalysisTestOK13",
			kind: speckindDelete,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a"},
				rec:   recordtype{base: typeinfo{kind: kindstruct}},
			},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &wherecolumn{colid: colid{name: "column_a"}, cmp: cmpne, colid2: colid{name: "column_b"}}},
				{op: booland, node: &wherecolumn{colid: colid{qual: "t", name: "column_c"}, cmp: cmpeq, colid2: colid{qual: "u", name: "column_d"}}},
				{op: booland, node: &wherecolumn{colid: colid{qual: "t", name: "column_e"}, cmp: cmpgt, lit: "123"}},
				{op: booland, node: &wherecolumn{colid: colid{qual: "t", name: "column_f"}, cmp: cmpeq, lit: "'active'"}},
				{op: booland, node: &wherecolumn{colid: colid{qual: "t", name: "column_g"}, cmp: cmpne, lit: "true"}},
			}},
		},
	}, {
		name: "DeleteAnalysisTestOK14",
		want: &typespec{
			name: "DeleteAnalysisTestOK14",
			kind: speckindDelete,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a"},
				rec:   recordtype{base: typeinfo{kind: kindstruct}},
			},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &wherebetween{
					name:  "a",
					colid: colid{name: "column_a"},
					cmp:   cmpisbetween,
					x:     &varinfo{name: "x", typ: typeinfo{kind: kindint}},
					y:     &varinfo{name: "y", typ: typeinfo{kind: kindint}},
				}},
				{op: booland, node: &wherebetween{
					name:  "b",
					colid: colid{name: "column_b"},
					cmp:   cmpisbetweensym,
					x:     colid{name: "column_x"},
					y:     colid{name: "column_y"},
				}},
				{op: booland, node: &wherebetween{
					name:  "c",
					colid: colid{name: "column_c"},
					cmp:   cmpnotbetweensym,
					x:     colid{name: "column_z"},
					y:     &varinfo{name: "z", typ: typeinfo{kind: kindint}},
				}},
				{op: booland, node: &wherebetween{
					name:  "d",
					colid: colid{name: "column_d"},
					cmp:   cmpnotbetween,
					x:     &varinfo{name: "z", typ: typeinfo{kind: kindint}},
					y:     colid{name: "column_z"},
				}},
			}},
		},
	}, {
		name: "DeleteAnalysisTestOK_DistinctFrom",
		want: &typespec{
			name: "DeleteAnalysisTestOK_DistinctFrom",
			kind: speckindDelete,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a"},
				rec:   recordtype{base: typeinfo{kind: kindstruct}},
			},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &wherefield{
					name:  "a",
					typ:   typeinfo{kind: kindint},
					colid: colid{name: "column_a"},
					cmp:   cmpisdistinct,
				}},
				{op: booland, node: &wherefield{
					name:  "b",
					typ:   typeinfo{kind: kindint},
					colid: colid{name: "column_b"},
					cmp:   cmpnotdistinct,
				}},
				{op: booland, node: &wherecolumn{colid: colid{name: "column_c"}, cmp: cmpisdistinct, colid2: colid{name: "column_x"}}},
				{op: booland, node: &wherecolumn{colid: colid{name: "column_d"}, cmp: cmpnotdistinct, colid2: colid{name: "column_y"}}},
			}},
		},
	}, {
		name: "DeleteAnalysisTestOK_ArrayComparisons",
		want: &typespec{
			name: "DeleteAnalysisTestOK_ArrayComparisons",
			kind: speckindDelete,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a"},
				rec:   recordtype{base: typeinfo{kind: kindstruct}},
			},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &wherefield{
					name: "a",
					typ: typeinfo{
						kind: kindslice,
						elem: &typeinfo{
							kind: kindint,
						},
					},
					colid: colid{name: "column_a"},
					cmp:   cmpisin,
				}},
				{op: booland, node: &wherefield{
					name: "b",
					typ: typeinfo{
						kind: kindarray,
						elem: &typeinfo{
							kind: kindint,
						},
						arraylen: 5,
					},
					colid: colid{name: "column_b"},
					cmp:   cmpnotin,
				}},
				{op: booland, node: &wherefield{
					name: "c",
					typ: typeinfo{
						kind: kindslice,
						elem: &typeinfo{
							kind: kindint,
						},
					},
					colid: colid{name: "column_c"},
					cmp:   cmpeq,
					sop:   scalarrany,
				}},
				{op: booland, node: &wherefield{
					name: "d",
					typ: typeinfo{
						kind: kindarray,
						elem: &typeinfo{
							kind: kindint,
						},
						arraylen: 10,
					},
					colid: colid{name: "column_d"},
					cmp:   cmpgt,
					sop:   scalarrsome,
				}},
				{op: booland, node: &wherefield{
					name: "e",
					typ: typeinfo{
						kind: kindslice,
						elem: &typeinfo{
							kind: kindint,
						},
					},
					colid: colid{name: "column_e"},
					cmp:   cmple,
					sop:   scalarrall,
				}},
			}},
		},
	}, {
		name: "DeleteAnalysisTestOK_PatternMatching",
		want: &typespec{
			name: "DeleteAnalysisTestOK_PatternMatching",
			kind: speckindDelete,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a"},
				rec:   recordtype{base: typeinfo{kind: kindstruct}},
			},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &wherefield{
					name:  "a",
					typ:   typeinfo{kind: kindstring},
					colid: colid{name: "column_a"},
					cmp:   cmpislike,
				}},
				{op: booland, node: &wherefield{
					name:  "b",
					typ:   typeinfo{kind: kindstring},
					colid: colid{name: "column_b"},
					cmp:   cmpnotlike,
				}},
				{op: booland, node: &wherefield{
					name:  "c",
					typ:   typeinfo{kind: kindstring},
					colid: colid{name: "column_c"},
					cmp:   cmpissimilar,
				}},
				{op: booland, node: &wherefield{
					name:  "d",
					typ:   typeinfo{kind: kindstring},
					colid: colid{name: "column_d"},
					cmp:   cmpnotsimilar,
				}},
				{op: booland, node: &wherefield{
					name:  "e",
					typ:   typeinfo{kind: kindstring},
					colid: colid{name: "column_e"},
					cmp:   cmprexp,
				}},
				{op: booland, node: &wherefield{
					name:  "f",
					typ:   typeinfo{kind: kindstring},
					colid: colid{name: "column_f"},
					cmp:   cmprexpi,
				}},
				{op: booland, node: &wherefield{
					name:  "g",
					typ:   typeinfo{kind: kindstring},
					colid: colid{name: "column_g"},
					cmp:   cmpnotrexp,
				}},
				{op: booland, node: &wherefield{
					name:  "h",
					typ:   typeinfo{kind: kindstring},
					colid: colid{name: "column_h"},
					cmp:   cmpnotrexpi,
				}},
			}},
		},
	}, {
		name: "DeleteAnalysisTestOK_Using",
		want: &typespec{
			name: "DeleteAnalysisTestOK_Using",
			kind: speckindDelete,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a", alias: "a"},
				rec:   recordtype{base: typeinfo{kind: kindstruct}},
			},
			join: &joinblock{rel: relid{name: "relation_b", alias: "b"}, items: []*joinitem{
				{typ: joinleft, rel: relid{name: "relation_c", alias: "c"}, conds: []*joincond{{
					col1: colid{qual: "c", name: "b_id"},
					col2: colid{qual: "b", name: "id"},
					cmp:  cmpeq,
				}}},
				{typ: joinright, rel: relid{name: "relation_d", alias: "d"}, conds: []*joincond{{
					col1: colid{qual: "d", name: "c_id"},
					col2: colid{qual: "c", name: "id"},
					cmp:  cmpeq,
				}, {
					op:   boolor,
					col1: colid{qual: "d", name: "num"},
					col2: colid{qual: "b", name: "num"},
					cmp:  cmpgt,
				}}},
				{typ: joinfull, rel: relid{name: "relation_e", alias: "e"}, conds: []*joincond{{
					col1: colid{qual: "e", name: "d_id"},
					col2: colid{qual: "d", name: "id"},
					cmp:  cmpeq,
				}, {
					op:   booland,
					col1: colid{qual: "e", name: "is_foo"},
					cmp:  cmpisfalse,
				}}},
				{typ: joincross, rel: relid{name: "relation_f", alias: "f"}},
			}},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &wherecolumn{
					colid:  colid{qual: "a", name: "id"},
					cmp:    cmpeq,
					colid2: colid{qual: "d", name: "a_id"},
				}},
			}},
		},
	}, {
		name: "UpdateAnalysisTestOK_From",
		want: &typespec{
			name: "UpdateAnalysisTestOK_From",
			kind: speckindUpdate,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a", alias: "a"},
				rec:   recordtype{base: typeinfo{kind: kindstruct}},
			},
			join: &joinblock{rel: relid{name: "relation_b", alias: "b"}, items: []*joinitem{
				{typ: joinleft, rel: relid{name: "relation_c", alias: "c"}, conds: []*joincond{{
					col1: colid{qual: "c", name: "b_id"},
					col2: colid{qual: "b", name: "id"},
					cmp:  cmpeq,
				}}},
				{typ: joinright, rel: relid{name: "relation_d", alias: "d"}, conds: []*joincond{{
					col1: colid{qual: "d", name: "c_id"},
					col2: colid{qual: "c", name: "id"},
					cmp:  cmpeq,
				}, {
					op:   boolor,
					col1: colid{qual: "d", name: "num"},
					col2: colid{qual: "b", name: "num"},
					cmp:  cmpgt,
				}}},
				{typ: joinfull, rel: relid{name: "relation_e", alias: "e"}, conds: []*joincond{{
					col1: colid{qual: "e", name: "d_id"},
					col2: colid{qual: "d", name: "id"},
					cmp:  cmpeq,
				}, {
					op:   booland,
					col1: colid{qual: "e", name: "is_foo"},
					cmp:  cmpisfalse,
				}}},
				{typ: joincross, rel: relid{name: "relation_f", alias: "f"}},
			}},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &wherecolumn{
					colid:  colid{qual: "a", name: "id"},
					cmp:    cmpeq,
					colid2: colid{qual: "d", name: "a_id"},
				}},
			}},
		},
	}, {
		name: "SelectAnalysisTestOK_Join",
		want: &typespec{
			name: "SelectAnalysisTestOK_Join",
			kind: speckindSelect,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a", alias: "a"},
				rec:   recordtype{base: typeinfo{kind: kindstruct}},
			},
			join: &joinblock{items: []*joinitem{
				{typ: joinleft, rel: relid{name: "relation_b", alias: "b"}, conds: []*joincond{{
					col1: colid{qual: "b", name: "a_id"},
					col2: colid{qual: "a", name: "id"},
					cmp:  cmpeq,
				}}},
				{typ: joinleft, rel: relid{name: "relation_c", alias: "c"}, conds: []*joincond{{
					col1: colid{qual: "c", name: "b_id"},
					col2: colid{qual: "b", name: "id"},
					cmp:  cmpeq,
				}}},
				{typ: joinright, rel: relid{name: "relation_d", alias: "d"}, conds: []*joincond{{
					col1: colid{qual: "d", name: "c_id"},
					col2: colid{qual: "c", name: "id"},
					cmp:  cmpeq,
				}, {
					op:   boolor,
					col1: colid{qual: "d", name: "num"},
					col2: colid{qual: "b", name: "num"},
					cmp:  cmpgt,
				}}},
				{typ: joinfull, rel: relid{name: "relation_e", alias: "e"}, conds: []*joincond{{
					col1: colid{qual: "e", name: "d_id"},
					col2: colid{qual: "d", name: "id"},
					cmp:  cmpeq,
				}, {
					op:   booland,
					col1: colid{qual: "e", name: "is_foo"},
					cmp:  cmpisfalse,
				}}},
				{typ: joincross, rel: relid{name: "relation_f", alias: "f"}},
			}},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &wherecolumn{
					colid:  colid{qual: "a", name: "id"},
					cmp:    cmpeq,
					colid2: colid{qual: "d", name: "a_id"},
				}},
			}},
		},
	}, {
		name: "UpdateAnalysisTestOK_All",
		want: &typespec{
			name: "UpdateAnalysisTestOK_All",
			kind: speckindUpdate,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a", alias: "a"},
				rec:   recordtype{base: typeinfo{kind: kindstruct}},
			},
			all: true,
		},
	}, {
		name: "DeleteAnalysisTestOK_All",
		want: &typespec{
			name: "DeleteAnalysisTestOK_All",
			kind: speckindDelete,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a", alias: "a"},
				rec:   recordtype{base: typeinfo{kind: kindstruct}},
			},
			all: true,
		},
	}, {
		name: "DeleteAnalysisTestOK_Return",
		want: &typespec{
			name: "DeleteAnalysisTestOK_Return",
			kind: speckindDelete,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a", alias: "a"},
				rec:   dummyrecord,
			},
			returning: &collist{all: true},
		},
	}, {
		name: "InsertAnalysisTestOK_Return",
		want: &typespec{
			name: "InsertAnalysisTestOK_Return",
			kind: speckindInsert,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a", alias: "a"},
				rec:   dummyrecord,
			},
			returning: &collist{items: []colid{
				{qual: "a", name: "foo"},
				{qual: "a", name: "bar"},
				{qual: "a", name: "baz"}}},
		},
	}, {
		name: "UpdateAnalysisTestOK_Return",
		want: &typespec{
			name: "UpdateAnalysisTestOK_Return",
			kind: speckindUpdate,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a", alias: "a"},
				rec:   dummyrecord,
			},
			returning: &collist{items: []colid{
				{qual: "a", name: "foo"},
				{qual: "a", name: "bar"},
				{qual: "a", name: "baz"}}},
		},
	}, {
		name: "InsertAnalysisTestOK_Default",
		want: &typespec{
			name: "InsertAnalysisTestOK_Default",
			kind: speckindInsert,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a", alias: "a"},
				rec:   recordtype{base: typeinfo{kind: kindstruct}},
			},
			defaults: &collist{all: true},
		},
	}, {
		name: "UpdateAnalysisTestOK_Default",
		want: &typespec{
			name: "UpdateAnalysisTestOK_Default",
			kind: speckindUpdate,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a", alias: "a"},
				rec:   recordtype{base: typeinfo{kind: kindstruct}},
			},
			defaults: &collist{items: []colid{
				{qual: "a", name: "foo"},
				{qual: "a", name: "bar"},
				{qual: "a", name: "baz"}}},
		},
	}, {
		name: "InsertAnalysisTestOK_Force",
		want: &typespec{
			name: "InsertAnalysisTestOK_Force",
			kind: speckindInsert,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a", alias: "a"},
				rec:   recordtype{base: typeinfo{kind: kindstruct}},
			},
			force: &collist{all: true},
		},
	}, {
		name: "UpdateAnalysisTestOK_Force",
		want: &typespec{
			name: "UpdateAnalysisTestOK_Force",
			kind: speckindUpdate,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a", alias: "a"},
				rec:   recordtype{base: typeinfo{kind: kindstruct}},
			},
			force: &collist{items: []colid{
				{qual: "a", name: "foo"},
				{qual: "a", name: "bar"},
				{qual: "a", name: "baz"}}},
		},
	}, {
		name: "SelectAnalysisTestOK_ErrorHandler",
		want: &typespec{
			name: "SelectAnalysisTestOK_ErrorHandler",
			kind: speckindSelect,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a", alias: "a"},
				rec:   recordtype{base: typeinfo{kind: kindstruct}},
			},
			erh: "eh",
		},
	}, {
		name: "InsertAnalysisTestOK_ErrorHandler",
		want: &typespec{
			name: "InsertAnalysisTestOK_ErrorHandler",
			kind: speckindInsert,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a", alias: "a"},
				rec:   recordtype{base: typeinfo{kind: kindstruct}},
			},
			erh: "myerrorhandler",
		},
	}, {
		name: "SelectAnalysisTestOK_Count",
		want: &typespec{
			name: "SelectAnalysisTestOK_Count",
			kind: speckindSelect,
			rel: &relfield{
				name:  "Count",
				relid: relid{name: "relation_a", alias: "a"},
			},
			selkind: selectcount,
		},
	}, {
		name: "SelectAnalysisTestOK_Exists",
		want: &typespec{
			name: "SelectAnalysisTestOK_Exists",
			kind: speckindSelect,
			rel: &relfield{
				name:  "Exists",
				relid: relid{name: "relation_a", alias: "a"},
			},
			selkind: selectexists,
		},
	}, {
		name: "SelectAnalysisTestOK_NotExists",
		want: &typespec{
			name: "SelectAnalysisTestOK_NotExists",
			kind: speckindSelect,
			rel: &relfield{
				name:  "NotExists",
				relid: relid{name: "relation_a", alias: "a"},
			},
			selkind: selectnotexists,
		},
	}, {
		name: "DeleteAnalysisTestOK_Relation",
		want: &typespec{
			name: "DeleteAnalysisTestOK_Relation",
			kind: speckindDelete,
			rel: &relfield{
				name:  "_",
				relid: relid{name: "relation_a", alias: "a"},
				isdir: true,
			},
		},
	}, {
		name: "SelectAnalysisTestOK_LimitDirective",
		want: &typespec{
			name:  "SelectAnalysisTestOK_LimitDirective",
			kind:  speckindSelect,
			rel:   reldummyslice,
			limit: &limitvar{value: 25},
		},
	}, {
		name: "SelectAnalysisTestOK_LimitField",
		want: &typespec{
			name:  "SelectAnalysisTestOK_LimitField",
			kind:  speckindSelect,
			rel:   reldummyslice,
			limit: &limitvar{value: 10, field: "Limit"},
		},
	}, {
		name: "SelectAnalysisTestOK_OffsetDirective",
		want: &typespec{
			name:   "SelectAnalysisTestOK_OffsetDirective",
			kind:   speckindSelect,
			rel:    reldummyslice,
			offset: &offsetvar{value: 25},
		},
	}, {
		name: "SelectAnalysisTestOK_OffsetField",
		want: &typespec{
			name:   "SelectAnalysisTestOK_OffsetField",
			kind:   speckindSelect,
			rel:    reldummyslice,
			offset: &offsetvar{value: 10, field: "Offset"},
		},
	}, {
		name: "SelectAnalysisTestOK_OrderByDirective",
		want: &typespec{
			name: "SelectAnalysisTestOK_OrderByDirective",
			kind: speckindSelect,
			rel:  reldummyslice,
			orderby: &orderbylist{items: []*orderbyitem{
				{col: colid{qual: "a", name: "foo"}, dir: orderasc, nulls: nullsfirst},
				{col: colid{qual: "a", name: "bar"}, dir: orderdesc, nulls: nullsfirst},
				{col: colid{qual: "a", name: "baz"}, dir: orderdesc, nulls: 0},
				{col: colid{qual: "a", name: "quux"}, dir: orderasc, nulls: nullslast},
			}},
		},
	}, {
		name: "InsertAnalysisTestOK_OverrideDirective",
		want: &typespec{
			name:     "InsertAnalysisTestOK_OverrideDirective",
			kind:     speckindInsert,
			rel:      reldummyslice,
			override: overridingsystem,
		},
	}, {
		name: "FilterAnalysisTestOK_TextSearchDirective",
		want: &typespec{
			name: "FilterAnalysisTestOK_TextSearchDirective",
			kind: speckindFilter,
			rel: &relfield{
				name:  "_",
				relid: relid{name: "relation_a", alias: "a"},
				rec:   dummyrecord,
			},
			textsearch: &colid{qual: "a", name: "ts_document"},
		},
	}, {
		name: "InsertAnalysisTestOK_OnConflict",
		want: &typespec{
			name:       "InsertAnalysisTestOK_OnConflict",
			kind:       speckindInsert,
			rel:        reldummyslice,
			onconflict: &onconflictblock{ignore: true},
		},
	}, {
		name: "InsertAnalysisTestOK_OnConflictColumn",
		want: &typespec{
			name: "InsertAnalysisTestOK_OnConflictColumn",
			kind: speckindInsert,
			rel:  reldummyslice,
			onconflict: &onconflictblock{
				column: []colid{{qual: "a", name: "id"}},
				ignore: true,
			},
		},
	}, {
		name: "InsertAnalysisTestOK_OnConflictConstraint",
		want: &typespec{
			name: "InsertAnalysisTestOK_OnConflictConstraint",
			kind: speckindInsert,
			rel:  reldummyslice,
			onconflict: &onconflictblock{
				constraint: "relation_constraint_xyz",
				update: &collist{items: []colid{
					{qual: "a", name: "foo"},
					{qual: "a", name: "bar"},
					{qual: "a", name: "baz"},
				}},
			},
		},
	}, {
		name: "InsertAnalysisTestOK_OnConflictIndex",
		want: &typespec{
			name: "InsertAnalysisTestOK_OnConflictIndex",
			kind: speckindInsert,
			rel:  reldummyslice,
			onconflict: &onconflictblock{
				index:  "relation_index_xyz",
				update: &collist{all: true},
			},
		},
	}, {
		name: "DeleteAnalysisTestOK_ResultField",
		want: &typespec{
			name: "DeleteAnalysisTestOK_ResultField",
			kind: speckindDelete,
			rel: &relfield{
				name:  "_",
				relid: relid{name: "relation_a", alias: "a"},
				isdir: true,
			},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &wherecolumn{colid: colid{qual: "a", name: "is_inactive"}, cmp: cmpistrue}},
			}},
			result: &resultfield{
				name: "Result",
				rec:  reldummyslice.rec,
			},
		},
	}, {
		name: "DeleteAnalysisTestOK_RowsAffected",
		want: &typespec{
			name: "DeleteAnalysisTestOK_RowsAffected",
			kind: speckindDelete,
			rel: &relfield{
				name:  "_",
				relid: relid{name: "relation_a", alias: "a"},
				isdir: true,
			},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &wherecolumn{colid: colid{qual: "a", name: "is_inactive"}, cmp: cmpistrue}},
			}},
			rowsaffected: "RowsAffected",
		},
	}, {
		name: "SelectAnalysisTestOK_FilterField",
		want: &typespec{
			name:   "SelectAnalysisTestOK_FilterField",
			kind:   speckindSelect,
			rel:    reldummyslice,
			filter: "Filter",
		},
	}, {
		name: "SelectAnalysisTestOK_FieldTypesBasic",
		want: &typespec{
			name: "SelectAnalysisTestOK_FieldTypesBasic",
			kind: speckindSelect,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a", alias: "a"},
				rec: recordtype{
					base: typeinfo{kind: kindstruct},
					fields: []*fieldinfo{{
						name: "f1", typ: typeinfo{kind: kindbool},
						colid: colid{name: "c1"},
						tag:   tagutil.Tag{"sql": {"c1"}},
					}, {
						name: "f2", typ: typeinfo{kind: kinduint8, isbyte: true},
						colid: colid{name: "c2"},
						tag:   tagutil.Tag{"sql": {"c2"}},
					}, {
						name: "f3", typ: typeinfo{kind: kindint32, isrune: true},
						colid: colid{name: "c3"},
						tag:   tagutil.Tag{"sql": {"c3"}},
					}, {
						name: "f4", typ: typeinfo{kind: kindint8},
						colid: colid{name: "c4"},
						tag:   tagutil.Tag{"sql": {"c4"}},
					}, {
						name: "f5", typ: typeinfo{kind: kindint16},
						colid: colid{name: "c5"},
						tag:   tagutil.Tag{"sql": {"c5"}},
					}, {
						name: "f6", typ: typeinfo{kind: kindint32},
						colid: colid{name: "c6"},
						tag:   tagutil.Tag{"sql": {"c6"}},
					}, {
						name: "f7", typ: typeinfo{kind: kindint64},
						colid: colid{name: "c7"},
						tag:   tagutil.Tag{"sql": {"c7"}},
					}, {
						name: "f8", typ: typeinfo{kind: kindint},
						colid: colid{name: "c8"},
						tag:   tagutil.Tag{"sql": {"c8"}},
					}, {
						name: "f9", typ: typeinfo{kind: kinduint8},
						colid: colid{name: "c9"},
						tag:   tagutil.Tag{"sql": {"c9"}},
					}, {
						name: "f10", typ: typeinfo{kind: kinduint16},
						colid: colid{name: "c10"},
						tag:   tagutil.Tag{"sql": {"c10"}},
					}, {
						name: "f11", typ: typeinfo{kind: kinduint32},
						colid: colid{name: "c11"},
						tag:   tagutil.Tag{"sql": {"c11"}},
					}, {
						name: "f12", typ: typeinfo{kind: kinduint64},
						colid: colid{name: "c12"},
						tag:   tagutil.Tag{"sql": {"c12"}},
					}, {
						name: "f13", typ: typeinfo{kind: kinduint},
						colid: colid{name: "c13"},
						tag:   tagutil.Tag{"sql": {"c13"}},
					}, {
						name: "f14", typ: typeinfo{kind: kinduintptr},
						colid: colid{name: "c14"},
						tag:   tagutil.Tag{"sql": {"c14"}},
					}, {
						name: "f15", typ: typeinfo{kind: kindfloat32},
						colid: colid{name: "c15"},
						tag:   tagutil.Tag{"sql": {"c15"}},
					}, {
						name: "f16", typ: typeinfo{kind: kindfloat64},
						colid: colid{name: "c16"},
						tag:   tagutil.Tag{"sql": {"c16"}},
					}, {
						name: "f17", typ: typeinfo{kind: kindcomplex64},
						colid: colid{name: "c17"},
						tag:   tagutil.Tag{"sql": {"c17"}},
					}, {
						name: "f18", typ: typeinfo{kind: kindcomplex128},
						colid: colid{name: "c18"},
						tag:   tagutil.Tag{"sql": {"c18"}},
					}, {
						name: "f19", typ: typeinfo{kind: kindstring},
						colid: colid{name: "c19"},
						tag:   tagutil.Tag{"sql": {"c19"}},
					}},
				},
			},
		},
	}, {
		name: "SelectAnalysisTestOK_FieldTypesSlices",
		want: &typespec{
			name: "SelectAnalysisTestOK_FieldTypesSlices",
			kind: speckindSelect,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a", alias: "a"},
				rec: recordtype{
					base: typeinfo{kind: kindstruct},
					fields: []*fieldinfo{{
						name: "f1", typ: typeinfo{
							kind: kindslice,
							elem: &typeinfo{kind: kindbool},
						},
						colid: colid{name: "c1"},
						tag:   tagutil.Tag{"sql": {"c1"}},
					}, {
						name: "f2", typ: typeinfo{
							kind: kindslice,
							elem: &typeinfo{kind: kinduint8, isbyte: true},
						},
						colid: colid{name: "c2"},
						tag:   tagutil.Tag{"sql": {"c2"}},
					}, {
						name: "f3", typ: typeinfo{
							kind: kindslice,
							elem: &typeinfo{kind: kindint32, isrune: true},
						},
						colid: colid{name: "c3"},
						tag:   tagutil.Tag{"sql": {"c3"}},
					}, {
						name: "f4", typ: typeinfo{
							name:       "HardwareAddr",
							kind:       kindslice,
							pkgpath:    "net",
							pkgname:    "net",
							pkglocal:   "net",
							isimported: true,
							elem:       &typeinfo{kind: kinduint8, isbyte: true},
						},
						colid: colid{name: "c4"},
						tag:   tagutil.Tag{"sql": {"c4"}},
					}, {
						name: "f5", typ: typeinfo{
							name:            "RawMessage",
							kind:            kindslice,
							pkgpath:         "encoding/json",
							pkgname:         "json",
							pkglocal:        "json",
							isimported:      true,
							isjsmarshaler:   true,
							isjsunmarshaler: true,
							elem:            &typeinfo{kind: kinduint8, isbyte: true},
						},
						colid: colid{name: "c5"},
						tag:   tagutil.Tag{"sql": {"c5"}},
					}, {
						name: "f6", typ: typeinfo{
							kind: kindslice,
							elem: &typeinfo{
								name:          "Marshaler",
								kind:          kindinterface,
								pkgpath:       "encoding/json",
								pkgname:       "json",
								pkglocal:      "json",
								isimported:    true,
								isjsmarshaler: true,
							},
						},
						colid: colid{name: "c6"},
						tag:   tagutil.Tag{"sql": {"c6"}},
					}, {
						name: "f7", typ: typeinfo{
							kind: kindslice,
							elem: &typeinfo{
								name:            "RawMessage",
								kind:            kindslice,
								pkgpath:         "encoding/json",
								pkgname:         "json",
								pkglocal:        "json",
								isimported:      true,
								isjsmarshaler:   true,
								isjsunmarshaler: true,
								elem:            &typeinfo{kind: kinduint8, isbyte: true},
							},
						},
						colid: colid{name: "c7"},
						tag:   tagutil.Tag{"sql": {"c7"}},
					}, {
						name: "f8", typ: typeinfo{
							kind: kindslice,
							elem: &typeinfo{
								kind: kindslice,
								elem: &typeinfo{kind: kinduint8, isbyte: true},
							},
						},
						colid: colid{name: "c8"},
						tag:   tagutil.Tag{"sql": {"c8"}},
					}, {
						name: "f9", typ: typeinfo{
							kind: kindslice,
							elem: &typeinfo{
								kind:     kindarray,
								arraylen: 2,
								elem: &typeinfo{
									kind:     kindarray,
									arraylen: 2,
									elem:     &typeinfo{kind: kindfloat64},
								},
							},
						},
						colid: colid{name: "c9"},
						tag:   tagutil.Tag{"sql": {"c9"}},
					}, {
						name: "f10", typ: typeinfo{
							kind: kindslice,
							elem: &typeinfo{
								kind: kindslice,
								elem: &typeinfo{
									kind:     kindarray,
									arraylen: 2,
									elem:     &typeinfo{kind: kindfloat64},
								},
							},
						},
						colid: colid{name: "c10"},
						tag:   tagutil.Tag{"sql": {"c10"}},
					}, {
						name: "f11", typ: typeinfo{
							kind: kindmap,
							key:  &typeinfo{kind: kindstring},
							elem: &typeinfo{
								name:       "NullString",
								kind:       kindstruct,
								pkgpath:    "database/sql",
								pkgname:    "sql",
								pkglocal:   "sql",
								isimported: true,
								isscanner:  true,
								isvaluer:   true,
							},
						},
						colid: colid{name: "c11"},
						tag:   tagutil.Tag{"sql": {"c11"}},
					}, {
						name: "f12", typ: typeinfo{
							kind: kindslice,
							elem: &typeinfo{
								kind: kindmap,
								key:  &typeinfo{kind: kindstring},
								elem: &typeinfo{
									kind: kindptr,
									elem: &typeinfo{kind: kindstring},
								},
							},
						},
						colid: colid{name: "c12"},
						tag:   tagutil.Tag{"sql": {"c12"}},
					}, {
						name: "f13", typ: typeinfo{
							kind: kindslice,
							elem: &typeinfo{
								kind:     kindarray,
								arraylen: 2,
								elem: &typeinfo{
									kind: kindptr,
									elem: &typeinfo{
										name:            "Int",
										kind:            kindstruct,
										pkgpath:         "math/big",
										pkgname:         "big",
										pkglocal:        "big",
										isimported:      true,
										isjsmarshaler:   true,
										isjsunmarshaler: true,
									},
								},
							},
						},
						colid: colid{name: "c13"},
						tag:   tagutil.Tag{"sql": {"c13"}},
					}},
				},
			},
		},
	}, {
		name: "SelectAnalysisTestOK_FieldTypesInterfaces",
		want: &typespec{
			name: "SelectAnalysisTestOK_FieldTypesInterfaces",
			kind: speckindSelect,
			rel: &relfield{
				name:  "Rel",
				relid: relid{name: "relation_a", alias: "a"},
				rec: recordtype{
					base: typeinfo{kind: kindstruct},
					fields: []*fieldinfo{{
						name: "f1", typ: typeinfo{
							name:          "Marshaler",
							kind:          kindinterface,
							pkgpath:       "encoding/json",
							pkgname:       "json",
							pkglocal:      "json",
							isimported:    true,
							isjsmarshaler: true,
						},
						colid: colid{name: "c1"},
						tag:   tagutil.Tag{"sql": {"c1"}},
					}, {
						name: "f2", typ: typeinfo{
							name:            "Unmarshaler",
							kind:            kindinterface,
							pkgpath:         "encoding/json",
							pkgname:         "json",
							pkglocal:        "json",
							isimported:      true,
							isjsunmarshaler: true,
						},
						colid: colid{name: "c2"},
						tag:   tagutil.Tag{"sql": {"c2"}},
					}, {
						name: "f3", typ: typeinfo{
							kind:            kindinterface,
							isjsmarshaler:   true,
							isjsunmarshaler: true,
						},
						colid: colid{name: "c3"},
						tag:   tagutil.Tag{"sql": {"c3"}},
					}},
				},
			},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := runAnalysis(tt.name, t)
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

	spec, err := runAnalysis("SelectAnalysisTestOK_typeinfo_string", t)
	if err != nil {
		t.Error(err)
	}
	fields := spec.rel.rec.fields
	for i := 0; i < len(fields); i++ {
		ff := fields[i]
		tt := tests[i]

		got := ff.typ.string(true)
		if ff.name != tt.name || got != tt.want {
			t.Errorf("got %s::%s, want %s::%s", ff.name, got, tt.name, tt.want)
		}
	}
}
