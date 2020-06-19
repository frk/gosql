package testdata

import (
	"github.com/frk/gosql/internal/testdata/common"
)

type UpdateWhereblockBasicSingle1Query struct {
	User  *common.User4 `rel:"test_user"`
	Where struct {
		Id int `sql:"id"`
	}
}
