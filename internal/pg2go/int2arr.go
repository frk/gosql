package pg2go

import (
	"database/sql/driver"
	"strconv"
)

type Int2ArrayFromIntSlice struct {
	Val []int
}

func (v Int2ArrayFromIntSlice) Value() (driver.Value, error) {
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

type Int2ArrayToIntSlice struct {
	Val *[]int
}

func (v Int2ArrayToIntSlice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		ints[i] = int(i64)
	}

	*v.Val = ints
	return nil
}

type Int2ArrayFromInt8Slice struct {
	Val []int8
}

func (v Int2ArrayFromInt8Slice) Value() (driver.Value, error) {
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

type Int2ArrayToInt8Slice struct {
	Val *[]int8
}

func (v Int2ArrayToInt8Slice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 8)
		if err != nil {
			return err
		}
		int8s[i] = int8(i64)
	}

	*v.Val = int8s
	return nil
}

type Int2ArrayFromInt16Slice struct {
	Val []int16
}

func (v Int2ArrayFromInt16Slice) Value() (driver.Value, error) {
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

type Int2ArrayToInt16Slice struct {
	Val *[]int16
}

func (v Int2ArrayToInt16Slice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		int16s[i] = int16(i64)
	}

	*v.Val = int16s
	return nil
}

type Int2ArrayFromInt32Slice struct {
	Val []int32
}

func (v Int2ArrayFromInt32Slice) Value() (driver.Value, error) {
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

type Int2ArrayToInt32Slice struct {
	Val *[]int32
}

func (v Int2ArrayToInt32Slice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		int32s[i] = int32(i64)
	}

	*v.Val = int32s
	return nil
}

type Int2ArrayFromInt64Slice struct {
	Val []int64
}

func (v Int2ArrayFromInt64Slice) Value() (driver.Value, error) {
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

type Int2ArrayToInt64Slice struct {
	Val *[]int64
}

func (v Int2ArrayToInt64Slice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		int64s[i] = i64
	}

	*v.Val = int64s
	return nil
}

////////////////////////////////////////////////////////////////////////////////

type Int2ArrayFromUintSlice struct {
	Val []uint
}

func (v Int2ArrayFromUintSlice) Value() (driver.Value, error) {
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

type Int2ArrayToUintSlice struct {
	Val *[]uint
}

func (v Int2ArrayToUintSlice) Scan(src interface{}) error {
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
		u64, err := strconv.ParseUint(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		uints[i] = uint(u64)
	}

	*v.Val = uints
	return nil
}

type Int2ArrayFromUint8Slice struct {
	Val []uint8
}

func (v Int2ArrayFromUint8Slice) Value() (driver.Value, error) {
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

type Int2ArrayToUint8Slice struct {
	Val *[]uint8
}

func (v Int2ArrayToUint8Slice) Scan(src interface{}) error {
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
		u64, err := strconv.ParseUint(string(elems[i]), 10, 8)
		if err != nil {
			return err
		}
		uint8s[i] = uint8(u64)
	}

	*v.Val = uint8s
	return nil
}

type Int2ArrayFromUint16Slice struct {
	Val []uint16
}

func (v Int2ArrayFromUint16Slice) Value() (driver.Value, error) {
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

type Int2ArrayToUint16Slice struct {
	Val *[]uint16
}

func (v Int2ArrayToUint16Slice) Scan(src interface{}) error {
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
		u64, err := strconv.ParseUint(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		uint16s[i] = uint16(u64)
	}

	*v.Val = uint16s
	return nil
}

type Int2ArrayFromUint32Slice struct {
	Val []uint32
}

func (v Int2ArrayFromUint32Slice) Value() (driver.Value, error) {
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

type Int2ArrayToUint32Slice struct {
	Val *[]uint32
}

func (v Int2ArrayToUint32Slice) Scan(src interface{}) error {
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
		u64, err := strconv.ParseUint(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		uint32s[i] = uint32(u64)
	}

	*v.Val = uint32s
	return nil
}

type Int2ArrayFromUint64Slice struct {
	Val []uint64
}

func (v Int2ArrayFromUint64Slice) Value() (driver.Value, error) {
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

type Int2ArrayToUint64Slice struct {
	Val *[]uint64
}

func (v Int2ArrayToUint64Slice) Scan(src interface{}) error {
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
		u64, err := strconv.ParseUint(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		uint64s[i] = u64
	}

	*v.Val = uint64s
	return nil
}

type Int2ArrayFromFloat32Slice struct {
	Val []float32
}

func (v Int2ArrayFromFloat32Slice) Value() (driver.Value, error) {
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

type Int2ArrayToFloat32Slice struct {
	Val *[]float32
}

func (v Int2ArrayToFloat32Slice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		float32s[i] = float32(i64)
	}

	*v.Val = float32s
	return nil
}

type Int2ArrayFromFloat64Slice struct {
	Val []float64
}

func (v Int2ArrayFromFloat64Slice) Value() (driver.Value, error) {
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

type Int2ArrayToFloat64Slice struct {
	Val *[]float64
}

func (v Int2ArrayToFloat64Slice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		float64s[i] = float64(i64)
	}

	*v.Val = float64s
	return nil
}
