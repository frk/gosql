package pg2go

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// NumRangeFromIntArray2 returns a driver.Valuer that produces a PostgreSQL numrange from the given Go [2]int.
func NumRangeFromIntArray2(val [2]int) driver.Valuer {
	return numRangeFromIntArray2{val: val}
}

// NumRangeToIntArray2 returns an sql.Scanner that converts a PostgreSQL numrange into a Go [2]int and sets it to val.
func NumRangeToIntArray2(val *[2]int) sql.Scanner {
	return numRangeToIntArray2{val: val}
}

// NumRangeFromInt8Array2 returns a driver.Valuer that produces a PostgreSQL numrange from the given Go [2]int8.
func NumRangeFromInt8Array2(val [2]int8) driver.Valuer {
	return numRangeFromInt8Array2{val: val}
}

// NumRangeToInt8Array2 returns an sql.Scanner that converts a PostgreSQL numrange into a Go [2]int8 and sets it to val.
func NumRangeToInt8Array2(val *[2]int8) sql.Scanner {
	return numRangeToInt8Array2{val: val}
}

// NumRangeFromInt16Array2 returns a driver.Valuer that produces a PostgreSQL numrange from the given Go [2]int16.
func NumRangeFromInt16Array2(val [2]int16) driver.Valuer {
	return numRangeFromInt16Array2{val: val}
}

// NumRangeToInt16Array2 returns an sql.Scanner that converts a PostgreSQL numrange into a Go [2]int16 and sets it to val.
func NumRangeToInt16Array2(val *[2]int16) sql.Scanner {
	return numRangeToInt16Array2{val: val}
}

// NumRangeFromInt32Array2 returns a driver.Valuer that produces a PostgreSQL numrange from the given Go [2]int32.
func NumRangeFromInt32Array2(val [2]int32) driver.Valuer {
	return numRangeFromInt32Array2{val: val}
}

// NumRangeToInt32Array2 returns an sql.Scanner that converts a PostgreSQL numrange into a Go [2]int32 and sets it to val.
func NumRangeToInt32Array2(val *[2]int32) sql.Scanner {
	return numRangeToInt32Array2{val: val}
}

// NumRangeFromInt64Array2 returns a driver.Valuer that produces a PostgreSQL numrange from the given Go [2]int64.
func NumRangeFromInt64Array2(val [2]int64) driver.Valuer {
	return numRangeFromInt64Array2{val: val}
}

// NumRangeToInt64Array2 returns an sql.Scanner that converts a PostgreSQL numrange into a Go [2]int64 and sets it to val.
func NumRangeToInt64Array2(val *[2]int64) sql.Scanner {
	return numRangeToInt64Array2{val: val}
}

// NumRangeFromUintArray2 returns a driver.Valuer that produces a PostgreSQL numrange from the given Go [2]uint.
func NumRangeFromUintArray2(val [2]uint) driver.Valuer {
	return numRangeFromUintArray2{val: val}
}

// NumRangeToUintArray2 returns an sql.Scanner that converts a PostgreSQL numrange into a Go [2]uint and sets it to val.
func NumRangeToUintArray2(val *[2]uint) sql.Scanner {
	return numRangeToUintArray2{val: val}
}

// NumRangeFromUint8Array2 returns a driver.Valuer that produces a PostgreSQL numrange from the given Go [2]uint8.
func NumRangeFromUint8Array2(val [2]uint8) driver.Valuer {
	return numRangeFromUint8Array2{val: val}
}

// NumRangeToUint8Array2 returns an sql.Scanner that converts a PostgreSQL numrange into a Go [2]uint8 and sets it to val.
func NumRangeToUint8Array2(val *[2]uint8) sql.Scanner {
	return numRangeToUint8Array2{val: val}
}

// NumRangeFromUint16Array2 returns a driver.Valuer that produces a PostgreSQL numrange from the given Go [2]uint16.
func NumRangeFromUint16Array2(val [2]uint16) driver.Valuer {
	return numRangeFromUint16Array2{val: val}
}

// NumRangeToUint16Array2 returns an sql.Scanner that converts a PostgreSQL numrange into a Go [2]uint16 and sets it to val.
func NumRangeToUint16Array2(val *[2]uint16) sql.Scanner {
	return numRangeToUint16Array2{val: val}
}

// NumRangeFromUint32Array2 returns a driver.Valuer that produces a PostgreSQL numrange from the given Go [2]uint32.
func NumRangeFromUint32Array2(val [2]uint32) driver.Valuer {
	return numRangeFromUint32Array2{val: val}
}

// NumRangeToUint32Array2 returns an sql.Scanner that converts a PostgreSQL numrange into a Go [2]uint32 and sets it to val.
func NumRangeToUint32Array2(val *[2]uint32) sql.Scanner {
	return numRangeToUint32Array2{val: val}
}

// NumRangeFromUint64Array2 returns a driver.Valuer that produces a PostgreSQL numrange from the given Go [2]uint64.
func NumRangeFromUint64Array2(val [2]uint64) driver.Valuer {
	return numRangeFromUint64Array2{val: val}
}

// NumRangeToUint64Array2 returns an sql.Scanner that converts a PostgreSQL numrange into a Go [2]uint64 and sets it to val.
func NumRangeToUint64Array2(val *[2]uint64) sql.Scanner {
	return numRangeToUint64Array2{val: val}
}

// NumRangeFromFloat32Array2 returns a driver.Valuer that produces a PostgreSQL numrange from the given Go [2]float32.
func NumRangeFromFloat32Array2(val [2]float32) driver.Valuer {
	return numRangeFromFloat32Array2{val: val}
}

// NumRangeToFloat32Array2 returns an sql.Scanner that converts a PostgreSQL numrange into a Go [2]float32 and sets it to val.
func NumRangeToFloat32Array2(val *[2]float32) sql.Scanner {
	return numRangeToFloat32Array2{val: val}
}

// NumRangeFromFloat64Array2 returns a driver.Valuer that produces a PostgreSQL numrange from the given Go [2]float64.
func NumRangeFromFloat64Array2(val [2]float64) driver.Valuer {
	return numRangeFromFloat64Array2{val: val}
}

// NumRangeToFloat64Array2 returns an sql.Scanner that converts a PostgreSQL numrange into a Go [2]float64 and sets it to val.
func NumRangeToFloat64Array2(val *[2]float64) sql.Scanner {
	return numRangeToFloat64Array2{val: val}
}

type numRangeFromIntArray2 struct {
	val [2]int
}

func (v numRangeFromIntArray2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type numRangeToIntArray2 struct {
	val *[2]int
}

func (v numRangeToIntArray2) Scan(src interface{}) error {
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

type numRangeFromInt8Array2 struct {
	val [2]int8
}

func (v numRangeFromInt8Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type numRangeToInt8Array2 struct {
	val *[2]int8
}

func (v numRangeToInt8Array2) Scan(src interface{}) error {
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

type numRangeFromInt16Array2 struct {
	val [2]int16
}

func (v numRangeFromInt16Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type numRangeToInt16Array2 struct {
	val *[2]int16
}

func (v numRangeToInt16Array2) Scan(src interface{}) error {
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

type numRangeFromInt32Array2 struct {
	val [2]int32
}

func (v numRangeFromInt32Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type numRangeToInt32Array2 struct {
	val *[2]int32
}

func (v numRangeToInt32Array2) Scan(src interface{}) error {
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

type numRangeFromInt64Array2 struct {
	val [2]int64
}

func (v numRangeFromInt64Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, v.val[0], 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, v.val[1], 10)
	out = append(out, ')')
	return out, nil
}

type numRangeToInt64Array2 struct {
	val *[2]int64
}

func (v numRangeToInt64Array2) Scan(src interface{}) error {
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

type numRangeFromUintArray2 struct {
	val [2]uint
}

func (v numRangeFromUintArray2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, uint64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, uint64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type numRangeToUintArray2 struct {
	val *[2]uint
}

func (v numRangeToUintArray2) Scan(src interface{}) error {
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

type numRangeFromUint8Array2 struct {
	val [2]uint8
}

func (v numRangeFromUint8Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, uint64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, uint64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type numRangeToUint8Array2 struct {
	val *[2]uint8
}

func (v numRangeToUint8Array2) Scan(src interface{}) error {
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

type numRangeFromUint16Array2 struct {
	val [2]uint16
}

func (v numRangeFromUint16Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, uint64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, uint64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type numRangeToUint16Array2 struct {
	val *[2]uint16
}

func (v numRangeToUint16Array2) Scan(src interface{}) error {
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

type numRangeFromUint32Array2 struct {
	val [2]uint32
}

func (v numRangeFromUint32Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, uint64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, uint64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type numRangeToUint32Array2 struct {
	val *[2]uint32
}

func (v numRangeToUint32Array2) Scan(src interface{}) error {
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

type numRangeFromUint64Array2 struct {
	val [2]uint64
}

func (v numRangeFromUint64Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, v.val[0], 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, v.val[1], 10)
	out = append(out, ')')
	return out, nil
}

type numRangeToUint64Array2 struct {
	val *[2]uint64
}

func (v numRangeToUint64Array2) Scan(src interface{}) error {
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

type numRangeFromFloat32Array2 struct {
	val [2]float32
}

func (v numRangeFromFloat32Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type numRangeToFloat32Array2 struct {
	val *[2]float32
}

func (v numRangeToFloat32Array2) Scan(src interface{}) error {
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

type numRangeFromFloat64Array2 struct {
	val [2]float64
}

func (v numRangeFromFloat64Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type numRangeToFloat64Array2 struct {
	val *[2]float64
}

func (v numRangeToFloat64Array2) Scan(src interface{}) error {
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
