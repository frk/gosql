package convert

//import "fmt"

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
