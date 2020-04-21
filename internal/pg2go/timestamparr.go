package pg2go

import (
	"database/sql/driver"
	"time"
)

type TimestampArrayFromTimeSlice struct {
	Val []time.Time
}

func (v TimestampArrayFromTimeSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	// len(`"yyyy-mm-dd hh:mm:ss"`) == 21
	size := (len(v.Val) * 21) + len(v.Val) - 1

	out := make([]byte, 1, size)
	out[0] = '{'

	for _, ts := range v.Val {
		out = append(out, '"')
		out = append(out, ts.Format(timestampLayout)...)
		out = append(out, '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type TimestampArrayToTimeSlice struct {
	Val *[]time.Time
}

func (v TimestampArrayToTimeSlice) Scan(src interface{}) error {
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
		t, err := time.ParseInLocation(timestampLayout, string(ts), noZone)
		if err != nil {
			return err
		}
		times[i] = t
	}

	*v.Val = times
	return nil
}
