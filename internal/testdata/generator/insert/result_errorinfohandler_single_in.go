package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

type InsertResultErrorInfoHandlerSingleQuery struct {
	User   *common.User2 `rel:"test_user:u"`
	Result *common.User2
	erh    common.ErrorInfoHandler
}
