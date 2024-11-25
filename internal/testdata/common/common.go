package common

import (
	"database/sql/driver"
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
	Id        int       `sql:"id" json:"id"`
	Email     string    `sql:"email" json:"email"`
	FullName  string    `sql:"full_name" json:"fullName"`
	CreatedAt time.Time `sql:"created_at" json:"createdAt"`
}

func (u *User2) AfterScan() {
	// ...
}

// for testing "ro", "wo", and "xf" tags
type User3 struct {
	Id        int       `sql:"id,pk,ro"`
	Email     string    `sql:"email"`
	Password  []byte    `sql:"password,wo,xf"`
	CreatedAt time.Time `sql:"created_at"`
	UpdatedAt time.Time `sql:"updated_at"`
}

type User3Iterator interface {
	NextUser(*User3) error
}

// for testing defaults
type User4 struct {
	Id        int       `sql:"id,ro"`
	Email     string    `sql:"email"`
	FullName  string    `sql:"full_name"`
	IsActive  bool      `sql:"is_active"`
	CreatedAt time.Time `sql:"created_at"`
	UpdatedAt time.Time `sql:"updated_at"`
}

type User4Iterator interface {
	NextUser(*User4) error
}

// for testing json
type User5 struct {
	Id        int                    `sql:"id,ro"`
	Email     string                 `sql:"email"`
	FullName  string                 `sql:"full_name"`
	IsActive  bool                   `sql:"is_active"`
	Metadata1 map[string]interface{} `sql:"metadata1,json"`
	Metadata2 ArbitraryStruct        `sql:"metadata2,json"`
	CreatedAt time.Time              `sql:"created_at"`
	UpdatedAt time.Time              `sql:"updated_at"`
}

type User5Iterator interface {
	NextUser(*User5) error
}

// for testing filter with json tags
type User6 struct {
	Id        int             `sql:"id,ro" json:"id"`
	Email     string          `sql:"email" json:"email"`
	FullName  string          `sql:"full_name" json:"fullName"`
	IsActive  bool            `sql:"is_active" json:"isActive"`
	Metadata1 ArbitraryStruct `sql:"metadata1,json" json:"-"` // should be omitted from filter map
	Metadata2 ArbitraryStruct `sql:"metadata2,json" json:""`  // should be omitted from filter map
	CreatedAt time.Time       `sql:"created_at" json:"createdAt"`
	UpdatedAt time.Time       `sql:"updated_at" json:"updatedAt"`
}

type ArbitraryStruct struct {
	// ...
}

// for testing conflict

type ConflictData struct {
	Id    int     `sql:"id,ro"`
	Key   int     `sql:"key"`
	Name  string  `sql:"name"`
	Fruit string  `sql:"fruit"`
	Value float64 `sql:"value"`
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

////////////////////////////////////////////////////////////////////////////////
// for testing embedded struct fields with *tag* FCKeys

type Embedded struct {
	FOO  *EFoo `json:"fooField" sql:">foo_"`
	EFoo `sql:">foo2_"`
}

type EFoo struct {
	Bar EBar  `json:"barField" sql:">bar_"`
	Baz *EBaz `json:"bazField" sql:">baz_"`
}

type EBar struct {
	EBaz `sql:">baz_"`
}

type EBaz struct {
	Val string `json:"value" sql:"val"`
}

////////////////////////////////////////////////////////////////////////////////

type BadIterator interface { // unexported method
	fn(*User) error
}

type ErrorHandler struct{}

func (ErrorHandler) HandleError(err error) error { return err }

type ErrorInfoHandler struct{}

func (ErrorInfoHandler) HandleErrorInfo(info *gosql.ErrorInfo) error { return nil }

type MyTime struct {
	time.Time
}

func (t MyTime) Value() (driver.Value, error) {
	return t.Time, nil
}

func (t MyTime) Scan(src interface{}) error {
	// ....
	return nil
}
