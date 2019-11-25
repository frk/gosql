// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

func (q *DeleteWithReturningSliceErrorInfoHandlerQuery) Exec(c gosql.Conn) error {
	const queryString = `DELETE FROM "test_user" AS u
	WHERE u."created_at" < $1
	RETURNING
	u."id"
	, u."email"
	, u."full_name"
	, u."created_at"` // `

	rows, err := c.Query(queryString, q.Where.CreatedBefore)
	if err != nil {
		return q.eh.HandleErrorInfo(&gosql.ErrorInfo{Error: err, Query: queryString, SpecKind: "Delete", SpecName: "DeleteWithReturningSliceErrorInfoHandlerQuery", SpecValue: q})
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
			return q.eh.HandleErrorInfo(&gosql.ErrorInfo{Error: err, Query: queryString, SpecKind: "Delete", SpecName: "DeleteWithReturningSliceErrorInfoHandlerQuery", SpecValue: q})
		}

		q.Users = append(q.Users, v)
	}
	if err := rows.Err(); err != nil {
		return q.eh.HandleErrorInfo(&gosql.ErrorInfo{Error: err, Query: queryString, SpecKind: "Delete", SpecName: "DeleteWithReturningSliceErrorInfoHandlerQuery", SpecValue: q})
	}
	return nil
}
