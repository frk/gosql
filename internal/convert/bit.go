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

// DO NOT USE: bit can be scanned into `bool` directly.
type BitToBool NO_TYPE

// DO NOT USE: bit can be scanned into `byte` directly.
type BitToByte NO_TYPE

// DO NOT USE: bit can be scanned into `string` directly.
type BitToString NO_TYPE

// DO NOT USE: bit can be created from `byte` directly.
type BitFromByte NO_TYPE

// DO NOT USE: bit can be created from `string` directly.
type BitFromString NO_TYPE
