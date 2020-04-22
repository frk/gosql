package pgsql

import (
	"database/sql"
	"database/sql/driver"
)

// TSQueryArrayFromStringSlice returns a driver.Valuer that produces a PostgreSQL tsquery[] from the given Go []string.
func TSQueryArrayFromStringSlice(val []string) driver.Valuer {
	return tsQueryArrayFromStringSlice{val: val}
}

// TSQueryArrayToStringSlice returns an sql.Scanner that converts a PostgreSQL tsquery[] into a Go []string and sets it to val.
func TSQueryArrayToStringSlice(val *[]string) sql.Scanner {
	return tsQueryArrayToStringSlice{val: val}
}

// TSQueryArrayFromByteSliceSlice returns a driver.Valuer that produces a PostgreSQL tsquery[] from the given Go [][]byte.
func TSQueryArrayFromByteSliceSlice(val [][]byte) driver.Valuer {
	return tsQueryArrayFromByteSliceSlice{val: val}
}

// TSQueryArrayToByteSliceSlice returns an sql.Scanner that converts a PostgreSQL tsquery[] into a Go [][]byte and sets it to val.
func TSQueryArrayToByteSliceSlice(val *[][]byte) sql.Scanner {
	return tsQueryArrayToByteSliceSlice{val: val}
}

type tsQueryArrayFromStringSlice struct {
	val []string
}

func (v tsQueryArrayFromStringSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := (len(v.val) * 2) + // number of surrounding double quotes
		(len(v.val) - 1) + // number of commas between elements
		2 // curly braces
	for i := 0; i < len(v.val); i++ {
		size += len(v.val[i]) // length of the element
	}

	out := make([]byte, size)
	out[0] = '{'

	var pos int
	for i := 0; i < len(v.val); i++ {
		pos += 1
		out[pos] = '"'
		pos += 1

		length := len(v.val[i])
		copy(out[pos:pos+length], v.val[i])
		pos += length

		out[pos] = '"'
		pos += 1
		out[pos] = ','
	}

	out[pos] = '}'
	return out, nil
}

type tsQueryArrayToStringSlice struct {
	val *[]string
}

func (v tsQueryArrayToStringSlice) Scan(src interface{}) error {
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

type tsQueryArrayFromByteSliceSlice struct {
	val [][]byte
}

func (v tsQueryArrayFromByteSliceSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := (len(v.val) * 2) + // number of surrounding double quotes
		(len(v.val) - 1) + // number of commas between elements
		2 // curly braces
	for i := 0; i < len(v.val); i++ {
		if v.val[i] == nil {
			size += 2 // len(`NULL`) - 2 double quotes
		} else {
			size += len(v.val[i]) // length of the element
		}
	}

	out := make([]byte, size)
	out[0] = '{'

	var pos int
	for i := 0; i < len(v.val); i++ {
		pos += 1

		if v.val[i] == nil {
			out[pos+0] = 'N'
			out[pos+1] = 'U'
			out[pos+2] = 'L'
			out[pos+3] = 'L'
			pos += 4
		} else {
			out[pos] = '"'
			pos += 1

			length := len(v.val[i])
			copy(out[pos:pos+length], v.val[i])
			pos += length

			out[pos] = '"'
			pos += 1
		}

		out[pos] = ','
	}

	out[pos] = '}'
	return out, nil
}

type tsQueryArrayToByteSliceSlice struct {
	val *[][]byte
}

func (v tsQueryArrayToByteSliceSlice) Scan(src interface{}) error {
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
		if len(elems[i]) == 4 && elems[i][0] == 'N' { // NULL?
			continue
		}

		bytess[i] = make([]byte, len(elems[i]))
		copy(bytess[i], elems[i])
	}

	*v.val = bytess
	return nil
}
