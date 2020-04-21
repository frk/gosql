package pg2go

import (
	"testing"
	"time"
)

func TestTimestampArray(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(TimestampArrayFromTimeSlice)
		},
		scanner: func() (interface{}, interface{}) {
			v := TimestampArrayToTimeSlice{Val: new([]time.Time)}
			return v, v.Val
		},
		data: []testdata{
			{input: []time.Time(nil), output: []time.Time(nil)},
			{input: []time.Time{}, output: []time.Time{}},
			{
				input:  []time.Time{timestamp(1999, 1, 8, 4, 5, 6, 0)},
				output: []time.Time{timestamp(1999, 1, 8, 4, 5, 6, 0)}},
			{
				input:  []time.Time{timestamp(1999, 1, 8, 4, 5, 6, 987), timestamp(2004, 10, 19, 10, 23, 54, 0)},
				output: []time.Time{timestamp(1999, 1, 8, 4, 5, 6, 987), timestamp(2004, 10, 19, 10, 23, 54, 0)}},
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
			{input: string(`{}`), output: string(`{}`)},
			{
				input:  string(`{"1999-01-08 04:05:06"}`),
				output: string(`{"1999-01-08 04:05:06"}`)},
			{
				input:  string(`{"1999-01-08 04:05:06.987","2004-10-19 10:23:54"}`),
				output: string(`{"1999-01-08 04:05:06.987","2004-10-19 10:23:54"}`)},
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
			{input: []byte(`{}`), output: []byte(`{}`)},
			{
				input:  []byte(`{"1999-01-08 04:05:06"}`),
				output: []byte(`{"1999-01-08 04:05:06"}`)},
			{
				input:  []byte(`{"1999-01-08 04:05:06.987","2004-10-19 10:23:54"}`),
				output: []byte(`{"1999-01-08 04:05:06.987","2004-10-19 10:23:54"}`)},
		},
	}}.execute(t, "timestamparr")
}
