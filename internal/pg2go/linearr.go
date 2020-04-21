package pg2go

import (
	"database/sql/driver"
	"strconv"
)

type LineArrayFromFloat64Array3Slice struct {
	Val [][3]float64
}

func (v LineArrayFromFloat64Array3Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := (len(v.Val) * 9) + // len(`"{a,b,c}"`) == 9
		(len(v.Val) - 1) + // number of commas between array elements
		2 // surrounding curly braces

	out := make([]byte, 1, size)
	out[0] = '{'

	for i := 0; i < len(v.Val); i++ {
		out = append(out, '"', '{')
		out = strconv.AppendFloat(out, v.Val[i][0], 'f', -1, 64)
		out = append(out, ',')
		out = strconv.AppendFloat(out, v.Val[i][1], 'f', -1, 64)
		out = append(out, ',')
		out = strconv.AppendFloat(out, v.Val[i][2], 'f', -1, 64)
		out = append(out, '}', '"', ',')
	}

	out[len(out)-1] = '}'
	return out, nil
}

type LineArrayToFloat64Array3Slice struct {
	Val *[][3]float64
}

func (v LineArrayToFloat64Array3Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	elems := pgParseLineArray(data)
	lines := make([][3]float64, len(elems))
	for i := 0; i < len(elems); i++ {
		if lines[i][0], err = strconv.ParseFloat(string(elems[i][0]), 64); err != nil {
			return err
		}
		if lines[i][1], err = strconv.ParseFloat(string(elems[i][1]), 64); err != nil {
			return err
		}
		if lines[i][2], err = strconv.ParseFloat(string(elems[i][2]), 64); err != nil {
			return err
		}
	}

	*v.Val = lines
	return nil
}
