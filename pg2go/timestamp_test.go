package pg2go

import (
	"testing"
)

func TestTimestamp(t *testing.T) {
	testlist2{{
		data: []testdata{
			{
				input:  timestamp(1999, 1, 8, 21, 5, 33, 0),
				output: timestamp(1999, 1, 8, 21, 5, 33, 0)},
		},
	}, {
		data: []testdata{
			{
				input:  string("1999-01-08 04:05:06"),
				output: string("1999-01-08T04:05:06Z")},
			{
				input:  string("1999-01-08T04:05:06Z"),
				output: string("1999-01-08T04:05:06Z")},
		},
	}, {
		data: []testdata{
			{
				input:  []byte("1999-01-08 04:05:06"),
				output: []byte("1999-01-08T04:05:06Z")},
			{
				input:  []byte("1999-01-08T04:05:06Z"),
				output: []byte("1999-01-08T04:05:06Z")},
		},
	}}.execute(t, "timestamp")
}
