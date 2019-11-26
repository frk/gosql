package testdata

import (
	"github.com/frk/gosql"
)

type DeleteWithAllDirectiveQuery struct {
	_ gosql.Relation `rel:"test_user"`
	_ gosql.All
}
