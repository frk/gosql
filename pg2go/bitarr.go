package pg2go

import (
	"database/sql"
	"database/sql/driver"
)

// BitArrayFromBoolSlice returns a driver.Valuer that produces a PostgreSQL bit[] from the given Go []bool.
func BitArrayFromBoolSlice(val []bool) driver.Valuer {
	return bitArrayFromBoolSlice{val: val}
}

// BitArrayToBoolSlice returns an sql.Scanner that converts a PostgreSQL bit[] into a Go []bool and sets it to val.
func BitArrayToBoolSlice(val *[]bool) sql.Scanner {
	return bitArrayToBoolSlice{val: val}
}

// BitArrayFromUint8Slice returns a driver.Valuer that produces a PostgreSQL bit[] from the given Go []uint8.
func BitArrayFromUint8Slice(val []uint8) driver.Valuer {
	return bitArrayFromUint8Slice{val: val}
}

// BitArrayToUint8Slice returns an sql.Scanner that converts a PostgreSQL bit[] to a Go []uint8 and sets it to val.
func BitArrayToUint8Slice(val *[]uint8) sql.Scanner {
	return bitArrayToUint8Slice{val: val}
}

// BitArrayFromUintSlice returns a driver.Valuer that produces a PostgreSQL bit[] from the given Go []uint.
func BitArrayFromUintSlice(val []uint) driver.Valuer {
	return bitArrayFromUintSlice{val: val}
}

// BitArrayToUintSlice returns an sql.Scanner that converts a PostgreSQL bit[] to a Go []uint and sets it to val.
func BitArrayToUintSlice(val *[]uint) sql.Scanner {
	return bitArrayToUintSlice{val: val}
}

type bitArrayFromBoolSlice struct {
	val []bool
}

func (v bitArrayFromBoolSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	}

	if n := len(v.val); n > 0 {
		out := make([]byte, 1+(n*2))
		out[0] = '{'

		j := 1
		for i := 0; i < n; i++ {
			if v.val[i] {
				out[j] = '1'
			} else {
				out[j] = '0'
			}
			out[j+1] = ','
			j += 2
		}

		out[len(out)-1] = '}'
		return out, nil
	}

	return []byte{'{', '}'}, nil
}

type bitArrayFromUint8Slice struct {
	val []uint8
}

func (v bitArrayFromUint8Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	}

	if n := len(v.val); n > 0 {
		out := make([]byte, 1+(n*2))
		out[0] = '{'

		j := 1
		for i := 0; i < n; i++ {
			if v.val[i] == 0 {
				out[j] = '0'
			} else {
				out[j] = '1'
			}
			out[j+1] = ','
			j += 2
		}

		out[len(out)-1] = '}'
		return out, nil
	}

	return []byte{'{', '}'}, nil
}

type bitArrayFromUintSlice struct {
	val []uint
}

func (v bitArrayFromUintSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	}

	if n := len(v.val); n > 0 {
		out := make([]byte, 1+(n*2))
		out[0] = '{'

		j := 1
		for i := 0; i < n; i++ {
			if v.val[i] == 0 {
				out[j] = '0'
			} else {
				out[j] = '1'
			}
			out[j+1] = ','
			j += 2
		}

		out[len(out)-1] = '}'
		return out, nil
	}

	return []byte{'{', '}'}, nil
	return nil, nil
}

type bitArrayToBoolSlice struct {
	val *[]bool
}

func (v bitArrayToBoolSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	bools := make([]bool, len(elems))
	for i := 0; i < len(elems); i++ {
		if elems[i][0] == '1' {
			bools[i] = true
		} else {
			bools[i] = false
		}
	}

	*v.val = bools
	return nil
}

type bitArrayToUint8Slice struct {
	val *[]uint8
}

func (v bitArrayToUint8Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	uint8s := make([]uint8, len(elems))
	for i := 0; i < len(elems); i++ {
		if elems[i][0] == '1' {
			uint8s[i] = 1
		} else {
			uint8s[i] = 0
		}
	}

	*v.val = uint8s
	return nil
}

type bitArrayToUintSlice struct {
	val *[]uint
}

func (v bitArrayToUintSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	uints := make([]uint, len(elems))
	for i := 0; i < len(elems); i++ {
		if elems[i][0] == '1' {
			uints[i] = 1
		} else {
			uints[i] = 0
		}
	}

	*v.val = uints
	return nil
}
