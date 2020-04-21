package pg2go

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// BoxFromFloat64Array2Array2 returns a driver.Valuer that produces a PostgreSQL box from the given Go [2][2]float64.
func BoxFromFloat64Array2Array2(val [2][2]float64) driver.Valuer {
	return boxFromFloat64Array2Array2{val: val}
}

// BoxToFloat64Array2Array2 returns an sql.Scanner that converts a PostgreSQL box into a Go [2][2]float64 and sets it to val.
func BoxToFloat64Array2Array2(val *[2][2]float64) sql.Scanner {
	return boxToFloat64Array2Array2{val: val}
}

type boxFromFloat64Array2Array2 struct {
	val [2][2]float64
}

func (v boxFromFloat64Array2Array2) Value() (driver.Value, error) {
	out := []byte{'(', '('}

	out = strconv.AppendFloat(out, v.val[0][0], 'f', -1, 64)
	out = append(out, ',')
	out = strconv.AppendFloat(out, v.val[0][1], 'f', -1, 64)

	out = append(out, ')', ',', '(')

	out = strconv.AppendFloat(out, v.val[1][0], 'f', -1, 64)
	out = append(out, ',')
	out = strconv.AppendFloat(out, v.val[1][1], 'f', -1, 64)

	return append(out, ')', ')'), nil
}

type boxToFloat64Array2Array2 struct {
	val *[2][2]float64
}

func (v boxToFloat64Array2Array2) Scan(src interface{}) error {
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

	*v.val = box
	return nil
}
