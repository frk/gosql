package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

type SelectWithWhereBlockQuery struct {
	User  *common.User `rel:"test_user"`
	Where struct {
		Id int `sql:"id"`
	}
}
