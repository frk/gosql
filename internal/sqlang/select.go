package sqlang

import (
	"strconv"

	"github.com/frk/gosql/internal/writer"
)

type SelectStatement struct {
	Columns ColumnExprSlice
	Table   Ident
	Join    JoinClause
	Where   *WhereClause
	Order   OrderClause
	Limit   LimitClause
	Offset  OffsetClause

	Close  bool
	Exists bool
	Count  bool
}

func (s SelectStatement) Walk(w *writer.Writer) {
	w.NoIndent()
	if s.Exists {
		w.Write("SELECT EXISTS(SELECT 1 FROM ")
		s.Table.Walk(w)
		w.Indent()
		s.Join.Walk(w)
		w.NewLine()
	} else if s.Count {
		w.Write("SELECT COUNT(1) FROM ")
		s.Table.Walk(w)
		w.Indent()
		s.Join.Walk(w)
		w.NewLine()
	} else {
		w.Write("SELECT")
		w.NewLine()
		w.Indent()
		s.Columns.Walk(w)
		w.NewLine()
		w.Write("FROM ")
		s.Table.Walk(w)
		s.Join.Walk(w)
		w.NewLine()
	}

	s.Where.Walk(w)
	s.Order.Walk(w)
	s.Limit.Walk(w)
	s.Offset.Walk(w)

	if s.Exists && s.Close {
		w.Write(")")
	}
}

type LimitClause struct {
	Count LimitCount
}

func (l LimitClause) Walk(w *writer.Writer) {
	if l.Count == nil {
		return
	}
	w.NewLine()
	w.Write("LIMIT ")
	l.Count.Walk(w)
}

type LimitCount interface {
	Expr
	limitCountNode()
}

type LimitInt int64

func (i LimitInt) Walk(w *writer.Writer) {
	w.Write(strconv.FormatInt(int64(i), 10))
}

type LimitUint uint64

func (u LimitUint) Walk(w *writer.Writer) {
	w.Write(strconv.FormatUint(uint64(u), 10))
}

type OffsetClause struct {
	Start OffsetStart
}

func (l OffsetClause) Walk(w *writer.Writer) {
	if l.Start == nil {
		return
	}
	w.NewLine()
	w.Write("OFFSET ")
	l.Start.Walk(w)
}

type OffsetStart interface {
	Expr
	offsetStartNode()
}

type OffsetInt int64

func (i OffsetInt) Walk(w *writer.Writer) {
	w.Write(strconv.FormatInt(int64(i), 10))
}

func (LimitInt) limitCountNode()            {}
func (LimitUint) limitCountNode()           {}
func (PositionalParameter) limitCountNode() {}

func (LimitInt) exprNode()  {}
func (LimitUint) exprNode() {}

func (OffsetInt) offsetStartNode()           {}
func (PositionalParameter) offsetStartNode() {}

type OrderClause struct {
	List []OrderBy
}

func (c OrderClause) Walk(w *writer.Writer) {
	if len(c.List) == 0 {
		return
	}
	w.NewLine()
	w.Write("ORDER BY ")
	for i, o := range c.List {
		if i > 0 {
			w.Write(", ")
		}
		o.Walk(w)
	}
}

type OrderBy struct {
	Column     ColumnIdent
	Desc       bool
	NullsFirst bool
}

func (o OrderBy) Walk(w *writer.Writer) {
	o.Column.Walk(w)
	if o.Desc {
		w.Write(" DESC")
	} else {
		w.Write(" ASC")
	}

	if o.NullsFirst {
		w.Write(" NULLS FIRST")
	} else {
		w.Write(" NULLS LAST")
	}
}

type JoinClause struct {
	List []*TableJoin
}

func (c JoinClause) Walk(w *writer.Writer) {
	if len(c.List) == 0 {
		return
	}

	for _, j := range c.List {
		w.NewLine()
		j.Walk(w)
	}
}