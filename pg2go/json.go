package pg2go

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// JSON returns a value that implements both the driver.Valuer and sql.Scanner
// interfaces. The driver.Valuer produces a PostgreSQL json(b) from the given val
// and the sql.Scanner unmarshals a PostgreSQL json(b) into the given val.
func JSON(val interface{}) interface {
	driver.Valuer
	sql.Scanner
} {
	return jsontype{val: val}
}

type jsontype struct {
	val interface{}
}

func (j jsontype) Value() (driver.Value, error) {
	if j.val == nil {
		return nil, nil
	}
	return json.Marshal(j.val)
}

func (j jsontype) Scan(src interface{}) error {
	if b, ok := src.([]byte); ok {
		return json.Unmarshal(b, j.val)
	}
	return nil
}
