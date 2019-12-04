package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

type SelectWithWhereBlockModifierFuncSingleQuery struct {
	User  *common.User `rel:"test_user:u"`
	Where struct {
		Email string `sql:"u.email,@lower"`
	}
}
