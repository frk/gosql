package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

type InsertOnConflictIndexUpdateSliceQuery struct {
	Data       []*common.ConflictData `rel:"test_onconflict:k"`
	OnConflict struct {
		_ gosql.Index  `sql:"test_onconflict_fruit_key_name_idx"`
		_ gosql.Update `sql:"*"`
	}
	_ gosql.Return `sql:"*"`
}
