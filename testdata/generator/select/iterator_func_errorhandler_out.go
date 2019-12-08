// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

func (q *SelectWithIteratorFuncErrorHandlerQuery) Exec(c gosql.Conn) error {
	const queryString = `SELECT
	u."id"
	, u."email"
	, u."full_name"
	, u."created_at"
	FROM "test_user" AS u
	WHERE u."created_at" > $1
	LIMIT $2
	OFFSET $3` // `

	rows, err := c.Query(queryString,
		q.where.createdafter,
		q.limit,
		q.offset,
	)
	if err != nil {
		return q.erh.HandleError(err)
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
			return q.erh.HandleError(err)
		}

		if err := q.next(v); err != nil {
			return q.erh.HandleError(err)
		}
	}
	return q.erh.HandleError(rows.Err())
}
