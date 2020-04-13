package convert

import (
	"database/sql/driver"
	"time"
)

type DateRangeArrayFromTimeArray2Slice struct {
	Val [][2]time.Time
}

func (v DateRangeArrayFromTimeArray2Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.Val {
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

type DateRangeArrayToTimeArray2Slice struct {
	Val *[][2]time.Time
}

func (v DateRangeArrayToTimeArray2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
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

	*v.Val = ranges
	return nil
}
