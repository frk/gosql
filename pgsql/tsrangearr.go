package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

// TsRangeArrayFromTimeArray2Slice returns a driver.Valuer that produces a PostgreSQL tsrange[] from the given Go [][2]time.Time.
func TsRangeArrayFromTimeArray2Slice(val [][2]time.Time) driver.Valuer {
	return tsRangeArrayFromTimeArray2Slice{val: val}
}

// TsRangeArrayToTimeArray2Slice returns an sql.Scanner that converts a PostgreSQL tsrange[] into a Go [][2]time.Time and sets it to val.
func TsRangeArrayToTimeArray2Slice(val *[][2]time.Time) sql.Scanner {
	return tsRangeArrayToTimeArray2Slice{val: val}
}

type tsRangeArrayFromTimeArray2Slice struct {
	val [][2]time.Time
}

func (v tsRangeArrayFromTimeArray2Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := (len(v.val) * 51) + // len(`"[\"yyyy-mm-dd hh:mm:ss\",\"yyyy-mm-dd hh:mm:ss\")"`) == 51
		(len(v.val) - 1) + // number of commas between array elements
		2 // surrounding curly braces

	out := make([]byte, 1, size)
	out[0] = '{'

	for _, a := range v.val {
		out = append(out, '"', '[', '\\', '"')
		out = append(out, a[0].Format(timestampLayout)...)
		out = append(out, '\\', '"', ',', '\\', '"')
		out = append(out, a[1].Format(timestampLayout)...)
		out = append(out, '\\', '"', ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type tsRangeArrayToTimeArray2Slice struct {
	val *[][2]time.Time
}

func (v tsRangeArrayToTimeArray2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseQuotedStringArray(data)
	slice := make([][2]time.Time, len(elems))
	for i := 0; i < len(elems); i++ {
		a := pgParseRange(elems[i])

		// drop surrounding escaped double quotes
		a[0] = a[0][2 : len(a[0])-2]
		a[1] = a[1][2 : len(a[1])-2]

		var t0, t1 time.Time
		t0, err = time.ParseInLocation(timestampLayout, string(a[0]), noZone)
		if err != nil {
			return err
		}
		t1, err = time.ParseInLocation(timestampLayout, string(a[1]), noZone)
		if err != nil {
			return err
		}

		slice[i][0] = t0
		slice[i][1] = t1
	}

	*v.val = slice
	return nil
}
