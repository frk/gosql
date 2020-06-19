package testdata

import (
	"github.com/frk/gosql/internal/testdata/common"
)

type InsertResultJSONSliceQuery struct {
	Users  []*common.User5 `rel:"test_user:u"`
	Result []*common.User5
}
