package pg2go

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// PolygonFromFloat64Array2Slice returns a driver.Valuer that produces a PostgreSQL polygon from the given Go [][2]float64.
func PolygonFromFloat64Array2Slice(val [][2]float64) driver.Valuer {
	return polygonFromFloat64Array2Slice{val: val}
}

// PolygonToFloat64Array2Slice returns an sql.Scanner that converts a PostgreSQL polygon into a Go [][2]float64 and sets it to val.
func PolygonToFloat64Array2Slice(val *[][2]float64) sql.Scanner {
	return polygonToFloat64Array2Slice{val: val}
}

type polygonFromFloat64Array2Slice struct {
	val [][2]float64
}

func (v polygonFromFloat64Array2Slice) Value() (driver.Value, error) {
	if v.val == nil || len(v.val) == 0 {
		return nil, nil
	}

	size := (len(v.val) * 5) + // len(`(x,y)`) == 5
		(len(v.val) - 1) + // number of commas between points
		2 // surrounding parentheses

	out := make([]byte, 1, size)
	out[0] = '('

	for i := 0; i < len(v.val); i++ {
		out = append(out, '(')
		out = strconv.AppendFloat(out, v.val[i][0], 'f', -1, 64)
		out = append(out, ',')
		out = strconv.AppendFloat(out, v.val[i][1], 'f', -1, 64)
		out = append(out, ')', ',')
	}

	out[len(out)-1] = ')'
	return out, nil
}

type polygonToFloat64Array2Slice struct {
	val *[][2]float64
}

func (v polygonToFloat64Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
		return nil
	}

	elems := pgParsePolygon(data)
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
