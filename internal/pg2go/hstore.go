package pg2go

import (
	"database/sql"
	"database/sql/driver"
)

type HStoreFromStringMap struct {
	Val map[string]string
}

func (v HStoreFromStringMap) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	}

	out := []byte{}
	for key, val := range v.Val {
		out = pgAppendQuote1(out, []byte(key))
		out = append(out, '=', '>')
		out = pgAppendQuote1(out, []byte(val))
		out = append(out, ',')
	}

	if len(out) > 0 {
		out = out[:len(out)-1] // drop the last ','
	}

	return out, nil
}

type HStoreToStringMap struct {
	Val *map[string]string
}

func (v HStoreToStringMap) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	elems := pgparsehstore(data)

	hash := make(map[string]string)
	for i := 0; i < len(elems); i++ {
		if value := elems[i][1]; value != nil {
			hash[string(elems[i][0])] = string(value)
		} else {
			hash[string(elems[i][0])] = ""
		}
	}

	*v.Val = hash
	return nil
}

type HStoreFromStringPtrMap struct {
	Val map[string]*string
}

func (v HStoreFromStringPtrMap) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	}

	out := []byte{}
	for key, val := range v.Val {
		out = pgAppendQuote1(out, []byte(key))
		out = append(out, '=', '>')
		if val != nil {
			out = pgAppendQuote1(out, []byte(*val))
		} else {
			out = append(out, 'N', 'U', 'L', 'L')
		}
		out = append(out, ',')
	}

	if len(out) > 0 {
		out = out[:len(out)-1] // drop the last ','
	}
	return out, nil
}

type HStoreToStringPtrMap struct {
	Val *map[string]*string
}

func (v HStoreToStringPtrMap) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	elems := pgparsehstore(data)

	hash := make(map[string]*string)
	for i := 0; i < len(elems); i++ {
		if value := elems[i][1]; value != nil {
			str := string(value)
			hash[string(elems[i][0])] = &str
		} else {
			hash[string(elems[i][0])] = nil
		}
	}

	*v.Val = hash
	return nil
}

type HStoreFromNullStringMap struct {
	Val map[string]sql.NullString
}

func (v HStoreFromNullStringMap) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	}

	out := []byte{}
	for key, val := range v.Val {
		out = pgAppendQuote1(out, []byte(key))
		out = append(out, '=', '>')
		if val.Valid {
			out = pgAppendQuote1(out, []byte(val.String))
		} else {
			out = append(out, 'N', 'U', 'L', 'L')
		}
		out = append(out, ',')
	}

	if len(out) > 0 {
		out = out[:len(out)-1] // drop the last ','
	}
	return out, nil
}

type HStoreToNullStringMap struct {
	Val *map[string]sql.NullString
}

func (v HStoreToNullStringMap) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	elems := pgparsehstore(data)

	hash := make(map[string]sql.NullString)
	for i := 0; i < len(elems); i++ {
		if value := elems[i][1]; value != nil {
			str := sql.NullString{String: string(value), Valid: true}
			hash[string(elems[i][0])] = str
		} else {
			hash[string(elems[i][0])] = sql.NullString{}
		}
	}

	*v.Val = hash
	return nil
}
