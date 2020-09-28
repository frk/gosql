// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

func (q *UpdateFilterResultSliceQuery) Exec(c gosql.Conn) error {
	var queryString = `UPDATE "test_user" AS u SET (
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
	)` // `

	filterString, params := q.Filter.ToSQL()
	queryString += filterString
	queryString += ` RETURNING
	u."id"
	, u."email"
	, u."full_name"
	, u."is_active"
	, u."created_at"
	, u."updated_at"` // `

	params = append([]interface{}{
		q.User.Email,
		q.User.FullName,
		q.User.IsActive,
		q.User.CreatedAt,
		q.User.UpdatedAt,
	}, params...)

	rows, err := c.Query(queryString, params...)
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
