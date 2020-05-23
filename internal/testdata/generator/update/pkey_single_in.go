package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

type UpdatePKeySingleQuery struct {
	User *common.User3 `rel:"test_user"`
}
