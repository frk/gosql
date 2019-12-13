package convert

type BitArr2BoolSlice struct {
	Ptr *[]bool
}

func (s BitArr2BoolSlice) Scan(src interface{}) error {
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
		if elems[i][0] == '1' {
			bools[i] = true
		} else {
			bools[i] = false
		}
	}

	*s.Ptr = bools
	return nil
}

type BitArr2Uint8Slice struct {
	Ptr *[]uint8
}

func (s BitArr2Uint8Slice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*s.Ptr = nil
		return nil
	}

	elems := pgparsearray1(arr)
	uint8s := make([]uint8, len(elems))
	for i := 0; i < len(elems); i++ {
		if elems[i][0] == '1' {
			uint8s[i] = 1
		} else {
			uint8s[i] = 0
		}
	}

	*s.Ptr = uint8s
	return nil
}

type BitArr2UintSlice struct {
	Ptr *[]uint
}

func (s BitArr2UintSlice) Scan(src interface{}) error {
	arr, err := srcbytes(src)
	if err != nil {
		return err
	} else if arr == nil {
		*s.Ptr = nil
		return nil
	}

	elems := pgparsearray1(arr)
	uints := make([]uint, len(elems))
	for i := 0; i < len(elems); i++ {
		if elems[i][0] == '1' {
			uints[i] = 1
		} else {
			uints[i] = 0
		}
	}

	*s.Ptr = uints
	return nil
}
