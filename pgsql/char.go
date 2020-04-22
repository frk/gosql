package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"unicode/utf8"
)

// CharFromByte returns a driver.Valuer that produces a PostgreSQL char from the given Go byte.
func CharFromByte(val byte) driver.Valuer {
	return charFromByte{val: val}
}

// CharToByte returns an sql.Scanner that converts a PostgreSQL char into a Go byte and sets it to val.
func CharToByte(val *byte) sql.Scanner {
	return charToByte{val: val}
}

// CharFromRune returns a driver.Valuer that produces a PostgreSQL char from the given Go rune.
func CharFromRune(val rune) driver.Valuer {
	return charFromRune{val: val}
}

// CharToRune returns an sql.Scanner that converts a PostgreSQL char into a Go rune and sets it to val.
func CharToRune(val *rune) sql.Scanner {
	return charToRune{val: val}
}

type charFromByte struct {
	val byte
}

func (v charFromByte) Value() (driver.Value, error) {
	return []byte{v.val}, nil
}

type charToByte struct {
	val *byte
}

func (v charToByte) Scan(src interface{}) error {
	switch data := src.(type) {
	case string:
		if len(data) > 0 {
			*v.val = data[0]
		}
	case []byte:
		if len(data) > 0 {
			*v.val = data[0]
		}
	}
	return nil
}

type charFromRune struct {
	val rune
}

func (v charFromRune) Value() (driver.Value, error) {
	return string(v.val), nil
}

type charToRune struct {
	val *rune
}

func (v charToRune) Scan(src interface{}) error {
	switch data := src.(type) {
	case string:
		if len(data) > 0 {
			*v.val, _ = utf8.DecodeRuneInString(data)
		}
	case []byte:
		if len(data) > 0 {
			*v.val, _ = utf8.DecodeRune(data)
		}
	}
	return nil
}
