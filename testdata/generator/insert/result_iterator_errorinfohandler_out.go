// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

func (q *InsertWithResultIteratorErrorInfoHandlerQuery) Exec(c gosql.Conn) error {
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
		return q.erh.HandleErrorInfo(&gosql.ErrorInfo{Error: err, Query: queryString, SpecKind: "Insert", SpecName: "InsertWithResultIteratorErrorInfoHandlerQuery", SpecValue: q})
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
			return q.erh.HandleErrorInfo(&gosql.ErrorInfo{Error: err, Query: queryString, SpecKind: "Insert", SpecName: "InsertWithResultIteratorErrorInfoHandlerQuery", SpecValue: q})
		}

		v.AfterScan()
		if err := q.result.NextUser(v); err != nil {
			return q.erh.HandleErrorInfo(&gosql.ErrorInfo{Error: err, Query: queryString, SpecKind: "Insert", SpecName: "InsertWithResultIteratorErrorInfoHandlerQuery", SpecValue: q})
		}
	}
	if err := rows.Err(); err != nil {
		return q.erh.HandleErrorInfo(&gosql.ErrorInfo{Error: err, Query: queryString, SpecKind: "Insert", SpecName: "InsertWithResultIteratorErrorInfoHandlerQuery", SpecValue: q})
	}
	return nil
}
