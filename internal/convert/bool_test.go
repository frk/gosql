package convert

import (
	"testing"
)

func TestBool(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return nil // bool
		},
		scanner: func() (interface{}, interface{}) {
			d := new(bool)
			return d, d
		},
		data: []testdata{
			{input: true, output: boolptr(true)},
			{input: false, output: boolptr(false)},
		},
	}, {
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			d := new(string)
			return d, d
		},
		data: []testdata{
			{input: string("true"), output: strptr(`true`)},
			{input: string("false"), output: strptr(`false`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			d := new([]byte)
			return d, d
		},
		data: []testdata{
			{input: []byte("true"), output: bytesptr(`true`)},
			{input: []byte("false"), output: bytesptr(`false`)},
		},
	}}.execute(t, "bool")
}
