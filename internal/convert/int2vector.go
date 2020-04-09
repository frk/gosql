package convert

import (
	"database/sql/driver"
	"strconv"
)

type Int2VectorFromIntSlice struct {
	Val []int
}

func (v Int2VectorFromIntSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{}, nil
	}

	out := []byte{}

	for _, i := range v.Val {
		out = strconv.AppendInt(out, int64(i), 10)
		out = append(out, ' ')
	}

	out = out[:len(out)-1] // drop last " "
	return out, nil
}

type Int2VectorToIntSlice struct {
	Val *[]int
}

func (v Int2VectorToIntSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVector(arr)
	ints := make([]int, len(elems))
	for i := 0; i < len(elems); i++ {
		i64, err := strconv.ParseInt(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		ints[i] = int(i64)
	}

	*v.Val = ints
	return nil
}

type Int2VectorFromInt8Slice struct {
	Val []int8
}

func (v Int2VectorFromInt8Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{}, nil
	}

	out := []byte{}

	for _, i8 := range v.Val {
		out = strconv.AppendInt(out, int64(i8), 10)
		out = append(out, ' ')
	}

	out = out[:len(out)-1] // drop last " "
	return out, nil
}

type Int2VectorToInt8Slice struct {
	Val *[]int8
}

func (v Int2VectorToInt8Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVector(arr)
	int8s := make([]int8, len(elems))
	for i := 0; i < len(elems); i++ {
		i64, err := strconv.ParseInt(string(elems[i]), 10, 8)
		if err != nil {
			return err
		}
		int8s[i] = int8(i64)
	}

	*v.Val = int8s
	return nil
}

type Int2VectorFromInt16Slice struct {
	Val []int16
}

func (v Int2VectorFromInt16Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{}, nil
	}

	out := []byte{}

	for _, i16 := range v.Val {
		out = strconv.AppendInt(out, int64(i16), 10)
		out = append(out, ' ')
	}

	out = out[:len(out)-1] // drop last " "
	return out, nil
}

type Int2VectorToInt16Slice struct {
	Val *[]int16
}

func (v Int2VectorToInt16Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVector(arr)
	int16s := make([]int16, len(elems))
	for i := 0; i < len(elems); i++ {
		i64, err := strconv.ParseInt(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		int16s[i] = int16(i64)
	}

	*v.Val = int16s
	return nil
}

type Int2VectorFromInt32Slice struct {
	Val []int32
}

func (v Int2VectorFromInt32Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{}, nil
	}

	out := []byte{}

	for _, i32 := range v.Val {
		out = strconv.AppendInt(out, int64(i32), 10)
		out = append(out, ' ')
	}

	out = out[:len(out)-1] // drop last " "
	return out, nil
}

type Int2VectorToInt32Slice struct {
	Val *[]int32
}

func (v Int2VectorToInt32Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVector(arr)
	int32s := make([]int32, len(elems))
	for i := 0; i < len(elems); i++ {
		i64, err := strconv.ParseInt(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		int32s[i] = int32(i64)
	}

	*v.Val = int32s
	return nil
}

type Int2VectorFromInt64Slice struct {
	Val []int64
}

func (v Int2VectorFromInt64Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{}, nil
	}

	out := []byte{}

	for _, i64 := range v.Val {
		out = strconv.AppendInt(out, i64, 10)
		out = append(out, ' ')
	}

	out = out[:len(out)-1] // drop last " "
	return out, nil
}

type Int2VectorToInt64Slice struct {
	Val *[]int64
}

func (v Int2VectorToInt64Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVector(arr)
	int64s := make([]int64, len(elems))
	for i := 0; i < len(elems); i++ {
		i64, err := strconv.ParseInt(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		int64s[i] = i64
	}

	*v.Val = int64s
	return nil
}

type Int2VectorFromUintSlice struct {
	Val []uint
}

func (v Int2VectorFromUintSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{}, nil
	}

	out := []byte{}

	for _, u := range v.Val {
		out = strconv.AppendUint(out, uint64(u), 10)
		out = append(out, ' ')
	}

	out = out[:len(out)-1] // drop last " "
	return out, nil
}

type Int2VectorToUintSlice struct {
	Val *[]uint
}

func (v Int2VectorToUintSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVector(arr)
	uints := make([]uint, len(elems))
	for i := 0; i < len(elems); i++ {
		u64, err := strconv.ParseUint(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		uints[i] = uint(u64)
	}

	*v.Val = uints
	return nil
}

type Int2VectorFromUint8Slice struct {
	Val []uint8
}

func (v Int2VectorFromUint8Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{}, nil
	}

	out := []byte{}

	for _, u8 := range v.Val {
		out = strconv.AppendUint(out, uint64(u8), 10)
		out = append(out, ' ')
	}

	out = out[:len(out)-1] // drop last " "
	return out, nil
}

type Int2VectorToUint8Slice struct {
	Val *[]uint8
}

func (v Int2VectorToUint8Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVector(arr)
	uint8s := make([]uint8, len(elems))
	for i := 0; i < len(elems); i++ {
		u64, err := strconv.ParseUint(string(elems[i]), 10, 8)
		if err != nil {
			return err
		}
		uint8s[i] = uint8(u64)
	}

	*v.Val = uint8s
	return nil
}

type Int2VectorFromUint16Slice struct {
	Val []uint16
}

func (v Int2VectorFromUint16Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{}, nil
	}

	out := []byte{}

	for _, u16 := range v.Val {
		out = strconv.AppendUint(out, uint64(u16), 10)
		out = append(out, ' ')
	}

	out = out[:len(out)-1] // drop last " "
	return out, nil
}

type Int2VectorToUint16Slice struct {
	Val *[]uint16
}

func (v Int2VectorToUint16Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVector(arr)
	uint16s := make([]uint16, len(elems))
	for i := 0; i < len(elems); i++ {
		u64, err := strconv.ParseUint(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		uint16s[i] = uint16(u64)
	}

	*v.Val = uint16s
	return nil
}

type Int2VectorFromUint32Slice struct {
	Val []uint32
}

func (v Int2VectorFromUint32Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{}, nil
	}

	out := []byte{}

	for _, u32 := range v.Val {
		out = strconv.AppendUint(out, uint64(u32), 10)
		out = append(out, ' ')
	}

	out = out[:len(out)-1] // drop last " "
	return out, nil
}

type Int2VectorToUint32Slice struct {
	Val *[]uint32
}

func (v Int2VectorToUint32Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVector(arr)
	uint32s := make([]uint32, len(elems))
	for i := 0; i < len(elems); i++ {
		u64, err := strconv.ParseUint(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		uint32s[i] = uint32(u64)
	}

	*v.Val = uint32s
	return nil
}

type Int2VectorFromUint64Slice struct {
	Val []uint64
}

func (v Int2VectorFromUint64Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{}, nil
	}

	out := []byte{}

	for _, u64 := range v.Val {
		out = strconv.AppendUint(out, u64, 10)
		out = append(out, ' ')
	}

	out = out[:len(out)-1] // drop last " "
	return out, nil
}

type Int2VectorToUint64Slice struct {
	Val *[]uint64
}

func (v Int2VectorToUint64Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVector(arr)
	uint64s := make([]uint64, len(elems))
	for i := 0; i < len(elems); i++ {
		u64, err := strconv.ParseUint(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		uint64s[i] = u64
	}

	*v.Val = uint64s
	return nil
}

type Int2VectorFromFloat32Slice struct {
	Val []float32
}

func (v Int2VectorFromFloat32Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{}, nil
	}

	out := []byte{}

	for _, f32 := range v.Val {
		out = strconv.AppendInt(out, int64(f32), 10)
		out = append(out, ' ')
	}

	out = out[:len(out)-1] // drop last " "
	return out, nil
}

type Int2VectorToFloat32Slice struct {
	Val *[]float32
}

func (v Int2VectorToFloat32Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVector(arr)
	float32s := make([]float32, len(elems))
	for i := 0; i < len(elems); i++ {
		i64, err := strconv.ParseInt(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		float32s[i] = float32(i64)
	}

	*v.Val = float32s
	return nil
}

type Int2VectorFromFloat64Slice struct {
	Val []float64
}

func (v Int2VectorFromFloat64Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{}, nil
	}

	out := []byte{}

	for _, f64 := range v.Val {
		out = strconv.AppendInt(out, int64(f64), 10)
		out = append(out, ' ')
	}

	out = out[:len(out)-1] // drop last " "
	return out, nil
}

type Int2VectorToFloat64Slice struct {
	Val *[]float64
}

func (v Int2VectorToFloat64Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVector(arr)
	float64s := make([]float64, len(elems))
	for i := 0; i < len(elems); i++ {
		i64, err := strconv.ParseInt(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		float64s[i] = float64(i64)
	}

	*v.Val = float64s
	return nil
}
