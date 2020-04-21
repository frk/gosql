package pg2go

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// PolygonArrayFromFloat64Array2SliceSlice returns a driver.Valuer that produces a PostgreSQL polygon[] from the given Go [][][2]float64.
func PolygonArrayFromFloat64Array2SliceSlice(val [][][2]float64) driver.Valuer {
	return polygonArrayFromFloat64Array2SliceSlice{val: val}
}

// PolygonArrayToFloat64Array2SliceSlice returns an sql.Scanner that converts a PostgreSQL polygon[] into a Go [][][2]float64 and sets it to val.
func PolygonArrayToFloat64Array2SliceSlice(val *[][][2]float64) sql.Scanner {
	return polygonArrayToFloat64Array2SliceSlice{val: val}
}

type polygonArrayFromFloat64Array2SliceSlice struct {
	val [][][2]float64
}

func (v polygonArrayFromFloat64Array2SliceSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := (len(v.val) * 4) + // len(`"()"`) == len(`NULL`) == 4
		(len(v.val) - 1) + // number of commas between elements
		2 // surrounding parentheses
	for i := 0; i < len(v.val); i++ {
		if v.val[i] != nil {
			size += (len(v.val[i]) * 5) + // len(`(x,y)`)
				(len(v.val[i]) - 1) // number of commas between points
		}

	}

	out := make([]byte, 1, size)
	out[0] = '{'

	for i := 0; i < len(v.val); i++ {
		if v.val[i] == nil {
			out = append(out, 'N', 'U', 'L', 'L', ',')
			continue
		}

		out = append(out, '"', '(')

		for j := 0; j < len(v.val[i]); j++ {
			out = append(out, '(')
			out = strconv.AppendFloat(out, v.val[i][j][0], 'f', -1, 64)
			out = append(out, ',')
			out = strconv.AppendFloat(out, v.val[i][j][1], 'f', -1, 64)
			out = append(out, ')', ',')
		}

		out[len(out)-1] = ')'
		out = append(out, '"', ',')
	}

	out[len(out)-1] = '}'
	return out, nil
}

type polygonArrayToFloat64Array2SliceSlice struct {
	val *[][][2]float64
}

func (v polygonArrayToFloat64Array2SliceSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
		return nil
	}

	elems := pgParsePolygonArray(data)
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

	*v.val = polygons
	return nil
}
