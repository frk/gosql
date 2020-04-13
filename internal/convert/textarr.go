package convert

import (
	"database/sql/driver"
)

type TextArrayFromStringSlice struct {
	Val []string
}

func (v TextArrayFromStringSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}
	for _, s := range v.Val {
		out = pgAppendQuote1(out, []byte(s))
		out = append(out, ',')
	}
	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type TextArrayToStringSlice struct {
	Val *[]string
}

func (v TextArrayToStringSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseStringArray(data)
	strings := make([]string, len(elems))
	for i := 0; i < len(elems); i++ {
		strings[i] = string(elems[i])
	}

	*v.Val = strings
	return nil
}

type TextArrayFromByteSliceSlice struct {
	Val [][]byte
}

func (v TextArrayFromByteSliceSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}
	for _, s := range v.Val {
		if s == nil {
			out = append(out, 'N', 'U', 'L', 'L')
		} else {
			out = pgAppendQuote1(out, s)
		}
		out = append(out, ',')
	}
	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type TextArrayToByteSliceSlice struct {
	Val *[][]byte
}

func (v TextArrayToByteSliceSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseStringArray(data)
	bytess := make([][]byte, len(elems))
	for i := 0; i < len(elems); i++ {
		if elems[i] != nil {
			bytess[i] = make([]byte, len(elems[i]))
			copy(bytess[i], elems[i])
		}
	}

	*v.Val = bytess
	return nil
}
