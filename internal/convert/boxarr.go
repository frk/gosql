package convert

import (
	"database/sql/driver"
	"strconv"
)

type BoxArrayToFloat64Array2Array2Slice struct {
	S *[][2][2]float64
}

func (s BoxArrayToFloat64Array2Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		s.S = nil
		return nil
	}

	elems := pgparseboxarr(data)

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

	*s.S = boxes
	return nil
}

type BoxArrayFromFloat64Array2Array2Slice struct {
	S [][2][2]float64
}

func (s BoxArrayFromFloat64Array2Array2Slice) Value() (driver.Value, error) {
	if s.S == nil {
		return nil, nil
	} else if len(s.S) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range s.S {
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
