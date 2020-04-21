package pg2go

import (
	"database/sql/driver"
	"strconv"
)

type Int4RangeFromIntArray2 struct {
	Val [2]int
}

func (v Int4RangeFromIntArray2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type Int4RangeToIntArray2 struct {
	Val *[2]int
}

func (v Int4RangeToIntArray2) Scan(src interface{}) error {
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

	v.Val[0] = int(lo)
	v.Val[1] = int(hi)
	return nil
}

type Int4RangeFromInt8Array2 struct {
	Val [2]int8
}

func (v Int4RangeFromInt8Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type Int4RangeToInt8Array2 struct {
	Val *[2]int8
}

func (v Int4RangeToInt8Array2) Scan(src interface{}) error {
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

type Int4RangeFromInt16Array2 struct {
	Val [2]int16
}

func (v Int4RangeFromInt16Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type Int4RangeToInt16Array2 struct {
	Val *[2]int16
}

func (v Int4RangeToInt16Array2) Scan(src interface{}) error {
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

type Int4RangeFromInt32Array2 struct {
	Val [2]int32
}

func (v Int4RangeFromInt32Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type Int4RangeToInt32Array2 struct {
	Val *[2]int32
}

func (v Int4RangeToInt32Array2) Scan(src interface{}) error {
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

type Int4RangeFromInt64Array2 struct {
	Val [2]int64
}

func (v Int4RangeFromInt64Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, v.Val[0], 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, v.Val[1], 10)
	out = append(out, ')')
	return out, nil
}

type Int4RangeToInt64Array2 struct {
	Val *[2]int64
}

func (v Int4RangeToInt64Array2) Scan(src interface{}) error {
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

	v.Val[0] = lo
	v.Val[1] = hi
	return nil
}

type Int4RangeFromUintArray2 struct {
	Val [2]uint
}

func (v Int4RangeFromUintArray2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, uint64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, uint64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type Int4RangeToUintArray2 struct {
	Val *[2]uint
}

func (v Int4RangeToUintArray2) Scan(src interface{}) error {
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

	v.Val[0] = uint(lo)
	v.Val[1] = uint(hi)
	return nil
}

type Int4RangeFromUint8Array2 struct {
	Val [2]uint8
}

func (v Int4RangeFromUint8Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, uint64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, uint64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type Int4RangeToUint8Array2 struct {
	Val *[2]uint8
}

func (v Int4RangeToUint8Array2) Scan(src interface{}) error {
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

type Int4RangeFromUint16Array2 struct {
	Val [2]uint16
}

func (v Int4RangeFromUint16Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, uint64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, uint64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type Int4RangeToUint16Array2 struct {
	Val *[2]uint16
}

func (v Int4RangeToUint16Array2) Scan(src interface{}) error {
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

type Int4RangeFromUint32Array2 struct {
	Val [2]uint32
}

func (v Int4RangeFromUint32Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, uint64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, uint64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type Int4RangeToUint32Array2 struct {
	Val *[2]uint32
}

func (v Int4RangeToUint32Array2) Scan(src interface{}) error {
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

type Int4RangeFromUint64Array2 struct {
	Val [2]uint64
}

func (v Int4RangeFromUint64Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, v.Val[0], 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, v.Val[1], 10)
	out = append(out, ')')
	return out, nil
}

type Int4RangeToUint64Array2 struct {
	Val *[2]uint64
}

func (v Int4RangeToUint64Array2) Scan(src interface{}) error {
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

	v.Val[0] = lo
	v.Val[1] = hi
	return nil
}

type Int4RangeFromFloat32Array2 struct {
	Val [2]float32
}

func (v Int4RangeFromFloat32Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type Int4RangeToFloat32Array2 struct {
	Val *[2]float32
}

func (v Int4RangeToFloat32Array2) Scan(src interface{}) error {
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

	v.Val[0] = float32(lo)
	v.Val[1] = float32(hi)
	return nil
}

type Int4RangeFromFloat64Array2 struct {
	Val [2]float64
}

func (v Int4RangeFromFloat64Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type Int4RangeToFloat64Array2 struct {
	Val *[2]float64
}

func (v Int4RangeToFloat64Array2) Scan(src interface{}) error {
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

	v.Val[0] = float64(lo)
	v.Val[1] = float64(hi)
	return nil
}
