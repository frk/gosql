package testdata

import (
	"github.com/frk/gosql/internal/testdata/common"
)

type InsertResultErrorHandlerSingleQuery struct {
	User   *common.User `rel:"test_user:u"`
	result *common.User
	erh    common.ErrorHandler
}
