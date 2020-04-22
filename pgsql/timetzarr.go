package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

// TimetzArrayFromTimeSlice returns a driver.Valuer that produces a PostgreSQL timetz[] from the given Go []time.Time.
func TimetzArrayFromTimeSlice(val []time.Time) driver.Valuer {
	return timetzArrayFromTimeSlice{val: val}
}

// TimetzArrayToTimeSlice returns an sql.Scanner that converts a PostgreSQL timetz[] into a Go []time.Time and sets it to val.
func TimetzArrayToTimeSlice(val *[]time.Time) sql.Scanner {
	return timetzArrayToTimeSlice{val: val}
}

type timetzArrayFromTimeSlice struct {
	val []time.Time
}

func (v timetzArrayFromTimeSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := 2 + (len(v.val) - 1) + (len(v.val) * 14) // len("hh:mm:ss-hh:ss") == 14

	out := make([]byte, 1, size)
	out[0] = '{'

	for _, t := range v.val {
		out = append(out, t.Format(timetzLayout)...)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type timetzArrayToTimeSlice struct {
	val *[]time.Time
}

func (v timetzArrayToTimeSlice) Scan(src interface{}) error {
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
		t, err := time.Parse(timetzLayout, string(elems[i]))
		if err != nil {
			return err
		}
		times[i] = t
	}

	*v.val = times
	return nil
}
