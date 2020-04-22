package pgsql

import (
	"database/sql"
	"database/sql/driver"
)

// BPCharArrayFromString returns a driver.Valuer that produces a PostgreSQL bpchar[] from the given Go string.
func BPCharArrayFromString(val string) driver.Valuer {
	return charArrayFromString{val: val}
}

// BPCharArrayToString returns an sql.Scanner that converts a PostgreSQL bpchar[] into a Go string and sets it to val.
func BPCharArrayToString(val *string) sql.Scanner {
	return charArrayToString{val: val}
}

// BPCharArrayFromByteSlice returns a driver.Valuer that produces a PostgreSQL bpchar[] from the given Go []byte.
func BPCharArrayFromByteSlice(val []byte) driver.Valuer {
	return charArrayFromByteSlice{val: val}
}

// BPCharArrayToByteSlice returns an sql.Scanner that converts a PostgreSQL bpchar[] into a Go []byte and sets it to val.
func BPCharArrayToByteSlice(val *[]byte) sql.Scanner {
	return charArrayToByteSlice{val: val}
}

// BPCharArrayFromRuneSlice returns a driver.Valuer that produces a PostgreSQL bpchar[] from the given Go []rune.
func BPCharArrayFromRuneSlice(val []rune) driver.Valuer {
	return charArrayFromRuneSlice{val: val}
}

// BPCharArrayToRuneSlice returns an sql.Scanner that converts a PostgreSQL bpchar[] into a Go []rune and sets it to val.
func BPCharArrayToRuneSlice(val *[]rune) sql.Scanner {
	return charArrayToRuneSlice{val: val}
}

// BPCharArrayFromStringSlice returns a driver.Valuer that produces a PostgreSQL bpchar[] from the given Go []string.
func BPCharArrayFromStringSlice(val []string) driver.Valuer {
	return charArrayFromStringSlice{val: val}
}

// BPCharArrayToStringSlice returns an sql.Scanner that converts a PostgreSQL bpchar[] into a Go []string and sets it to val.
func BPCharArrayToStringSlice(val *[]string) sql.Scanner {
	return charArrayToStringSlice{val: val}
}
