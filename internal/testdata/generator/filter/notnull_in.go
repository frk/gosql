package testdata

import (
	"github.com/frk/gosql/internal/testdata"
	"github.com/frk/gosql/internal/testdata/common"
)

type FilterNotNULLRecords struct {
	CT1                *testdata.CT1 `rel:"view_test:v"`
	common.FilterMaker `filter:"v2"`
}
