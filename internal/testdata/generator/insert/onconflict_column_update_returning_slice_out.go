// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
)

func (q *InsertOnConflictColumnUpdateReturningSliceQuery) Exec(c gosql.Conn) error {
	var queryString = `INSERT INTO "test_onconflict" AS k (
		"key"
		, "name"
		, "fruit"
		, "value"
	) VALUES ` // `

	params := make([]interface{}, len(q.Data)*4)
	for i, v := range q.Data {
		pos := i * 4

		params[pos+0] = v.Key
		params[pos+1] = v.Name
		params[pos+2] = v.Fruit
		params[pos+3] = v.Value

		queryString += `(NULLIF(` + gosql.OrdinalParameters[pos+0] + `, 0)::integer` +
			`, NULLIF(` + gosql.OrdinalParameters[pos+1] + `, '')::text` +
			`, NULLIF(` + gosql.OrdinalParameters[pos+2] + `, '')::text` +
			`, NULLIF(` + gosql.OrdinalParameters[pos+3] + `, 0)::double precision` +
			`),`
	}

	queryString = queryString[:len(queryString)-1]
	queryString += ` ON CONFLICT (key)
	DO UPDATE SET "fruit" = EXCLUDED."fruit"
	RETURNING
	k."id"
	, COALESCE(k."key", 0::integer)
	, COALESCE(k."name", ''::text)
	, COALESCE(k."fruit", ''::text)
	, COALESCE(k."value", 0::double precision)` // `

	rows, err := c.Query(queryString, params...)
	if err != nil {
		return err
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		err := rows.Scan(
			&q.Data[i].Id,
			&q.Data[i].Key,
			&q.Data[i].Name,
			&q.Data[i].Fruit,
			&q.Data[i].Value,
		)
		if err != nil {
			return err
		}

		i += 1
	}
	return rows.Err()
}
