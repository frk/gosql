package convert

import (
	"time"
)

type TimeToString struct {
	Val *string
}

func (v TimeToString) Scan(src interface{}) error {
	switch data := src.(type) {
	case time.Time:
		*v.Val = data.Format(timeLayout)
	case []byte:
		*v.Val = string(data)
	case string:
		*v.Val = data
	}
	return nil
}

type TimeToByteSlice struct {
	Val *[]byte
}

func (v TimeToByteSlice) Scan(src interface{}) error {
	switch data := src.(type) {
	case time.Time:
		*v.Val = []byte(data.Format(timeLayout))
	case []byte:
		*v.Val = data
	case string:
		*v.Val = []byte(data)
	}
	return nil
}
