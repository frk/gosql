package testdata

import (
	"context"

	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

type InsertReturningContextSingleQuery struct {
	User *common.User3 `rel:"test_user:u"`
	_    gosql.Return  `sql:"*"`
	ctx  context.Context
}
