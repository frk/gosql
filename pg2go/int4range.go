package pg2go

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// Int4RangeFromIntArray2 returns a driver.Valuer that produces a PostgreSQL int4range from the given Go [2]int.
func Int4RangeFromIntArray2(val [2]int) driver.Valuer {
	return int4RangeFromIntArray2{val: val}
}

// Int4RangeToIntArray2 returns an sql.Scanner that converts a PostgreSQL int4range into a Go [2]int and sets it to val.
func Int4RangeToIntArray2(val *[2]int) sql.Scanner {
	return int4RangeToIntArray2{val: val}
}

// Int4RangeFromInt8Array2 returns a driver.Valuer that produces a PostgreSQL int4range from the given Go [2]int8.
func Int4RangeFromInt8Array2(val [2]int8) driver.Valuer {
	return int4RangeFromInt8Array2{val: val}
}

// Int4RangeToInt8Array2 returns an sql.Scanner that converts a PostgreSQL int4range into a Go [2]int8 and sets it to val.
func Int4RangeToInt8Array2(val *[2]int8) sql.Scanner {
	return int4RangeToInt8Array2{val: val}
}

// Int4RangeFromInt16Array2 returns a driver.Valuer that produces a PostgreSQL int4range from the given Go [2]int16.
func Int4RangeFromInt16Array2(val [2]int16) driver.Valuer {
	return int4RangeFromInt16Array2{val: val}
}

// Int4RangeToInt16Array2 returns an sql.Scanner that converts a PostgreSQL int4range into a Go [2]int16 and sets it to val.
func Int4RangeToInt16Array2(val *[2]int16) sql.Scanner {
	return int4RangeToInt16Array2{val: val}
}

// Int4RangeFromInt32Array2 returns a driver.Valuer that produces a PostgreSQL int4range from the given Go [2]int32.
func Int4RangeFromInt32Array2(val [2]int32) driver.Valuer {
	return int4RangeFromInt32Array2{val: val}
}

// Int4RangeToInt32Array2 returns an sql.Scanner that converts a PostgreSQL int4range into a Go [2]int32 and sets it to val.
func Int4RangeToInt32Array2(val *[2]int32) sql.Scanner {
	return int4RangeToInt32Array2{val: val}
}

// Int4RangeFromInt64Array2 returns a driver.Valuer that produces a PostgreSQL int4range from the given Go [2]int64.
func Int4RangeFromInt64Array2(val [2]int64) driver.Valuer {
	return int4RangeFromInt64Array2{val: val}
}

// Int4RangeToInt64Array2 returns an sql.Scanner that converts a PostgreSQL int4range into a Go [2]int64 and sets it to val.
func Int4RangeToInt64Array2(val *[2]int64) sql.Scanner {
	return int4RangeToInt64Array2{val: val}
}

// Int4RangeFromUintArray2 returns a driver.Valuer that produces a PostgreSQL int4range from the given Go [2]uint.
func Int4RangeFromUintArray2(val [2]uint) driver.Valuer {
	return int4RangeFromUintArray2{val: val}
}

// Int4RangeToUintArray2 returns an sql.Scanner that converts a PostgreSQL int4range into a Go [2]uint and sets it to val.
func Int4RangeToUintArray2(val *[2]uint) sql.Scanner {
	return int4RangeToUintArray2{val: val}
}

// Int4RangeFromUint8Array2 returns a driver.Valuer that produces a PostgreSQL int4range from the given Go [2]uint8.
func Int4RangeFromUint8Array2(val [2]uint8) driver.Valuer {
	return int4RangeFromUint8Array2{val: val}
}

// Int4RangeToUint8Array2 returns an sql.Scanner that converts a PostgreSQL int4range into a Go [2]uint8 and sets it to val.
func Int4RangeToUint8Array2(val *[2]uint8) sql.Scanner {
	return int4RangeToUint8Array2{val: val}
}

// Int4RangeFromUint16Array2 returns a driver.Valuer that produces a PostgreSQL int4range from the given Go [2]uint16.
func Int4RangeFromUint16Array2(val [2]uint16) driver.Valuer {
	return int4RangeFromUint16Array2{val: val}
}

// Int4RangeToUint16Array2 returns an sql.Scanner that converts a PostgreSQL int4range into a Go [2]uint16 and sets it to val.
func Int4RangeToUint16Array2(val *[2]uint16) sql.Scanner {
	return int4RangeToUint16Array2{val: val}
}

// Int4RangeFromUint32Array2 returns a driver.Valuer that produces a PostgreSQL int4range from the given Go [2]uint32.
func Int4RangeFromUint32Array2(val [2]uint32) driver.Valuer {
	return int4RangeFromUint32Array2{val: val}
}

// Int4RangeToUint32Array2 returns an sql.Scanner that converts a PostgreSQL int4range into a Go [2]uint32 and sets it to val.
func Int4RangeToUint32Array2(val *[2]uint32) sql.Scanner {
	return int4RangeToUint32Array2{val: val}
}

// Int4RangeFromUint64Array2 returns a driver.Valuer that produces a PostgreSQL int4range from the given Go [2]uint64.
func Int4RangeFromUint64Array2(val [2]uint64) driver.Valuer {
	return int4RangeFromUint64Array2{val: val}
}

// Int4RangeToUint64Array2 returns an sql.Scanner that converts a PostgreSQL int4range into a Go [2]uint64 and sets it to val.
func Int4RangeToUint64Array2(val *[2]uint64) sql.Scanner {
	return int4RangeToUint64Array2{val: val}
}

// Int4RangeFromFloat32Array2 returns a driver.Valuer that produces a PostgreSQL int4range from the given Go [2]float32.
func Int4RangeFromFloat32Array2(val [2]float32) driver.Valuer {
	return int4RangeFromFloat32Array2{val: val}
}

// Int4RangeToFloat32Array2 returns an sql.Scanner that converts a PostgreSQL int4range into a Go [2]float32 and sets it to val.
func Int4RangeToFloat32Array2(val *[2]float32) sql.Scanner {
	return int4RangeToFloat32Array2{val: val}
}

// Int4RangeFromFloat64Array2 returns a driver.Valuer that produces a PostgreSQL int4range from the given Go [2]float64.
func Int4RangeFromFloat64Array2(val [2]float64) driver.Valuer {
	return int4RangeFromFloat64Array2{val: val}
}

// Int4RangeToFloat64Array2 returns an sql.Scanner that converts a PostgreSQL int4range into a Go [2]float64 and sets it to val.
func Int4RangeToFloat64Array2(val *[2]float64) sql.Scanner {
	return int4RangeToFloat64Array2{val: val}
}

type int4RangeFromIntArray2 struct {
	val [2]int
}

func (v int4RangeFromIntArray2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type int4RangeToIntArray2 struct {
	val *[2]int
}

func (v int4RangeToIntArray2) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	var lo, hi int64
	elems := pgParseRange(data)
	if len(elems[0]) > 0 {
		if lo, err = strconv.ParseInt(string(elems[0]), 10, 32); err != nil {
			return err
		}
	}
	if len(elems[1]) > 0 {
		if hi, err = strconv.ParseInt(string(elems[1]), 10, 32); err != nil {
			return err
		}
	}

	v.val[0] = int(lo)
	v.val[1] = int(hi)
	return nil
}

type int4RangeFromInt8Array2 struct {
	val [2]int8
}

func (v int4RangeFromInt8Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type int4RangeToInt8Array2 struct {
	val *[2]int8
}

func (v int4RangeToInt8Array2) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	var lo, hi int64
	elems := pgParseRange(data)
	if len(elems[0]) > 0 {
		if lo, err = strconv.ParseInt(string(elems[0]), 10, 8); err != nil {
			return err
		}
	}
	if len(elems[1]) > 0 {
		if hi, err = strconv.ParseInt(string(elems[1]), 10, 8); err != nil {
			return err
		}
	}

	v.val[0] = int8(lo)
	v.val[1] = int8(hi)
	return nil
}

type int4RangeFromInt16Array2 struct {
	val [2]int16
}

func (v int4RangeFromInt16Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type int4RangeToInt16Array2 struct {
	val *[2]int16
}

func (v int4RangeToInt16Array2) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	var lo, hi int64
	elems := pgParseRange(data)
	if len(elems[0]) > 0 {
		if lo, err = strconv.ParseInt(string(elems[0]), 10, 16); err != nil {
			return err
		}
	}
	if len(elems[1]) > 0 {
		if hi, err = strconv.ParseInt(string(elems[1]), 10, 16); err != nil {
			return err
		}
	}

	v.val[0] = int16(lo)
	v.val[1] = int16(hi)
	return nil
}

type int4RangeFromInt32Array2 struct {
	val [2]int32
}

func (v int4RangeFromInt32Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type int4RangeToInt32Array2 struct {
	val *[2]int32
}

func (v int4RangeToInt32Array2) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	var lo, hi int64
	elems := pgParseRange(data)
	if len(elems[0]) > 0 {
		if lo, err = strconv.ParseInt(string(elems[0]), 10, 32); err != nil {
			return err
		}
	}
	if len(elems[1]) > 0 {
		if hi, err = strconv.ParseInt(string(elems[1]), 10, 32); err != nil {
			return err
		}
	}

	v.val[0] = int32(lo)
	v.val[1] = int32(hi)
	return nil
}

type int4RangeFromInt64Array2 struct {
	val [2]int64
}

func (v int4RangeFromInt64Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, v.val[0], 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, v.val[1], 10)
	out = append(out, ')')
	return out, nil
}

type int4RangeToInt64Array2 struct {
	val *[2]int64
}

func (v int4RangeToInt64Array2) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	var lo, hi int64
	elems := pgParseRange(data)
	if len(elems[0]) > 0 {
		if lo, err = strconv.ParseInt(string(elems[0]), 10, 32); err != nil {
			return err
		}
	}
	if len(elems[1]) > 0 {
		if hi, err = strconv.ParseInt(string(elems[1]), 10, 32); err != nil {
			return err
		}
	}

	v.val[0] = lo
	v.val[1] = hi
	return nil
}

type int4RangeFromUintArray2 struct {
	val [2]uint
}

func (v int4RangeFromUintArray2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, uint64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, uint64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type int4RangeToUintArray2 struct {
	val *[2]uint
}

func (v int4RangeToUintArray2) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	var lo, hi uint64
	elems := pgParseRange(data)
	if len(elems[0]) > 0 {
		if lo, err = strconv.ParseUint(string(elems[0]), 10, 32); err != nil {
			return err
		}
	}
	if len(elems[1]) > 0 {
		if hi, err = strconv.ParseUint(string(elems[1]), 10, 32); err != nil {
			return err
		}
	}

	v.val[0] = uint(lo)
	v.val[1] = uint(hi)
	return nil
}

type int4RangeFromUint8Array2 struct {
	val [2]uint8
}

func (v int4RangeFromUint8Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, uint64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, uint64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type int4RangeToUint8Array2 struct {
	val *[2]uint8
}

func (v int4RangeToUint8Array2) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	var lo, hi uint64
	elems := pgParseRange(data)
	if len(elems[0]) > 0 {
		if lo, err = strconv.ParseUint(string(elems[0]), 10, 8); err != nil {
			return err
		}
	}
	if len(elems[1]) > 0 {
		if hi, err = strconv.ParseUint(string(elems[1]), 10, 8); err != nil {
			return err
		}
	}

	v.val[0] = uint8(lo)
	v.val[1] = uint8(hi)
	return nil
}

type int4RangeFromUint16Array2 struct {
	val [2]uint16
}

func (v int4RangeFromUint16Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, uint64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, uint64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type int4RangeToUint16Array2 struct {
	val *[2]uint16
}

func (v int4RangeToUint16Array2) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	var lo, hi uint64
	elems := pgParseRange(data)
	if len(elems[0]) > 0 {
		if lo, err = strconv.ParseUint(string(elems[0]), 10, 16); err != nil {
			return err
		}
	}
	if len(elems[1]) > 0 {
		if hi, err = strconv.ParseUint(string(elems[1]), 10, 16); err != nil {
			return err
		}
	}

	v.val[0] = uint16(lo)
	v.val[1] = uint16(hi)
	return nil
}

type int4RangeFromUint32Array2 struct {
	val [2]uint32
}

func (v int4RangeFromUint32Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, uint64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, uint64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type int4RangeToUint32Array2 struct {
	val *[2]uint32
}

func (v int4RangeToUint32Array2) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	var lo, hi uint64
	elems := pgParseRange(data)
	if len(elems[0]) > 0 {
		if lo, err = strconv.ParseUint(string(elems[0]), 10, 32); err != nil {
			return err
		}
	}
	if len(elems[1]) > 0 {
		if hi, err = strconv.ParseUint(string(elems[1]), 10, 32); err != nil {
			return err
		}
	}

	v.val[0] = uint32(lo)
	v.val[1] = uint32(hi)
	return nil
}

type int4RangeFromUint64Array2 struct {
	val [2]uint64
}

func (v int4RangeFromUint64Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, v.val[0], 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, v.val[1], 10)
	out = append(out, ')')
	return out, nil
}

type int4RangeToUint64Array2 struct {
	val *[2]uint64
}

func (v int4RangeToUint64Array2) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	var lo, hi uint64
	elems := pgParseRange(data)
	if len(elems[0]) > 0 {
		if lo, err = strconv.ParseUint(string(elems[0]), 10, 32); err != nil {
			return err
		}
	}
	if len(elems[1]) > 0 {
		if hi, err = strconv.ParseUint(string(elems[1]), 10, 32); err != nil {
			return err
		}
	}

	v.val[0] = lo
	v.val[1] = hi
	return nil
}

type int4RangeFromFloat32Array2 struct {
	val [2]float32
}

func (v int4RangeFromFloat32Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type int4RangeToFloat32Array2 struct {
	val *[2]float32
}

func (v int4RangeToFloat32Array2) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	var lo, hi int64
	elems := pgParseRange(data)
	if len(elems[0]) > 0 {
		if lo, err = strconv.ParseInt(string(elems[0]), 10, 32); err != nil {
			return err
		}
	}
	if len(elems[1]) > 0 {
		if hi, err = strconv.ParseInt(string(elems[1]), 10, 32); err != nil {
			return err
		}
	}

	v.val[0] = float32(lo)
	v.val[1] = float32(hi)
	return nil
}

type int4RangeFromFloat64Array2 struct {
	val [2]float64
}

func (v int4RangeFromFloat64Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type int4RangeToFloat64Array2 struct {
	val *[2]float64
}

func (v int4RangeToFloat64Array2) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	var lo, hi int64
	elems := pgParseRange(data)
	if len(elems[0]) > 0 {
		if lo, err = strconv.ParseInt(string(elems[0]), 10, 32); err != nil {
			return err
		}
	}
	if len(elems[1]) > 0 {
		if hi, err = strconv.ParseInt(string(elems[1]), 10, 32); err != nil {
			return err
		}
	}

	v.val[0] = float64(lo)
	v.val[1] = float64(hi)
	return nil
}
