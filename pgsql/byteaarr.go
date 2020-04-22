package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
)

// ByteaArrayFromByteSliceSlice returns a driver.Valuer that produces a PostgreSQL bytea[] from the given Go [][]byte.
func ByteaArrayFromByteSliceSlice(val [][]byte) driver.Valuer {
	return byteaArrayFromByteSliceSlice{val: val}
}

// ByteaArrayToByteSliceSlice returns an sql.Scanner that converts a PostgreSQL bytea[] into a Go [][]byte and sets it to val.
func ByteaArrayToByteSliceSlice(val *[][]byte) sql.Scanner {
	return byteaArrayToByteSliceSlice{val: val}
}

// ByteaArrayFromStringSlice returns a driver.Valuer that produces a PostgreSQL bytea[] from the given Go []string.
func ByteaArrayFromStringSlice(val []string) driver.Valuer {
	return byteaArrayFromStringSlice{val: val}
}

// ByteaArrayToStringSlice returns an sql.Scanner that converts a PostgreSQL bytea[] into a Go []string and sets it to val.
func ByteaArrayToStringSlice(val *[]string) sql.Scanner {
	return byteaArrayToStringSlice{val: val}
}

type byteaArrayFromByteSliceSlice struct {
	val [][]byte
}

func (v byteaArrayFromByteSliceSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := make([]byte, 1, (len(v.val)*3)+1)
	out[0] = '{'

	for i := 0; i < len(v.val); i++ {
		out = append(out, '"', '\\', '\\', 'x')

		dst := make([]byte, hex.EncodedLen(len(v.val[i])))
		_ = hex.Encode(dst, v.val[i])
		out = append(out, dst...)

		out = append(out, '"', ',')
	}

	out[len(out)-1] = '}'
	return out, nil
}

type byteaArrayToByteSliceSlice struct {
	val *[][]byte
}

func (v byteaArrayToByteSliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	out := make([][]byte, len(elems))
	for i := 0; i < len(elems); i++ {
		src := elems[i]

		// drop the initial "\\x and the last "
		src = src[4 : len(src)-1]

		dst := make([]byte, hex.DecodedLen(len(src)))
		if _, err := hex.Decode(dst, src); err != nil {
			return err
		}

		out[i] = dst
	}

	*v.val = out
	return nil
}

type byteaArrayFromStringSlice struct {
	val []string
}

func (v byteaArrayFromStringSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := make([]byte, 1, (len(v.val)*3)+1)
	out[0] = '{'

	for i := 0; i < len(v.val); i++ {
		out = append(out, '"', '\\', '\\', 'x')

		src := []byte(v.val[i])
		dst := make([]byte, hex.EncodedLen(len(src)))
		_ = hex.Encode(dst, src)
		out = append(out, dst...)

		out = append(out, '"', ',')
	}

	out[len(out)-1] = '}'
	return out, nil
}

type byteaArrayToStringSlice struct {
	val *[]string
}

func (v byteaArrayToStringSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	out := make([]string, len(elems))
	for i := 0; i < len(elems); i++ {
		src := elems[i]

		// drop the initial "\\x and the last "
		src = src[4 : len(src)-1]

		dst := make([]byte, hex.DecodedLen(len(src)))
		if _, err := hex.Decode(dst, src); err != nil {
			return err
		}

		out[i] = string(dst)
	}

	*v.val = out
	return nil
}
