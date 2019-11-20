package golang

import (
	"github.com/frk/gosql/internal/writer"
)

type NL struct{}

func (NL) Walk(w *writer.Writer) {
	w.Write("\n")
}

func (NL) stmtNode() {}

type StmtList []Stmt

func (list StmtList) Walk(w *writer.Writer) {
	for _, stmt := range list {
		stmt.Walk(w)
		w.NewLine()
	}
	w.NoNewLine()
}

func (list *StmtList) Add(ss ...Stmt) {
	*list = append(*list, ss...)
}

func (StmtList) stmtNode() {}

var iferrreturn = IfStmt{
	Cond: BinaryExpr{X: Ident{"err"}, Op: BINARY_NEQ, Y: Ident{"nil"}},
	Body: BlockStmt{List: []Stmt{ReturnStmt{Ident{"err"}}}},
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
