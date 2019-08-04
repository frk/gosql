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
				name:       "Time",
				kind:       kindstruct,
				pkgpath:    "time",
				pkgname:    "time",
				pkglocal:   "time",
				isimported: true,
				istime:     true,
			},
			isexported: true,
			colid:      colid{name: "created_at"},
			tag:        tagutil.Tag{"sql": {"created_at"}},
		}},
	}

	reldummyslice := &relinfo{
		field: "Rel",
		relid: relid{name: "relation_a", alias: "a"},
		datatype: datatype{typeinfo: typeinfo{
			kind:     kindstruct,
			name:     "T",
			pkgpath:  "path/to/test",
			pkgname:  "testdata",
			pkglocal: "testdata",
			isslice:  true,
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
			relid:    relid{name: "users_table"},
			datatype: datatype{typeinfo: commonUserTypeinfo},
		}},
	}, {
		name: "InsertTestOK2",
		want: &command{name: "InsertTestOK2", typ: cmdtypeInsert, rel: &relinfo{
			field: "UserRec",
			relid: relid{name: "users_table"},
			datatype: datatype{
				typeinfo: typeinfo{
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
		want: &command{name: "SelectTestOK3", typ: cmdtypeSelect, rel: &relinfo{
			field: "User",
			relid: relid{name: "users_table"},
			datatype: datatype{
				typeinfo: commonUserTypeinfo,
				useiter:  true,
			},
		}},
	}, {
		name: "SelectTestOK4",
		want: &command{name: "SelectTestOK4", typ: cmdtypeSelect, rel: &relinfo{
			field: "User",
			relid: relid{name: "users_table"},
			datatype: datatype{
				typeinfo: commonUserTypeinfo,
				useiter:  true,
			},
		}},
	}, {
		name: "SelectTestOK5",
		want: &command{name: "SelectTestOK5", typ: cmdtypeSelect, rel: &relinfo{
			field: "User",
			relid: relid{name: "users_table"},
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
			relid: relid{name: "users_table"},
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
			relid: relid{name: "relation_a"},
			datatype: datatype{
				typeinfo: typeinfo{
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
		want: &command{name: "InsertTestOK8", typ: cmdtypeInsert, rel: &relinfo{
			field: "Rel",
			relid: relid{name: "relation_a"},
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
			rel: &relinfo{
				field:    "Rel",
				relid:    relid{name: "relation_a"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
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
			rel: &relinfo{
				field:    "Rel",
				relid:    relid{name: "relation_a"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
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
			rel: &relinfo{
				field:    "Rel",
				relid:    relid{name: "relation_a"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
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
			rel: &relinfo{
				field:    "Rel",
				relid:    relid{name: "relation_a"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
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
			rel: &relinfo{
				field:    "Rel",
				relid:    relid{name: "relation_a"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
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
			rel: &relinfo{
				field:    "Rel",
				relid:    relid{name: "relation_a"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
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
			rel: &relinfo{
				field:    "Rel",
				relid:    relid{name: "relation_a"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
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
			rel: &relinfo{
				field:    "Rel",
				relid:    relid{name: "relation_a"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
			},
			where: &whereblock{name: "Where", items: []*whereitem{
				{node: &wherefield{
					name:  "a",
					typ:   typeinfo{kind: kindint, isslice: true},
					colid: colid{name: "column_a"},
					cmp:   cmpisin,
				}},
				{op: booland, node: &wherefield{
					name:  "b",
					typ:   typeinfo{kind: kindint, isarray: true, arraylen: 5},
					colid: colid{name: "column_b"},
					cmp:   cmpnotin,
				}},
				{op: booland, node: &wherefield{
					name:  "c",
					typ:   typeinfo{kind: kindint, isslice: true},
					colid: colid{name: "column_c"},
					cmp:   cmpeq,
					sop:   scalarrany,
				}},
				{op: booland, node: &wherefield{
					name:  "d",
					typ:   typeinfo{kind: kindint, isarray: true, arraylen: 10},
					colid: colid{name: "column_d"},
					cmp:   cmpgt,
					sop:   scalarrsome,
				}},
				{op: booland, node: &wherefield{
					name:  "e",
					typ:   typeinfo{kind: kindint, isslice: true},
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
			rel: &relinfo{
				field:    "Rel",
				relid:    relid{name: "relation_a"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
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
			rel: &relinfo{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
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
			rel: &relinfo{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
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
			rel: &relinfo{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
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
			rel: &relinfo{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
			},
			all: true,
		},
	}, {
		name: "DeleteTestOK_All",
		want: &command{
			name: "DeleteTestOK_All",
			typ:  cmdtypeDelete,
			rel: &relinfo{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
			},
			all: true,
		},
	}, {
		name: "DeleteTestOK_Return",
		want: &command{
			name: "DeleteTestOK_Return",
			typ:  cmdtypeDelete,
			rel: &relinfo{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
			},
			returning: &collist{all: true},
		},
	}, {
		name: "InsertTestOK_Return",
		want: &command{
			name: "InsertTestOK_Return",
			typ:  cmdtypeInsert,
			rel: &relinfo{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
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
			rel: &relinfo{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
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
			rel: &relinfo{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
			},
			defaults: &collist{all: true},
		},
	}, {
		name: "UpdateTestOK_Default",
		want: &command{
			name: "UpdateTestOK_Default",
			typ:  cmdtypeUpdate,
			rel: &relinfo{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
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
			rel: &relinfo{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
			},
			force: &collist{all: true},
		},
	}, {
		name: "UpdateTestOK_Force",
		want: &command{
			name: "UpdateTestOK_Force",
			typ:  cmdtypeUpdate,
			rel: &relinfo{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
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
			rel: &relinfo{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
			},
			erh: "eh",
		},
	}, {
		name: "InsertTestOK_ErrorHandler",
		want: &command{
			name: "InsertTestOK_ErrorHandler",
			typ:  cmdtypeInsert,
			rel: &relinfo{
				field:    "Rel",
				relid:    relid{name: "relation_a", alias: "a"},
				datatype: datatype{typeinfo: typeinfo{kind: kindstruct}},
			},
			erh: "myerrorhandler",
		},
	}, {
		name: "SelectTestOK_Count",
		want: &command{
			name: "SelectTestOK_Count",
			typ:  cmdtypeSelect,
			rel: &relinfo{
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
			rel: &relinfo{
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
			rel: &relinfo{
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
			rel: &relinfo{
				field:    "_",
				relid:    relid{name: "relation_a", alias: "a"},
				isreldir: true,
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
			name:       "FilterTestOK_TextSearchDirective",
			typ:        cmdtypeFilter,
			rel:        reldummyslice,
			textsearch: &colid{qual: "a", name: "ts_document"},
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
