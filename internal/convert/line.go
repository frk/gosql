package convert

import (
	"database/sql/driver"
	"strconv"
)

type LineFromFloat64Array3 struct {
	Val [3]float64
}

func (v LineFromFloat64Array3) Value() (driver.Value, error) {
	out := make([]byte, 1, 7) // len(`{a,b,c}`) == 7 (min size)
	out[0] = '{'

	out = strconv.AppendFloat(out, v.Val[0], 'f', -1, 64)
	out = append(out, ',')
	out = strconv.AppendFloat(out, v.Val[1], 'f', -1, 64)
	out = append(out, ',')
	out = strconv.AppendFloat(out, v.Val[2], 'f', -1, 64)
	out = append(out, '}')
	return out, nil
}

type LineToFloat64Array3 struct {
	Val *[3]float64
}

func (v LineToFloat64Array3) Scan(src interface{}) error {
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

	*v.Val = [3]float64{f0, f1, f2}
	return nil
}
