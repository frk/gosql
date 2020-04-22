package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// PointFromFloat64Array2 returns a driver.Valuer that produces a PostgreSQL point from the given Go [2]float64.
func PointFromFloat64Array2(val [2]float64) driver.Valuer {
	return pointFromFloat64Array2{val: val}
}

// PointToFloat64Array2 returns an sql.Scanner that converts a PostgreSQL point into a Go [2]float64 and sets it to val.
func PointToFloat64Array2(val *[2]float64) sql.Scanner {
	return pointToFloat64Array2{val: val}
}

type pointFromFloat64Array2 struct {
	val [2]float64
}

func (v pointFromFloat64Array2) Value() (driver.Value, error) {
	out := make([]byte, 1, 5) // len(`(x,y)`) == 5 (min size)
	out[0] = '('

	out = strconv.AppendFloat(out, v.val[0], 'f', -1, 64)
	out = append(out, ',')
	out = strconv.AppendFloat(out, v.val[1], 'f', -1, 64)
	out = append(out, ')')
	return out, nil
}

type pointToFloat64Array2 struct {
	val *[2]float64
}

func (v pointToFloat64Array2) Scan(src interface{}) error {
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

	*v.val = [2]float64{f0, f1}
	return nil
}
