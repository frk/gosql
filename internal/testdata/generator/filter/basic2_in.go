package testdata

import (
	"github.com/frk/gosql/internal/testdata/common"
)

type FilterBasic2Records struct {
	User   *common.User6 `rel:"test_user"`
	Filter common.FilterMaker
}
