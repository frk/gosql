package testdata

import (
	"time"

	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

type DeleteWithReturningIteratorAfterScanAllQuery struct {
	Iter  common.User2Iterator `rel:"test_user:u"`
	Where struct {
		CreatedBefore time.Time `sql:"u.created_at <"`
	}
	_ gosql.Return `sql:"*"`
}
