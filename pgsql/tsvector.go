package pgsql

import (
	"database/sql"
	"database/sql/driver"
)

// TSVectorFromStringSlice returns a driver.Valuer that produces a PostgreSQL tsvector from the given Go []string.
func TSVectorFromStringSlice(val []string) driver.Valuer {
	return tsVectorFromStringSlice{val: val}
}

// TSVectorToStringSlice returns an sql.Scanner that converts a PostgreSQL tsvector into a Go []string and sets it to val.
func TSVectorToStringSlice(val *[]string) sql.Scanner {
	return tsVectorToStringSlice{val: val}
}

// TSVectorFromByteSliceSlice returns a driver.Valuer that produces a PostgreSQL tsvector from the given Go [][]byte.
func TSVectorFromByteSliceSlice(val [][]byte) driver.Valuer {
	return tsVectorFromByteSliceSlice{val: val}
}

// TSVectorToByteSliceSlice returns an sql.Scanner that converts a PostgreSQL tsvector into a Go [][]byte and sets it to val.
func TSVectorToByteSliceSlice(val *[][]byte) sql.Scanner {
	return tsVectorToByteSliceSlice{val: val}
}

type tsVectorFromStringSlice struct {
	val []string
}

func (v tsVectorFromStringSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{}, nil
	}

	size := len(v.val) - 1
	for i := 0; i < len(v.val); i++ {
		size += len(v.val[i])
	}
	out := make([]byte, size)

	var pos int
	for i := 0; i < len(v.val); i++ {
		length := len(v.val[i])
		copy(out[pos:pos+length], v.val[i])
		pos += length

		if pos < size {
			out[pos] = ' '
			pos += 1
		}
	}
	return out, nil
}

type tsVectorToStringSlice struct {
	val *[]string
}

func (v tsVectorToStringSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseVector(data)
	strings := make([]string, len(elems))
	for i := 0; i < len(elems); i++ {
		strings[i] = string(elems[i])
	}

	*v.val = strings
	return nil
}

type tsVectorFromByteSliceSlice struct {
	val [][]byte
}

func (v tsVectorFromByteSliceSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{}, nil
	}

	size := len(v.val) - 1
	for i := 0; i < len(v.val); i++ {
		size += len(v.val[i])
	}
	out := make([]byte, size)

	var pos int
	for i := 0; i < len(v.val); i++ {
		length := len(v.val[i])
		copy(out[pos:pos+length], v.val[i])
		pos += length

		if pos < size {
			out[pos] = ' '
			pos += 1
		}
	}
	return out, nil
}

type tsVectorToByteSliceSlice struct {
	val *[][]byte
}

func (v tsVectorToByteSliceSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseVector(data)
	bytess := make([][]byte, len(elems))
	for i := 0; i < len(elems); i++ {
		bytess[i] = make([]byte, len(elems[i]))
		copy(bytess[i], elems[i])
	}

	*v.val = bytess
	return nil
}
