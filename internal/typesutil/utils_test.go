package typesutil

import (
	"go/types"
	"strconv"
	"testing"

	"github.com/frk/gosql/internal/testutil"
)

func TestGetDirectiveName(t *testing.T) {
	tests := []struct {
		index int
		want  string
	}{
		{index: 0, want: ""},
		{index: 1, want: ""},
		{index: 2, want: ""},
		{index: 3, want: ""},
		{index: 4, want: ""},

		{index: 5, want: "Column"},
		{index: 6, want: "Relation"},
		{index: 7, want: "RightJoin"},
	}

	name := "GetDirectiveNameTest"
	named, _ := testutil.FindNamedType(name, tdata)
	if named == nil {
		t.Errorf("%q named type not found", name)
		return
	}

	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.index), func(t *testing.T) {
			f := (named.Underlying().(*types.Struct)).Field(tt.index)
			got := GetDirectiveName(f)
			if got != tt.want {
				t.Errorf("got=%q; want=%q", got, tt.want)
			}
		})
	}
}
