package filter

import (
	"github.com/frk/gosql"
)

// An sqlNode representing a IN comparison.
type sqlIn[T any] struct {
	col  string // the LHS column
	vals []T    // the list of values
}

func (sqlIn[T]) canAndOr() bool { return true }

func (n sqlIn[T]) write(w *sqlWriter) {
	w.WriteString(n.col)
	w.WriteString(" IN (")
	for i := range n.vals {
		if i > 0 {
			w.WriteString(",")
		}
		w.WriteString(gosql.OrdinalParameters[w.p])
		w.params = append(w.params, n.vals[i])
		w.p += 1
	}
	w.WriteString(")")

}

func (c *Constructor) ColInInt64s(column string, values []int64) {
	if c.canAndOr() {
		c.filter.where = append(c.filter.where, sqlAnd{})
	}
	c.filter.where = append(c.filter.where, sqlIn[int64]{
		col:  column,
		vals: values,
	})
}

func (c *Constructor) ColInInts(column string, values []int) {
	if c.canAndOr() {
		c.filter.where = append(c.filter.where, sqlAnd{})
	}
	c.filter.where = append(c.filter.where, sqlIn[int]{
		col:  column,
		vals: values,
	})
}

func (c *Constructor) ColInInt16s(column string, values []int16) {
	if c.canAndOr() {
		c.filter.where = append(c.filter.where, sqlAnd{})
	}
	c.filter.where = append(c.filter.where, sqlIn[int16]{
		col:  column,
		vals: values,
	})
}

func (c *Constructor) ColInStrings(column string, values []string) {
	if c.canAndOr() {
		c.filter.where = append(c.filter.where, sqlAnd{})
	}
	c.filter.where = append(c.filter.where, sqlIn[string]{
		col:  column,
		vals: values,
	})
}

////////////////////////////////////////////////////////////////////////////////

// An sqlNode representing a "= ANY(array)" comparison.
type sqlAny struct {
	col  string // the LHS column
	vals any    // the list of values
	cast string // the required type CAST
}

func (sqlAny) canAndOr() bool { return true }

func (n sqlAny) write(w *sqlWriter) {
	w.WriteString(n.col)
	w.WriteString(" = ANY(")
	w.WriteString(gosql.OrdinalParameters[w.p] + n.cast)
	w.WriteString(")")

	w.params = append(w.params, n.vals)
	w.p += 1
}

func (c *Constructor) ColAnyInt64s(column string, values []int64) {
	if c.canAndOr() {
		c.filter.where = append(c.filter.where, sqlAnd{})
	}
	c.filter.where = append(c.filter.where, sqlAny{
		col:  column,
		vals: values,
		cast: "::int8[]",
	})
}

func (c *Constructor) ColAnyInts(column string, values []int) {
	if c.canAndOr() {
		c.filter.where = append(c.filter.where, sqlAnd{})
	}
	c.filter.where = append(c.filter.where, sqlAny{
		col:  column,
		vals: values,
		cast: "::int4[]",
	})
}

func (c *Constructor) ColAnyInt16s(column string, values []int16) {
	if c.canAndOr() {
		c.filter.where = append(c.filter.where, sqlAnd{})
	}
	c.filter.where = append(c.filter.where, sqlAny{
		col:  column,
		vals: values,
		cast: "::int2[]",
	})
}

func (c *Constructor) ColAnyStrings(column string, values []string) {
	if c.canAndOr() {
		c.filter.where = append(c.filter.where, sqlAnd{})
	}
	c.filter.where = append(c.filter.where, sqlAny{
		col:  column,
		vals: values,
		cast: "::text[]",
	})
}
