package convert

import (
	"database/sql/driver"
	"strconv"
)

type Int8RangeFromIntArray2 struct {
	Val [2]int
}

func (v Int8RangeFromIntArray2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type Int8RangeToIntArray2 struct {
	Val *[2]int
}

func (v Int8RangeToIntArray2) Scan(src interface{}) error {
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

	v.Val[0] = int(lo)
	v.Val[1] = int(hi)
	return nil
}

type Int8RangeFromInt8Array2 struct {
	Val [2]int8
}

func (v Int8RangeFromInt8Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type Int8RangeToInt8Array2 struct {
	Val *[2]int8
}

func (v Int8RangeToInt8Array2) Scan(src interface{}) error {
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

	v.Val[0] = int8(lo)
	v.Val[1] = int8(hi)
	return nil
}

type Int8RangeFromInt16Array2 struct {
	Val [2]int16
}

func (v Int8RangeFromInt16Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type Int8RangeToInt16Array2 struct {
	Val *[2]int16
}

func (v Int8RangeToInt16Array2) Scan(src interface{}) error {
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

	v.Val[0] = int16(lo)
	v.Val[1] = int16(hi)
	return nil
}

type Int8RangeFromInt32Array2 struct {
	Val [2]int32
}

func (v Int8RangeFromInt32Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type Int8RangeToInt32Array2 struct {
	Val *[2]int32
}

func (v Int8RangeToInt32Array2) Scan(src interface{}) error {
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

	v.Val[0] = int32(lo)
	v.Val[1] = int32(hi)
	return nil
}

type Int8RangeFromInt64Array2 struct {
	Val [2]int64
}

func (v Int8RangeFromInt64Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, v.Val[0], 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, v.Val[1], 10)
	out = append(out, ')')
	return out, nil
}

type Int8RangeToInt64Array2 struct {
	Val *[2]int64
}

func (v Int8RangeToInt64Array2) Scan(src interface{}) error {
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

	v.Val[0] = lo
	v.Val[1] = hi
	return nil
}

type Int8RangeFromUintArray2 struct {
	Val [2]uint
}

func (v Int8RangeFromUintArray2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, uint64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, uint64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type Int8RangeToUintArray2 struct {
	Val *[2]uint
}

func (v Int8RangeToUintArray2) Scan(src interface{}) error {
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

	v.Val[0] = uint(lo)
	v.Val[1] = uint(hi)
	return nil
}

type Int8RangeFromUint8Array2 struct {
	Val [2]uint8
}

func (v Int8RangeFromUint8Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, uint64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, uint64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type Int8RangeToUint8Array2 struct {
	Val *[2]uint8
}

func (v Int8RangeToUint8Array2) Scan(src interface{}) error {
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

	v.Val[0] = uint8(lo)
	v.Val[1] = uint8(hi)
	return nil
}

type Int8RangeFromUint16Array2 struct {
	Val [2]uint16
}

func (v Int8RangeFromUint16Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, uint64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, uint64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type Int8RangeToUint16Array2 struct {
	Val *[2]uint16
}

func (v Int8RangeToUint16Array2) Scan(src interface{}) error {
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

	v.Val[0] = uint16(lo)
	v.Val[1] = uint16(hi)
	return nil
}

type Int8RangeFromUint32Array2 struct {
	Val [2]uint32
}

func (v Int8RangeFromUint32Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, uint64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, uint64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type Int8RangeToUint32Array2 struct {
	Val *[2]uint32
}

func (v Int8RangeToUint32Array2) Scan(src interface{}) error {
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

	v.Val[0] = uint32(lo)
	v.Val[1] = uint32(hi)
	return nil
}

type Int8RangeFromUint64Array2 struct {
	Val [2]uint64
}

func (v Int8RangeFromUint64Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, v.Val[0], 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, v.Val[1], 10)
	out = append(out, ')')
	return out, nil
}

type Int8RangeToUint64Array2 struct {
	Val *[2]uint64
}

func (v Int8RangeToUint64Array2) Scan(src interface{}) error {
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

	v.Val[0] = lo
	v.Val[1] = hi
	return nil
}

type Int8RangeFromFloat32Array2 struct {
	Val [2]float32
}

func (v Int8RangeFromFloat32Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type Int8RangeToFloat32Array2 struct {
	Val *[2]float32
}

func (v Int8RangeToFloat32Array2) Scan(src interface{}) error {
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

	v.Val[0] = float32(lo)
	v.Val[1] = float32(hi)
	return nil
}

type Int8RangeFromFloat64Array2 struct {
	Val [2]float64
}

func (v Int8RangeFromFloat64Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type Int8RangeToFloat64Array2 struct {
	Val *[2]float64
}

func (v Int8RangeToFloat64Array2) Scan(src interface{}) error {
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

	v.Val[0] = float64(lo)
	v.Val[1] = float64(hi)
	return nil
}
