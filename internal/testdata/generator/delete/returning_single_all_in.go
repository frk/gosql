package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

type DeleteWithReturningSingleAllQuery struct {
	User  *common.User `rel:"test_user:u"`
	Where struct {
		Id int `sql:"u.id"`
	}
	_ gosql.Return `sql:"*"`
}
