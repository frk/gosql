package golang

import (
	"bytes"
	"testing"

	"github.com/frk/compare"
)

func TestGenDecl(t *testing.T) {
	tests := []struct {
		decl GenDecl
		want string
	}{{
		decl: GenDecl{Token: DECL_CONST, Specs: ValueSpec{
			Names: Ident{"K"}, Values: String("foo"),
		}},
		want: "const K = \"foo\"",
	}, {
		decl: GenDecl{Token: DECL_CONST, Specs: ValueSpec{
			Names: IdentList{{"K"}, {"L"}, {"M"}}, Type: Ident{"SomeType"}, Values: String("some_value"),
		}},
		want: "const K, L, M SomeType = \"some_value\"",
	}, {
		decl: GenDecl{Token: DECL_CONST, Specs: SpecList{
			ValueSpec{Names: IdentList{{"K"}, {"L"}, {"M"}}, Type: Ident{"SomeType"}, Values: String("some_value")},
			ValueSpec{Names: Ident{"N"}, Type: Ident{"int64"}, Values: BasicLit{"123"}},
		}},
		want: "const (\nK, L, M SomeType = \"some_value\"\nN int64 = 123\n)",
	}, {
		decl: GenDecl{Token: DECL_VAR, Specs: ValueSpec{
			Names: Ident{"V"}, Values: String("foo"),
		}},
		want: "var V = \"foo\"",
	}, {
		decl: GenDecl{Token: DECL_VAR, Specs: ValueSpec{
			Names: IdentList{{"V"}, {"W"}, {"X"}}, Type: Ident{"SomeType"}, Values: String("some_value"),
		}},
		want: "var V, W, X SomeType = \"some_value\"",
	}, {
		decl: GenDecl{Token: DECL_VAR, Specs: SpecList{
			ValueSpec{Names: IdentList{{"V"}, {"W"}, {"X"}}, Type: Ident{"SomeType"}, Values: String("some_value")},
			ValueSpec{Names: Ident{"Y"}, Type: Ident{"int64"}, Values: BasicLit{"123"}},
		}},
		want: "var (\nV, W, X SomeType = \"some_value\"\nY int64 = 123\n)",
	}, {
		decl: GenDecl{Token: DECL_TYPE, Specs: TypeSpec{
			Name: Ident{"T"}, Type: Ident{"int8"},
		}},
		want: "type T int8",
	}, {
		decl: GenDecl{Token: DECL_TYPE, Specs: TypeSpec{
			Name: Ident{"T"}, Type: SelectorExpr{X: Ident{"time"}, Sel: Ident{"Time"}},
		}},
		want: "type T time.Time",
	}, {
		decl: GenDecl{Token: DECL_TYPE, Specs: TypeSpec{
			Name: Ident{"T"}, Alias: true, Type: SelectorExpr{X: Ident{"time"}, Sel: Ident{"Time"}},
		}},
		want: "type T = time.Time",
	}, {
		decl: GenDecl{Token: DECL_TYPE, Specs: TypeSpec{
			Name: Ident{"T"}, Type: StructType{Fields: FieldList{
				{Names: Ident{"F1"}, Type: Ident{"string"}},
				{Names: Ident{"F2"}, Type: Ident{"int"}},
				{Names: Ident{"F3"}, Type: Ident{"bool"}},
			}},
		}},
		want: "type T struct {\nF1 string\nF2 int\nF3 bool\n}",
	}, {
		decl: GenDecl{Token: DECL_TYPE, Specs: SpecList{
			TypeSpec{Name: Ident{"S"}, Type: ArrayType{Elt: StructType{}}},
			TypeSpec{Name: Ident{"T"}, Type: StructType{Fields: FieldList{
				{Names: Ident{"F1"}, Type: Ident{"string"}},
				{Names: Ident{"F2"}, Type: Ident{"int"}},
				{Names: Ident{"F3"}, Type: Ident{"bool"}},
			}}},
		}},
		want: "type (\nS []struct{}\nT struct {\nF1 string\nF2 int\nF3 bool\n}\n)",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.decl, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestFuncDecl(t *testing.T) {
	tests := []struct {
		decl FuncDecl
		want string
	}{{
		decl: FuncDecl{
			Name: Ident{"Foo"},
			Type: FuncType{},
		},
		want: "func Foo() {}",
	}, {
		decl: FuncDecl{
			Name: Ident{"Foo"},
			Recv: RecvParam{Name: Ident{"t"}, Type: StarExpr{X: Ident{"Type"}}},
			Type: FuncType{},
		},
		want: "func (t *Type) Foo() {}",
	}, {
		decl: FuncDecl{
			Name: Ident{"Foo"},
			Recv: RecvParam{Name: Ident{"t"}, Type: StarExpr{X: Ident{"Type"}}},
			Type: FuncType{
				Params: ParamList{
					{Names: IdentList{{"foo"}, {"bar"}}, Type: Ident{"string"}},
					{Names: Ident{"baz"}, Type: Ident{"bool"}},
				},
			},
		},
		want: "func (t *Type) Foo(foo, bar string, baz bool) {}",
	}, {
		decl: FuncDecl{
			Name: Ident{"Foo"},
			Recv: RecvParam{Name: Ident{"t"}, Type: StarExpr{X: Ident{"Type"}}},
			Type: FuncType{
				Params: ParamList{
					{Names: IdentList{{"foo"}, {"bar"}}, Type: Ident{"string"}},
					{Names: Ident{"baz"}, Type: Ident{"bool"}},
				},
				Results: ParamList{
					{Names: Ident{"num"}, Type: Ident{"int"}},
					{Names: Ident{"err"}, Type: Ident{"error"}},
				},
			},
		},
		want: "func (t *Type) Foo(foo, bar string, baz bool) (num int, err error) {}",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.decl, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}
