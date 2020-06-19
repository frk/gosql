package testdata

import (
	"time"

	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

type DeleteWithReturningSliceErrorHandlerQuery struct {
	Users []*common.User `rel:"test_user:u"`
	Where struct {
		CreatedBefore time.Time `sql:"u.created_at <"`
	}
	_  gosql.Return `sql:"*"`
	eh common.ErrorHandler
}
