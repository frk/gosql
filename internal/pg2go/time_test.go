package pg2go

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return nil // time.Time
		},
		scanner: func() (interface{}, interface{}) {
			v := new(time.Time)
			return v, v
		},
		data: []testdata{
			{
				input:  timeval(21, 5, 33, 0),
				output: timeval(21, 5, 33, 0)},
			{
				input:  timeval(4, 5, 6, 789),
				output: timeval(4, 5, 6, 789)},
		},
	}, {
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			v := TimeToString{Val: new(string)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  string("21:05:33"),
				output: string("21:05:33")},
			{
				input:  string("04:05:06.789"),
				output: string("04:05:06.789")},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			v := TimeToByteSlice{Val: new([]byte)}
			return v, v.Val
		},
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
