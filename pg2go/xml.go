package pg2go

import (
	"database/sql"
	"database/sql/driver"
	"encoding/xml"
)

// XML returns a value that implements both the driver.Valuer and sql.Scanner
// interfaces. The driver.Valuer produces a PostgreSQL xml from the given val
// and the sql.Scanner unmarshals a PostgreSQL xml into the given val.
func XML(val interface{}) interface {
	driver.Valuer
	sql.Scanner
} {
	return xmltype{val: val}
}

type xmltype struct {
	val interface{}
}

func (x xmltype) Scan(src interface{}) error {
	if b, ok := src.([]byte); ok {
		return xml.Unmarshal(b, x.val)
	}
	return nil
}

func (x xmltype) Value() (driver.Value, error) {
	if x.val == nil {
		return nil, nil
	}
	return xml.Marshal(x.val)
}
