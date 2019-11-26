package testdata

import (
	"github.com/frk/gosql"
)

type DeleteWithFilterQuery struct {
	_ gosql.Relation `rel:"test_user"`
	gosql.Filter
}
