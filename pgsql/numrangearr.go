package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// NumRangeArrayFromIntArray2Slice returns a driver.Valuer that produces a PostgreSQL numrange[] from the given Go [][2]int.
func NumRangeArrayFromIntArray2Slice(val [][2]int) driver.Valuer {
	return numRangeArrayFromIntArray2Slice{val: val}
}

// NumRangeArrayToIntArray2Slice returns an sql.Scanner that converts a PostgreSQL numrange[] into a Go [][2]int and sets it to val.
func NumRangeArrayToIntArray2Slice(val *[][2]int) sql.Scanner {
	return numRangeArrayToIntArray2Slice{val: val}
}

// NumRangeArrayFromInt8Array2Slice returns a driver.Valuer that produces a PostgreSQL numrange[] from the given Go [][2]int8.
func NumRangeArrayFromInt8Array2Slice(val [][2]int8) driver.Valuer {
	return numRangeArrayFromInt8Array2Slice{val: val}
}

// NumRangeArrayToInt8Array2Slice returns an sql.Scanner that converts a PostgreSQL numrange[] into a Go [][2]int8 and sets it to val.
func NumRangeArrayToInt8Array2Slice(val *[][2]int8) sql.Scanner {
	return numRangeArrayToInt8Array2Slice{val: val}
}

// NumRangeArrayFromInt16Array2Slice returns a driver.Valuer that produces a PostgreSQL numrange[] from the given Go [][2]int16.
func NumRangeArrayFromInt16Array2Slice(val [][2]int16) driver.Valuer {
	return numRangeArrayFromInt16Array2Slice{val: val}
}

// NumRangeArrayToInt16Array2Slice returns an sql.Scanner that converts a PostgreSQL numrange[] into a Go [][2]int16 and sets it to val.
func NumRangeArrayToInt16Array2Slice(val *[][2]int16) sql.Scanner {
	return numRangeArrayToInt16Array2Slice{val: val}
}

// NumRangeArrayFromInt32Array2Slice returns a driver.Valuer that produces a PostgreSQL numrange[] from the given Go [][2]int32.
func NumRangeArrayFromInt32Array2Slice(val [][2]int32) driver.Valuer {
	return numRangeArrayFromInt32Array2Slice{val: val}
}

// NumRangeArrayToInt32Array2Slice returns an sql.Scanner that converts a PostgreSQL numrange[] into a Go [][2]int32 and sets it to val.
func NumRangeArrayToInt32Array2Slice(val *[][2]int32) sql.Scanner {
	return numRangeArrayToInt32Array2Slice{val: val}
}

// NumRangeArrayFromInt64Array2Slice returns a driver.Valuer that produces a PostgreSQL numrange[] from the given Go [][2]int64.
func NumRangeArrayFromInt64Array2Slice(val [][2]int64) driver.Valuer {
	return numRangeArrayFromInt64Array2Slice{val: val}
}

// NumRangeArrayToInt64Array2Slice returns an sql.Scanner that converts a PostgreSQL numrange[] into a Go [][2]int64 and sets it to val.
func NumRangeArrayToInt64Array2Slice(val *[][2]int64) sql.Scanner {
	return numRangeArrayToInt64Array2Slice{val: val}
}

// NumRangeArrayFromUintArray2Slice returns a driver.Valuer that produces a PostgreSQL numrange[] from the given Go [][2]uint.
func NumRangeArrayFromUintArray2Slice(val [][2]uint) driver.Valuer {
	return numRangeArrayFromUintArray2Slice{val: val}
}

// NumRangeArrayToUintArray2Slice returns an sql.Scanner that converts a PostgreSQL numrange[] into a Go [][2]uint and sets it to val.
func NumRangeArrayToUintArray2Slice(val *[][2]uint) sql.Scanner {
	return numRangeArrayToUintArray2Slice{val: val}
}

// NumRangeArrayFromUint8Array2Slice returns a driver.Valuer that produces a PostgreSQL numrange[] from the given Go [][2]uint8.
func NumRangeArrayFromUint8Array2Slice(val [][2]uint8) driver.Valuer {
	return numRangeArrayFromUint8Array2Slice{val: val}
}

// NumRangeArrayToUint8Array2Slice returns an sql.Scanner that converts a PostgreSQL numrange[] into a Go [][2]uint8 and sets it to val.
func NumRangeArrayToUint8Array2Slice(val *[][2]uint8) sql.Scanner {
	return numRangeArrayToUint8Array2Slice{val: val}
}

// NumRangeArrayFromUint16Array2Slice returns a driver.Valuer that produces a PostgreSQL numrange[] from the given Go [][2]uint16.
func NumRangeArrayFromUint16Array2Slice(val [][2]uint16) driver.Valuer {
	return numRangeArrayFromUint16Array2Slice{val: val}
}

// NumRangeArrayToUint16Array2Slice returns an sql.Scanner that converts a PostgreSQL numrange[] into a Go [][2]uint16 and sets it to val.
func NumRangeArrayToUint16Array2Slice(val *[][2]uint16) sql.Scanner {
	return numRangeArrayToUint16Array2Slice{val: val}
}

// NumRangeArrayFromUint32Array2Slice returns a driver.Valuer that produces a PostgreSQL numrange[] from the given Go [][2]uint32.
func NumRangeArrayFromUint32Array2Slice(val [][2]uint32) driver.Valuer {
	return numRangeArrayFromUint32Array2Slice{val: val}
}

// NumRangeArrayToUint32Array2Slice returns an sql.Scanner that converts a PostgreSQL numrange[] into a Go [][2]uint32 and sets it to val.
func NumRangeArrayToUint32Array2Slice(val *[][2]uint32) sql.Scanner {
	return numRangeArrayToUint32Array2Slice{val: val}
}

// NumRangeArrayFromUint64Array2Slice returns a driver.Valuer that produces a PostgreSQL numrange[] from the given Go [][2]uint64.
func NumRangeArrayFromUint64Array2Slice(val [][2]uint64) driver.Valuer {
	return numRangeArrayFromUint64Array2Slice{val: val}
}

// NumRangeArrayToUint64Array2Slice returns an sql.Scanner that converts a PostgreSQL numrange[] into a Go [][2]uint64 and sets it to val.
func NumRangeArrayToUint64Array2Slice(val *[][2]uint64) sql.Scanner {
	return numRangeArrayToUint64Array2Slice{val: val}
}

// NumRangeArrayFromFloat32Array2Slice returns a driver.Valuer that produces a PostgreSQL numrange[] from the given Go [][2]float32.
func NumRangeArrayFromFloat32Array2Slice(val [][2]float32) driver.Valuer {
	return numRangeArrayFromFloat32Array2Slice{val: val}
}

// NumRangeArrayToFloat32Array2Slice returns an sql.Scanner that converts a PostgreSQL numrange[] into a Go [][2]float32 and sets it to val.
func NumRangeArrayToFloat32Array2Slice(val *[][2]float32) sql.Scanner {
	return numRangeArrayToFloat32Array2Slice{val: val}
}

// NumRangeArrayFromFloat64Array2Slice returns a driver.Valuer that produces a PostgreSQL numrange[] from the given Go [][2]float64.
func NumRangeArrayFromFloat64Array2Slice(val [][2]float64) driver.Valuer {
	return numRangeArrayFromFloat64Array2Slice{val: val}
}

// NumRangeArrayToFloat64Array2Slice returns an sql.Scanner that converts a PostgreSQL numrange[] into a Go [][2]float64 and sets it to val.
func NumRangeArrayToFloat64Array2Slice(val *[][2]float64) sql.Scanner {
	return numRangeArrayToFloat64Array2Slice{val: val}
}

type numRangeArrayFromIntArray2Slice struct {
	val [][2]int
}

func (v numRangeArrayFromIntArray2Slice) Value() (driver.Value, error) {
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

type numRangeArrayToIntArray2Slice struct {
	val *[][2]int
}

func (v numRangeArrayToIntArray2Slice) Scan(src interface{}) error {
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

type numRangeArrayFromInt8Array2Slice struct {
	val [][2]int8
}

func (v numRangeArrayFromInt8Array2Slice) Value() (driver.Value, error) {
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

type numRangeArrayToInt8Array2Slice struct {
	val *[][2]int8
}

func (v numRangeArrayToInt8Array2Slice) Scan(src interface{}) error {
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

type numRangeArrayFromInt16Array2Slice struct {
	val [][2]int16
}

func (v numRangeArrayFromInt16Array2Slice) Value() (driver.Value, error) {
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

type numRangeArrayToInt16Array2Slice struct {
	val *[][2]int16
}

func (v numRangeArrayToInt16Array2Slice) Scan(src interface{}) error {
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

type numRangeArrayFromInt32Array2Slice struct {
	val [][2]int32
}

func (v numRangeArrayFromInt32Array2Slice) Value() (driver.Value, error) {
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

type numRangeArrayToInt32Array2Slice struct {
	val *[][2]int32
}

func (v numRangeArrayToInt32Array2Slice) Scan(src interface{}) error {
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

type numRangeArrayFromInt64Array2Slice struct {
	val [][2]int64
}

func (v numRangeArrayFromInt64Array2Slice) Value() (driver.Value, error) {
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

type numRangeArrayToInt64Array2Slice struct {
	val *[][2]int64
}

func (v numRangeArrayToInt64Array2Slice) Scan(src interface{}) error {
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

type numRangeArrayFromUintArray2Slice struct {
	val [][2]uint
}

func (v numRangeArrayFromUintArray2Slice) Value() (driver.Value, error) {
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

type numRangeArrayToUintArray2Slice struct {
	val *[][2]uint
}

func (v numRangeArrayToUintArray2Slice) Scan(src interface{}) error {
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

type numRangeArrayFromUint8Array2Slice struct {
	val [][2]uint8
}

func (v numRangeArrayFromUint8Array2Slice) Value() (driver.Value, error) {
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

type numRangeArrayToUint8Array2Slice struct {
	val *[][2]uint8
}

func (v numRangeArrayToUint8Array2Slice) Scan(src interface{}) error {
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

type numRangeArrayFromUint16Array2Slice struct {
	val [][2]uint16
}

func (v numRangeArrayFromUint16Array2Slice) Value() (driver.Value, error) {
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

type numRangeArrayToUint16Array2Slice struct {
	val *[][2]uint16
}

func (v numRangeArrayToUint16Array2Slice) Scan(src interface{}) error {
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

type numRangeArrayFromUint32Array2Slice struct {
	val [][2]uint32
}

func (v numRangeArrayFromUint32Array2Slice) Value() (driver.Value, error) {
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

type numRangeArrayToUint32Array2Slice struct {
	val *[][2]uint32
}

func (v numRangeArrayToUint32Array2Slice) Scan(src interface{}) error {
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

type numRangeArrayFromUint64Array2Slice struct {
	val [][2]uint64
}

func (v numRangeArrayFromUint64Array2Slice) Value() (driver.Value, error) {
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

type numRangeArrayToUint64Array2Slice struct {
	val *[][2]uint64
}

func (v numRangeArrayToUint64Array2Slice) Scan(src interface{}) error {
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

type numRangeArrayFromFloat32Array2Slice struct {
	val [][2]float32
}

func (v numRangeArrayFromFloat32Array2Slice) Value() (driver.Value, error) {
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

type numRangeArrayToFloat32Array2Slice struct {
	val *[][2]float32
}

func (v numRangeArrayToFloat32Array2Slice) Scan(src interface{}) error {
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

type numRangeArrayFromFloat64Array2Slice struct {
	val [][2]float64
}

func (v numRangeArrayFromFloat64Array2Slice) Value() (driver.Value, error) {
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

type numRangeArrayToFloat64Array2Slice struct {
	val *[][2]float64
}

func (v numRangeArrayToFloat64Array2Slice) Scan(src interface{}) error {
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
