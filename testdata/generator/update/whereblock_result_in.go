package testdata

import (
	"time"

	"github.com/frk/gosql/testdata/common"
)

type UpdateWhereblockResultQuery struct {
	User  *common.User4 `rel:"test_user:u"`
	Where struct {
		CreatedAfter  time.Time `sql:"u.created_at >"`
		CreatedBefore time.Time `sql:"u.created_at <"`
	}
	Result []*common.User4
}
