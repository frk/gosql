package convert

import (
	"database/sql/driver"
)

type XMLArrayFromByteSliceSlice struct {
	Val [][]byte
}

func (v XMLArrayFromByteSliceSlice) Value() (driver.Value, error) {
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

type XMLArrayToByteSliceSlice struct {
	Val *[][]byte
}

func (v XMLArrayToByteSliceSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseStringArray(data)
	xmls := make([][]byte, len(elems))
	for i := 0; i < len(elems); i++ {
		if elems[i] == nil {
			xmls[i] = nil
		} else {
			xmls[i] = make([]byte, len(elems[i]))
			copy(xmls[i], elems[i])
		}
	}

	*v.Val = xmls
	return nil
}
