package pgsql

import (
	"testing"
	"time"
)

func TestTstzRange(t *testing.T) {
	dublin, err := time.LoadLocation("Europe/Dublin")
	if err != nil {
		t.Fatal(err)
	}

	testlist2{{
		valuer:  TstzRangeFromTimeArray2,
		scanner: TstzRangeToTimeArray2,
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
		data: []testdata{
			{
				input:  string(`["1999-01-08 04:05:06-08:00","2004-10-19 10:23:54+04:00")`),
				output: string(`["1999-01-08 13:05:06+01","2004-10-19 08:23:54+02")`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte(`["1999-01-08 04:05:06-08:00","2004-10-19 10:23:54+04:00")`),
				output: []byte(`["1999-01-08 13:05:06+01","2004-10-19 08:23:54+02")`)},
		},
	}}.execute(t, "tstzrange")
}
