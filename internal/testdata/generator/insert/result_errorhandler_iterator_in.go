package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

type InsertResultErrorHandlerIteratorQuery struct {
	Users  []*common.User `rel:"test_user:u"`
	result common.User2Iterator
	erh    common.ErrorHandler
}
