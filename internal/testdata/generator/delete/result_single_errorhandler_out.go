// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql/testdata/common"

	"github.com/frk/gosql"
)

func (q *DeleteWithResultSingleErrorHandlerQuery) Exec(c gosql.Conn) error {
	const queryString = `DELETE FROM "test_user" AS u
	WHERE u."id" = $1
	RETURNING
	u."id"
	, u."email"
	, u."full_name"
	, u."created_at"` // `

	row := c.QueryRow(queryString, q.Where.Id)

	q.Result = new(common.User)
	return q.eh.HandleError(row.Scan(
		&q.Result.Id,
		&q.Result.Email,
		&q.Result.FullName,
		&q.Result.CreatedAt,
	))
}
