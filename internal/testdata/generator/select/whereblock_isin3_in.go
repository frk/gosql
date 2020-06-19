package testdata

import (
	"time"

	"github.com/frk/gosql/internal/testdata/common"
)

type SelectWithWhereBlockInPredicate3Query struct {
	Users []*common.User `rel:"test_user:u"`
	Where struct {
		IDs       []int `sql:"u.id isin"`
		CreatedAt struct {
			After  time.Time `sql:"x"`
			Before time.Time `sql:"y"`
		} `sql:"u.created_at isbetween"`
		Or struct {
			IDs           []int     `sql:"u.id isin"`
			CreatedBefore time.Time `sql:"u.created_at <"`
		} `sql:">" bool:"or"`
		Emails    []string `sql:"u.email isin" bool:"or"`
		FullNames []string `sql:"u.full_name isin"`
		IsActive  bool     `sql:"u.is_active"`
	}
}
