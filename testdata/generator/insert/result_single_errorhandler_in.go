package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

type InsertWithResultSingleErrorHandlerQuery struct {
	User   *common.User `rel:"test_user:u"`
	result *common.User
	erh    common.ErrorHandler
}
