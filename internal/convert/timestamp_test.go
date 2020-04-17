package convert

import (
	"testing"
	"time"
)

func TestTimestamp(t *testing.T) {
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
				input:  timestamp(1999, 1, 8, 21, 5, 33, 0),
				output: timestamp(1999, 1, 8, 21, 5, 33, 0)},
		},
	}, {
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			v := new(string)
			return v, v
		},
		data: []testdata{
			{
				input:  string("1999-01-08 04:05:06"),
				output: string("1999-01-08T04:05:06Z")},
			{
				input:  string("1999-01-08T04:05:06Z"),
				output: string("1999-01-08T04:05:06Z")},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			v := new([]byte)
			return v, v
		},
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
