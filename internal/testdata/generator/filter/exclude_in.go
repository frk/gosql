package testdata

import (
	"github.com/frk/gosql/internal/testdata/common"
)

type FilterWithExcludeRecords struct {
	User   *common.User3 `rel:"test_user"`
	Filter common.FilterMaker
}
