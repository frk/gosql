package convert

import (
	"database/sql/driver"
	"time"
)

type TimestamptzArrayFromTimeSlice struct {
	Val []time.Time
}

func (v TimestamptzArrayFromTimeSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := (len(v.Val) * 24) + // len(`"yyyy-mm-dd hh:mm:ss-tz"`) == 24
		(len(v.Val) - 1) + // number of commas between array elements
		2 // surrounding curly braces

	out := make([]byte, 1, size)
	out[0] = '{'

	for _, ts := range v.Val {
		out = append(out, '"')
		out = append(out, ts.Format(timestamptzLayout)...)
		out = append(out, '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type TimestamptzArrayToTimeSlice struct {
	Val *[]time.Time
}

func (v TimestamptzArrayToTimeSlice) Scan(src interface{}) error {
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
		ts := elems[i]
		if len(ts) > 0 && ts[0] == 'N' { // NULL
			continue
		}

		ts = ts[1 : len(ts)-1] // drop surrounding double quotes
		t, err := time.Parse(timestamptzLayout, string(ts))
		if err != nil {
			return err
		}
		times[i] = t
	}

	*v.Val = times
	return nil
}
