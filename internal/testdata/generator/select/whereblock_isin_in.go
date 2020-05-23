package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

type SelectWithWhereBlockInPredicateQuery struct {
	Users []*common.User `rel:"test_user:u"`
	Where struct {
		IDs []int `sql:"u.id isin"`
	}
}
