package pg2go

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
)

// VarBitFromInt64 returns a driver.Valuer that produces a PostgreSQL varbit from the given Go int64.
func VarBitFromInt64(val int64) driver.Valuer {
	return varBitFromInt64{val: val}
}

// VarBitToInt64 returns an sql.Scanner that converts a PostgreSQL varbit into a Go int64 and sets it to val.
func VarBitToInt64(val *int64) sql.Scanner {
	return varBitToInt64{val: val}
}

// VarBitFromBoolSlice returns a driver.Valuer that produces a PostgreSQL varbit from the given Go []bool.
func VarBitFromBoolSlice(val []bool) driver.Valuer {
	return varBitFromBoolSlice{val: val}
}

// VarBitToBoolSlice returns an sql.Scanner that converts a PostgreSQL varbit into a Go []bool and sets it to val.
func VarBitToBoolSlice(val *[]bool) sql.Scanner {
	return varBitToBoolSlice{val: val}
}

// VarBitFromUint8Slice returns a driver.Valuer that produces a PostgreSQL varbit from the given Go []uint8.
func VarBitFromUint8Slice(val []uint8) driver.Valuer {
	return varBitFromUint8Slice{val: val}
}

// VarBitToUint8Slice returns an sql.Scanner that converts a PostgreSQL varbit into a Go []uint8 and sets it to val.
func VarBitToUint8Slice(val *[]uint8) sql.Scanner {
	return varBitToUint8Slice{val: val}
}

type varBitFromInt64 struct {
	val int64
}

func (v varBitFromInt64) Value() (driver.Value, error) {
	out := strconv.AppendInt([]byte(nil), v.val, 2)
	return out, nil
}

type varBitToInt64 struct {
	val *int64
}

func (v varBitToInt64) Scan(src interface{}) error {
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

	*v.val = i64
	return nil
}

type varBitFromBoolSlice struct {
	val []bool
}

func (v varBitFromBoolSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte(""), nil
	}

	out := make([]byte, len(v.val))
	for i := 0; i < len(v.val); i++ {
		if v.val[i] {
			out[i] = '1'
		} else {
			out[i] = '0'
		}
	}
	return out, nil
}

type varBitToBoolSlice struct {
	val *[]bool
}

func (v varBitToBoolSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
		return nil
	}

	bools := make([]bool, len(data))
	for i := 0; i < len(data); i++ {
		if data[i] == '1' {
			bools[i] = true
		}
	}

	*v.val = bools
	return nil
}

type varBitFromUint8Slice struct {
	val []uint8
}

func (v varBitFromUint8Slice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte(""), nil
	}

	out := make([]byte, len(v.val))
	for i := 0; i < len(v.val); i++ {
		if v.val[i] == 1 {
			out[i] = '1'
		} else {
			out[i] = '0'
		}
	}
	return out, nil
}

type varBitToUint8Slice struct {
	val *[]uint8
}

func (v varBitToUint8Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
		return nil
	}

	uint8s := make([]uint8, len(data))
	for i := 0; i < len(data); i++ {
		if data[i] == '1' {
			uint8s[i] = 1
		}
	}

	*v.val = uint8s
	return nil
}
