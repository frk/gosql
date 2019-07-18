package gosql

import (
	"testing"

	"github.com/frk/compare"
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
		kind:       kindStruct,
		pkgpath:    "github.com/frk/gosql/testdata/common",
		pkgname:    "common",
		pkglocal:   "common",
		isimported: true,
		ispointer:  true,
		fields: []*fieldinfo{{
			name:       "Id",
			typ:        typeinfo{kind: kindInt},
			isexported: true,
			colid:      objid{name: "id"},
			tag:        tagutil.Tag{"sql": {"id"}},
		}, {
			name:       "Email",
			typ:        typeinfo{kind: kindString},
			isexported: true,
			colid:      objid{name: "email"},
			tag:        tagutil.Tag{"sql": {"email"}},
		}, {
			name:       "FullName",
			typ:        typeinfo{kind: kindString},
			isexported: true,
			colid:      objid{name: "full_name"},
			tag:        tagutil.Tag{"sql": {"full_name"}},
		}, {
			name: "CreatedAt",
			typ: typeinfo{
				name:       "Time",
				kind:       kindStruct,
				pkgpath:    "time",
				pkgname:    "time",
				pkglocal:   "time",
				isimported: true,
				istime:     true,
			},
			isexported: true,
			colid:      objid{name: "created_at"},
			tag:        tagutil.Tag{"sql": {"created_at"}},
		}},
	}

	tests := []struct {
		name string
		want *command
		err  error
	}{{
		name: "InsertTestBAD1",
		err:  &analysisError{code: noRecordError, args: args{"InsertTestBAD1"}},
	}, {
		name: "InsertTestBAD2",
		err:  &analysisError{code: noRecordError, args: args{"InsertTestBAD2"}},
	}, {
		name: "InsertTestBAD3",
		err:  &analysisError{code: badRecordTypeError, args: args{"InsertTestBAD3"}},
	}, {
		name: "InsertTestOK1",
		want: &command{name: "InsertTestOK1", typ: cmdtypeInsert, rel: &relinfo{
			field:    "UserRec",
			relid:    objid{name: "users_table"},
			datatype: datatype{typeinfo: commonUserTypeinfo},
		}},
	}, {
		name: "InsertTestOK2",
		want: &command{name: "InsertTestOK2", typ: cmdtypeInsert, rel: &relinfo{
			field: "UserRec",
			relid: objid{name: "users_table"},
			datatype: datatype{
				typeinfo: typeinfo{
					kind: kindStruct,
					fields: []*fieldinfo{{
						name:       "Name3",
						typ:        typeinfo{kind: kindString},
						isexported: true,
						colid:      objid{name: "name"},
						tag:        tagutil.Tag{"sql": {"name"}},
					}},
				},
			},
		}},
	}, {
		name: "SelectTestOK3",
		want: &command{name: "SelectTestOK3", typ: cmdtypeSelect, rel: &relinfo{
			field: "User",
			relid: objid{name: "users_table"},
			datatype: datatype{
				typeinfo: commonUserTypeinfo,
				useiter:  true,
			},
		}},
	}, {
		name: "SelectTestOK4",
		want: &command{name: "SelectTestOK4", typ: cmdtypeSelect, rel: &relinfo{
			field: "User",
			relid: objid{name: "users_table"},
			datatype: datatype{
				typeinfo: commonUserTypeinfo,
				useiter:  true,
			},
		}},
	}, {
		name: "SelectTestOK5",
		want: &command{name: "SelectTestOK5", typ: cmdtypeSelect, rel: &relinfo{
			field: "User",
			relid: objid{name: "users_table"},
			datatype: datatype{
				typeinfo:   commonUserTypeinfo,
				useiter:    true,
				itermethod: "Fn",
			},
		}},
	}, {
		name: "SelectTestOK6",
		want: &command{name: "SelectTestOK6", typ: cmdtypeSelect, rel: &relinfo{
			field: "User",
			relid: objid{name: "users_table"},
			datatype: datatype{
				typeinfo:   commonUserTypeinfo,
				useiter:    true,
				itermethod: "Fn",
			},
		}},
	}, {
		name: "SelectTestOK7",
		want: &command{name: "SelectTestOK7", typ: cmdtypeSelect, rel: &relinfo{
			field: "Rel",
			relid: objid{name: "a_relation"},
			datatype: datatype{
				typeinfo: typeinfo{
					kind: kindStruct,
					fields: []*fieldinfo{{
						name:   "a",
						typ:    typeinfo{kind: kindInt},
						colid:  objid{name: "a"},
						tag:    tagutil.Tag{"sql": {"a", "pk", "auto"}},
						ispkey: true,
						auto:   true,
					}, {
						name:      "b",
						typ:       typeinfo{kind: kindInt},
						colid:     objid{name: "b"},
						tag:       tagutil.Tag{"sql": {"b", "nullempty"}},
						nullempty: true,
					}, {
						name:     "c",
						typ:      typeinfo{kind: kindInt},
						colid:    objid{name: "c"},
						tag:      tagutil.Tag{"sql": {"c", "ro", "json"}},
						readonly: true,
						usejson:  true,
					}, {
						name:      "d",
						typ:       typeinfo{kind: kindInt},
						colid:     objid{name: "d"},
						tag:       tagutil.Tag{"sql": {"d", "wo"}},
						writeonly: true,
					}, {
						name:   "e",
						typ:    typeinfo{kind: kindInt},
						colid:  objid{name: "e"},
						tag:    tagutil.Tag{"sql": {"e", "+"}},
						binadd: true,
					}, {
						name:        "f",
						typ:         typeinfo{kind: kindInt},
						colid:       objid{name: "f"},
						tag:         tagutil.Tag{"sql": {"f", "coalesce"}},
						usecoalesce: true,
					}, {
						name:        "g",
						typ:         typeinfo{kind: kindInt},
						colid:       objid{name: "g"},
						tag:         tagutil.Tag{"sql": {"g", "coalesce(-1)"}},
						usecoalesce: true,
						coalesceval: "-1",
					}},
				},
			},
		}},
	}, {
		name: "InsertTestOK8",
		want: &command{name: "InsertTestOK8", typ: cmdtypeInsert, rel: &relinfo{
			field: "Rel",
			relid: objid{name: "a_relation"},
			datatype: datatype{
				typeinfo: typeinfo{
					kind: kindStruct,
					fields: []*fieldinfo{{
						name:       "Foobar",
						isexported: true,
						typ: typeinfo{
							name:       "Foo",
							kind:       kindStruct,
							pkgpath:    "github.com/frk/gosql/testdata/common",
							pkgname:    "common",
							pkglocal:   "common",
							isimported: true,
							fields: []*fieldinfo{{
								name:       "Bar",
								isexported: true,
								typ: typeinfo{
									name:       "Bar",
									kind:       kindStruct,
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
											kind:       kindStruct,
											pkgpath:    "github.com/frk/gosql/testdata/common",
											pkgname:    "common",
											pkglocal:   "common",
											isimported: true,
											fields: []*fieldinfo{{
												name:       "Val",
												isexported: true,
												typ:        typeinfo{kind: kindString},
												colid:      objid{name: "foo_bar_baz_val"},
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
			rel: &relinfo{
				field:    "Rel",
				relid:    objid{name: "a_relation"},
				datatype: datatype{typeinfo: typeinfo{kind: kindStruct}},
			},
			where: &whereblock{name: "Where", items: []*whereitem{{
				node: &wherefield{
					name:  "ID",
					colid: objid{name: "id"},
					typ:   typeinfo{kind: kindInt},
					cmp:   cmpeq,
				},
			}}},
		},
	}, {
		name: "DeleteTestOK10",
		want: &command{
			name: "DeleteTestOK10",
			typ:  cmdtypeDelete,
			rel: &relinfo{
				field:    "Rel",
				relid:    objid{name: "a_relation"},
				datatype: datatype{typeinfo: typeinfo{kind: kindStruct}},
			},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &wherecolumn{colid: objid{name: "column_a"}, pred: prednotnull}},
				{op: booland, node: &wherecolumn{colid: objid{name: "column_b"}, pred: predisnull}},
				{op: boolor, node: &wherecolumn{colid: objid{name: "column_c"}, pred: prednottrue}},
				{op: booland, node: &wherecolumn{colid: objid{name: "column_d"}, pred: predistrue}},
				{op: boolor, node: &wherecolumn{colid: objid{name: "column_e"}, pred: prednotfalse}},
				{op: boolor, node: &wherecolumn{colid: objid{name: "column_f"}, pred: predisfalse}},
				{op: booland, node: &wherecolumn{colid: objid{name: "column_g"}, pred: prednotunknown}},
				{op: booland, node: &wherecolumn{colid: objid{name: "column_h"}, pred: predisunknown}},
			}},
		},
	}, {
		name: "DeleteTestOK11",
		want: &command{
			name: "DeleteTestOK11",
			typ:  cmdtypeDelete,
			rel: &relinfo{
				field:    "Rel",
				relid:    objid{name: "a_relation"},
				datatype: datatype{typeinfo: typeinfo{kind: kindStruct}},
			},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &whereblock{name: "x", items: []*whereitem{
					{node: &wherefield{
						name:  "foo",
						typ:   typeinfo{kind: kindInt},
						colid: objid{name: "column_foo"},
						cmp:   cmpeq,
					}},
					{op: booland, node: &wherecolumn{colid: objid{name: "column_a"}, pred: predisnull}},
				}}},
				{op: boolor, node: &whereblock{name: "y", items: []*whereitem{
					{node: &wherecolumn{colid: objid{name: "column_b"}, pred: prednottrue}},
					{op: boolor, node: &wherefield{
						name:  "bar",
						typ:   typeinfo{kind: kindString},
						colid: objid{name: "column_bar"},
						cmp:   cmpeq,
					}},
					{op: booland, node: &whereblock{name: "z", items: []*whereitem{
						{node: &wherefield{
							name:  "baz",
							typ:   typeinfo{kind: kindBool},
							colid: objid{name: "column_baz"},
							cmp:   cmpeq,
						}},
						{op: booland, node: &wherefield{
							name:  "quux",
							typ:   typeinfo{kind: kindString},
							colid: objid{name: "column_quux"},
							cmp:   cmpeq,
						}},
						{op: boolor, node: &wherecolumn{colid: objid{name: "column_c"}, pred: predistrue}},
					}}},
				}}},
				{op: boolor, node: &wherecolumn{colid: objid{name: "column_d"}, pred: prednotfalse}},
				{op: booland, node: &wherecolumn{colid: objid{name: "column_e"}, pred: predisfalse}},
				{op: booland, node: &wherefield{
					name:  "foo",
					typ:   typeinfo{kind: kindInt},
					colid: objid{name: "column_foo"},
					cmp:   cmpeq,
				}},
			}},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := runAnalysis(tt.name, t)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Errorf("%v - %#v", e, err)
			}
			if e := compare.Compare(got, tt.want); e != nil {
				t.Error(e)
			}
		})
	}
}
