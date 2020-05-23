package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

type InsertReturningCollistSliceQuery struct {
	Users []*common.User `rel:"test_user:u"`
	_     gosql.Return   `sql:"u.id"`
}
