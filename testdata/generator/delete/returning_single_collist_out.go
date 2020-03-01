// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
)

func (q *DeleteWithReturningSingleCollistQuery) Exec(c gosql.Conn) error {
	const queryString = `DELETE FROM "test_user" AS u
	WHERE u."id" = $1
	RETURNING u."email", u."full_name"` // `

	row := c.QueryRow(queryString, q.Where.Id)
	return row.Scan(&q.User.Email, &q.User.FullName)
}
