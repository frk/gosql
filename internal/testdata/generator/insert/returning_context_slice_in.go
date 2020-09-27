package testdata

import (
	"context"

	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

type InsertReturningContextSliceQuery struct {
	Users []*common.User `rel:"test_user:u"`
	_     gosql.Return   `sql:"*"`
	ctx   context.Context
}
