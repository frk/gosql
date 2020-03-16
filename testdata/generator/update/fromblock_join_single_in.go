package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

type UpdateFromblockJoinSingleQuery struct {
	User *common.User4 `rel:"test_user:u"`
	From struct {
		_ gosql.Relation  `sql:"test_post:p"`
		_ gosql.LeftJoin  `sql:"test_join1:j1,j1.post_id = p.id"`
		_ gosql.RightJoin `sql:"test_join2:j2,j2.join1_id = j1.id"`
		_ gosql.FullJoin  `sql:"test_join3:j3,j3.join2_id = j2.id"`
		_ gosql.CrossJoin `sql:"test_join4:j4"`
	}
	Where struct {
		_ gosql.Column `sql:"u.id=p.user_id"`
		_ gosql.Column `sql:"p.is_spam"`
	}
}
