package testdata

import (
	"github.com/frk/gosql/internal/testdata/common"
)

type InsertResultAfterScanSliceQuery struct {
	Users  []*common.User `rel:"test_user:u"`
	Result []*common.User2
}
