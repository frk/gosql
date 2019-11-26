// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

func (q *SelectWithJoinBlockSliceQuery) Exec(c gosql.Conn) error {
	const queryString = `SELECT
	u."id"
	, u."email"
	, u."full_name"
	, u."created_at"
	FROM "test_user" AS u
	LEFT JOIN "test_post" AS p ON p."user_id" = u."id"
	LEFT JOIN "test_join1" AS j1 ON j1."post_id" = p."id"
	RIGHT JOIN "test_join2" AS j2 ON j2."join1_id" = j1."id"
	FULL JOIN "test_join3" AS j3 ON j3."join2_id" = j2."id"
	CROSS JOIN "test_join4" AS j4
	WHERE p."is_spam"` // `

	rows, err := c.Query(queryString)
	if err != nil {
		return err
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
			return err
		}

		q.Users = append(q.Users, v)
	}
	return rows.Err()
}
