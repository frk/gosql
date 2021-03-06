package testdata

import (
	"github.com/frk/gosql/internal/testdata/common"
)

type SelectWithWhereBlockSingleQuery struct {
	User  *common.User `rel:"test_user:u"`
	Where struct {
		Id int `sql:"u.id"`
	}
}
