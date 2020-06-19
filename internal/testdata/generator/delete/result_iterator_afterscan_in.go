package testdata

import (
	"time"

	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

type DeleteWithResultIteratorAfterScanQuery struct {
	_     gosql.Relation `rel:"test_user:u"`
	Where struct {
		CreatedBefore time.Time `sql:"u.created_at <"`
	}
	Result common.User2Iterator
}
