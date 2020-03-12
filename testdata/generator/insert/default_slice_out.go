// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
)

func (q *InsertSliceDefaultQuery) Exec(c gosql.Conn) error {
	var queryString = `INSERT INTO "test_user_with_defaults" AS u (
		"email"
		, "full_name"
		, "is_active"
		, "created_at"
		, "updated_at"
	) VALUES ` // `

	params := make([]interface{}, len(q.Users)*3)
	for i, v := range q.Users {
		pos := i * 3

		params[pos+0] = v.Email
		params[pos+1] = v.IsActive
		params[pos+2] = v.UpdatedAt

		queryString += `(` + gosql.OrdinalParameters[pos+0] +
			`, DEFAULT` +
			`, ` + gosql.OrdinalParameters[pos+1] +
			`, DEFAULT` +
			`, ` + gosql.OrdinalParameters[pos+2] +
			`),`
	}

	queryString = queryString[:len(queryString)-1]

	_, err := c.Exec(queryString, params...)
	return err
}
