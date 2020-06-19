package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

type InsertReturningCollistSingleQuery struct {
	User *common.User3 `rel:"test_user:u"`
	_    gosql.Return  `sql:"u.id,u.created_at"`
}
