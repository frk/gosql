package convert

import (
	"strconv"
)

type BoxArr2Float64a2a2Slice struct {
	Ptr *[][2][2]float64
}

func (s BoxArr2Float64a2a2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		s.Ptr = nil
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

	*s.Ptr = boxes
	return nil
}
