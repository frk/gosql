package testdata

import (
	"github.com/frk/gosql/internal/testdata/common"
)

type InsertResultBasicSliceQuery struct {
	Users  []*common.User `rel:"test_user:u"`
	Result []*common.User
}
