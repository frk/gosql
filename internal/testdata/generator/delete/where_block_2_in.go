package testdata

import (
	"time"

	"github.com/frk/gosql"
)

type DeleteWithWhereBlock2Query struct {
	_     gosql.Relation `rel:"test_user"`
	Where struct {
		CreatedAfter  time.Time `sql:"created_at >"`
		CreatedBefore time.Time `sql:"created_at <"`
		FullName      struct {
			_ gosql.Column `sql:"full_name = 'John Doe'"`
			_ gosql.Column `sql:"full_name = 'Jane Doe'" bool:"or"`
		} `sql:">"`
	}
}
