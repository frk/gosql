package testdata

import (
	"time"

	"github.com/frk/gosql/testdata/common"
)

type SelectWithWhereBlockArrayComparisonPredicate3Query struct {
	Users []*common.User `rel:"test_user:u"`
	Where struct {
		IDs       []int `sql:"u.id = any"`
		CreatedAt struct {
			After  time.Time `sql:"x"`
			Before time.Time `sql:"y"`
		} `sql:"u.created_at isbetween"`
		Or struct {
			IDs           []int     `sql:"u.id = any"`
			CreatedBefore time.Time `sql:"u.created_at <"`
		} `sql:">" bool:"or"`
		Emails    []string `sql:"u.email = some" bool:"or"`
		FullNames []string `sql:"u.full_name = some"`
		IsActive  bool     `sql:"u.is_active"`
	}
}
