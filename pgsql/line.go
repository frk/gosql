package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// LineFromFloat64Array3 returns a driver.Valuer that produces a PostgreSQL line from the given Go [3]float64.
func LineFromFloat64Array3(val [3]float64) driver.Valuer {
	return lineFromFloat64Array3{val: val}
}

// LineToFloat64Array3 returns an sql.Scanner that converts a PostgreSQL line into a Go [3]float64 and sets it to val.
func LineToFloat64Array3(val *[3]float64) sql.Scanner {
	return lineToFloat64Array3{val: val}
}

type lineFromFloat64Array3 struct {
	val [3]float64
}

func (v lineFromFloat64Array3) Value() (driver.Value, error) {
	out := make([]byte, 1, 7) // len(`{a,b,c}`) == 7 (min size)
	out[0] = '{'

	out = strconv.AppendFloat(out, v.val[0], 'f', -1, 64)
	out = append(out, ',')
	out = strconv.AppendFloat(out, v.val[1], 'f', -1, 64)
	out = append(out, ',')
	out = strconv.AppendFloat(out, v.val[2], 'f', -1, 64)
	out = append(out, '}')
	return out, nil
}

type lineToFloat64Array3 struct {
	val *[3]float64
}

func (v lineToFloat64Array3) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	line := pgParseLine(data)
	f0, err := strconv.ParseFloat(string(line[0]), 64)
	if err != nil {
		return err
	}
	f1, err := strconv.ParseFloat(string(line[1]), 64)
	if err != nil {
		return err
	}
	f2, err := strconv.ParseFloat(string(line[2]), 64)
	if err != nil {
		return err
	}

	*v.val = [3]float64{f0, f1, f2}
	return nil
}
