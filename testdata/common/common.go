package common

import (
	"time"
)

type User struct {
	Id        int       `sql:"id"`
	Email     string    `sql:"email"`
	FullName  string    `sql:"full_name"`
	CreatedAt time.Time `sql:"created_at"`
}

// for testing nested struct fields
type Foo struct {
	Bar Bar `sql:">bar_"`
}

type Bar struct {
	Baz `sql:">baz_"`
}

type Baz struct {
	Val string `sql:"val"`
}
