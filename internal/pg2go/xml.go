package pg2go

import (
	"database/sql/driver"
	"encoding/xml"
)

type XML struct {
	Val interface{}
}

func (x XML) Scan(src interface{}) error {
	if b, ok := src.([]byte); ok {
		return xml.Unmarshal(b, x.Val)
	}
	return nil
}

func (x XML) Value() (driver.Value, error) {
	if x.Val == nil {
		return nil, nil
	}
	return xml.Marshal(x.Val)
}
