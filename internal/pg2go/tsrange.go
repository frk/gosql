package pg2go

import (
	"database/sql/driver"
	"time"
)

type TsRangeFromTimeArray2 struct {
	Val [2]time.Time
}

func (v TsRangeFromTimeArray2) Value() (driver.Value, error) {
	// len(`"yyyy-mm-dd hh:mm:ss"`) == 21
	size := 3 + (len(v.Val) * 21)

	out := make([]byte, 2, size)
	out[0], out[1] = '[', '"'

	out = append(out, v.Val[0].Format(timestampLayout)...)
	out = append(out, '"', ',', '"')
	out = append(out, v.Val[1].Format(timestampLayout)...)
	out = append(out, '"', ')')
	return out, nil
}

type TsRangeToTimeArray2 struct {
	Val *[2]time.Time
}

func (v TsRangeToTimeArray2) Scan(src interface{}) error {
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

	v.Val[0] = lo
	v.Val[1] = hi
	return nil
}
