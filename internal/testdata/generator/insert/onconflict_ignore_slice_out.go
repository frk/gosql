// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
)

func (q *InsertOnConflictIgnoreSliceQuery) Exec(c gosql.Conn) error {
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

		queryString += `(` + gosql.OrdinalParameters[pos+0] +
			`, ` + gosql.OrdinalParameters[pos+1] +
			`, ` + gosql.OrdinalParameters[pos+2] +
			`, ` + gosql.OrdinalParameters[pos+3] +
			`),`
	}

	queryString = queryString[:len(queryString)-1]
	queryString += ` ON CONFLICT DO NOTHING` // `

	_, err := c.Exec(queryString, params...)
	return err
}