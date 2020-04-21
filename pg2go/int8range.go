package pg2go

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// Int8RangeFromIntArray2 returns a driver.Valuer that produces a PostgreSQL int8range from the given Go [2]int.
func Int8RangeFromIntArray2(val [2]int) driver.Valuer {
	return int8RangeFromIntArray2{val: val}
}

// Int8RangeToIntArray2 returns an sql.Scanner that converts a PostgreSQL int8range into a Go [2]int and sets it to val.
func Int8RangeToIntArray2(val *[2]int) sql.Scanner {
	return int8RangeToIntArray2{val: val}
}

// Int8RangeFromInt8Array2 returns a driver.Valuer that produces a PostgreSQL int8range from the given Go [2]int8.
func Int8RangeFromInt8Array2(val [2]int8) driver.Valuer {
	return int8RangeFromInt8Array2{val: val}
}

// Int8RangeToInt8Array2 returns an sql.Scanner that converts a PostgreSQL int8range into a Go [2]int8 and sets it to val.
func Int8RangeToInt8Array2(val *[2]int8) sql.Scanner {
	return int8RangeToInt8Array2{val: val}
}

// Int8RangeFromInt16Array2 returns a driver.Valuer that produces a PostgreSQL int8range from the given Go [2]int16.
func Int8RangeFromInt16Array2(val [2]int16) driver.Valuer {
	return int8RangeFromInt16Array2{val: val}
}

// Int8RangeToInt16Array2 returns an sql.Scanner that converts a PostgreSQL int8range into a Go [2]int16 and sets it to val.
func Int8RangeToInt16Array2(val *[2]int16) sql.Scanner {
	return int8RangeToInt16Array2{val: val}
}

// Int8RangeFromInt32Array2 returns a driver.Valuer that produces a PostgreSQL int8range from the given Go [2]int32.
func Int8RangeFromInt32Array2(val [2]int32) driver.Valuer {
	return int8RangeFromInt32Array2{val: val}
}

// Int8RangeToInt32Array2 returns an sql.Scanner that converts a PostgreSQL int8range into a Go [2]int32 and sets it to val.
func Int8RangeToInt32Array2(val *[2]int32) sql.Scanner {
	return int8RangeToInt32Array2{val: val}
}

// Int8RangeFromInt64Array2 returns a driver.Valuer that produces a PostgreSQL int8range from the given Go [2]int64.
func Int8RangeFromInt64Array2(val [2]int64) driver.Valuer {
	return int8RangeFromInt64Array2{val: val}
}

// Int8RangeToInt64Array2 returns an sql.Scanner that converts a PostgreSQL int8range into a Go [2]int64 and sets it to val.
func Int8RangeToInt64Array2(val *[2]int64) sql.Scanner {
	return int8RangeToInt64Array2{val: val}
}

// Int8RangeFromUintArray2 returns a driver.Valuer that produces a PostgreSQL int8range from the given Go [2]uint.
func Int8RangeFromUintArray2(val [2]uint) driver.Valuer {
	return int8RangeFromUintArray2{val: val}
}

// Int8RangeToUintArray2 returns an sql.Scanner that converts a PostgreSQL int8range into a Go [2]uint and sets it to val.
func Int8RangeToUintArray2(val *[2]uint) sql.Scanner {
	return int8RangeToUintArray2{val: val}
}

// Int8RangeFromUint8Array2 returns a driver.Valuer that produces a PostgreSQL int8range from the given Go [2]uint8.
func Int8RangeFromUint8Array2(val [2]uint8) driver.Valuer {
	return int8RangeFromUint8Array2{val: val}
}

// Int8RangeToUint8Array2 returns an sql.Scanner that converts a PostgreSQL int8range into a Go [2]uint8 and sets it to val.
func Int8RangeToUint8Array2(val *[2]uint8) sql.Scanner {
	return int8RangeToUint8Array2{val: val}
}

// Int8RangeFromUint16Array2 returns a driver.Valuer that produces a PostgreSQL int8range from the given Go [2]uint16.
func Int8RangeFromUint16Array2(val [2]uint16) driver.Valuer {
	return int8RangeFromUint16Array2{val: val}
}

// Int8RangeToUint16Array2 returns an sql.Scanner that converts a PostgreSQL int8range into a Go [2]uint16 and sets it to val.
func Int8RangeToUint16Array2(val *[2]uint16) sql.Scanner {
	return int8RangeToUint16Array2{val: val}
}

// Int8RangeFromUint32Array2 returns a driver.Valuer that produces a PostgreSQL int8range from the given Go [2]uint32.
func Int8RangeFromUint32Array2(val [2]uint32) driver.Valuer {
	return int8RangeFromUint32Array2{val: val}
}

// Int8RangeToUint32Array2 returns an sql.Scanner that converts a PostgreSQL int8range into a Go [2]uint32 and sets it to val.
func Int8RangeToUint32Array2(val *[2]uint32) sql.Scanner {
	return int8RangeToUint32Array2{val: val}
}

// Int8RangeFromUint64Array2 returns a driver.Valuer that produces a PostgreSQL int8range from the given Go [2]uint64.
func Int8RangeFromUint64Array2(val [2]uint64) driver.Valuer {
	return int8RangeFromUint64Array2{val: val}
}

// Int8RangeToUint64Array2 returns an sql.Scanner that converts a PostgreSQL int8range into a Go [2]uint64 and sets it to val.
func Int8RangeToUint64Array2(val *[2]uint64) sql.Scanner {
	return int8RangeToUint64Array2{val: val}
}

// Int8RangeFromFloat32Array2 returns a driver.Valuer that produces a PostgreSQL int8range from the given Go [2]float32.
func Int8RangeFromFloat32Array2(val [2]float32) driver.Valuer {
	return int8RangeFromFloat32Array2{val: val}
}

// Int8RangeToFloat32Array2 returns an sql.Scanner that converts a PostgreSQL int8range into a Go [2]float32 and sets it to val.
func Int8RangeToFloat32Array2(val *[2]float32) sql.Scanner {
	return int8RangeToFloat32Array2{val: val}
}

// Int8RangeFromFloat64Array2 returns a driver.Valuer that produces a PostgreSQL int8range from the given Go [2]float64.
func Int8RangeFromFloat64Array2(val [2]float64) driver.Valuer {
	return int8RangeFromFloat64Array2{val: val}
}

// Int8RangeToFloat64Array2 returns an sql.Scanner that converts a PostgreSQL int8range into a Go [2]float64 and sets it to val.
func Int8RangeToFloat64Array2(val *[2]float64) sql.Scanner {
	return int8RangeToFloat64Array2{val: val}
}

type int8RangeFromIntArray2 struct {
	val [2]int
}

func (v int8RangeFromIntArray2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type int8RangeToIntArray2 struct {
	val *[2]int
}

func (v int8RangeToIntArray2) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	var lo, hi int64
	elems := pgParseRange(data)
	if len(elems[0]) > 0 {
		if lo, err = strconv.ParseInt(string(elems[0]), 10, 64); err != nil {
			return err
		}
	}
	if len(elems[1]) > 0 {
		if hi, err = strconv.ParseInt(string(elems[1]), 10, 64); err != nil {
			return err
		}
	}

	v.val[0] = int(lo)
	v.val[1] = int(hi)
	return nil
}

type int8RangeFromInt8Array2 struct {
	val [2]int8
}

func (v int8RangeFromInt8Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type int8RangeToInt8Array2 struct {
	val *[2]int8
}

func (v int8RangeToInt8Array2) Scan(src interface{}) error {
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

type int8RangeFromInt16Array2 struct {
	val [2]int16
}

func (v int8RangeFromInt16Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type int8RangeToInt16Array2 struct {
	val *[2]int16
}

func (v int8RangeToInt16Array2) Scan(src interface{}) error {
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

type int8RangeFromInt32Array2 struct {
	val [2]int32
}

func (v int8RangeFromInt32Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type int8RangeToInt32Array2 struct {
	val *[2]int32
}

func (v int8RangeToInt32Array2) Scan(src interface{}) error {
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

type int8RangeFromInt64Array2 struct {
	val [2]int64
}

func (v int8RangeFromInt64Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, v.val[0], 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, v.val[1], 10)
	out = append(out, ')')
	return out, nil
}

type int8RangeToInt64Array2 struct {
	val *[2]int64
}

func (v int8RangeToInt64Array2) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	var lo, hi int64
	elems := pgParseRange(data)
	if len(elems[0]) > 0 {
		if lo, err = strconv.ParseInt(string(elems[0]), 10, 64); err != nil {
			return err
		}
	}
	if len(elems[1]) > 0 {
		if hi, err = strconv.ParseInt(string(elems[1]), 10, 64); err != nil {
			return err
		}
	}

	v.val[0] = lo
	v.val[1] = hi
	return nil
}

type int8RangeFromUintArray2 struct {
	val [2]uint
}

func (v int8RangeFromUintArray2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, uint64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, uint64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type int8RangeToUintArray2 struct {
	val *[2]uint
}

func (v int8RangeToUintArray2) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	var lo, hi uint64
	elems := pgParseRange(data)
	if len(elems[0]) > 0 {
		if lo, err = strconv.ParseUint(string(elems[0]), 10, 64); err != nil {
			return err
		}
	}
	if len(elems[1]) > 0 {
		if hi, err = strconv.ParseUint(string(elems[1]), 10, 64); err != nil {
			return err
		}
	}

	v.val[0] = uint(lo)
	v.val[1] = uint(hi)
	return nil
}

type int8RangeFromUint8Array2 struct {
	val [2]uint8
}

func (v int8RangeFromUint8Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, uint64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, uint64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type int8RangeToUint8Array2 struct {
	val *[2]uint8
}

func (v int8RangeToUint8Array2) Scan(src interface{}) error {
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

type int8RangeFromUint16Array2 struct {
	val [2]uint16
}

func (v int8RangeFromUint16Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, uint64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, uint64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type int8RangeToUint16Array2 struct {
	val *[2]uint16
}

func (v int8RangeToUint16Array2) Scan(src interface{}) error {
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

type int8RangeFromUint32Array2 struct {
	val [2]uint32
}

func (v int8RangeFromUint32Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, uint64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, uint64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type int8RangeToUint32Array2 struct {
	val *[2]uint32
}

func (v int8RangeToUint32Array2) Scan(src interface{}) error {
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

type int8RangeFromUint64Array2 struct {
	val [2]uint64
}

func (v int8RangeFromUint64Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, v.val[0], 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, v.val[1], 10)
	out = append(out, ')')
	return out, nil
}

type int8RangeToUint64Array2 struct {
	val *[2]uint64
}

func (v int8RangeToUint64Array2) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	var lo, hi uint64
	elems := pgParseRange(data)
	if len(elems[0]) > 0 {
		if lo, err = strconv.ParseUint(string(elems[0]), 10, 64); err != nil {
			return err
		}
	}
	if len(elems[1]) > 0 {
		if hi, err = strconv.ParseUint(string(elems[1]), 10, 64); err != nil {
			return err
		}
	}

	v.val[0] = lo
	v.val[1] = hi
	return nil
}

type int8RangeFromFloat32Array2 struct {
	val [2]float32
}

func (v int8RangeFromFloat32Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type int8RangeToFloat32Array2 struct {
	val *[2]float32
}

func (v int8RangeToFloat32Array2) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	var lo, hi int64
	elems := pgParseRange(data)
	if len(elems[0]) > 0 {
		if lo, err = strconv.ParseInt(string(elems[0]), 10, 64); err != nil {
			return err
		}
	}
	if len(elems[1]) > 0 {
		if hi, err = strconv.ParseInt(string(elems[1]), 10, 64); err != nil {
			return err
		}
	}

	v.val[0] = float32(lo)
	v.val[1] = float32(hi)
	return nil
}

type int8RangeFromFloat64Array2 struct {
	val [2]float64
}

func (v int8RangeFromFloat64Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type int8RangeToFloat64Array2 struct {
	val *[2]float64
}

func (v int8RangeToFloat64Array2) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	var lo, hi int64
	elems := pgParseRange(data)
	if len(elems[0]) > 0 {
		if lo, err = strconv.ParseInt(string(elems[0]), 10, 64); err != nil {
			return err
		}
	}
	if len(elems[1]) > 0 {
		if hi, err = strconv.ParseInt(string(elems[1]), 10, 64); err != nil {
			return err
		}
	}

	v.val[0] = float64(lo)
	v.val[1] = float64(hi)
	return nil
}
