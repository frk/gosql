package pg2go

import (
	"database/sql/driver"
	"encoding/json"
)

type JSON struct {
	Val interface{}
}

func (j JSON) Scan(src interface{}) error {
	if b, ok := src.([]byte); ok {
		return json.Unmarshal(b, j.Val)
	}
	return nil
}

func (j JSON) Value() (driver.Value, error) {
	if j.Val == nil {
		return nil, nil
	}
	return json.Marshal(j.Val)
}
