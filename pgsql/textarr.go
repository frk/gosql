package pgsql

import (
	"database/sql"
	"database/sql/driver"
)

// TextArrayFromStringSlice returns a driver.Valuer that produces a PostgreSQL text[] from the given Go []string.
func TextArrayFromStringSlice(val []string) driver.Valuer {
	return textArrayFromStringSlice{val: val}
}

// TextArrayToStringSlice returns an sql.Scanner that converts a PostgreSQL text[] into a Go []string and sets it to val.
func TextArrayToStringSlice(val *[]string) sql.Scanner {
	return textArrayToStringSlice{val: val}
}

// TextArrayFromByteSliceSlice returns a driver.Valuer that produces a PostgreSQL text[] from the given Go [][]byte.
func TextArrayFromByteSliceSlice(val [][]byte) driver.Valuer {
	return textArrayFromByteSliceSlice{val: val}
}

// TextArrayToByteSliceSlice returns an sql.Scanner that converts a PostgreSQL text[] into a Go [][]byte and sets it to val.
func TextArrayToByteSliceSlice(val *[][]byte) sql.Scanner {
	return textArrayToByteSliceSlice{val: val}
}

type textArrayFromStringSlice struct {
	val []string
}

func (v textArrayFromStringSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}
	for _, s := range v.val {
		out = pgAppendQuote1(out, []byte(s))
		out = append(out, ',')
	}
	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type textArrayToStringSlice struct {
	val *[]string
}

func (v textArrayToStringSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseStringArray(data)
	strings := make([]string, len(elems))
	for i := 0; i < len(elems); i++ {
		strings[i] = string(elems[i])
	}

	*v.val = strings
	return nil
}

type textArrayFromByteSliceSlice struct {
	val [][]byte
}

func (v textArrayFromByteSliceSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}
	for _, s := range v.val {
		if s == nil {
			out = append(out, 'N', 'U', 'L', 'L')
		} else {
			out = pgAppendQuote1(out, s)
		}
		out = append(out, ',')
	}
	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type textArrayToByteSliceSlice struct {
	val *[][]byte
}

func (v textArrayToByteSliceSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseStringArray(data)
	bytess := make([][]byte, len(elems))
	for i := 0; i < len(elems); i++ {
		if elems[i] != nil {
			bytess[i] = make([]byte, len(elems[i]))
			copy(bytess[i], elems[i])
		}
	}

	*v.val = bytess
	return nil
}
