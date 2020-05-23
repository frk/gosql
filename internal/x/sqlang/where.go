package sqlang

import (
	"github.com/frk/gosql/internal/x/writer"
)

type WhereClause struct {
	SearchCondition BoolValueExpr
}

func (where WhereClause) Walk(w *writer.Writer) {
	if where.SearchCondition != nil {
		w.Write("WHERE ")
		where.SearchCondition.Walk(w)
	}
}

type BoolValueExpr interface {
	Node
	boolValueExpr()
}

type BoolValueExprList struct {
	Parenthesized bool
	Initial       BoolValueExpr
	Items         []BoolOpExpr
	ListStyle     bool
}

func (list BoolValueExprList) Walk(w *writer.Writer) {
	if list.ListStyle {

		// list style
		if list.Parenthesized {
			w.Write("(")
			w.Indent()
			w.NewLine()
		}
		list.Initial.Walk(w)
		for _, x := range list.Items {
			w.NewLine()
			x.Walk(w)
		}
		if list.Parenthesized {
			w.Unindent()
			w.NewLine()
			w.Write(")")
		}
	} else {

		// compact style
		if list.Parenthesized {
			w.Write("(")
		}
		list.Initial.Walk(w)
		for _, x := range list.Items {
			w.Write(" ")
			x.Walk(w)
		}
		if list.Parenthesized {
			w.Write(")")
		}
	}
}

type BoolOpExpr interface {
	Node
	boolOpExpr()
}

type AND struct {
	Not     bool
	Operand BoolValueExpr
}

func (op AND) Walk(w *writer.Writer) {
	if op.Not {
		w.Write("AND NOT ")
	} else {
		w.Write("AND ")
	}
	op.Operand.Walk(w)
}

type OR struct {
	Not     bool
	Operand BoolValueExpr
}

func (op OR) Walk(w *writer.Writer) {
	if op.Not {
		w.Write("OR NOT ")
	} else {
		w.Write("OR ")
	}
	op.Operand.Walk(w)
}

func (BoolValueExprList) boolValueExpr() {}
func (ColumnReference) boolValueExpr()   {}

func (AND) boolOpExpr() {}
func (OR) boolOpExpr()  {}
