package testdata

import (
	"github.com/frk/gosql"
)

type SelectCountWithFilterQuery struct {
	Count int `rel:"test_user:u"`
	gosql.Filter
}
