package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

type UpdateFilterResultSliceQuery struct {
	User *common.User4 `rel:"test_user:u"`
	gosql.Filter
	Result []*common.User4
}
