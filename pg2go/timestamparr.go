package pg2go

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

// TimestampArrayFromTimeSlice returns a driver.Valuer that produces a PostgreSQL timestamp[] from the given Go []time.Time.
func TimestampArrayFromTimeSlice(val []time.Time) driver.Valuer {
	return timestampArrayFromTimeSlice{val: val}
}

// TimestampArrayToTimeSlice returns an sql.Scanner that converts a PostgreSQL timestamp[] into a Go []time.Time and sets it to val.
func TimestampArrayToTimeSlice(val *[]time.Time) sql.Scanner {
	return timestampArrayToTimeSlice{val: val}
}

type timestampArrayFromTimeSlice struct {
	val []time.Time
}

func (v timestampArrayFromTimeSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	// len(`"yyyy-mm-dd hh:mm:ss"`) == 21
	size := (len(v.val) * 21) + len(v.val) - 1

	out := make([]byte, 1, size)
	out[0] = '{'

	for _, ts := range v.val {
		out = append(out, '"')
		out = append(out, ts.Format(timestampLayout)...)
		out = append(out, '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type timestampArrayToTimeSlice struct {
	val *[]time.Time
}

func (v timestampArrayToTimeSlice) Scan(src interface{}) error {
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
		t, err := time.ParseInLocation(timestampLayout, string(ts), noZone)
		if err != nil {
			return err
		}
		times[i] = t
	}

	*v.val = times
	return nil
}
