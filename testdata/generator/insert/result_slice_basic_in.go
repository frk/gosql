package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

type InsertWithResultSliceBasicQuery struct {
	Users  []*common.User `rel:"test_user:u"`
	Result []*common.User
}
