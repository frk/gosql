package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// BoxArrayFromFloat64Array2Array2Slice returns a driver.Valuer that produces a PostgreSQL box[] from the given Go [][2][2]float64.
func BoxArrayFromFloat64Array2Array2Slice(val [][2][2]float64) driver.Valuer {
	return boxArrayFromFloat64Array2Array2Slice{val: val}
}

// BoxArrayToFloat64Array2Array2Slice returns an sql.Scanner that converts a PostgreSQL box[] into a Go [][2][2]float64 and sets it to val.
func BoxArrayToFloat64Array2Array2Slice(val *[][2][2]float64) sql.Scanner {
	return boxArrayToFloat64Array2Array2Slice{val: val}
}

type boxArrayFromFloat64Array2Array2Slice struct {
	val [][2][2]float64
}

func (v boxArrayFromFloat64Array2Array2Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.val {
		out = append(out, '(')
		out = strconv.AppendFloat(out, a[0][0], 'f', -1, 64)
		out = append(out, ',')
		out = strconv.AppendFloat(out, a[0][1], 'f', -1, 64)

		out = append(out, ')', ',', '(')

		out = strconv.AppendFloat(out, a[1][0], 'f', -1, 64)
		out = append(out, ',')
		out = strconv.AppendFloat(out, a[1][1], 'f', -1, 64)
		out = append(out, ')', ';')
	}

	out[len(out)-1] = '}' // replace last ";" with "}"
	return out, nil
}

type boxArrayToFloat64Array2Array2Slice struct {
	val *[][2][2]float64
}

func (v boxArrayToFloat64Array2Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseBoxArray(data)

	n := len(elems) / 4
	boxes := make([][2][2]float64, n)
	for i := 0; i < n; i++ {
		j := i * 4

		x1, err := strconv.ParseFloat(string(elems[j]), 64)
		if err != nil {
			return err
		}
		y1, err := strconv.ParseFloat(string(elems[j+1]), 64)
		if err != nil {
			return err
		}
		x2, err := strconv.ParseFloat(string(elems[j+2]), 64)
		if err != nil {
			return err
		}
		y2, err := strconv.ParseFloat(string(elems[j+3]), 64)
		if err != nil {
			return err
		}

		boxes[i][0][0] = x1
		boxes[i][0][1] = y1
		boxes[i][1][0] = x2
		boxes[i][1][1] = y2
	}

	*v.val = boxes
	return nil
}
