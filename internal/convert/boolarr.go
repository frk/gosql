package convert

import (
	"database/sql/driver"
)

type BoolArrayFromBoolSlice struct {
	S []bool
}

func (s BoolArrayFromBoolSlice) Value() (driver.Value, error) {
	if s.S == nil {
		return nil, nil
	}
	if n := len(s.S); n > 0 {
		out := make([]byte, 1+(n*2))
		out[0] = '{'

		j := 1
		for i := 0; i < n; i++ {
			if s.S[i] {
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
	S *[]bool
}

func (s BoolArrayToBoolSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*s.S = nil
		return nil
	}

	elems := pgparsearray1(arr)
	bools := make([]bool, len(elems))
	for i := 0; i < len(elems); i++ {
		if elems[i][0] == 't' {
			bools[i] = true
		} else {
			bools[i] = false
		}
	}

	*s.S = bools
	return nil
}
