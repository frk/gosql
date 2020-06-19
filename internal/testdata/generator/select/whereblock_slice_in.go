package testdata

import (
	"time"

	"github.com/frk/gosql/internal/testdata/common"
)

type SelectWithWhereBlockSliceQuery struct {
	Users []*common.User `rel:"test_user:u"`
	Where struct {
		CreatedAfter time.Time `sql:"u.created_at >"`
	}
}
