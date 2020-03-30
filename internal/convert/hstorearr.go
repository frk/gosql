package convert

import (
	"database/sql"
	"database/sql/driver"
)

type HStoreArrayFromStringMapSlice struct {
	Val []map[string]string
}

func (v HStoreArrayFromStringMapSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, m := range v.Val {
		if m == nil {
			out = append(out, 'N', 'U', 'L', 'L', ',')
			continue
		}

		i, size := 0, 0
		pairs := make([][]byte, len(m))
		for key, val := range m {
			// len(`\"\"=>\"\"`) == 10
			pair := make([]byte, 0, 10+len(key)+len(val))
			pair = pgAppendQuote2(pair, []byte(key))
			pair = append(pair, '=', '>')
			pair = pgAppendQuote2(pair, []byte(val))

			pairs[i] = pair
			i += 1
			size += len(pair) + 1
		}

		var hstore []byte
		if size == 0 {
			hstore = []byte{'"', '"'}
		} else {
			hstore = make([]byte, 1, size+1)
			hstore[0] = '"'
			for _, pair := range pairs {
				hstore = append(hstore, pair...)
				hstore = append(hstore, ',')
			}
			hstore[len(hstore)-1] = '"'
		}

		out = append(out, hstore...)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type HStoreArrayToStringMapSlice struct {
	Val *[]map[string]string
}

func (v HStoreArrayToStringMapSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	elems := pgparsehstorearr(data)
	hashes := make([]map[string]string, len(elems))

	for i := 0; i < len(elems); i++ {
		pairs := elems[i]
		if pairs == nil {
			continue
		}

		hash := make(map[string]string)
		for j := 0; j < len(pairs); j++ {
			if value := pairs[j][1]; value != nil {
				hash[string(pairs[j][0])] = string(value)
			} else {
				hash[string(pairs[j][0])] = ""
			}
		}

		hashes[i] = hash
	}

	*v.Val = hashes
	return nil
}

type HStoreArrayFromStringPtrMapSlice struct {
	Val []map[string]*string
}

func (v HStoreArrayFromStringPtrMapSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, m := range v.Val {
		if m == nil {
			out = append(out, 'N', 'U', 'L', 'L', ',')
			continue
		}

		i, size := 0, 0
		pairs := make([][]byte, len(m))
		for key, val := range m {
			if val == nil {
				// len(`\"\"=>NULL`) == 10
				pair := make([]byte, 0, 10+len(key))
				pair = pgAppendQuote2(pair, []byte(key))
				pair = append(pair, '=', '>', 'N', 'U', 'L', 'L')

				pairs[i] = pair
				i += 1
				size += len(pair) + 1
				continue
			}

			// len(`\"\"=>\"\"`) == 10
			pair := make([]byte, 0, 10+len(key)+len(*val))
			pair = pgAppendQuote2(pair, []byte(key))
			pair = append(pair, '=', '>')
			pair = pgAppendQuote2(pair, []byte(*val))

			pairs[i] = pair
			i += 1
			size += len(pair) + 1
		}

		var hstore []byte
		if size == 0 {
			hstore = []byte{'"', '"'}
		} else {
			hstore = make([]byte, 1, size+1)
			hstore[0] = '"'
			for _, pair := range pairs {
				hstore = append(hstore, pair...)
				hstore = append(hstore, ',')
			}
			hstore[len(hstore)-1] = '"'
		}

		out = append(out, hstore...)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type HStoreArrayToStringPtrMapSlice struct {
	Val *[]map[string]*string
}

func (v HStoreArrayToStringPtrMapSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	elems := pgparsehstorearr(data)
	hashes := make([]map[string]*string, len(elems))

	for i := 0; i < len(elems); i++ {
		pairs := elems[i]
		if pairs == nil {
			continue
		}

		hash := make(map[string]*string)
		for j := 0; j < len(pairs); j++ {
			if value := pairs[j][1]; value != nil {
				str := string(value)
				hash[string(pairs[j][0])] = &str
			} else {
				hash[string(pairs[j][0])] = nil
			}
		}

		hashes[i] = hash
	}

	*v.Val = hashes
	return nil
}

type HStoreArrayFromNullStringMapSlice struct {
	Val []map[string]sql.NullString
}

func (v HStoreArrayFromNullStringMapSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, m := range v.Val {
		if m == nil {
			out = append(out, 'N', 'U', 'L', 'L', ',')
			continue
		}

		i, size := 0, 0
		pairs := make([][]byte, len(m))
		for key, val := range m {
			if !val.Valid {
				// len(`\"\"=>NULL`) == 10
				pair := make([]byte, 0, 10+len(key))
				pair = pgAppendQuote2(pair, []byte(key))
				pair = append(pair, '=', '>', 'N', 'U', 'L', 'L')

				pairs[i] = pair
				i += 1
				size += len(pair) + 1
				continue
			}

			// len(`\"\"=>\"\"`) == 10
			pair := make([]byte, 0, 10+len(key)+len(val.String))
			pair = pgAppendQuote2(pair, []byte(key))
			pair = append(pair, '=', '>')
			pair = pgAppendQuote2(pair, []byte(val.String))

			pairs[i] = pair
			i += 1
			size += len(pair) + 1
		}

		var hstore []byte
		if size == 0 {
			hstore = []byte{'"', '"'}
		} else {
			hstore = make([]byte, 1, size+1)
			hstore[0] = '"'
			for _, pair := range pairs {
				hstore = append(hstore, pair...)
				hstore = append(hstore, ',')
			}
			hstore[len(hstore)-1] = '"'
		}

		out = append(out, hstore...)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type HStoreArrayToNullStringMapSlice struct {
	Val *[]map[string]sql.NullString
}

func (v HStoreArrayToNullStringMapSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	elems := pgparsehstorearr(data)
	hashes := make([]map[string]sql.NullString, len(elems))

	for i := 0; i < len(elems); i++ {
		pairs := elems[i]
		if pairs == nil {
			continue
		}

		hash := make(map[string]sql.NullString)
		for j := 0; j < len(pairs); j++ {
			if value := pairs[j][1]; value != nil {
				hash[string(pairs[j][0])] = sql.NullString{String: string(value), Valid: true}
			} else {
				hash[string(pairs[j][0])] = sql.NullString{String: "", Valid: false}
			}
		}

		hashes[i] = hash
	}

	*v.Val = hashes
	return nil
}
