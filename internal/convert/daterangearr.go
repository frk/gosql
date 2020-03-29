package convert

import (
	"database/sql/driver"
	"time"
)

type DateRangeArrayFromTimeArray2Slice struct {
	V [][2]time.Time
}

func (v DateRangeArrayFromTimeArray2Slice) Value() (driver.Value, error) {
	if v.V == nil {
		return nil, nil
	} else if len(v.V) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.V {
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
	V *[][2]time.Time
}

func (v DateRangeArrayToTimeArray2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.V = nil
		return nil
	}

	elems := pgparsearray0(data)
	ranges := make([][2]time.Time, len(elems))

	for i, elem := range elems {
		elem = elem[1 : len(elem)-1] // drop surrounding double quotes

		var t0, t1 time.Time
		arr := pgparserange(elem)

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

	*v.V = ranges
	return nil
}
