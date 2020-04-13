package convert

import (
	"database/sql/driver"
	"strconv"
)

type VarBitFromInt64 struct {
	Val int64
}

func (v VarBitFromInt64) Value() (driver.Value, error) {
	out := strconv.AppendInt([]byte(nil), v.Val, 2)
	return out, nil
}

type VarBitToInt64 struct {
	Val *int64
}

func (v VarBitToInt64) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	i64, err := strconv.ParseInt(string(data), 2, 64)
	if err != nil {
		return err
	}

	*v.Val = i64
	return nil
}

type VarBitFromBoolSlice struct {
	Val []bool
}

func (v VarBitFromBoolSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte(""), nil
	}

	out := make([]byte, len(v.Val))
	for i := 0; i < len(v.Val); i++ {
		if v.Val[i] {
			out[i] = '1'
		} else {
			out[i] = '0'
		}
	}
	return out, nil
}

type VarBitToBoolSlice struct {
	Val *[]bool
}

func (v VarBitToBoolSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.Val = nil
		return nil
	}

	bools := make([]bool, len(data))
	for i := 0; i < len(data); i++ {
		if data[i] == '1' {
			bools[i] = true
		}
	}

	*v.Val = bools
	return nil
}

type VarBitFromUint8Slice struct {
	Val []uint8
}

func (v VarBitFromUint8Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte(""), nil
	}

	out := make([]byte, len(v.Val))
	for i := 0; i < len(v.Val); i++ {
		if v.Val[i] == 1 {
			out[i] = '1'
		} else {
			out[i] = '0'
		}
	}
	return out, nil
}

type VarBitToUint8Slice struct {
	Val *[]uint8
}

func (v VarBitToUint8Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.Val = nil
		return nil
	}

	uint8s := make([]uint8, len(data))
	for i := 0; i < len(data); i++ {
		if data[i] == '1' {
			uint8s[i] = 1
		}
	}

	*v.Val = uint8s
	return nil
}
