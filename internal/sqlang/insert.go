package sqlang

import (
	"github.com/frk/gosql/internal/writer"
)

type InsertStatement struct {
	Header InsertHeader
	Tail   InsertTail
}

func (s InsertStatement) Walk(w *writer.Writer) {
	w.NoIndent()
	s.Header.Walk(w)
	w.NewLine()
	s.Tail.Walk(w)
	w.NoNewLine()
}

func (s *InsertStatement) SetSourceValues(xs []Expr) {
	s.Header.Source.Values = &ValuesClause{xs}
}

type InsertHeader struct {
	Table      Ident
	Columns    NameGroup
	Overriding OverridingClause
	Source     InsertSource
}

func (h InsertHeader) Walk(w *writer.Writer) {
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
	t.OnConflict.Walk(w)
	w.NewLine()
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
	Exprs []Expr
}

func (vc ValuesClause) Walk(w *writer.Writer) {
	w.Write("VALUES ")
	if len(vc.Exprs) == 0 {
		return
	}

	w.Write("(")
	w.NewLine()
	w.Indent()

	for i, x := range vc.Exprs {
		if i > 0 {
			w.NewLine()
			w.Write(", ")
		}
		x.Walk(w)
	}

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

	w.Write("ON CONFLICT")
	if oc.Target != nil {
		oc.Target.Walk(w)
	}
	oc.Action.Walk(w)
}

type ConflictTarget interface {
	Node
	conflictTargetNode()
}

type ConflictColumns []Name

func (cc ConflictColumns) Walk(w *writer.Writer) {
	w.Write(" ")
	for i, c := range cc {
		if i > 0 {
			w.Write(", ")
		}
		c.Walk(w)
	}
}

type ConflictIndexes []Name

func (ci ConflictIndexes) Walk(w *writer.Writer) {
	w.Write(" (")
	for i, idx := range ci {
		if i > 0 {
			w.Write(", ")
		}
		idx.Walk(w)
	}
	w.Write(")")
}

type ConflictConstraint string

type ConflictAction struct {
	Update UpdateExcluded
}

func (a *ConflictAction) Walk(w *writer.Writer) {
	if a == nil {
		w.Write(" DO NOTHING")
		return
	}
	w.Write(" DO UPDATE SET ")
	a.Update.Walk(w)
}

type UpdateExcluded []Name

func (ex UpdateExcluded) Walk(w *writer.Writer) {
	for i, c := range ex {
		if i > 0 {
			w.Write(", ")
		}
		c.Walk(w)
		w.Write(" = EXCLUDED.")
		c.Walk(w)
	}
}

type ReturningClause []ColumnExpr

func (rc ReturningClause) Walk(w *writer.Writer) {
	if len(rc) < 1 {
		return
	}

	w.Write("RETURNING ")
	for i, c := range rc {
		if i > 0 {
			w.Write(", ")
		}
		c.Walk(w)
	}
}

func (ConflictColumns) conflictTargetNode()    {}
func (ConflictIndexes) conflictTargetNode()    {}
func (ConflictConstraint) conflictTargetNode() {}
