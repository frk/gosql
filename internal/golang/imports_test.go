package golang

import (
	"bytes"
	"testing"

	"github.com/frk/compare"
)

func TestImportDecl(t *testing.T) {
	tests := []struct {
		imp  ImportDecl
		want string
	}{{
		imp: ImportDecl{}, want: ``,
	}, {
		imp: ImportDecl{
			Doc:   LineCommentList{" comment 1"},
			Specs: []ImportSpec{{Path: "foo/bar/baz"}},
		},
		want: "// comment 1\nimport (\n\"foo/bar/baz\"\n)",
	}, {
		imp: ImportDecl{
			Doc: LineCommentList{" comment 1"},
			Specs: []ImportSpec{{
				Doc:     LineCommentList{" comment 2", " comment 3"},
				Name:    Ident{"abc"},
				Path:    "foo/bar/baz",
				Comment: LineComment{" comment 4"},
			}, {
				Path: "baz/bar/foo",
			}},
		},
		want: "// comment 1\nimport (\n" +
			"// comment 2\n// comment 3\n" +
			"abc \"foo/bar/baz\"// comment 4\n" +
			"\"baz/bar/foo\"\n)",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.imp, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestSingleImportDecl(t *testing.T) {
	tests := []struct {
		imp  SingleImportDecl
		want string
	}{{
		imp:  SingleImportDecl{},
		want: `import ""`,
	}, {
		imp: SingleImportDecl{
			Doc:  LineCommentList{" comment line 1", " comment line 2"},
			Spec: ImportSpec{Path: "foo/bar/baz"},
		},
		want: "// comment line 1\n// comment line 2\nimport \"foo/bar/baz\"",
	}, {
		imp: SingleImportDecl{
			Doc: LineCommentList{" comment line 1"},
			Spec: ImportSpec{
				Doc:  LineCommentList{" comment line 2"},
				Path: "foo/bar/baz",
			},
		},
		want: "// comment line 1\n// comment line 2\nimport \"foo/bar/baz\"",
	}, {
		imp: SingleImportDecl{
			Spec: ImportSpec{
				Doc:     LineCommentList{" doc comment"},
				Name:    Ident{"abc"},
				Path:    "foo/bar/baz",
				Comment: LineComment{" line comment"},
			},
		},
		want: "// doc comment\nimport abc \"foo/bar/baz\"// line comment",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.imp, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestImportSpec(t *testing.T) {
	tests := []struct {
		spec ImportSpec
		want string
	}{{
		spec: ImportSpec{},
		want: `""`,
	}, {
		spec: ImportSpec{Path: "foo/bar/baz"},
		want: `"foo/bar/baz"`,
	}, {
		spec: ImportSpec{Name: Ident{"."}, Path: "foo/bar/baz"},
		want: `. "foo/bar/baz"`,
	}, {
		spec: ImportSpec{Name: Ident{"_"}, Path: "foo/bar/baz"},
		want: `_ "foo/bar/baz"`,
	}, {
		spec: ImportSpec{Name: Ident{"baz"}, Path: "foo/bar/baz/v2"},
		want: `baz "foo/bar/baz/v2"`,
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.spec, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}
