package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

type InsertOnConflictIgnoreSliceQuery struct {
	Data       []*common.ConflictData `rel:"test_onconflict:k"`
	OnConflict struct {
		_ gosql.Ignore
	}
}
