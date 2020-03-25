package convert

import (
	"database/sql/driver"
	"encoding/hex"
)

type ByteaArrayFromByteSliceSlice struct {
	S [][]byte
}

func (v ByteaArrayFromByteSliceSlice) Value() (driver.Value, error) {
	if v.S == nil {
		return nil, nil
	} else if len(v.S) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := make([]byte, 1, (len(v.S)*3)+1)
	out[0] = '{'

	for _, src := range v.S {
		out = append(out, '"', '\\', '\\', 'x')

		dst := make([]byte, hex.EncodedLen(len(src)))
		_ = hex.Encode(dst, src)
		out = append(out, dst...)

		out = append(out, '"', ',')
	}

	out[len(out)-1] = '}'
	return out, nil
}

type ByteaArrayFromStringSlice struct {
	S []string
}

func (v ByteaArrayFromStringSlice) Value() (driver.Value, error) {
	if v.S == nil {
		return nil, nil
	} else if len(v.S) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := make([]byte, 1, (len(v.S)*3)+1)
	out[0] = '{'

	for _, s := range v.S {
		out = append(out, '"', '\\', '\\', 'x')

		src := []byte(s)
		dst := make([]byte, hex.EncodedLen(len(src)))
		_ = hex.Encode(dst, src)
		out = append(out, dst...)

		out = append(out, '"', ',')
	}

	out[len(out)-1] = '}'
	return out, nil
}

type ByteaArrayToByteSliceSlice struct {
	S *[][]byte
}

func (s ByteaArrayToByteSliceSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*s.S = nil
		return nil
	}

	elems := pgparsearray1(arr)
	out := make([][]byte, len(elems))
	for i := 0; i < len(elems); i++ {
		src := elems[i]

		// drop the initial "\\x and the last "
		src = src[4 : len(src)-1]

		dst := make([]byte, hex.DecodedLen(len(src)))
		if _, err := hex.Decode(dst, src); err != nil {
			return err
		}

		out[i] = dst
	}

	*s.S = out
	return nil
}

type ByteaArrayToStringSlice struct {
	S *[]string
}

func (s ByteaArrayToStringSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*s.S = nil
		return nil
	}

	elems := pgparsearray1(arr)
	out := make([]string, len(elems))
	for i := 0; i < len(elems); i++ {
		src := elems[i]

		// drop the initial "\\x and the last "
		src = src[4 : len(src)-1]

		dst := make([]byte, hex.DecodedLen(len(src)))
		if _, err := hex.Decode(dst, src); err != nil {
			return err
		}

		out[i] = string(dst)
	}

	*s.S = out
	return nil
}
