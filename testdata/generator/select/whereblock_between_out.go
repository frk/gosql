// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

func (q *SelectWithWhereBlockBetweenQuery) Exec(c gosql.Conn) error {
	const queryString = `SELECT
	u."id"
	, u."email"
	, u."full_name"
	, u."created_at"
	FROM "test_user" AS u
	WHERE u."created_at" BETWEEN $1 AND $2` // `

	rows, err := c.Query(queryString, q.Where.CreatedAt.After, q.Where.CreatedAt.Before)
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
