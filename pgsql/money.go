package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// MoneyFromInt64 returns a driver.Valuer that produces a PostgreSQL money from the given Go int64.
func MoneyFromInt64(val int64) driver.Valuer {
	return moneyFromInt64{val: val}
}

// MoneyToInt64 returns an sql.Scanner that converts a PostgreSQL money into a Go int64 and sets it to val.
func MoneyToInt64(val *int64) sql.Scanner {
	return moneyToInt64{val: val}
}

type moneyFromInt64 struct {
	val int64
}

func (v moneyFromInt64) Value() (driver.Value, error) {
	out := []byte{'$'}
	out = strconv.AppendFloat(out, float64(v.val)/100.0, 'f', 2, 64)
	return out, nil
}

type moneyToInt64 struct {
	val *int64
}

func (v moneyToInt64) Scan(src interface{}) error {
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

	*v.val = int64(f64 * 100.0)
	return nil
}
