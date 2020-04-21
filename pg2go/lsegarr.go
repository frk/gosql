package pg2go

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// LsegArrayFromFloat64Array2Array2Slice returns a driver.Valuer that produces a PostgreSQL lseg[] from the given Go [][2][2]float64.
func LsegArrayFromFloat64Array2Array2Slice(val [][2][2]float64) driver.Valuer {
	return lsegArrayFromFloat64Array2Array2Slice{val: val}
}

// LsegArrayToFloat64Array2Array2Slice returns an sql.Scanner that converts a PostgreSQL lseg[] into a Go [][2][2]float64 and sets it to val.
func LsegArrayToFloat64Array2Array2Slice(val *[][2][2]float64) sql.Scanner {
	return lsegArrayToFloat64Array2Array2Slice{val: val}
}

type lsegArrayFromFloat64Array2Array2Slice struct {
	val [][2][2]float64
}

func (v lsegArrayFromFloat64Array2Array2Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := (len(v.val) * 15) + // len(`"[(x,y),(x,y)]"`) == 15
		(len(v.val) - 1) + // number of commas between array elements
		2 // surrounding curly braces

	out := make([]byte, 1, size)
	out[0] = '{'

	for i := 0; i < len(v.val); i++ {
		out = append(out, '"', '[', '(')
		out = strconv.AppendFloat(out, v.val[i][0][0], 'f', -1, 64)
		out = append(out, ',')
		out = strconv.AppendFloat(out, v.val[i][0][1], 'f', -1, 64)
		out = append(out, ')', ',', '(')
		out = strconv.AppendFloat(out, v.val[i][1][0], 'f', -1, 64)
		out = append(out, ',')
		out = strconv.AppendFloat(out, v.val[i][1][1], 'f', -1, 64)
		out = append(out, ')', ']', '"', ',')
	}

	out[len(out)-1] = '}'
	return out, nil
}

type lsegArrayToFloat64Array2Array2Slice struct {
	val *[][2][2]float64
}

func (v lsegArrayToFloat64Array2Array2Slice) Scan(src interface{}) error {
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

	*v.val = lsegs
	return nil
}
