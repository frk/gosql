// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

func (q *InsertWithResultSingleAfterScanQuery) Exec(c gosql.Conn) error {
	const queryString = `INSERT INTO "test_user" AS u (
		"id"
		, "email"
		, "full_name"
		, "created_at"
	) VALUES (
		$1
		, $2
		, $3
		, $4
	)
	RETURNING
	u."id"
	, u."email"
	, u."full_name"
	, u."created_at"` // `

	row := c.QueryRow(queryString,
		q.User.Id,
		q.User.Email,
		q.User.FullName,
		q.User.CreatedAt,
	)

	q.result = new(common.User2)
	err := row.Scan(
		&q.result.Id,
		&q.result.Email,
		&q.result.FullName,
		&q.result.CreatedAt,
	)
	if err != nil {
		return err
	}

	q.result.AfterScan()
	return nil
}
