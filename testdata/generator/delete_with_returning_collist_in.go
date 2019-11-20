package testdata

import (
	"github.com/frk/gosql"
)

type DeleteWithReturningCollistQuery struct {
	User struct {
		Email    string `sql:"email"`
		FullName string `sql:"full_name"`
	} `rel:"test_user:u"`
	Where struct {
		Id int `sql:"u.id"`
	}
	_ gosql.Return `sql:"u.email,u.full_name"`
}
