// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
)

func (q *SelectNotExistsWithFilterQuery) Exec(c gosql.Conn) error {
	var queryString = `SELECT NOT EXISTS(SELECT 1 FROM "test_user" AS u
	` // `

	queryString += q.Filter.ToSQL()
	queryString += `)`

	params := q.Filter.Params()
	row := c.QueryRow(queryString, params...)
	return row.Scan(&q.NotExists)
}