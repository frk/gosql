package convert

import (
	"database/sql/driver"
	"strconv"
)

type IntArr2Int8Slice struct {
	Ptr *[]int8
}

func (s IntArr2Int8Slice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		int8s[i] = int8(i64)
	}

	*s.Ptr = int8s
	return nil
}

type IntArr2Int16Slice struct {
	Ptr *[]int16
}

func (s IntArr2Int16Slice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		int16s[i] = int16(i64)
	}

	*s.Ptr = int16s
	return nil
}

type IntArr2Int32Slice struct {
	Ptr *[]int32
}

func (s IntArr2Int32Slice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		int32s[i] = int32(i64)
	}

	*s.Ptr = int32s
	return nil
}

type IntArr2Int64Slice struct {
	Ptr *[]int64
}

func (s IntArr2Int64Slice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		int64s[i] = i64
	}

	*s.Ptr = int64s
	return nil
}

type IntArr2IntSlice struct {
	Ptr *[]int
}

func (s IntArr2IntSlice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		ints[i] = int(i64)
	}

	*s.Ptr = ints
	return nil
}

type IntArr2Uint8Slice struct {
	Ptr *[]uint8
}

func (s IntArr2Uint8Slice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		uint8s[i] = uint8(i64)
	}

	*s.Ptr = uint8s
	return nil
}

type IntArr2Uint16Slice struct {
	Ptr *[]uint16
}

func (s IntArr2Uint16Slice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		uint16s[i] = uint16(i64)
	}

	*s.Ptr = uint16s
	return nil
}

type IntArr2Uint32Slice struct {
	Ptr *[]uint32
}

func (s IntArr2Uint32Slice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		uint32s[i] = uint32(i64)
	}

	*s.Ptr = uint32s
	return nil
}

type IntArr2Uint64Slice struct {
	Ptr *[]uint64
}

func (s IntArr2Uint64Slice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		uint64s[i] = uint64(i64)
	}

	*s.Ptr = uint64s
	return nil
}

type IntArr2UintSlice struct {
	Ptr *[]uint
}

func (s IntArr2UintSlice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 10, 64)
		if err != nil {
			return err
		}
		uints[i] = uint(i64)
	}

	*s.Ptr = uints
	return nil
}

type IntSlice2IntArray struct {
	S []int
}

func (s IntSlice2IntArray) Value() (driver.Value, error) {
	if s.S == nil {
		return nil, nil
	}

	if n := len(s.S); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '{'

		b = strconv.AppendInt(b, int64(s.S[0]), 10)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendInt(b, int64(s.S[i]), 10)
		}

		return append(b, '}'), nil
	}

	return []byte{'{', '}'}, nil
}
