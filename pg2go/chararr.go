package pg2go

import (
	"database/sql"
	"database/sql/driver"
	"unicode/utf8"
)

// CharArrayFromString returns a driver.Valuer that produces a PostgreSQL char[] from the given Go string.
func CharArrayFromString(val string) driver.Valuer {
	return charArrayFromString{val: val}
}

// CharArrayToString returns an sql.Scanner that converts a PostgreSQL char[] into a Go string and sets it to val.
func CharArrayToString(val *string) sql.Scanner {
	return charArrayToString{val: val}
}

// CharArrayFromByteSlice returns a driver.Valuer that produces a PostgreSQL char[] from the given Go []byte.
func CharArrayFromByteSlice(val []byte) driver.Valuer {
	return charArrayFromByteSlice{val: val}
}

// CharArrayToByteSlice returns an sql.Scanner that converts a PostgreSQL char[] into a Go []byte and sets it to val.
func CharArrayToByteSlice(val *[]byte) sql.Scanner {
	return charArrayToByteSlice{val: val}
}

// CharArrayFromRuneSlice returns a driver.Valuer that produces a PostgreSQL char[] from the given Go []rune.
func CharArrayFromRuneSlice(val []rune) driver.Valuer {
	return charArrayFromRuneSlice{val: val}
}

// CharArrayToRuneSlice returns an sql.Scanner that converts a PostgreSQL char[] into a Go []rune and sets it to val.
func CharArrayToRuneSlice(val *[]rune) sql.Scanner {
	return charArrayToRuneSlice{val: val}
}

// CharArrayFromStringSlice returns a driver.Valuer that produces a PostgreSQL char[] from the given Go []string.
func CharArrayFromStringSlice(val []string) driver.Valuer {
	return charArrayFromStringSlice{val: val}
}

// CharArrayToStringSlice returns an sql.Scanner that converts a PostgreSQL char[] into a Go []string and sets it to val.
func CharArrayToStringSlice(val *[]string) sql.Scanner {
	return charArrayToStringSlice{val: val}
}

type charArrayFromString struct {
	val string
}

func (v charArrayFromString) Value() (driver.Value, error) {
	if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := make([]byte, 1, (len(v.val)*2)+1)
	out[0] = '{'

	for _, r := range v.val {
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

type charArrayToString struct {
	val *string
}

func (v charArrayToString) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	bytes := pgparsearray2(data)
	*v.val = string(bytes)
	return nil
}

type charArrayFromByteSlice struct {
	val []byte
}

func (v charArrayFromByteSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := make([]byte, 1, (len(v.val)*2)+1)
	out[0] = '{'

	for _, b := range v.val {
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

type charArrayToByteSlice struct {
	val *[]byte
}

func (v charArrayToByteSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
		return nil
	}

	bytes := pgparsearray2(data)
	*v.val = bytes
	return nil
}

type charArrayFromRuneSlice struct {
	val []rune
}

func (v charArrayFromRuneSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := make([]byte, 1, (len(v.val)*2)+1)
	out[0] = '{'

	for _, r := range v.val {
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

type charArrayToRuneSlice struct {
	val *[]rune
}

func (v charArrayToRuneSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
		return nil
	}

	runes := pgparsearray3(data)
	*v.val = runes
	return nil
}

type charArrayFromStringSlice struct {
	val []string
}

func (v charArrayFromStringSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := make([]byte, 1, (len(v.val)*2)+1)
	out[0] = '{'

	for _, v := range v.val {
		r, size := utf8.DecodeRuneInString(v)
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

type charArrayToStringSlice struct {
	val *[]string
}

func (v charArrayToStringSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
		return nil
	}

	runes := pgparsearray3(data)
	val := make([]string, len(runes))
	for i := 0; i < len(val); i++ {
		val[i] = string(runes[i])
	}
	*v.val = val
	return nil
}
