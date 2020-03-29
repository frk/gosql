package convert

import (
	"time"
)

type DateToTime struct {
	V *time.Time
}

func (v DateToTime) Scan(src interface{}) error {
	if t, ok := src.(time.Time); ok {
		*v.V = t.In(time.UTC)
	} else {
		*v.V = time.Time{}
	}
	return nil
}

type DateToString struct {
	V *string
}

func (v DateToString) Scan(src interface{}) error {
	if t, ok := src.(time.Time); ok {
		*v.V = t.Format(dateLayout)
	} else {
		*v.V = ""
	}
	return nil
}

type DateToByteSlice struct {
	V *[]byte
}

func (v DateToByteSlice) Scan(src interface{}) error {
	if t, ok := src.(time.Time); ok {
		*v.V = []byte(t.Format(dateLayout))
	} else {
		*v.V = nil
	}
	return nil
}
