package pgsql

import (
	"database/sql"
	"database/sql/driver"
)

// UUIDArrayFromByteArray16Slice returns a driver.Valuer that produces a PostgreSQL uuid[] from the given Go [][16]byte.
func UUIDArrayFromByteArray16Slice(val [][16]byte) driver.Valuer {
	return uuidArrayFromByteArray16Slice{val: val}
}

// UUIDArrayToByteArray16Slice returns an sql.Scanner that converts a PostgreSQL uuid[] into a Go [][16]byte and sets it to val.
func UUIDArrayToByteArray16Slice(val *[][16]byte) sql.Scanner {
	return uuidArrayToByteArray16Slice{val: val}
}

// UUIDArrayFromStringSlice returns a driver.Valuer that produces a PostgreSQL uuid[] from the given Go []string.
func UUIDArrayFromStringSlice(val []string) driver.Valuer {
	return uuidArrayFromStringSlice{val: val}
}

// UUIDArrayToStringSlice returns an sql.Scanner that converts a PostgreSQL uuid[] into a Go []string and sets it to val.
func UUIDArrayToStringSlice(val *[]string) sql.Scanner {
	return uuidArrayToStringSlice{val: val}
}

// UUIDArrayFromByteSliceSlice returns a driver.Valuer that produces a PostgreSQL uuid[] from the given Go [][]byte.
func UUIDArrayFromByteSliceSlice(val [][]byte) driver.Valuer {
	return uuidArrayFromByteSliceSlice{val: val}
}

// UUIDArrayToByteSliceSlice returns an sql.Scanner that converts a PostgreSQL uuid[] into a Go [][]byte and sets it to val.
func UUIDArrayToByteSliceSlice(val *[][]byte) sql.Scanner {
	return uuidArrayToByteSliceSlice{val: val}
}

type uuidArrayFromByteArray16Slice struct {
	val [][16]byte
}

func (v uuidArrayFromByteArray16Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := 1 + len(v.val) + (len(v.val) * 36)

	out := make([]byte, size)
	out[0] = '{'

	var pos int
	for i := 0; i < len(v.val); i++ {
		pos += 1

		uuid := pgFormatUUID(v.val[i])
		copy(out[pos:pos+36], uuid)

		pos += 36
		out[pos] = ','
	}

	out[pos] = '}'
	return out, nil
}

type uuidArrayToByteArray16Slice struct {
	val *[][16]byte
}

func (v uuidArrayToByteArray16Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
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

	*v.val = uuids
	return nil
}

type uuidArrayFromStringSlice struct {
	val []string
}

func (v uuidArrayFromStringSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := 2 + (len(v.val) - 1)
	for i := 0; i < len(v.val); i++ {
		size += len(v.val[i])
	}

	out := make([]byte, size)
	out[0] = '{'

	var pos int
	for i := 0; i < len(v.val); i++ {
		pos += 1

		length := len(v.val[i])
		copy(out[pos:pos+length], v.val[i])

		pos += length
		out[pos] = ','
	}

	out[pos] = '}'
	return out, nil
}

type uuidArrayToStringSlice struct {
	val *[]string
}

func (v uuidArrayToStringSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseCommaArray(data)
	strings := make([]string, len(elems))
	for i := 0; i < len(elems); i++ {
		strings[i] = string(elems[i])
	}

	*v.val = strings
	return nil
}

type uuidArrayFromByteSliceSlice struct {
	val [][]byte
}

func (v uuidArrayFromByteSliceSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := 2 + (len(v.val) - 1)
	for i := 0; i < len(v.val); i++ {
		if v.val[i] == nil {
			size += 4 // len("NULL")
		} else {
			size += len(v.val[i])
		}
	}

	out := make([]byte, size)
	out[0] = '{'

	var pos int
	for i := 0; i < len(v.val); i++ {
		pos += 1

		if v.val[i] == nil {
			out[pos+0] = 'N'
			out[pos+1] = 'U'
			out[pos+2] = 'L'
			out[pos+3] = 'L'
			pos += 4
		} else {
			length := len(v.val[i])
			copy(out[pos:pos+length], v.val[i])
			pos += length
		}

		out[pos] = ','
	}

	out[pos] = '}'
	return out, nil
}

type uuidArrayToByteSliceSlice struct {
	val *[][]byte
}

func (v uuidArrayToByteSliceSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
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

	*v.val = bytess
	return nil
}
