package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// Int2ArrayFromIntSlice returns a driver.Valuer that produces a PostgreSQL int2[] from the given Go []int.
func Int2ArrayFromIntSlice(val []int) driver.Valuer {
	return int2ArrayFromIntSlice{val: val}
}

// Int2ArrayToIntSlice returns an sql.Scanner that converts a PostgreSQL int2[] into a Go []int and sets it to val.
func Int2ArrayToIntSlice(val *[]int) sql.Scanner {
	return int2ArrayToIntSlice{val: val}
}

// Int2ArrayFromInt8Slice returns a driver.Valuer that produces a PostgreSQL int2[] from the given Go []int8.
func Int2ArrayFromInt8Slice(val []int8) driver.Valuer {
	return int2ArrayFromInt8Slice{val: val}
}

// Int2ArrayToInt8Slice returns an sql.Scanner that converts a PostgreSQL int2[] into a Go []int8 and sets it to val.
func Int2ArrayToInt8Slice(val *[]int8) sql.Scanner {
	return int2ArrayToInt8Slice{val: val}
}

// Int2ArrayFromInt16Slice returns a driver.Valuer that produces a PostgreSQL int2[] from the given Go []int16.
func Int2ArrayFromInt16Slice(val []int16) driver.Valuer {
	return int2ArrayFromInt16Slice{val: val}
}

// Int2ArrayToInt16Slice returns an sql.Scanner that converts a PostgreSQL int2[] into a Go []int16 and sets it to val.
func Int2ArrayToInt16Slice(val *[]int16) sql.Scanner {
	return int2ArrayToInt16Slice{val: val}
}

// Int2ArrayFromInt32Slice returns a driver.Valuer that produces a PostgreSQL int2[] from the given Go []int32.
func Int2ArrayFromInt32Slice(val []int32) driver.Valuer {
	return int2ArrayFromInt32Slice{val: val}
}

// Int2ArrayToInt32Slice returns an sql.Scanner that converts a PostgreSQL int2[] into a Go []int32 and sets it to val.
func Int2ArrayToInt32Slice(val *[]int32) sql.Scanner {
	return int2ArrayToInt32Slice{val: val}
}

// Int2ArrayFromInt64Slice returns a driver.Valuer that produces a PostgreSQL int2[] from the given Go []int64.
func Int2ArrayFromInt64Slice(val []int64) driver.Valuer {
	return int2ArrayFromInt64Slice{val: val}
}

// Int2ArrayToInt64Slice returns an sql.Scanner that converts a PostgreSQL int2[] into a Go []int64 and sets it to val.
func Int2ArrayToInt64Slice(val *[]int64) sql.Scanner {
	return int2ArrayToInt64Slice{val: val}
}

// Int2ArrayFromUintSlice returns a driver.Valuer that produces a PostgreSQL int2[] from the given Go []uint.
func Int2ArrayFromUintSlice(val []uint) driver.Valuer {
	return int2ArrayFromUintSlice{val: val}
}

// Int2ArrayToUintSlice returns an sql.Scanner that converts a PostgreSQL int2[] into a Go []uint and sets it to val.
func Int2ArrayToUintSlice(val *[]uint) sql.Scanner {
	return int2ArrayToUintSlice{val: val}
}

// Int2ArrayFromUint8Slice returns a driver.Valuer that produces a PostgreSQL int2[] from the given Go []uint8.
func Int2ArrayFromUint8Slice(val []uint8) driver.Valuer {
	return int2ArrayFromUint8Slice{val: val}
}

// Int2ArrayToUint8Slice returns an sql.Scanner that converts a PostgreSQL int2[] into a Go []uint8 and sets it to val.
func Int2ArrayToUint8Slice(val *[]uint8) sql.Scanner {
	return int2ArrayToUint8Slice{val: val}
}

// Int2ArrayFromUint16Slice returns a driver.Valuer that produces a PostgreSQL int2[] from the given Go []uint16.
func Int2ArrayFromUint16Slice(val []uint16) driver.Valuer {
	return int2ArrayFromUint16Slice{val: val}
}

// Int2ArrayToUint16Slice returns an sql.Scanner that converts a PostgreSQL int2[] into a Go []uint16 and sets it to val.
func Int2ArrayToUint16Slice(val *[]uint16) sql.Scanner {
	return int2ArrayToUint16Slice{val: val}
}

// Int2ArrayFromUint32Slice returns a driver.Valuer that produces a PostgreSQL int2[] from the given Go []uint32.
func Int2ArrayFromUint32Slice(val []uint32) driver.Valuer {
	return int2ArrayFromUint32Slice{val: val}
}

// Int2ArrayToUint32Slice returns an sql.Scanner that converts a PostgreSQL int2[] into a Go []uint32 and sets it to val.
func Int2ArrayToUint32Slice(val *[]uint32) sql.Scanner {
	return int2ArrayToUint32Slice{val: val}
}

// Int2ArrayFromUint64Slice returns a driver.Valuer that produces a PostgreSQL int2[] from the given Go []uint64.
func Int2ArrayFromUint64Slice(val []uint64) driver.Valuer {
	return int2ArrayFromUint64Slice{val: val}
}

// Int2ArrayToUint64Slice returns an sql.Scanner that converts a PostgreSQL int2[] into a Go []uint64 and sets it to val.
func Int2ArrayToUint64Slice(val *[]uint64) sql.Scanner {
	return int2ArrayToUint64Slice{val: val}
}

// Int2ArrayFromFloat32Slice returns a driver.Valuer that produces a PostgreSQL int2[] from the given Go []float32.
func Int2ArrayFromFloat32Slice(val []float32) driver.Valuer {
	return int2ArrayFromFloat32Slice{val: val}
}

// Int2ArrayToFloat32Slice returns an sql.Scanner that converts a PostgreSQL int2[] into a Go []float32 and sets it to val.
func Int2ArrayToFloat32Slice(val *[]float32) sql.Scanner {
	return int2ArrayToFloat32Slice{val: val}
}

// Int2ArrayFromFloat64Slice returns a driver.Valuer that produces a PostgreSQL int2[] from the given Go []float64.
func Int2ArrayFromFloat64Slice(val []float64) driver.Valuer {
	return int2ArrayFromFloat64Slice{val: val}
}

// Int2ArrayToFloat64Slice returns an sql.Scanner that converts a PostgreSQL int2[] into a Go []float64 and sets it to val.
func Int2ArrayToFloat64Slice(val *[]float64) sql.Scanner {
	return int2ArrayToFloat64Slice{val: val}
}

type int2ArrayFromIntSlice struct {
	val []int
}

func (v int2ArrayFromIntSlice) Value() (driver.Value, error) {
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

type int2ArrayToIntSlice struct {
	val *[]int
}

func (v int2ArrayToIntSlice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		ints[i] = int(i64)
	}

	*v.val = ints
	return nil
}

type int2ArrayFromInt8Slice struct {
	val []int8
}

func (v int2ArrayFromInt8Slice) Value() (driver.Value, error) {
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

type int2ArrayToInt8Slice struct {
	val *[]int8
}

func (v int2ArrayToInt8Slice) Scan(src interface{}) error {
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

type int2ArrayFromInt16Slice struct {
	val []int16
}

func (v int2ArrayFromInt16Slice) Value() (driver.Value, error) {
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

type int2ArrayToInt16Slice struct {
	val *[]int16
}

func (v int2ArrayToInt16Slice) Scan(src interface{}) error {
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

type int2ArrayFromInt32Slice struct {
	val []int32
}

func (v int2ArrayFromInt32Slice) Value() (driver.Value, error) {
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

type int2ArrayToInt32Slice struct {
	val *[]int32
}

func (v int2ArrayToInt32Slice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		int32s[i] = int32(i64)
	}

	*v.val = int32s
	return nil
}

type int2ArrayFromInt64Slice struct {
	val []int64
}

func (v int2ArrayFromInt64Slice) Value() (driver.Value, error) {
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

type int2ArrayToInt64Slice struct {
	val *[]int64
}

func (v int2ArrayToInt64Slice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		int64s[i] = i64
	}

	*v.val = int64s
	return nil
}

////////////////////////////////////////////////////////////////////////////////

type int2ArrayFromUintSlice struct {
	val []uint
}

func (v int2ArrayFromUintSlice) Value() (driver.Value, error) {
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

type int2ArrayToUintSlice struct {
	val *[]uint
}

func (v int2ArrayToUintSlice) Scan(src interface{}) error {
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
		u64, err := strconv.ParseUint(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		uints[i] = uint(u64)
	}

	*v.val = uints
	return nil
}

type int2ArrayFromUint8Slice struct {
	val []uint8
}

func (v int2ArrayFromUint8Slice) Value() (driver.Value, error) {
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

type int2ArrayToUint8Slice struct {
	val *[]uint8
}

func (v int2ArrayToUint8Slice) Scan(src interface{}) error {
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

type int2ArrayFromUint16Slice struct {
	val []uint16
}

func (v int2ArrayFromUint16Slice) Value() (driver.Value, error) {
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

type int2ArrayToUint16Slice struct {
	val *[]uint16
}

func (v int2ArrayToUint16Slice) Scan(src interface{}) error {
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

type int2ArrayFromUint32Slice struct {
	val []uint32
}

func (v int2ArrayFromUint32Slice) Value() (driver.Value, error) {
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

type int2ArrayToUint32Slice struct {
	val *[]uint32
}

func (v int2ArrayToUint32Slice) Scan(src interface{}) error {
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
		u64, err := strconv.ParseUint(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		uint32s[i] = uint32(u64)
	}

	*v.val = uint32s
	return nil
}

type int2ArrayFromUint64Slice struct {
	val []uint64
}

func (v int2ArrayFromUint64Slice) Value() (driver.Value, error) {
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

type int2ArrayToUint64Slice struct {
	val *[]uint64
}

func (v int2ArrayToUint64Slice) Scan(src interface{}) error {
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
		u64, err := strconv.ParseUint(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		uint64s[i] = u64
	}

	*v.val = uint64s
	return nil
}

type int2ArrayFromFloat32Slice struct {
	val []float32
}

func (v int2ArrayFromFloat32Slice) Value() (driver.Value, error) {
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

type int2ArrayToFloat32Slice struct {
	val *[]float32
}

func (v int2ArrayToFloat32Slice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		float32s[i] = float32(i64)
	}

	*v.val = float32s
	return nil
}

type int2ArrayFromFloat64Slice struct {
	val []float64
}

func (v int2ArrayFromFloat64Slice) Value() (driver.Value, error) {
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

type int2ArrayToFloat64Slice struct {
	val *[]float64
}

func (v int2ArrayToFloat64Slice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 16)
		if err != nil {
			return err
		}
		float64s[i] = float64(i64)
	}

	*v.val = float64s
	return nil
}
