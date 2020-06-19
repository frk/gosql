package testdata

import (
	"github.com/frk/gosql/internal/testdata/common"
)

type InsertResultAfterScanSingleQuery struct {
	User   *common.User2 `rel:"test_user:u"`
	result *common.User2
}
