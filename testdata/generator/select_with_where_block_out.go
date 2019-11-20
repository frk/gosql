// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

func (q *SelectWithWhereBlockQuery) Exec(c gosql.Conn) error {
	const queryString = `SELECT
	u."id"
	, u."email"
	, u."full_name"
	, u."created_at"
	FROM "users_table" AS u
	WHERE u."id" = $1` //`

	row := c.QueryRow(queryString, q.Where.Id)

	q.User = new(common.User)
	return row.Scan(
		&q.User.Id,
		&q.User.Email,
		&q.User.FullName,
		&q.User.CreatedAt,
	)
}
