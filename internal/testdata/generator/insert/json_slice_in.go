package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

type InsertJSONSliceQuery struct {
	Users []*common.User5 `rel:"test_user:u"`
}
