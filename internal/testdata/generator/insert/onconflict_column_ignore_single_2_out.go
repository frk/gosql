// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
)

func (q *InsertOnConflictColumnIgnoreSingle2Query) Exec(c gosql.Conn) error {
	const queryString = `INSERT INTO "test_onconflict" AS k (
		"key"
		, "name"
		, "fruit"
		, "value"
	) VALUES (
		NULLIF($1, 0)::integer
		, NULLIF($2, '')::text
		, NULLIF($3, '')::text
		, NULLIF($4, 0)::double precision
	)
	ON CONFLICT (key, name)
	DO NOTHING` // `

	_, err := c.Exec(queryString,
		q.Data.Key,
		q.Data.Name,
		q.Data.Fruit,
		q.Data.Value,
	)
	return err
}
