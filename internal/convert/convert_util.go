package convert

import (
	"bytes"
	"time"
)

const dateLayout = "2006-01-02"

func srcbytes(src interface{}) ([]byte, error) {
	switch src := src.(type) {
	case []byte:
		return src, nil
	case string:
		return []byte(src), nil
	case nil:
		return nil, nil
	}
	return nil, nil // TODO error
}

// Expected format: '{S [, ...]}' where S is double quoted string that can contain a comma.
func pgparsearray0(a []byte) (out [][]byte) {
	a = a[1 : len(a)-1] // drop curly braces

	n := 0     // start of next elem
	q := false // in quotes?
	for i := 0; i < len(a); i++ {
		if a[i] == '"' {
			q = !q
		} else if !q && a[i] == ',' {
			out = append(out, a[n:i])
			n = i + 1
		}
	}

	// append the last element
	if len(a) > 0 {
		out = append(out, a[n:])
	}
	return out
}

// Expected format: '{X [, ...]}' where X is anything that doesn't contain a comma.
func pgparsearray1(a []byte) (out [][]byte) {
	a = a[1 : len(a)-1] // drop curly braces

	n := 0 // start of next elem
	for i := 0; i < len(a); i++ {
		if a[i] == ',' {
			out = append(out, a[n:i])
			n = i + 1
		}
	}

	// append the last element
	if len(a) > 0 {
		out = append(out, a[n:])
	}
	return out
}

// parse as single bytes
func pgparsearray2(a []byte) (out []byte) {
	out = make([]byte, 0, len(a))

	a = a[1 : len(a)-1] // drop curly braces

	for i := 0; i < len(a); i++ {
		if a[i] == '"' {
			if a[i+1] == '\\' {
				out = append(out, a[i+2])
				i += 3
				continue
			} else {
				out = append(out, a[i+1])
				i += 2
				continue
			}
		}

		if a[i] != ',' {
			out = append(out, a[i])
		}
	}

	return out
}

// parse as single runes
func pgparsearray3(a []byte) (out []rune) {
	out = make([]rune, 0)

	a = a[1 : len(a)-1] // drop curly braces
	r := bytes.Runes(a)

	for i := 0; i < len(r); i++ {
		if r[i] == '"' {
			if r[i+1] == '\\' {
				out = append(out, r[i+2])
				i += 3
				continue
			} else {
				out = append(out, r[i+1])
				i += 2
				continue
			}
		}

		if r[i] != ',' {
			out = append(out, r[i])
		}
	}

	return out
}

// Expected format: '(X1,Y1),(X2,Y2)' where X and Y are numbers.
func pgparsebox(a []byte) (out [][]byte) {
	a = a[1 : len(a)-1] // drop the first '(' and last ')'

	n := 0 // start of next elem
	for i := 0; i < len(a); i++ {
		if a[i] == ',' {
			out = append(out, a[n:i])
			n = i + 1
		}

		if a[i] == ')' {
			out = append(out, a[n:i])
			i += 2 // skip over ",("
			n = i + 1
		}
	}

	// append the last element
	if n > 0 {
		out = append(out, a[n:])
	}
	return out
}

// Expected format: '{(X1,Y1),(X2,Y2)[; ...]}' where X and Y are numbers.
func pgparseboxarr(a []byte) (out [][]byte) {
	if len(a) == 2 {
		return out // nothing to do if empty "{}"
	}

	a = a[2 : len(a)-2] // drop the first "{(" and last ")}"

	n := 0 // start of next elem
	for i := 0; i < len(a); i++ {
		if a[i] == ',' {
			out = append(out, a[n:i])
			n = i + 1
		}

		if a[i] == ')' {
			out = append(out, a[n:i])
			i += 2 // skip over ",(" or ";("
			n = i + 1
		}
	}

	// append the last element
	if n > 0 {
		out = append(out, a[n:])
	}
	return out
}

func pgparserange(a []byte) (out [2][]byte) {
	a = a[1 : len(a)-1] // drop range delimiters

	for i := 0; i < len(a); i++ {
		if a[i] == ',' {
			out[0] = a[:i]
			out[1] = a[i+1:]
			break
		}
	}

	return out
}

func pgparsedate(a []byte) (time.Time, error) {
	return time.ParseInLocation(dateLayout, string(a), time.UTC)
}
