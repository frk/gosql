package pg2go

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// LineArrayFromFloat64Array3Slice returns a driver.Valuer that produces a PostgreSQL line[] from the given Go [][3]float64.
func LineArrayFromFloat64Array3Slice(val [][3]float64) driver.Valuer {
	return lineArrayFromFloat64Array3Slice{val: val}
}

// LineArrayToFloat64Array3Slice returns an sql.Scanner that converts a PostgreSQL line[] into a Go [][3]float64 and sets it to val.
func LineArrayToFloat64Array3Slice(val *[][3]float64) sql.Scanner {
	return lineArrayToFloat64Array3Slice{val: val}
}

type lineArrayFromFloat64Array3Slice struct {
	val [][3]float64
}

func (v lineArrayFromFloat64Array3Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := (len(v.val) * 9) + // len(`"{a,b,c}"`) == 9
		(len(v.val) - 1) + // number of commas between array elements
		2 // surrounding curly braces

	out := make([]byte, 1, size)
	out[0] = '{'

	for i := 0; i < len(v.val); i++ {
		out = append(out, '"', '{')
		out = strconv.AppendFloat(out, v.val[i][0], 'f', -1, 64)
		out = append(out, ',')
		out = strconv.AppendFloat(out, v.val[i][1], 'f', -1, 64)
		out = append(out, ',')
		out = strconv.AppendFloat(out, v.val[i][2], 'f', -1, 64)
		out = append(out, '}', '"', ',')
	}

	out[len(out)-1] = '}'
	return out, nil
}

type lineArrayToFloat64Array3Slice struct {
	val *[][3]float64
}

func (v lineArrayToFloat64Array3Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	elems := pgParseLineArray(data)
	lines := make([][3]float64, len(elems))
	for i := 0; i < len(elems); i++ {
		if lines[i][0], err = strconv.ParseFloat(string(elems[i][0]), 64); err != nil {
			return err
		}
		if lines[i][1], err = strconv.ParseFloat(string(elems[i][1]), 64); err != nil {
			return err
		}
		if lines[i][2], err = strconv.ParseFloat(string(elems[i][2]), 64); err != nil {
			return err
		}
	}

	*v.val = lines
	return nil
}
