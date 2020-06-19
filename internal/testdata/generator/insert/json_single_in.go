package testdata

import (
	"github.com/frk/gosql/internal/testdata/common"
)

type InsertJSONSingleQuery struct {
	User *common.User5 `rel:"test_user:u"`
}
