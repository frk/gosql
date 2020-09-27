package testdata

import (
	"context"
	"time"

	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

type DeleteWithReturningSliceContextQuery struct {
	context.Context
	Users []*common.User `rel:"test_user:u"`
	Where struct {
		CreatedBefore time.Time `sql:"u.created_at <"`
	}
	_ gosql.Return `sql:"*"`
}
