package convert

import (
	"database/sql/driver"
	"strconv"
)

type MoneyFromInt64 struct {
	Val int64
}

func (v MoneyFromInt64) Value() (driver.Value, error) {
	out := []byte{'$'}
	out = strconv.AppendFloat(out, float64(v.Val)/100.0, 'f', 2, 64)
	return out, nil
}

type MoneyToInt64 struct {
	Val *int64
}

func (v MoneyToInt64) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	data = data[1:] // drop $
	f64, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return err
	}

	*v.Val = int64(f64 * 100.0)
	return nil
}
