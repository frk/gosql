package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

type InsertWithReturningSingleAllQuery struct {
	User *common.User3 `rel:"test_user:u"`
	_    gosql.Return  `sql:"*"`
}
