package testdata

import (
	"time"

	"github.com/frk/gosql/internal/testdata/common"
)

type SelectWithIteratorFuncErrorHandlerQuery struct {
	next  func(*common.User) error `rel:"test_user:u"`
	where struct {
		createdafter time.Time `sql:"u.created_at >"`
	}
	limit  int
	offset int
	erh    common.ErrorHandler
}
