package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// NumericArrayFromIntSlice returns a driver.Valuer that produces a PostgreSQL numeric[] from the given Go []int.
func NumericArrayFromIntSlice(val []int) driver.Valuer {
	return numericArrayFromIntSlice{val: val}
}

// NumericArrayToIntSlice returns an sql.Scanner that converts a PostgreSQL numeric[] into a Go []int and sets it to val.
func NumericArrayToIntSlice(val *[]int) sql.Scanner {
	return numericArrayToIntSlice{val: val}
}

// NumericArrayFromInt8Slice returns a driver.Valuer that produces a PostgreSQL numeric[] from the given Go []int8.
func NumericArrayFromInt8Slice(val []int8) driver.Valuer {
	return numericArrayFromInt8Slice{val: val}
}

// NumericArrayToInt8Slice returns an sql.Scanner that converts a PostgreSQL numeric[] into a Go []int8 and sets it to val.
func NumericArrayToInt8Slice(val *[]int8) sql.Scanner {
	return numericArrayToInt8Slice{val: val}
}

// NumericArrayFromInt16Slice returns a driver.Valuer that produces a PostgreSQL numeric[] from the given Go []int16.
func NumericArrayFromInt16Slice(val []int16) driver.Valuer {
	return numericArrayFromInt16Slice{val: val}
}

// NumericArrayToInt16Slice returns an sql.Scanner that converts a PostgreSQL numeric[] into a Go []int16 and sets it to val.
func NumericArrayToInt16Slice(val *[]int16) sql.Scanner {
	return numericArrayToInt16Slice{val: val}
}

// NumericArrayFromInt32Slice returns a driver.Valuer that produces a PostgreSQL numeric[] from the given Go []int32.
func NumericArrayFromInt32Slice(val []int32) driver.Valuer {
	return numericArrayFromInt32Slice{val: val}
}

// NumericArrayToInt32Slice returns an sql.Scanner that converts a PostgreSQL numeric[] into a Go []int32 and sets it to val.
func NumericArrayToInt32Slice(val *[]int32) sql.Scanner {
	return numericArrayToInt32Slice{val: val}
}

// NumericArrayFromInt64Slice returns a driver.Valuer that produces a PostgreSQL numeric[] from the given Go []int64.
func NumericArrayFromInt64Slice(val []int64) driver.Valuer {
	return numericArrayFromInt64Slice{val: val}
}

// NumericArrayToInt64Slice returns an sql.Scanner that converts a PostgreSQL numeric[] into a Go []int64 and sets it to val.
func NumericArrayToInt64Slice(val *[]int64) sql.Scanner {
	return numericArrayToInt64Slice{val: val}
}

// NumericArrayFromUintSlice returns a driver.Valuer that produces a PostgreSQL numeric[] from the given Go []uint.
func NumericArrayFromUintSlice(val []uint) driver.Valuer {
	return numericArrayFromUintSlice{val: val}
}

// NumericArrayToUintSlice returns an sql.Scanner that converts a PostgreSQL numeric[] into a Go []uint and sets it to val.
func NumericArrayToUintSlice(val *[]uint) sql.Scanner {
	return numericArrayToUintSlice{val: val}
}

// NumericArrayFromUint8Slice returns a driver.Valuer that produces a PostgreSQL numeric[] from the given Go []uint8.
func NumericArrayFromUint8Slice(val []uint8) driver.Valuer {
	return numericArrayFromUint8Slice{val: val}
}

// NumericArrayToUint8Slice returns an sql.Scanner that converts a PostgreSQL numeric[] into a Go []uint8 and sets it to val.
func NumericArrayToUint8Slice(val *[]uint8) sql.Scanner {
	return numericArrayToUint8Slice{val: val}
}

// NumericArrayFromUint16Slice returns a driver.Valuer that produces a PostgreSQL numeric[] from the given Go []uint16.
func NumericArrayFromUint16Slice(val []uint16) driver.Valuer {
	return numericArrayFromUint16Slice{val: val}
}

// NumericArrayToUint16Slice returns an sql.Scanner that converts a PostgreSQL numeric[] into a Go []uint16 and sets it to val.
func NumericArrayToUint16Slice(val *[]uint16) sql.Scanner {
	return numericArrayToUint16Slice{val: val}
}

// NumericArrayFromUint32Slice returns a driver.Valuer that produces a PostgreSQL numeric[] from the given Go []uint32.
func NumericArrayFromUint32Slice(val []uint32) driver.Valuer {
	return numericArrayFromUint32Slice{val: val}
}

// NumericArrayToUint32Slice returns an sql.Scanner that converts a PostgreSQL numeric[] into a Go []uint32 and sets it to val.
func NumericArrayToUint32Slice(val *[]uint32) sql.Scanner {
	return numericArrayToUint32Slice{val: val}
}

// NumericArrayFromUint64Slice returns a driver.Valuer that produces a PostgreSQL numeric[] from the given Go []uint64.
func NumericArrayFromUint64Slice(val []uint64) driver.Valuer {
	return numericArrayFromUint64Slice{val: val}
}

// NumericArrayToUint64Slice returns an sql.Scanner that converts a PostgreSQL numeric[] into a Go []uint64 and sets it to val.
func NumericArrayToUint64Slice(val *[]uint64) sql.Scanner {
	return numericArrayToUint64Slice{val: val}
}

// NumericArrayFromFloat32Slice returns a driver.Valuer that produces a PostgreSQL numeric[] from the given Go []float32.
func NumericArrayFromFloat32Slice(val []float32) driver.Valuer {
	return numericArrayFromFloat32Slice{val: val}
}

// NumericArrayToFloat32Slice returns an sql.Scanner that converts a PostgreSQL numeric[] into a Go []float32 and sets it to val.
func NumericArrayToFloat32Slice(val *[]float32) sql.Scanner {
	return numericArrayToFloat32Slice{val: val}
}

// NumericArrayFromFloat64Slice returns a driver.Valuer that produces a PostgreSQL numeric[] from the given Go []float64.
func NumericArrayFromFloat64Slice(val []float64) driver.Valuer {
	return numericArrayFromFloat64Slice{val: val}
}

// NumericArrayToFloat64Slice returns an sql.Scanner that converts a PostgreSQL numeric[] into a Go []float64 and sets it to val.
func NumericArrayToFloat64Slice(val *[]float64) sql.Scanner {
	return numericArrayToFloat64Slice{val: val}
}

type numericArrayFromIntSlice struct {
	val []int
}

func (v numericArrayFromIntSlice) Value() (driver.Value, error) {
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

type numericArrayToIntSlice struct {
	val *[]int
}

func (v numericArrayToIntSlice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		ints[i] = int(i64)
	}

	*v.val = ints
	return nil
}

type numericArrayFromInt8Slice struct {
	val []int8
}

func (v numericArrayFromInt8Slice) Value() (driver.Value, error) {
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

type numericArrayToInt8Slice struct {
	val *[]int8
}

func (v numericArrayToInt8Slice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		int8s[i] = int8(i64)
	}

	*v.val = int8s
	return nil
}

type numericArrayFromInt16Slice struct {
	val []int16
}

func (v numericArrayFromInt16Slice) Value() (driver.Value, error) {
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

type numericArrayToInt16Slice struct {
	val *[]int16
}

func (v numericArrayToInt16Slice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		int16s[i] = int16(i64)
	}

	*v.val = int16s
	return nil
}

type numericArrayFromInt32Slice struct {
	val []int32
}

func (v numericArrayFromInt32Slice) Value() (driver.Value, error) {
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

type numericArrayToInt32Slice struct {
	val *[]int32
}

func (v numericArrayToInt32Slice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		int32s[i] = int32(i64)
	}

	*v.val = int32s
	return nil
}

type numericArrayFromInt64Slice struct {
	val []int64
}

func (v numericArrayFromInt64Slice) Value() (driver.Value, error) {
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

type numericArrayToInt64Slice struct {
	val *[]int64
}

func (v numericArrayToInt64Slice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		int64s[i] = i64
	}

	*v.val = int64s
	return nil
}

////////////////////////////////////////////////////////////////////////////////

type numericArrayFromUintSlice struct {
	val []uint
}

func (v numericArrayFromUintSlice) Value() (driver.Value, error) {
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

type numericArrayToUintSlice struct {
	val *[]uint
}

func (v numericArrayToUintSlice) Scan(src interface{}) error {
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
		u64, err := strconv.ParseUint(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		uints[i] = uint(u64)
	}

	*v.val = uints
	return nil
}

type numericArrayFromUint8Slice struct {
	val []uint8
}

func (v numericArrayFromUint8Slice) Value() (driver.Value, error) {
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

type numericArrayToUint8Slice struct {
	val *[]uint8
}

func (v numericArrayToUint8Slice) Scan(src interface{}) error {
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
		u64, err := strconv.ParseUint(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		uint8s[i] = uint8(u64)
	}

	*v.val = uint8s
	return nil
}

type numericArrayFromUint16Slice struct {
	val []uint16
}

func (v numericArrayFromUint16Slice) Value() (driver.Value, error) {
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

type numericArrayToUint16Slice struct {
	val *[]uint16
}

func (v numericArrayToUint16Slice) Scan(src interface{}) error {
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
		u64, err := strconv.ParseUint(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		uint16s[i] = uint16(u64)
	}

	*v.val = uint16s
	return nil
}

type numericArrayFromUint32Slice struct {
	val []uint32
}

func (v numericArrayFromUint32Slice) Value() (driver.Value, error) {
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

type numericArrayToUint32Slice struct {
	val *[]uint32
}

func (v numericArrayToUint32Slice) Scan(src interface{}) error {
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
		u64, err := strconv.ParseUint(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		uint32s[i] = uint32(u64)
	}

	*v.val = uint32s
	return nil
}

type numericArrayFromUint64Slice struct {
	val []uint64
}

func (v numericArrayFromUint64Slice) Value() (driver.Value, error) {
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

type numericArrayToUint64Slice struct {
	val *[]uint64
}

func (v numericArrayToUint64Slice) Scan(src interface{}) error {
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
		u64, err := strconv.ParseUint(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		uint64s[i] = u64
	}

	*v.val = uint64s
	return nil
}

type numericArrayFromFloat32Slice struct {
	val []float32
}

func (v numericArrayFromFloat32Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, f := range v.val {
		out = strconv.AppendFloat(out, float64(f), 'f', -1, 32)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type numericArrayToFloat32Slice struct {
	val *[]float32
}

func (v numericArrayToFloat32Slice) Scan(src interface{}) error {
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
		f64, err := strconv.ParseFloat(string(elems[i]), 32)
		if err != nil {
			return err
		}
		float32s[i] = float32(f64)
	}

	*v.val = float32s
	return nil
}

type numericArrayFromFloat64Slice struct {
	val []float64
}

func (v numericArrayFromFloat64Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, f := range v.val {
		out = strconv.AppendFloat(out, f, 'f', -1, 64)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type numericArrayToFloat64Slice struct {
	val *[]float64
}

func (v numericArrayToFloat64Slice) Scan(src interface{}) error {
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
		f64, err := strconv.ParseFloat(string(elems[i]), 64)
		if err != nil {
			return err
		}
		float64s[i] = f64
	}

	*v.val = float64s
	return nil
}
