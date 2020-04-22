package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// PointArrayFromFloat64Array2Slice returns a driver.Valuer that produces a PostgreSQL point[] from the given Go [][2]float64.
func PointArrayFromFloat64Array2Slice(val [][2]float64) driver.Valuer {
	return pointArrayFromFloat64Array2Slice{val: val}
}

// PointArrayToFloat64Array2Slice returns an sql.Scanner that converts a PostgreSQL point[] into a Go [][2]float64 and sets it to val.
func PointArrayToFloat64Array2Slice(val *[][2]float64) sql.Scanner {
	return pointArrayToFloat64Array2Slice{val: val}
}

type pointArrayFromFloat64Array2Slice struct {
	val [][2]float64
}

func (v pointArrayFromFloat64Array2Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := (len(v.val) * 7) + // len(`"(x,y)"`) == 7
		(len(v.val) - 1) + // number of commas between array elements
		2 // surrounding curly braces

	out := make([]byte, 1, size)
	out[0] = '{'

	for i := 0; i < len(v.val); i++ {
		out = append(out, '"', '(')
		out = strconv.AppendFloat(out, v.val[i][0], 'f', -1, 64)
		out = append(out, ',')
		out = strconv.AppendFloat(out, v.val[i][1], 'f', -1, 64)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}'
	return out, nil
}

type pointArrayToFloat64Array2Slice struct {
	val *[][2]float64
}

func (v pointArrayToFloat64Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	elems := pgParsePointArray(data)
	points := make([][2]float64, len(elems))
	for i := 0; i < len(elems); i++ {
		if points[i][0], err = strconv.ParseFloat(string(elems[i][0]), 64); err != nil {
			return err
		}
		if points[i][1], err = strconv.ParseFloat(string(elems[i][1]), 64); err != nil {
			return err
		}
	}

	*v.val = points
	return nil
}
