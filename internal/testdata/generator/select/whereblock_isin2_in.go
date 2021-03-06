package testdata

import (
	"time"

	"github.com/frk/gosql/internal/testdata/common"
)

type SelectWithWhereBlockInPredicate2Query struct {
	Users []*common.User `rel:"test_user:u"`
	Where struct {
		CreatedAt struct {
			After  time.Time `sql:"x"`
			Before time.Time `sql:"y"`
		} `sql:"u.created_at isbetween"`
		Or struct {
			IDs           []int     `sql:"u.id isin"`
			CreatedBefore time.Time `sql:"u.created_at <"`
		} `sql:">" bool:"or"`
	}
}
