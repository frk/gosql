package sqlang

import (
	"github.com/frk/gosql/internal/writer"
)

type WhereClause struct {
	Conds []SearchCondition
}

func (c *WhereClause) Walk(w *writer.Writer) {
	if c == nil || len(c.Conds) == 0 {
		return
	}
	w.Write("WHERE ")
	for _, cond := range c.Conds {
		cond.Walk(w)
	}
}

type SearchCondition struct {
	Bool BoolOp
	Lhs  Expr
	Op   CmpOp
	Rhs  Expr
}

func (c SearchCondition) Walk(w *writer.Writer) {
	c.Bool.Walk(w)
	c.Lhs.Walk(w)
	c.Op.Walk(w)
	c.Rhs.Walk(w)
}

type CmpOp string

const (
	CmpNone CmpOp = ""
	CmpEq   CmpOp = "="
	CmpNe   CmpOp = "<>"
	CmpGt   CmpOp = ">"
	CmpLt   CmpOp = "<"
	CmpGe   CmpOp = ">="
	CmpLe   CmpOp = "<="
)

func (op CmpOp) Walk(w *writer.Writer) {
	if op == CmpNone {
		return
	}

	w.Write(" ")
	w.Write(string(op))
	w.Write(" ")
}

type BoolOp string

const (
	BoolNone BoolOp = ""
	BoolNot  BoolOp = "NOT"
	BoolAnd  BoolOp = "AND"
	BoolOr   BoolOp = "OR"
)

func (op BoolOp) Walk(w *writer.Writer) {
	if op == BoolNone {
		return
	}

	w.Write(" ")
	w.Write(string(op))
	w.Write(" ")
}

type Modifier func(Expr) Expr

func ModIdent(x Expr) Expr { return x }

func ModLower(x Expr) Expr { return FuncExpr{"lower", x} }

func ModUpper(x Expr) Expr { return FuncExpr{"upper", x} }

type FuncExpr struct {
	Name string
	X    Expr
}

func (fx FuncExpr) Walk(w *writer.Writer) {
	w.Write(fx.Name)
	w.Write("(")
	fx.X.Walk(w)
	w.Write(")")
}

func (FuncExpr) exprNode() {}
