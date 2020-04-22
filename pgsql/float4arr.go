package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// Float4ArrayFromFloat32Slice returns a driver.Valuer that produces a PostgreSQL float4[] from the given Go []float32.
func Float4ArrayFromFloat32Slice(val []float32) driver.Valuer {
	return float4ArrayFromFloat32Slice{val: val}
}

// Float4ArrayToFloat32Slice returns an sql.Scanner that converts a PostgreSQL float4[] into a Go []float32 and sets it to val.
func Float4ArrayToFloat32Slice(val *[]float32) sql.Scanner {
	return float4ArrayToFloat32Slice{val: val}
}

// Float4ArrayFromFloat64Slice returns a driver.Valuer that produces a PostgreSQL float4[] from the given Go []float64.
func Float4ArrayFromFloat64Slice(val []float64) driver.Valuer {
	return float4ArrayFromFloat64Slice{val: val}
}

// Float4ArrayToFloat64Slice returns an sql.Scanner that converts a PostgreSQL float4[] into a Go []float64 and sets it to val.
func Float4ArrayToFloat64Slice(val *[]float64) sql.Scanner {
	return float4ArrayToFloat64Slice{val: val}
}

type float4ArrayFromFloat32Slice struct {
	val []float32
}

func (v float4ArrayFromFloat32Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for i := 0; i < len(v.val); i++ {
		out = strconv.AppendFloat(out, float64(v.val[i]), 'f', -1, 32)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type float4ArrayToFloat32Slice struct {
	val *[]float32
}

func (v float4ArrayToFloat32Slice) Scan(src interface{}) error {
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
		f32, err := strconv.ParseFloat(string(elems[i]), 32)
		if err != nil {
			return err
		}
		float32s[i] = float32(f32)
	}

	*v.val = float32s
	return nil
}

type float4ArrayFromFloat64Slice struct {
	val []float64
}

func (v float4ArrayFromFloat64Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for i := 0; i < len(v.val); i++ {
		out = strconv.AppendFloat(out, v.val[i], 'f', -1, 64)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type float4ArrayToFloat64Slice struct {
	val *[]float64
}

func (v float4ArrayToFloat64Slice) Scan(src interface{}) error {
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
		float64s[i] = float64(f64)
	}

	*v.val = float64s
	return nil
}
