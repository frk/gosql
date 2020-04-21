package pg2go

import (
	"database/sql/driver"
	"strconv"
)

type BoxFromFloat64Array2Array2 struct {
	Val [2][2]float64
}

func (s BoxFromFloat64Array2Array2) Value() (driver.Value, error) {
	out := []byte{'(', '('}

	out = strconv.AppendFloat(out, s.Val[0][0], 'f', -1, 64)
	out = append(out, ',')
	out = strconv.AppendFloat(out, s.Val[0][1], 'f', -1, 64)

	out = append(out, ')', ',', '(')

	out = strconv.AppendFloat(out, s.Val[1][0], 'f', -1, 64)
	out = append(out, ',')
	out = strconv.AppendFloat(out, s.Val[1][1], 'f', -1, 64)

	return append(out, ')', ')'), nil
}

type BoxToFloat64Array2Array2 struct {
	Val *[2][2]float64
}

func (s BoxToFloat64Array2Array2) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	elems := pgParseBox(data)

	var box [2][2]float64
	x1, err := strconv.ParseFloat(string(elems[0]), 64)
	if err != nil {
		return err
	}
	y1, err := strconv.ParseFloat(string(elems[1]), 64)
	if err != nil {
		return err
	}
	x2, err := strconv.ParseFloat(string(elems[2]), 64)
	if err != nil {
		return err
	}
	y2, err := strconv.ParseFloat(string(elems[3]), 64)
	if err != nil {
		return err
	}

	box[0][0] = x1
	box[0][1] = y1
	box[1][0] = x2
	box[1][1] = y2

	*s.Val = box
	return nil
}
