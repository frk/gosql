// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

type InsertDefaultAllSliceQuery struct {
	Users []*common.User4 `rel:"test_user_with_defaults:u"`
	_     gosql.Default   `sql:"*"`
}
