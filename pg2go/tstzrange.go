package pg2go

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

// TstzRangeFromTimeArray2 returns a driver.Valuer that produces a PostgreSQL tstzrange from the given Go [2]time.Time.
func TstzRangeFromTimeArray2(val [2]time.Time) driver.Valuer {
	return tstzRangeFromTimeArray2{val: val}
}

// TstzRangeToTimeArray2 returns an sql.Scanner that converts a PostgreSQL tstzrange into a Go [2]time.Time and sets it to val.
func TstzRangeToTimeArray2(val *[2]time.Time) sql.Scanner {
	return tstzRangeToTimeArray2{val: val}
}

type tstzRangeFromTimeArray2 struct {
	val [2]time.Time
}

func (v tstzRangeFromTimeArray2) Value() (driver.Value, error) {
	// len(`"yyyy-mm-dd hh:mm:ss-tz"`) == 24
	size := 3 + (2 * 24)

	out := make([]byte, 2, size)
	out[0], out[1] = '[', '"'

	out = append(out, v.val[0].Format(timestamptzLayout)...)
	out = append(out, '"', ',', '"')
	out = append(out, v.val[1].Format(timestamptzLayout)...)
	out = append(out, '"', ')')
	return out, nil
}

type tstzRangeToTimeArray2 struct {
	val *[2]time.Time
}

func (v tstzRangeToTimeArray2) Scan(src interface{}) error {
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
		s = s[1 : len(s)-1] // remove surrounding double quotes
		if lo, err = time.Parse(timestamptzLayout, string(s)); err != nil {
			return err
		}
	}
	if len(elems[1]) > 0 {
		s := elems[1]
		s = s[1 : len(s)-1] // remove surrounding double quotes
		if hi, err = time.Parse(timestamptzLayout, string(s)); err != nil {
			return err
		}
	}

	v.val[0] = lo
	v.val[1] = hi
	return nil
}
