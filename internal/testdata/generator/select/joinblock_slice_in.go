package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

type SelectWithJoinBlockSliceQuery struct {
	Users []*common.User `rel:"test_user:u"`
	Join  struct {
		_ gosql.LeftJoin  `sql:"test_post:p,p.user_id = u.id"`
		_ gosql.LeftJoin  `sql:"test_join1:j1,j1.post_id = p.id"`
		_ gosql.RightJoin `sql:"test_join2:j2,j2.join1_id = j1.id"`
		_ gosql.FullJoin  `sql:"test_join3:j3,j3.join2_id = j2.id"`
		_ gosql.CrossJoin `sql:"test_join4:j4"`
	}
	Where struct {
		_ gosql.Column `sql:"p.is_spam"`
	}
}
