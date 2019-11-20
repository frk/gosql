package sqlang

import (
	"github.com/frk/gosql/internal/writer"
)

type TableExpr interface {
	Node
	tableExprNode()
}

type TableJoin struct {
	Type JoinType
	Rel  Ident
	On   JoinOn // can be left empty for cross joins
}

func (tj TableJoin) Walk(w *writer.Writer) {
	tj.Type.Walk(w)
	w.Write(" ")
	tj.Rel.Walk(w)

	// write the join_condition if not cross join
	if tj.Type != JoinCross {
		w.Write(" ")
		tj.On.Walk(w)
	}
}

func (Ident) tableExprNode()     {}
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
	Cmp  CmpOp
	Rhs  Expr
}

func (x JoinOnExpr) Walk(w *writer.Writer) {
	x.Bool.Walk(w)
	x.Lhs.Walk(w)
	x.Cmp.Walk(w)
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
