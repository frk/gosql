package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

type InsertReturningAfterScanSliceQuery struct {
	Users []*common.User2 `rel:"test_user:u"`
	_     gosql.Return    `sql:"*"`
}
