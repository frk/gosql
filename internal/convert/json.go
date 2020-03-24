package convert

import (
	"database/sql/driver"
	"encoding/json"
)

type JSON struct {
	V interface{}
}

func (j JSON) Scan(src interface{}) error {
	if b, ok := src.([]byte); ok {
		return json.Unmarshal(b, j.V)
	}
	return nil
}

func (j JSON) Value() (driver.Value, error) {
	if j.V == nil {
		return nil, nil
	}
	return json.Marshal(j.V)
}
