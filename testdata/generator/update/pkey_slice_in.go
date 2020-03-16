package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

type UpdatePKeySliceQuery struct {
	Users []*common.User3 `rel:"test_user:u"`
}
