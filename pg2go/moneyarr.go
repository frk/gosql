package pg2go

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// MoneyArrayFromInt64Slice returns a driver.Valuer that produces a PostgreSQL money[] from the given Go []int64.
func MoneyArrayFromInt64Slice(val []int64) driver.Valuer {
	return moneyArrayFromInt64Slice{val: val}
}

// MoneyArrayToInt64Slice returns an sql.Scanner that converts a PostgreSQL money[] into a Go []int64 and sets it to val.
func MoneyArrayToInt64Slice(val *[]int64) sql.Scanner {
	return moneyArrayToInt64Slice{val: val}
}

type moneyArrayFromInt64Slice struct {
	val []int64
}

func (v moneyArrayFromInt64Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := (len(v.val) * 5) + // len(`$0.00`) == 5 (min length)
		(len(v.val) - 1) + // number of commas between array elements
		2 // curly braces

	out := make([]byte, 1, size)
	out[0] = '{'

	for i := 0; i < len(v.val); i++ {
		out = append(out, '$')
		out = strconv.AppendFloat(out, float64(v.val[i])/100.0, 'f', 2, 64)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type moneyArrayToInt64Slice struct {
	val *[]int64
}

func (v moneyArrayToInt64Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
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

	*v.val = int64s
	return nil
}
