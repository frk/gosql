package convert

import (
	"database/sql/driver"
	"strconv"
)

type PointFromFloat64Array2 struct {
	Val [2]float64
}

func (v PointFromFloat64Array2) Value() (driver.Value, error) {
	out := make([]byte, 1, 5) // len(`(x,y)`) == 5 (min size)
	out[0] = '('

	out = strconv.AppendFloat(out, v.Val[0], 'f', -1, 64)
	out = append(out, ',')
	out = strconv.AppendFloat(out, v.Val[1], 'f', -1, 64)
	out = append(out, ')')
	return out, nil
}

type PointToFloat64Array2 struct {
	Val *[2]float64
}

func (v PointToFloat64Array2) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	point := pgParsePoint(data)
	f0, err := strconv.ParseFloat(string(point[0]), 64)
	if err != nil {
		return err
	}
	f1, err := strconv.ParseFloat(string(point[1]), 64)
	if err != nil {
		return err
	}

	*v.Val = [2]float64{f0, f1}
	return nil
}
