// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
)

func (q *DeleteWithFilterQuery) Exec(c gosql.Conn) error {
	var queryString = `DELETE FROM "test_user"` // `

	filterString, params := q.Filter.ToSQL()
	queryString += filterString

	_, err := c.Exec(queryString, params...)
	return err
}
