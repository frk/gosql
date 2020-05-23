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
		stmt: DeclStmt{Decl: ConstDecl{Spec: ValueSpec{
			Names: Ident{"K"}, Values: StringLit("foo"),
		}}},
		want: "const K = \"foo\"",
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
		stmt: IncDecStmt{X: Ident{"i"}, Token: IncDecIncrement},
		want: "i++",
	}, {
		stmt: IncDecStmt{X: Ident{"j"}, Token: IncDecDecrement},
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
		stmt: AssignStmt{Lhs: Ident{"a"}, Rhs: StringLit("foo"), Token: Assign},
		want: `a = "foo"`,
	}, {
		stmt: AssignStmt{Lhs: ExprList{Ident{"a"}, Ident{"b"}}, Rhs: ExprList{StringLit("foo"), IntLit(123)}, Token: Assign},
		want: `a, b = "foo", 123`,
	}, {
		stmt: AssignStmt{Lhs: ExprList{Ident{"a"}, Ident{"b"}}, Rhs: ExprList{StringLit("foo"), IntLit(123)}, Token: AssignDefine},
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
		stmt: BranchStmt{Token: BranchBreak, Label: Ident{"Loop"}},
		want: "break Loop",
	}, {
		stmt: BranchStmt{Token: BranchCont, Label: Ident{"Loop"}},
		want: "continue Loop",
	}, {
		stmt: BranchStmt{Token: BranchCont},
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
		stmt: BlockStmt{List: []StmtNode{
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
		stmt: IfStmt{Cond: BinaryExpr{X: Ident{"a"}, Op: BinaryLss, Y: Ident{"b"}}},
		want: "if a < b {}",
	}, {
		stmt: IfStmt{
			Init: AssignStmt{
				Lhs:   ExprList{Ident{"_"}, Ident{"err"}},
				Rhs:   CallExpr{Fun: Ident{"f"}},
				Token: AssignDefine,
			},
			Cond: BinaryExpr{X: Ident{"err"}, Op: BinaryNeq, Y: Ident{"nil"}}},
		want: "if _, err := f(); err != nil {}",
	}, {
		stmt: IfStmt{
			Cond: BinaryExpr{X: Ident{"a"}, Op: BinaryLss, Y: Ident{"b"}},
			Body: BlockStmt{List: []StmtNode{ReturnStmt{Ident{"a"}}}},
		},
		want: "if a < b {\nreturn a\n}",
	}, {
		stmt: IfStmt{
			Cond: BinaryExpr{X: Ident{"a"}, Op: BinaryLss, Y: Ident{"b"}},
			Body: BlockStmt{List: []StmtNode{ReturnStmt{Ident{"a"}}}},
			Else: BlockStmt{List: []StmtNode{ReturnStmt{Ident{"b"}}}},
		},
		want: "if a < b {\nreturn a\n} else {\nreturn b\n}",
	}, {
		stmt: IfStmt{
			Cond: BinaryExpr{X: Ident{"a"}, Op: BinaryLss, Y: Ident{"b"}},
			Body: BlockStmt{List: []StmtNode{ReturnStmt{Ident{"b"}}}},
			Else: IfStmt{
				Cond: BinaryExpr{X: Ident{"a"}, Op: BinaryLss, Y: Ident{"c"}},
				Body: BlockStmt{List: []StmtNode{ReturnStmt{Ident{"c"}}}},
				Else: IfStmt{
					Cond: BinaryExpr{X: Ident{"a"}, Op: BinaryLss, Y: Ident{"d"}},
					Body: BlockStmt{List: []StmtNode{ReturnStmt{Ident{"d"}}}},
					Else: BlockStmt{List: []StmtNode{ReturnStmt{Ident{"a"}}}},
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
		stmt: SwitchStmt{X: Ident{"tag"}},
		want: "switch tag {}",
	}, {
		stmt: SwitchStmt{Init: AssignStmt{
			Lhs:   Ident{"x"},
			Rhs:   CallExpr{Fun: Ident{"f"}},
			Token: AssignDefine,
		}},
		want: "switch x := f(); {}",
	}, {
		stmt: SwitchStmt{Init: AssignStmt{
			Lhs:   Ident{"x"},
			Rhs:   CallExpr{Fun: Ident{"f"}},
			Token: AssignDefine,
		}, X: Ident{"x"}},
		want: "switch x := f(); x {}",
	}, {
		stmt: SwitchStmt{Cases: []CaseClause{
			{},
		}},
		want: "switch {\ndefault:\n}",
	}, {
		stmt: SwitchStmt{Cases: []CaseClause{{
			List: BinaryExpr{X: Ident{"a"}, Op: BinaryLss, Y: Ident{"b"}},
			Body: []StmtNode{ReturnStmt{Ident{"a"}}},
		}, {
			List: BinaryExpr{X: Ident{"a"}, Op: BinaryGtr, Y: Ident{"b"}},
			Body: []StmtNode{ReturnStmt{Ident{"b"}}},
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
			Init:  AssignStmt{Lhs: Ident{"x"}, Rhs: Ident{"y"}, Token: AssignDefine},
			Guard: TypeSwitchGuard{Name: Ident{"v"}, X: Ident{"x"}},
		},
		want: "switch x := y; v := x.(type) {}",
	}, {
		stmt: TypeSwitchStmt{
			Guard: TypeSwitchGuard{Name: Ident{"v"}, X: Ident{"x"}},
			Cases: []TypeCaseClause{{}},
		},
		want: "switch v := x.(type) {\ndefault:\n}",
	}, {
		stmt: TypeSwitchStmt{
			Guard: TypeSwitchGuard{Name: Ident{"v"}, X: Ident{"x"}},
			Cases: []TypeCaseClause{
				{List: Ident{"nil"}},
				{List: TypeList{Ident{"int64"}, Ident{"int32"}, Ident{"int"}}, Body: []StmtNode{AssignStmt{
					Lhs:   Ident{"_"},
					Rhs:   Ident{"v"},
					Token: Assign,
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
		stmt: SelectStmt{Cases: []CommClause{
			{},
		}},
		want: "select {\ndefault:\n}",
	}, {
		stmt: SelectStmt{Cases: []CommClause{
			{Comm: SendStmt{Chan: Ident{"c"}, Value: IntLit(0)}},
		}},
		want: "select {\ncase c <- 0:\n}",
	}, {
		stmt: SelectStmt{Cases: []CommClause{
			{Comm: SendStmt{Chan: Ident{"c"}, Value: IntLit(0)}},
			{Comm: SendStmt{Chan: Ident{"c"}, Value: IntLit(1)}},
		}},
		want: "select {\ncase c <- 0:\ncase c <- 1:\n}",
	}, {
		stmt: SelectStmt{Cases: []CommClause{
			{Comm: AssignStmt{
				Lhs:   Ident{"i"},
				Rhs:   UnaryExpr{Op: UnaryRecv, X: Ident{"c"}},
				Token: AssignDefine,
			}, Body: []StmtNode{AssignStmt{
				Lhs:   Ident{"_"},
				Rhs:   Ident{"i"},
				Token: Assign,
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
		stmt: ForStmt{Clause: ForCondition{Ident{"false"}}},
		want: "for false {}",
	}, {
		stmt: ForStmt{Clause: ForCondition{BinaryExpr{
			X:  Ident{"a"},
			Op: BinaryLss,
			Y:  Ident{"b"},
		}}},
		want: "for a < b {}",
	}, {
		stmt: ForStmt{Clause: ForClause{
			Init: AssignStmt{Lhs: Ident{"i"}, Rhs: IntLit(0), Token: AssignDefine},
			Cond: BinaryExpr{X: Ident{"a"}, Op: BinaryLss, Y: Ident{"b"}},
			Post: IncDecStmt{X: Ident{"i"}, Token: IncDecIncrement},
		}},
		want: "for i := 0; a < b; i++ {}",
	}, {
		stmt: ForStmt{Clause: ForRangeClause{
			X: Ident{"ch"},
		}},
		want: "for range ch {}",
	}, {
		stmt: ForStmt{Clause: ForRangeClause{
			Key:    Ident{"w"},
			Define: true,
			X:      Ident{"ch"},
		}},
		want: "for w := range ch {}",
	}, {
		stmt: ForStmt{Clause: ForRangeClause{
			Key:   Ident{"key"},
			Value: Ident{"val"},
			X:     Ident{"m"},
		}},
		want: "for key, val = range m {}",
	}, {
		stmt: ForStmt{Clause: ForRangeClause{
			Value: Ident{"val"},
			X:     Ident{"m"},
		}},
		want: "for _, val = range m {}",
	}, {
		stmt: ForStmt{Clause: ForRangeClause{
			Key:    Ident{"i"},
			Value:  Ident{"s"},
			Define: true,
			X:      Ident{"a"},
		}},
		want: "for i, s := range a {}",
	}, {
		stmt: ForStmt{Clause: ForRangeClause{
			Key:    Ident{"i"},
			Value:  Ident{"_"},
			Define: true,
			X:      SelectorExpr{X: Ident{"testdata"}, Sel: Ident{"a"}},
		}},
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
