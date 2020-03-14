package sqlang

import (
	"github.com/frk/gosql/internal/writer"
)

type InsertStatement struct {
	Head InsertHead
	Tail InsertTail
}

func (s InsertStatement) Walk(w *writer.Writer) {
	w.NoIndent()
	s.Head.Walk(w)
	if s.Tail.OnConflict != nil || len(s.Tail.Returning) > 0 {
		w.NewLine()
		s.Tail.Walk(w)
	}
	w.NoNewLine()
}

type InsertHead struct {
	Table      Ident
	Columns    NameGroup
	Overriding OverridingClause
	Source     InsertSource
}

func (h InsertHead) Walk(w *writer.Writer) {
	w.Write("INSERT INTO ")
	h.Table.Walk(w)
	w.Write(" ")

	w.Indent()
	h.Columns.Walk(w)
	w.Write(" ")

	if len(h.Overriding) > 0 {
		h.Overriding.Walk(w)
		w.Write(" ")
	}

	h.Source.Walk(w)
}

type InsertTail struct {
	OnConflict *OnConflictClause
	Returning  ReturningClause
}

func (t InsertTail) Walk(w *writer.Writer) {
	if t.OnConflict != nil {
		t.OnConflict.Walk(w)
		if len(t.Returning) > 0 {
			w.NewLine()
		}
	}
	t.Returning.Walk(w)
}

type OverridingClause string

func (c OverridingClause) Walk(w *writer.Writer) {
	if len(c) < 1 {
		return
	}
	w.Write("OVERRIDING ")
	w.Write(string(c))
	w.Write(" VALUE")
}

type InsertSource struct {
	// InsertSource expects either that both of its fields are nil or that
	// one of its field is set. If both fields are set it will produce invalid
	// sql. In the case in which both of its fields are nil it produces the
	// `DEFAULT VALUES` clause and in the case where one of its field is set
	// it will produce the clause corresponding to the field.
	Values *ValuesClause
	Select *SelectStatement
}

func (s InsertSource) Walk(w *writer.Writer) {
	if s.Values == nil && s.Select == nil {
		w.Write("DEFAULT VALUES")
		return
	}
	if s.Values != nil {
		s.Values.Walk(w)
	} else if s.Select != nil {
		s.Select.Walk(w)
	}
}

type ValuesClause struct {
	Exprs ValueExprList
}

func (vc ValuesClause) Walk(w *writer.Writer) {
	w.Write("VALUES ")
	if len(vc.Exprs) == 0 {
		return
	}

	w.Write("(")
	w.NewLine()
	w.Indent()
	vc.Exprs.Walk(w)
	w.NewLine()
	w.Unindent()
	w.Write(")")
}

type ExprList struct {
	List []Expr
}

func (el ExprList) Walk(w *writer.Writer) {
	w.Write("(")
	for i, x := range el.List {
		if i > 0 {
			w.Write(", ")
		}
		x.Walk(w)
	}
	w.Write(")")
}

type OnConflictClause struct {
	Target ConflictTarget
	Action *ConflictAction
}

func (oc *OnConflictClause) Walk(w *writer.Writer) {
	if oc == nil {
		return
	}

	w.Write("ON CONFLICT ")
	if oc.Target != nil {
		oc.Target.Walk(w)
		w.NewLine()
	}
	oc.Action.Walk(w)
}

type ConflictTarget interface {
	Node
	conflictTargetNode()
}

type ConflictColumns []Name

func (cc ConflictColumns) Walk(w *writer.Writer) {
	w.Write("(")
	for i, c := range cc {
		if i > 0 {
			w.Write(", ")
		}
		c.Walk(w)
	}
	w.Write(")")
}

type ConflictIndex struct {
	Expr string
	Pred string
}

func (ind ConflictIndex) Walk(w *writer.Writer) {
	w.Write("(" + ind.Expr + ")")
	if len(ind.Pred) > 0 {
		w.Write(" WHERE " + ind.Pred)
	}
}

type ConflictConstraint string

func (cc ConflictConstraint) Walk(w *writer.Writer) {
	w.Write("ON CONSTRAINT ")
	w.Write(`"` + string(cc) + `"`)
}

type ConflictAction struct {
	Update UpdateExcluded
}

func (a *ConflictAction) Walk(w *writer.Writer) {
	if a == nil {
		w.Write("DO NOTHING")
		return
	}
	w.Write("DO UPDATE SET")
	a.Update.Walk(w)
}

type UpdateExcluded struct {
	Columns []Name
	Compact bool
}

func (x UpdateExcluded) Walk(w *writer.Writer) {
	if len(x.Columns) > 0 {
		if !x.Compact {
			w.NewLine()
		} else {
			w.Write(" ")
		}
		x.Columns[0].Walk(w)
		w.Write(" = EXCLUDED.")
		x.Columns[0].Walk(w)
	}

	for _, c := range x.Columns[1:] {
		if !x.Compact {
			w.NewLine()
		} else {
			w.Write(" ")
		}
		w.Write(", ")
		c.Walk(w)
		w.Write(" = EXCLUDED.")
		c.Walk(w)
	}
}

type ReturningClause []ValueExpr

func (rc ReturningClause) Walk(w *writer.Writer) {
	if len(rc) < 1 {
		return
	}

	w.Write("RETURNING")

	//
	if len(rc) < 3 {
		w.Write(" ")
		for i, c := range rc {
			if i > 0 {
				w.Write(", ")
			}
			c.Walk(w)
		}
		return // exit
	}

	//
	w.NewLine()
	for i, c := range rc {
		if i > 0 {
			w.NewLine()
			w.Write(", ")
		}
		c.Walk(w)
	}
}

func (ConflictColumns) conflictTargetNode()    {}
func (ConflictIndex) conflictTargetNode()      {}
func (ConflictConstraint) conflictTargetNode() {}
