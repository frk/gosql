package pg2go

import (
	"database/sql/driver"
)

type TSVectorFromStringSlice struct {
	Val []string
}

func (v TSVectorFromStringSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{}, nil
	}

	size := len(v.Val) - 1
	for i := 0; i < len(v.Val); i++ {
		size += len(v.Val[i])
	}
	out := make([]byte, size)

	var pos int
	for i := 0; i < len(v.Val); i++ {
		length := len(v.Val[i])
		copy(out[pos:pos+length], v.Val[i])
		pos += length

		if pos < size {
			out[pos] = ' '
			pos += 1
		}
	}
	return out, nil
}

type TSVectorToStringSlice struct {
	Val *[]string
}

func (v TSVectorToStringSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVector(data)
	strings := make([]string, len(elems))
	for i := 0; i < len(elems); i++ {
		strings[i] = string(elems[i])
	}

	*v.Val = strings
	return nil
}

type TSVectorFromByteSliceSlice struct {
	Val [][]byte
}

func (v TSVectorFromByteSliceSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{}, nil
	}

	size := len(v.Val) - 1
	for i := 0; i < len(v.Val); i++ {
		size += len(v.Val[i])
	}
	out := make([]byte, size)

	var pos int
	for i := 0; i < len(v.Val); i++ {
		length := len(v.Val[i])
		copy(out[pos:pos+length], v.Val[i])
		pos += length

		if pos < size {
			out[pos] = ' '
			pos += 1
		}
	}
	return out, nil
}

type TSVectorToByteSliceSlice struct {
	Val *[][]byte
}

func (v TSVectorToByteSliceSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVector(data)
	bytess := make([][]byte, len(elems))
	for i := 0; i < len(elems); i++ {
		bytess[i] = make([]byte, len(elems[i]))
		copy(bytess[i], elems[i])
	}

	*v.Val = bytess
	return nil
}
