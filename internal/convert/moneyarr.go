package convert

import (
	"database/sql/driver"
	"strconv"
)

type MoneyArrayFromInt64Slice struct {
	Val []int64
}

func (v MoneyArrayFromInt64Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := (len(v.Val) * 5) + // len(`$0.00`) == 5 (min length)
		(len(v.Val) - 1) + // number of commas between array elements
		2 // curly braces

	out := make([]byte, 1, size)
	out[0] = '{'

	for i := 0; i < len(v.Val); i++ {
		out = append(out, '$')
		out = strconv.AppendFloat(out, float64(v.Val[i])/100.0, 'f', 2, 64)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type MoneyArrayToInt64Slice struct {
	Val *[]int64
}

func (v MoneyArrayToInt64Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseCommaArray(data)
	int64s := make([]int64, len(elems))
	for i := 0; i < len(elems); i++ {
		elems[i] = elems[i][1:] // drop $
		f64, err := strconv.ParseFloat(string(elems[i]), 64)
		if err != nil {
			return err
		}

		int64s[i] = int64(f64 * 100.0)
	}

	*v.Val = int64s
	return nil
}
