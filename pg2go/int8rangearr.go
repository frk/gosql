package pg2go

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// Int8RangeArrayFromIntArray2Slice returns a driver.Valuer that produces a PostgreSQL int8range[] from the given Go [][2]int.
func Int8RangeArrayFromIntArray2Slice(val [][2]int) driver.Valuer {
	return int8RangeArrayFromIntArray2Slice{val: val}
}

// Int8RangeArrayToIntArray2Slice returns an sql.Scanner that converts a PostgreSQL int8range[] into a Go [][2]int and sets it to val.
func Int8RangeArrayToIntArray2Slice(val *[][2]int) sql.Scanner {
	return int8RangeArrayToIntArray2Slice{val: val}
}

// Int8RangeArrayFromInt8Array2Slice returns a driver.Valuer that produces a PostgreSQL int8range[] from the given Go [][2]int8.
func Int8RangeArrayFromInt8Array2Slice(val [][2]int8) driver.Valuer {
	return int8RangeArrayFromInt8Array2Slice{val: val}
}

// Int8RangeArrayToInt8Array2Slice returns an sql.Scanner that converts a PostgreSQL int8range[] into a Go [][2]int8 and sets it to val.
func Int8RangeArrayToInt8Array2Slice(val *[][2]int8) sql.Scanner {
	return int8RangeArrayToInt8Array2Slice{val: val}
}

// Int8RangeArrayFromInt16Array2Slice returns a driver.Valuer that produces a PostgreSQL int8range[] from the given Go [][2]int16.
func Int8RangeArrayFromInt16Array2Slice(val [][2]int16) driver.Valuer {
	return int8RangeArrayFromInt16Array2Slice{val: val}
}

// Int8RangeArrayToInt16Array2Slice returns an sql.Scanner that converts a PostgreSQL int8range[] into a Go [][2]int16 and sets it to val.
func Int8RangeArrayToInt16Array2Slice(val *[][2]int16) sql.Scanner {
	return int8RangeArrayToInt16Array2Slice{val: val}
}

// Int8RangeArrayFromInt32Array2Slice returns a driver.Valuer that produces a PostgreSQL int8range[] from the given Go [][2]int32.
func Int8RangeArrayFromInt32Array2Slice(val [][2]int32) driver.Valuer {
	return int8RangeArrayFromInt32Array2Slice{val: val}
}

// Int8RangeArrayToInt32Array2Slice returns an sql.Scanner that converts a PostgreSQL int8range[] into a Go [][2]int32 and sets it to val.
func Int8RangeArrayToInt32Array2Slice(val *[][2]int32) sql.Scanner {
	return int8RangeArrayToInt32Array2Slice{val: val}
}

// Int8RangeArrayFromInt64Array2Slice returns a driver.Valuer that produces a PostgreSQL int8range[] from the given Go [][2]int64.
func Int8RangeArrayFromInt64Array2Slice(val [][2]int64) driver.Valuer {
	return int8RangeArrayFromInt64Array2Slice{val: val}
}

// Int8RangeArrayToInt64Array2Slice returns an sql.Scanner that converts a PostgreSQL int8range[] into a Go [][2]int64 and sets it to val.
func Int8RangeArrayToInt64Array2Slice(val *[][2]int64) sql.Scanner {
	return int8RangeArrayToInt64Array2Slice{val: val}
}

// Int8RangeArrayFromUintArray2Slice returns a driver.Valuer that produces a PostgreSQL int8range[] from the given Go [][2]uint.
func Int8RangeArrayFromUintArray2Slice(val [][2]uint) driver.Valuer {
	return int8RangeArrayFromUintArray2Slice{val: val}
}

// Int8RangeArrayToUintArray2Slice returns an sql.Scanner that converts a PostgreSQL int8range[] into a Go [][2]uint and sets it to val.
func Int8RangeArrayToUintArray2Slice(val *[][2]uint) sql.Scanner {
	return int8RangeArrayToUintArray2Slice{val: val}
}

// Int8RangeArrayFromUint8Array2Slice returns a driver.Valuer that produces a PostgreSQL int8range[] from the given Go [][2]uint8.
func Int8RangeArrayFromUint8Array2Slice(val [][2]uint8) driver.Valuer {
	return int8RangeArrayFromUint8Array2Slice{val: val}
}

// Int8RangeArrayToUint8Array2Slice returns an sql.Scanner that converts a PostgreSQL int8range[] into a Go [][2]uint8 and sets it to val.
func Int8RangeArrayToUint8Array2Slice(val *[][2]uint8) sql.Scanner {
	return int8RangeArrayToUint8Array2Slice{val: val}
}

// Int8RangeArrayFromUint16Array2Slice returns a driver.Valuer that produces a PostgreSQL int8range[] from the given Go [][2]uint16.
func Int8RangeArrayFromUint16Array2Slice(val [][2]uint16) driver.Valuer {
	return int8RangeArrayFromUint16Array2Slice{val: val}
}

// Int8RangeArrayToUint16Array2Slice returns an sql.Scanner that converts a PostgreSQL int8range[] into a Go [][2]uint16 and sets it to val.
func Int8RangeArrayToUint16Array2Slice(val *[][2]uint16) sql.Scanner {
	return int8RangeArrayToUint16Array2Slice{val: val}
}

// Int8RangeArrayFromUint32Array2Slice returns a driver.Valuer that produces a PostgreSQL int8range[] from the given Go [][2]uint32.
func Int8RangeArrayFromUint32Array2Slice(val [][2]uint32) driver.Valuer {
	return int8RangeArrayFromUint32Array2Slice{val: val}
}

// Int8RangeArrayToUint32Array2Slice returns an sql.Scanner that converts a PostgreSQL int8range[] into a Go [][2]uint32 and sets it to val.
func Int8RangeArrayToUint32Array2Slice(val *[][2]uint32) sql.Scanner {
	return int8RangeArrayToUint32Array2Slice{val: val}
}

// Int8RangeArrayFromUint64Array2Slice returns a driver.Valuer that produces a PostgreSQL int8range[] from the given Go [][2]uint64.
func Int8RangeArrayFromUint64Array2Slice(val [][2]uint64) driver.Valuer {
	return int8RangeArrayFromUint64Array2Slice{val: val}
}

// Int8RangeArrayToUint64Array2Slice returns an sql.Scanner that converts a PostgreSQL int8range[] into a Go [][2]uint64 and sets it to val.
func Int8RangeArrayToUint64Array2Slice(val *[][2]uint64) sql.Scanner {
	return int8RangeArrayToUint64Array2Slice{val: val}
}

// Int8RangeArrayFromFloat32Array2Slice returns a driver.Valuer that produces a PostgreSQL int8range[] from the given Go [][2]float32.
func Int8RangeArrayFromFloat32Array2Slice(val [][2]float32) driver.Valuer {
	return int8RangeArrayFromFloat32Array2Slice{val: val}
}

// Int8RangeArrayToFloat32Array2Slice returns an sql.Scanner that converts a PostgreSQL int8range[] into a Go [][2]float32 and sets it to val.
func Int8RangeArrayToFloat32Array2Slice(val *[][2]float32) sql.Scanner {
	return int8RangeArrayToFloat32Array2Slice{val: val}
}

// Int8RangeArrayFromFloat64Array2Slice returns a driver.Valuer that produces a PostgreSQL int8range[] from the given Go [][2]float64.
func Int8RangeArrayFromFloat64Array2Slice(val [][2]float64) driver.Valuer {
	return int8RangeArrayFromFloat64Array2Slice{val: val}
}

// Int8RangeArrayToFloat64Array2Slice returns an sql.Scanner that converts a PostgreSQL int8range[] into a Go [][2]float64 and sets it to val.
func Int8RangeArrayToFloat64Array2Slice(val *[][2]float64) sql.Scanner {
	return int8RangeArrayToFloat64Array2Slice{val: val}
}

type int8RangeArrayFromIntArray2Slice struct {
	val [][2]int
}

func (v int8RangeArrayFromIntArray2Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.val {
		out = append(out, '"', '[')
		out = strconv.AppendInt(out, int64(a[0]), 10)
		out = append(out, ',')
		out = strconv.AppendInt(out, int64(a[1]), 10)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type int8RangeArrayToIntArray2Slice struct {
	val *[][2]int
}

func (v int8RangeArrayToIntArray2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.val = nil
		return nil
	}

	elems := pgParseQuotedStringArray(data)
	ranges := make([][2]int, len(elems))

	for i, elem := range elems {
		var lo, hi int64
		arr := pgParseRange(elem)
		if len(arr[0]) > 0 {
			if lo, err = strconv.ParseInt(string(arr[0]), 10, 64); err != nil {
				return err
			}
		}
		if len(arr[1]) > 0 {
			if hi, err = strconv.ParseInt(string(arr[1]), 10, 64); err != nil {
				return err
			}
		}

		ranges[i][0] = int(lo)
		ranges[i][1] = int(hi)
	}

	*v.val = ranges
	return nil
}

type int8RangeArrayFromInt8Array2Slice struct {
	val [][2]int8
}

func (v int8RangeArrayFromInt8Array2Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.val {
		out = append(out, '"', '[')
		out = strconv.AppendInt(out, int64(a[0]), 10)
		out = append(out, ',')
		out = strconv.AppendInt(out, int64(a[1]), 10)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type int8RangeArrayToInt8Array2Slice struct {
	val *[][2]int8
}

func (v int8RangeArrayToInt8Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.val = nil
		return nil
	}

	elems := pgParseQuotedStringArray(data)
	ranges := make([][2]int8, len(elems))

	for i, elem := range elems {
		var lo, hi int64
		arr := pgParseRange(elem)
		if len(arr[0]) > 0 {
			if lo, err = strconv.ParseInt(string(arr[0]), 10, 8); err != nil {
				return err
			}
		}
		if len(arr[1]) > 0 {
			if hi, err = strconv.ParseInt(string(arr[1]), 10, 8); err != nil {
				return err
			}
		}

		ranges[i][0] = int8(lo)
		ranges[i][1] = int8(hi)
	}

	*v.val = ranges
	return nil
}

type int8RangeArrayFromInt16Array2Slice struct {
	val [][2]int16
}

func (v int8RangeArrayFromInt16Array2Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.val {
		out = append(out, '"', '[')
		out = strconv.AppendInt(out, int64(a[0]), 10)
		out = append(out, ',')
		out = strconv.AppendInt(out, int64(a[1]), 10)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type int8RangeArrayToInt16Array2Slice struct {
	val *[][2]int16
}

func (v int8RangeArrayToInt16Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.val = nil
		return nil
	}

	elems := pgParseQuotedStringArray(data)
	ranges := make([][2]int16, len(elems))

	for i, elem := range elems {
		var lo, hi int64
		arr := pgParseRange(elem)
		if len(arr[0]) > 0 {
			if lo, err = strconv.ParseInt(string(arr[0]), 10, 16); err != nil {
				return err
			}
		}
		if len(arr[1]) > 0 {
			if hi, err = strconv.ParseInt(string(arr[1]), 10, 16); err != nil {
				return err
			}
		}

		ranges[i][0] = int16(lo)
		ranges[i][1] = int16(hi)
	}

	*v.val = ranges
	return nil
}

type int8RangeArrayFromInt32Array2Slice struct {
	val [][2]int32
}

func (v int8RangeArrayFromInt32Array2Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.val {
		out = append(out, '"', '[')
		out = strconv.AppendInt(out, int64(a[0]), 10)
		out = append(out, ',')
		out = strconv.AppendInt(out, int64(a[1]), 10)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type int8RangeArrayToInt32Array2Slice struct {
	val *[][2]int32
}

func (v int8RangeArrayToInt32Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.val = nil
		return nil
	}

	elems := pgParseQuotedStringArray(data)
	ranges := make([][2]int32, len(elems))

	for i, elem := range elems {
		var lo, hi int64
		arr := pgParseRange(elem)
		if len(arr[0]) > 0 {
			if lo, err = strconv.ParseInt(string(arr[0]), 10, 32); err != nil {
				return err
			}
		}
		if len(arr[1]) > 0 {
			if hi, err = strconv.ParseInt(string(arr[1]), 10, 32); err != nil {
				return err
			}
		}

		ranges[i][0] = int32(lo)
		ranges[i][1] = int32(hi)
	}

	*v.val = ranges
	return nil
}

type int8RangeArrayFromInt64Array2Slice struct {
	val [][2]int64
}

func (v int8RangeArrayFromInt64Array2Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.val {
		out = append(out, '"', '[')
		out = strconv.AppendInt(out, int64(a[0]), 10)
		out = append(out, ',')
		out = strconv.AppendInt(out, int64(a[1]), 10)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type int8RangeArrayToInt64Array2Slice struct {
	val *[][2]int64
}

func (v int8RangeArrayToInt64Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.val = nil
		return nil
	}

	elems := pgParseQuotedStringArray(data)
	ranges := make([][2]int64, len(elems))

	for i, elem := range elems {
		var lo, hi int64
		arr := pgParseRange(elem)
		if len(arr[0]) > 0 {
			if lo, err = strconv.ParseInt(string(arr[0]), 10, 64); err != nil {
				return err
			}
		}
		if len(arr[1]) > 0 {
			if hi, err = strconv.ParseInt(string(arr[1]), 10, 64); err != nil {
				return err
			}
		}

		ranges[i][0] = lo
		ranges[i][1] = hi
	}

	*v.val = ranges
	return nil
}

type int8RangeArrayFromUintArray2Slice struct {
	val [][2]uint
}

func (v int8RangeArrayFromUintArray2Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.val {
		out = append(out, '"', '[')
		out = strconv.AppendUint(out, uint64(a[0]), 10)
		out = append(out, ',')
		out = strconv.AppendUint(out, uint64(a[1]), 10)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type int8RangeArrayToUintArray2Slice struct {
	val *[][2]uint
}

func (v int8RangeArrayToUintArray2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.val = nil
		return nil
	}

	elems := pgParseQuotedStringArray(data)
	ranges := make([][2]uint, len(elems))

	for i, elem := range elems {
		var lo, hi uint64
		arr := pgParseRange(elem)
		if len(arr[0]) > 0 {
			if lo, err = strconv.ParseUint(string(arr[0]), 10, 64); err != nil {
				return err
			}
		}
		if len(arr[1]) > 0 {
			if hi, err = strconv.ParseUint(string(arr[1]), 10, 64); err != nil {
				return err
			}
		}

		ranges[i][0] = uint(lo)
		ranges[i][1] = uint(hi)
	}

	*v.val = ranges
	return nil
}

type int8RangeArrayFromUint8Array2Slice struct {
	val [][2]uint8
}

func (v int8RangeArrayFromUint8Array2Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.val {
		out = append(out, '"', '[')
		out = strconv.AppendUint(out, uint64(a[0]), 10)
		out = append(out, ',')
		out = strconv.AppendUint(out, uint64(a[1]), 10)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type int8RangeArrayToUint8Array2Slice struct {
	val *[][2]uint8
}

func (v int8RangeArrayToUint8Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.val = nil
		return nil
	}

	elems := pgParseQuotedStringArray(data)
	ranges := make([][2]uint8, len(elems))

	for i, elem := range elems {
		var lo, hi uint64
		arr := pgParseRange(elem)
		if len(arr[0]) > 0 {
			if lo, err = strconv.ParseUint(string(arr[0]), 10, 8); err != nil {
				return err
			}
		}
		if len(arr[1]) > 0 {
			if hi, err = strconv.ParseUint(string(arr[1]), 10, 8); err != nil {
				return err
			}
		}

		ranges[i][0] = uint8(lo)
		ranges[i][1] = uint8(hi)
	}

	*v.val = ranges
	return nil
}

type int8RangeArrayFromUint16Array2Slice struct {
	val [][2]uint16
}

func (v int8RangeArrayFromUint16Array2Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.val {
		out = append(out, '"', '[')
		out = strconv.AppendUint(out, uint64(a[0]), 10)
		out = append(out, ',')
		out = strconv.AppendUint(out, uint64(a[1]), 10)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type int8RangeArrayToUint16Array2Slice struct {
	val *[][2]uint16
}

func (v int8RangeArrayToUint16Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.val = nil
		return nil
	}

	elems := pgParseQuotedStringArray(data)
	ranges := make([][2]uint16, len(elems))

	for i, elem := range elems {
		var lo, hi uint64
		arr := pgParseRange(elem)
		if len(arr[0]) > 0 {
			if lo, err = strconv.ParseUint(string(arr[0]), 10, 16); err != nil {
				return err
			}
		}
		if len(arr[1]) > 0 {
			if hi, err = strconv.ParseUint(string(arr[1]), 10, 16); err != nil {
				return err
			}
		}

		ranges[i][0] = uint16(lo)
		ranges[i][1] = uint16(hi)
	}

	*v.val = ranges
	return nil
}

type int8RangeArrayFromUint32Array2Slice struct {
	val [][2]uint32
}

func (v int8RangeArrayFromUint32Array2Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.val {
		out = append(out, '"', '[')
		out = strconv.AppendUint(out, uint64(a[0]), 10)
		out = append(out, ',')
		out = strconv.AppendUint(out, uint64(a[1]), 10)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type int8RangeArrayToUint32Array2Slice struct {
	val *[][2]uint32
}

func (v int8RangeArrayToUint32Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.val = nil
		return nil
	}

	elems := pgParseQuotedStringArray(data)
	ranges := make([][2]uint32, len(elems))

	for i, elem := range elems {
		var lo, hi uint64
		arr := pgParseRange(elem)
		if len(arr[0]) > 0 {
			if lo, err = strconv.ParseUint(string(arr[0]), 10, 32); err != nil {
				return err
			}
		}
		if len(arr[1]) > 0 {
			if hi, err = strconv.ParseUint(string(arr[1]), 10, 32); err != nil {
				return err
			}
		}

		ranges[i][0] = uint32(lo)
		ranges[i][1] = uint32(hi)
	}

	*v.val = ranges
	return nil
}

type int8RangeArrayFromUint64Array2Slice struct {
	val [][2]uint64
}

func (v int8RangeArrayFromUint64Array2Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.val {
		out = append(out, '"', '[')
		out = strconv.AppendUint(out, uint64(a[0]), 10)
		out = append(out, ',')
		out = strconv.AppendUint(out, uint64(a[1]), 10)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type int8RangeArrayToUint64Array2Slice struct {
	val *[][2]uint64
}

func (v int8RangeArrayToUint64Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.val = nil
		return nil
	}

	elems := pgParseQuotedStringArray(data)
	ranges := make([][2]uint64, len(elems))

	for i, elem := range elems {
		var lo, hi uint64
		arr := pgParseRange(elem)
		if len(arr[0]) > 0 {
			if lo, err = strconv.ParseUint(string(arr[0]), 10, 64); err != nil {
				return err
			}
		}
		if len(arr[1]) > 0 {
			if hi, err = strconv.ParseUint(string(arr[1]), 10, 64); err != nil {
				return err
			}
		}

		ranges[i][0] = lo
		ranges[i][1] = hi
	}

	*v.val = ranges
	return nil
}

type int8RangeArrayFromFloat32Array2Slice struct {
	val [][2]float32
}

func (v int8RangeArrayFromFloat32Array2Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.val {
		out = append(out, '"', '[')
		out = strconv.AppendInt(out, int64(a[0]), 10)
		out = append(out, ',')
		out = strconv.AppendInt(out, int64(a[1]), 10)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type int8RangeArrayToFloat32Array2Slice struct {
	val *[][2]float32
}

func (v int8RangeArrayToFloat32Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.val = nil
		return nil
	}

	elems := pgParseQuotedStringArray(data)
	ranges := make([][2]float32, len(elems))

	for i, elem := range elems {
		var lo, hi int64
		arr := pgParseRange(elem)
		if len(arr[0]) > 0 {
			if lo, err = strconv.ParseInt(string(arr[0]), 10, 32); err != nil {
				return err
			}
		}
		if len(arr[1]) > 0 {
			if hi, err = strconv.ParseInt(string(arr[1]), 10, 32); err != nil {
				return err
			}
		}

		ranges[i][0] = float32(lo)
		ranges[i][1] = float32(hi)
	}

	*v.val = ranges
	return nil
}

type int8RangeArrayFromFloat64Array2Slice struct {
	val [][2]float64
}

func (v int8RangeArrayFromFloat64Array2Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.val {
		out = append(out, '"', '[')
		out = strconv.AppendInt(out, int64(a[0]), 10)
		out = append(out, ',')
		out = strconv.AppendInt(out, int64(a[1]), 10)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type int8RangeArrayToFloat64Array2Slice struct {
	val *[][2]float64
}

func (v int8RangeArrayToFloat64Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.val = nil
		return nil
	}

	elems := pgParseQuotedStringArray(data)
	ranges := make([][2]float64, len(elems))

	for i, elem := range elems {
		var lo, hi int64
		arr := pgParseRange(elem)
		if len(arr[0]) > 0 {
			if lo, err = strconv.ParseInt(string(arr[0]), 10, 64); err != nil {
				return err
			}
		}
		if len(arr[1]) > 0 {
			if hi, err = strconv.ParseInt(string(arr[1]), 10, 64); err != nil {
				return err
			}
		}

		ranges[i][0] = float64(lo)
		ranges[i][1] = float64(hi)
	}

	*v.val = ranges
	return nil
}
