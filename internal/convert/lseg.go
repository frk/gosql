package convert

import (
	"database/sql/driver"
	"strconv"
)

type LsegFromFloat64Array2Array2 struct {
	Val [2][2]float64
}

func (v LsegFromFloat64Array2Array2) Value() (driver.Value, error) {
	out := make([]byte, 2, 13) // len(`[(x,y),(x,y)]`) == 13 (min size)
	out[0], out[1] = '[', '('

	out = strconv.AppendFloat(out, v.Val[0][0], 'f', -1, 64)
	out = append(out, ',')
	out = strconv.AppendFloat(out, v.Val[0][1], 'f', -1, 64)
	out = append(out, ')', ',', '(')
	out = strconv.AppendFloat(out, v.Val[1][0], 'f', -1, 64)
	out = append(out, ',')
	out = strconv.AppendFloat(out, v.Val[1][1], 'f', -1, 64)
	out = append(out, ')', ']')
	return out, nil
}

type LsegToFloat64Array2Array2 struct {
	Val *[2][2]float64
}

func (v LsegToFloat64Array2Array2) Scan(src interface{}) error {
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

	*v.Val = [2][2]float64{{f0, f1}, {f2, f3}}
	return nil
}
