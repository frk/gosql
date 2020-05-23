package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

type SelectWithAfterScanSingleQuery struct {
	User *common.User2 `rel:"test_user:u"`
}
