package pg2go

import (
	"database/sql/driver"
	"time"
)

type TimeArrayFromTimeSlice struct {
	Val []time.Time
}

func (v TimeArrayFromTimeSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := 2 + (len(v.Val) - 1) + (len(v.Val) * 8) // len("hh:mm:ss") == 8

	out := make([]byte, 1, size)
	out[0] = '{'

	for _, t := range v.Val {
		out = append(out, t.Format(timeLayout)...)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type TimeArrayToTimeSlice struct {
	Val *[]time.Time
}

func (v TimeArrayToTimeSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseCommaArray(data)
	times := make([]time.Time, len(elems))
	for i := 0; i < len(elems); i++ {
		t, err := time.Parse(timeLayout, string(elems[i]))
		if err != nil {
			return err
		}
		times[i] = t
	}

	*v.Val = times
	return nil
}
