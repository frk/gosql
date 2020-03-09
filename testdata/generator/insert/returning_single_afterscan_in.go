package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

type InsertWithReturningSingleAfterScanQuery struct {
	User *common.User2 `rel:"test_user:u"`
	_    gosql.Return  `sql:"*"`
}
