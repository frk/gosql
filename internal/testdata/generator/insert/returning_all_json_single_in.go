package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

type InsertReturningAllJSONSingleQuery struct {
	User *common.User5 `rel:"test_user:u"`
	_    gosql.Return  `sql:"*"`
}
