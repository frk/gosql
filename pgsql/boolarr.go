package pgsql

import (
	"database/sql"
	"database/sql/driver"
)

// BoolArrayFromBoolSlice returns a driver.Valuer that produces a PostgreSQL boolean[] from the given Go []bool.
func BoolArrayFromBoolSlice(val []bool) driver.Valuer {
	return boolArrayFromBoolSlice{val: val}
}

// BoolArrayToBoolSlice returns an sql.Scanner that converts a PostgreSQL boolean[] into a Go []bool and sets it to val.
func BoolArrayToBoolSlice(val *[]bool) sql.Scanner {
	return boolArrayToBoolSlice{val: val}
}

type boolArrayFromBoolSlice struct {
	val []bool
}

func (v boolArrayFromBoolSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	}
	if n := len(v.val); n > 0 {
		out := make([]byte, 1+(n*2))
		out[0] = '{'

		j := 1
		for i := 0; i < n; i++ {
			if v.val[i] {
				out[j] = 't'
			} else {
				out[j] = 'f'
			}
			out[j+1] = ','
			j += 2
		}

		out[len(out)-1] = '}'
		return out, nil
	}
	return []byte{'{', '}'}, nil
}

type boolArrayToBoolSlice struct {
	val *[]bool
}

func (v boolArrayToBoolSlice) Scan(src interface{}) error {
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
		if elems[i][0] == 't' {
			bools[i] = true
		} else {
			bools[i] = false
		}
	}

	*v.val = bools
	return nil
}
