package convert

import (
	"database/sql/driver"
	"unicode/utf8"
)

type CharFromByte struct {
	Val byte
}

func (c CharFromByte) Value() (driver.Value, error) {
	return []byte{c.Val}, nil
}

type CharFromRune struct {
	Val rune
}

func (c CharFromRune) Value() (driver.Value, error) {
	return string(c.Val), nil
}

type CharToByte struct {
	Val *byte
}

func (c CharToByte) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		if len(v) > 0 {
			*c.Val = v[0]
		}
	case []byte:
		if len(v) > 0 {
			*c.Val = v[0]
		}
	}
	return nil
}

type CharToRune struct {
	Val *rune
}

func (c CharToRune) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		if len(v) > 0 {
			*c.Val, _ = utf8.DecodeRuneInString(v)
		}
	case []byte:
		if len(v) > 0 {
			*c.Val, _ = utf8.DecodeRune(v)
		}
	}
	return nil
}
