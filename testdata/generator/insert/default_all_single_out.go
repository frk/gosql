// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
)

func (q *InsertDefaultAllSingleQuery) Exec(c gosql.Conn) error {
	const queryString = `INSERT INTO "test_user_with_defaults" AS u (
		"email"
		, "full_name"
		, "is_active"
		, "created_at"
		, "updated_at"
	) VALUES (
		DEFAULT
		, DEFAULT
		, DEFAULT
		, DEFAULT
		, DEFAULT
	)` // `

	_, err := c.Exec(queryString)
	return err
}