// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
)

func (q *InsertDefaultAllSliceQuery) Exec(c gosql.Conn) error {
	var queryString = `INSERT INTO "test_user_with_defaults" AS u (
		"email"
		, "full_name"
		, "is_active"
		, "created_at"
		, "updated_at"
	) VALUES ` // `

	for _, _ = range q.Users {
		queryString += `(DEFAULT` +
			`, DEFAULT` +
			`, DEFAULT` +
			`, DEFAULT` +
			`, DEFAULT` +
			`),`
	}

	queryString = queryString[:len(queryString)-1]

	_, err := c.Exec(queryString)
	return err
}
