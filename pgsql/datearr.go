package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

// DateArrayFromTimeSlice returns a driver.Valuer that produces a PostgreSQL date[] from the given Go []time.Time.
func DateArrayFromTimeSlice(val []time.Time) driver.Valuer {
	return dateArrayFromTimeSlice{val: val}
}

// DateArrayToTimeSlice returns an sql.Scanner that converts a PostgreSQL date[] into a Go []time.Time and sets it to val.
func DateArrayToTimeSlice(val *[]time.Time) sql.Scanner {
	return dateArrayToTimeSlice{val: val}
}

type dateArrayFromTimeSlice struct {
	val []time.Time
}

func (v dateArrayFromTimeSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, t := range v.val {
		out = append(out, []byte(t.Format(dateLayout))...)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last ";" with "}"
	return out, nil
}

type dateArrayToTimeSlice struct {
	val *[]time.Time
}

func (v dateArrayToTimeSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseCommaArray(data)
	dates := make([]time.Time, len(elems))
	for i := 0; i < len(elems); i++ {
		t, err := pgparsedate(elems[i])
		if err != nil {
			return err
		}
		dates[i] = t
	}

	*v.val = dates
	return nil
}
