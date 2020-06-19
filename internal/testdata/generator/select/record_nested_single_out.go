// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

func (q *SelectWithRecordNestedSingleQuery) Exec(c gosql.Conn) error {
	const queryString = `SELECT
	n."foo_bar_baz_val"
	, n."foo_baz_val"
	, n."foo2_bar_baz_val"
	, n."foo2_baz_val"
	FROM "test_nested" AS n
	LIMIT 1` // `

	row := c.QueryRow(queryString)

	q.Nested = new(common.Nested)
	q.Nested.FOO = new(common.Foo)
	q.Nested.FOO.Baz = new(common.Baz)
	q.Nested.Foo.Baz = new(common.Baz)
	return row.Scan(
		&q.Nested.FOO.Bar.Baz.Val,
		&q.Nested.FOO.Baz.Val,
		&q.Nested.Foo.Bar.Baz.Val,
		&q.Nested.Foo.Baz.Val,
	)
}
