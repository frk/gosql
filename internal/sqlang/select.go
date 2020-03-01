package sqlang

import (
	"strconv"

	"github.com/frk/gosql/internal/writer"
)

type SelectStatement struct {
	Columns ValueExpr
	Table   Ident
	Join    JoinClause
	Where   WhereClause
	Order   OrderClause
	Limit   LimitClause
	Offset  OffsetClause
}

func (s SelectStatement) Walk(w *writer.Writer) {
	w.NoIndent()
	w.Write("SELECT")
	w.NewLine()
	w.Indent()
	s.Columns.Walk(w)
	w.NewLine()
	w.Write("FROM ")
	s.Table.Walk(w)
	s.Join.Walk(w)
	w.NewLine()
	s.Where.Walk(w)
	s.Order.Walk(w)
	s.Limit.Walk(w)
	s.Offset.Walk(w)
}

type SelectExistsStatement struct {
	Table  Ident
	Join   JoinClause
	Where  WhereClause
	Order  OrderClause
	Limit  LimitClause
	Offset OffsetClause
	Open   bool
	Not    bool
}

func (s SelectExistsStatement) Walk(w *writer.Writer) {
	w.NoIndent()
	w.Write("SELECT ")
	if s.Not {
		w.Write("NOT ")
	}
	w.Write("EXISTS(SELECT 1 FROM ")
	s.Table.Walk(w)
	w.Indent()
	s.Join.Walk(w)
	w.NewLine()
	s.Where.Walk(w)
	s.Order.Walk(w)
	s.Limit.Walk(w)
	s.Offset.Walk(w)

	if !s.Open {
		w.Write(")")
	}
}

type SelectCountStatement struct {
	Table  Ident
	Join   JoinClause
	Where  WhereClause
	Order  OrderClause
	Limit  LimitClause
	Offset OffsetClause
}

func (s SelectCountStatement) Walk(w *writer.Writer) {
	w.NoIndent()
	w.Write("SELECT COUNT(*) FROM ")
	s.Table.Walk(w)
	w.Indent()
	s.Join.Walk(w)
	w.NewLine()
	s.Where.Walk(w)
	s.Order.Walk(w)
	s.Limit.Walk(w)
	s.Offset.Walk(w)
}

type LimitClause struct {
	Value LimitValue
}

func (l LimitClause) Walk(w *writer.Writer) {
	if l.Value == nil {
		return
	}
	w.NewLine()
	w.Write("LIMIT ")
	l.Value.Walk(w)
}

type LimitValue interface {
	Node
	limitValueNode()
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
	Value OffsetValue
}

func (l OffsetClause) Walk(w *writer.Writer) {
	if l.Value == nil {
		return
	}
	w.NewLine()
	w.Write("OFFSET ")
	l.Value.Walk(w)
}

type OffsetValue interface {
	Node
	offsetValueNode()
}

type OffsetInt int64

func (i OffsetInt) Walk(w *writer.Writer) {
	w.Write(strconv.FormatInt(int64(i), 10))
}

type OffsetUint uint64

func (u OffsetUint) Walk(w *writer.Writer) {
	w.Write(strconv.FormatUint(uint64(u), 10))
}

func (LimitInt) exprNode()   {}
func (LimitUint) exprNode()  {}
func (OffsetInt) exprNode()  {}
func (OffsetUint) exprNode() {}

func (LimitInt) limitValueNode()    {}
func (LimitUint) limitValueNode()   {}
func (OffsetInt) offsetValueNode()  {}
func (OffsetUint) offsetValueNode() {}

func (OrdinalParameterSpec) limitValueNode()  {}
func (OrdinalParameterSpec) offsetValueNode() {}

func (DynamicParmeterSpec) limitValueNode()  {}
func (DynamicParmeterSpec) offsetValueNode() {}

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
	Column     ColumnReference
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
	List []TableJoin
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
