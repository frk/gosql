package convert

import (
	"database/sql/driver"
)

type JSONArrayFromByteSliceSlice struct {
	Val [][]byte
}

func (v JSONArrayFromByteSliceSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.Val {
		if a == nil {
			out = append(out, 'N', 'U', 'L', 'L', ',')
			continue
		}

		out = pgAppendQuote1(out, a)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last ";" with "}"
	return out, nil
}

type JSONArrayToByteSliceSlice struct {
	Val *[][]byte
}

func (v JSONArrayToByteSliceSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseStringArray(data)
	jsons := make([][]byte, len(elems))
	for i := 0; i < len(elems); i++ {
		if elems[i] == nil {
			jsons[i] = nil
		} else {
			jsons[i] = make([]byte, len(elems[i]))
			copy(jsons[i], elems[i])
		}
	}

	*v.Val = jsons
	return nil
}
