package gosql

import (
	"testing"

	"github.com/frk/compare"
	"github.com/frk/gosql/internal/testutil"
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
		want: &command{name: "InsertTestOK1", typ: cmdtypeInsert, rec: &record{
			field: "UserRec",
			typ: gotype{
				name:       "User",
				kind:       gokindStruct,
				pkgpath:    "github.com/frk/gosql/testdata/common",
				pkgname:    "common",
				pkglocal:   "common",
				isimported: true,
				ispointer:  true,
			},
			rel: relation{ident: ident{name: "users_table"}},
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
