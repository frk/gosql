package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// Int4ArrayFromIntSlice returns a driver.Valuer that produces a PostgreSQL int4[] from the given Go []int.
func Int4ArrayFromIntSlice(val []int) driver.Valuer {
	return int4ArrayFromIntSlice{val: val}
}

// Int4ArrayToIntSlice returns an sql.Scanner that converts a PostgreSQL int4[] into a Go []int and sets it to val.
func Int4ArrayToIntSlice(val *[]int) sql.Scanner {
	return int4ArrayToIntSlice{val: val}
}

// Int4ArrayFromInt8Slice returns a driver.Valuer that produces a PostgreSQL int4[] from the given Go []int8.
func Int4ArrayFromInt8Slice(val []int8) driver.Valuer {
	return int4ArrayFromInt8Slice{val: val}
}

// Int4ArrayToInt8Slice returns an sql.Scanner that converts a PostgreSQL int4[] into a Go []int8 and sets it to val.
func Int4ArrayToInt8Slice(val *[]int8) sql.Scanner {
	return int4ArrayToInt8Slice{val: val}
}

// Int4ArrayFromInt16Slice returns a driver.Valuer that produces a PostgreSQL int4[] from the given Go []int16.
func Int4ArrayFromInt16Slice(val []int16) driver.Valuer {
	return int4ArrayFromInt16Slice{val: val}
}

// Int4ArrayToInt16Slice returns an sql.Scanner that converts a PostgreSQL int4[] into a Go []int16 and sets it to val.
func Int4ArrayToInt16Slice(val *[]int16) sql.Scanner {
	return int4ArrayToInt16Slice{val: val}
}

// Int4ArrayFromInt32Slice returns a driver.Valuer that produces a PostgreSQL int4[] from the given Go []int32.
func Int4ArrayFromInt32Slice(val []int32) driver.Valuer {
	return int4ArrayFromInt32Slice{val: val}
}

// Int4ArrayToInt32Slice returns an sql.Scanner that converts a PostgreSQL int4[] into a Go []int32 and sets it to val.
func Int4ArrayToInt32Slice(val *[]int32) sql.Scanner {
	return int4ArrayToInt32Slice{val: val}
}

// Int4ArrayFromInt64Slice returns a driver.Valuer that produces a PostgreSQL int4[] from the given Go []int64.
func Int4ArrayFromInt64Slice(val []int64) driver.Valuer {
	return int4ArrayFromInt64Slice{val: val}
}

// Int4ArrayToInt64Slice returns an sql.Scanner that converts a PostgreSQL int4[] into a Go []int64 and sets it to val.
func Int4ArrayToInt64Slice(val *[]int64) sql.Scanner {
	return int4ArrayToInt64Slice{val: val}
}

// Int4ArrayFromUintSlice returns a driver.Valuer that produces a PostgreSQL int4[] from the given Go []uint.
func Int4ArrayFromUintSlice(val []uint) driver.Valuer {
	return int4ArrayFromUintSlice{val: val}
}

// Int4ArrayToUintSlice returns an sql.Scanner that converts a PostgreSQL int4[] into a Go []uint and sets it to val.
func Int4ArrayToUintSlice(val *[]uint) sql.Scanner {
	return int4ArrayToUintSlice{val: val}
}

// Int4ArrayFromUint8Slice returns a driver.Valuer that produces a PostgreSQL int4[] from the given Go []uint8.
func Int4ArrayFromUint8Slice(val []uint8) driver.Valuer {
	return int4ArrayFromUint8Slice{val: val}
}

// Int4ArrayToUint8Slice returns an sql.Scanner that converts a PostgreSQL int4[] into a Go []uint8 and sets it to val.
func Int4ArrayToUint8Slice(val *[]uint8) sql.Scanner {
	return int4ArrayToUint8Slice{val: val}
}

// Int4ArrayFromUint16Slice returns a driver.Valuer that produces a PostgreSQL int4[] from the given Go []uint16.
func Int4ArrayFromUint16Slice(val []uint16) driver.Valuer {
	return int4ArrayFromUint16Slice{val: val}
}

// Int4ArrayToUint16Slice returns an sql.Scanner that converts a PostgreSQL int4[] into a Go []uint16 and sets it to val.
func Int4ArrayToUint16Slice(val *[]uint16) sql.Scanner {
	return int4ArrayToUint16Slice{val: val}
}

// Int4ArrayFromUint32Slice returns a driver.Valuer that produces a PostgreSQL int4[] from the given Go []uint32.
func Int4ArrayFromUint32Slice(val []uint32) driver.Valuer {
	return int4ArrayFromUint32Slice{val: val}
}

// Int4ArrayToUint32Slice returns an sql.Scanner that converts a PostgreSQL int4[] into a Go []uint32 and sets it to val.
func Int4ArrayToUint32Slice(val *[]uint32) sql.Scanner {
	return int4ArrayToUint32Slice{val: val}
}

// Int4ArrayFromUint64Slice returns a driver.Valuer that produces a PostgreSQL int4[] from the given Go []uint64.
func Int4ArrayFromUint64Slice(val []uint64) driver.Valuer {
	return int4ArrayFromUint64Slice{val: val}
}

// Int4ArrayToUint64Slice returns an sql.Scanner that converts a PostgreSQL int4[] into a Go []uint64 and sets it to val.
func Int4ArrayToUint64Slice(val *[]uint64) sql.Scanner {
	return int4ArrayToUint64Slice{val: val}
}

// Int4ArrayFromFloat32Slice returns a driver.Valuer that produces a PostgreSQL int4[] from the given Go []float32.
func Int4ArrayFromFloat32Slice(val []float32) driver.Valuer {
	return int4ArrayFromFloat32Slice{val: val}
}

// Int4ArrayToFloat32Slice returns an sql.Scanner that converts a PostgreSQL int4[] into a Go []float32 and sets it to val.
func Int4ArrayToFloat32Slice(val *[]float32) sql.Scanner {
	return int4ArrayToFloat32Slice{val: val}
}

// Int4ArrayFromFloat64Slice returns a driver.Valuer that produces a PostgreSQL int4[] from the given Go []float64.
func Int4ArrayFromFloat64Slice(val []float64) driver.Valuer {
	return int4ArrayFromFloat64Slice{val: val}
}

// Int4ArrayToFloat64Slice returns an sql.Scanner that converts a PostgreSQL int4[] into a Go []float64 and sets it to val.
func Int4ArrayToFloat64Slice(val *[]float64) sql.Scanner {
	return int4ArrayToFloat64Slice{val: val}
}

type int4ArrayFromIntSlice struct {
	val []int
}

func (v int4ArrayFromIntSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, i := range v.val {
		out = strconv.AppendInt(out, int64(i), 10)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type int4ArrayToIntSlice struct {
	val *[]int
}

func (v int4ArrayToIntSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	ints := make([]int, len(elems))
	for i := 0; i < len(elems); i++ {
		i64, err := strconv.ParseInt(string(elems[i]), 10, 32)
		if err != nil {
			return err
		}
		ints[i] = int(i64)
	}

	*v.val = ints
	return nil
}

type int4ArrayFromInt8Slice struct {
	val []int8
}

func (v int4ArrayFromInt8Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, i8 := range v.val {
		out = strconv.AppendInt(out, int64(i8), 10)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type int4ArrayToInt8Slice struct {
	val *[]int8
}

func (v int4ArrayToInt8Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = int8s
	return nil
}

type int4ArrayFromInt16Slice struct {
	val []int16
}

func (v int4ArrayFromInt16Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, i16 := range v.val {
		out = strconv.AppendInt(out, int64(i16), 10)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type int4ArrayToInt16Slice struct {
	val *[]int16
}

func (v int4ArrayToInt16Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = int16s
	return nil
}

type int4ArrayFromInt32Slice struct {
	val []int32
}

func (v int4ArrayFromInt32Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, i32 := range v.val {
		out = strconv.AppendInt(out, int64(i32), 10)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type int4ArrayToInt32Slice struct {
	val *[]int32
}

func (v int4ArrayToInt32Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	int32s := make([]int32, len(elems))
	for i := 0; i < len(elems); i++ {
		i64, err := strconv.ParseInt(string(elems[i]), 10, 32)
		if err != nil {
			return err
		}
		int32s[i] = int32(i64)
	}

	*v.val = int32s
	return nil
}

type int4ArrayFromInt64Slice struct {
	val []int64
}

func (v int4ArrayFromInt64Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, i64 := range v.val {
		out = strconv.AppendInt(out, i64, 10)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type int4ArrayToInt64Slice struct {
	val *[]int64
}

func (v int4ArrayToInt64Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	int64s := make([]int64, len(elems))
	for i := 0; i < len(elems); i++ {
		i64, err := strconv.ParseInt(string(elems[i]), 10, 32)
		if err != nil {
			return err
		}
		int64s[i] = i64
	}

	*v.val = int64s
	return nil
}

////////////////////////////////////////////////////////////////////////////////

type int4ArrayFromUintSlice struct {
	val []uint
}

func (v int4ArrayFromUintSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, u := range v.val {
		out = strconv.AppendUint(out, uint64(u), 10)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type int4ArrayToUintSlice struct {
	val *[]uint
}

func (v int4ArrayToUintSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	uints := make([]uint, len(elems))
	for i := 0; i < len(elems); i++ {
		u64, err := strconv.ParseUint(string(elems[i]), 10, 32)
		if err != nil {
			return err
		}
		uints[i] = uint(u64)
	}

	*v.val = uints
	return nil
}

type int4ArrayFromUint8Slice struct {
	val []uint8
}

func (v int4ArrayFromUint8Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, u := range v.val {
		out = strconv.AppendUint(out, uint64(u), 10)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type int4ArrayToUint8Slice struct {
	val *[]uint8
}

func (v int4ArrayToUint8Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = uint8s
	return nil
}

type int4ArrayFromUint16Slice struct {
	val []uint16
}

func (v int4ArrayFromUint16Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, u := range v.val {
		out = strconv.AppendUint(out, uint64(u), 10)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type int4ArrayToUint16Slice struct {
	val *[]uint16
}

func (v int4ArrayToUint16Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
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

	*v.val = uint16s
	return nil
}

type int4ArrayFromUint32Slice struct {
	val []uint32
}

func (v int4ArrayFromUint32Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, u := range v.val {
		out = strconv.AppendUint(out, uint64(u), 10)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type int4ArrayToUint32Slice struct {
	val *[]uint32
}

func (v int4ArrayToUint32Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	uint32s := make([]uint32, len(elems))
	for i := 0; i < len(elems); i++ {
		u64, err := strconv.ParseUint(string(elems[i]), 10, 32)
		if err != nil {
			return err
		}
		uint32s[i] = uint32(u64)
	}

	*v.val = uint32s
	return nil
}

type int4ArrayFromUint64Slice struct {
	val []uint64
}

func (v int4ArrayFromUint64Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, u := range v.val {
		out = strconv.AppendUint(out, u, 10)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type int4ArrayToUint64Slice struct {
	val *[]uint64
}

func (v int4ArrayToUint64Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	uint64s := make([]uint64, len(elems))
	for i := 0; i < len(elems); i++ {
		u64, err := strconv.ParseUint(string(elems[i]), 10, 32)
		if err != nil {
			return err
		}
		uint64s[i] = u64
	}

	*v.val = uint64s
	return nil
}

type int4ArrayFromFloat32Slice struct {
	val []float32
}

func (v int4ArrayFromFloat32Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, f := range v.val {
		out = strconv.AppendInt(out, int64(f), 10)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type int4ArrayToFloat32Slice struct {
	val *[]float32
}

func (v int4ArrayToFloat32Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	float32s := make([]float32, len(elems))
	for i := 0; i < len(elems); i++ {
		i64, err := strconv.ParseInt(string(elems[i]), 10, 32)
		if err != nil {
			return err
		}
		float32s[i] = float32(i64)
	}

	*v.val = float32s
	return nil
}

type int4ArrayFromFloat64Slice struct {
	val []float64
}

func (v int4ArrayFromFloat64Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, f := range v.val {
		out = strconv.AppendInt(out, int64(f), 10)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type int4ArrayToFloat64Slice struct {
	val *[]float64
}

func (v int4ArrayToFloat64Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseCommaArray(arr)
	float64s := make([]float64, len(elems))
	for i := 0; i < len(elems); i++ {
		i64, err := strconv.ParseInt(string(elems[i]), 10, 32)
		if err != nil {
			return err
		}
		float64s[i] = float64(i64)
	}

	*v.val = float64s
	return nil
}
