package pg2go

import (
	"testing"
)

func TestTime(t *testing.T) {
	testlist2{{
		data: []testdata{
			{
				input:  timeval(21, 5, 33, 0),
				output: timeval(21, 5, 33, 0)},
			{
				input:  timeval(4, 5, 6, 789),
				output: timeval(4, 5, 6, 789)},
		},
	}, {
		scanner: TimeToString,
		data: []testdata{
			{
				input:  string("21:05:33"),
				output: string("21:05:33")},
			{
				input:  string("04:05:06.789"),
				output: string("04:05:06.789")},
		},
	}, {
		scanner: TimeToByteSlice,
		data: []testdata{
			{
				input:  []byte("21:05:33"),
				output: []byte("21:05:33")},
			{
				input:  []byte("04:05:06.789"),
				output: []byte("04:05:06.789")},
		},
	}}.execute(t, "time")
}
