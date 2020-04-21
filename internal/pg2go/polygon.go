package pg2go

import (
	"database/sql/driver"
	"strconv"
)

type PolygonFromFloat64Array2Slice struct {
	Val [][2]float64
}

func (v PolygonFromFloat64Array2Slice) Value() (driver.Value, error) {
	if v.Val == nil || len(v.Val) == 0 {
		return nil, nil
	}

	size := (len(v.Val) * 5) + // len(`(x,y)`) == 5
		(len(v.Val) - 1) + // number of commas between points
		2 // surrounding parentheses

	out := make([]byte, 1, size)
	out[0] = '('

	for i := 0; i < len(v.Val); i++ {
		out = append(out, '(')
		out = strconv.AppendFloat(out, v.Val[i][0], 'f', -1, 64)
		out = append(out, ',')
		out = strconv.AppendFloat(out, v.Val[i][1], 'f', -1, 64)
		out = append(out, ')', ',')
	}

	out[len(out)-1] = ')'
	return out, nil
}

type PolygonToFloat64Array2Slice struct {
	Val *[][2]float64
}

func (v PolygonToFloat64Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.Val = nil
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

	*v.Val = points
	return nil
}
