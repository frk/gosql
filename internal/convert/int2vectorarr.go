package convert

import (
	"database/sql/driver"
	"strconv"
)

type Int2VectorArrayFromIntSliceSlice struct {
	Val [][]int
}

func (v Int2VectorArrayFromIntSliceSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, ints := range v.Val {
		out = append(out, '"')
		for _, i := range ints {
			out = strconv.AppendInt(out, int64(i), 10)
			out = append(out, ' ')
		}
		out[len(out)-1] = '"' // replace last " " with `"`
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Int2VectorArrayToIntSliceSlice struct {
	Val *[][]int
}

func (v Int2VectorArrayToIntSliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVectorArray(arr)
	intss := make([][]int, len(elems))
	for i := 0; i < len(elems); i++ {
		ints := make([]int, len(elems[i]))
		for j := 0; j < len(elems[i]); j++ {
			i64, err := strconv.ParseInt(string(elems[i][j]), 10, 16)
			if err != nil {
				return err
			}
			ints[j] = int(i64)
		}
		intss[i] = ints
	}

	*v.Val = intss
	return nil
}

type Int2VectorArrayFromInt8SliceSlice struct {
	Val [][]int8
}

func (v Int2VectorArrayFromInt8SliceSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, int8s := range v.Val {
		out = append(out, '"')
		for _, i8 := range int8s {
			out = strconv.AppendInt(out, int64(i8), 10)
			out = append(out, ' ')
		}
		out[len(out)-1] = '"' // replace last " " with `"`
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Int2VectorArrayToInt8SliceSlice struct {
	Val *[][]int8
}

func (v Int2VectorArrayToInt8SliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVectorArray(arr)
	int8ss := make([][]int8, len(elems))
	for i := 0; i < len(elems); i++ {
		int8s := make([]int8, len(elems[i]))
		for j := 0; j < len(elems[i]); j++ {
			i64, err := strconv.ParseInt(string(elems[i][j]), 10, 8)
			if err != nil {
				return err
			}
			int8s[j] = int8(i64)
		}
		int8ss[i] = int8s
	}

	*v.Val = int8ss
	return nil
}

type Int2VectorArrayFromInt16SliceSlice struct {
	Val [][]int16
}

func (v Int2VectorArrayFromInt16SliceSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, int16s := range v.Val {
		out = append(out, '"')
		for _, i16 := range int16s {
			out = strconv.AppendInt(out, int64(i16), 10)
			out = append(out, ' ')
		}
		out[len(out)-1] = '"' // replace last " " with `"`
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Int2VectorArrayToInt16SliceSlice struct {
	Val *[][]int16
}

func (v Int2VectorArrayToInt16SliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVectorArray(arr)
	int16ss := make([][]int16, len(elems))
	for i := 0; i < len(elems); i++ {
		int16s := make([]int16, len(elems[i]))
		for j := 0; j < len(elems[i]); j++ {
			i64, err := strconv.ParseInt(string(elems[i][j]), 10, 16)
			if err != nil {
				return err
			}
			int16s[j] = int16(i64)
		}
		int16ss[i] = int16s
	}

	*v.Val = int16ss
	return nil
}

type Int2VectorArrayFromInt32SliceSlice struct {
	Val [][]int32
}

func (v Int2VectorArrayFromInt32SliceSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, int32s := range v.Val {
		out = append(out, '"')
		for _, i32 := range int32s {
			out = strconv.AppendInt(out, int64(i32), 10)
			out = append(out, ' ')
		}
		out[len(out)-1] = '"' // replace last " " with `"`
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Int2VectorArrayToInt32SliceSlice struct {
	Val *[][]int32
}

func (v Int2VectorArrayToInt32SliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVectorArray(arr)
	int32ss := make([][]int32, len(elems))
	for i := 0; i < len(elems); i++ {
		int32s := make([]int32, len(elems[i]))
		for j := 0; j < len(elems[i]); j++ {
			i64, err := strconv.ParseInt(string(elems[i][j]), 10, 16)
			if err != nil {
				return err
			}
			int32s[j] = int32(i64)
		}
		int32ss[i] = int32s
	}

	*v.Val = int32ss
	return nil
}

type Int2VectorArrayFromInt64SliceSlice struct {
	Val [][]int64
}

func (v Int2VectorArrayFromInt64SliceSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, int64s := range v.Val {
		out = append(out, '"')
		for _, i64 := range int64s {
			out = strconv.AppendInt(out, i64, 10)
			out = append(out, ' ')
		}
		out[len(out)-1] = '"' // replace last " " with `"`
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Int2VectorArrayToInt64SliceSlice struct {
	Val *[][]int64
}

func (v Int2VectorArrayToInt64SliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVectorArray(arr)
	int64ss := make([][]int64, len(elems))
	for i := 0; i < len(elems); i++ {
		int64s := make([]int64, len(elems[i]))
		for j := 0; j < len(elems[i]); j++ {
			i64, err := strconv.ParseInt(string(elems[i][j]), 10, 16)
			if err != nil {
				return err
			}
			int64s[j] = i64
		}
		int64ss[i] = int64s
	}

	*v.Val = int64ss
	return nil
}

type Int2VectorArrayFromUintSliceSlice struct {
	Val [][]uint
}

func (v Int2VectorArrayFromUintSliceSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, uints := range v.Val {
		out = append(out, '"')
		for _, u := range uints {
			out = strconv.AppendUint(out, uint64(u), 10)
			out = append(out, ' ')
		}
		out[len(out)-1] = '"' // replace last " " with `"`
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Int2VectorArrayToUintSliceSlice struct {
	Val *[][]uint
}

func (v Int2VectorArrayToUintSliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVectorArray(arr)
	uintss := make([][]uint, len(elems))
	for i := 0; i < len(elems); i++ {
		uints := make([]uint, len(elems[i]))
		for j := 0; j < len(elems[i]); j++ {
			u64, err := strconv.ParseUint(string(elems[i][j]), 10, 16)
			if err != nil {
				return err
			}
			uints[j] = uint(u64)
		}
		uintss[i] = uints
	}

	*v.Val = uintss
	return nil
}

type Int2VectorArrayFromUint8SliceSlice struct {
	Val [][]uint8
}

func (v Int2VectorArrayFromUint8SliceSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, uint8s := range v.Val {
		out = append(out, '"')
		for _, u8 := range uint8s {
			out = strconv.AppendUint(out, uint64(u8), 10)
			out = append(out, ' ')
		}
		out[len(out)-1] = '"' // replace last " " with `"`
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Int2VectorArrayToUint8SliceSlice struct {
	Val *[][]uint8
}

func (v Int2VectorArrayToUint8SliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVectorArray(arr)
	uint8ss := make([][]uint8, len(elems))
	for i := 0; i < len(elems); i++ {
		uint8s := make([]uint8, len(elems[i]))
		for j := 0; j < len(elems[i]); j++ {
			u64, err := strconv.ParseUint(string(elems[i][j]), 10, 8)
			if err != nil {
				return err
			}
			uint8s[j] = uint8(u64)
		}
		uint8ss[i] = uint8s
	}

	*v.Val = uint8ss
	return nil
}

type Int2VectorArrayFromUint16SliceSlice struct {
	Val [][]uint16
}

func (v Int2VectorArrayFromUint16SliceSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, uint16s := range v.Val {
		out = append(out, '"')
		for _, u16 := range uint16s {
			out = strconv.AppendUint(out, uint64(u16), 10)
			out = append(out, ' ')
		}
		out[len(out)-1] = '"' // replace last " " with `"`
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Int2VectorArrayToUint16SliceSlice struct {
	Val *[][]uint16
}

func (v Int2VectorArrayToUint16SliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVectorArray(arr)
	uint16ss := make([][]uint16, len(elems))
	for i := 0; i < len(elems); i++ {
		uint16s := make([]uint16, len(elems[i]))
		for j := 0; j < len(elems[i]); j++ {
			u64, err := strconv.ParseUint(string(elems[i][j]), 10, 16)
			if err != nil {
				return err
			}
			uint16s[j] = uint16(u64)
		}
		uint16ss[i] = uint16s
	}

	*v.Val = uint16ss
	return nil
}

type Int2VectorArrayFromUint32SliceSlice struct {
	Val [][]uint32
}

func (v Int2VectorArrayFromUint32SliceSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, uint32s := range v.Val {
		out = append(out, '"')
		for _, u32 := range uint32s {
			out = strconv.AppendUint(out, uint64(u32), 10)
			out = append(out, ' ')
		}
		out[len(out)-1] = '"' // replace last " " with `"`
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Int2VectorArrayToUint32SliceSlice struct {
	Val *[][]uint32
}

func (v Int2VectorArrayToUint32SliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVectorArray(arr)
	uint32ss := make([][]uint32, len(elems))
	for i := 0; i < len(elems); i++ {
		uint32s := make([]uint32, len(elems[i]))
		for j := 0; j < len(elems[i]); j++ {
			u64, err := strconv.ParseUint(string(elems[i][j]), 10, 16)
			if err != nil {
				return err
			}
			uint32s[j] = uint32(u64)
		}
		uint32ss[i] = uint32s
	}

	*v.Val = uint32ss
	return nil
}

type Int2VectorArrayFromUint64SliceSlice struct {
	Val [][]uint64
}

func (v Int2VectorArrayFromUint64SliceSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, uint64s := range v.Val {
		out = append(out, '"')
		for _, u64 := range uint64s {
			out = strconv.AppendUint(out, u64, 10)
			out = append(out, ' ')
		}
		out[len(out)-1] = '"' // replace last " " with `"`
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Int2VectorArrayToUint64SliceSlice struct {
	Val *[][]uint64
}

func (v Int2VectorArrayToUint64SliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVectorArray(arr)
	uint64ss := make([][]uint64, len(elems))
	for i := 0; i < len(elems); i++ {
		uint64s := make([]uint64, len(elems[i]))
		for j := 0; j < len(elems[i]); j++ {
			u64, err := strconv.ParseUint(string(elems[i][j]), 10, 16)
			if err != nil {
				return err
			}
			uint64s[j] = u64
		}
		uint64ss[i] = uint64s
	}

	*v.Val = uint64ss
	return nil
}

type Int2VectorArrayFromFloat32SliceSlice struct {
	Val [][]float32
}

func (v Int2VectorArrayFromFloat32SliceSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, float32s := range v.Val {
		out = append(out, '"')
		for _, f32 := range float32s {
			out = strconv.AppendInt(out, int64(f32), 10)
			out = append(out, ' ')
		}
		out[len(out)-1] = '"' // replace last " " with `"`
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Int2VectorArrayToFloat32SliceSlice struct {
	Val *[][]float32
}

func (v Int2VectorArrayToFloat32SliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVectorArray(arr)
	float32ss := make([][]float32, len(elems))
	for i := 0; i < len(elems); i++ {
		float32s := make([]float32, len(elems[i]))
		for j := 0; j < len(elems[i]); j++ {
			i64, err := strconv.ParseInt(string(elems[i][j]), 10, 16)
			if err != nil {
				return err
			}
			float32s[j] = float32(i64)
		}
		float32ss[i] = float32s
	}

	*v.Val = float32ss
	return nil
}

type Int2VectorArrayFromFloat64SliceSlice struct {
	Val [][]float64
}

func (v Int2VectorArrayFromFloat64SliceSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, float64s := range v.Val {
		out = append(out, '"')
		for _, f64 := range float64s {
			out = strconv.AppendInt(out, int64(f64), 10)
			out = append(out, ' ')
		}
		out[len(out)-1] = '"' // replace last " " with `"`
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Int2VectorArrayToFloat64SliceSlice struct {
	Val *[][]float64
}

func (v Int2VectorArrayToFloat64SliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseVectorArray(arr)
	float64ss := make([][]float64, len(elems))
	for i := 0; i < len(elems); i++ {
		float64s := make([]float64, len(elems[i]))
		for j := 0; j < len(elems[i]); j++ {
			i64, err := strconv.ParseInt(string(elems[i][j]), 10, 16)
			if err != nil {
				return err
			}
			float64s[j] = float64(i64)
		}
		float64ss[i] = float64s
	}

	*v.Val = float64ss
	return nil
}
