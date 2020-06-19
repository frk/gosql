package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

type FilterNestedRecords struct {
	_ *common.Nested `rel:"test_nested:n"`
	gosql.Filter
}
