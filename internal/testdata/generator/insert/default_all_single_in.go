// DO NOT EDIT. This file was generated by "github.com/frk/gosql".

package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

type InsertDefaultAllSingleQuery struct {
	User *common.User4 `rel:"test_user_with_defaults:u"`
	_    gosql.Default `sql:"*"`
}
