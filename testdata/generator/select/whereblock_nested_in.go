package testdata

import (
	"time"

	"github.com/frk/gosql/testdata/common"
)

type SelectWithWhereBlockNestedQuery struct {
	Users []*common.User `rel:"test_user:u"`
	Where struct {
		FullName     string    `sql:"u.full_name islike"`
		CreatedAfter time.Time `sql:"u.created_at >"`
		Or           struct {
			FullName      string    `sql:"u.full_name islike"`
			CreatedBefore time.Time `sql:"u.created_at <"`
		} `sql:">" bool:"or"`
	}
}
