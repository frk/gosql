package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

type UpdateAllSingleQuery struct {
	User *common.User4 `rel:"test_user"`
	_    gosql.All
}
