package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

// TsRangeFromTimeArray2 returns a driver.Valuer that produces a PostgreSQL tsrange from the given Go [2]time.Time.
func TsRangeFromTimeArray2(val [2]time.Time) driver.Valuer {
	return tsRangeFromTimeArray2{val: val}
}

// TsRangeToTimeArray2 returns an sql.Scanner that converts a PostgreSQL tsrange into a Go [2]time.Time and sets it to val.
func TsRangeToTimeArray2(val *[2]time.Time) sql.Scanner {
	return tsRangeToTimeArray2{val: val}
}

type tsRangeFromTimeArray2 struct {
	val [2]time.Time
}

func (v tsRangeFromTimeArray2) Value() (driver.Value, error) {
	// len(`"yyyy-mm-dd hh:mm:ss"`) == 21
	size := 3 + (len(v.val) * 21)

	out := make([]byte, 2, size)
	out[0], out[1] = '[', '"'

	out = append(out, v.val[0].Format(timestampLayout)...)
	out = append(out, '"', ',', '"')
	out = append(out, v.val[1].Format(timestampLayout)...)
	out = append(out, '"', ')')
	return out, nil
}

type tsRangeToTimeArray2 struct {
	val *[2]time.Time
}

func (v tsRangeToTimeArray2) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	var lo, hi time.Time
	elems := pgParseRange(data)
	if len(elems[0]) > 0 {
		s := elems[0]
		s = s[1 : len(s)-1]
		if lo, err = time.Parse(timestampLayout, string(s)); err != nil {
			return err
		}
	}
	if len(elems[1]) > 0 {
		s := elems[1]
		s = s[1 : len(s)-1]
		if hi, err = time.Parse(timestampLayout, string(s)); err != nil {
			return err
		}
	}

	v.val[0] = lo
	v.val[1] = hi
	return nil
}
