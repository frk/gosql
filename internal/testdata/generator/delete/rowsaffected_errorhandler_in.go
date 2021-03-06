package testdata

import (
	"time"

	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

type DeleteWithRowsAffectedErrorHandlerQuery struct {
	_     gosql.Relation `rel:"test_user:u"`
	Where struct {
		CreatedBefore time.Time `sql:"u.created_at <"`
	}
	RowsAffected int
	common.ErrorHandler
}
