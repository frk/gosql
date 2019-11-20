package golang

import (
	"bytes"
	"testing"

	"github.com/frk/compare"
)

func TestDeclStmt(t *testing.T) {
	tests := []struct {
		stmt DeclStmt
		want string
	}{{
		stmt: DeclStmt{Decl: ImportDecl{{Path: "foo/bar/baz"}}},
		want: "import (\n\"foo/bar/baz\"\n)",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.stmt, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestLabeledStmt(t *testing.T) {
	tests := []struct {
		stmt LabeledStmt
		want string
	}{{
		stmt: LabeledStmt{Label: Ident{"Loop"}, Stmt: ForStmt{}},
		want: "Loop:\nfor {}",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.stmt, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestExprStmt(t *testing.T) {
	tests := []struct {
		stmt ExprStmt
		want string
	}{{
		stmt: ExprStmt{X: CallExpr{Fun: Ident{"f"}}},
		want: "f()",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.stmt, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestSendStmt(t *testing.T) {
	tests := []struct {
		stmt SendStmt
		want string
	}{{
		stmt: SendStmt{Chan: Ident{"c"}, Value: Ident{"v"}},
		want: "c <- v",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.stmt, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestIncDecStmt(t *testing.T) {
	tests := []struct {
		stmt IncDecStmt
		want string
	}{{
		stmt: IncDecStmt{X: Ident{"i"}, Token: INCDEC_INC},
		want: "i++",
	}, {
		stmt: IncDecStmt{X: Ident{"j"}, Token: INCDEC_DEC},
		want: "j--",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.stmt, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestAssignStmt(t *testing.T) {
	tests := []struct {
		stmt AssignStmt
		want string
	}{{
		stmt: AssignStmt{Lhs: []Expr{Ident{"a"}}, Rhs: []Expr{String("foo")}, Token: ASSIGN},
		want: `a = "foo"`,
	}, {
		stmt: AssignStmt{Lhs: []Expr{Ident{"a"}, Ident{"b"}}, Rhs: []Expr{String("foo"), BasicLit{"123"}}, Token: ASSIGN},
		want: `a, b = "foo", 123`,
	}, {
		stmt: AssignStmt{Lhs: []Expr{Ident{"a"}, Ident{"b"}}, Rhs: []Expr{String("foo"), BasicLit{"123"}}, Token: ASSIGN_DEFINE},
		want: `a, b := "foo", 123`,
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.stmt, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestGoStmt(t *testing.T) {
	tests := []struct {
		stmt GoStmt
		want string
	}{{
		stmt: GoStmt{Call: CallExpr{Fun: Ident{"f"}}},
		want: "go f()",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.stmt, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestDeferStmt(t *testing.T) {
	tests := []struct {
		stmt DeferStmt
		want string
	}{{
		stmt: DeferStmt{Call: CallExpr{Fun: Ident{"f"}}},
		want: "defer f()",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.stmt, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestReturnStmt(t *testing.T) {
	tests := []struct {
		stmt ReturnStmt
		want string
	}{{
		stmt: ReturnStmt{},
		want: "return",
	}, {
		stmt: ReturnStmt{CallExpr{Fun: Ident{"f"}}},
		want: "return f()",
	}, {
		stmt: ReturnStmt{ExprList{Ident{"num"}, Ident{"err"}}},
		want: "return num, err",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.stmt, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestBranchStmt(t *testing.T) {
	tests := []struct {
		stmt BranchStmt
		want string
	}{{
		stmt: BranchStmt{Token: BRANCH_BREAK, Label: Ident{"Loop"}},
		want: "break Loop",
	}, {
		stmt: BranchStmt{Token: BRANCH_CONT, Label: Ident{"Loop"}},
		want: "continue Loop",
	}, {
		stmt: BranchStmt{Token: BRANCH_CONT},
		want: "continue",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.stmt, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestBlockStmt(t *testing.T) {
	tests := []struct {
		stmt BlockStmt
		want string
	}{{
		stmt: BlockStmt{},
		want: "{}",
	}, {
		stmt: BlockStmt{List: []Stmt{
			ExprStmt{CallExpr{Fun: Ident{"foo"}}},
			ExprStmt{CallExpr{Fun: Ident{"bar"}}},
		}},
		want: "{\nfoo()\nbar()\n}",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.stmt, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestIfStmt(t *testing.T) {
	tests := []struct {
		stmt IfStmt
		want string
	}{{
		stmt: IfStmt{Cond: Ident{"true"}},
		want: "if true {}",
	}, {
		stmt: IfStmt{Cond: BinaryExpr{X: Ident{"a"}, Op: BINARY_LSS, Y: Ident{"b"}}},
		want: "if a < b {}",
	}, {
		stmt: IfStmt{
			Init: AssignStmt{
				Lhs:   []Expr{Ident{"_"}, Ident{"err"}},
				Rhs:   []Expr{CallExpr{Fun: Ident{"f"}}},
				Token: ASSIGN_DEFINE,
			},
			Cond: BinaryExpr{X: Ident{"err"}, Op: BINARY_NEQ, Y: Ident{"nil"}}},
		want: "if _, err := f(); err != nil {}",
	}, {
		stmt: IfStmt{
			Cond: BinaryExpr{X: Ident{"a"}, Op: BINARY_LSS, Y: Ident{"b"}},
			Body: BlockStmt{List: []Stmt{ReturnStmt{Ident{"a"}}}},
		},
		want: "if a < b {\nreturn a\n}",
	}, {
		stmt: IfStmt{
			Cond: BinaryExpr{X: Ident{"a"}, Op: BINARY_LSS, Y: Ident{"b"}},
			Body: BlockStmt{List: []Stmt{ReturnStmt{Ident{"a"}}}},
			Else: BlockStmt{List: []Stmt{ReturnStmt{Ident{"b"}}}},
		},
		want: "if a < b {\nreturn a\n} else {\nreturn b\n}",
	}, {
		stmt: IfStmt{
			Cond: BinaryExpr{X: Ident{"a"}, Op: BINARY_LSS, Y: Ident{"b"}},
			Body: BlockStmt{List: []Stmt{ReturnStmt{Ident{"b"}}}},
			Else: IfStmt{
				Cond: BinaryExpr{X: Ident{"a"}, Op: BINARY_LSS, Y: Ident{"c"}},
				Body: BlockStmt{List: []Stmt{ReturnStmt{Ident{"c"}}}},
				Else: IfStmt{
					Cond: BinaryExpr{X: Ident{"a"}, Op: BINARY_LSS, Y: Ident{"d"}},
					Body: BlockStmt{List: []Stmt{ReturnStmt{Ident{"d"}}}},
					Else: BlockStmt{List: []Stmt{ReturnStmt{Ident{"a"}}}},
				},
			},
		},
		want: "if a < b {\nreturn b\n} else if a < c {\nreturn c\n} else if a < d {\nreturn d\n} else {\nreturn a\n}",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.stmt, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestSwitchStmt(t *testing.T) {
	tests := []struct {
		stmt SwitchStmt
		want string
	}{{
		stmt: SwitchStmt{},
		want: "switch {}",
	}, {
		stmt: SwitchStmt{Tag: Ident{"tag"}},
		want: "switch tag {}",
	}, {
		stmt: SwitchStmt{Init: AssignStmt{
			Lhs:   []Expr{Ident{"x"}},
			Rhs:   []Expr{CallExpr{Fun: Ident{"f"}}},
			Token: ASSIGN_DEFINE,
		}},
		want: "switch x := f(); {}",
	}, {
		stmt: SwitchStmt{Init: AssignStmt{
			Lhs:   []Expr{Ident{"x"}},
			Rhs:   []Expr{CallExpr{Fun: Ident{"f"}}},
			Token: ASSIGN_DEFINE,
		}, Tag: Ident{"x"}},
		want: "switch x := f(); x {}",
	}, {
		stmt: SwitchStmt{Body: []CaseClause{
			{},
		}},
		want: "switch {\ndefault:\n}",
	}, {
		stmt: SwitchStmt{Body: []CaseClause{{
			List: []Expr{BinaryExpr{X: Ident{"a"}, Op: BINARY_LSS, Y: Ident{"b"}}},
			Body: []Stmt{ReturnStmt{Ident{"a"}}},
		}, {
			List: []Expr{BinaryExpr{X: Ident{"a"}, Op: BINARY_GTR, Y: Ident{"b"}}},
			Body: []Stmt{ReturnStmt{Ident{"b"}}},
		}}},
		want: "switch {\ncase a < b:\nreturn a\ncase a > b:\nreturn b\n}",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.stmt, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestTypeSwitchStmt(t *testing.T) {
	tests := []struct {
		stmt TypeSwitchStmt
		want string
	}{{
		stmt: TypeSwitchStmt{
			Guard: TypeSwitchGuard{X: Ident{"x"}},
		},
		want: "switch x.(type) {}",
	}, {
		stmt: TypeSwitchStmt{
			Guard: TypeSwitchGuard{Name: Ident{"v"}, X: Ident{"x"}},
		},
		want: "switch v := x.(type) {}",
	}, {
		stmt: TypeSwitchStmt{
			Init:  AssignStmt{Lhs: []Expr{Ident{"x"}}, Rhs: []Expr{Ident{"y"}}, Token: ASSIGN_DEFINE},
			Guard: TypeSwitchGuard{Name: Ident{"v"}, X: Ident{"x"}},
		},
		want: "switch x := y; v := x.(type) {}",
	}, {
		stmt: TypeSwitchStmt{
			Guard: TypeSwitchGuard{Name: Ident{"v"}, X: Ident{"x"}},
			Body:  []CaseClause{{}},
		},
		want: "switch v := x.(type) {\ndefault:\n}",
	}, {
		stmt: TypeSwitchStmt{
			Guard: TypeSwitchGuard{Name: Ident{"v"}, X: Ident{"x"}},
			Body: []CaseClause{
				{List: []Expr{Ident{"nil"}}},
				{List: []Expr{Ident{"int64"}, Ident{"int32"}, Ident{"int"}}, Body: []Stmt{AssignStmt{
					Lhs:   []Expr{Ident{"_"}},
					Rhs:   []Expr{Ident{"v"}},
					Token: ASSIGN,
				}}},
				{},
			},
		},
		want: "switch v := x.(type) {\ncase nil:\ncase int64, int32, int:\n_ = v\ndefault:\n}",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.stmt, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestSelectStmt(t *testing.T) {
	tests := []struct {
		stmt SelectStmt
		want string
	}{{
		stmt: SelectStmt{},
		want: "select {}",
	}, {
		stmt: SelectStmt{Body: []CommClause{
			{},
		}},
		want: "select {\ndefault:\n}",
	}, {
		stmt: SelectStmt{Body: []CommClause{
			{Comm: SendStmt{Chan: Ident{"c"}, Value: BasicLit{"0"}}},
		}},
		want: "select {\ncase c <- 0:\n}",
	}, {
		stmt: SelectStmt{Body: []CommClause{
			{Comm: SendStmt{Chan: Ident{"c"}, Value: BasicLit{"0"}}},
			{Comm: SendStmt{Chan: Ident{"c"}, Value: BasicLit{"1"}}},
		}},
		want: "select {\ncase c <- 0:\ncase c <- 1:\n}",
	}, {
		stmt: SelectStmt{Body: []CommClause{
			{Comm: AssignStmt{
				Lhs:   []Expr{Ident{"i"}},
				Rhs:   []Expr{UnaryExpr{Op: UNARY_RECV, X: Ident{"c"}}},
				Token: ASSIGN_DEFINE,
			}, Body: []Stmt{AssignStmt{
				Lhs:   []Expr{Ident{"_"}},
				Rhs:   []Expr{Ident{"i"}},
				Token: ASSIGN,
			}}},
			{},
		}},
		want: "select {\ncase i := <-c:\n_ = i\ndefault:\n}",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.stmt, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestForStmt(t *testing.T) {
	tests := []struct {
		stmt ForStmt
		want string
	}{{
		stmt: ForStmt{},
		want: "for {}",
	}, {
		stmt: ForStmt{Cond: Ident{"false"}},
		want: "for false {}",
	}, {
		stmt: ForStmt{Cond: BinaryExpr{
			X:  Ident{"a"},
			Op: BINARY_LSS,
			Y:  Ident{"b"},
		}},
		want: "for a < b {}",
	}, {
		stmt: ForStmt{
			Init: AssignStmt{Lhs: []Expr{Ident{"i"}}, Rhs: []Expr{BasicLit{"0"}}, Token: ASSIGN_DEFINE},
			Cond: BinaryExpr{X: Ident{"a"}, Op: BINARY_LSS, Y: Ident{"b"}},
			Post: IncDecStmt{X: Ident{"i"}, Token: INCDEC_INC},
		},
		want: "for i := 0; a < b; i++ {}",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.stmt, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestRangeStmt(t *testing.T) {
	tests := []struct {
		stmt RangeStmt
		want string
	}{{
		stmt: RangeStmt{
			X: Ident{"ch"},
		},
		want: "for range ch {}",
	}, {
		stmt: RangeStmt{
			Key:    Ident{"w"},
			Define: true,
			X:      Ident{"ch"},
		},
		want: "for w := range ch {}",
	}, {
		stmt: RangeStmt{
			Key:   Ident{"key"},
			Value: Ident{"val"},
			X:     Ident{"m"},
		},
		want: "for key, val = range m {}",
	}, {
		stmt: RangeStmt{
			Key:    Ident{"i"},
			Value:  Ident{"s"},
			Define: true,
			X:      Ident{"a"},
		},
		want: "for i, s := range a {}",
	}, {
		stmt: RangeStmt{
			Key:    Ident{"i"},
			Value:  Ident{"_"},
			Define: true,
			X:      SelectorExpr{X: Ident{"testdata"}, Sel: Ident{"a"}},
		},
		want: "for i, _ := range testdata.a {}",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.stmt, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}
