// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
)

func (q *UpdatePKeySingleQuery) Exec(c gosql.Conn) error {
	const queryString = `UPDATE "test_user" SET (
		"email"
		, "password"
		, "created_at"
		, "updated_at"
	) = (
		$1
		, $2
		, $3
		, $4
	)
	WHERE "id" = $5` // `

	_, err := c.Exec(queryString,
		q.User.Email,
		q.User.Password,
		q.User.CreatedAt,
		q.User.UpdatedAt,
		q.User.Id,
	)
	return err
}
