// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
)

func (q *InsertWithReturningSliceErrorHandlerQuery) Exec(c gosql.Conn) error {
	var queryString = `INSERT INTO "test_user" AS u (
		"email"
		, "password"
		, "created_at"
		, "updated_at"
	) VALUES ` // `

	params := make([]interface{}, len(q.Users)*4)
	for i, v := range q.Users {
		pos := i * 4

		params[pos+0] = v.Email
		params[pos+1] = v.Password
		params[pos+2] = v.CreatedAt
		params[pos+3] = v.UpdatedAt

		queryString += `(` + gosql.OrdinalParameters[pos+0] +
			`, ` + gosql.OrdinalParameters[pos+1] +
			`, ` + gosql.OrdinalParameters[pos+2] +
			`, ` + gosql.OrdinalParameters[pos+3] +
			`),`
	}

	queryString = queryString[:len(queryString)-1]
	queryString += ` RETURNING u."id"` // `

	rows, err := c.Query(queryString, params...)
	if err != nil {
		return q.erh.HandleError(err)
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		err := rows.Scan(&q.Users[i].Id)
		if err != nil {
			return q.erh.HandleError(err)
		}

		i += 1
	}
	return q.erh.HandleError(rows.Err())
}
