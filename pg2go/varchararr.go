package pg2go

import (
	"database/sql"
	"database/sql/driver"
)

// VarCharArrayFromStringSlice returns a driver.Valuer that produces a PostgreSQL varchar[] from the given Go []string.
func VarCharArrayFromStringSlice(val []string) driver.Valuer {
	return textArrayFromStringSlice{val: val}
}

// VarCharArrayToStringSlice returns an sql.Scanner that converts a PostgreSQL varchar[] into a Go []string and sets it to val.
func VarCharArrayToStringSlice(val *[]string) sql.Scanner {
	return textArrayToStringSlice{val: val}
}

// VarCharArrayFromByteSliceSlice returns a driver.Valuer that produces a PostgreSQL varchar[] from the given Go [][]byte.
func VarCharArrayFromByteSliceSlice(val [][]byte) driver.Valuer {
	return textArrayFromByteSliceSlice{val: val}
}

// VarCharArrayToByteSliceSlice returns an sql.Scanner that converts a PostgreSQL varchar[] into a Go [][]byte and sets it to val.
func VarCharArrayToByteSliceSlice(val *[][]byte) sql.Scanner {
	return textArrayToByteSliceSlice{val: val}
}
