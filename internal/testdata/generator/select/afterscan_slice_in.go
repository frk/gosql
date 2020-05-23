package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

type SelectWithAfterScanSliceQuery struct {
	Users []*common.User2 `rel:"test_user:u"`
}
