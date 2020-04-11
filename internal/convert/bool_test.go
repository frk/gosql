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
			{input: true, output: bool(true)},
			{input: false, output: bool(false)},
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
			{input: string("true"), output: string(`true`)},
			{input: string("false"), output: string(`false`)},
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
			{input: []byte("true"), output: []byte(`true`)},
			{input: []byte("false"), output: []byte(`false`)},
		},
	}}.execute(t, "bool")
}
