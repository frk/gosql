package convert

import (
	"database/sql/driver"
)

type TSVectorArrayFromStringSliceSlice struct {
	Val [][]string
}

func (v TSVectorArrayFromStringSliceSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := (2 + (len(v.Val) - 1)) // curly braces + number of commas
	for i := 0; i < len(v.Val); i++ {
		if v.Val[i] == nil {
			size += 4 // len(`NULL`)
		} else if len(v.Val[i]) == 0 {
			size += 2 // len(`""`)
		} else {
			size += (2 + (len(v.Val[i]) - 1)) // double quotes + number of spaces
			for j := 0; j < len(v.Val[i]); j++ {
				size += len(v.Val[i][j])
			}
		}
	}

	out := make([]byte, size)
	out[0] = '{'

	var pos int
	for i := 0; i < len(v.Val); i++ {
		if v.Val[i] == nil {
			out[pos+1] = 'N'
			out[pos+2] = 'U'
			out[pos+3] = 'L'
			out[pos+4] = 'L'
			out[pos+5] = ','
			pos += 5
			continue
		}

		if len(v.Val[i]) == 0 {
			out[pos+1] = '"'
			out[pos+2] = '"'
			out[pos+3] = ','
			pos += 3
			continue
		}

		pos += 1
		out[pos] = '"'

		for j, num := 0, len(v.Val[i]); j < num; j++ {
			pos += 1

			length := len(v.Val[i][j])
			copy(out[pos:pos+length], v.Val[i][j])

			pos += length
			out[pos] = ' '
		}

		out[pos] = '"'
		pos += 1

		out[pos] = ','
	}

	out[pos] = '}'
	return out, nil
}

type TSVectorArrayToStringSliceSlice struct {
	Val *[][]string
}

func (v TSVectorArrayToStringSliceSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVectorArray(data)
	stringss := make([][]string, len(elems))
	for i := 0; i < len(elems); i++ {
		if elems[i] == nil {
			continue
		}
		if len(elems[i]) == 1 && len(elems[i][0]) == 0 {
			stringss[i] = []string{}
			continue
		}

		strings := make([]string, len(elems[i]))
		for j := 0; j < len(elems[i]); j++ {
			strings[j] = string(elems[i][j])
		}
		stringss[i] = strings
	}

	*v.Val = stringss
	return nil
}

type TSVectorArrayFromByteSliceSliceSlice struct {
	Val [][][]byte
}

func (v TSVectorArrayFromByteSliceSliceSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := (2 + (len(v.Val) - 1)) // curly braces + number of commas
	for i := 0; i < len(v.Val); i++ {
		if v.Val[i] == nil {
			size += 4 // len(`NULL`)
		} else if len(v.Val[i]) == 0 {
			size += 2 // len(`""`)
		} else {
			size += (2 + (len(v.Val[i]) - 1)) // double quotes + number of spaces
			for j := 0; j < len(v.Val[i]); j++ {
				size += len(v.Val[i][j])
			}
		}
	}

	out := make([]byte, size)
	out[0] = '{'

	var pos int
	for i := 0; i < len(v.Val); i++ {
		if v.Val[i] == nil {
			out[pos+1] = 'N'
			out[pos+2] = 'U'
			out[pos+3] = 'L'
			out[pos+4] = 'L'
			out[pos+5] = ','
			pos += 5
			continue
		}

		if len(v.Val[i]) == 0 {
			out[pos+1] = '"'
			out[pos+2] = '"'
			out[pos+3] = ','
			pos += 3
			continue
		}

		pos += 1
		out[pos] = '"'

		for j, num := 0, len(v.Val[i]); j < num; j++ {
			pos += 1

			length := len(v.Val[i][j])
			copy(out[pos:pos+length], v.Val[i][j])

			pos += length
			out[pos] = ' '
		}

		out[pos] = '"'
		pos += 1

		out[pos] = ','
	}

	out[pos] = '}'
	return out, nil
}

type TSVectorArrayToByteSliceSliceSlice struct {
	Val *[][][]byte
}

func (v TSVectorArrayToByteSliceSliceSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVectorArray(data)
	bytesss := make([][][]byte, len(elems))
	for i := 0; i < len(elems); i++ {
		if elems[i] == nil {
			continue
		}
		if len(elems[i]) == 1 && len(elems[i][0]) == 0 {
			bytesss[i] = [][]byte{}
			continue
		}

		bytess := make([][]byte, len(elems[i]))
		for j := 0; j < len(elems[i]); j++ {
			bytess[j] = make([]byte, len(elems[i][j]))
			copy(bytess[j], elems[i][j])
		}
		bytesss[i] = bytess
	}

	*v.Val = bytesss
	return nil
}
