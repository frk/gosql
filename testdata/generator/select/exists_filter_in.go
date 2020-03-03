package testdata

import (
	"github.com/frk/gosql"
)

type SelectExistsWithFilterQuery struct {
	Exists bool `rel:"test_user:u"`
	gosql.Filter
}
