package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// [ ( x1 , y1 ) , ... , ( xn , yn ) ]
// ( ( x1 , y1 ) , ... , ( xn , yn ) )
//
// Square brackets "[]" indicate an open path, while parentheses "()" indicate
// a closed path.
//
// - will require struct tag option to specify whether the path should be closed
//   or open, this unfortunately will make it a "static" path type

// PathFromFloat64Array2Slice returns a driver.Valuer that produces a PostgreSQL path from the given Go [][2]float64.
func PathFromFloat64Array2Slice(val [][2]float64) driver.Valuer {
	return pathFromFloat64Array2Slice{val: val}
}

// PathToFloat64Array2Slice returns an sql.Scanner that converts a PostgreSQL path into a Go [][2]float64 and sets it to val.
func PathToFloat64Array2Slice(val *[][2]float64) sql.Scanner {
	return pathToFloat64Array2Slice{val: val}
}

type pathFromFloat64Array2Slice struct {
	val [][2]float64
}

func (v pathFromFloat64Array2Slice) Value() (driver.Value, error) {
	if v.val == nil || len(v.val) == 0 {
		return nil, nil
	}

	size := (len(v.val) * 5) + // len(`(x,y)`) == 5
		(len(v.val) - 1) + // number of commas between points
		2 // surrounding parentheses

	out := make([]byte, 1, size)
	out[0] = '['

	for i := 0; i < len(v.val); i++ {
		out = append(out, '(')
		out = strconv.AppendFloat(out, v.val[i][0], 'f', -1, 64)
		out = append(out, ',')
		out = strconv.AppendFloat(out, v.val[i][1], 'f', -1, 64)
		out = append(out, ')', ',')
	}

	out[len(out)-1] = ']'
	return out, nil
}

type pathToFloat64Array2Slice struct {
	val *[][2]float64
}

func (v pathToFloat64Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
		return nil
	}

	elems, _ := pgParsePath(data)
	points := make([][2]float64, len(elems))
	for i := 0; i < len(elems); i++ {
		f0, err := strconv.ParseFloat(string(elems[i][0]), 64)
		if err != nil {
			return err
		}
		f1, err := strconv.ParseFloat(string(elems[i][1]), 64)
		if err != nil {
			return err
		}

		points[i][0] = f0
		points[i][1] = f1
	}

	*v.val = points
	return nil
}
