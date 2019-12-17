package convert

// []byte
// []rune
// string

import (
//"bytes"
//"fmt"
)

type CharArr2ByteSlice struct {
	Ptr *[]byte
}

func (s CharArr2ByteSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*s.Ptr = nil
		return nil
	}

	bytes := pgparsearray2(data)
	*s.Ptr = bytes
	return nil
}

type CharArr2RuneSlice struct {
	Ptr *[]rune
}

func (s CharArr2RuneSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*s.Ptr = nil
		return nil
	}

	runes := pgparsearray3(data)
	*s.Ptr = runes
	return nil
}

type CharArr2String struct {
	Ptr *string
}

func (s CharArr2String) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	bytes := pgparsearray2(data)
	*s.Ptr = string(bytes)
	return nil
}
