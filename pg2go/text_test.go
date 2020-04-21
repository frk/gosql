package pg2go

import (
	"testing"
)

func TestText(t *testing.T) {
	testlist2{{
		data: []testdata{
			{input: string("foo bar"), output: string("foo bar")},
		},
	}, {
		data: []testdata{
			{input: []byte("foo bar"), output: []byte("foo bar")},
		},
	}}.execute(t, "text")
}
