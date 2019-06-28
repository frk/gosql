package testdata

import (
	"database/sql/driver"
	"time"
)

// Var is not error type, should return false
type IsErrorTest1 struct {
	Var string
}

// Var is not error type, should return false
type IsErrorTest2 struct {
	Var []error
}

// Var is error type, should return true
type IsErrorTest3 struct {
	Var error
}

////////////////////////////////////////////////////////////////////////////////

// Var is not interface{} type, should return false
type IsEmptyInterfaceTest1 struct {
	Var string
}

// Var is not interface{} type, should return false
type IsEmptyInterfaceTest2 struct {
	Var []interface{}
}

// Var is not interface{} type, should return false
type IsEmptyInterfaceTest3 struct {
	Var interface {
		M()
	}
}

// Var is interface{} type, should return true
type IsEmptyInterfaceTest4 struct {
	Var interface{}
}

////////////////////////////////////////////////////////////////////////////////

// Var is not time.Time type, should return false
type IsTimeTest1 struct {
	Var string
}

// Var is not time.Time type, should return false
type IsTimeTest2 struct {
	Var []time.Time
}

// Var is time.Time type, should return true
type IsTimeTest3 struct {
	Var time.Time
}

type CustomTime struct {
	time.Timer
	time.Time
}

// Var is a custom type that embeds time.Time type *directly*, should return true
type IsTimeTest4 struct {
	Var CustomTime
}

type CustomTime2 struct {
	CustomTime
}

// Var is a custom type that embeds time.Time type *indirectly*, should return false
type IsTimeTest5 struct {
	Var CustomTime2
}

////////////////////////////////////////////////////////////////////////////////

// Var is not driver.Value type, should return false
type IsSqlDriverValueTest1 struct {
	Var string
}

// Var is not driver.Value type, should return false
type IsSqlDriverValueTest2 struct {
	Var []driver.Value
}

// Var is driver.Value type, should return true
type IsSqlDriverValueTest3 struct {
	Var driver.Value
}

////////////////////////////////////////////////////////////////////////////////

// Does not implement the "database/sql.Scanner" interface, should return false
type ImplementsScannerTest1 struct{}

// Does not implement the "database/sql.Scanner" interface, should return false
type ImplementsScannerTest2 struct{}

func (ImplementsScannerTest2) Scan() {}

// Does not implement the "database/sql.Scanner" interface, should return false
type ImplementsScannerTest3 struct{}

func (ImplementsScannerTest3) Scan(src interface{}) (error, error) { return nil, nil }

// Does not implement the "database/sql.Scanner" interface, should return false
type ImplementsScannerTest4 struct{}

func (ImplementsScannerTest4) Scan(src ...interface{}) error { return nil }

// Does implement the "database/sql.Scanner" interface, should return true
type ImplementsScannerTest5 struct{}

func (ImplementsScannerTest5) Scan(src interface{}) error { return nil }

////////////////////////////////////////////////////////////////////////////////

// Does not implement the "database/sql/driver.Valuer" interface, should return false
type ImplementsValuerTest1 struct{}

// Does not implement the "database/sql/driver.Valuer" interface, should return false
type ImplementsValuerTest2 struct{}

func (ImplementsValuerTest2) Value() {}

// Does not implement the "database/sql/driver.Valuer" interface, should return false
type ImplementsValuerTest3 struct{}

func (ImplementsValuerTest3) Value(x interface{}) (driver.Value, error) { return nil, nil }

// Does not implement the "database/sql/driver.Valuer" interface, should return false
type ImplementsValuerTest4 struct{}

func (ImplementsValuerTest4) Value() driver.Value { return nil }

// Does not implement the "database/sql/driver.Valuer" interface, should return false
type ImplementsValuerTest5 struct{}

func (ImplementsValuerTest5) Value() error { return nil }

// Does not implement the "database/sql/driver.Valuer" interface, should return false
type ImplementsValuerTest6 struct{}

func (ImplementsValuerTest6) Value() (error, driver.Value) { return nil, nil }

// Does implement the "database/sql/driver.Valuer" interface, should return true
type ImplementsValuerTest7 struct{}

func (ImplementsValuerTest7) Value() (driver.Value, error) { return nil, nil }
