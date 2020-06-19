package testdata

import (
	"time"

	"github.com/frk/gosql/internal/testdata/common"
)

type SelectWithIteratorIfaceQuery struct {
	iter  common.UserIterator `rel:"test_user:u"`
	where struct {
		createdafter time.Time `sql:"u.created_at >"`
	}
	limit  int
	offset int
}
