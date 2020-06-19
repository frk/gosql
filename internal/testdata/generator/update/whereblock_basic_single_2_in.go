package testdata

import (
	"time"

	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

type UpdateWhereblockBasicSingle2Query struct {
	User  *common.User4 `rel:"test_user:u"`
	Where struct {
		CreatedAfter  time.Time `sql:"u.created_at >"`
		CreatedBefore time.Time `sql:"u.created_at <"`
		FullName      struct {
			_ gosql.Column `sql:"u.full_name = 'John Doe'"`
			_ gosql.Column `sql:"u.full_name = 'Jane Doe'" bool:"or"`
		} `sql:">"`
	}
}
