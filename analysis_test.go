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
		kind:       kindstruct,
		pkgpath:    "github.com/frk/gosql/testdata/common",
		pkgname:    "common",
		pkglocal:   "common",
		isimported: true,
		ispointer:  true,
		fields: []*fieldinfo{{
			name:       "Id",
			typ:        typeinfo{kind: kindint},
			isexported: true,
			colid:      objid{name: "id"},
			tag:        tagutil.Tag{"sql": {"id"}},
		}, {
			name:       "Email",
			typ:        typeinfo{kind: kindstring},
			isexported: true,
			colid:      objid{name: "email"},
			tag:        tagutil.Tag{"sql": {"email"}},
		}, {
			name:       "FullName",
			typ:        typeinfo{kind: kindstring},
			isexported: true,
			colid:      objid{name: "full_name"},
			tag:        tagutil.Tag{"sql": {"full_name"}},
		}, {
			name: "CreatedAt",
			typ: typeinfo{
				name:       "Time",
				kind:       kindstruct,
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
		err:  &analysisError{code: errNoRelation, args: []interface{}{"InsertTestBAD1"}},
	}, {
		name: "InsertTestBAD2",
		err:  &analysisError{code: errNoRelation, args: []interface{}{"InsertTestBAD2"}},
	}, {
		name: "InsertTestBAD3",
		err:  &analysisError{code: errBadRelationType, args: []interface{}{"InsertTestBAD3", "User"}},
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
					kind: kindstruct,
					fields: []*fieldinfo{{
						name:       "Name3",
						typ:        typeinfo{kind: kindstring},
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
					kind: kindstruct,
					fields: []*fieldinfo{{
						name:   "a",
						typ:    typeinfo{kind: kindint},
						colid:  objid{name: "a"},
						tag:    tagutil.Tag{"sql": {"a", "pk", "auto"}},
						ispkey: true,
						auto:   true,
					}, {
						name:      "b",
						typ:       typeinfo{kind: kindint},
						colid:     objid{name: "b"},
						tag:       tagutil.Tag{"sql": {"b", "nullempty"}},
						nullempty: true,
					}, {
						name:     "c",
						typ:      typeinfo{kind: kindint},
						colid:    objid{name: "c"},
						tag:      tagutil.Tag{"sql": {"c", "ro", "json"}},
						readonly: true,
						usejson:  true,
					}, {
						name:      "d",
						typ:       typeinfo{kind: kindint},
						colid:     objid{name: "d"},
						tag:       tagutil.Tag{"sql": {"d", "wo"}},
						writeonly: true,
					}, {
						name:   "e",
						typ:    typeinfo{kind: kindint},
						colid:  objid{name: "e"},
						tag:    tagutil.Tag{"sql": {"e", "+"}},
						binadd: true,
					}, {
						name:        "f",
						typ:         typeinfo{kind: kindint},
						colid:       objid{name: "f"},
						tag:         tagutil.Tag{"sql": {"f", "coalesce"}},
						usecoalesce: true,
					}, {
						name:        "g",
						typ:         typeinfo{kind: kindint},
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
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
			},
			where: &whereblock{name: "Where", items: []*whereitem{{
				node: &wherefield{
					name:  "ID",
					colid: objid{name: "id"},
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
			rel: &relinfo{
				field:    "Rel",
				relid:    objid{name: "a_relation"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
			},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &wherecolumn{colid: objid{name: "column_a"}, cmp: cmpnotnull}},
				{op: booland, node: &wherecolumn{colid: objid{name: "column_b"}, cmp: cmpisnull}},
				{op: boolor, node: &wherecolumn{colid: objid{name: "column_c"}, cmp: cmpnottrue}},
				{op: booland, node: &wherecolumn{colid: objid{name: "column_d"}, cmp: cmpistrue}},
				{op: boolor, node: &wherecolumn{colid: objid{name: "column_e"}, cmp: cmpnotfalse}},
				{op: boolor, node: &wherecolumn{colid: objid{name: "column_f"}, cmp: cmpisfalse}},
				{op: booland, node: &wherecolumn{colid: objid{name: "column_g"}, cmp: cmpnotunknown}},
				{op: booland, node: &wherecolumn{colid: objid{name: "column_h"}, cmp: cmpisunknown}},
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
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
			},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &whereblock{name: "x", items: []*whereitem{
					{node: &wherefield{
						name:  "foo",
						typ:   typeinfo{kind: kindint},
						colid: objid{name: "column_foo"},
						cmp:   cmpeq,
					}},
					{op: booland, node: &wherecolumn{colid: objid{name: "column_a"}, cmp: cmpisnull}},
				}}},
				{op: boolor, node: &whereblock{name: "y", items: []*whereitem{
					{node: &wherecolumn{colid: objid{name: "column_b"}, cmp: cmpnottrue}},
					{op: boolor, node: &wherefield{
						name:  "bar",
						typ:   typeinfo{kind: kindstring},
						colid: objid{name: "column_bar"},
						cmp:   cmpeq,
					}},
					{op: booland, node: &whereblock{name: "z", items: []*whereitem{
						{node: &wherefield{
							name:  "baz",
							typ:   typeinfo{kind: kindbool},
							colid: objid{name: "column_baz"},
							cmp:   cmpeq,
						}},
						{op: booland, node: &wherefield{
							name:  "quux",
							typ:   typeinfo{kind: kindstring},
							colid: objid{name: "column_quux"},
							cmp:   cmpeq,
						}},
						{op: boolor, node: &wherecolumn{colid: objid{name: "column_c"}, cmp: cmpistrue}},
					}}},
				}}},
				{op: boolor, node: &wherecolumn{colid: objid{name: "column_d"}, cmp: cmpnotfalse}},
				{op: booland, node: &wherecolumn{colid: objid{name: "column_e"}, cmp: cmpisfalse}},
				{op: booland, node: &wherefield{
					name:  "foo",
					typ:   typeinfo{kind: kindint},
					colid: objid{name: "column_foo"},
					cmp:   cmpeq,
				}},
			}},
		},
	}, {
		name: "DeleteTestOK12",
		want: &command{
			name: "DeleteTestOK12",
			typ:  cmdtypeDelete,
			rel: &relinfo{
				field:    "Rel",
				relid:    objid{name: "a_relation"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
			},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &wherefield{name: "a", typ: typeinfo{kind: kindint}, colid: objid{name: "column_a"}, cmp: cmplt}},
				{op: booland, node: &wherefield{name: "b", typ: typeinfo{kind: kindint}, colid: objid{name: "column_b"}, cmp: cmpgt}},
				{op: booland, node: &wherefield{name: "c", typ: typeinfo{kind: kindint}, colid: objid{name: "column_c"}, cmp: cmple}},
				{op: booland, node: &wherefield{name: "d", typ: typeinfo{kind: kindint}, colid: objid{name: "column_d"}, cmp: cmpge}},
				{op: booland, node: &wherefield{name: "e", typ: typeinfo{kind: kindint}, colid: objid{name: "column_e"}, cmp: cmpeq}},
				{op: booland, node: &wherefield{name: "f", typ: typeinfo{kind: kindint}, colid: objid{name: "column_f"}, cmp: cmpne}},
				{op: booland, node: &wherefield{name: "g", typ: typeinfo{kind: kindint}, colid: objid{name: "column_g"}, cmp: cmpeq}},
			}},
		},
	}, {
		name: "DeleteTestOK13",
		want: &command{
			name: "DeleteTestOK13",
			typ:  cmdtypeDelete,
			rel: &relinfo{
				field:    "Rel",
				relid:    objid{name: "a_relation"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
			},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &wherecolumn{colid: objid{name: "column_a"}, cmp: cmpne, colid2: objid{name: "column_b"}}},
				{op: booland, node: &wherecolumn{colid: objid{qual: "t", name: "column_c"}, cmp: cmpeq, colid2: objid{qual: "u", name: "column_d"}}},
				{op: booland, node: &wherecolumn{colid: objid{qual: "t", name: "column_e"}, cmp: cmpgt, lit: "123"}},
				{op: booland, node: &wherecolumn{colid: objid{qual: "t", name: "column_f"}, cmp: cmpeq, lit: "'active'"}},
				{op: booland, node: &wherecolumn{colid: objid{qual: "t", name: "column_g"}, cmp: cmpne, lit: "true"}},
			}},
		},
	}, {
		name: "DeleteTestOK14",
		want: &command{
			name: "DeleteTestOK14",
			typ:  cmdtypeDelete,
			rel: &relinfo{
				field:    "Rel",
				relid:    objid{name: "a_relation"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
			},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &wherebetween{
					name:  "a",
					colid: objid{name: "column_a"},
					cmp:   cmpisbetween,
					x:     &varinfo{name: "x", typ: typeinfo{kind: kindint}},
					y:     &varinfo{name: "y", typ: typeinfo{kind: kindint}},
				}},
				{op: booland, node: &wherebetween{
					name:  "b",
					colid: objid{name: "column_b"},
					cmp:   cmpisbetweensym,
					x:     objid{name: "column_x"},
					y:     objid{name: "column_y"},
				}},
				{op: booland, node: &wherebetween{
					name:  "c",
					colid: objid{name: "column_c"},
					cmp:   cmpnotbetweensym,
					x:     objid{name: "column_z"},
					y:     &varinfo{name: "z", typ: typeinfo{kind: kindint}},
				}},
				{op: booland, node: &wherebetween{
					name:  "d",
					colid: objid{name: "column_d"},
					cmp:   cmpnotbetween,
					x:     &varinfo{name: "z", typ: typeinfo{kind: kindint}},
					y:     objid{name: "column_z"},
				}},
			}},
		},
	}, {
		name: "DeleteTestOK_DistinctFrom",
		want: &command{
			name: "DeleteTestOK_DistinctFrom",
			typ:  cmdtypeDelete,
			rel: &relinfo{
				field:    "Rel",
				relid:    objid{name: "a_relation"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
			},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &wherefield{
					name:  "a",
					typ:   typeinfo{kind: kindint},
					colid: objid{name: "column_a"},
					cmp:   cmpisdistinct,
				}},
				{op: booland, node: &wherefield{
					name:  "b",
					typ:   typeinfo{kind: kindint},
					colid: objid{name: "column_b"},
					cmp:   cmpnotdistinct,
				}},
				{op: booland, node: &wherecolumn{colid: objid{name: "column_c"}, cmp: cmpisdistinct, colid2: objid{name: "column_x"}}},
				{op: booland, node: &wherecolumn{colid: objid{name: "column_d"}, cmp: cmpnotdistinct, colid2: objid{name: "column_y"}}},
			}},
		},
	}, {
		name: "DeleteTestOK_ArrayComparisons",
		want: &command{
			name: "DeleteTestOK_ArrayComparisons",
			typ:  cmdtypeDelete,
			rel: &relinfo{
				field:    "Rel",
				relid:    objid{name: "a_relation"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
			},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &wherefield{
					name:  "a",
					typ:   typeinfo{kind: kindint, isslice: true},
					colid: objid{name: "column_a"},
					cmp:   cmpisin,
				}},
				{op: booland, node: &wherefield{
					name:  "b",
					typ:   typeinfo{kind: kindint, isarray: true, arraylen: 5},
					colid: objid{name: "column_b"},
					cmp:   cmpnotin,
				}},
				{op: booland, node: &wherefield{
					name:  "c",
					typ:   typeinfo{kind: kindint, isslice: true},
					colid: objid{name: "column_c"},
					cmp:   cmpeq,
					saop:  scalarrany,
				}},
				{op: booland, node: &wherefield{
					name:  "d",
					typ:   typeinfo{kind: kindint, isarray: true, arraylen: 10},
					colid: objid{name: "column_d"},
					cmp:   cmpgt,
					saop:  scalarrsome,
				}},
				{op: booland, node: &wherefield{
					name:  "e",
					typ:   typeinfo{kind: kindint, isslice: true},
					colid: objid{name: "column_e"},
					cmp:   cmple,
					saop:  scalarrall,
				}},
			}},
		},
	}, {
		name: "DeleteTestOK_PatternMatching",
		want: &command{
			name: "DeleteTestOK_PatternMatching",
			typ:  cmdtypeDelete,
			rel: &relinfo{
				field:    "Rel",
				relid:    objid{name: "a_relation"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
			},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &wherefield{
					name:  "a",
					typ:   typeinfo{kind: kindstring},
					colid: objid{name: "column_a"},
					cmp:   cmpislike,
				}},
				{op: booland, node: &wherefield{
					name:  "b",
					typ:   typeinfo{kind: kindstring},
					colid: objid{name: "column_b"},
					cmp:   cmpnotlike,
				}},
				{op: booland, node: &wherefield{
					name:  "c",
					typ:   typeinfo{kind: kindstring},
					colid: objid{name: "column_c"},
					cmp:   cmpissimilar,
				}},
				{op: booland, node: &wherefield{
					name:  "d",
					typ:   typeinfo{kind: kindstring},
					colid: objid{name: "column_d"},
					cmp:   cmpnotsimilar,
				}},
				{op: booland, node: &wherefield{
					name:  "e",
					typ:   typeinfo{kind: kindstring},
					colid: objid{name: "column_e"},
					cmp:   cmprexp,
				}},
				{op: booland, node: &wherefield{
					name:  "f",
					typ:   typeinfo{kind: kindstring},
					colid: objid{name: "column_f"},
					cmp:   cmprexpi,
				}},
				{op: booland, node: &wherefield{
					name:  "g",
					typ:   typeinfo{kind: kindstring},
					colid: objid{name: "column_g"},
					cmp:   cmpnotrexp,
				}},
				{op: booland, node: &wherefield{
					name:  "h",
					typ:   typeinfo{kind: kindstring},
					colid: objid{name: "column_h"},
					cmp:   cmpnotrexpi,
				}},
			}},
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
