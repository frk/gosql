package testdata

import (
	"github.com/frk/gosql/internal/testdata/common"
)

type DeleteWithDatatype2Query struct {
	_     *common.User `rel:"test_user:u"`
	Where struct {
		Id int `sql:"u.id"`
	}
}
