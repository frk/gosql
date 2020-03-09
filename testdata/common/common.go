package common

import (
	"time"

	"github.com/frk/gosql"
)

type UserIterator interface {
	NextUser(*User) error
}

type User struct {
	Id        int       `sql:"id"`
	Email     string    `sql:"email"`
	FullName  string    `sql:"full_name"`
	CreatedAt time.Time `sql:"created_at"`
}

// for testing AfterScan

type User2Iterator interface {
	NextUser(*User2) error
}

type User2 struct {
	Id        int       `sql:"id"`
	Email     string    `sql:"email"`
	FullName  string    `sql:"full_name"`
	CreatedAt time.Time `sql:"created_at"`
}

func (u *User2) AfterScan() {
	// ...
}

// for testing ro and wo tags
type User3 struct {
	Id        int       `sql:"id,ro"`
	Email     string    `sql:"email"`
	Password  []byte    `sql:"password,wo"`
	CreatedAt time.Time `sql:"created_at"`
	UpdatedAt time.Time `sql:"updated_at"`
}

type User3Iterator interface {
	NextUser(*User3) error
}

// for testing nested struct fields

type Nested struct {
	FOO *Foo `sql:">foo_"`
	Foo `sql:">foo2_"`
}

type Foo struct {
	Bar Bar  `sql:">bar_"`
	Baz *Baz `sql:">baz_"`
}

type Bar struct {
	Baz `sql:">baz_"`
}

type Baz struct {
	Val string `sql:"val"`
}

type BadIterator interface { // unexported method
	fn(*User) error
}

type ErrorHandler struct{}

func (ErrorHandler) HandleError(err error) error { return err }

type ErrorInfoHandler struct{}

func (ErrorInfoHandler) HandleErrorInfo(info *gosql.ErrorInfo) error { return nil }
