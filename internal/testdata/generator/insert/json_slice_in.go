package testdata

import (
	"github.com/frk/gosql/internal/testdata/common"
)

type InsertJSONSliceQuery struct {
	Users []*common.User5 `rel:"test_user:u"`
}
