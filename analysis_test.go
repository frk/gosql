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
			field: "UserRec",
			ident: ident{name: "users_table"},
			datatype: datatype{
				typeinfo: typeinfo{
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
						tag:        tagutil.Tag{"sql": []string{"id"}},
					}, {
						name:       "Email",
						typ:        typeinfo{kind: kindString},
						isexported: true,
						tag:        tagutil.Tag{"sql": []string{"email"}},
					}, {
						name:       "FullName",
						typ:        typeinfo{kind: kindString},
						isexported: true,
						tag:        tagutil.Tag{"sql": []string{"full_name"}},
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
						tag:        tagutil.Tag{"sql": []string{"created_at"}},
					}},
				},
			},
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
						tag:        tagutil.Tag{"sql": []string{"name"}},
					}},
				},
			},
		}},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := runAnalysis(tt.name, t)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Error(e)
			}
			if e := compare.Compare(got, tt.want); e != nil {
				t.Error(e)
			}
		})
	}
}
