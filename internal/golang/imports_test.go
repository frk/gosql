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
			Doc:   LineComment{" comment 1"},
			Specs: []ImportSpec{{Path: "foo/bar/baz"}},
		},
		want: "// comment 1\nimport (\n\"foo/bar/baz\"\n)",
	}, {
		imp: ImportDecl{
			Doc: LineComment{" comment 1"},
			Specs: []ImportSpec{{
				Doc:     LineComment{" comment 2", " comment 3"},
				Name:    Ident{"abc"},
				Path:    "foo/bar/baz",
				Comment: LineComment{" comment 4"},
			}, {
				Path: "baz/bar/foo",
			}},
		},
		want: "// comment 1\nimport (\n" +
			"// comment 2\n// comment 3\n" +
			"abc \"foo/bar/baz\" // comment 4\n" +
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
			Doc:  LineComment{" comment line 1", " comment line 2"},
			Spec: ImportSpec{Path: "foo/bar/baz"},
		},
		want: "// comment line 1\n// comment line 2\nimport \"foo/bar/baz\"",
	}, {
		imp: SingleImportDecl{
			Doc: LineComment{" comment line 1"},
			Spec: ImportSpec{
				Doc:  LineComment{" comment line 2"},
				Path: "foo/bar/baz",
			},
		},
		want: "// comment line 1\n// comment line 2\nimport \"foo/bar/baz\"",
	}, {
		imp: SingleImportDecl{
			Spec: ImportSpec{
				Doc:     LineComment{" doc comment"},
				Name:    Ident{"abc"},
				Path:    "foo/bar/baz",
				Comment: LineComment{" line comment"},
			},
		},
		want: "// doc comment\nimport abc \"foo/bar/baz\" // line comment",
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
