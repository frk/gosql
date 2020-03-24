package convert

import (
	"database/sql/driver"
	"unicode/utf8"
)

type CharFromByte struct {
	B byte
}

func (c CharFromByte) Value() (driver.Value, error) {
	return []byte{c.B}, nil
}

type CharFromRune struct {
	R rune
}

func (c CharFromRune) Value() (driver.Value, error) {
	return string(c.R), nil
}

type CharToByte struct {
	B *byte
}

func (c CharToByte) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		if len(v) > 0 {
			*c.B = v[0]
		}
	case []byte:
		if len(v) > 0 {
			*c.B = v[0]
		}
	}
	return nil
}

type CharToRune struct {
	R *rune
}

func (c CharToRune) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		if len(v) > 0 {
			*c.R, _ = utf8.DecodeRuneInString(v)
		}
	case []byte:
		if len(v) > 0 {
			*c.R, _ = utf8.DecodeRune(v)
		}
	}
	return nil
}
