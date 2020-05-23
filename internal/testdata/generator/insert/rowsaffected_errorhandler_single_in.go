package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

type InsertRowsAffectedErrorHandlerSingleQuery struct {
	User         *common.User3 `rel:"test_user:u"`
	RowsAffected int
	common.ErrorHandler
}
