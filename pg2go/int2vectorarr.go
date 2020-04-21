package pg2go

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// Int2VectorArrayFromIntSliceSlice returns a driver.Valuer that produces a PostgreSQL int2vector[] from the given Go [][]int.
func Int2VectorArrayFromIntSliceSlice(val [][]int) driver.Valuer {
	return int2VectorArrayFromIntSliceSlice{val: val}
}

// Int2VectorArrayToIntSliceSlice returns an sql.Scanner that converts a PostgreSQL int2vector[] into a Go [][]int and sets it to val.
func Int2VectorArrayToIntSliceSlice(val *[][]int) sql.Scanner {
	return int2VectorArrayToIntSliceSlice{val: val}
}

// Int2VectorArrayFromInt8SliceSlice returns a driver.Valuer that produces a PostgreSQL int2vector[] from the given Go [][]int8.
func Int2VectorArrayFromInt8SliceSlice(val [][]int8) driver.Valuer {
	return int2VectorArrayFromInt8SliceSlice{val: val}
}

// Int2VectorArrayToInt8SliceSlice returns an sql.Scanner that converts a PostgreSQL int2vector[] into a Go [][]int8 and sets it to val.
func Int2VectorArrayToInt8SliceSlice(val *[][]int8) sql.Scanner {
	return int2VectorArrayToInt8SliceSlice{val: val}
}

// Int2VectorArrayFromInt16SliceSlice returns a driver.Valuer that produces a PostgreSQL int2vector[] from the given Go [][]int16.
func Int2VectorArrayFromInt16SliceSlice(val [][]int16) driver.Valuer {
	return int2VectorArrayFromInt16SliceSlice{val: val}
}

// Int2VectorArrayToInt16SliceSlice returns an sql.Scanner that converts a PostgreSQL int2vector[] into a Go [][]int16 and sets it to val.
func Int2VectorArrayToInt16SliceSlice(val *[][]int16) sql.Scanner {
	return int2VectorArrayToInt16SliceSlice{val: val}
}

// Int2VectorArrayFromInt32SliceSlice returns a driver.Valuer that produces a PostgreSQL int2vector[] from the given Go [][]int32.
func Int2VectorArrayFromInt32SliceSlice(val [][]int32) driver.Valuer {
	return int2VectorArrayFromInt32SliceSlice{val: val}
}

// Int2VectorArrayToInt32SliceSlice returns an sql.Scanner that converts a PostgreSQL int2vector[] into a Go [][]int32 and sets it to val.
func Int2VectorArrayToInt32SliceSlice(val *[][]int32) sql.Scanner {
	return int2VectorArrayToInt32SliceSlice{val: val}
}

// Int2VectorArrayFromInt64SliceSlice returns a driver.Valuer that produces a PostgreSQL int2vector[] from the given Go [][]int64.
func Int2VectorArrayFromInt64SliceSlice(val [][]int64) driver.Valuer {
	return int2VectorArrayFromInt64SliceSlice{val: val}
}

// Int2VectorArrayToInt64SliceSlice returns an sql.Scanner that converts a PostgreSQL int2vector[] into a Go [][]int64 and sets it to val.
func Int2VectorArrayToInt64SliceSlice(val *[][]int64) sql.Scanner {
	return int2VectorArrayToInt64SliceSlice{val: val}
}

// Int2VectorArrayFromUintSliceSlice returns a driver.Valuer that produces a PostgreSQL int2vector[] from the given Go [][]uint.
func Int2VectorArrayFromUintSliceSlice(val [][]uint) driver.Valuer {
	return int2VectorArrayFromUintSliceSlice{val: val}
}

// Int2VectorArrayToUintSliceSlice returns an sql.Scanner that converts a PostgreSQL int2vector[] into a Go [][]uint and sets it to val.
func Int2VectorArrayToUintSliceSlice(val *[][]uint) sql.Scanner {
	return int2VectorArrayToUintSliceSlice{val: val}
}

// Int2VectorArrayFromUint8SliceSlice returns a driver.Valuer that produces a PostgreSQL int2vector[] from the given Go [][]uint8.
func Int2VectorArrayFromUint8SliceSlice(val [][]uint8) driver.Valuer {
	return int2VectorArrayFromUint8SliceSlice{val: val}
}

// Int2VectorArrayToUint8SliceSlice returns an sql.Scanner that converts a PostgreSQL int2vector[] into a Go [][]uint8 and sets it to val.
func Int2VectorArrayToUint8SliceSlice(val *[][]uint8) sql.Scanner {
	return int2VectorArrayToUint8SliceSlice{val: val}
}

// Int2VectorArrayFromUint16SliceSlice returns a driver.Valuer that produces a PostgreSQL int2vector[] from the given Go [][]uint16.
func Int2VectorArrayFromUint16SliceSlice(val [][]uint16) driver.Valuer {
	return int2VectorArrayFromUint16SliceSlice{val: val}
}

// Int2VectorArrayToUint16SliceSlice returns an sql.Scanner that converts a PostgreSQL int2vector[] into a Go [][]uint16 and sets it to val.
func Int2VectorArrayToUint16SliceSlice(val *[][]uint16) sql.Scanner {
	return int2VectorArrayToUint16SliceSlice{val: val}
}

// Int2VectorArrayFromUint32SliceSlice returns a driver.Valuer that produces a PostgreSQL int2vector[] from the given Go [][]uint32.
func Int2VectorArrayFromUint32SliceSlice(val [][]uint32) driver.Valuer {
	return int2VectorArrayFromUint32SliceSlice{val: val}
}

// Int2VectorArrayToUint32SliceSlice returns an sql.Scanner that converts a PostgreSQL int2vector[] into a Go [][]uint32 and sets it to val.
func Int2VectorArrayToUint32SliceSlice(val *[][]uint32) sql.Scanner {
	return int2VectorArrayToUint32SliceSlice{val: val}
}

// Int2VectorArrayFromUint64SliceSlice returns a driver.Valuer that produces a PostgreSQL int2vector[] from the given Go [][]uint64.
func Int2VectorArrayFromUint64SliceSlice(val [][]uint64) driver.Valuer {
	return int2VectorArrayFromUint64SliceSlice{val: val}
}

// Int2VectorArrayToUint64SliceSlice returns an sql.Scanner that converts a PostgreSQL int2vector[] into a Go [][]uint64 and sets it to val.
func Int2VectorArrayToUint64SliceSlice(val *[][]uint64) sql.Scanner {
	return int2VectorArrayToUint64SliceSlice{val: val}
}

// Int2VectorArrayFromFloat32SliceSlice returns a driver.Valuer that produces a PostgreSQL int2vector[] from the given Go [][]float32.
func Int2VectorArrayFromFloat32SliceSlice(val [][]float32) driver.Valuer {
	return int2VectorArrayFromFloat32SliceSlice{val: val}
}

// Int2VectorArrayToFloat32SliceSlice returns an sql.Scanner that converts a PostgreSQL int2vector[] into a Go [][]float32 and sets it to val.
func Int2VectorArrayToFloat32SliceSlice(val *[][]float32) sql.Scanner {
	return int2VectorArrayToFloat32SliceSlice{val: val}
}

// Int2VectorArrayFromFloat64SliceSlice returns a driver.Valuer that produces a PostgreSQL int2vector[] from the given Go [][]float64.
func Int2VectorArrayFromFloat64SliceSlice(val [][]float64) driver.Valuer {
	return int2VectorArrayFromFloat64SliceSlice{val: val}
}

// Int2VectorArrayToFloat64SliceSlice returns an sql.Scanner that converts a PostgreSQL int2vector[] into a Go [][]float64 and sets it to val.
func Int2VectorArrayToFloat64SliceSlice(val *[][]float64) sql.Scanner {
	return int2VectorArrayToFloat64SliceSlice{val: val}
}

type int2VectorArrayFromIntSliceSlice struct {
	val [][]int
}

func (v int2VectorArrayFromIntSliceSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, ints := range v.val {
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

type int2VectorArrayToIntSliceSlice struct {
	val *[][]int
}

func (v int2VectorArrayToIntSliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = intss
	return nil
}

type int2VectorArrayFromInt8SliceSlice struct {
	val [][]int8
}

func (v int2VectorArrayFromInt8SliceSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, int8s := range v.val {
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

type int2VectorArrayToInt8SliceSlice struct {
	val *[][]int8
}

func (v int2VectorArrayToInt8SliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = int8ss
	return nil
}

type int2VectorArrayFromInt16SliceSlice struct {
	val [][]int16
}

func (v int2VectorArrayFromInt16SliceSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, int16s := range v.val {
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

type int2VectorArrayToInt16SliceSlice struct {
	val *[][]int16
}

func (v int2VectorArrayToInt16SliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = int16ss
	return nil
}

type int2VectorArrayFromInt32SliceSlice struct {
	val [][]int32
}

func (v int2VectorArrayFromInt32SliceSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, int32s := range v.val {
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

type int2VectorArrayToInt32SliceSlice struct {
	val *[][]int32
}

func (v int2VectorArrayToInt32SliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = int32ss
	return nil
}

type int2VectorArrayFromInt64SliceSlice struct {
	val [][]int64
}

func (v int2VectorArrayFromInt64SliceSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, int64s := range v.val {
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

type int2VectorArrayToInt64SliceSlice struct {
	val *[][]int64
}

func (v int2VectorArrayToInt64SliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = int64ss
	return nil
}

type int2VectorArrayFromUintSliceSlice struct {
	val [][]uint
}

func (v int2VectorArrayFromUintSliceSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, uints := range v.val {
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

type int2VectorArrayToUintSliceSlice struct {
	val *[][]uint
}

func (v int2VectorArrayToUintSliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = uintss
	return nil
}

type int2VectorArrayFromUint8SliceSlice struct {
	val [][]uint8
}

func (v int2VectorArrayFromUint8SliceSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, uint8s := range v.val {
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

type int2VectorArrayToUint8SliceSlice struct {
	val *[][]uint8
}

func (v int2VectorArrayToUint8SliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = uint8ss
	return nil
}

type int2VectorArrayFromUint16SliceSlice struct {
	val [][]uint16
}

func (v int2VectorArrayFromUint16SliceSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, uint16s := range v.val {
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

type int2VectorArrayToUint16SliceSlice struct {
	val *[][]uint16
}

func (v int2VectorArrayToUint16SliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = uint16ss
	return nil
}

type int2VectorArrayFromUint32SliceSlice struct {
	val [][]uint32
}

func (v int2VectorArrayFromUint32SliceSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, uint32s := range v.val {
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

type int2VectorArrayToUint32SliceSlice struct {
	val *[][]uint32
}

func (v int2VectorArrayToUint32SliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = uint32ss
	return nil
}

type int2VectorArrayFromUint64SliceSlice struct {
	val [][]uint64
}

func (v int2VectorArrayFromUint64SliceSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, uint64s := range v.val {
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

type int2VectorArrayToUint64SliceSlice struct {
	val *[][]uint64
}

func (v int2VectorArrayToUint64SliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = uint64ss
	return nil
}

type int2VectorArrayFromFloat32SliceSlice struct {
	val [][]float32
}

func (v int2VectorArrayFromFloat32SliceSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, float32s := range v.val {
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

type int2VectorArrayToFloat32SliceSlice struct {
	val *[][]float32
}

func (v int2VectorArrayToFloat32SliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = float32ss
	return nil
}

type int2VectorArrayFromFloat64SliceSlice struct {
	val [][]float64
}

func (v int2VectorArrayFromFloat64SliceSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, float64s := range v.val {
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

type int2VectorArrayToFloat64SliceSlice struct {
	val *[][]float64
}

func (v int2VectorArrayToFloat64SliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = float64ss
	return nil
}
