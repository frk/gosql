package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

type InsertWithReturningSliceErrorHandlerQuery struct {
	Users []*common.User3 `rel:"test_user:u"`
	_     gosql.Return    `sql:"u.id"`
	erh   common.ErrorHandler
}
