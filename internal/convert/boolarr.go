package convert

import (
	"database/sql/driver"
)

type BoolArrayFromBoolSlice struct {
	Val []bool
}

func (s BoolArrayFromBoolSlice) Value() (driver.Value, error) {
	if s.Val == nil {
		return nil, nil
	}
	if n := len(s.Val); n > 0 {
		out := make([]byte, 1+(n*2))
		out[0] = '{'

		j := 1
		for i := 0; i < n; i++ {
			if s.Val[i] {
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

type BoolArrayToBoolSlice struct {
	Val *[]bool
}

func (s BoolArrayToBoolSlice) Scan(src interface{}) error {
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
		if elems[i][0] == 't' {
			bools[i] = true
		} else {
			bools[i] = false
		}
	}

	*s.Val = bools
	return nil
}
