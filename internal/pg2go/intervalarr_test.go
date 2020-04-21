package pg2go

import (
	"testing"
)

func TestIntervalArray(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			s := new(string)
			return s, s
		},
		data: []testdata{
			{
				input:  string(`{"1 day","-5 years -4 mons -00:34:00"}`),
				output: string(`{"1 day","-5 years -4 mons -00:34:00"}`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			s := new([]byte)
			return s, s
		},
		data: []testdata{
			{
				input:  []byte(`{"1 day","-5 years -4 mons -00:34:00"}`),
				output: []byte(`{"1 day","-5 years -4 mons -00:34:00"}`)},
		},
	}}.execute(t, "intervalarr")
}
