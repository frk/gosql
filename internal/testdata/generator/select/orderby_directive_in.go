package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

type SelectWithOrderByDirectiveQuery struct {
	Users []*common.User `rel:"test_user:u"`
	_     gosql.OrderBy  `sql:"-u.created_at"`
}
