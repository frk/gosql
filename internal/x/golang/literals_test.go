package golang

import (
	"bytes"
	"testing"

	"github.com/frk/compare"
)

func TestIntLit(t *testing.T) {
	tests := []struct {
		lit  IntLit
		want string
	}{
		{lit: 0, want: "0"},
		{lit: 42, want: "42"},
		{lit: -41, want: "-41"},
	}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.lit, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestStringLit(t *testing.T) {
	tests := []struct {
		lit  StringLit
		want string
	}{
		{lit: "", want: `""`},
		{lit: "hello", want: `"hello"`},
		{lit: "hello world", want: `"hello world"`},
	}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.lit, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestRawStringLit(t *testing.T) {
	tests := []struct {
		lit  RawStringLit
		want string
	}{
		{lit: "", want: "``"},
		{lit: "hello", want: "`hello`"},
		{lit: "hello world", want: "`hello world`"},
	}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.lit, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestFuncLit(t *testing.T) {
	tests := []struct {
		lit  FuncLit
		want string
	}{{
		lit: FuncLit{
			Type: FuncType{},
			Body: BlockStmt{List: []StmtNode{
				LineComment{" ..."},
			}},
		},
		want: "func() {\n// ...\n}",
	}, {
		lit: FuncLit{
			Type: FuncType{
				Params:  ParamList{{Names: Ident{"s"}, Type: Ident{"string"}}},
				Results: ParamList{{Type: Ident{"error"}}},
			},
			Body: BlockStmt{List: []StmtNode{
				ReturnStmt{Ident{"nil"}},
			}},
		},
		want: "func(s string) error {\nreturn nil\n}",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.lit, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestSliceLit(t *testing.T) {
	tests := []struct {
		lit  SliceLit
		want string
	}{{
		lit:  SliceLit{},
		want: "{}",
	}, {
		lit:  SliceLit{Type: SliceType{Ident{"string"}}},
		want: "[]string{}",
	}, {
		lit: SliceLit{Type: SliceType{
			Ident{"string"}},
			Elems:   ExprList{StringLit(""), StringLit("abc")},
			Compact: true,
		},
		want: `[]string{"", "abc"}`,
	}, {
		lit: SliceLit{
			Type:    SliceType{Ident{"string"}},
			Elems:   ExprList{StringLit(""), StringLit("abc")},
			Compact: false,
		},
		want: "[]string{\n\"\", \n\"abc\", \n}",
	}, {
		lit: SliceLit{
			Type: SliceType{SliceType{Ident{"string"}}},
			Elems: ExprList{
				SliceLit{Elems: ExprList{StringLit(""), StringLit("abc")}, Compact: true},
				SliceLit{Elems: ExprList{StringLit("foobar")}, Compact: true},
			},
			Compact: true,
		},
		want: `[][]string{{"", "abc"}, {"foobar"}}`,
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.lit, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestStructLit(t *testing.T) {
	tests := []struct {
		lit  StructLit
		want string
	}{{
		lit:  StructLit{},
		want: "{}",
	}, {
		lit:  StructLit{Type: Ident{"Foo"}},
		want: "Foo{}",
	}, {
		lit: StructLit{Type: Ident{"Foo"}, Compact: true,
			Elems: []FieldElement{{"F", StringLit("foobar")}, {"T", IntLit(123)}}},
		want: `Foo{F: "foobar", T: 123}`,
	}, {
		lit: StructLit{Type: Ident{"Foo"}, Compact: true,
			Elems: []FieldElement{{"F", RawStringLit("foobar")}}},
		want: "Foo{F: `foobar`}",
	}, {
		lit: StructLit{Type: Ident{"Foo"},
			Elems: []FieldElement{{"F", StringLit("foobar")}, {"T", IntLit(123)}}},
		want: "Foo{\nF: \"foobar\", \nT: 123, \n}",
	}, {
		lit: StructLit{Type: Ident{"Foo"}, Elems: []FieldElement{
			{"T", StructLit{Type: QualifiedIdent{"time", "Time"}}},
		}},
		want: "Foo{\nT: time.Time{}, \n}",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.lit, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestMapLit(t *testing.T) {
	tests := []struct {
		lit  MapLit
		want string
	}{{
		lit:  MapLit{},
		want: "{}",
	}, {
		lit:  MapLit{Type: MapType{Ident{"string"}, Ident{"string"}}},
		want: "map[string]string{}",
	}, {
		lit: MapLit{Type: MapType{Ident{"string"}, Ident{"string"}},
			Compact: true,
			Elems:   []KeyElement{{StringLit("foo"), StringLit("bar")}}},
		want: `map[string]string{"foo": "bar"}`,
	}, {
		lit: MapLit{Type: MapType{Ident{"string"}, Ident{"interface{}"}},
			Elems: []KeyElement{
				{StringLit("foo"), IntLit(123)},
				{StringLit("bar"), StringLit("abc")},
			}},
		want: "map[string]interface{}{\n\"foo\": 123, \n\"bar\": \"abc\", \n}",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.lit, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}
