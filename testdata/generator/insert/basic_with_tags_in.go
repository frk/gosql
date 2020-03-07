package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

type InsertBasicWithTagsQuery struct {
	User *common.User3 `rel:"test_user:u"`
}
