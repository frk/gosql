package pgsql

import (
	"database/sql"
	"database/sql/driver"
)

// HStoreFromStringMap returns a driver.Valuer that produces a PostgreSQL hstore from the given Go map[string]string.
func HStoreFromStringMap(val map[string]string) driver.Valuer {
	return hstoreFromStringMap{val: val}
}

// HStoreToStringMap returns an sql.Scanner that converts a PostgreSQL hstore into a Go map[string]string and sets it to val.
func HStoreToStringMap(val *map[string]string) sql.Scanner {
	return hstoreToStringMap{val: val}
}

// HStoreFromStringPtrMap returns a driver.Valuer that produces a PostgreSQL hstore from the given Go map[string]*string.
func HStoreFromStringPtrMap(val map[string]*string) driver.Valuer {
	return hstoreFromStringPtrMap{val: val}
}

// HStoreToStringPtrMap returns an sql.Scanner that converts a PostgreSQL hstore into a Go map[string]*string and sets it to val.
func HStoreToStringPtrMap(val *map[string]*string) sql.Scanner {
	return hstoreToStringPtrMap{val: val}
}

// HStoreFromNullStringMap returns a driver.Valuer that produces a PostgreSQL hstore from the given Go map[string]sql.NullString.
func HStoreFromNullStringMap(val map[string]sql.NullString) driver.Valuer {
	return hstoreFromNullStringMap{val: val}
}

// HStoreToNullStringMap returns an sql.Scanner that converts a PostgreSQL hstore into a Go map[string]sql.NullString and sets it to val.
func HStoreToNullStringMap(val *map[string]sql.NullString) sql.Scanner {
	return hstoreToNullStringMap{val: val}
}

type hstoreFromStringMap struct {
	val map[string]string
}

func (v hstoreFromStringMap) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	}

	out := []byte{}
	for key, val := range v.val {
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

type hstoreToStringMap struct {
	val *map[string]string
}

func (v hstoreToStringMap) Scan(src interface{}) error {
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

	*v.val = hash
	return nil
}

type hstoreFromStringPtrMap struct {
	val map[string]*string
}

func (v hstoreFromStringPtrMap) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	}

	out := []byte{}
	for key, val := range v.val {
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

type hstoreToStringPtrMap struct {
	val *map[string]*string
}

func (v hstoreToStringPtrMap) Scan(src interface{}) error {
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

	*v.val = hash
	return nil
}

type hstoreFromNullStringMap struct {
	val map[string]sql.NullString
}

func (v hstoreFromNullStringMap) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	}

	out := []byte{}
	for key, val := range v.val {
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

type hstoreToNullStringMap struct {
	val *map[string]sql.NullString
}

func (v hstoreToNullStringMap) Scan(src interface{}) error {
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

	*v.val = hash
	return nil
}
