package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

type InsertWithRowsAffectedQuery struct {
	User         *common.User3 `rel:"test_user:u"`
	RowsAffected int
}
