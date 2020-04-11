package convert

import (
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return nil // time.Time
		},
		scanner: func() (interface{}, interface{}) {
			s := &DateToTime{Val: new(time.Time)}
			return s, s.Val
		},
		data: []testdata{
			{input: dateptr(1999, 1, 8), output: dateptr(1999, 1, 8)},
			{input: dateptr(2001, 5, 5), output: dateptr(2001, 5, 5)},
			{input: dateptr(2020, 3, 28), output: dateptr(2020, 3, 28)},
		},
	}, {
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			s := &DateToString{Val: new(string)}
			return s, s.Val
		},
		data: []testdata{
			{input: string(`1999-01-08`), output: strptr(`1999-01-08`)},
			{input: string(`2001-05-05`), output: strptr(`2001-05-05`)},
			{input: string(`2020-03-28`), output: strptr(`2020-03-28`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			s := &DateToByteSlice{Val: new([]byte)}
			return s, s.Val
		},
		data: []testdata{
			{input: []byte(`1999-01-08`), output: bytesptr(`1999-01-08`)},
			{input: []byte(`2001-05-05`), output: bytesptr(`2001-05-05`)},
			{input: []byte(`2020-03-28`), output: bytesptr(`2020-03-28`)},
		},
	}}.execute(t, "date")
}
