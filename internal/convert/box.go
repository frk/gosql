package convert

import (
	"strconv"
)

type Box2Float64a2a2 struct {
	Ptr *[2][2]float64
}

func (s Box2Float64a2a2) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	elems := pgparsebox(data)

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

	*s.Ptr = box
	return nil
}
