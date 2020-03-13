package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

type InsertJSONSingleQuery struct {
	User *common.User5 `rel:"test_user:u"`
}
