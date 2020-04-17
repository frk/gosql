package convert

import (
	"testing"
	"time"
)

func TestTstzRange(t *testing.T) {
	dublin, err := time.LoadLocation("Europe/Dublin")
	if err != nil {
		t.Fatal(err)
	}

	testlist{{
		valuer: func() interface{} {
			return new(TstzRangeFromTimeArray2)
		},
		scanner: func() (interface{}, interface{}) {
			v := TstzRangeToTimeArray2{Val: new([2]time.Time)}
			return v, v.Val
		},
		data: []testdata{
			{
				input: [2]time.Time{
					timestamptz(1999, 1, 8, 21, 5, 33, 0, dublin),
					timestamptz(2004, 10, 19, 10, 23, 54, 0, time.UTC)},
				output: [2]time.Time{
					timestamptz(1999, 1, 8, 21, 5, 33, 0, dublin),
					timestamptz(2004, 10, 19, 10, 23, 54, 0, time.UTC)}},
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
				input:  string(`["1999-01-08 04:05:06-08:00","2004-10-19 10:23:54+04:00")`),
				output: string(`["1999-01-08 13:05:06+01","2004-10-19 08:23:54+02")`)},
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
				input:  []byte(`["1999-01-08 04:05:06-08:00","2004-10-19 10:23:54+04:00")`),
				output: []byte(`["1999-01-08 13:05:06+01","2004-10-19 08:23:54+02")`)},
		},
	}}.execute(t, "tstzrange")
}
