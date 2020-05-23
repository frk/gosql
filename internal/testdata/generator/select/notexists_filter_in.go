package testdata

import (
	"github.com/frk/gosql"
)

type SelectNotExistsWithFilterQuery struct {
	NotExists bool `rel:"test_user:u"`
	gosql.Filter
}
