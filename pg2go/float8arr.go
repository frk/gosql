package pg2go

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// Float8ArrayFromFloat32Slice returns a driver.Valuer that produces a PostgreSQL float8[] from the given Go []float32.
func Float8ArrayFromFloat32Slice(val []float32) driver.Valuer {
	return float8ArrayFromFloat32Slice{val: val}
}

// Float8ArrayToFloat32Slice returns an sql.Scanner that converts a PostgreSQL float8[] into a Go []float32 and sets it to val.
func Float8ArrayToFloat32Slice(val *[]float32) sql.Scanner {
	return float8ArrayToFloat32Slice{val: val}
}

// Float8ArrayFromFloat64Slice returns a driver.Valuer that produces a PostgreSQL float8[] from the given Go []float64.
func Float8ArrayFromFloat64Slice(val []float64) driver.Valuer {
	return float8ArrayFromFloat64Slice{val: val}
}

// Float8ArrayToFloat64Slice returns an sql.Scanner that converts a PostgreSQL float8[] into a Go []float64 and sets it to val.
func Float8ArrayToFloat64Slice(val *[]float64) sql.Scanner {
	return float8ArrayToFloat64Slice{val: val}
}

type float8ArrayFromFloat32Slice struct {
	val []float32
}

func (v float8ArrayFromFloat32Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for i := 0; i < len(v.val); i++ {
		out = strconv.AppendFloat(out, float64(v.val[i]), 'f', -1, 64)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type float8ArrayToFloat32Slice struct {
	val *[]float32
}

func (v float8ArrayToFloat32Slice) Scan(src interface{}) error {
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

type float8ArrayFromFloat64Slice struct {
	val []float64
}

func (v float8ArrayFromFloat64Slice) Value() (driver.Value, error) {
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

type float8ArrayToFloat64Slice struct {
	val *[]float64
}

func (v float8ArrayToFloat64Slice) Scan(src interface{}) error {
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
