// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
)

func (q *SelectCoalesceTableQuery) Exec(c gosql.Conn) error {
	const queryString = `SELECT
	COALESCE(c."col_a", ''::text)
	, c."col_b"
	, COALESCE(c."col_c", 0::integer)
	, COALESCE(c."col_d", '0001-01-01 00:00:00'::timestamp without time zone)
	, c."col_e"
	, COALESCE(c."col_a", 'foo'::text)
	, COALESCE(c."col_b", ''::text)
	, COALESCE(c."col_e", 'red'::color_enum)
	FROM "test_coalesce" AS c
	LIMIT 1` // `

	row := c.QueryRow(queryString)
	return row.Scan(
		&q.Columns.A,
		&q.Columns.B,
		&q.Columns.C,
		&q.Columns.D,
		&q.Columns.E,
		&q.Columns.A2,
		&q.Columns.B2,
		&q.Columns.E2,
	)
}
