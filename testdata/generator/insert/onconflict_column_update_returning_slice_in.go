package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

type InsertOnConflictColumnUpdateReturningSliceQuery struct {
	Data       []*common.ConflictData `rel:"test_onconflict:k"`
	OnConflict struct {
		_ gosql.Column `sql:"key"`
		_ gosql.Update `sql:"fruit"`
	}
	_ gosql.Return `sql:"*"`
}
