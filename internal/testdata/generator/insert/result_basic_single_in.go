package testdata

import (
	"github.com/frk/gosql/internal/testdata/common"
)

type InsertResultBasicSingleQuery struct {
	User   *common.User3 `rel:"test_user:u"`
	Result struct {
		Id int `sql:"id"`
	}
}
