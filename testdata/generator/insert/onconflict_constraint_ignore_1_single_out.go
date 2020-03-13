// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
)

func (q *InsertOnConflictConstraintIgnore1SingleQuery) Exec(c gosql.Conn) error {
	const queryString = `INSERT INTO "test_onconflict" AS k (
		"key"
		, "name"
		, "fruit"
		, "value"
	) VALUES (
		$1
		, $2
		, $3
		, $4
	)
	ON CONFLICT ON CONSTRAINT "test_onconflict_key_value_key"
	DO NOTHING` // `

	_, err := c.Exec(queryString,
		q.Data.Key,
		q.Data.Name,
		q.Data.Fruit,
		q.Data.Value,
	)
	return err
}
