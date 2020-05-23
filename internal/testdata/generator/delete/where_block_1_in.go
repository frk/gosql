package testdata

import (
	"github.com/frk/gosql"
)

type DeleteWithWhereBlock1Query struct {
	_     gosql.Relation `rel:"test_user"`
	Where struct {
		Id int `sql:"id"`
	}
}
