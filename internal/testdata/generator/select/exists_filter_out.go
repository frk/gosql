// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
)

func (q *SelectExistsWithFilterQuery) Exec(c gosql.Conn) error {
	var queryString = `SELECT EXISTS(SELECT 1 FROM "test_user" AS u
	` // `

	filterString, params := q.Filter.ToSQL(0)
	queryString += filterString + `)`

	row := c.QueryRow(queryString, params...)
	return row.Scan(&q.Exists)
}
