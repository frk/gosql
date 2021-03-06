package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

type UpdateWhereblockReturningAllSingleQuery struct {
	User  *common.User4 `rel:"test_user:u"`
	Where struct {
		Id int `sql:"u.id"`
	}
	_ gosql.Return `sql:"*"`
}
