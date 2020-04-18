package convert

import (
	"database/sql/driver"
	"strconv"
)

type PathArrayFromFloat64Array2SliceSlice struct {
	Val [][][2]float64
}

func (v PathArrayFromFloat64Array2SliceSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := (len(v.Val) * 4) + // len(`"[]"`) == len(`NULL`) == 4
		(len(v.Val) - 1) + // number of commas between elements
		2 // surrounding parentheses
	for i := 0; i < len(v.Val); i++ {
		if v.Val[i] != nil {
			size += (len(v.Val[i]) * 5) + // len(`(x,y)`)
				(len(v.Val[i]) - 1) // number of commas between points
		}

	}

	out := make([]byte, 1, size)
	out[0] = '{'

	for i := 0; i < len(v.Val); i++ {
		if v.Val[i] == nil {
			out = append(out, 'N', 'U', 'L', 'L', ',')
			continue
		}

		out = append(out, '"', '[')

		for j := 0; j < len(v.Val[i]); j++ {
			out = append(out, '(')
			out = strconv.AppendFloat(out, v.Val[i][j][0], 'f', -1, 64)
			out = append(out, ',')
			out = strconv.AppendFloat(out, v.Val[i][j][1], 'f', -1, 64)
			out = append(out, ')', ',')
		}

		out[len(out)-1] = ']'
		out = append(out, '"', ',')
	}

	out[len(out)-1] = '}'
	return out, nil
}

type PathArrayToFloat64Array2SliceSlice struct {
	Val *[][][2]float64
}

func (v PathArrayToFloat64Array2SliceSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParsePathArray(data)
	polygons := make([][][2]float64, len(elems))
	for i := 0; i < len(elems); i++ {
		if elems[i] == nil {
			continue
		}

		polygon := make([][2]float64, len(elems[i]))

		for j := 0; j < len(elems[i]); j++ {
			f0, err := strconv.ParseFloat(string(elems[i][j][0]), 64)
			if err != nil {
				return err
			}
			f1, err := strconv.ParseFloat(string(elems[i][j][1]), 64)
			if err != nil {
				return err
			}

			polygon[j][0] = f0
			polygon[j][1] = f1
		}

		polygons[i] = polygon
	}

	*v.Val = polygons
	return nil
}
