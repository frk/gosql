package testdata

import (
	"time"

	"github.com/frk/gosql/testdata/common"
)

type SelectWithOffsetFieldDefaultQuery struct {
	Users []*common.User `rel:"test_user:u"`
	Where struct {
		CreatedAfter time.Time `sql:"u.created_at >"`
	}
	Offset int `sql:"50"`
}
