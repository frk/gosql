package pg2go

import (
	"database/sql/driver"
	"unicode/utf8"
)

type CharArrayFromString struct {
	Val string
}

func (v CharArrayFromString) Value() (driver.Value, error) {
	if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := make([]byte, 1, (len(v.Val)*2)+1)
	out[0] = '{'

	for _, r := range v.Val {
		switch r {
		case ' ', '\t', '\r', '\n', '\f', '\v':
			out = append(out, '"', byte(r), '"', ',')
		case '"', '\\':
			out = append(out, '"', '\\', byte(r), '"', ',')
		case ',':
			out = append(out, '"', ',', '"', ',')
		case '\b', '\a':
			// TODO handle if possible
		default:
			if size := utf8.RuneLen(r); size > 1 {
				p := make([]byte, size+1, size+1)
				utf8.EncodeRune(p, r)
				p[size] = ','

				out = append(out, p...)
			} else {
				out = append(out, byte(r), ',')
			}
		}
	}

	out[len(out)-1] = '}'
	return out, nil
}

type CharArrayFromByteSlice struct {
	Val []byte
}

func (v CharArrayFromByteSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := make([]byte, 1, (len(v.Val)*2)+1)
	out[0] = '{'

	for _, b := range v.Val {
		switch b {
		case ' ', '\t', '\r', '\n', '\f', '\v':
			out = append(out, '"', b, '"', ',')
		case '"', '\\':
			out = append(out, '"', '\\', b, '"', ',')
		case ',':
			out = append(out, '"', ',', '"', ',')
		case '\b', '\a':
			// TODO handle if possible
		default:
			out = append(out, b, ',')
		}
	}

	out[len(out)-1] = '}'
	return out, nil
}

type CharArrayFromRuneSlice struct {
	Val []rune
}

func (v CharArrayFromRuneSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := make([]byte, 1, (len(v.Val)*2)+1)
	out[0] = '{'

	for _, r := range v.Val {
		switch r {
		case ' ', '\t', '\r', '\n', '\f', '\v':
			out = append(out, '"', byte(r), '"', ',')
		case '"', '\\':
			out = append(out, '"', '\\', byte(r), '"', ',')
		case ',':
			out = append(out, '"', ',', '"', ',')
		case '\b', '\a':
			// TODO handle if possible
		default:
			if size := utf8.RuneLen(r); size > 1 {
				p := make([]byte, size+1, size+1)
				utf8.EncodeRune(p, r)
				p[size] = ','

				out = append(out, p...)
			} else {
				out = append(out, byte(r), ',')
			}
		}
	}

	out[len(out)-1] = '}'
	return out, nil
}

type CharArrayFromStringSlice struct {
	Val []string
}

func (v CharArrayFromStringSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := make([]byte, 1, (len(v.Val)*2)+1)
	out[0] = '{'

	for _, s := range v.Val {
		r, size := utf8.DecodeRuneInString(s)
		switch r {
		case ' ', '\t', '\r', '\n', '\f', '\v':
			out = append(out, '"', byte(r), '"', ',')
		case '"', '\\':
			out = append(out, '"', '\\', byte(r), '"', ',')
		case ',':
			out = append(out, '"', ',', '"', ',')
		case '\b', '\a':
			// TODO handle if possible
		default:
			if size > 1 {
				p := make([]byte, size+1, size+1)
				utf8.EncodeRune(p, r)
				p[size] = ','

				out = append(out, p...)
			} else {
				out = append(out, byte(r), ',')
			}
		}
	}

	out[len(out)-1] = '}'
	return out, nil
}

type CharArrayToString struct {
	Val *string
}

func (s CharArrayToString) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	bytes := pgparsearray2(data)
	*s.Val = string(bytes)
	return nil
}

type CharArrayToByteSlice struct {
	Val *[]byte
}

func (s CharArrayToByteSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*s.Val = nil
		return nil
	}

	bytes := pgparsearray2(data)
	*s.Val = bytes
	return nil
}

type CharArrayToRuneSlice struct {
	Val *[]rune
}

func (s CharArrayToRuneSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*s.Val = nil
		return nil
	}

	runes := pgparsearray3(data)
	*s.Val = runes
	return nil
}

type CharArrayToStringSlice struct {
	Val *[]string
}

func (s CharArrayToStringSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*s.Val = nil
		return nil
	}

	runes := pgparsearray3(data)
	Val := make([]string, len(runes))
	for i := 0; i < len(Val); i++ {
		Val[i] = string(runes[i])
	}
	*s.Val = Val
	return nil
}
