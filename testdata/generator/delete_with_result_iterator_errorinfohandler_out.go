// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

func (q *DeleteWithResultIteratorErrorInfoHandlerQuery) Exec(c gosql.Conn) error {
	const queryString = `DELETE FROM "test_user" AS u
	WHERE u."created_at" < $1
	RETURNING
	u."id"
	, u."email"
	, u."full_name"
	, u."created_at"` // `

	rows, err := c.Query(queryString, q.Where.CreatedBefore)
	if err != nil {
		return q.errhandler.HandleErrorInfo(&gosql.ErrorInfo{Error: err, Query: queryString, SpecKind: "Delete", SpecName: "DeleteWithResultIteratorErrorInfoHandlerQuery", SpecValue: q})
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
			return q.errhandler.HandleErrorInfo(&gosql.ErrorInfo{Error: err, Query: queryString, SpecKind: "Delete", SpecName: "DeleteWithResultIteratorErrorInfoHandlerQuery", SpecValue: q})
		}

		v.AfterScan()
		if err := q.Result.NextUser(v); err != nil {
			return q.errhandler.HandleErrorInfo(&gosql.ErrorInfo{Error: err, Query: queryString, SpecKind: "Delete", SpecName: "DeleteWithResultIteratorErrorInfoHandlerQuery", SpecValue: q})
		}
	}
	if err := rows.Err(); err != nil {
		return q.errhandler.HandleErrorInfo(&gosql.ErrorInfo{Error: err, Query: queryString, SpecKind: "Delete", SpecName: "DeleteWithResultIteratorErrorInfoHandlerQuery", SpecValue: q})
	}
	return nil
}
