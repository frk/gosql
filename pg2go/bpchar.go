package pg2go

import (
	"database/sql"
	"database/sql/driver"
)

// BPCharFromByte returns a driver.Valuer that produces a PostgreSQL bpchar from the given Go byte.
func BPCharFromByte(val byte) driver.Valuer {
	return charFromByte{val: val}
}

// BPCharToByte returns an sql.Scanner that converts a PostgreSQL bpchar into a Go byte and sets it to val.
func BPCharToByte(val *byte) sql.Scanner {
	return charToByte{val: val}
}

// BPCharFromRune returns a driver.Valuer that produces a PostgreSQL bpchar from the given Go rune.
func BPCharFromRune(val rune) driver.Valuer {
	return charFromRune{val: val}
}

// BPCharToRune returns an sql.Scanner that converts a PostgreSQL bpchar into a Go rune and sets it to val.
func BPCharToRune(val *rune) sql.Scanner {
	return charToRune{val: val}
}
