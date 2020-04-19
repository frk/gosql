package convert

import (
	"database/sql/driver"
	"strconv"
)

type NumRangeFromIntArray2 struct {
	Val [2]int
}

func (v NumRangeFromIntArray2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type NumRangeToIntArray2 struct {
	Val *[2]int
}

func (v NumRangeToIntArray2) Scan(src interface{}) error {
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

type NumRangeFromInt8Array2 struct {
	Val [2]int8
}

func (v NumRangeFromInt8Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type NumRangeToInt8Array2 struct {
	Val *[2]int8
}

func (v NumRangeToInt8Array2) Scan(src interface{}) error {
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

type NumRangeFromInt16Array2 struct {
	Val [2]int16
}

func (v NumRangeFromInt16Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type NumRangeToInt16Array2 struct {
	Val *[2]int16
}

func (v NumRangeToInt16Array2) Scan(src interface{}) error {
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

type NumRangeFromInt32Array2 struct {
	Val [2]int32
}

func (v NumRangeFromInt32Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type NumRangeToInt32Array2 struct {
	Val *[2]int32
}

func (v NumRangeToInt32Array2) Scan(src interface{}) error {
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

type NumRangeFromInt64Array2 struct {
	Val [2]int64
}

func (v NumRangeFromInt64Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, v.Val[0], 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, v.Val[1], 10)
	out = append(out, ')')
	return out, nil
}

type NumRangeToInt64Array2 struct {
	Val *[2]int64
}

func (v NumRangeToInt64Array2) Scan(src interface{}) error {
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

type NumRangeFromUintArray2 struct {
	Val [2]uint
}

func (v NumRangeFromUintArray2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, uint64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, uint64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type NumRangeToUintArray2 struct {
	Val *[2]uint
}

func (v NumRangeToUintArray2) Scan(src interface{}) error {
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

type NumRangeFromUint8Array2 struct {
	Val [2]uint8
}

func (v NumRangeFromUint8Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, uint64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, uint64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type NumRangeToUint8Array2 struct {
	Val *[2]uint8
}

func (v NumRangeToUint8Array2) Scan(src interface{}) error {
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

type NumRangeFromUint16Array2 struct {
	Val [2]uint16
}

func (v NumRangeFromUint16Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, uint64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, uint64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type NumRangeToUint16Array2 struct {
	Val *[2]uint16
}

func (v NumRangeToUint16Array2) Scan(src interface{}) error {
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

type NumRangeFromUint32Array2 struct {
	Val [2]uint32
}

func (v NumRangeFromUint32Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, uint64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, uint64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type NumRangeToUint32Array2 struct {
	Val *[2]uint32
}

func (v NumRangeToUint32Array2) Scan(src interface{}) error {
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

type NumRangeFromUint64Array2 struct {
	Val [2]uint64
}

func (v NumRangeFromUint64Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendUint(out, v.Val[0], 10)
	out = append(out, ',')
	out = strconv.AppendUint(out, v.Val[1], 10)
	out = append(out, ')')
	return out, nil
}

type NumRangeToUint64Array2 struct {
	Val *[2]uint64
}

func (v NumRangeToUint64Array2) Scan(src interface{}) error {
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

type NumRangeFromFloat32Array2 struct {
	Val [2]float32
}

func (v NumRangeFromFloat32Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type NumRangeToFloat32Array2 struct {
	Val *[2]float32
}

func (v NumRangeToFloat32Array2) Scan(src interface{}) error {
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

type NumRangeFromFloat64Array2 struct {
	Val [2]float64
}

func (v NumRangeFromFloat64Array2) Value() (driver.Value, error) {
	out := []byte{'['}
	out = strconv.AppendInt(out, int64(v.Val[0]), 10)
	out = append(out, ',')
	out = strconv.AppendInt(out, int64(v.Val[1]), 10)
	out = append(out, ')')
	return out, nil
}

type NumRangeToFloat64Array2 struct {
	Val *[2]float64
}

func (v NumRangeToFloat64Array2) Scan(src interface{}) error {
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
