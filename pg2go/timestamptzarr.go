package pg2go

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

// TimestamptzArrayFromTimeSlice returns a driver.Valuer that produces a PostgreSQL timestamptz[] from the given Go []time.Time.
func TimestamptzArrayFromTimeSlice(val []time.Time) driver.Valuer {
	return timestamptzArrayFromTimeSlice{val: val}
}

// TimestamptzArrayToTimeSlice returns an sql.Scanner that converts a PostgreSQL timestamptz[] into a Go []time.Time and sets it to val.
func TimestamptzArrayToTimeSlice(val *[]time.Time) sql.Scanner {
	return timestamptzArrayToTimeSlice{val: val}
}

type timestamptzArrayFromTimeSlice struct {
	val []time.Time
}

func (v timestamptzArrayFromTimeSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := (len(v.val) * 24) + // len(`"yyyy-mm-dd hh:mm:ss-tz"`) == 24
		(len(v.val) - 1) + // number of commas between array elements
		2 // surrounding curly braces

	out := make([]byte, 1, size)
	out[0] = '{'

	for _, ts := range v.val {
		out = append(out, '"')
		out = append(out, ts.Format(timestamptzLayout)...)
		out = append(out, '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type timestamptzArrayToTimeSlice struct {
	val *[]time.Time
}

func (v timestamptzArrayToTimeSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
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

	*v.val = times
	return nil
}
