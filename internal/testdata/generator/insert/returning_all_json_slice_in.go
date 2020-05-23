package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

type InsertReturningAllJSONSliceQuery struct {
	Users []*common.User5 `rel:"test_user:u"`
	_     gosql.Return    `sql:"*"`
}
