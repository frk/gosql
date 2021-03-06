// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

func (q *DeleteWithReturningSliceCollistQuery) Exec(c gosql.Conn) error {
	const queryString = `DELETE FROM "test_user" AS u
	WHERE u."created_at" < $1
	RETURNING u."email", u."full_name"` // `

	rows, err := c.Query(queryString, q.Where.CreatedBefore)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var v common.User
		err := rows.Scan(&v.Email, &v.FullName)
		if err != nil {
			return err
		}

		q.Users = append(q.Users, v)
	}
	return rows.Err()
}
