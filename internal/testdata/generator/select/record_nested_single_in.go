package testdata

import (
	"github.com/frk/gosql/internal/testdata/common"
)

type SelectWithRecordNestedSingleQuery struct {
	Nested *common.Nested `rel:"test_nested:n"`
}
