package testdata

import (
	"github.com/frk/gosql/internal/testdata/common"
)

type UpdatePKeySliceQuery struct {
	Users []*common.User3 `rel:"test_user:u"`
}
