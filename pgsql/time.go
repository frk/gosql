package pgsql

import (
	"database/sql"
	"time"
)

// TimeToString returns an sql.Scanner that converts a PostgreSQL time into a Go string and sets it to val.
func TimeToString(val *string) sql.Scanner {
	return timeToString{val: val}
}

// TimeToByteSlice returns an sql.Scanner that converts a PostgreSQL time into a Go []byte and sets it to val.
func TimeToByteSlice(val *[]byte) sql.Scanner {
	return timeToByteSlice{val: val}
}

type timeToString struct {
	val *string
}

func (v timeToString) Scan(src interface{}) error {
	switch data := src.(type) {
	case time.Time:
		*v.val = data.Format(timeLayout)
	case []byte:
		*v.val = string(data)
	case string:
		*v.val = data
	}
	return nil
}

type timeToByteSlice struct {
	val *[]byte
}

func (v timeToByteSlice) Scan(src interface{}) error {
	switch data := src.(type) {
	case time.Time:
		*v.val = []byte(data.Format(timeLayout))
	case []byte:
		*v.val = data
	case string:
		*v.val = []byte(data)
	}
	return nil
}
