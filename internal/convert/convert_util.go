package convert

import (
	"bytes"
	"encoding/hex"
	"time"
)

const (
	dateLayout        = "2006-01-02"
	timeLayout        = "15:04:05.999"
	timetzLayout      = "15:04:05.999-07:00"
	timestampLayout   = "2006-01-02 15:04:05.999"
	timestamptzLayout = "2006-01-02 15:04:05.999-07"
)

var noZone = time.FixedZone("", 0)

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

func pgparsehstore(a []byte) (out [][2][]byte) {
	if len(a) == 0 {
		return out // nothing to do if empty
	}

	var (
		idx  uint8     // pair index 0=key 1=val
		pair [2][]byte // current pair
	)

	for i := 0; i < len(a); i++ {
		if a[i] == '"' { // start of key or value?
			for j := i + 1; j < len(a); j++ {
				if a[j] == '\\' {
					j++ // skip escaped char
					pair[idx] = append(pair[idx], a[j])
					continue
				}

				if a[j] == '"' { // end of key or value?
					i = j
					break
				}

				pair[idx] = append(pair[idx], a[j])
			}

			if idx == 1 { // is the pair done?
				out = append(out, pair)
				pair = [2][]byte{} // next pair
			}

			// flip the index
			idx ^= 1
		}

		if a[i] == 'N' {
			i += 4
			idx = 0
			out = append(out, pair)
			pair = [2][]byte{}
		}
	}
	return out
}

func pgparsehstorearr(a []byte) (out [][][2][]byte) {
	if len(a) == 0 {
		return out // nothing to do if empty
	}
	a = a[1 : len(a)-1] // drop array delimiters

	var (
		idx  uint8     // pair index 0=key 1=val
		pair [2][]byte // current pair
	)

	for i := 0; i < len(a); i++ {
		if a[i] == '"' { // start of hstore in array?
			hstore := [][2][]byte{}

			for j := i + 1; j < len(a); j++ {

				if a[j] == '\\' && a[j+1] == '"' { // start of key or value?
					for k := j + 2; k < len(a); k++ {
						if a[k] == '\\' {
							if a[k+1] == '\\' { // escape??
								k += 3 // skip escaped char
								pair[idx] = append(pair[idx], a[k])
								continue
							}

							if a[k+1] == '"' { // key or value ending quote?
								j = k + 2
								break
							}
						}

						pair[idx] = append(pair[idx], a[k])
					}

					if idx == 1 { // is the pair done?
						hstore = append(hstore, pair)
						pair = [2][]byte{} // next pair
					}

					// flip the index
					idx ^= 1
				}

				// handle NULL value in pair
				if a[j] == 'N' {
					j += 4
					idx = 0
					hstore = append(hstore, pair)
					pair = [2][]byte{} // next pair
				}

				if a[j] == '"' { // end of hstore in array?
					out = append(out, hstore)
					i = j
					break
				}
			}
		}

		// handle NULL hstore
		if a[i] == 'N' {
			out = append(out, nil)
			i += 4
		}
	}
	return out
}

func pgparsedate(a []byte) (time.Time, error) {
	return time.ParseInLocation(dateLayout, string(a), time.UTC)
}

func pgParseRange(a []byte) (out [2][]byte) {
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

// Expected format: '{STRING [, ...]}' where STRING is a double quoted string.
func pgParseQuotedStringArray(a []byte) (out [][]byte) {
	a = a[1 : len(a)-1] // drop curly braces

mainloop:
	for i := 0; i < len(a); i++ {
		if a[i] == '"' { // start of string
			for j := i + 1; j < len(a); j++ {
				if a[j] == '\\' {
					j++ // skip escaped char
					continue
				}

				if a[j] == '"' { // end of string
					out = append(out, a[i+1:j])
					i = j + 1
					continue mainloop
				}
			}
		}

		if a[i] == 'N' { // NULL?
			out = append(out, []byte(nil))
			i += 4
		}
	}

	return out
}

// Expected format: '{STRING [, ...]}' where STRING is either an unquoted or a double quoted string.
func pgParseStringArray(a []byte) (out [][]byte) {
	a = a[1 : len(a)-1] // drop curly braces

mainloop:
	for i := 0; i < len(a); i++ {
		switch a[i] {
		case ',': // element separator?
			continue mainloop
		case '"': // start of double quoted string
			str := []byte{}
			for j := i + 1; j < len(a); j++ {
				if a[j] == '\\' { // escape sequence
					str = append(str, a[i+1:j]...)
					i = j
					j++
					continue
				}

				if a[j] == '"' { // end of string
					str = append(str, a[i+1:j]...)
					out = append(out, str)
					i = j
					continue mainloop
				}
			}
		case 'N': // NULL?
			if len(a) > i+3 && string(a[i:i+4]) == `NULL` {
				out = append(out, []byte(nil))
				i = i + 3
				continue mainloop
			}
			fallthrough
		default: // start of unquoted string
			var j int
			for j = i + 1; j < len(a); j++ {
				if a[j] == ',' { // end of string
					out = append(out, a[i:j])
					i = j
					continue mainloop
				}
			}

			// end of the last element
			out = append(out, a[i:j])
			i = j
		}
	}

	return out
}

// Expected format: '(X1,Y1),(X2,Y2)' where X and Y are numbers.
func pgParseBox(a []byte) (out [][]byte) {
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
func pgParseBoxArray(a []byte) (out [][]byte) {
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

// Expected format: '{X [, ...]}' where X is anything that doesn't contain a comma.
func pgParseCommaArray(a []byte) (out [][]byte) {
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

// Expected format: 'E[ ...]' (space separated list of elements).
func pgParseVector(a []byte) (out [][]byte) {
	var j int
	for i := 0; i < len(a); i++ {
		if a[i] == ' ' {
			out = append(out, a[j:i])
			j = i + 1
		}
	}
	if len(a) > 0 {
		out = append(out, a[j:]) // last
	}
	return out
}

// Expected format: '{VEC [, ...]}' where VEC is a space separated list of elements,
// the list being enclosed in double quotes, or a single unquoted element, or NULL.
func pgParseVectorArray(a []byte) (out [][][]byte) {
	if len(a) == 0 {
		return out // nothing to do if empty
	}
	a = a[1 : len(a)-1] // drop array delimiters

	for i := 0; i < len(a); i++ {
		switch a[i] {
		case '"': // quoted element
			vector := [][]byte{}

			k := i + 1 // vector start
			for j := k; j < len(a); j++ {
				if a[j] == '"' { // vector end
					vector = append(vector, a[k:j])
					i = j + 1
					break
				}
				if a[j] == ' ' {
					vector = append(vector, a[k:j])
					k = j + 1
				}
			}

			out = append(out, vector)
		case 'N': // NULL element?
			if len(a) > i+3 && string(a[i:i+4]) == `NULL` {
				out = append(out, [][]byte(nil))
				i = i + 3
			}
		default: // unquoted element
			k, j := i, i // elem start
			for ; j < len(a); j++ {
				if a[j] == ',' { // elem end
					i = j
					break
				}
			}
			out = append(out, [][]byte{a[k:j]})
		}
	}

	return out
}

// Expected format: 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx'.
func pgParseUUID(data []byte) (arr [16]byte, err error) {
	if _, err = hex.Decode(arr[0:4], data[0:8]); err != nil {
		return arr, err
	}
	if _, err = hex.Decode(arr[4:6], data[9:13]); err != nil {
		return arr, err
	}
	if _, err = hex.Decode(arr[6:8], data[14:18]); err != nil {
		return arr, err
	}
	if _, err = hex.Decode(arr[8:10], data[19:23]); err != nil {
		return arr, err
	}
	if _, err = hex.Decode(arr[10:], data[24:]); err != nil {
		return arr, err
	}
	return arr, nil
}

// pgFormatUUID converts the given array to a slice of bytes in the following
// format "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" representing a uuid.
func pgFormatUUID(arr [16]byte) (out []byte) {
	out = make([]byte, 36)

	_ = hex.Encode(out[0:8], arr[0:4])
	_ = hex.Encode(out[9:13], arr[4:6])
	_ = hex.Encode(out[14:18], arr[6:8])
	_ = hex.Encode(out[19:23], arr[8:10])
	_ = hex.Encode(out[24:], arr[10:])

	out[8] = '-'
	out[13] = '-'
	out[18] = '-'
	out[23] = '-'

	return out
}

// pgAppendQuote1
func pgAppendQuote1(buf, elem []byte) []byte {
	buf = append(buf, '"')
	for i := 0; i < len(elem); i++ {
		switch elem[i] {
		case '"', '\\':
			buf = append(buf, '\\', elem[i])
		case '\a':
			buf = append(buf, '\\', '\a')
		case '\b':
			buf = append(buf, '\\', '\b')
		case '\f':
			buf = append(buf, '\\', '\f')
		case '\n':
			buf = append(buf, '\\', '\n')
		case '\r':
			buf = append(buf, '\\', '\r')
		case '\t':
			buf = append(buf, '\\', '\t')
		case '\v':
			buf = append(buf, '\\', '\v')
		default:
			buf = append(buf, elem[i])
		}
	}
	return append(buf, '"')
}

// pgAppendQuote2
func pgAppendQuote2(buf, elem []byte) []byte {
	buf = append(buf, '\\', '"')
	for i := 0; i < len(elem); i++ {
		switch elem[i] {
		case '"', '\\':
			buf = append(buf, '\\', '\\', '\\', elem[i])
		case '\a':
			buf = append(buf, '\\', '\\', '\\', '\a')
		case '\b':
			buf = append(buf, '\\', '\\', '\\', '\b')
		case '\f':
			buf = append(buf, '\\', '\\', '\\', '\f')
		case '\n':
			buf = append(buf, '\\', '\\', '\\', '\n')
		case '\r':
			buf = append(buf, '\\', '\\', '\\', '\r')
		case '\t':
			buf = append(buf, '\\', '\\', '\\', '\t')
		case '\v':
			buf = append(buf, '\\', '\\', '\\', '\v')
		default:
			buf = append(buf, elem[i])
		}
	}
	return append(buf, '\\', '"')
}
