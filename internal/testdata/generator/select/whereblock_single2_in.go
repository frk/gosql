package testdata

import (
	"github.com/frk/gosql/internal/testdata"
	"github.com/frk/gosql/internal/testdata/common"
)

type SelectWithWhereBlockSingle2Query struct {
	Rel   *testdata.CT3b `rel:"column_tests_3"`
	Where struct {
		SomeTime common.MyTime `sql:"some_time <"`
	}
}
