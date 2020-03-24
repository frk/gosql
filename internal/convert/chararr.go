package convert

import (
	"database/sql/driver"
	"unicode/utf8"
)

type CharArrayFromString struct {
	S string
}

func (c CharArrayFromString) Value() (driver.Value, error) {
	if len(c.S) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := make([]byte, 1, (len(c.S)*2)+1)
	out[0] = '{'

	for _, r := range c.S {
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
	S []byte
}

func (c CharArrayFromByteSlice) Value() (driver.Value, error) {
	if c.S == nil {
		return nil, nil
	} else if len(c.S) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := make([]byte, 1, (len(c.S)*2)+1)
	out[0] = '{'

	for _, b := range c.S {
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
	S []rune
}

func (c CharArrayFromRuneSlice) Value() (driver.Value, error) {
	if c.S == nil {
		return nil, nil
	} else if len(c.S) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := make([]byte, 1, (len(c.S)*2)+1)
	out[0] = '{'

	for _, r := range c.S {
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
	S []string
}

func (c CharArrayFromStringSlice) Value() (driver.Value, error) {
	if c.S == nil {
		return nil, nil
	} else if len(c.S) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := make([]byte, 1, (len(c.S)*2)+1)
	out[0] = '{'

	for _, s := range c.S {
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
	S *string
}

func (s CharArrayToString) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	bytes := pgparsearray2(data)
	*s.S = string(bytes)
	return nil
}

type CharArrayToByteSlice struct {
	S *[]byte
}

func (s CharArrayToByteSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*s.S = nil
		return nil
	}

	bytes := pgparsearray2(data)
	*s.S = bytes
	return nil
}

type CharArrayToRuneSlice struct {
	S *[]rune
}

func (s CharArrayToRuneSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*s.S = nil
		return nil
	}

	runes := pgparsearray3(data)
	*s.S = runes
	return nil
}

type CharArrayToStringSlice struct {
	S *[]string
}

func (s CharArrayToStringSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*s.S = nil
		return nil
	}

	runes := pgparsearray3(data)
	S := make([]string, len(runes))
	for i := 0; i < len(S); i++ {
		S[i] = string(runes[i])
	}
	*s.S = S
	return nil
}
