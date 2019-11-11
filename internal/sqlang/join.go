package sqlang

import (
	"github.com/frk/gosql/internal/writer"
)

type TableExpr interface {
	Node
	tableExprNode()
}

type TableRef struct {
	Rel Ident
}

func (tr *TableRef) Walk(w *writer.Writer) {
	tr.Rel.Walk(w)
}

type TableJoin struct {
	Rel  Ident
	Join JoinType
	On   JoinOn
}

func (tj *TableJoin) Walk(w *writer.Writer) {
	tj.Join.Walk(w)
	w.Write(" ")
	tj.Rel.Walk(w)
	w.Write(" ")
	tj.On.Walk(w)
}

func (TableRef) tableExprNode()  {}
func (TableJoin) tableExprNode() {}

type JoinOn struct {
	List []JoinOnExpr
}

func (jo JoinOn) Walk(w *writer.Writer) {
	w.Write("ON ")
	for _, x := range jo.List {
		x.Walk(w)
	}
}

type JoinOnExpr struct {
	Bool BoolOp
	Lhs  Expr
	Op   CmpOp
	Rhs  Expr
}

func (x JoinOnExpr) Walk(w *writer.Writer) {
	x.Bool.Walk(w)
	x.Lhs.Walk(w)
	x.Op.Walk(w)
	x.Rhs.Walk(w)
}

type JoinType string

func (typ JoinType) Walk(w *writer.Writer) {
	w.Write(string(typ))
}

const (
	JoinNone  JoinType = ""
	JoinLeft  JoinType = "LEFT JOIN"
	JoinRight JoinType = "RIGHT JOIN"
	JoinFull  JoinType = "FULL JOIN"
	JoinCross JoinType = "CROSS JOIN"
)
