package convert

import (
	"database/sql/driver"
	"strconv"
)

type Int8RangeArrayFromIntArray2Slice struct {
	Val [][2]int
}

func (v Int8RangeArrayFromIntArray2Slice) Value() (driver.Value, error) {
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

type Int8RangeArrayToIntArray2Slice struct {
	Val *[][2]int
}

func (v Int8RangeArrayToIntArray2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
		return nil
	}

	elems := pgParseStringArray(data)
	ranges := make([][2]int, len(elems))

	for i, elem := range elems {
		elem = elem[1 : len(elem)-1] // drop surrounding double quotes

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

	*v.Val = ranges
	return nil
}

type Int8RangeArrayFromInt8Array2Slice struct {
	Val [][2]int8
}

func (v Int8RangeArrayFromInt8Array2Slice) Value() (driver.Value, error) {
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

type Int8RangeArrayToInt8Array2Slice struct {
	Val *[][2]int8
}

func (v Int8RangeArrayToInt8Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
		return nil
	}

	elems := pgParseStringArray(data)
	ranges := make([][2]int8, len(elems))

	for i, elem := range elems {
		elem = elem[1 : len(elem)-1] // drop surrounding double quotes

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

type Int8RangeArrayFromInt16Array2Slice struct {
	Val [][2]int16
}

func (v Int8RangeArrayFromInt16Array2Slice) Value() (driver.Value, error) {
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

type Int8RangeArrayToInt16Array2Slice struct {
	Val *[][2]int16
}

func (v Int8RangeArrayToInt16Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
		return nil
	}

	elems := pgParseStringArray(data)
	ranges := make([][2]int16, len(elems))

	for i, elem := range elems {
		elem = elem[1 : len(elem)-1] // drop surrounding double quotes

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

type Int8RangeArrayFromInt32Array2Slice struct {
	Val [][2]int32
}

func (v Int8RangeArrayFromInt32Array2Slice) Value() (driver.Value, error) {
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

type Int8RangeArrayToInt32Array2Slice struct {
	Val *[][2]int32
}

func (v Int8RangeArrayToInt32Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
		return nil
	}

	elems := pgParseStringArray(data)
	ranges := make([][2]int32, len(elems))

	for i, elem := range elems {
		elem = elem[1 : len(elem)-1] // drop surrounding double quotes

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

type Int8RangeArrayFromInt64Array2Slice struct {
	Val [][2]int64
}

func (v Int8RangeArrayFromInt64Array2Slice) Value() (driver.Value, error) {
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

type Int8RangeArrayToInt64Array2Slice struct {
	Val *[][2]int64
}

func (v Int8RangeArrayToInt64Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
		return nil
	}

	elems := pgParseStringArray(data)
	ranges := make([][2]int64, len(elems))

	for i, elem := range elems {
		elem = elem[1 : len(elem)-1] // drop surrounding double quotes

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

	*v.Val = ranges
	return nil
}

type Int8RangeArrayFromUintArray2Slice struct {
	Val [][2]uint
}

func (v Int8RangeArrayFromUintArray2Slice) Value() (driver.Value, error) {
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

type Int8RangeArrayToUintArray2Slice struct {
	Val *[][2]uint
}

func (v Int8RangeArrayToUintArray2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
		return nil
	}

	elems := pgParseStringArray(data)
	ranges := make([][2]uint, len(elems))

	for i, elem := range elems {
		elem = elem[1 : len(elem)-1] // drop surrounding double quotes

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

	*v.Val = ranges
	return nil
}

type Int8RangeArrayFromUint8Array2Slice struct {
	Val [][2]uint8
}

func (v Int8RangeArrayFromUint8Array2Slice) Value() (driver.Value, error) {
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

type Int8RangeArrayToUint8Array2Slice struct {
	Val *[][2]uint8
}

func (v Int8RangeArrayToUint8Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
		return nil
	}

	elems := pgParseStringArray(data)
	ranges := make([][2]uint8, len(elems))

	for i, elem := range elems {
		elem = elem[1 : len(elem)-1] // drop surrounding double quotes

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

type Int8RangeArrayFromUint16Array2Slice struct {
	Val [][2]uint16
}

func (v Int8RangeArrayFromUint16Array2Slice) Value() (driver.Value, error) {
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

type Int8RangeArrayToUint16Array2Slice struct {
	Val *[][2]uint16
}

func (v Int8RangeArrayToUint16Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
		return nil
	}

	elems := pgParseStringArray(data)
	ranges := make([][2]uint16, len(elems))

	for i, elem := range elems {
		elem = elem[1 : len(elem)-1] // drop surrounding double quotes

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

type Int8RangeArrayFromUint32Array2Slice struct {
	Val [][2]uint32
}

func (v Int8RangeArrayFromUint32Array2Slice) Value() (driver.Value, error) {
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

type Int8RangeArrayToUint32Array2Slice struct {
	Val *[][2]uint32
}

func (v Int8RangeArrayToUint32Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
		return nil
	}

	elems := pgParseStringArray(data)
	ranges := make([][2]uint32, len(elems))

	for i, elem := range elems {
		elem = elem[1 : len(elem)-1] // drop surrounding double quotes

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

type Int8RangeArrayFromUint64Array2Slice struct {
	Val [][2]uint64
}

func (v Int8RangeArrayFromUint64Array2Slice) Value() (driver.Value, error) {
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

type Int8RangeArrayToUint64Array2Slice struct {
	Val *[][2]uint64
}

func (v Int8RangeArrayToUint64Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
		return nil
	}

	elems := pgParseStringArray(data)
	ranges := make([][2]uint64, len(elems))

	for i, elem := range elems {
		elem = elem[1 : len(elem)-1] // drop surrounding double quotes

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

	*v.Val = ranges
	return nil
}

type Int8RangeArrayFromFloat32Array2Slice struct {
	Val [][2]float32
}

func (v Int8RangeArrayFromFloat32Array2Slice) Value() (driver.Value, error) {
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

type Int8RangeArrayToFloat32Array2Slice struct {
	Val *[][2]float32
}

func (v Int8RangeArrayToFloat32Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
		return nil
	}

	elems := pgParseStringArray(data)
	ranges := make([][2]float32, len(elems))

	for i, elem := range elems {
		elem = elem[1 : len(elem)-1] // drop surrounding double quotes

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

type Int8RangeArrayFromFloat64Array2Slice struct {
	Val [][2]float64
}

func (v Int8RangeArrayFromFloat64Array2Slice) Value() (driver.Value, error) {
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

type Int8RangeArrayToFloat64Array2Slice struct {
	Val *[][2]float64
}

func (v Int8RangeArrayToFloat64Array2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		v.Val = nil
		return nil
	}

	elems := pgParseStringArray(data)
	ranges := make([][2]float64, len(elems))

	for i, elem := range elems {
		elem = elem[1 : len(elem)-1] // drop surrounding double quotes

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

	*v.Val = ranges
	return nil
}
