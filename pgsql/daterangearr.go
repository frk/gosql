package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

// DateRangeArrayFromTimeArray2Slice returns a driver.Valuer that produces a PostgreSQL daterange[] from the given Go [][2]time.Time.
func DateRangeArrayFromTimeArray2Slice(val [][2]time.Time) driver.Valuer {
	return dateRangeArrayFromTimeArray2Slice{val: val}
}

// DateRangeArrayToTimeArray2Slice returns an sql.Scanner that converts a PostgreSQL daterange[] into a Go [][2]time.Time and sets it to val.
func DateRangeArrayToTimeArray2Slice(val *[][2]time.Time) sql.Scanner {
	return dateRangeArrayToTimeArray2Slice{val: val}
}

type dateRangeArrayFromTimeArray2Slice struct {
	val [][2]time.Time
}

func (v dateRangeArrayFromTimeArray2Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.val {
		if !a[0].IsZero() {
			out = append(out, '"', '[')
			out = append(out, []byte(a[0].Format(dateLayout))...)
		} else {
			out = append(out, '"', '(')
		}

		out = append(out, ',')

		if !a[1].IsZero() {
			out = append(out, []byte(a[1].Format(dateLayout))...)
		}

		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type dateRangeArrayToTimeArray2Slice struct {
	val *[][2]time.Time
}

func (v dateRangeArrayToTimeArray2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.val = nil
		return nil
	}

	elems := pgParseQuotedStringArray(data)
	ranges := make([][2]time.Time, len(elems))

	for i, elem := range elems {
		var t0, t1 time.Time
		arr := pgParseRange(elem)

		if len(arr[0]) > 0 {
			if t0, err = pgparsedate(arr[0]); err != nil {
				return err
			}
		}
		if len(arr[1]) > 0 {
			if t1, err = pgparsedate(arr[1]); err != nil {
				return err
			}
		}

		ranges[i][0] = t0
		ranges[i][1] = t1
	}

	*v.val = ranges
	return nil
}
