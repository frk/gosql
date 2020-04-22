package pgsql

import (
	"database/sql"
	"time"
)

// TimetzToString returns an sql.Scanner that converts a PostgreSQL timetz into a Go string and sets it to val.
func TimetzToString(val *string) sql.Scanner {
	return timetzToString{val: val}
}

// TimetzToByteSlice returns an sql.Scanner that converts a PostgreSQL timetz into a Go []byte and sets it to val.
func TimetzToByteSlice(val *[]byte) sql.Scanner {
	return timetzToByteSlice{val: val}
}

type timetzToString struct {
	val *string
}

func (v timetzToString) Scan(src interface{}) error {
	switch data := src.(type) {
	case time.Time:
		*v.val = data.Format(timetzLayout)
	case []byte:
		*v.val = string(data)
	case string:
		*v.val = data
	}
	return nil
}

type timetzToByteSlice struct {
	val *[]byte
}

func (v timetzToByteSlice) Scan(src interface{}) error {
	switch data := src.(type) {
	case time.Time:
		*v.val = []byte(data.Format(timetzLayout))
	case []byte:
		*v.val = data
	case string:
		*v.val = []byte(data)
	}
	return nil
}
