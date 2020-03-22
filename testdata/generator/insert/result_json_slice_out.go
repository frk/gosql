// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql/testdata/common"

	"github.com/frk/gosql"
)

func (q *InsertResultJSONSliceQuery) Exec(c gosql.Conn) error {
	var queryString = `INSERT INTO "test_user" AS u (
		"email"
		, "full_name"
		, "is_active"
		, "metadata1"
		, "metadata2"
		, "created_at"
		, "updated_at"
	) VALUES ` // `

	params := make([]interface{}, len(q.Users)*7)
	for i, v := range q.Users {
		pos := i * 7

		params[pos+0] = v.Email
		params[pos+1] = v.FullName
		params[pos+2] = v.IsActive
		params[pos+3] = gosql.JSON(v.Metadata1)
		params[pos+4] = gosql.JSON(v.Metadata2)
		params[pos+5] = v.CreatedAt
		params[pos+6] = v.UpdatedAt

		queryString += `(` + gosql.OrdinalParameters[pos+0] +
			`, ` + gosql.OrdinalParameters[pos+1] +
			`, ` + gosql.OrdinalParameters[pos+2] +
			`, ` + gosql.OrdinalParameters[pos+3] +
			`, ` + gosql.OrdinalParameters[pos+4] +
			`, ` + gosql.OrdinalParameters[pos+5] +
			`, ` + gosql.OrdinalParameters[pos+6] +
			`),`
	}

	queryString = queryString[:len(queryString)-1]
	queryString += ` RETURNING
	u."id"
	, u."email"
	, u."full_name"
	, u."is_active"
	, u."metadata1"
	, u."metadata2"
	, u."created_at"
	, u."updated_at"` // `

	rows, err := c.Query(queryString, params...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		v := new(common.User5)
		err := rows.Scan(
			&v.Id,
			&v.Email,
			&v.FullName,
			&v.IsActive,
			gosql.JSON(&v.Metadata1),
			gosql.JSON(&v.Metadata2),
			&v.CreatedAt,
			&v.UpdatedAt,
		)
		if err != nil {
			return err
		}

		q.Result = append(q.Result, v)
	}
	return rows.Err()
}
