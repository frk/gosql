package pg2go

import (
	"time"
)

type DateToTime struct {
	Val *time.Time
}

func (v DateToTime) Scan(src interface{}) error {
	if t, ok := src.(time.Time); ok {
		*v.Val = t.In(time.UTC)
	} else {
		*v.Val = time.Time{}
	}
	return nil
}

type DateToString struct {
	Val *string
}

func (v DateToString) Scan(src interface{}) error {
	if t, ok := src.(time.Time); ok {
		*v.Val = t.Format(dateLayout)
	} else {
		*v.Val = ""
	}
	return nil
}

type DateToByteSlice struct {
	Val *[]byte
}

func (v DateToByteSlice) Scan(src interface{}) error {
	if t, ok := src.(time.Time); ok {
		*v.Val = []byte(t.Format(dateLayout))
	} else {
		*v.Val = nil
	}
	return nil
}
