package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

type FilterBasicRecords struct {
	User *common.User2 `rel:"test_user"`
	gosql.Filter
}
