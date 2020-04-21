package pg2go

import (
	"database/sql/driver"
)

type BitFromBool struct {
	Val bool
}

func (v BitFromBool) Value() (driver.Value, error) {
	if v.Val {
		return []byte(`1`), nil
	}
	return []byte(`0`), nil
}
