package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

type UpdatePKeyReturningAllSingleQuery struct {
	User *common.User3 `rel:"test_user"`
	_    gosql.Return  `sql:"*"`
}
