package pg2go

import (
	"database/sql/driver"
)

type UUIDArrayFromByteArray16Slice struct {
	Val [][16]byte
}

func (v UUIDArrayFromByteArray16Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := 1 + len(v.Val) + (len(v.Val) * 36)

	out := make([]byte, size)
	out[0] = '{'

	var pos int
	for i := 0; i < len(v.Val); i++ {
		pos += 1

		uuid := pgFormatUUID(v.Val[i])
		copy(out[pos:pos+36], uuid)

		pos += 36
		out[pos] = ','
	}

	out[pos] = '}'
	return out, nil
}

type UUIDArrayToByteArray16Slice struct {
	Val *[][16]byte
}

func (v UUIDArrayToByteArray16Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseCommaArray(data)
	uuids := make([][16]byte, len(elems))
	for i := 0; i < len(elems); i++ {
		uuid, err := pgParseUUID(elems[i])
		if err != nil {
			return err
		}
		uuids[i] = uuid
	}

	*v.Val = uuids
	return nil
}

type UUIDArrayFromStringSlice struct {
	Val []string
}

func (v UUIDArrayFromStringSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := 2 + (len(v.Val) - 1)
	for i := 0; i < len(v.Val); i++ {
		size += len(v.Val[i])
	}

	out := make([]byte, size)
	out[0] = '{'

	var pos int
	for i := 0; i < len(v.Val); i++ {
		pos += 1

		length := len(v.Val[i])
		copy(out[pos:pos+length], v.Val[i])

		pos += length
		out[pos] = ','
	}

	out[pos] = '}'
	return out, nil
}

type UUIDArrayToStringSlice struct {
	Val *[]string
}

func (v UUIDArrayToStringSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseCommaArray(data)
	strings := make([]string, len(elems))
	for i := 0; i < len(elems); i++ {
		strings[i] = string(elems[i])
	}

	*v.Val = strings
	return nil
}

type UUIDArrayFromByteSliceSlice struct {
	Val [][]byte
}

func (v UUIDArrayFromByteSliceSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := 2 + (len(v.Val) - 1)
	for i := 0; i < len(v.Val); i++ {
		if v.Val[i] == nil {
			size += 4 // len("NULL")
		} else {
			size += len(v.Val[i])
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
			length := len(v.Val[i])
			copy(out[pos:pos+length], v.Val[i])
			pos += length
		}

		out[pos] = ','
	}

	out[pos] = '}'
	return out, nil
}

type UUIDArrayToByteSliceSlice struct {
	Val *[][]byte
}

func (v UUIDArrayToByteSliceSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseCommaArray(data)
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
