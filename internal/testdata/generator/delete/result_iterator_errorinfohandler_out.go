// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
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
		return q.errhandler.HandleErrorInfo(&gosql.ErrorInfo{Error: err, QueryString: queryString, QueryKind: "Delete", QueryName: "DeleteWithResultIteratorErrorInfoHandlerQuery", QueryValue: q})
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
			return q.errhandler.HandleErrorInfo(&gosql.ErrorInfo{Error: err, QueryString: queryString, QueryKind: "Delete", QueryName: "DeleteWithResultIteratorErrorInfoHandlerQuery", QueryValue: q})
		}

		v.AfterScan()
		if err := q.Result.NextUser(v); err != nil {
			return q.errhandler.HandleErrorInfo(&gosql.ErrorInfo{Error: err, QueryString: queryString, QueryKind: "Delete", QueryName: "DeleteWithResultIteratorErrorInfoHandlerQuery", QueryValue: q})
		}
	}
	if err := rows.Err(); err != nil {
		return q.errhandler.HandleErrorInfo(&gosql.ErrorInfo{Error: err, QueryString: queryString, QueryKind: "Delete", QueryName: "DeleteWithResultIteratorErrorInfoHandlerQuery", QueryValue: q})
	}
	return nil
}
