package pg2go

import (
	"testing"
)

func TestVarChar(t *testing.T) {
	testlist2{{
		data: []testdata{
			{input: string(""), output: string("")},
			{input: string("foo bar"), output: string("foo bar")},
		},
	}, {
		data: []testdata{
			{input: []byte(""), output: []byte("")},
			{input: []byte("foo bar"), output: []byte("foo bar")},
		},
	}}.execute(t, "varchar")
}
