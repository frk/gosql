package convert

import (
	"testing"
)

func TestVarChar(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			d := new(string)
			return d, d
		},
		data: []testdata{
			{input: string(""), output: string("")},
			{input: string("foo bar"), output: string("foo bar")},
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
			{input: []byte(""), output: []byte("")},
			{input: []byte("foo bar"), output: []byte("foo bar")},
		},
	}}.execute(t, "varchar")
}
