// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
)

func (q *UpdateWhereblockBasicSingle2Query) Exec(c gosql.Conn) error {
	const queryString = `UPDATE "test_user" AS u SET (
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
	WHERE u."created_at" > $6
	AND u."created_at" < $7
	AND (u."full_name" = 'John Doe' OR u."full_name" = 'Jane Doe')` // `

	_, err := c.Exec(queryString,
		q.User.Email,
		q.User.FullName,
		q.User.IsActive,
		q.User.CreatedAt,
		q.User.UpdatedAt,
		q.Where.CreatedAfter,
		q.Where.CreatedBefore,
	)
	return err
}
