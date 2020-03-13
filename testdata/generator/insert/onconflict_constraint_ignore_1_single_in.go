package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

type InsertOnConflictConstraintIgnore1SingleQuery struct {
	Data       *common.ConflictData `rel:"test_onconflict:k"`
	OnConflict struct {
		_ gosql.Constraint `sql:"test_onconflict_key_value_key"`
		_ gosql.Ignore
	}
}
