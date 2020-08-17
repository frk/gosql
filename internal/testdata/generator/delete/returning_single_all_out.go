// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

func (q *DeleteWithReturningSingleAllQuery) Exec(c gosql.Conn) error {
	const queryString = `DELETE FROM "test_user" AS u
	WHERE u."id" = $1
	RETURNING
	u."id"
	, u."email"
	, u."full_name"
	, u."created_at"` // `

	row := c.QueryRow(queryString, q.Where.Id)

	q.User = new(common.User)
	return row.Scan(
		&q.User.Id,
		&q.User.Email,
		&q.User.FullName,
		&q.User.CreatedAt,
	)
}