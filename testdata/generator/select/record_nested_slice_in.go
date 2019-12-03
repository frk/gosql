package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

type SelectWithRecordNestedSliceQuery struct {
	Nesteds []*common.Nested `rel:"test_nested:n"`
}
