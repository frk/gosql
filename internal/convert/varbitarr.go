package convert

import (
	"database/sql/driver"
	"strconv"
)

type VarBitArrayFromBoolSliceSlice struct {
	Val [][]bool
}

func (v VarBitArrayFromBoolSliceSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := 2 + (len(v.Val) - 1)
	for i := 0; i < len(v.Val); i++ {
		if v.Val[i] == nil {
			size += 4 // len("NULL")
		} else if len(v.Val[i]) == 0 {
			size += 2 // len(`""`)
		} else {
			size += len(v.Val[i])
		}
	}

	out := make([]byte, size)
	out[0] = '{'

	idx := 1
	for i := 0; i < len(v.Val); i++ {
		if v.Val[i] == nil {
			out[idx+0] = 'N'
			out[idx+1] = 'U'
			out[idx+2] = 'L'
			out[idx+3] = 'L'
			idx += 4
		} else if len(v.Val[i]) == 0 {
			out[idx+0] = '"'
			out[idx+1] = '"'
			idx += 2
		} else {
			for j := 0; j < len(v.Val[i]); j++ {
				if v.Val[i][j] {
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

type VarBitArrayToBoolSliceSlice struct {
	Val *[][]bool
}

func (v VarBitArrayToBoolSliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
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

	*v.Val = boolss
	return nil
}

type VarBitArrayFromUint8SliceSlice struct {
	Val [][]uint8
}

func (v VarBitArrayFromUint8SliceSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := 2 + (len(v.Val) - 1)
	for i := 0; i < len(v.Val); i++ {
		if v.Val[i] == nil {
			size += 4 // len("NULL")
		} else if len(v.Val[i]) == 0 {
			size += 2 // len(`""`)
		} else {
			size += len(v.Val[i])
		}
	}

	out := make([]byte, size)
	out[0] = '{'

	idx := 1
	for i := 0; i < len(v.Val); i++ {
		if v.Val[i] == nil {
			out[idx+0] = 'N'
			out[idx+1] = 'U'
			out[idx+2] = 'L'
			out[idx+3] = 'L'
			idx += 4
		} else if len(v.Val[i]) == 0 {
			out[idx+0] = '"'
			out[idx+1] = '"'
			idx += 2
		} else {
			for j := 0; j < len(v.Val[i]); j++ {
				if v.Val[i][j] == 1 {
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

type VarBitArrayToUint8SliceSlice struct {
	Val *[][]uint8
}

func (v VarBitArrayToUint8SliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
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

	*v.Val = uint8ss
	return nil
}

type VarBitArrayFromStringSlice struct {
	Val []string
}

func (v VarBitArrayFromStringSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := 2 + (len(v.Val) - 1)
	for i := 0; i < len(v.Val); i++ {
		if len(v.Val[i]) == 0 {
			size += 2 // len(`""`)
		} else {
			size += len(v.Val[i])
		}
	}

	out := make([]byte, size)
	out[0] = '{'

	idx := 1
	for i := 0; i < len(v.Val); i++ {
		if length := len(v.Val[i]); length == 0 {
			out[idx+0] = '"'
			out[idx+1] = '"'
			idx += 2
		} else {
			copy(out[idx:idx+length], []byte(v.Val[i]))
			idx += length
		}

		out[idx] = ','
		idx += 1
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type VarBitArrayToStringSlice struct {
	Val *[]string
}

func (v VarBitArrayToStringSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
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

	*v.Val = strings
	return nil
}

type VarBitArrayFromInt64Slice struct {
	Val []int64
}

func (v VarBitArrayFromInt64Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}
	for i := 0; i < len(v.Val); i++ {
		out = strconv.AppendInt(out, v.Val[i], 2)
		out = append(out, ',')
	}
	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type VarBitArrayToInt64Slice struct {
	Val *[]int64
}

func (v VarBitArrayToInt64Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
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

	*v.Val = int64s
	return nil
}
