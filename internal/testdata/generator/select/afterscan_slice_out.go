// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

func (q *SelectWithAfterScanSliceQuery) Exec(c gosql.Conn) error {
	const queryString = `SELECT
	u."id"
	, u."email"
	, u."full_name"
	, u."created_at"
	FROM "test_user" AS u
	` // `

	rows, err := c.Query(queryString)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		v := new(common.User2)
		err := rows.Scan(
			&v.Id,
			&v.Email,
			&v.FullName,
			&v.CreatedAt,
		)
		if err != nil {
			return err
		}

		v.AfterScan()
		q.Users = append(q.Users, v)
	}
	return rows.Err()
}
