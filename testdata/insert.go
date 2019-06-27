package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

//BAD: missing record type field
type InsertTestBAD1 struct {
	// no record type ...
}

//BAD: missing `rel` tag
type InsertTestBAD2 struct {
	User *common.User
}

//BAD: invalid record type
type InsertTestBAD3 struct {
	User string `rel:"users_table"`
}

//OK: user record
type InsertTestOK1 struct {
	UserRec *common.User `rel:"users_table"`
}
