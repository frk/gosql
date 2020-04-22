package pgsql

import (
	"database/sql"
	"time"
)

// DateToTime returns an sql.Scanner that converts a PostgreSQL date into a Go time.Time and sets it to val.
func DateToTime(val *time.Time) sql.Scanner {
	return dateToTime{val: val}
}

// DateToString returns an sql.Scanner that converts a PostgreSQL date into a Go string and sets it to val.
func DateToString(val *string) sql.Scanner {
	return dateToString{val: val}
}

// DateToByteSlice returns an sql.Scanner that converts a PostgreSQL date into a Go []byte and sets it to val.
func DateToByteSlice(val *[]byte) sql.Scanner {
	return dateToByteSlice{val: val}
}

type dateToTime struct {
	val *time.Time
}

func (v dateToTime) Scan(src interface{}) error {
	if t, ok := src.(time.Time); ok {
		*v.val = t.In(time.UTC)
	} else {
		*v.val = time.Time{}
	}
	return nil
}

type dateToString struct {
	val *string
}

func (v dateToString) Scan(src interface{}) error {
	if t, ok := src.(time.Time); ok {
		*v.val = t.Format(dateLayout)
	} else {
		*v.val = ""
	}
	return nil
}

type dateToByteSlice struct {
	val *[]byte
}

func (v dateToByteSlice) Scan(src interface{}) error {
	if t, ok := src.(time.Time); ok {
		*v.val = []byte(t.Format(dateLayout))
	} else {
		*v.val = nil
	}
	return nil
}
