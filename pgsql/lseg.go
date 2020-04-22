package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// LsegFromFloat64Array2Array2 returns a driver.Valuer that produces a PostgreSQL lseg from the given Go [2][2]float64.
func LsegFromFloat64Array2Array2(val [2][2]float64) driver.Valuer {
	return lsegFromFloat64Array2Array2{val: val}
}

// LsegToFloat64Array2Array2 returns an sql.Scanner that converts a PostgreSQL lseg into a Go [2][2]float64 and sets it to val.
func LsegToFloat64Array2Array2(val *[2][2]float64) sql.Scanner {
	return lsegToFloat64Array2Array2{val: val}
}

type lsegFromFloat64Array2Array2 struct {
	val [2][2]float64
}

func (v lsegFromFloat64Array2Array2) Value() (driver.Value, error) {
	out := make([]byte, 2, 13) // len(`[(x,y),(x,y)]`) == 13 (min size)
	out[0], out[1] = '[', '('

	out = strconv.AppendFloat(out, v.val[0][0], 'f', -1, 64)
	out = append(out, ',')
	out = strconv.AppendFloat(out, v.val[0][1], 'f', -1, 64)
	out = append(out, ')', ',', '(')
	out = strconv.AppendFloat(out, v.val[1][0], 'f', -1, 64)
	out = append(out, ',')
	out = strconv.AppendFloat(out, v.val[1][1], 'f', -1, 64)
	out = append(out, ')', ']')
	return out, nil
}

type lsegToFloat64Array2Array2 struct {
	val *[2][2]float64
}

func (v lsegToFloat64Array2Array2) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	lseg := pgParseLseg(data)
	f0, err := strconv.ParseFloat(string(lseg[0][0]), 64)
	if err != nil {
		return err
	}
	f1, err := strconv.ParseFloat(string(lseg[0][1]), 64)
	if err != nil {
		return err
	}
	f2, err := strconv.ParseFloat(string(lseg[1][0]), 64)
	if err != nil {
		return err
	}
	f3, err := strconv.ParseFloat(string(lseg[1][1]), 64)
	if err != nil {
		return err
	}

	*v.val = [2][2]float64{{f0, f1}, {f2, f3}}
	return nil
}
