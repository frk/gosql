package pgsql

import (
	"database/sql/driver"
)

// BitFromBool returns a driver.Valuer that produces a PostgreSQL bit from the given Go bool.
func BitFromBool(val bool) driver.Valuer {
	return bitFromBool{val: val}
}

type bitFromBool struct {
	val bool
}

func (v bitFromBool) Value() (driver.Value, error) {
	if v.val {
		return []byte(`1`), nil
	}
	return []byte(`0`), nil
}
