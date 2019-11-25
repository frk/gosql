package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

type DeleteWithResultAfterScanQuery struct {
	_     gosql.Relation `rel:"test_user:u"`
	Where struct {
		Id int `sql:"u.id"`
	}
	Result *common.User2
}
