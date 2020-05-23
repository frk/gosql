package testdata

import (
	"time"

	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

type DeleteWithResultIteratorErrorInfoHandlerQuery struct {
	_     gosql.Relation `rel:"test_user:u"`
	Where struct {
		CreatedBefore time.Time `sql:"u.created_at <"`
	}
	Result     common.User2Iterator
	errhandler common.ErrorInfoHandler
}
