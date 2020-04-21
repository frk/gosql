package pg2go

import (
	"database/sql/driver"
	"strconv"
)

type NumRangeArrayFromIntArray2Slice struct {
	Val [][2]int
}

func (v NumRangeArrayFromIntArray2Slice) Value() (driver.Value, error) {
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

type NumRangeArrayToIntArray2Slice struct {
	Val *[][2]int
}

func (v NumRangeArrayToIntArray2Slice) Scan(src interface{}) error {
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

type NumRangeArrayFromInt8Array2Slice struct {
	Val [][2]int8
}

func (v NumRangeArrayFromInt8Array2Slice) Value() (driver.Value, error) {
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

type NumRangeArrayToInt8Array2Slice struct {
	Val *[][2]int8
}

func (v NumRangeArrayToInt8Array2Slice) Scan(src interface{}) error {
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

type NumRangeArrayFromInt16Array2Slice struct {
	Val [][2]int16
}

func (v NumRangeArrayFromInt16Array2Slice) Value() (driver.Value, error) {
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

type NumRangeArrayToInt16Array2Slice struct {
	Val *[][2]int16
}

func (v NumRangeArrayToInt16Array2Slice) Scan(src interface{}) error {
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

type NumRangeArrayFromInt32Array2Slice struct {
	Val [][2]int32
}

func (v NumRangeArrayFromInt32Array2Slice) Value() (driver.Value, error) {
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

type NumRangeArrayToInt32Array2Slice struct {
	Val *[][2]int32
}

func (v NumRangeArrayToInt32Array2Slice) Scan(src interface{}) error {
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

type NumRangeArrayFromInt64Array2Slice struct {
	Val [][2]int64
}

func (v NumRangeArrayFromInt64Array2Slice) Value() (driver.Value, error) {
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

type NumRangeArrayToInt64Array2Slice struct {
	Val *[][2]int64
}

func (v NumRangeArrayToInt64Array2Slice) Scan(src interface{}) error {
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

type NumRangeArrayFromUintArray2Slice struct {
	Val [][2]uint
}

func (v NumRangeArrayFromUintArray2Slice) Value() (driver.Value, error) {
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

type NumRangeArrayToUintArray2Slice struct {
	Val *[][2]uint
}

func (v NumRangeArrayToUintArray2Slice) Scan(src interface{}) error {
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

type NumRangeArrayFromUint8Array2Slice struct {
	Val [][2]uint8
}

func (v NumRangeArrayFromUint8Array2Slice) Value() (driver.Value, error) {
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

type NumRangeArrayToUint8Array2Slice struct {
	Val *[][2]uint8
}

func (v NumRangeArrayToUint8Array2Slice) Scan(src interface{}) error {
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

type NumRangeArrayFromUint16Array2Slice struct {
	Val [][2]uint16
}

func (v NumRangeArrayFromUint16Array2Slice) Value() (driver.Value, error) {
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

type NumRangeArrayToUint16Array2Slice struct {
	Val *[][2]uint16
}

func (v NumRangeArrayToUint16Array2Slice) Scan(src interface{}) error {
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

type NumRangeArrayFromUint32Array2Slice struct {
	Val [][2]uint32
}

func (v NumRangeArrayFromUint32Array2Slice) Value() (driver.Value, error) {
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

type NumRangeArrayToUint32Array2Slice struct {
	Val *[][2]uint32
}

func (v NumRangeArrayToUint32Array2Slice) Scan(src interface{}) error {
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

type NumRangeArrayFromUint64Array2Slice struct {
	Val [][2]uint64
}

func (v NumRangeArrayFromUint64Array2Slice) Value() (driver.Value, error) {
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

type NumRangeArrayToUint64Array2Slice struct {
	Val *[][2]uint64
}

func (v NumRangeArrayToUint64Array2Slice) Scan(src interface{}) error {
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

type NumRangeArrayFromFloat32Array2Slice struct {
	Val [][2]float32
}

func (v NumRangeArrayFromFloat32Array2Slice) Value() (driver.Value, error) {
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

type NumRangeArrayToFloat32Array2Slice struct {
	Val *[][2]float32
}

func (v NumRangeArrayToFloat32Array2Slice) Scan(src interface{}) error {
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

type NumRangeArrayFromFloat64Array2Slice struct {
	Val [][2]float64
}

func (v NumRangeArrayFromFloat64Array2Slice) Value() (driver.Value, error) {
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

type NumRangeArrayToFloat64Array2Slice struct {
	Val *[][2]float64
}

func (v NumRangeArrayToFloat64Array2Slice) Scan(src interface{}) error {
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
