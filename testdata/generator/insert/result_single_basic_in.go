package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

type InsertWithResultSingleQuery struct {
	User   *common.User3 `rel:"test_user:u"`
	Result struct {
		Id int `sql:"id"`
	}
}
