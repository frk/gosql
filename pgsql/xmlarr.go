package pgsql

import (
	"database/sql"
	"database/sql/driver"
)

// XMLArrayFromByteSliceSlice returns a driver.Valuer that produces a PostgreSQL xml[] from the given Go [][]byte.
func XMLArrayFromByteSliceSlice(val [][]byte) driver.Valuer {
	return xmlArrayFromByteSliceSlice{val: val}
}

// XMLArrayToByteSliceSlice returns an sql.Scanner that converts a PostgreSQL xml[] into a Go [][]byte and sets it to val.
func XMLArrayToByteSliceSlice(val *[][]byte) sql.Scanner {
	return xmlArrayToByteSliceSlice{val: val}
}

type xmlArrayFromByteSliceSlice struct {
	val [][]byte
}

func (v xmlArrayFromByteSliceSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.val {
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

type xmlArrayToByteSliceSlice struct {
	val *[][]byte
}

func (v xmlArrayToByteSliceSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
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

	*v.val = xmls
	return nil
}
