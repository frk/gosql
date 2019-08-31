package testdata

import (
	"database/sql/driver"
	"encoding/xml"
	"time"

	"github.com/frk/gosql"
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

// Var is not gosql directive type, should return false
type IsDirectiveTest1 struct {
	Var string
}

// Var is not gosql directive type, should return false
type IsDirectiveTest2 struct {
	Var []gosql.Column
}

type Column struct{}

// Var is not gosql directive type, should return false
type IsDirectiveTest3 struct {
	Var Column
}

type Relation struct {
	gosql.Relation
}

// Var is not gosql directive type, should return false
type IsDirectiveTest4 struct {
	Var Relation
}

// Var is gosql directive type, should return true
type IsDirectiveTest5 struct {
	Var gosql.Column
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

////////////////////////////////////////////////////////////////////////////////

// Does not implement the "gosql.AfterScanner" interface, should return false
type ImplementsAfterScannerTest1 struct{}

// Does not implement the "gosql.AfterScanner" interface, should return false
type ImplementsAfterScannerTest2 struct{}

func (ImplementsAfterScannerTest2) AfterScan() error { return nil }

// Does not implement the "gosql.AfterScanner" interface, should return false
type ImplementsAfterScannerTest3 struct{}

func (ImplementsAfterScannerTest3) AfterScan(x interface{}) {}

// Does implement the "gosql.AfterScanner" interface, should return true
type ImplementsAfterScannerTest4 struct{}

func (ImplementsAfterScannerTest4) AfterScan() {}

////////////////////////////////////////////////////////////////////////////////

// Does not implement the "gosql.ErrorHandler" interface, should return false
type ImplementsErrorHandlerTest1 struct{}

// Does not implement the "gosql.ErrorHandler" interface, should return false
type ImplementsErrorHandlerTest2 struct{}

func (ImplementsErrorHandlerTest2) HandleError() error { return nil }

// Does not implement the "gosql.ErrorHandler" interface, should return false
type ImplementsErrorHandlerTest3 struct{}

func (ImplementsErrorHandlerTest3) HandleError(err error) {}

// Does not implement the "gosql.ErrorHandler" interface, should return false
type ImplementsErrorHandlerTest4 struct{}

func (ImplementsErrorHandlerTest4) HandleError(err interface{}) error { return nil }

// Does not implement the "gosql.ErrorHandler" interface, should return false
type ImplementsErrorHandlerTest5 struct{}

func (ImplementsErrorHandlerTest5) HandleError(err error) interface{} { return nil }

// Does implement the "gosql.ErrorHandler" interface, should return true
type ImplementsErrorHandlerTest6 struct{}

func (ImplementsErrorHandlerTest6) HandleError(err error) error { return nil }

////////////////////////////////////////////////////////////////////////////////

// Does not implement the "json.Marshaler" interface, should return false
type ImplementsJSONMarshalerTest1 struct{}

// Does not implement the "json.Marshaler" interface, should return false
type ImplementsJSONMarshalerTest2 struct{}

func (ImplementsJSONMarshalerTest2) UnmarshalJSON([]byte) error { return nil }

// Does implement the "json.Marshaler" interface, should return true
type ImplementsJSONMarshalerTest3 struct{}

func (ImplementsJSONMarshalerTest3) MarshalJSON() ([]byte, error) { return nil, nil }

////////////////////////////////////////////////////////////////////////////////

// Does not implement the "json.Unmarshaler" interface, should return false
type ImplementsJSONUnmarshalerTest1 struct{}

// Does not implement the "json.Unmarshaler" interface, should return false
type ImplementsJSONUnmarshalerTest2 struct{}

func (ImplementsJSONUnmarshalerTest2) MarshalJSON() ([]byte, error) { return nil, nil }

// Does implement the "json.Unmarshaler" interface, should return true
type ImplementsJSONUnmarshalerTest3 struct{}

func (ImplementsJSONUnmarshalerTest3) UnmarshalJSON([]byte) error { return nil }

////////////////////////////////////////////////////////////////////////////////

// Does not implement the "xml.Marshaler" interface, should return false
type ImplementsXMLMarshalerTest1 struct{}

// Does not implement the "xml.Marshaler" interface, should return false
type ImplementsXMLMarshalerTest2 struct{}

func (ImplementsXMLMarshalerTest2) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return nil
}

// Does implement the "xml.Marshaler" interface, should return true
type ImplementsXMLMarshalerTest3 struct{}

func (ImplementsXMLMarshalerTest3) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return nil
}

////////////////////////////////////////////////////////////////////////////////

// Does not implement the "xml.Unmarshaler" interface, should return false
type ImplementsXMLUnmarshalerTest1 struct{}

// Does not implement the "xml.Unmarshaler" interface, should return false
type ImplementsXMLUnmarshalerTest2 struct{}

func (ImplementsXMLUnmarshalerTest2) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return nil
}

// Does implement the "xml.Unmarshaler" interface, should return true
type ImplementsXMLUnmarshalerTest3 struct{}

func (ImplementsXMLUnmarshalerTest3) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return nil
}
