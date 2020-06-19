package testdata

import (
	"github.com/frk/gosql/internal/testdata/common"
)

type InsertRowsAffectedSingleQuery struct {
	User         *common.User3 `rel:"test_user:u"`
	RowsAffected int
}
