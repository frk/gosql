package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// VarBitArrayFromBoolSliceSlice returns a driver.Valuer that produces a PostgreSQL varbit[] from the given Go [][]bool.
func VarBitArrayFromBoolSliceSlice(val [][]bool) driver.Valuer {
	return varBitArrayFromBoolSliceSlice{val: val}
}

// VarBitArrayToBoolSliceSlice returns an sql.Scanner that converts a PostgreSQL varbit[] into a Go [][]bool and sets it to val.
func VarBitArrayToBoolSliceSlice(val *[][]bool) sql.Scanner {
	return varBitArrayToBoolSliceSlice{val: val}
}

// VarBitArrayFromUint8SliceSlice returns a driver.Valuer that produces a PostgreSQL varbit[] from the given Go [][]uint8.
func VarBitArrayFromUint8SliceSlice(val [][]uint8) driver.Valuer {
	return varBitArrayFromUint8SliceSlice{val: val}
}

// VarBitArrayToUint8SliceSlice returns an sql.Scanner that converts a PostgreSQL varbit[] into a Go [][]uint8 and sets it to val.
func VarBitArrayToUint8SliceSlice(val *[][]uint8) sql.Scanner {
	return varBitArrayToUint8SliceSlice{val: val}
}

// VarBitArrayFromStringSlice returns a driver.Valuer that produces a PostgreSQL varbit[] from the given Go []string.
func VarBitArrayFromStringSlice(val []string) driver.Valuer {
	return varBitArrayFromStringSlice{val: val}
}

// VarBitArrayToStringSlice returns an sql.Scanner that converts a PostgreSQL varbit[] into a Go []string and sets it to val.
func VarBitArrayToStringSlice(val *[]string) sql.Scanner {
	return varBitArrayToStringSlice{val: val}
}

// VarBitArrayFromInt64Slice returns a driver.Valuer that produces a PostgreSQL varbit[] from the given Go []int64.
func VarBitArrayFromInt64Slice(val []int64) driver.Valuer {
	return varBitArrayFromInt64Slice{val: val}
}

// VarBitArrayToInt64Slice returns an sql.Scanner that converts a PostgreSQL varbit[] into a Go []int64 and sets it to val.
func VarBitArrayToInt64Slice(val *[]int64) sql.Scanner {
	return varBitArrayToInt64Slice{val: val}
}

type varBitArrayFromBoolSliceSlice struct {
	val [][]bool
}

func (v varBitArrayFromBoolSliceSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := 2 + (len(v.val) - 1)
	for i := 0; i < len(v.val); i++ {
		if v.val[i] == nil {
			size += 4 // len("NULL")
		} else if len(v.val[i]) == 0 {
			size += 2 // len(`""`)
		} else {
			size += len(v.val[i])
		}
	}

	out := make([]byte, size)
	out[0] = '{'

	idx := 1
	for i := 0; i < len(v.val); i++ {
		if v.val[i] == nil {
			out[idx+0] = 'N'
			out[idx+1] = 'U'
			out[idx+2] = 'L'
			out[idx+3] = 'L'
			idx += 4
		} else if len(v.val[i]) == 0 {
			out[idx+0] = '"'
			out[idx+1] = '"'
			idx += 2
		} else {
			for j := 0; j < len(v.val[i]); j++ {
				if v.val[i][j] {
					out[idx] = '1'
				} else {
					out[idx] = '0'
				}
				idx += 1
			}
		}

		out[idx] = ','
		idx += 1
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type varBitArrayToBoolSliceSlice struct {
	val *[][]bool
}

func (v varBitArrayToBoolSliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	boolss := make([][]bool, len(elems))
	for i := 0; i < len(elems); i++ {
		if len(elems[i]) == 4 && elems[i][0] == 'N' { // NULL?
			continue
		}
		if len(elems[i]) == 2 && elems[i][0] == '"' { // ""?
			boolss[i] = []bool{}
			continue
		}

		bools := make([]bool, len(elems[i]))
		for j := 0; j < len(elems[i]); j++ {
			if elems[i][j] == '1' {
				bools[j] = true
			}
		}

		boolss[i] = bools
	}

	*v.val = boolss
	return nil
}

type varBitArrayFromUint8SliceSlice struct {
	val [][]uint8
}

func (v varBitArrayFromUint8SliceSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := 2 + (len(v.val) - 1)
	for i := 0; i < len(v.val); i++ {
		if v.val[i] == nil {
			size += 4 // len("NULL")
		} else if len(v.val[i]) == 0 {
			size += 2 // len(`""`)
		} else {
			size += len(v.val[i])
		}
	}

	out := make([]byte, size)
	out[0] = '{'

	idx := 1
	for i := 0; i < len(v.val); i++ {
		if v.val[i] == nil {
			out[idx+0] = 'N'
			out[idx+1] = 'U'
			out[idx+2] = 'L'
			out[idx+3] = 'L'
			idx += 4
		} else if len(v.val[i]) == 0 {
			out[idx+0] = '"'
			out[idx+1] = '"'
			idx += 2
		} else {
			for j := 0; j < len(v.val[i]); j++ {
				if v.val[i][j] == 1 {
					out[idx] = '1'
				} else {
					out[idx] = '0'
				}
				idx += 1
			}
		}

		out[idx] = ','
		idx += 1
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type varBitArrayToUint8SliceSlice struct {
	val *[][]uint8
}

func (v varBitArrayToUint8SliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	uint8ss := make([][]uint8, len(elems))
	for i := 0; i < len(elems); i++ {
		if len(elems[i]) == 4 && elems[i][0] == 'N' { // NULL?
			continue
		}
		if len(elems[i]) == 2 && elems[i][0] == '"' { // ""?
			uint8ss[i] = []uint8{}
			continue
		}

		uint8s := make([]uint8, len(elems[i]))
		for j := 0; j < len(elems[i]); j++ {
			if elems[i][j] == '1' {
				uint8s[j] = 1
			}
		}

		uint8ss[i] = uint8s
	}

	*v.val = uint8ss
	return nil
}

type varBitArrayFromStringSlice struct {
	val []string
}

func (v varBitArrayFromStringSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := 2 + (len(v.val) - 1)
	for i := 0; i < len(v.val); i++ {
		if len(v.val[i]) == 0 {
			size += 2 // len(`""`)
		} else {
			size += len(v.val[i])
		}
	}

	out := make([]byte, size)
	out[0] = '{'

	idx := 1
	for i := 0; i < len(v.val); i++ {
		if length := len(v.val[i]); length == 0 {
			out[idx+0] = '"'
			out[idx+1] = '"'
			idx += 2
		} else {
			copy(out[idx:idx+length], []byte(v.val[i]))
			idx += length
		}

		out[idx] = ','
		idx += 1
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type varBitArrayToStringSlice struct {
	val *[]string
}

func (v varBitArrayToStringSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	strings := make([]string, len(elems))
	for i := 0; i < len(elems); i++ {
		if len(elems[i]) == 2 && elems[i][0] == '"' { // ""?
			continue
		}

		strings[i] = string(elems[i])
	}

	*v.val = strings
	return nil
}

type varBitArrayFromInt64Slice struct {
	val []int64
}

func (v varBitArrayFromInt64Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}
	for i := 0; i < len(v.val); i++ {
		out = strconv.AppendInt(out, v.val[i], 2)
		out = append(out, ',')
	}
	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type varBitArrayToInt64Slice struct {
	val *[]int64
}

func (v varBitArrayToInt64Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	int64s := make([]int64, len(elems))
	for i := 0; i < len(elems); i++ {
		i64, err := strconv.ParseInt(string(elems[i]), 2, 64)
		if err != nil {
			return err
		}
		int64s[i] = i64
	}

	*v.val = int64s
	return nil
}
