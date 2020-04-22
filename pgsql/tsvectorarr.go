package pgsql

import (
	"database/sql"
	"database/sql/driver"
)

// TSVectorArrayFromStringSliceSlice returns a driver.Valuer that produces a PostgreSQL tsvector[] from the given Go [][]string.
func TSVectorArrayFromStringSliceSlice(val [][]string) driver.Valuer {
	return tsVectorArrayFromStringSliceSlice{val: val}
}

// TSVectorArrayToStringSliceSlice returns an sql.Scanner that converts a PostgreSQL tsvector[] into a Go [][]string and sets it to val.
func TSVectorArrayToStringSliceSlice(val *[][]string) sql.Scanner {
	return tsVectorArrayToStringSliceSlice{val: val}
}

// TSVectorArrayFromByteSliceSliceSlice returns a driver.Valuer that produces a PostgreSQL tsvector[] from the given Go [][][]byte.
func TSVectorArrayFromByteSliceSliceSlice(val [][][]byte) driver.Valuer {
	return tsVectorArrayFromByteSliceSliceSlice{val: val}
}

// TSVectorArrayToByteSliceSliceSlice returns an sql.Scanner that converts a PostgreSQL tsvector[] into a Go [][][]byte and sets it to val.
func TSVectorArrayToByteSliceSliceSlice(val *[][][]byte) sql.Scanner {
	return tsVectorArrayToByteSliceSliceSlice{val: val}
}

type tsVectorArrayFromStringSliceSlice struct {
	val [][]string
}

func (v tsVectorArrayFromStringSliceSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := (2 + (len(v.val) - 1)) // curly braces + number of commas
	for i := 0; i < len(v.val); i++ {
		if v.val[i] == nil {
			size += 4 // len(`NULL`)
		} else if len(v.val[i]) == 0 {
			size += 2 // len(`""`)
		} else {
			size += (2 + (len(v.val[i]) - 1)) // double quotes + number of spaces
			for j := 0; j < len(v.val[i]); j++ {
				size += len(v.val[i][j])
			}
		}
	}

	out := make([]byte, size)
	out[0] = '{'

	var pos int
	for i := 0; i < len(v.val); i++ {
		if v.val[i] == nil {
			out[pos+1] = 'N'
			out[pos+2] = 'U'
			out[pos+3] = 'L'
			out[pos+4] = 'L'
			out[pos+5] = ','
			pos += 5
			continue
		}

		if len(v.val[i]) == 0 {
			out[pos+1] = '"'
			out[pos+2] = '"'
			out[pos+3] = ','
			pos += 3
			continue
		}

		pos += 1
		out[pos] = '"'

		for j, num := 0, len(v.val[i]); j < num; j++ {
			pos += 1

			length := len(v.val[i][j])
			copy(out[pos:pos+length], v.val[i][j])

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

type tsVectorArrayToStringSliceSlice struct {
	val *[][]string
}

func (v tsVectorArrayToStringSliceSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
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

	*v.val = stringss
	return nil
}

type tsVectorArrayFromByteSliceSliceSlice struct {
	val [][][]byte
}

func (v tsVectorArrayFromByteSliceSliceSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := (2 + (len(v.val) - 1)) // curly braces + number of commas
	for i := 0; i < len(v.val); i++ {
		if v.val[i] == nil {
			size += 4 // len(`NULL`)
		} else if len(v.val[i]) == 0 {
			size += 2 // len(`""`)
		} else {
			size += (2 + (len(v.val[i]) - 1)) // double quotes + number of spaces
			for j := 0; j < len(v.val[i]); j++ {
				size += len(v.val[i][j])
			}
		}
	}

	out := make([]byte, size)
	out[0] = '{'

	var pos int
	for i := 0; i < len(v.val); i++ {
		if v.val[i] == nil {
			out[pos+1] = 'N'
			out[pos+2] = 'U'
			out[pos+3] = 'L'
			out[pos+4] = 'L'
			out[pos+5] = ','
			pos += 5
			continue
		}

		if len(v.val[i]) == 0 {
			out[pos+1] = '"'
			out[pos+2] = '"'
			out[pos+3] = ','
			pos += 3
			continue
		}

		pos += 1
		out[pos] = '"'

		for j, num := 0, len(v.val[i]); j < num; j++ {
			pos += 1

			length := len(v.val[i][j])
			copy(out[pos:pos+length], v.val[i][j])

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

type tsVectorArrayToByteSliceSliceSlice struct {
	val *[][][]byte
}

func (v tsVectorArrayToByteSliceSliceSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
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

	*v.val = bytesss
	return nil
}
