package convert

import (
	"database/sql/driver"
	"strconv"
)

type Int4RangeArrayFromIntArray2Slice struct {
	Val [][2]int
}

func (v Int4RangeArrayFromIntArray2Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.Val {
		out = append(out, '"', '[')
		out = strconv.AppendInt(out, int64(a[0]), 10)
		out = append(out, ',')
		out = strconv.AppendInt(out, int64(a[1]), 10)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Int4RangeArrayToIntArray2Slice struct {
	Val *[][2]int
}

func (v Int4RangeArrayToIntArray2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
		return nil
	}

	elems := pgParseQuotedStringArray(data)
	ranges := make([][2]int, len(elems))

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

		ranges[i][0] = int(lo)
		ranges[i][1] = int(hi)
	}

	*v.Val = ranges
	return nil
}

type Int4RangeArrayFromInt8Array2Slice struct {
	Val [][2]int8
}

func (v Int4RangeArrayFromInt8Array2Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.Val {
		out = append(out, '"', '[')
		out = strconv.AppendInt(out, int64(a[0]), 10)
		out = append(out, ',')
		out = strconv.AppendInt(out, int64(a[1]), 10)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Int4RangeArrayToInt8Array2Slice struct {
	Val *[][2]int8
}

func (v Int4RangeArrayToInt8Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
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

	*v.Val = ranges
	return nil
}

type Int4RangeArrayFromInt16Array2Slice struct {
	Val [][2]int16
}

func (v Int4RangeArrayFromInt16Array2Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.Val {
		out = append(out, '"', '[')
		out = strconv.AppendInt(out, int64(a[0]), 10)
		out = append(out, ',')
		out = strconv.AppendInt(out, int64(a[1]), 10)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Int4RangeArrayToInt16Array2Slice struct {
	Val *[][2]int16
}

func (v Int4RangeArrayToInt16Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
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

	*v.Val = ranges
	return nil
}

type Int4RangeArrayFromInt32Array2Slice struct {
	Val [][2]int32
}

func (v Int4RangeArrayFromInt32Array2Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.Val {
		out = append(out, '"', '[')
		out = strconv.AppendInt(out, int64(a[0]), 10)
		out = append(out, ',')
		out = strconv.AppendInt(out, int64(a[1]), 10)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Int4RangeArrayToInt32Array2Slice struct {
	Val *[][2]int32
}

func (v Int4RangeArrayToInt32Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
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

	*v.Val = ranges
	return nil
}

type Int4RangeArrayFromInt64Array2Slice struct {
	Val [][2]int64
}

func (v Int4RangeArrayFromInt64Array2Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.Val {
		out = append(out, '"', '[')
		out = strconv.AppendInt(out, int64(a[0]), 10)
		out = append(out, ',')
		out = strconv.AppendInt(out, int64(a[1]), 10)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Int4RangeArrayToInt64Array2Slice struct {
	Val *[][2]int64
}

func (v Int4RangeArrayToInt64Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
		return nil
	}

	elems := pgParseQuotedStringArray(data)
	ranges := make([][2]int64, len(elems))

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

		ranges[i][0] = lo
		ranges[i][1] = hi
	}

	*v.Val = ranges
	return nil
}

type Int4RangeArrayFromUintArray2Slice struct {
	Val [][2]uint
}

func (v Int4RangeArrayFromUintArray2Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.Val {
		out = append(out, '"', '[')
		out = strconv.AppendUint(out, uint64(a[0]), 10)
		out = append(out, ',')
		out = strconv.AppendUint(out, uint64(a[1]), 10)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Int4RangeArrayToUintArray2Slice struct {
	Val *[][2]uint
}

func (v Int4RangeArrayToUintArray2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
		return nil
	}

	elems := pgParseQuotedStringArray(data)
	ranges := make([][2]uint, len(elems))

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

		ranges[i][0] = uint(lo)
		ranges[i][1] = uint(hi)
	}

	*v.Val = ranges
	return nil
}

type Int4RangeArrayFromUint8Array2Slice struct {
	Val [][2]uint8
}

func (v Int4RangeArrayFromUint8Array2Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.Val {
		out = append(out, '"', '[')
		out = strconv.AppendUint(out, uint64(a[0]), 10)
		out = append(out, ',')
		out = strconv.AppendUint(out, uint64(a[1]), 10)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Int4RangeArrayToUint8Array2Slice struct {
	Val *[][2]uint8
}

func (v Int4RangeArrayToUint8Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
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

	*v.Val = ranges
	return nil
}

type Int4RangeArrayFromUint16Array2Slice struct {
	Val [][2]uint16
}

func (v Int4RangeArrayFromUint16Array2Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.Val {
		out = append(out, '"', '[')
		out = strconv.AppendUint(out, uint64(a[0]), 10)
		out = append(out, ',')
		out = strconv.AppendUint(out, uint64(a[1]), 10)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Int4RangeArrayToUint16Array2Slice struct {
	Val *[][2]uint16
}

func (v Int4RangeArrayToUint16Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
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

	*v.Val = ranges
	return nil
}

type Int4RangeArrayFromUint32Array2Slice struct {
	Val [][2]uint32
}

func (v Int4RangeArrayFromUint32Array2Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.Val {
		out = append(out, '"', '[')
		out = strconv.AppendUint(out, uint64(a[0]), 10)
		out = append(out, ',')
		out = strconv.AppendUint(out, uint64(a[1]), 10)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Int4RangeArrayToUint32Array2Slice struct {
	Val *[][2]uint32
}

func (v Int4RangeArrayToUint32Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
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

	*v.Val = ranges
	return nil
}

type Int4RangeArrayFromUint64Array2Slice struct {
	Val [][2]uint64
}

func (v Int4RangeArrayFromUint64Array2Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.Val {
		out = append(out, '"', '[')
		out = strconv.AppendUint(out, uint64(a[0]), 10)
		out = append(out, ',')
		out = strconv.AppendUint(out, uint64(a[1]), 10)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Int4RangeArrayToUint64Array2Slice struct {
	Val *[][2]uint64
}

func (v Int4RangeArrayToUint64Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
		return nil
	}

	elems := pgParseQuotedStringArray(data)
	ranges := make([][2]uint64, len(elems))

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

		ranges[i][0] = lo
		ranges[i][1] = hi
	}

	*v.Val = ranges
	return nil
}

type Int4RangeArrayFromFloat32Array2Slice struct {
	Val [][2]float32
}

func (v Int4RangeArrayFromFloat32Array2Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.Val {
		out = append(out, '"', '[')
		out = strconv.AppendInt(out, int64(a[0]), 10)
		out = append(out, ',')
		out = strconv.AppendInt(out, int64(a[1]), 10)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Int4RangeArrayToFloat32Array2Slice struct {
	Val *[][2]float32
}

func (v Int4RangeArrayToFloat32Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
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

	*v.Val = ranges
	return nil
}

type Int4RangeArrayFromFloat64Array2Slice struct {
	Val [][2]float64
}

func (v Int4RangeArrayFromFloat64Array2Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, a := range v.Val {
		out = append(out, '"', '[')
		out = strconv.AppendInt(out, int64(a[0]), 10)
		out = append(out, ',')
		out = strconv.AppendInt(out, int64(a[1]), 10)
		out = append(out, ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Int4RangeArrayToFloat64Array2Slice struct {
	Val *[][2]float64
}

func (v Int4RangeArrayToFloat64Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
		return nil
	}

	elems := pgParseQuotedStringArray(data)
	ranges := make([][2]float64, len(elems))

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

		ranges[i][0] = float64(lo)
		ranges[i][1] = float64(hi)
	}

	*v.Val = ranges
	return nil
}
