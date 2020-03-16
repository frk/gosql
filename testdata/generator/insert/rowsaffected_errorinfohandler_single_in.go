package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

type InsertRowsAffectedErrorInfoHandlerSingleQuery struct {
	User         *common.User3 `rel:"test_user:u"`
	RowsAffected int
	common.ErrorInfoHandler
}