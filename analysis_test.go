package gosql

import (
	"testing"

	"github.com/frk/compare"
	"github.com/frk/gosql/internal/errors"
	"github.com/frk/gosql/internal/testutil"
	"github.com/frk/tagutil"
)

var tdata = testutil.ParseTestdata("testdata")

func runAnalysis(name string, t *testing.T) (*command, error) {
	named := testutil.FindNamedType(name, tdata)
	if named == nil {
		// Stop the test if no type with the given name was found.
		t.Fatal(name, " not found")
		return nil, nil
	}
	return analyze(named)
}

func TestAnalysis_InsertCommand(t *testing.T) {

	// for reuse, analyzed common.User typeinfo
	commonUserTypeinfo := typeinfo{
		name:       "User",
		kind:       kindstruct,
		pkgpath:    "github.com/frk/gosql/testdata/common",
		pkgname:    "common",
		pkglocal:   "common",
		isimported: true,
		fields: []*fieldinfo{{
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
		}},
	}

	reldummyslice := &relfield{
		field: "Rel",
		relid: relid{name: "relation_a", alias: "a"},
		datatype: datatype{
			base: typeinfo{
				name:     "T",
				kind:     kindstruct,
				pkgpath:  "path/to/test",
				pkgname:  "testdata",
				pkglocal: "testdata",
			},
			isslice: true,
		},
	}

	dummytype := datatype{
		base: typeinfo{
			name:     "T",
			kind:     kindstruct,
			pkgpath:  "path/to/test",
			pkgname:  "testdata",
			pkglocal: "testdata",
		},
	}

	tests := []struct {
		name string
		want *command
		err  error
	}{{
		name: "InsertTestBAD_NoRelfield",
		err:  errors.NoRelfieldError,
	}, {
		name: "InsertTestBAD3",
		err:  errors.BadRelfieldTypeError,
	}, {
		name: "DeleteTestBAD_BadRelId",
		err:  errors.BadRelIdError,
	}, {
		name: "DeleteTestBAD_IllegalCountField",
		err:  errors.IllegalCountFieldError,
	}, {
		name: "UpdateTestBAD_IllegalExistsField",
		err:  errors.IllegalExistsFieldError,
	}, {
		name: "InsertTestBAD_IllegalNotExistsField",
		err:  errors.IllegalNotExistsFieldError,
	}, {
		name: "SelectTestBAD_IllegalRelationDirective",
		err:  errors.IllegalRelationDirectiveError,
	}, {
		name: "SelectTestBAD_UnnamedBaseStructType",
		err:  errors.BadRelfieldTypeError,
	}, {
		name: "SelectTestBAD_IllegalAllDirective",
		err:  errors.IllegalAllDirectiveError,
	}, {
		name: "InsertTestBAD_IllegalAllDirective",
		err:  errors.IllegalAllDirectiveError,
	}, {
		name: "UpdateTestBAD_ConflictWhereProducer",
		err:  errors.ConflictWhereProducerError,
	}, {
		name: "DeleteTestBAD_IllegalDefaultDirective",
		err:  errors.IllegalDefaultDirectiveError,
	}, {
		name: "UpdateTestBAD_EmptyDefaultDirectiveCollist",
		err:  errors.EmptyColListError,
	}, {
		name: "SelectTestBAD_IllegalForceDirective",
		err:  errors.IllegalForceDirectiveError,
	}, {
		name: "UpdateTestBAD_BadForceDirectiveColId",
		err:  errors.BadColIdError,
	}, {
		name: "FilterTestBAD_IllegalReturnDirective",
		err:  errors.IllegalReturnDirectiveError,
	}, {
		name: "DeleteTestBAD_ConflictResultProducer",
		err:  errors.ConflictResultProducerError,
	}, {
		name: "UpdateTestBAD_EmptyReturnDirectiveCollist",
		err:  errors.EmptyColListError,
	}, {
		name: "InsertTestBAD_IllegalLimitField",
		err:  errors.IllegalLimitFieldOrDirectiveError,
	}, {
		name: "UpdateTestBAD_IllegalOffsetField",
		err:  errors.IllegalOffsetFieldOrDirectiveError,
	}, {
		name: "InsertTestBAD_IllegalOrderByDirective",
		err:  errors.IllegalOrderByDirectiveError,
	}, {
		name: "DeleteTestBAD_IllegalOverrideDirective",
		err:  errors.IllegalOverrideDirectiveError,
	}, {
		name: "SelectTestBAD_IllegalTextSearchDirective",
		err:  errors.IllegalTextSearchDirectiveError,
	}, {
		name: "SelectTestBAD_IllegalColumnDirective",
		err:  errors.IllegalCommandDirectiveError,
	}, {
		name: "InsertTestBAD_IllegalWhereBlock",
		err:  errors.IllegalWhereBlockError,
	}, {
		name: "UpdateTestBAD_IllegalJoinBlock",
		err:  errors.IllegalJoinBlockError,
	}, {
		name: "DeleteTestBAD_IllegalFromBlock",
		err:  errors.IllegalFromBlockError,
	}, {
		name: "SelectTestBAD_IllegalUsingBlock",
		err:  errors.IllegalUsingBlockError,
	}, {
		name: "UpdateTestBAD_IllegalOnConflictBlock",
		err:  errors.IllegalOnConflictBlockError,
	}, {
		name: "SelectTestBAD_IllegalResultField",
		err:  errors.IllegalResultFieldError,
	}, {
		name: "SelectTestBAD_ConflictLimitProducer",
		err:  errors.ConflictLimitProducerError,
	}, {
		name: "SelectTestBAD_ConflictOffsetProducer",
		err:  errors.ConflictOffsetProducerError,
	}, {
		name: "SelectTestBAD_IllegalRowsAffectedField",
		err:  errors.IllegalRowsAffectedFieldError,
	}, {
		name: "InsertTestBAD_IllegalFilterField",
		err:  errors.IllegalFilterFieldError,
	}, {
		name: "SelectTestBAD_ConflictWhereProducer",
		err:  errors.ConflictWhereProducerError,
	}, {
		name: "DeleteTestBAD_ConflictWhereProducer",
		err:  errors.ConflictErrorHandlerFieldError,
	}, {
		name: "SelectTestBAD_IteratorWithTooManyMethods",
		err:  errors.BadIteratorTypeError,
	}, {
		name: "SelectTestBAD_IteratorWithBadSignature",
		err:  errors.BadIteratorTypeError,
	}, {
		name: "SelectTestBAD_IteratorWithBadSignatureIface",
		err:  errors.BadIteratorTypeError,
	}, {
		name: "SelectTestBAD_IteratorWithUnexportedMethod",
		err:  errors.BadIteratorTypeError,
	}, {
		name: "SelectTestBAD_IteratorWithUnnamedArgument",
		err:  errors.BadIteratorTypeError,
	}, {
		name: "SelectTestBAD_IteratorWithNonStructArgument",
		err:  errors.BadIteratorTypeError,
	}, {
		name: "InsertTestBAD_BadRelfiedlStructBaseType",
		err:  errors.BadRelfieldTypeError,
	}, {
		name: "UpdateTestBAD_BadRelTypeFieldColId",
		err:  errors.BadColIdError,
	}, {
		name: "UpdateTestBAD_ConflictWhereProducer2",
		err:  errors.ConflictWhereProducerError,
	}, {
		name: "DeleteTestBAD_BadWhereBlockType",
		err:  errors.BadWhereBlockTypeError,
	}, {
		name: "SelectTestBAD_BadBoolTagValue",
		err:  errors.BadBoolTagValueError,
	}, {
		name: "SelectTestBAD_BadNestedWhereBlockType",
		err:  errors.BadWhereBlockTypeError,
	}, {
		name: "SelectTestBAD_BadColumnExpressionLHS",
		err:  errors.BadColIdError,
	}, {
		name: "SelectTestBAD_BadColumnCmpopCombo",
		err:  errors.BadCmpopComboError,
	}, {
		name: "DeleteTestBAD_BadColumnExpressionLHS",
		err:  errors.BadColIdError,
	}, {
		name: "UpdateTestBAD_BadUnaryOp",
		err:  errors.BadUnaryCmpopError,
	}, {
		name: "UpdateTestBAD_ExtraScalarrop",
		err:  errors.ExtraScalarropError,
	}, {
		name: "SelectTestBAD_BadBetweenFieldType",
		err:  errors.BadBetweenTypeError,
	}, {
		name: "SelectTestBAD_BadBetweenFieldType2",
		err:  errors.BadBetweenTypeError,
	}, {
		name: "SelectTestBAD_BadBetweenArgColId",
		err:  errors.BadColIdError,
	}, {
		name: "SelectTestBAD_NoBetweenXYArg",
		err:  errors.NoBetweenXYArgsError,
	}, {
		name: "SelectTestBAD_BadBetweenColId",
		err:  errors.BadColIdError,
	}, {
		name: "DeleteTestBAD_BadWhereFieldColId",
		err:  errors.BadColIdError,
	}, {
		name: "DeleteTestBAD_BadWhereFieldCmpopCombo",
		err:  errors.BadCmpopComboError,
	}, {
		name: "UpdateTestBAD_BadWhereFieldTypeForScalarrop",
		err:  errors.BadScalarFieldTypeError,
	}, {
		name: "SelectTestBAD_BadJoinBlockType",
		err:  errors.BadJoinBlockTypeError,
	}, {
		name: "SelectTestBAD_IllegalJoinBlockRelationDirective",
		err:  errors.IllegalJoinBlockRelationDirectiveError,
	}, {
		name: "DeleteTestBAD_ConflictRelationDirective",
		err:  errors.ConflictJoinBlockRelationDirectiveError,
	}, {
		name: "UpdateTestBAD_BadFromRelationRelId",
		err:  errors.BadRelIdError,
	}, {
		name: "SelectTestBAD_BadJoinDirectiveRelId",
		err:  errors.BadRelIdError,
	}, {
		name: "SelectTestBAD_BadJoinDirectiveExpressionColId",
		err:  errors.BadColIdError,
	}, {
		name: "SelectTestBAD_BadJoinDirectiveExpressionCmpop",
		err:  errors.BadUnaryCmpopError,
	}, {
		name: "SelectTestBAD_BadJoinDirectiveExpressionExtraScalarrop",
		err:  errors.ExtraScalarropError,
	}, {
		name: "SelectTestBAD_BadJoinDirectiveExpressionCmpopCombo",
		err:  errors.BadCmpopComboError,
	}, {
		name: "DeleteTestBAD_IllegalJoinBlockDirective",
		err:  errors.IllegalJoinBlockDirectiveError,
	}, {
		name: "InsertTestBAD_BadOnConflictBlockType",
		err:  errors.BadOnConflictBlockTypeError,
	}, {
		name: "InsertTestBAD_ConflictOnConflictBlockTargetProducer",
		err:  errors.ConflictOnConflictBlockTargetProducerError,
	}, {
		name: "InsertTestBAD_ConflictOnConflictBlockTargetProducer2",
		err:  errors.ConflictOnConflictBlockTargetProducerError,
	}, {
		name: "InsertTestBAD_ConflictOnConflictBlockTargetProducer3",
		err:  errors.ConflictOnConflictBlockTargetProducerError,
	}, {
		name: "InsertTestBAD_ConflictOnConflictBlockActionProducer",
		err:  errors.ConflictOnConflictBlockActionProducerError,
	}, {
		name: "InsertTestBAD_ConflictOnConflictBlockActionProducer2",
		err:  errors.ConflictOnConflictBlockActionProducerError,
	}, {
		name: "InsertTestBAD_BadOnConflictColumnTargetValue",
		err:  errors.BadColIdError,
	}, {
		name: "InsertTestBAD_BadOnConflictIndexTargetIdent",
		err:  errors.BadIndexIdentifierValueError,
	}, {
		name: "InsertTestBAD_BadOnConflictConstraintTargetIdent",
		err:  errors.BadConstraintIdentifierValueError,
	}, {
		name: "InsertTestBAD_BadOnConflictUpdateActionCollist",
		err:  errors.BadColIdError,
	}, {
		name: "InsertTestBAD_IllegalOnConflictDirective",
		err:  errors.IllegalOnConflictBlockDirectiveError,
	}, {
		name: "InsertTestBAD_NoOnConflictTarget",
		err:  errors.NoOnConflictTargetError,
	}, {
		name: "SelectTestBAD_BadLimitFieldType",
		err:  errors.BadLimitTypeError,
	}, {
		name: "SelectTestBAD_NoLimitDirectiveValue",
		err:  errors.NoLimitDirectiveValueError,
	}, {
		name: "SelectTestBAD_BadLimitDirectiveValue",
		err:  errors.BadLimitValueError,
	}, {
		name: "SelectTestBAD_BadOffsetFieldType",
		err:  errors.BadOffsetTypeError,
	}, {
		name: "SelectTestBAD_NoOffsetDirectiveValue",
		err:  errors.NoOffsetDirectiveValueError,
	}, {
		name: "SelectTestBAD_BadOffsetDirectiveValue",
		err:  errors.BadOffsetValueError,
	}, {
		name: "SelectTestBAD_EmptyOrderByDirectiveCollist",
		err:  errors.EmptyOrderByListError,
	}, {
		name: "SelectTestBAD_BadOrderByDirectiveNullsOrderValue",
		err:  errors.BadNullsOrderOptionValueError,
	}, {
		name: "SelectTestBAD_BadOrderByDirectiveCollist",
		err:  errors.BadColIdError,
	}, {
		name: "InsertTestBAD_BadOverrideDirectiveKindValue",
		err:  errors.BadOverrideKindValueError,
	}, {
		name: "UpdateTestBAD_ConflictResultProducer",
		err:  errors.ConflictResultProducerError,
	}, {
		name: "UpdateTestBAD_BadResultFieldType",
		err:  errors.BadRelfieldTypeError,
	}, {
		name: "DeleteTestBAD_ConflictResultProducer2",
		err:  errors.ConflictResultProducerError,
	}, {
		name: "DeleteTestBAD_BadRowsAffecteFieldType",
		err:  errors.BadRowsAffectedTypeError,
	}, {
		name: "FilterTestBAD_BadTextSearchDirectiveColId",
		err:  errors.BadColIdError,
	}, {
		name: "InsertTestOK1",
		want: &command{name: "InsertTestOK1", typ: cmdtypeInsert, rel: &relfield{
			field: "UserRec",
			relid: relid{name: "users_table"},
			datatype: datatype{
				base:      commonUserTypeinfo,
				ispointer: true,
			},
		}},
	}, {
		name: "InsertTestOK2",
		want: &command{name: "InsertTestOK2", typ: cmdtypeInsert, rel: &relfield{
			field: "UserRec",
			relid: relid{name: "users_table"},
			datatype: datatype{
				base: typeinfo{
					kind: kindstruct,
					fields: []*fieldinfo{{
						name:       "Name3",
						typ:        typeinfo{kind: kindstring},
						isexported: true,
						colid:      colid{name: "name"},
						tag:        tagutil.Tag{"sql": {"name"}},
					}},
				},
			},
		}},
	}, {
		name: "SelectTestOK3",
		want: &command{name: "SelectTestOK3", typ: cmdtypeSelect, rel: &relfield{
			field: "User",
			relid: relid{name: "users_table"},
			datatype: datatype{
				base:      commonUserTypeinfo,
				ispointer: true,
				isiter:    true,
			},
		}},
	}, {
		name: "SelectTestOK4",
		want: &command{name: "SelectTestOK4", typ: cmdtypeSelect, rel: &relfield{
			field: "User",
			relid: relid{name: "users_table"},
			datatype: datatype{
				base:      commonUserTypeinfo,
				ispointer: true,
				isiter:    true,
			},
		}},
	}, {
		name: "SelectTestOK5",
		want: &command{name: "SelectTestOK5", typ: cmdtypeSelect, rel: &relfield{
			field: "User",
			relid: relid{name: "users_table"},
			datatype: datatype{
				base:       commonUserTypeinfo,
				ispointer:  true,
				isiter:     true,
				itermethod: "Fn",
			},
		}},
	}, {
		name: "SelectTestOK6",
		want: &command{name: "SelectTestOK6", typ: cmdtypeSelect, rel: &relfield{
			field: "User",
			relid: relid{name: "users_table"},
			datatype: datatype{
				base:       commonUserTypeinfo,
				ispointer:  true,
				isiter:     true,
				itermethod: "Fn",
			},
		}},
	}, {
		name: "SelectTestOK7",
		want: &command{name: "SelectTestOK7", typ: cmdtypeSelect, rel: &relfield{
			field: "Rel",
			relid: relid{name: "relation_a"},
			datatype: datatype{
				base: typeinfo{
					kind: kindstruct,
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
			},
		}},
	}, {
		name: "InsertTestOK8",
		want: &command{name: "InsertTestOK8", typ: cmdtypeInsert, rel: &relfield{
			field: "Rel",
			relid: relid{name: "relation_a"},
			datatype: datatype{
				base: typeinfo{
					kind: kindstruct,
					fields: []*fieldinfo{{
						name:       "Foobar",
						isexported: true,
						typ: typeinfo{
							name:       "Foo",
							kind:       kindstruct,
							pkgpath:    "github.com/frk/gosql/testdata/common",
							pkgname:    "common",
							pkglocal:   "common",
							isimported: true,
							fields: []*fieldinfo{{
								name:       "Bar",
								isexported: true,
								typ: typeinfo{
									name:       "Bar",
									kind:       kindstruct,
									pkgpath:    "github.com/frk/gosql/testdata/common",
									pkgname:    "common",
									pkglocal:   "common",
									isimported: true,
									fields: []*fieldinfo{{
										name:       "Baz",
										isexported: true,
										isembedded: true,
										typ: typeinfo{
											name:       "Baz",
											kind:       kindstruct,
											pkgpath:    "github.com/frk/gosql/testdata/common",
											pkgname:    "common",
											pkglocal:   "common",
											isimported: true,
											fields: []*fieldinfo{{
												name:       "Val",
												isexported: true,
												typ:        typeinfo{kind: kindstring},
												colid:      colid{name: "foo_bar_baz_val"},
												tag:        tagutil.Tag{"sql": {"val"}},
											}},
										},
										tag: tagutil.Tag{"sql": {">baz_"}},
									}},
								},
								tag: tagutil.Tag{"sql": {">bar_"}},
							}},
						},
						tag: tagutil.Tag{"sql": {">foo_"}},
					}},
				},
			},
		}},
	}, {
		name: "DeleteTestOK9",
		want: &command{
			name: "DeleteTestOK9",
			typ:  cmdtypeDelete,
			rel: &relfield{
				field:    "Rel",
				relid:    relid{name: "relation_a"},
				datatype: datatype{base: typeinfo{kind: kindstruct}},
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
		name: "DeleteTestOK10",
		want: &command{
			name: "DeleteTestOK10",
			typ:  cmdtypeDelete,
			rel: &relfield{
				field:    "Rel",
				relid:    relid{name: "relation_a"},
				datatype: datatype{base: typeinfo{kind: kindstruct}},
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
			}},
		},
	}, {
		name: "DeleteTestOK11",
		want: &command{
			name: "DeleteTestOK11",
			typ:  cmdtypeDelete,
			rel: &relfield{
				field:    "Rel",
				relid:    relid{name: "relation_a"},
				datatype: datatype{base: typeinfo{kind: kindstruct}},
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
		name: "DeleteTestOK12",
		want: &command{
			name: "DeleteTestOK12",
			typ:  cmdtypeDelete,
			rel: &relfield{
				field:    "Rel",
				relid:    relid{name: "relation_a"},
				datatype: datatype{base: typeinfo{kind: kindstruct}},
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
		name: "DeleteTestOK13",
		want: &command{
			name: "DeleteTestOK13",
			typ:  cmdtypeDelete,
			rel: &relfield{
				field:    "Rel",
				relid:    relid{name: "relation_a"},
				datatype: datatype{base: typeinfo{kind: kindstruct}},
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
		name: "DeleteTestOK14",
		want: &command{
			name: "DeleteTestOK14",
			typ:  cmdtypeDelete,
			rel: &relfield{
				field:    "Rel",
				relid:    relid{name: "relation_a"},
				datatype: datatype{base: typeinfo{kind: kindstruct}},
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
		name: "DeleteTestOK_DistinctFrom",
		want: &command{
			name: "DeleteTestOK_DistinctFrom",
			typ:  cmdtypeDelete,
			rel: &relfield{
				field:    "Rel",
				relid:    relid{name: "relation_a"},
				datatype: datatype{base: typeinfo{kind: kindstruct}},
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
		name: "DeleteTestOK_ArrayComparisons",
		want: &command{
			name: "DeleteTestOK_ArrayComparisons",
			typ:  cmdtypeDelete,
			rel: &relfield{
				field:    "Rel",
				relid:    relid{name: "relation_a"},
				datatype: datatype{base: typeinfo{kind: kindstruct}},
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
		name: "DeleteTestOK_PatternMatching",
		want: &command{
			name: "DeleteTestOK_PatternMatching",
			typ:  cmdtypeDelete,
			rel: &relfield{
				field:    "Rel",
				relid:    relid{name: "relation_a"},
				datatype: datatype{base: typeinfo{kind: kindstruct}},
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
		name: "DeleteTestOK_Using",
		want: &command{
			name: "DeleteTestOK_Using",
			typ:  cmdtypeDelete,
			rel: &relfield{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{base: typeinfo{kind: kindstruct}},
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
		name: "UpdateTestOK_From",
		want: &command{
			name: "UpdateTestOK_From",
			typ:  cmdtypeUpdate,
			rel: &relfield{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{base: typeinfo{kind: kindstruct}},
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
		name: "SelectTestOK_Join",
		want: &command{
			name: "SelectTestOK_Join",
			typ:  cmdtypeSelect,
			rel: &relfield{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{base: typeinfo{kind: kindstruct}},
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
		name: "UpdateTestOK_All",
		want: &command{
			name: "UpdateTestOK_All",
			typ:  cmdtypeUpdate,
			rel: &relfield{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{base: typeinfo{kind: kindstruct}},
			},
			all: true,
		},
	}, {
		name: "DeleteTestOK_All",
		want: &command{
			name: "DeleteTestOK_All",
			typ:  cmdtypeDelete,
			rel: &relfield{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{base: typeinfo{kind: kindstruct}},
			},
			all: true,
		},
	}, {
		name: "DeleteTestOK_Return",
		want: &command{
			name: "DeleteTestOK_Return",
			typ:  cmdtypeDelete,
			rel: &relfield{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{base: typeinfo{kind: kindstruct}},
			},
			returning: &collist{all: true},
		},
	}, {
		name: "InsertTestOK_Return",
		want: &command{
			name: "InsertTestOK_Return",
			typ:  cmdtypeInsert,
			rel: &relfield{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{base: typeinfo{kind: kindstruct}},
			},
			returning: &collist{items: []colid{
				{qual: "a", name: "foo"},
				{qual: "a", name: "bar"},
				{qual: "a", name: "baz"}}},
		},
	}, {
		name: "UpdateTestOK_Return",
		want: &command{
			name: "UpdateTestOK_Return",
			typ:  cmdtypeUpdate,
			rel: &relfield{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{base: typeinfo{kind: kindstruct}},
			},
			returning: &collist{items: []colid{
				{qual: "a", name: "foo"},
				{qual: "a", name: "bar"},
				{qual: "a", name: "baz"}}},
		},
	}, {
		name: "InsertTestOK_Default",
		want: &command{
			name: "InsertTestOK_Default",
			typ:  cmdtypeInsert,
			rel: &relfield{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{base: typeinfo{kind: kindstruct}},
			},
			defaults: &collist{all: true},
		},
	}, {
		name: "UpdateTestOK_Default",
		want: &command{
			name: "UpdateTestOK_Default",
			typ:  cmdtypeUpdate,
			rel: &relfield{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{base: typeinfo{kind: kindstruct}},
			},
			defaults: &collist{items: []colid{
				{qual: "a", name: "foo"},
				{qual: "a", name: "bar"},
				{qual: "a", name: "baz"}}},
		},
	}, {
		name: "InsertTestOK_Force",
		want: &command{
			name: "InsertTestOK_Force",
			typ:  cmdtypeInsert,
			rel: &relfield{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{base: typeinfo{kind: kindstruct}},
			},
			force: &collist{all: true},
		},
	}, {
		name: "UpdateTestOK_Force",
		want: &command{
			name: "UpdateTestOK_Force",
			typ:  cmdtypeUpdate,
			rel: &relfield{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{base: typeinfo{kind: kindstruct}},
			},
			force: &collist{items: []colid{
				{qual: "a", name: "foo"},
				{qual: "a", name: "bar"},
				{qual: "a", name: "baz"}}},
		},
	}, {
		name: "SelectTestOK_ErrorHandler",
		want: &command{
			name: "SelectTestOK_ErrorHandler",
			typ:  cmdtypeSelect,
			rel: &relfield{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{base: typeinfo{kind: kindstruct}},
			},
			erh: "eh",
		},
	}, {
		name: "InsertTestOK_ErrorHandler",
		want: &command{
			name: "InsertTestOK_ErrorHandler",
			typ:  cmdtypeInsert,
			rel: &relfield{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{base: typeinfo{kind: kindstruct}},
			},
			erh: "myerrorhandler",
		},
	}, {
		name: "SelectTestOK_Count",
		want: &command{
			name: "SelectTestOK_Count",
			typ:  cmdtypeSelect,
			rel: &relfield{
				field: "Count",
				relid: relid{name: "relation_a", alias: "a"},
			},
			sel: selcount,
		},
	}, {
		name: "SelectTestOK_Exists",
		want: &command{
			name: "SelectTestOK_Exists",
			typ:  cmdtypeSelect,
			rel: &relfield{
				field: "Exists",
				relid: relid{name: "relation_a", alias: "a"},
			},
			sel: selexists,
		},
	}, {
		name: "SelectTestOK_NotExists",
		want: &command{
			name: "SelectTestOK_NotExists",
			typ:  cmdtypeSelect,
			rel: &relfield{
				field: "NotExists",
				relid: relid{name: "relation_a", alias: "a"},
			},
			sel: selnotexists,
		},
	}, {
		name: "DeleteTestOK_Relation",
		want: &command{
			name: "DeleteTestOK_Relation",
			typ:  cmdtypeDelete,
			rel: &relfield{
				field: "_",
				relid: relid{name: "relation_a", alias: "a"},
				isdir: true,
			},
		},
	}, {
		name: "SelectTestOK_LimitDirective",
		want: &command{
			name:  "SelectTestOK_LimitDirective",
			typ:   cmdtypeSelect,
			rel:   reldummyslice,
			limit: &limitvar{value: 25},
		},
	}, {
		name: "SelectTestOK_LimitField",
		want: &command{
			name:  "SelectTestOK_LimitField",
			typ:   cmdtypeSelect,
			rel:   reldummyslice,
			limit: &limitvar{value: 10, field: "Limit"},
		},
	}, {
		name: "SelectTestOK_OffsetDirective",
		want: &command{
			name:   "SelectTestOK_OffsetDirective",
			typ:    cmdtypeSelect,
			rel:    reldummyslice,
			offset: &offsetvar{value: 25},
		},
	}, {
		name: "SelectTestOK_OffsetField",
		want: &command{
			name:   "SelectTestOK_OffsetField",
			typ:    cmdtypeSelect,
			rel:    reldummyslice,
			offset: &offsetvar{value: 10, field: "Offset"},
		},
	}, {
		name: "SelectTestOK_OrderByDirective",
		want: &command{
			name: "SelectTestOK_OrderByDirective",
			typ:  cmdtypeSelect,
			rel:  reldummyslice,
			orderby: &orderbylist{items: []*orderbyitem{
				{col: colid{qual: "a", name: "foo"}, dir: orderasc, nulls: nullsfirst},
				{col: colid{qual: "a", name: "bar"}, dir: orderdesc, nulls: nullsfirst},
				{col: colid{qual: "a", name: "baz"}, dir: orderdesc, nulls: 0},
				{col: colid{qual: "a", name: "quux"}, dir: orderasc, nulls: nullslast},
			}},
		},
	}, {
		name: "InsertTestOK_OverrideDirective",
		want: &command{
			name:     "InsertTestOK_OverrideDirective",
			typ:      cmdtypeInsert,
			rel:      reldummyslice,
			override: overridingsystem,
		},
	}, {
		name: "FilterTestOK_TextSearchDirective",
		want: &command{
			name: "FilterTestOK_TextSearchDirective",
			typ:  cmdtypeFilter,
			rel: &relfield{
				field:    "_",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: dummytype,
			},
			textsearch: &colid{qual: "a", name: "ts_document"},
		},
	}, {
		name: "InsertTestOK_OnConflict",
		want: &command{
			name:       "InsertTestOK_OnConflict",
			typ:        cmdtypeInsert,
			rel:        reldummyslice,
			onconflict: &onconflictblock{ignore: true},
		},
	}, {
		name: "InsertTestOK_OnConflictColumn",
		want: &command{
			name: "InsertTestOK_OnConflictColumn",
			typ:  cmdtypeInsert,
			rel:  reldummyslice,
			onconflict: &onconflictblock{
				column: []colid{{qual: "a", name: "id"}},
				ignore: true,
			},
		},
	}, {
		name: "InsertTestOK_OnConflictConstraint",
		want: &command{
			name: "InsertTestOK_OnConflictConstraint",
			typ:  cmdtypeInsert,
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
		name: "InsertTestOK_OnConflictIndex",
		want: &command{
			name: "InsertTestOK_OnConflictIndex",
			typ:  cmdtypeInsert,
			rel:  reldummyslice,
			onconflict: &onconflictblock{
				index:  "relation_index_xyz",
				update: &collist{all: true},
			},
		},
	}, {
		name: "DeleteTestOK_ResultField",
		want: &command{
			name: "DeleteTestOK_ResultField",
			typ:  cmdtypeDelete,
			rel: &relfield{
				field: "_",
				relid: relid{name: "relation_a", alias: "a"},
				isdir: true,
			},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &wherecolumn{colid: colid{qual: "a", name: "is_inactive"}, cmp: cmpistrue}},
			}},
			result: &resultfield{
				name:     "Result",
				datatype: reldummyslice.datatype,
			},
		},
	}, {
		name: "DeleteTestOK_RowsAffected",
		want: &command{
			name: "DeleteTestOK_RowsAffected",
			typ:  cmdtypeDelete,
			rel: &relfield{
				field: "_",
				relid: relid{name: "relation_a", alias: "a"},
				isdir: true,
			},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &wherecolumn{colid: colid{qual: "a", name: "is_inactive"}, cmp: cmpistrue}},
			}},
			rowsaffected: "RowsAffected",
		},
	}, {
		name: "SelectTestOK_FilterField",
		want: &command{
			name:   "SelectTestOK_FilterField",
			typ:    cmdtypeSelect,
			rel:    reldummyslice,
			filter: "Filter",
		},
	}, {
		name: "SelectTestOK_FieldTypesBasic",
		want: &command{
			name: "SelectTestOK_FieldTypesBasic",
			typ:  cmdtypeSelect,
			rel: &relfield{
				field: "Rel",
				relid: relid{name: "relation_a", alias: "a"},
				datatype: datatype{base: typeinfo{
					kind: kindstruct,
					fields: []*fieldinfo{{
						name: "f1", typ: typeinfo{kind: kindbool},
						colid: colid{name: "c1"},
						tag:   tagutil.Tag{"sql": {"c1"}},
					}, {
						name: "f2", typ: typeinfo{kind: kinduint8},
						colid: colid{name: "c2"},
						tag:   tagutil.Tag{"sql": {"c2"}},
					}, {
						name: "f3", typ: typeinfo{kind: kindint32},
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
				}},
			},
		},
	}, {
		name: "SelectTestOK_FieldTypesSlices",
		want: &command{
			name: "SelectTestOK_FieldTypesSlices",
			typ:  cmdtypeSelect,
			rel: &relfield{
				field: "Rel",
				relid: relid{name: "relation_a", alias: "a"},
				datatype: datatype{base: typeinfo{
					kind: kindstruct,
					fields: []*fieldinfo{{
						name: "f1", typ: typeinfo{
							kind: kindslice,
							elem: &typeinfo{
								kind: kindbool,
							},
						},
						colid: colid{name: "c1"},
						tag:   tagutil.Tag{"sql": {"c1"}},
					}, {
						name: "f2", typ: typeinfo{
							kind: kindslice,
							elem: &typeinfo{
								kind: kinduint8,
							},
						},
						colid: colid{name: "c2"},
						tag:   tagutil.Tag{"sql": {"c2"}},
					}, {
						name: "f3", typ: typeinfo{
							kind: kindslice,
							elem: &typeinfo{
								kind: kindint32,
							},
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
							elem: &typeinfo{
								kind: kinduint8,
							},
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
							elem: &typeinfo{
								kind: kinduint8,
							},
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
								elem: &typeinfo{
									kind: kinduint8,
								},
							},
						},
						colid: colid{name: "c7"},
						tag:   tagutil.Tag{"sql": {"c7"}},
					}, {
						name: "f8", typ: typeinfo{
							kind: kindslice,
							elem: &typeinfo{
								kind: kindslice,
								elem: &typeinfo{
									kind: kinduint8,
								},
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
									elem: &typeinfo{
										kind: kindfloat64,
									},
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
									elem: &typeinfo{
										kind: kindfloat64,
									},
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
				}},
			},
		},
	}, {
		name: "SelectTestOK_FieldTypesInterfaces",
		want: &command{
			name: "SelectTestOK_FieldTypesInterfaces",
			typ:  cmdtypeSelect,
			rel: &relfield{
				field: "Rel",
				relid: relid{name: "relation_a", alias: "a"},
				datatype: datatype{base: typeinfo{
					kind: kindstruct,
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
				}},
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
