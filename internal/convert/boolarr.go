package convert

type BoolArr2BoolSlice struct {
	Ptr *[]bool
}

func (s BoolArr2BoolSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*s.Ptr = nil
		return nil
	}

	elems := pgparsearray1(arr)
	bools := make([]bool, len(elems))
	for i := 0; i < len(elems); i++ {
		if elems[i][0] == 't' {
			bools[i] = true
		} else {
			bools[i] = false
		}
	}

	*s.Ptr = bools
	return nil
}
