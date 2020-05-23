package golang

import (
	"bytes"
	"testing"

	"github.com/frk/compare"
)

func TestLineComment(t *testing.T) {
	tests := []struct {
		line LineComment
		want string
	}{{
		line: LineComment{},
		want: "",
	}, {
		line: LineComment{" this is a comment"},
		want: "// this is a comment",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.line, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestLineCommentList(t *testing.T) {
	tests := []struct {
		list LineCommentList
		want string
	}{{
		list: LineCommentList{},
		want: "",
	}, {
		list: LineCommentList{" this is a comment"},
		want: "// this is a comment",
	}, {
		list: LineCommentList{" line 1", "", " line 2"},
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

func TestBlockComment(t *testing.T) {
	tests := []struct {
		gc   BlockComment
		want string
	}{{
		gc:   BlockComment{},
		want: "",
	}, {
		gc:   BlockComment{"this is a comment"},
		want: "/*this is a comment*/",
	}, {
		gc:   BlockComment{"line 1", "line 2"},
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
