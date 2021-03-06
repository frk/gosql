// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

func (q *UpdateWhereblockResultSliceQuery) Exec(c gosql.Conn) error {
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
	WHERE u."created_at" > $6 AND u."created_at" < $7
	RETURNING
	u."id"
	, u."email"
	, u."full_name"
	, u."is_active"
	, u."created_at"
	, u."updated_at"` // `

	rows, err := c.Query(queryString,
		q.User.Email,
		q.User.FullName,
		q.User.IsActive,
		q.User.CreatedAt,
		q.User.UpdatedAt,
		q.Where.CreatedAfter,
		q.Where.CreatedBefore,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		v := new(common.User4)
		err := rows.Scan(
			&v.Id,
			&v.Email,
			&v.FullName,
			&v.IsActive,
			&v.CreatedAt,
			&v.UpdatedAt,
		)
		if err != nil {
			return err
		}

		q.Result = append(q.Result, v)
	}
	return rows.Err()
}
