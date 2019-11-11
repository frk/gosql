package golang

import (
	"bytes"
	"testing"

	"github.com/frk/compare"
)

func TestIdent(t *testing.T) {
	tests := []struct {
		id   Ident
		want string
	}{{
		id:   Ident{},
		want: "",
	}, {
		id:   Ident{"Name"},
		want: "Name",
	}, {
		id:   Ident{"value"},
		want: "value",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.id, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestStarExpr(t *testing.T) {
	tests := []struct {
		star StarExpr
		want string
	}{{
		star: StarExpr{X: Ident{"a"}},
		want: "*a",
	}, {
		star: StarExpr{X: StarExpr{X: Ident{"b"}}},
		want: "**b",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.star, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestUnaryExpr(t *testing.T) {
	tests := []struct {
		unary UnaryExpr
		want  string
	}{{
		unary: UnaryExpr{Op: UNARY_ADD, X: Ident{"a"}},
		want:  "+a",
	}, {
		unary: UnaryExpr{Op: UNARY_SUB, X: UnaryExpr{Op: UNARY_SUB, X: Ident{"b"}}},
		want:  "--b",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.unary, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestSelectorExpr(t *testing.T) {
	tests := []struct {
		selector SelectorExpr
		want     string
	}{{
		selector: SelectorExpr{Sel: Ident{"b"}, X: Ident{"a"}},
		want:     "a.b",
	}, {
		selector: SelectorExpr{Sel: Ident{"c"}, X: SelectorExpr{Sel: Ident{"b"}, X: Ident{"a"}}},
		want:     "a.b.c",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.selector, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestBinaryExpr(t *testing.T) {
	tests := []struct {
		binary BinaryExpr
		want   string
	}{{
		binary: BinaryExpr{X: Ident{"a"}, Op: BINARY_ADD, Y: Ident{"b"}},
		want:   "a + b",
	}, {
		binary: BinaryExpr{X: Ident{"a"}, Op: BINARY_LEQ, Y: Ident{"b"}},
		want:   "a <= b",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.binary, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestCallExpr(t *testing.T) {
	tests := []struct {
		call CallExpr
		want string
	}{{
		call: CallExpr{Fun: Ident{"foo"}},
		want: "foo()",
	}, {
		call: CallExpr{Fun: Ident{"foo"}, Args: ArgsList{
			List: []Expr{Ident{"arg1"}, Ident{"arg2"}},
		}},
		want: "foo(arg1, arg2)",
	}, {
		call: CallExpr{Fun: Ident{"foo"}, Args: ArgsList{
			List:     []Expr{Ident{"arg1"}, Ident{"arg2"}},
			Ellipsis: true,
		}},
		want: "foo(arg1, arg2...)",
	}, {
		call: CallExpr{Fun: Ident{"foo"}, Args: ArgsList{
			List:       []Expr{Ident{"arg1"}, Ident{"arg2"}},
			Ellipsis:   true,
			OnePerLine: 1,
		}},
		want: "foo(\narg1, \narg2...,\n)",
	}, {
		call: CallExpr{Fun: Ident{"foo"}, Args: ArgsList{
			List:       []Expr{Ident{"arg1"}, Ident{"arg2"}},
			Ellipsis:   true,
			OnePerLine: 2,
		}},
		want: "foo(arg1, \narg2...,\n)",
	}, {
		call: CallExpr{
			Fun: SelectorExpr{Sel: Ident{"Sqrt"}, X: Ident{"math"}},
			Args: ArgsList{List: []Expr{
				BinaryExpr{X: Ident{"0.34"}, Op: BINARY_MUL, Y: Ident{"25"}},
				Ident{"arg2"},
			}},
		},
		want: "math.Sqrt(0.34 * 25, arg2)",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.call, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestParenExpr(t *testing.T) {
	tests := []struct {
		paren ParenExpr
		want  string
	}{{
		paren: ParenExpr{},
		want:  "()",
	}, {
		paren: ParenExpr{X: Ident{"123"}},
		want:  "(123)",
	}, {
		paren: ParenExpr{X: BinaryExpr{
			X:  ParenExpr{X: BinaryExpr{X: Ident{"0.34"}, Op: BINARY_QUO, Y: Ident{"54"}}},
			Op: BINARY_MUL,
			Y:  Ident{"25"},
		}},
		want: "((0.34 / 54) * 25)",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.paren, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestIndexExpr(t *testing.T) {
	tests := []struct {
		index IndexExpr
		want  string
	}{{
		index: IndexExpr{X: Ident{"array"}, Index: Ident{"index"}},
		want:  "array[index]",
	}, {
		index: IndexExpr{X: Ident{"slice"}, Index: Ident{"0"}},
		want:  "slice[0]",
	}, {
		index: IndexExpr{X: Ident{"map"}, Index: String("key")},
		want:  `map["key"]`,
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.index, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestKeyValueExpr(t *testing.T) {
	tests := []struct {
		kv   KeyValueExpr
		want string
	}{{
		kv:   KeyValueExpr{Key: String("key"), Value: Ident{"value"}},
		want: `"key": value`,
	}, {
		kv:   KeyValueExpr{Key: Ident{"key"}, Value: String("value")},
		want: `key: "value"`,
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.kv, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestSliceExpr(t *testing.T) {
	tests := []struct {
		slice SliceExpr
		want  string
	}{{
		slice: SliceExpr{X: Ident{"array"}, Low: Ident{"low"}, High: Ident{"high"}},
		want:  `array[low:high]`,
	}, {
		slice: SliceExpr{X: Ident{"array"}, Low: Ident{"low"}, High: Ident{"high"}, Max: Ident{"max"}},
		want:  `array[low:high:max]`,
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.slice, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestTypeAssertExpr(t *testing.T) {
	tests := []struct {
		assert TypeAssertExpr
		want   string
	}{{
		assert: TypeAssertExpr{X: Ident{"value"}},
		want:   `value.(type)`,
	}, {
		assert: TypeAssertExpr{X: Ident{"value"}, Type: Ident{"string"}},
		want:   `value.(string)`,
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.assert, w); err != nil {
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
		f:    Field{Names: []Ident{{"a"}}, Type: Ident{"int"}},
		want: "a int",
	}, {
		f: Field{
			Names: []Ident{{"foo"}, {"bar"}, {"baz"}},
			Type:  StarExpr{X: Ident{"string"}},
		},
		want: "foo, bar, baz *string",
	}, {
		f: Field{
			Names: []Ident{{"Foo"}},
			Type:  Ident{"string"},
			Tag:   RawString(`json:"foo"`),
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
