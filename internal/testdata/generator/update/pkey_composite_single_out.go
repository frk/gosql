// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
)

func (q *UpdatePKeyCompositeSingleQuery) Exec(c gosql.Conn) error {
	const queryString = `UPDATE "test_composite_pkey" AS p SET (
		"key"
		, "name"
		, "fruit"
		, "value"
	) = (
		$1
		, $2
		, NULLIF($3, '')::text
		, NULLIF($4, 0)::double precision
	)
	WHERE p."id" = $5 AND p."key" = $1 AND p."name" = $2` // `

	_, err := c.Exec(queryString,
		q.Data.Key,
		q.Data.Name,
		q.Data.Fruit,
		q.Data.Value,
		q.Data.Id,
	)
	return err
}
