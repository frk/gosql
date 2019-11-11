package sqlang

import (
	"github.com/frk/gosql/internal/writer"
)

type UpdateStatement struct {
	Head UpdateHead
	Tail UpdateTail
}

func (s UpdateStatement) Walk(w *writer.Writer) {
	s.Head.Walk(w)
	w.NewLine()
	s.Tail.Walk(w)
	w.NoNewLine()
}

type UpdateHead struct {
	Table Ident
	Set   SetClause
	From  *FromClause
}

func (s UpdateHead) Walk(w *writer.Writer) {
	w.NoIndent()
	w.Write("UPDATE ")
	s.Table.Walk(w)
	w.Indent()
	s.Set.Walk(w)
	w.NewLine()
	s.From.Walk(w)
}

type UpdateTail struct {
	Where     *WhereClause
	Returning ReturningClause
}

func (s UpdateTail) Walk(w *writer.Writer) {
	s.Where.Walk(w)
	w.NewLine()
	s.Returning.Walk(w)
	w.NoNewLine()
}

type SetClause struct {
	Targets NameGroup
	Values  ExprGroup
}

func (c SetClause) Walk(w *writer.Writer) {
	w.Write(" SET ")
	if len(c.Targets) == 1 {
		c.Targets[0].Walk(w)
		w.Write(" = ")
		c.Values[0].Walk(w)
		return
	}

	c.Targets.Walk(w)
	w.Write(" = ")
	c.Values.Walk(w)
}

type NameGroup []Name

func (g NameGroup) Walk(w *writer.Writer) {
	w.Write("(")
	w.NewLine()
	w.Indent()

	for i, v := range g {
		if i > 0 {
			w.NewLine()
			w.Write(", ")
		}
		v.Walk(w)
	}

	w.Unindent()
	w.NewLine()
	w.Write(")")
}

type ExprGroup []Expr

func (g ExprGroup) Walk(w *writer.Writer) {
	w.Write("(")
	w.NewLine()
	w.Indent()

	for i, v := range g {
		if i > 0 {
			w.NewLine()
			w.Write(", ")
		}
		v.Walk(w)
	}

	w.Unindent()
	w.NewLine()
	w.Write(")")
}

type FromClause struct {
	List []TableExpr
}

func (c *FromClause) Walk(w *writer.Writer) {
	if c == nil {
		return
	}

	w.Write("FROM ")
	for i, x := range c.List {
		if i > 0 {
			w.NewLine()
		}
		x.Walk(w)
	}
}

// "constant table"
// table alias list

type ValuesList struct {
	Clause ValuesListClausePartial
	List   [][]Expr
	Alias  ValuesListAliasPartial
}

func (ls *ValuesList) Walk(w *writer.Writer) {
	if ls == nil {
		return
	}
	ls.Clause.Walk(w)
	for i, vs := range ls.List {
		if i > 0 {
			w.Write(",")
			w.NewLine()
		}
		w.Write("(")
		for j, x := range vs {
			if j > 0 {
				w.Write(", ")
			}
			x.Walk(w)
		}
		w.Write(")")
	}
	ls.Alias.Walk(w)
}

type ValuesListClausePartial struct{}

func (v ValuesListClausePartial) Walk(w *writer.Writer) {
	w.Write("(VALUES")
}

type ValuesListAliasPartial struct {
	Alias   string
	Columns NameGroup
}

func (v ValuesListAliasPartial) Walk(w *writer.Writer) {
	w.Write(")")

	if len(v.Alias) > 0 {
		w.Write(" AS ")
		w.Write(v.Alias)
		if len(v.Columns) > 0 {
			w.Write(" ")
			v.Columns.Walk(w)
		}
	}
}
