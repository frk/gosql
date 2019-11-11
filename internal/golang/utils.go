package golang

import (
	"github.com/frk/gosql/internal/writer"
)

type NL struct{}

func (NL) Walk(w *writer.Writer) {
	w.Write("\n")
}

func (NL) stmtNode() {}

type StmtList struct {
	List []Stmt
}

func (s StmtList) Walk(w *writer.Writer) {
	for _, stmt := range s.List {
		stmt.Walk(w)
		w.NewLine()
	}
	w.NoNewLine()
}

func (s *StmtList) Add(ss ...Stmt) {
	s.List = append(s.List, ss...)
}

func (StmtList) stmtNode() {}

var iferrreturn = IfStmt{
	Cond: BinaryExpr{X: Ident{"err"}, Op: BINARY_NEQ, Y: Ident{"nil"}},
	Body: BlockStmt{List: []Stmt{ReturnStmt{[]Expr{Ident{"err"}}}}},
}

type IfErrReturn struct{}

func (IfErrReturn) Walk(w *writer.Writer) {
	iferrreturn.Walk(w)
}

func (IfErrReturn) stmtNode() {}

type NoOp struct{}

func (NoOp) Walk(w *writer.Writer) {}

func (NoOp) stmtNode() {}
func (NoOp) exprNode() {}
func (NoOp) declNode() {}
func (NoOp) specNode() {}
