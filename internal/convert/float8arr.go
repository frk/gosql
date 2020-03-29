package convert

import (
	"database/sql/driver"
	"strconv"
)

type Float8ArrayFromFloat32Slice struct {
	V []float32
}

func (v Float8ArrayFromFloat32Slice) Value() (driver.Value, error) {
	if v.V == nil {
		return nil, nil
	} else if len(v.V) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, f := range v.V {
		out = strconv.AppendFloat(out, float64(f), 'f', -1, 64)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Float8ArrayToFloat32Slice struct {
	V *[]float32
}

func (v Float8ArrayToFloat32Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.V = nil
		return nil
	}

	elems := pgparsearray1(arr)
	float32s := make([]float32, len(elems))
	for i := 0; i < len(elems); i++ {
		f32, err := strconv.ParseFloat(string(elems[i]), 32)
		if err != nil {
			return err
		}
		float32s[i] = float32(f32)
	}

	*v.V = float32s
	return nil
}

type Float8ArrayFromFloat64Slice struct {
	V []float64
}

func (v Float8ArrayFromFloat64Slice) Value() (driver.Value, error) {
	if v.V == nil {
		return nil, nil
	} else if len(v.V) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, f := range v.V {
		out = strconv.AppendFloat(out, f, 'f', -1, 64)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type Float8ArrayToFloat64Slice struct {
	V *[]float64
}

func (v Float8ArrayToFloat64Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*v.V = nil
		return nil
	}

	elems := pgparsearray1(arr)
	float64s := make([]float64, len(elems))
	for i := 0; i < len(elems); i++ {
		f64, err := strconv.ParseFloat(string(elems[i]), 64)
		if err != nil {
			return err
		}
		float64s[i] = float64(f64)
	}

	*v.V = float64s
	return nil
}
