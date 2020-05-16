package pgsql

import (
	"database/sql"
)

// AnyToEmptyInterface returns an sql.Scanner that sets val from
// the value provided to the Scan method.
func AnyToEmptyInterface(val *interface{}) sql.Scanner {
	return anyToEmptyInterface{val: val}
}

type anyToEmptyInterface struct {
	val *interface{}
}

func (v anyToEmptyInterface) Scan(src interface{}) error {
	*v.val = src
	return nil
}
