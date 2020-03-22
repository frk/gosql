package golang

import (
	"bytes"
	"testing"

	"github.com/frk/compare"
)

func TestConstDecl(t *testing.T) {
	tests := []struct {
		decl ConstDecl
		want string
	}{{
		decl: ConstDecl{Spec: ValueSpec{
			Names: Ident{"K"}, Values: StringLit("foo"),
		}},
		want: "const K = \"foo\"",
	}, {
		decl: ConstDecl{Spec: ValueSpec{
			Names: IdentList{{"K"}, {"L"}, {"M"}}, Type: Ident{"SomeType"}, Values: StringLit("some_value"),
		}},
		want: "const K, L, M SomeType = \"some_value\"",
	}, {
		decl: ConstDecl{Spec: ValueSpecList{
			{Names: IdentList{{"K"}, {"L"}, {"M"}}, Type: Ident{"SomeType"}, Values: StringLit("some_value")},
			{Names: Ident{"N"}, Type: Ident{"int64"}, Values: IntLit(123)},
		}},
		want: "const (\nK, L, M SomeType = \"some_value\"\nN int64 = 123\n)",
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

func TestVarDecl(t *testing.T) {
	tests := []struct {
		decl VarDecl
		want string
	}{{
		decl: VarDecl{Spec: ValueSpec{
			Names: Ident{"V"}, Values: StringLit("foo"),
		}},
		want: "var V = \"foo\"",
	}, {
		decl: VarDecl{Spec: ValueSpec{
			Names: IdentList{{"V"}, {"W"}, {"X"}}, Type: Ident{"SomeType"}, Values: StringLit("some_value"),
		}},
		want: "var V, W, X SomeType = \"some_value\"",
	}, {
		decl: VarDecl{Spec: ValueSpecList{
			{Names: IdentList{{"V"}, {"W"}, {"X"}}, Type: Ident{"SomeType"}, Values: StringLit("some_value")},
			{Names: Ident{"Y"}, Type: Ident{"int64"}, Values: IntLit(123)},
		}},
		want: "var (\nV, W, X SomeType = \"some_value\"\nY int64 = 123\n)",
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

func TestTypeDecl(t *testing.T) {
	tests := []struct {
		decl TypeDecl
		want string
	}{{
		decl: TypeDecl{Spec: TypeSpec{
			Name: Ident{"T"}, Type: Ident{"int8"},
		}},
		want: "type T int8",
	}, {
		decl: TypeDecl{Spec: TypeSpec{
			Name: Ident{"T"}, Type: QualifiedIdent{"time", "Time"},
		}},
		want: "type T time.Time",
	}, {
		decl: TypeDecl{Spec: TypeSpec{
			Name: Ident{"T"}, Alias: true, Type: QualifiedIdent{"time", "Time"},
		}},
		want: "type T = time.Time",
	}, {
		decl: TypeDecl{Spec: TypeSpec{
			Name: Ident{"T"}, Type: StructType{Fields: FieldList{
				{Names: Ident{"F1"}, Type: Ident{"string"}},
				{Names: Ident{"F2"}, Type: Ident{"int"}},
				{Names: Ident{"F3"}, Type: Ident{"bool"}},
			}},
		}},
		want: "type T struct {\nF1 string\nF2 int\nF3 bool\n}",
	}, {
		decl: TypeDecl{Spec: TypeSpecList{
			{Name: Ident{"S"}, Type: SliceType{Elem: StructType{}}},
			{Name: Ident{"T"}, Type: StructType{Fields: FieldList{
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
			Type: Signature{},
		},
		want: "func Foo() {}",
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

func TestMethodDecl(t *testing.T) {
	tests := []struct {
		decl MethodDecl
		want string
	}{{
		decl: MethodDecl{
			Name: Ident{"Foo"},
			Recv: RecvParam{Name: Ident{"t"}, Type: PointerRecvType{"Type"}},
			Type: Signature{},
		},
		want: "func (t *Type) Foo() {}",
	}, {
		decl: MethodDecl{
			Name: Ident{"Foo"},
			Recv: RecvParam{Name: Ident{"t"}, Type: PointerRecvType{"Type"}},
			Type: Signature{
				Params: ParamList{
					{Names: IdentList{{"foo"}, {"bar"}}, Type: Ident{"string"}},
					{Names: Ident{"baz"}, Type: Ident{"bool"}},
				},
			},
		},
		want: "func (t *Type) Foo(foo, bar string, baz bool) {}",
	}, {
		decl: MethodDecl{
			Name: Ident{"Foo"},
			Recv: RecvParam{Name: Ident{"t"}, Type: PointerRecvType{"Type"}},
			Type: Signature{
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
