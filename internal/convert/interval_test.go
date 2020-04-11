package convert

import (
	"testing"
)

func TestInterval(t *testing.T) {
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
				input:  `1 day`,
				output: strptr(`1 day`)},
			{
				input:  `-5 years -4 mons -00:34:00`,
				output: strptr(`-5 years -4 mons -00:34:00`)},
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
				input:  []byte(`1 day`),
				output: bytesptr(`1 day`)},
			{
				input:  []byte(`-5 years -4 mons -00:34:00`),
				output: bytesptr(`-5 years -4 mons -00:34:00`)},
		},
	}}.execute(t, "interval")
}
