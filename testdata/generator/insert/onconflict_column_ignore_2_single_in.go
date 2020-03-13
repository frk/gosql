package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

type InsertOnConflictColumnIgnore2SingleQuery struct {
	Data       *common.ConflictData `rel:"test_onconflict:k"`
	OnConflict struct {
		_ gosql.Column `sql:"key,name"`
		_ gosql.Ignore
	}
}
