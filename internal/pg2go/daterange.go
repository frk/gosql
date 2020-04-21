package pg2go

import (
	"database/sql/driver"
	"time"
)

// NOTE(mkopriva): Given the simple structure of [2]time.Time there's no way
// to indicate the bounds of the range, therefore DateRangeFromTimeArray2 defaults
// to using the bounds configuration that is used by postgres for output, which is:
//
// The array's first element as the inclusive-lower-bound and the array's second
// element as the exclusive-upper-bound, that is, an array of <date1> and <date2>
// will be stored as: "[<date1>,<date2>)".

// TODO(mkopriva): To allow some control over the limitation mentioned above, a
// tag option to specify the bounds' configuration could be supported by the gosql
// package, then the generator could produce code that passes that option to
// DateRangeFromTimeArray2 which would be extended by a field that would hold
// that option and use it to determine how to format the driver.Value.
//
// In the same manner the option would be passed to DateRangeToTimeArray2 so that
// it knows how to interpret the date values -- whether to add/subtract a day the
// time.Time value after parsing.
//
// XXX(mkopriva): If the above is implemented then the programmer must be made
// aware that if the bounds' configuration is changed during the life-time of the
// app the dates will be interpreted differently and may therefore produce unexpected
// reusults.

type DateRangeFromTimeArray2 struct {
	Val [2]time.Time
}

func (v DateRangeFromTimeArray2) Value() (driver.Value, error) {
	out := make([]byte, 1)
	if !v.Val[0].IsZero() {
		out[0] = '['
		out = append(out, []byte(v.Val[0].Format(dateLayout))...)
	} else {
		out[0] = '('
	}

	out = append(out, ',')

	if !v.Val[1].IsZero() {
		out = append(out, []byte(v.Val[1].Format(dateLayout))...)
	}

	out = append(out, ')')

	return out, nil
}

type DateRangeToTimeArray2 struct {
	Val *[2]time.Time
}

func (v DateRangeToTimeArray2) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	var t0, t1 time.Time
	elems := pgParseRange(data)
	if len(elems[0]) > 0 {
		if t0, err = pgparsedate(elems[0]); err != nil {
			return err
		}
	}
	if len(elems[1]) > 0 {
		if t1, err = pgparsedate(elems[1]); err != nil {
			return err
		}
	}

	v.Val[0] = t0
	v.Val[1] = t1
	return nil
}
