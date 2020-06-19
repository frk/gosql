package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

type SelectWithFilterIteratorQuery struct {
	Iter common.User2Iterator `rel:"test_user:u"`
	gosql.Filter
}
