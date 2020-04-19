package convert

import (
	"database/sql/driver"
	"strconv"
)

type NumericArrayFromIntSlice struct {
	Val []int
}

func (v NumericArrayFromIntSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, i := range v.Val {
		out = strconv.AppendInt(out, int64(i), 10)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type NumericArrayToIntSlice struct {
	Val *[]int
}

func (v NumericArrayToIntSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	ints := make([]int, len(elems))
	for i := 0; i < len(elems); i++ {
		i64, err := strconv.ParseInt(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		ints[i] = int(i64)
	}

	*v.Val = ints
	return nil
}

type NumericArrayFromInt8Slice struct {
	Val []int8
}

func (v NumericArrayFromInt8Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, i8 := range v.Val {
		out = strconv.AppendInt(out, int64(i8), 10)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type NumericArrayToInt8Slice struct {
	Val *[]int8
}

func (v NumericArrayToInt8Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	int8s := make([]int8, len(elems))
	for i := 0; i < len(elems); i++ {
		i64, err := strconv.ParseInt(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		int8s[i] = int8(i64)
	}

	*v.Val = int8s
	return nil
}

type NumericArrayFromInt16Slice struct {
	Val []int16
}

func (v NumericArrayFromInt16Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, i16 := range v.Val {
		out = strconv.AppendInt(out, int64(i16), 10)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type NumericArrayToInt16Slice struct {
	Val *[]int16
}

func (v NumericArrayToInt16Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	int16s := make([]int16, len(elems))
	for i := 0; i < len(elems); i++ {
		i64, err := strconv.ParseInt(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		int16s[i] = int16(i64)
	}

	*v.Val = int16s
	return nil
}

type NumericArrayFromInt32Slice struct {
	Val []int32
}

func (v NumericArrayFromInt32Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, i32 := range v.Val {
		out = strconv.AppendInt(out, int64(i32), 10)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type NumericArrayToInt32Slice struct {
	Val *[]int32
}

func (v NumericArrayToInt32Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	int32s := make([]int32, len(elems))
	for i := 0; i < len(elems); i++ {
		i64, err := strconv.ParseInt(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		int32s[i] = int32(i64)
	}

	*v.Val = int32s
	return nil
}

type NumericArrayFromInt64Slice struct {
	Val []int64
}

func (v NumericArrayFromInt64Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, i64 := range v.Val {
		out = strconv.AppendInt(out, i64, 10)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type NumericArrayToInt64Slice struct {
	Val *[]int64
}

func (v NumericArrayToInt64Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	int64s := make([]int64, len(elems))
	for i := 0; i < len(elems); i++ {
		i64, err := strconv.ParseInt(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		int64s[i] = i64
	}

	*v.Val = int64s
	return nil
}

////////////////////////////////////////////////////////////////////////////////

type NumericArrayFromUintSlice struct {
	Val []uint
}

func (v NumericArrayFromUintSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, u := range v.Val {
		out = strconv.AppendUint(out, uint64(u), 10)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type NumericArrayToUintSlice struct {
	Val *[]uint
}

func (v NumericArrayToUintSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	uints := make([]uint, len(elems))
	for i := 0; i < len(elems); i++ {
		u64, err := strconv.ParseUint(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		uints[i] = uint(u64)
	}

	*v.Val = uints
	return nil
}

type NumericArrayFromUint8Slice struct {
	Val []uint8
}

func (v NumericArrayFromUint8Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, u := range v.Val {
		out = strconv.AppendUint(out, uint64(u), 10)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type NumericArrayToUint8Slice struct {
	Val *[]uint8
}

func (v NumericArrayToUint8Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	uint8s := make([]uint8, len(elems))
	for i := 0; i < len(elems); i++ {
		u64, err := strconv.ParseUint(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		uint8s[i] = uint8(u64)
	}

	*v.Val = uint8s
	return nil
}

type NumericArrayFromUint16Slice struct {
	Val []uint16
}

func (v NumericArrayFromUint16Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, u := range v.Val {
		out = strconv.AppendUint(out, uint64(u), 10)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type NumericArrayToUint16Slice struct {
	Val *[]uint16
}

func (v NumericArrayToUint16Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	uint16s := make([]uint16, len(elems))
	for i := 0; i < len(elems); i++ {
		u64, err := strconv.ParseUint(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		uint16s[i] = uint16(u64)
	}

	*v.Val = uint16s
	return nil
}

type NumericArrayFromUint32Slice struct {
	Val []uint32
}

func (v NumericArrayFromUint32Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, u := range v.Val {
		out = strconv.AppendUint(out, uint64(u), 10)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type NumericArrayToUint32Slice struct {
	Val *[]uint32
}

func (v NumericArrayToUint32Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	uint32s := make([]uint32, len(elems))
	for i := 0; i < len(elems); i++ {
		u64, err := strconv.ParseUint(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		uint32s[i] = uint32(u64)
	}

	*v.Val = uint32s
	return nil
}

type NumericArrayFromUint64Slice struct {
	Val []uint64
}

func (v NumericArrayFromUint64Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, u := range v.Val {
		out = strconv.AppendUint(out, u, 10)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type NumericArrayToUint64Slice struct {
	Val *[]uint64
}

func (v NumericArrayToUint64Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	uint64s := make([]uint64, len(elems))
	for i := 0; i < len(elems); i++ {
		u64, err := strconv.ParseUint(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		uint64s[i] = u64
	}

	*v.Val = uint64s
	return nil
}

type NumericArrayFromFloat32Slice struct {
	Val []float32
}

func (v NumericArrayFromFloat32Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, f := range v.Val {
		out = strconv.AppendInt(out, int64(f), 10)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type NumericArrayToFloat32Slice struct {
	Val *[]float32
}

func (v NumericArrayToFloat32Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	float32s := make([]float32, len(elems))
	for i := 0; i < len(elems); i++ {
		f64, err := strconv.ParseFloat(string(elems[i]), 64)
		if err != nil {
			return err
		}
		float32s[i] = float32(f64)
	}

	*v.Val = float32s
	return nil
}

type NumericArrayFromFloat64Slice struct {
	Val []float64
}

func (v NumericArrayFromFloat64Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, f := range v.Val {
		out = strconv.AppendInt(out, int64(f), 10)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type NumericArrayToFloat64Slice struct {
	Val *[]float64
}

func (v NumericArrayToFloat64Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	float64s := make([]float64, len(elems))
	for i := 0; i < len(elems); i++ {
		f64, err := strconv.ParseFloat(string(elems[i]), 64)
		if err != nil {
			return err
		}
		float64s[i] = f64
	}

	*v.Val = float64s
	return nil
}
