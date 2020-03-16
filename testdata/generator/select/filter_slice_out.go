// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

func (q *SelectWithFilterSliceQuery) Exec(c gosql.Conn) error {
	var queryString = `SELECT
	u."id"
	, u."email"
	, u."full_name"
	, u."created_at"
	FROM "test_user" AS u
	` // `

	queryString += q.Filter.ToSQL()

	params := q.Filter.Params()
	rows, err := c.Query(queryString, params...)
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
