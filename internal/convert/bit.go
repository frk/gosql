package convert

import (
	"database/sql/driver"
)

type BitFromBool struct {
	V bool
}

func (v BitFromBool) Value() (driver.Value, error) {
	if v.V {
		return []byte(`1`), nil
	}
	return []byte(`0`), nil
}
