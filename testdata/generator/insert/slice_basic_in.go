package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

type InsertSliceBasicQuery struct {
	Users []*common.User2 `rel:"test_user:u"`
}
