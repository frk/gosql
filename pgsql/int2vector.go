package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// Int2VectorFromIntSlice returns a driver.Valuer that produces a PostgreSQL int2vector from the given Go []int.
func Int2VectorFromIntSlice(val []int) driver.Valuer {
	return int2VectorFromIntSlice{val: val}
}

// Int2VectorToIntSlice returns an sql.Scanner that converts a PostgreSQL int2vector into a Go []int and sets it to val.
func Int2VectorToIntSlice(val *[]int) sql.Scanner {
	return int2VectorToIntSlice{val: val}
}

// Int2VectorFromInt8Slice returns a driver.Valuer that produces a PostgreSQL int2vector from the given Go []int8.
func Int2VectorFromInt8Slice(val []int8) driver.Valuer {
	return int2VectorFromInt8Slice{val: val}
}

// Int2VectorToInt8Slice returns an sql.Scanner that converts a PostgreSQL int2vector into a Go []int8 and sets it to val.
func Int2VectorToInt8Slice(val *[]int8) sql.Scanner {
	return int2VectorToInt8Slice{val: val}
}

// Int2VectorFromInt16Slice returns a driver.Valuer that produces a PostgreSQL int2vector from the given Go []int16.
func Int2VectorFromInt16Slice(val []int16) driver.Valuer {
	return int2VectorFromInt16Slice{val: val}
}

// Int2VectorToInt16Slice returns an sql.Scanner that converts a PostgreSQL int2vector into a Go []int16 and sets it to val.
func Int2VectorToInt16Slice(val *[]int16) sql.Scanner {
	return int2VectorToInt16Slice{val: val}
}

// Int2VectorFromInt32Slice returns a driver.Valuer that produces a PostgreSQL int2vector from the given Go []int32.
func Int2VectorFromInt32Slice(val []int32) driver.Valuer {
	return int2VectorFromInt32Slice{val: val}
}

// Int2VectorToInt32Slice returns an sql.Scanner that converts a PostgreSQL int2vector into a Go []int32 and sets it to val.
func Int2VectorToInt32Slice(val *[]int32) sql.Scanner {
	return int2VectorToInt32Slice{val: val}
}

// Int2VectorFromInt64Slice returns a driver.Valuer that produces a PostgreSQL int2vector from the given Go []int64.
func Int2VectorFromInt64Slice(val []int64) driver.Valuer {
	return int2VectorFromInt64Slice{val: val}
}

// Int2VectorToInt64Slice returns an sql.Scanner that converts a PostgreSQL int2vector into a Go []int64 and sets it to val.
func Int2VectorToInt64Slice(val *[]int64) sql.Scanner {
	return int2VectorToInt64Slice{val: val}
}

// Int2VectorFromUintSlice returns a driver.Valuer that produces a PostgreSQL int2vector from the given Go []uint.
func Int2VectorFromUintSlice(val []uint) driver.Valuer {
	return int2VectorFromUintSlice{val: val}
}

// Int2VectorToUintSlice returns an sql.Scanner that converts a PostgreSQL int2vector into a Go []uint and sets it to val.
func Int2VectorToUintSlice(val *[]uint) sql.Scanner {
	return int2VectorToUintSlice{val: val}
}

// Int2VectorFromUint8Slice returns a driver.Valuer that produces a PostgreSQL int2vector from the given Go []uint8.
func Int2VectorFromUint8Slice(val []uint8) driver.Valuer {
	return int2VectorFromUint8Slice{val: val}
}

// Int2VectorToUint8Slice returns an sql.Scanner that converts a PostgreSQL int2vector into a Go []uint8 and sets it to val.
func Int2VectorToUint8Slice(val *[]uint8) sql.Scanner {
	return int2VectorToUint8Slice{val: val}
}

// Int2VectorFromUint16Slice returns a driver.Valuer that produces a PostgreSQL int2vector from the given Go []uint16.
func Int2VectorFromUint16Slice(val []uint16) driver.Valuer {
	return int2VectorFromUint16Slice{val: val}
}

// Int2VectorToUint16Slice returns an sql.Scanner that converts a PostgreSQL int2vector into a Go []uint16 and sets it to val.
func Int2VectorToUint16Slice(val *[]uint16) sql.Scanner {
	return int2VectorToUint16Slice{val: val}
}

// Int2VectorFromUint32Slice returns a driver.Valuer that produces a PostgreSQL int2vector from the given Go []uint32.
func Int2VectorFromUint32Slice(val []uint32) driver.Valuer {
	return int2VectorFromUint32Slice{val: val}
}

// Int2VectorToUint32Slice returns an sql.Scanner that converts a PostgreSQL int2vector into a Go []uint32 and sets it to val.
func Int2VectorToUint32Slice(val *[]uint32) sql.Scanner {
	return int2VectorToUint32Slice{val: val}
}

// Int2VectorFromUint64Slice returns a driver.Valuer that produces a PostgreSQL int2vector from the given Go []uint64.
func Int2VectorFromUint64Slice(val []uint64) driver.Valuer {
	return int2VectorFromUint64Slice{val: val}
}

// Int2VectorToUint64Slice returns an sql.Scanner that converts a PostgreSQL int2vector into a Go []uint64 and sets it to val.
func Int2VectorToUint64Slice(val *[]uint64) sql.Scanner {
	return int2VectorToUint64Slice{val: val}
}

// Int2VectorFromFloat32Slice returns a driver.Valuer that produces a PostgreSQL int2vector from the given Go []float32.
func Int2VectorFromFloat32Slice(val []float32) driver.Valuer {
	return int2VectorFromFloat32Slice{val: val}
}

// Int2VectorToFloat32Slice returns an sql.Scanner that converts a PostgreSQL int2vector into a Go []float32 and sets it to val.
func Int2VectorToFloat32Slice(val *[]float32) sql.Scanner {
	return int2VectorToFloat32Slice{val: val}
}

// Int2VectorFromFloat64Slice returns a driver.Valuer that produces a PostgreSQL int2vector from the given Go []float64.
func Int2VectorFromFloat64Slice(val []float64) driver.Valuer {
	return int2VectorFromFloat64Slice{val: val}
}

// Int2VectorToFloat64Slice returns an sql.Scanner that converts a PostgreSQL int2vector into a Go []float64 and sets it to val.
func Int2VectorToFloat64Slice(val *[]float64) sql.Scanner {
	return int2VectorToFloat64Slice{val: val}
}

type int2VectorFromIntSlice struct {
	val []int
}

func (v int2VectorFromIntSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{}, nil
	}

	out := []byte{}

	for _, i := range v.val {
		out = strconv.AppendInt(out, int64(i), 10)
		out = append(out, ' ')
	}

	out = out[:len(out)-1] // drop last " "
	return out, nil
}

type int2VectorToIntSlice struct {
	val *[]int
}

func (v int2VectorToIntSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = ints
	return nil
}

type int2VectorFromInt8Slice struct {
	val []int8
}

func (v int2VectorFromInt8Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{}, nil
	}

	out := []byte{}

	for _, i8 := range v.val {
		out = strconv.AppendInt(out, int64(i8), 10)
		out = append(out, ' ')
	}

	out = out[:len(out)-1] // drop last " "
	return out, nil
}

type int2VectorToInt8Slice struct {
	val *[]int8
}

func (v int2VectorToInt8Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = int8s
	return nil
}

type int2VectorFromInt16Slice struct {
	val []int16
}

func (v int2VectorFromInt16Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{}, nil
	}

	out := []byte{}

	for _, i16 := range v.val {
		out = strconv.AppendInt(out, int64(i16), 10)
		out = append(out, ' ')
	}

	out = out[:len(out)-1] // drop last " "
	return out, nil
}

type int2VectorToInt16Slice struct {
	val *[]int16
}

func (v int2VectorToInt16Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = int16s
	return nil
}

type int2VectorFromInt32Slice struct {
	val []int32
}

func (v int2VectorFromInt32Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{}, nil
	}

	out := []byte{}

	for _, i32 := range v.val {
		out = strconv.AppendInt(out, int64(i32), 10)
		out = append(out, ' ')
	}

	out = out[:len(out)-1] // drop last " "
	return out, nil
}

type int2VectorToInt32Slice struct {
	val *[]int32
}

func (v int2VectorToInt32Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = int32s
	return nil
}

type int2VectorFromInt64Slice struct {
	val []int64
}

func (v int2VectorFromInt64Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{}, nil
	}

	out := []byte{}

	for _, i64 := range v.val {
		out = strconv.AppendInt(out, i64, 10)
		out = append(out, ' ')
	}

	out = out[:len(out)-1] // drop last " "
	return out, nil
}

type int2VectorToInt64Slice struct {
	val *[]int64
}

func (v int2VectorToInt64Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = int64s
	return nil
}

type int2VectorFromUintSlice struct {
	val []uint
}

func (v int2VectorFromUintSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{}, nil
	}

	out := []byte{}

	for _, u := range v.val {
		out = strconv.AppendUint(out, uint64(u), 10)
		out = append(out, ' ')
	}

	out = out[:len(out)-1] // drop last " "
	return out, nil
}

type int2VectorToUintSlice struct {
	val *[]uint
}

func (v int2VectorToUintSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = uints
	return nil
}

type int2VectorFromUint8Slice struct {
	val []uint8
}

func (v int2VectorFromUint8Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{}, nil
	}

	out := []byte{}

	for _, u8 := range v.val {
		out = strconv.AppendUint(out, uint64(u8), 10)
		out = append(out, ' ')
	}

	out = out[:len(out)-1] // drop last " "
	return out, nil
}

type int2VectorToUint8Slice struct {
	val *[]uint8
}

func (v int2VectorToUint8Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = uint8s
	return nil
}

type int2VectorFromUint16Slice struct {
	val []uint16
}

func (v int2VectorFromUint16Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{}, nil
	}

	out := []byte{}

	for _, u16 := range v.val {
		out = strconv.AppendUint(out, uint64(u16), 10)
		out = append(out, ' ')
	}

	out = out[:len(out)-1] // drop last " "
	return out, nil
}

type int2VectorToUint16Slice struct {
	val *[]uint16
}

func (v int2VectorToUint16Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = uint16s
	return nil
}

type int2VectorFromUint32Slice struct {
	val []uint32
}

func (v int2VectorFromUint32Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{}, nil
	}

	out := []byte{}

	for _, u32 := range v.val {
		out = strconv.AppendUint(out, uint64(u32), 10)
		out = append(out, ' ')
	}

	out = out[:len(out)-1] // drop last " "
	return out, nil
}

type int2VectorToUint32Slice struct {
	val *[]uint32
}

func (v int2VectorToUint32Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = uint32s
	return nil
}

type int2VectorFromUint64Slice struct {
	val []uint64
}

func (v int2VectorFromUint64Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{}, nil
	}

	out := []byte{}

	for _, u64 := range v.val {
		out = strconv.AppendUint(out, u64, 10)
		out = append(out, ' ')
	}

	out = out[:len(out)-1] // drop last " "
	return out, nil
}

type int2VectorToUint64Slice struct {
	val *[]uint64
}

func (v int2VectorToUint64Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = uint64s
	return nil
}

type int2VectorFromFloat32Slice struct {
	val []float32
}

func (v int2VectorFromFloat32Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{}, nil
	}

	out := []byte{}

	for _, f32 := range v.val {
		out = strconv.AppendInt(out, int64(f32), 10)
		out = append(out, ' ')
	}

	out = out[:len(out)-1] // drop last " "
	return out, nil
}

type int2VectorToFloat32Slice struct {
	val *[]float32
}

func (v int2VectorToFloat32Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = float32s
	return nil
}

type int2VectorFromFloat64Slice struct {
	val []float64
}

func (v int2VectorFromFloat64Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{}, nil
	}

	out := []byte{}

	for _, f64 := range v.val {
		out = strconv.AppendInt(out, int64(f64), 10)
		out = append(out, ' ')
	}

	out = out[:len(out)-1] // drop last " "
	return out, nil
}

type int2VectorToFloat64Slice struct {
	val *[]float64
}

func (v int2VectorToFloat64Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = float64s
	return nil
}
