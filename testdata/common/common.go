package common

import (
	"time"
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

// for testing nested struct fields
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
