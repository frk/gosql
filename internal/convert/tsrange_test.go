package convert

import (
	"testing"
	"time"
)

func TestTsRange(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(TsRangeFromTimeArray2)
		},
		scanner: func() (interface{}, interface{}) {
			v := TsRangeToTimeArray2{Val: new([2]time.Time)}
			return v, v.Val
		},
		data: []testdata{
			{
				input: [2]time.Time{
					timestamp(1999, 1, 8, 21, 5, 33, 0),
					timestamp(2004, 10, 19, 10, 23, 54, 0)},
				output: [2]time.Time{
					timestamp(1999, 1, 8, 21, 5, 33, 0),
					timestamp(2004, 10, 19, 10, 23, 54, 0)}},
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
				input:  string(`["1999-01-08 04:05:06","2004-10-19 10:23:54")`),
				output: string(`["1999-01-08 04:05:06","2004-10-19 10:23:54")`)},
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
				input:  []byte(`["1999-01-08 04:05:06","2004-10-19 10:23:54")`),
				output: []byte(`["1999-01-08 04:05:06","2004-10-19 10:23:54")`)},
		},
	}}.execute(t, "tsrange")
}
