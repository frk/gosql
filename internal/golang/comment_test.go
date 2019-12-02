package golang

import (
	"bytes"
	"testing"

	"github.com/frk/compare"
)

func TestLineComment(t *testing.T) {
	tests := []struct {
		list LineComment
		want string
	}{{
		list: LineComment{},
		want: "",
	}, {
		list: LineComment{" this is a comment"},
		want: "// this is a comment",
	}, {
		list: LineComment{" line 1", "", " line 2"},
		want: "// line 1\n//\n// line 2",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.list, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestGeneralComment(t *testing.T) {
	tests := []struct {
		gc   GeneralComment
		want string
	}{{
		gc:   GeneralComment{},
		want: "",
	}, {
		gc:   GeneralComment{"this is a comment"},
		want: "/*this is a comment*/",
	}, {
		gc:   GeneralComment{"line 1", "line 2"},
		want: "/*\n\tline 1\n\tline 2\n*/",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.gc, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}
