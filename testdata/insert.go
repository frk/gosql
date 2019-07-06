package testdata

import (
	"github.com/frk/gosql/testdata/common"
)

//BAD: missing relation field
type InsertTestBAD1 struct {
	// no record type ...
}

//BAD: missing `rel` tag
type InsertTestBAD2 struct {
	User *common.User
}

//BAD: invalid datatype kind
type InsertTestBAD3 struct {
	User string `rel:"users_table"`
}

//OK: user datatype
type InsertTestOK1 struct {
	UserRec *common.User `rel:"users_table"`
}

//OK: ignored datatype fields
type InsertTestOK2 struct {
	UserRec struct {
		_     string `sql:"name"` // ignore blank fields
		Name  string `sql:"-"`    // ignore "-" tags
		Name2 string ``           // ignore no `sql` tag
		Name3 string `sql:"name"` // all good
	} `rel:"users_table"`
}
