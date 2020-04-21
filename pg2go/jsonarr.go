package pg2go

import (
	"database/sql"
	"database/sql/driver"
)

// JSONArrayFromByteSliceSlice returns a driver.Valuer that produces a PostgreSQL json(b)[] from the given Go [][]byte.
func JSONArrayFromByteSliceSlice(val [][]byte) driver.Valuer {
	return jsonArrayFromByteSliceSlice{val: val}
}

// JSONArrayToByteSliceSlice returns an sql.Scanner that converts a PostgreSQL json(b)[] into a Go [][]byte and sets it to val.
func JSONArrayToByteSliceSlice(val *[][]byte) sql.Scanner {
	return jsonArrayToByteSliceSlice{val: val}
}

type jsonArrayFromByteSliceSlice struct {
	val [][]byte
}

func (v jsonArrayFromByteSliceSlice) Value() (driver.Value, error) {
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

type jsonArrayToByteSliceSlice struct {
	val *[][]byte
}

func (v jsonArrayToByteSliceSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
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

	*v.val = jsons
	return nil
}
