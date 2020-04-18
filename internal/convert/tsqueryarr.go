package convert

import (
	"database/sql/driver"
)

type TSQueryArrayFromStringSlice struct {
	Val []string
}

func (v TSQueryArrayFromStringSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := (len(v.Val) * 2) + // number of surrounding double quotes
		(len(v.Val) - 1) + // number of commas between elements
		2 // curly braces
	for i := 0; i < len(v.Val); i++ {
		size += len(v.Val[i]) // length of the element
	}

	out := make([]byte, size)
	out[0] = '{'

	var pos int
	for i := 0; i < len(v.Val); i++ {
		pos += 1
		out[pos] = '"'
		pos += 1

		length := len(v.Val[i])
		copy(out[pos:pos+length], v.Val[i])
		pos += length

		out[pos] = '"'
		pos += 1
		out[pos] = ','
	}

	out[pos] = '}'
	return out, nil
}

type TSQueryArrayToStringSlice struct {
	Val *[]string
}

func (v TSQueryArrayToStringSlice) Scan(src interface{}) error {
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

type TSQueryArrayFromByteSliceSlice struct {
	Val [][]byte
}

func (v TSQueryArrayFromByteSliceSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := (len(v.Val) * 2) + // number of surrounding double quotes
		(len(v.Val) - 1) + // number of commas between elements
		2 // curly braces
	for i := 0; i < len(v.Val); i++ {
		if v.Val[i] == nil {
			size += 2 // len(`NULL`) - 2 double quotes
		} else {
			size += len(v.Val[i]) // length of the element
		}
	}

	out := make([]byte, size)
	out[0] = '{'

	var pos int
	for i := 0; i < len(v.Val); i++ {
		pos += 1

		if v.Val[i] == nil {
			out[pos+0] = 'N'
			out[pos+1] = 'U'
			out[pos+2] = 'L'
			out[pos+3] = 'L'
			pos += 4
		} else {
			out[pos] = '"'
			pos += 1

			length := len(v.Val[i])
			copy(out[pos:pos+length], v.Val[i])
			pos += length

			out[pos] = '"'
			pos += 1
		}

		out[pos] = ','
	}

	out[pos] = '}'
	return out, nil
}

type TSQueryArrayToByteSliceSlice struct {
	Val *[][]byte
}

func (v TSQueryArrayToByteSliceSlice) Scan(src interface{}) error {
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
		if len(elems[i]) == 4 && elems[i][0] == 'N' { // NULL?
			continue
		}

		bytess[i] = make([]byte, len(elems[i]))
		copy(bytess[i], elems[i])
	}

	*v.Val = bytess
	return nil
}
