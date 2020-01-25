package golang

import (
	"bytes"
	"testing"

	"github.com/frk/compare"
)

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

func TestValueSpec(t *testing.T) {
	tests := []struct {
		spec ValueSpec
		want string
	}{{
		spec: ValueSpec{
			Names:  Ident{"a"},
			Values: BasicLit{"123"},
		},
		want: `a = 123`,
	}, {
		spec: ValueSpec{
			Names:  IdentList{{"a"}, {"b"}},
			Values: BasicLit{"123"},
		},
		want: `a, b = 123`,
	}, {
		spec: ValueSpec{
			Names:  IdentList{{"a"}, {"b"}},
			Values: ExprList{BasicLit{"123"}, String("123")},
		},
		want: `a, b = 123, "123"`,
	}, {
		spec: ValueSpec{
			Names: IdentList{{"a"}, {"b"}},
			Type:  SelectorExpr{X: Ident{"time"}, Sel: Ident{"Time"}},
		},
		want: `a, b time.Time`,
	}, {
		spec: ValueSpec{
			Names:  IdentList{{"a"}, {"b"}},
			Type:   SelectorExpr{X: Ident{"time"}, Sel: Ident{"Time"}},
			Values: CallExpr{Fun: SelectorExpr{X: Ident{"time"}, Sel: Ident{"Now"}}},
		},
		want: `a, b time.Time = time.Now()`,
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

func TestTypeSpec(t *testing.T) {
	tests := []struct {
		spec TypeSpec
		want string
	}{{
		spec: TypeSpec{
			Name: Ident{"ByteSlice"},
			Type: SliceType{Elt: Ident{"byte"}},
		},
		want: `ByteSlice []byte`,
	}, {
		spec: TypeSpec{
			Name:  Ident{"SomeTypeAlias"},
			Alias: true,
			Type:  Ident{"SomeType"},
		},
		want: `SomeTypeAlias = SomeType`,
	}, {
		spec: TypeSpec{
			Name: Ident{"Object"},
			Type: StructType{},
		},
		want: `Object struct{}`,
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
