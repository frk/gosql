package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

type InsertOnConflictIndexIgnore1SingleQuery struct {
	Data       *common.ConflictData `rel:"test_onconflict:k"`
	OnConflict struct {
		_ gosql.Index `sql:"test_onconflict_name_fruit_idx"`
		_ gosql.Ignore
	}
}
