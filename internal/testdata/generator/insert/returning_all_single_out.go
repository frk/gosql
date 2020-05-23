// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
)

func (q *InsertReturningAllSingleQuery) Exec(c gosql.Conn) error {
	const queryString = `INSERT INTO "test_user" AS u (
		"email"
		, "password"
		, "created_at"
		, "updated_at"
	) VALUES (
		$1
		, $2
		, $3
		, $4
	)
	RETURNING
	u."id"
	, u."email"
	, u."created_at"
	, u."updated_at"` // `

	row := c.QueryRow(queryString,
		q.User.Email,
		q.User.Password,
		q.User.CreatedAt,
		q.User.UpdatedAt,
	)
	return row.Scan(
		&q.User.Id,
		&q.User.Email,
		&q.User.CreatedAt,
		&q.User.UpdatedAt,
	)
}
