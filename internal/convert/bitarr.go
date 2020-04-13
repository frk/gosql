package convert

import (
	"database/sql/driver"
)

type BitArrayFromBoolSlice struct {
	Val []bool
}

func (s BitArrayFromBoolSlice) Value() (driver.Value, error) {
	if s.Val == nil {
		return nil, nil
	}

	if n := len(s.Val); n > 0 {
		out := make([]byte, 1+(n*2))
		out[0] = '{'

		j := 1
		for i := 0; i < n; i++ {
			if s.Val[i] {
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

type BitArrayFromUint8Slice struct {
	Val []uint8
}

func (s BitArrayFromUint8Slice) Value() (driver.Value, error) {
	if s.Val == nil {
		return nil, nil
	}

	if n := len(s.Val); n > 0 {
		out := make([]byte, 1+(n*2))
		out[0] = '{'

		j := 1
		for i := 0; i < n; i++ {
			if s.Val[i] == 0 {
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

type BitArrayFromUintSlice struct {
	Val []uint
}

func (s BitArrayFromUintSlice) Value() (driver.Value, error) {
	if s.Val == nil {
		return nil, nil
	}

	if n := len(s.Val); n > 0 {
		out := make([]byte, 1+(n*2))
		out[0] = '{'

		j := 1
		for i := 0; i < n; i++ {
			if s.Val[i] == 0 {
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

type BitArrayToBoolSlice struct {
	Val *[]bool
}

func (s BitArrayToBoolSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*s.Val = nil
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

	*s.Val = bools
	return nil
}

type BitArrayToUint8Slice struct {
	Val *[]uint8
}

func (s BitArrayToUint8Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*s.Val = nil
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

	*s.Val = uint8s
	return nil
}

type BitArrayToUintSlice struct {
	Val *[]uint
}

func (s BitArrayToUintSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*s.Val = nil
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

	*s.Val = uints
	return nil
}
