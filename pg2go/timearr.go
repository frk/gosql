package pg2go

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

// TimeArrayFromTimeSlice returns a driver.Valuer that produces a PostgreSQL time[] from the given Go []time.Time.
func TimeArrayFromTimeSlice(val []time.Time) driver.Valuer {
	return timeArrayFromTimeSlice{val: val}
}

// TimeArrayToTimeSlice returns an sql.Scanner that converts a PostgreSQL time[] into a Go []time.Time and sets it to val.
func TimeArrayToTimeSlice(val *[]time.Time) sql.Scanner {
	return timeArrayToTimeSlice{val: val}
}

type timeArrayFromTimeSlice struct {
	val []time.Time
}

func (v timeArrayFromTimeSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := 2 + (len(v.val) - 1) + (len(v.val) * 8) // len("hh:mm:ss") == 8

	out := make([]byte, 1, size)
	out[0] = '{'

	for _, t := range v.val {
		out = append(out, t.Format(timeLayout)...)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type timeArrayToTimeSlice struct {
	val *[]time.Time
}

func (v timeArrayToTimeSlice) Scan(src interface{}) error {
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
		t, err := time.Parse(timeLayout, string(elems[i]))
		if err != nil {
			return err
		}
		times[i] = t
	}

	*v.val = times
	return nil
}
