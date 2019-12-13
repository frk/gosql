package convert

import (
	"strconv"
)

type FloatArr2Float32Slice struct {
	Ptr *[]float32
}

func (s FloatArr2Float32Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*s.Ptr = nil
		return nil
	}

	elems := pgparsearray1(arr)
	float32s := make([]float32, len(elems))
	for i := 0; i < len(elems); i++ {
		f64, err := strconv.ParseFloat(string(elems[i]), 64)
		if err != nil {
			return err
		}
		float32s[i] = float32(f64)
	}

	*s.Ptr = float32s
	return nil
}

type FloatArr2Float64Slice struct {
	Ptr *[]float64
}

func (s FloatArr2Float64Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*s.Ptr = nil
		return nil
	}

	elems := pgparsearray1(arr)
	float64s := make([]float64, len(elems))
	for i := 0; i < len(elems); i++ {
		f64, err := strconv.ParseFloat(string(elems[i]), 64)
		if err != nil {
			return err
		}
		float64s[i] = f64
	}

	*s.Ptr = float64s
	return nil
}

type FloatArr2Int8Slice struct {
	Ptr *[]int8
}

func (s FloatArr2Int8Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*s.Ptr = nil
		return nil
	}

	elems := pgparsearray1(arr)
	int8s := make([]int8, len(elems))
	for i := 0; i < len(elems); i++ {
		f64, err := strconv.ParseFloat(string(elems[i]), 64)
		if err != nil {
			return err
		}
		int8s[i] = int8(f64)
	}

	*s.Ptr = int8s
	return nil
}

type FloatArr2Int16Slice struct {
	Ptr *[]int16
}

func (s FloatArr2Int16Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*s.Ptr = nil
		return nil
	}

	elems := pgparsearray1(arr)
	int16s := make([]int16, len(elems))
	for i := 0; i < len(elems); i++ {
		f64, err := strconv.ParseFloat(string(elems[i]), 64)
		if err != nil {
			return err
		}
		int16s[i] = int16(f64)
	}

	*s.Ptr = int16s
	return nil
}

type FloatArr2Int32Slice struct {
	Ptr *[]int32
}

func (s FloatArr2Int32Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*s.Ptr = nil
		return nil
	}

	elems := pgparsearray1(arr)
	int32s := make([]int32, len(elems))
	for i := 0; i < len(elems); i++ {
		f64, err := strconv.ParseFloat(string(elems[i]), 64)
		if err != nil {
			return err
		}
		int32s[i] = int32(f64)
	}

	*s.Ptr = int32s
	return nil
}

type FloatArr2Int64Slice struct {
	Ptr *[]int64
}

func (s FloatArr2Int64Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*s.Ptr = nil
		return nil
	}

	elems := pgparsearray1(arr)
	int64s := make([]int64, len(elems))
	for i := 0; i < len(elems); i++ {
		f64, err := strconv.ParseFloat(string(elems[i]), 64)
		if err != nil {
			return err
		}
		int64s[i] = int64(f64)
	}

	*s.Ptr = int64s
	return nil
}

type FloatArr2IntSlice struct {
	Ptr *[]int
}

func (s FloatArr2IntSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*s.Ptr = nil
		return nil
	}

	elems := pgparsearray1(arr)
	ints := make([]int, len(elems))
	for i := 0; i < len(elems); i++ {
		f64, err := strconv.ParseFloat(string(elems[i]), 64)
		if err != nil {
			return err
		}
		ints[i] = int(f64)
	}

	*s.Ptr = ints
	return nil
}

type FloatArr2Uint8Slice struct {
	Ptr *[]uint8
}

func (s FloatArr2Uint8Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*s.Ptr = nil
		return nil
	}

	elems := pgparsearray1(arr)
	uint8s := make([]uint8, len(elems))
	for i := 0; i < len(elems); i++ {
		f64, err := strconv.ParseFloat(string(elems[i]), 64)
		if err != nil {
			return err
		}
		uint8s[i] = uint8(f64)
	}

	*s.Ptr = uint8s
	return nil
}

type FloatArr2Uint16Slice struct {
	Ptr *[]uint16
}

func (s FloatArr2Uint16Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*s.Ptr = nil
		return nil
	}

	elems := pgparsearray1(arr)
	uint16s := make([]uint16, len(elems))
	for i := 0; i < len(elems); i++ {
		f64, err := strconv.ParseFloat(string(elems[i]), 64)
		if err != nil {
			return err
		}
		uint16s[i] = uint16(f64)
	}

	*s.Ptr = uint16s
	return nil
}

type FloatArr2Uint32Slice struct {
	Ptr *[]uint32
}

func (s FloatArr2Uint32Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*s.Ptr = nil
		return nil
	}

	elems := pgparsearray1(arr)
	uint32s := make([]uint32, len(elems))
	for i := 0; i < len(elems); i++ {
		f64, err := strconv.ParseFloat(string(elems[i]), 64)
		if err != nil {
			return err
		}
		uint32s[i] = uint32(f64)
	}

	*s.Ptr = uint32s
	return nil
}

type FloatArr2Uint64Slice struct {
	Ptr *[]uint64
}

func (s FloatArr2Uint64Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*s.Ptr = nil
		return nil
	}

	elems := pgparsearray1(arr)
	uint64s := make([]uint64, len(elems))
	for i := 0; i < len(elems); i++ {
		f64, err := strconv.ParseFloat(string(elems[i]), 64)
		if err != nil {
			return err
		}
		uint64s[i] = uint64(f64)
	}

	*s.Ptr = uint64s
	return nil
}

type FloatArr2UintSlice struct {
	Ptr *[]uint
}

func (s FloatArr2UintSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*s.Ptr = nil
		return nil
	}

	elems := pgparsearray1(arr)
	uints := make([]uint, len(elems))
	for i := 0; i < len(elems); i++ {
		f64, err := strconv.ParseFloat(string(elems[i]), 64)
		if err != nil {
			return err
		}
		uints[i] = uint(f64)
	}

	*s.Ptr = uints
	return nil
}
