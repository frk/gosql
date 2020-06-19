package testdata

import (
	"time"

	"github.com/frk/gosql/internal/testdata/common"
)

type SelectWithWhereBlockBetweenQuery struct {
	Users []*common.User `rel:"test_user:u"`
	Where struct {
		CreatedAt struct {
			After  time.Time `sql:"x"`
			Before time.Time `sql:"y"`
		} `sql:"u.created_at isbetween"`
	}
}
