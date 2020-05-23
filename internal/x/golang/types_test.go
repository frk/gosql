package golang

import (
	"bytes"
	"testing"

	"github.com/frk/compare"
)

func TestArrayType(t *testing.T) {
	tests := []struct {
		typ  ArrayType
		want string
	}{{
		typ:  ArrayType{Len: Ident{"num"}, Elem: Ident{"int"}},
		want: "[num]int",
	}, {
		typ:  ArrayType{Len: Ellipsis, Elem: Ident{"float64"}},
		want: "[...]float64",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.typ, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestSliceType(t *testing.T) {
	tests := []struct {
		typ  SliceType
		want string
	}{{
		typ:  SliceType{Elem: Ident{"int"}},
		want: "[]int",
	}, {
		typ:  SliceType{Elem: QualifiedIdent{"time", "Time"}},
		want: "[]time.Time",
	}, {
		typ: SliceType{Elem: PointerType{StructType{Fields: FieldList{
			{Names: Ident{"Foo"}, Type: Ident{"string"}},
		}}}},
		want: "[]*struct {\nFoo string\n}",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.typ, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestStructType(t *testing.T) {
	tests := []struct {
		typ  StructType
		want string
	}{{
		typ: StructType{Fields: FieldList{
			{Names: Ident{"Foo"}, Type: Ident{"string"}},
		}},
		want: "struct {\nFoo string\n}",
	}, {
		typ: StructType{Fields: FieldList{
			{Names: Ident{"Foo"}, Type: Ident{"string"}},
			{Names: Ident{"Bar"}, Type: Ident{"float64"}},
			{Names: Ident{"Baz"}, Type: Ident{"bool"}},
		}},
		want: "struct {\nFoo string\nBar float64\nBaz bool\n}",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.typ, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestFuncType(t *testing.T) {
	tests := []struct {
		typ  FuncType
		want string
	}{{
		typ: FuncType{
			Params: ParamList{},
		},
		want: "func()",
	}, {
		typ: FuncType{
			Params: ParamList{
				{Names: Ident{"foo"}, Type: Ident{"string"}},
			},
		},
		want: "func(foo string)",
	}, {
		typ: FuncType{
			Params: ParamList{
				{Names: IdentList{{"foo"}, {"bar"}, {"baz"}}, Type: Ident{"string"}},
			},
		},
		want: "func(foo, bar, baz string)",
	}, {
		typ: FuncType{
			Params: ParamList{
				{Names: Ident{"foo"}, Type: Ident{"string"}},
				{Names: Ident{"bar"}, Type: Ident{"int"}},
				{Names: Ident{"baz"}, Type: Ident{"bool"}},
			},
		},
		want: "func(foo string, bar int, baz bool)",
	}, {
		typ: FuncType{
			Params: ParamList{
				{Names: Ident{"foo"}, Type: Ident{"string"}},
				{Names: Ident{"bar"}, Type: Ident{"int"}, Variadic: true},
			},
		},
		want: "func(foo string, bar ...int)",
	}, {
		typ: FuncType{
			Params: ParamList{
				{Names: IdentList{{"foo"}, {"bar"}}, Type: Ident{"string"}},
			},
			Results: ParamList{
				{Names: nil, Type: Ident{"error"}},
			},
		},
		want: "func(foo, bar string) error",
	}, {
		typ: FuncType{
			Params: ParamList{
				{Names: IdentList{{"foo"}, {"bar"}}, Type: Ident{"string"}},
			},
			Results: ParamList{
				{Names: Ident{"err"}, Type: Ident{"error"}},
			},
		},
		want: "func(foo, bar string) (err error)",
	}, {
		typ: FuncType{
			Params: ParamList{
				{Names: IdentList{{"foo"}, {"bar"}}, Type: Ident{"string"}},
			},
			Results: ParamList{
				{Names: nil, Type: Ident{"int"}},
				{Names: nil, Type: Ident{"error"}},
			},
		},
		want: "func(foo, bar string) (int, error)",
	}, {
		typ: FuncType{
			Params: ParamList{
				{Names: IdentList{{"foo"}, {"bar"}}, Type: Ident{"string"}},
			},
			Results: ParamList{
				{Names: Ident{"num"}, Type: Ident{"int"}},
				{Names: Ident{"err"}, Type: Ident{"error"}},
			},
		},
		want: "func(foo, bar string) (num int, err error)",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.typ, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestSignature(t *testing.T) {
	tests := []struct {
		typ  Signature
		want string
	}{{
		typ: Signature{
			Params: ParamList{},
		},
		want: "()",
	}, {
		typ: Signature{
			Params: ParamList{
				{Names: Ident{"foo"}, Type: Ident{"string"}},
			},
		},
		want: "(foo string)",
	}, {
		typ: Signature{
			Params: ParamList{
				{Names: IdentList{{"foo"}, {"bar"}, {"baz"}}, Type: Ident{"string"}},
			},
		},
		want: "(foo, bar, baz string)",
	}, {
		typ: Signature{
			Params: ParamList{
				{Names: Ident{"foo"}, Type: Ident{"string"}},
				{Names: Ident{"bar"}, Type: Ident{"int"}},
				{Names: Ident{"baz"}, Type: Ident{"bool"}},
			},
		},
		want: "(foo string, bar int, baz bool)",
	}, {
		typ: Signature{
			Params: ParamList{
				{Names: Ident{"foo"}, Type: Ident{"string"}},
				{Names: Ident{"bar"}, Type: Ident{"int"}, Variadic: true},
			},
		},
		want: "(foo string, bar ...int)",
	}, {
		typ: Signature{
			Params: ParamList{
				{Names: IdentList{{"foo"}, {"bar"}}, Type: Ident{"string"}},
			},
			Results: ParamList{
				{Names: nil, Type: Ident{"error"}},
			},
		},
		want: "(foo, bar string) error",
	}, {
		typ: Signature{
			Params: ParamList{
				{Names: IdentList{{"foo"}, {"bar"}}, Type: Ident{"string"}},
			},
			Results: ParamList{
				{Names: Ident{"err"}, Type: Ident{"error"}},
			},
		},
		want: "(foo, bar string) (err error)",
	}, {
		typ: Signature{
			Params: ParamList{
				{Names: IdentList{{"foo"}, {"bar"}}, Type: Ident{"string"}},
			},
			Results: ParamList{
				{Names: nil, Type: Ident{"int"}},
				{Names: nil, Type: Ident{"error"}},
			},
		},
		want: "(foo, bar string) (int, error)",
	}, {
		typ: Signature{
			Params: ParamList{
				{Names: IdentList{{"foo"}, {"bar"}}, Type: Ident{"string"}},
			},
			Results: ParamList{
				{Names: Ident{"num"}, Type: Ident{"int"}},
				{Names: Ident{"err"}, Type: Ident{"error"}},
			},
		},
		want: "(foo, bar string) (num int, err error)",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.typ, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestInterfaceType(t *testing.T) {
	tests := []struct {
		typ  InterfaceType
		want string
	}{{
		typ:  InterfaceType{},
		want: "interface{}",
	}, {
		typ: InterfaceType{Methods: MethodList{
			{Name: Ident{"Foo"}, Type: Signature{}},
		}},
		want: "interface {\nFoo()\n}",
	}, {
		typ: InterfaceType{Methods: MethodList{
			{Name: Ident{"Foo"}, Type: Signature{}},
			{Name: Ident{"Bar"}, Type: Signature{
				Params: ParamList{
					{Names: Ident{"a"}, Type: Ident{"string"}},
				},
			}},
			{Name: Ident{"Baz"}, Type: Signature{
				Results: ParamList{
					{Names: nil, Type: Ident{"error"}},
				},
			}},
		}},
		want: "interface {\nFoo()\nBar(a string)\nBaz() error\n}",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.typ, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestMapType(t *testing.T) {
	tests := []struct {
		typ  MapType
		want string
	}{{
		typ:  MapType{Key: Ident{"int"}, Value: Ident{"string"}},
		want: "map[int]string",
	}, {
		typ:  MapType{Key: Ident{"string"}, Value: Ident{"string"}},
		want: "map[string]string",
	}, {
		typ:  MapType{Key: Ident{"string"}, Value: StructType{}},
		want: "map[string]struct{}",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.typ, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestChanType(t *testing.T) {
	tests := []struct {
		typ  ChanType
		want string
	}{{
		typ:  ChanType{Elem: Ident{"string"}},
		want: "chan string",
	}, {
		typ:  ChanType{Dir: ChanRecv, Elem: StructType{}},
		want: "<-chan struct{}",
	}, {
		typ:  ChanType{Dir: ChanSend, Elem: StructType{}},
		want: "chan<- struct{}",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.typ, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestField(t *testing.T) {
	tests := []struct {
		f    Field
		want string
	}{{
		f:    Field{Type: Ident{"int"}},
		want: "int",
	}, {
		f:    Field{Names: Ident{"a"}, Type: Ident{"int"}},
		want: "a int",
	}, {
		f: Field{
			Names: IdentList{{"foo"}, {"bar"}, {"baz"}},
			Type:  PointerType{Ident{"string"}},
		},
		want: "foo, bar, baz *string",
	}, {
		f: Field{
			Names: Ident{"Foo"},
			Type:  Ident{"string"},
			Tag:   RawStringLit(`json:"foo"`),
		},
		want: "Foo string `json:\"foo\"`",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.f, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}
