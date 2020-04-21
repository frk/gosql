package pg2go

import (
	"database/sql/driver"
	"strconv"
)

type LsegArrayFromFloat64Array2Array2Slice struct {
	Val [][2][2]float64
}

func (v LsegArrayFromFloat64Array2Array2Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := (len(v.Val) * 15) + // len(`"[(x,y),(x,y)]"`) == 15
		(len(v.Val) - 1) + // number of commas between array elements
		2 // surrounding curly braces

	out := make([]byte, 1, size)
	out[0] = '{'

	for i := 0; i < len(v.Val); i++ {
		out = append(out, '"', '[', '(')
		out = strconv.AppendFloat(out, v.Val[i][0][0], 'f', -1, 64)
		out = append(out, ',')
		out = strconv.AppendFloat(out, v.Val[i][0][1], 'f', -1, 64)
		out = append(out, ')', ',', '(')
		out = strconv.AppendFloat(out, v.Val[i][1][0], 'f', -1, 64)
		out = append(out, ',')
		out = strconv.AppendFloat(out, v.Val[i][1][1], 'f', -1, 64)
		out = append(out, ')', ']', '"', ',')
	}

	out[len(out)-1] = '}'
	return out, nil
}

type LsegArrayToFloat64Array2Array2Slice struct {
	Val *[][2][2]float64
}

func (v LsegArrayToFloat64Array2Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	elems := pgParseLsegArray(data)
	lsegs := make([][2][2]float64, len(elems))
	for i := 0; i < len(elems); i++ {
		if lsegs[i][0][0], err = strconv.ParseFloat(string(elems[i][0][0]), 64); err != nil {
			return err
		}
		if lsegs[i][0][1], err = strconv.ParseFloat(string(elems[i][0][1]), 64); err != nil {
			return err
		}
		if lsegs[i][1][0], err = strconv.ParseFloat(string(elems[i][1][0]), 64); err != nil {
			return err
		}
		if lsegs[i][1][1], err = strconv.ParseFloat(string(elems[i][1][1]), 64); err != nil {
			return err
		}
	}

	*v.Val = lsegs
	return nil
}
