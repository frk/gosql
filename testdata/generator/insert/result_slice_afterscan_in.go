package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

type InsertWithResultSliceAfterScanQuery struct {
	Users  []*common.User `rel:"test_user:u"`
	Result []*common.User2
}
