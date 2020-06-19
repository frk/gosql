package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

type UpdateFilterSingleQuery struct {
	User *common.User4 `rel:"test_user:u"`
	gosql.Filter
}
