package golang

import (
	"bytes"
	"testing"

	"github.com/frk/compare"
)

func TestPointerIndirectionExpr(t *testing.T) {
	tests := []struct {
		pix  PointerIndirectionExpr
		want string
	}{{
		pix:  PointerIndirectionExpr{X: Ident{"a"}},
		want: "*a",
	}, {
		pix:  PointerIndirectionExpr{X: PointerIndirectionExpr{X: Ident{"b"}}},
		want: "**b",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.pix, w); err != nil {
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
		unary: UnaryExpr{Op: UnaryAdd, X: Ident{"a"}},
		want:  "+a",
	}, {
		unary: UnaryExpr{Op: UnarySub, X: UnaryExpr{Op: UnarySub, X: Ident{"b"}}},
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
		binary: BinaryExpr{X: Ident{"a"}, Op: BinaryAdd, Y: Ident{"b"}},
		want:   "a + b",
	}, {
		binary: BinaryExpr{X: Ident{"a"}, Op: BinaryLeq, Y: Ident{"b"}},
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
			List: ExprList{Ident{"arg1"}, Ident{"arg2"}},
		}},
		want: "foo(arg1, arg2)",
	}, {
		call: CallExpr{Fun: Ident{"foo"}, Args: ArgsList{
			List:     ExprList{Ident{"arg1"}, Ident{"arg2"}},
			Ellipsis: true,
		}},
		want: "foo(arg1, arg2...)",
	}, {
		call: CallExpr{Fun: Ident{"foo"}, Args: ArgsList{
			List:       ExprList{Ident{"arg1"}, Ident{"arg2"}},
			Ellipsis:   true,
			OnePerLine: 1,
		}},
		want: "foo(\narg1, \narg2...,\n)",
	}, {
		call: CallExpr{Fun: Ident{"foo"}, Args: ArgsList{
			List:       ExprList{Ident{"arg1"}, Ident{"arg2"}},
			Ellipsis:   true,
			OnePerLine: 2,
		}},
		want: "foo(arg1, \narg2...,\n)",
	}, {
		call: CallExpr{
			Fun: QualifiedIdent{"math", "Sqrt"},
			Args: ArgsList{List: ExprList{
				BinaryExpr{X: Ident{"0.34"}, Op: BinaryMul, Y: Ident{"25"}},
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

func TestCallNewExpr(t *testing.T) {
	tests := []struct {
		call CallNewExpr
		want string
	}{{
		call: CallNewExpr{Ident{"T"}},
		want: "new(T)",
	}, {
		call: CallNewExpr{QualifiedIdent{"time", "Time"}},
		want: "new(time.Time)",
	}, {
		call: CallNewExpr{StructType{Field{Names: Ident{"F"}, Type: Ident{"string"}}}},
		want: "new(struct {\nF string\n})",
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

func TestCallMakeExpr(t *testing.T) {
	tests := []struct {
		call CallMakeExpr
		want string
	}{{
		call: CallMakeExpr{Ident{"S"}, nil, nil},
		want: "make(S)",
	}, {
		call: CallMakeExpr{Ident{"S"}, IntLit(32), nil},
		want: "make(S, 32)",
	}, {
		call: CallMakeExpr{SliceType{Ident{"T"}}, IntLit(32), IntLit(128)},
		want: "make([]T, 32, 128)",
	}, {
		call: CallMakeExpr{MapType{Ident{"string"}, Ident{"interface{}"}}, IntLit(1000), nil},
		want: "make(map[string]interface{}, 1000)",
	}, {
		call: CallMakeExpr{ChanType{0, Ident{"struct{}"}}, IntLit(1000), nil},
		want: "make(chan struct{}, 1000)",
	}, {
		call: CallMakeExpr{ChanType{0, Ident{"struct{}"}}, nil, nil},
		want: "make(chan struct{})",
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

func TestCallLenExpr(t *testing.T) {
	tests := []struct {
		call CallLenExpr
		want string
	}{{
		call: CallLenExpr{Ident{"someSliceValue"}},
		want: "len(someSliceValue)",
	}, {
		call: CallLenExpr{QualifiedIdent{"list", "Items"}},
		want: "len(list.Items)",
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
			X:  ParenExpr{X: BinaryExpr{X: Ident{"0.34"}, Op: BinaryQuo, Y: Ident{"54"}}},
			Op: BinaryMul,
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
		index: IndexExpr{X: Ident{"slice"}, Index: IntLit(0)},
		want:  "slice[0]",
	}, {
		index: IndexExpr{X: Ident{"map"}, Index: StringLit("key")},
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
