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
			column:     column{ident: ident{name: "id"}},
			tag:        tagutil.Tag{"sql": {"id"}},
		}, {
			name:       "Email",
			typ:        typeinfo{kind: kindString},
			isexported: true,
			column:     column{ident: ident{name: "email"}},
			tag:        tagutil.Tag{"sql": {"email"}},
		}, {
			name:       "FullName",
			typ:        typeinfo{kind: kindString},
			isexported: true,
			column:     column{ident: ident{name: "full_name"}},
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
			column:     column{ident: ident{name: "created_at"}},
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
			ident:    ident{name: "users_table"},
			datatype: datatype{typeinfo: commonUserTypeinfo},
		}},
	}, {
		name: "InsertTestOK2",
		want: &command{name: "InsertTestOK2", typ: cmdtypeInsert, rel: &relinfo{
			field: "UserRec",
			ident: ident{name: "users_table"},
			datatype: datatype{
				typeinfo: typeinfo{
					kind: kindStruct,
					fields: []*fieldinfo{{
						name:       "Name3",
						typ:        typeinfo{kind: kindString},
						isexported: true,
						column:     column{ident: ident{name: "name"}},
						tag:        tagutil.Tag{"sql": {"name"}},
					}},
				},
			},
		}},
	}, {
		name: "SelectTestOK3",
		want: &command{name: "SelectTestOK3", typ: cmdtypeSelect, rel: &relinfo{
			field: "User",
			ident: ident{name: "users_table"},
			datatype: datatype{
				typeinfo: commonUserTypeinfo,
				iter:     &iterator{},
			},
		}},
	}, {
		name: "SelectTestOK4",
		want: &command{name: "SelectTestOK4", typ: cmdtypeSelect, rel: &relinfo{
			field: "User",
			ident: ident{name: "users_table"},
			datatype: datatype{
				typeinfo: commonUserTypeinfo,
				iter:     &iterator{},
			},
		}},
	}, {
		name: "SelectTestOK5",
		want: &command{name: "SelectTestOK5", typ: cmdtypeSelect, rel: &relinfo{
			field: "User",
			ident: ident{name: "users_table"},
			datatype: datatype{
				typeinfo: commonUserTypeinfo,
				iter:     &iterator{method: "Fn"},
			},
		}},
	}, {
		name: "SelectTestOK6",
		want: &command{name: "SelectTestOK6", typ: cmdtypeSelect, rel: &relinfo{
			field: "User",
			ident: ident{name: "users_table"},
			datatype: datatype{
				typeinfo: commonUserTypeinfo,
				iter:     &iterator{method: "Fn"},
			},
		}},
	}, {
		name: "SelectTestOK7",
		want: &command{name: "SelectTestOK7", typ: cmdtypeSelect, rel: &relinfo{
			field: "Rel",
			ident: ident{name: "a_relation"},
			datatype: datatype{
				typeinfo: typeinfo{
					kind: kindStruct,
					fields: []*fieldinfo{{
						name:   "a",
						typ:    typeinfo{kind: kindInt},
						column: column{ident: ident{name: "a"}},
						tag:    tagutil.Tag{"sql": {"a", "pk", "auto"}},
						ispkey: true,
						auto:   true,
					}, {
						name:      "b",
						typ:       typeinfo{kind: kindInt},
						column:    column{ident: ident{name: "b"}},
						tag:       tagutil.Tag{"sql": {"b", "nullempty"}},
						nullempty: true,
					}, {
						name:     "c",
						typ:      typeinfo{kind: kindInt},
						column:   column{ident: ident{name: "c"}},
						tag:      tagutil.Tag{"sql": {"c", "ro", "json"}},
						readonly: true,
						usejson:  true,
					}, {
						name:      "d",
						typ:       typeinfo{kind: kindInt},
						column:    column{ident: ident{name: "d"}},
						tag:       tagutil.Tag{"sql": {"d", "wo"}},
						writeonly: true,
					}, {
						name:   "e",
						typ:    typeinfo{kind: kindInt},
						column: column{ident: ident{name: "e"}},
						tag:    tagutil.Tag{"sql": {"e", "+"}},
						binadd: true,
					}, {
						name:     "f",
						typ:      typeinfo{kind: kindInt},
						column:   column{ident: ident{name: "f"}},
						tag:      tagutil.Tag{"sql": {"f", "coalesce"}},
						coalesce: &coalesceinfo{},
					}, {
						name:     "g",
						typ:      typeinfo{kind: kindInt},
						column:   column{ident: ident{name: "g"}},
						tag:      tagutil.Tag{"sql": {"g", "coalesce(-1)"}},
						coalesce: &coalesceinfo{"-1"},
					}},
				},
			},
		}},
	}, {
		name: "InsertTestOK8",
		want: &command{name: "InsertTestOK8", typ: cmdtypeInsert, rel: &relinfo{
			field: "Rel",
			ident: ident{name: "a_relation"},
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
												column:     column{ident: ident{name: "foo_bar_baz_val"}},
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
