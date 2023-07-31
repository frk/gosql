package testdata

import (
	"github.com/frk/gosql/internal/testdata/common"
)

type FilterEmbeddedRecords struct {
	_ *common.Embedded `rel:"test_nested:n"`
	common.FilterMaker
}
