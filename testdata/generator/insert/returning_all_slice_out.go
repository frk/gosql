// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
)

func (q *InsertReturningAllSliceQuery) Exec(c gosql.Conn) error {
	var queryString = `INSERT INTO "test_user" AS u (
		"id"
		, "email"
		, "full_name"
		, "created_at"
	) VALUES ` // `

	params := make([]interface{}, len(q.Users)*4)
	for i, v := range q.Users {
		pos := i * 4

		params[pos+0] = v.Id
		params[pos+1] = v.Email
		params[pos+2] = v.FullName
		params[pos+3] = v.CreatedAt

		queryString += `(` + gosql.OrdinalParameters[pos+0] +
			`, ` + gosql.OrdinalParameters[pos+1] +
			`, ` + gosql.OrdinalParameters[pos+2] +
			`, ` + gosql.OrdinalParameters[pos+3] +
			`),`
	}

	queryString = queryString[:len(queryString)-1]
	queryString += ` RETURNING
	u."id"
	, u."email"
	, u."full_name"
	, u."created_at"` // `

	rows, err := c.Query(queryString, params...)
	if err != nil {
		return err
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		err := rows.Scan(
			&q.Users[i].Id,
			&q.Users[i].Email,
			&q.Users[i].FullName,
			&q.Users[i].CreatedAt,
		)
		if err != nil {
			return err
		}

		i += 1
	}
	return rows.Err()
}