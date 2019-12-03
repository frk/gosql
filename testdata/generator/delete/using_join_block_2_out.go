// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
)

func (q *DeleteWithUsingJoinBlock2Query) Exec(c gosql.Conn) error {
	const queryString = `DELETE FROM "test_user" AS u
	USING "test_post" AS p
	LEFT JOIN "test_join1" AS j1 ON j1."post_id" = p."id"
	RIGHT JOIN "test_join2" AS j2 ON j2."join1_id" = j1."id"
	FULL JOIN "test_join3" AS j3 ON j3."join2_id" = j2."id"
	CROSS JOIN "test_join4" AS j4
	WHERE u."id" = p."user_id" AND p."is_spam" IS TRUE` // `

	_, err := c.Exec(queryString)
	return err
}
