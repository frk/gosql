package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

type UpdateFromblockBasicSingleQuery struct {
	User *common.User4 `rel:"test_user:u"`
	From struct {
		_ gosql.Relation `sql:"test_post:p"`
	}
	Where struct {
		_ gosql.Column `sql:"u.id=p.user_id"`
		_ gosql.Column `sql:"p.is_spam"`
	}
}
