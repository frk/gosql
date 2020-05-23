package golang

import (
	"bytes"
	"testing"

	"github.com/frk/compare"
)

func TestValueSpec(t *testing.T) {
	tests := []struct {
		spec ValueSpec
		want string
	}{{
		spec: ValueSpec{
			Names:  Ident{"a"},
			Values: IntLit(123),
		},
		want: `a = 123`,
	}, {
		spec: ValueSpec{
			Names:  IdentList{{"a"}, {"b"}},
			Values: IntLit(123),
		},
		want: `a, b = 123`,
	}, {
		spec: ValueSpec{
			Names:  IdentList{{"a"}, {"b"}},
			Values: ExprList{IntLit(123), StringLit("123")},
		},
		want: `a, b = 123, "123"`,
	}, {
		spec: ValueSpec{
			Names: IdentList{{"a"}, {"b"}},
			Type:  QualifiedIdent{"time", "Time"},
		},
		want: `a, b time.Time`,
	}, {
		spec: ValueSpec{
			Names:  IdentList{{"a"}, {"b"}},
			Type:   QualifiedIdent{"time", "Time"},
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
			Type: SliceType{Elem: Ident{"byte"}},
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

func TestTypeSpecList(t *testing.T) {
	tests := []struct {
		spec TypeSpecList
		want string
	}{{
		spec: TypeSpecList{},
		want: ``,
	}, {
		spec: TypeSpecList{{
			Name: Ident{"ByteSlice"},
			Type: SliceType{Elem: Ident{"byte"}},
		}},
		want: `ByteSlice []byte`,
	}, {
		spec: TypeSpecList{{
			Name: Ident{"ByteSlice"},
			Type: SliceType{Elem: Ident{"byte"}},
		}, {
			Name:  Ident{"SomeTypeAlias"},
			Alias: true,
			Type:  Ident{"SomeType"},
		}},
		want: "(\nByteSlice []byte\nSomeTypeAlias = SomeType\n)",
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
