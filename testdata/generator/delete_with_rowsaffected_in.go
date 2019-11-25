package testdata

import (
	"time"

	"github.com/frk/gosql"
)

type DeleteWithRowsAffectedQuery struct {
	_     gosql.Relation `rel:"test_user:u"`
	Where struct {
		CreatedBefore time.Time `sql:"u.created_at <"`
	}
	RowsAffected int
}
