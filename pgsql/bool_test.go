package pgsql

import (
	"testing"
)

func TestBool(t *testing.T) {
	testlist2{{
		data: []testdata{
			{input: bool(true), output: bool(true)},
			{input: bool(false), output: bool(false)},
		},
	}, {
		data: []testdata{
			{input: string("true"), output: string(`true`)},
			{input: string("false"), output: string(`false`)},
		},
	}, {
		data: []testdata{
			{input: []byte("true"), output: []byte(`true`)},
			{input: []byte("false"), output: []byte(`false`)},
		},
	}}.execute(t, "bool")
}
