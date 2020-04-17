package convert

import (
	"database/sql/driver"
	"time"
)

type TstzRangeFromTimeArray2 struct {
	Val [2]time.Time
}

func (v TstzRangeFromTimeArray2) Value() (driver.Value, error) {
	// len(`"yyyy-mm-dd hh:mm:ss-tz"`) == 24
	size := 3 + (2 * 24)

	out := make([]byte, 2, size)
	out[0], out[1] = '[', '"'

	out = append(out, v.Val[0].Format(timestamptzLayout)...)
	out = append(out, '"', ',', '"')
	out = append(out, v.Val[1].Format(timestamptzLayout)...)
	out = append(out, '"', ')')
	return out, nil
}

type TstzRangeToTimeArray2 struct {
	Val *[2]time.Time
}

func (v TstzRangeToTimeArray2) Scan(src interface{}) error {
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

	v.Val[0] = lo
	v.Val[1] = hi
	return nil
}
