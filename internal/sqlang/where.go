package sqlang

import (
	"github.com/frk/gosql/internal/writer"
)

type WhereClause struct {
	Conds []SearchCondition
}

func (wc *WhereClause) Walk(w *writer.Writer) {
	if wc == nil || len(wc.Conds) == 0 {
		return
	}
	w.Write("WHERE ")
	if len(wc.Conds) > 2 {
		wc.Conds[0].Walk(w)
		for _, cond := range wc.Conds[1:] {
			w.NewLine()
			cond.Walk(w)
		}
	} else {
		for i, cond := range wc.Conds {
			if i > 0 {
				w.Write(" ")
			}
			cond.Walk(w)
		}
	}
}

type SearchCondition interface {
	Node
	searchConditionNode()
}

type BoolExpr struct {
	Bool BoolOp
	Lhs  Expr
	Cmp  CmpOp
	Rhs  Expr
}

func (x BoolExpr) Walk(w *writer.Writer) {
	x.Bool.Walk(w)
	if len(x.Bool) > 0 {
		w.Write(" ")
	}

	x.Lhs.Walk(w)
	x.Cmp.Walk(w)

	if x.Rhs != nil {
		x.Rhs.Walk(w)
	}
}

type BoolExprList struct {
	Bool BoolOp
	List []SearchCondition
}

func (xl BoolExprList) Walk(w *writer.Writer) {
	xl.Bool.Walk(w)
	if len(xl.Bool) > 0 {
		w.Write(" ")
	}

	w.Write("(")
	for i, x := range xl.List {
		if i > 0 {
			w.Write(" ")
		}
		x.Walk(w)
	}
	w.Write(")")
}

func (BoolExpr) searchConditionNode()     {}
func (BoolExprList) searchConditionNode() {}

type CmpOp string

const (
	CmpNone CmpOp = ""
	CmpEq   CmpOp = "="
	CmpNe   CmpOp = "<>"
	CmpNe2  CmpOp = "!="
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

	//w.Write(" ")
	w.Write(string(op))
	//w.Write(" ")
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
