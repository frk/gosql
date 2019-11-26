package testdata

import (
	"github.com/frk/gosql"
)

type DeleteWithUsingJoinBlock1Query struct {
	_     gosql.Relation `rel:"test_user:u"`
	Using struct {
		_ gosql.Relation `sql:"test_post:p"`
	}
	Where struct {
		_ gosql.Column `sql:"u.id=p.user_id"`
		_ gosql.Column `sql:"p.is_spam"`
	}
}
