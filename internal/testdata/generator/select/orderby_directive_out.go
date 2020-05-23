// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql/testdata/common"

	"github.com/frk/gosql"
)

func (q *SelectWithOrderByDirectiveQuery) Exec(c gosql.Conn) error {
	const queryString = `SELECT
	u."id"
	, u."email"
	, u."full_name"
	, u."created_at"
	FROM "test_user" AS u
	ORDER BY u."created_at" DESC NULLS LAST` // `

	rows, err := c.Query(queryString)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		v := new(common.User)
		err := rows.Scan(
			&v.Id,
			&v.Email,
			&v.FullName,
			&v.CreatedAt,
		)
		if err != nil {
			return err
		}

		q.Users = append(q.Users, v)
	}
	return rows.Err()
}
