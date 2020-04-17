package convert

import (
	"time"
)

type TimetzToString struct {
	Val *string
}

func (v TimetzToString) Scan(src interface{}) error {
	switch data := src.(type) {
	case time.Time:
		*v.Val = data.Format(timetzLayout)
	case []byte:
		*v.Val = string(data)
	case string:
		*v.Val = data
	}
	return nil
}

type TimetzToByteSlice struct {
	Val *[]byte
}

func (v TimetzToByteSlice) Scan(src interface{}) error {
	switch data := src.(type) {
	case time.Time:
		*v.Val = []byte(data.Format(timetzLayout))
	case []byte:
		*v.Val = data
	case string:
		*v.Val = []byte(data)
	}
	return nil
}
