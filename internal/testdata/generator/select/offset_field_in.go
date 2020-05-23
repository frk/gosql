package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

type SelectWithOffsetFieldQuery struct {
	Users  []*common.User `rel:"test_user:u"`
	Offset int
}
