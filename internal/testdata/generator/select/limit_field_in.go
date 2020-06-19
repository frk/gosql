package testdata

import (
	"github.com/frk/gosql/internal/testdata/common"
)

type SelectWithLimitFieldQuery struct {
	Users []*common.User `rel:"test_user:u"`
	Limit int
}
