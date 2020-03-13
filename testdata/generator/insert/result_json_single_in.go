package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

type InsertResultJSONSingleQuery struct {
	User   *common.User5 `rel:"test_user:u"`
	Result *common.User5
}
