package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

type InsertReturningErrorHandlerSliceQuery struct {
	Users []*common.User3 `rel:"test_user:u"`
	_     gosql.Return    `sql:"u.id"`
	erh   common.ErrorHandler
}
