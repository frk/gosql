package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

type SelectWithLimitDirectiveQuery struct {
	Users []*common.User `rel:"test_user:u"`
	_     gosql.Limit    `sql:"25"`
}
