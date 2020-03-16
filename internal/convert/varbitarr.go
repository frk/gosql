package convert

import (
	"strconv"
)

type VarBitArr2StringSlice struct {
	Ptr *[]string
}

func (s VarBitArr2StringSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*s.Ptr = nil
		return nil
	}

	elems := pgparsearray1(arr)
	strings := make([]string, len(elems))
	for i := 0; i < len(elems); i++ {
		strings[i] = string(elems[i])
	}

	*s.Ptr = strings
	return nil
}

type VarBitArr2Int64Slice struct {
	Ptr *[]int64
}

func (s VarBitArr2Int64Slice) Scan(src interface{}) error {
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
		i64, err := strconv.ParseInt(string(elems[i]), 2, 64)
		if err != nil {
			return err
		}
		int64s[i] = i64
	}

	*s.Ptr = int64s
	return nil
}