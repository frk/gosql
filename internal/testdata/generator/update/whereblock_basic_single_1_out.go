// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
)

func (q *UpdateWhereblockBasicSingle1Query) Exec(c gosql.Conn) error {
	const queryString = `UPDATE "test_user" SET (
		"email"
		, "full_name"
		, "is_active"
		, "created_at"
		, "updated_at"
	) = (
		$1
		, $2
		, $3
		, $4
		, $5
	)
	WHERE "id" = $6` // `

	_, err := c.Exec(queryString,
		q.User.Email,
		q.User.FullName,
		q.User.IsActive,
		q.User.CreatedAt,
		q.User.UpdatedAt,
		q.Where.Id,
	)
	return err
}
