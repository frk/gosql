package pgsql

import (
	"testing"
	"time"
)

func TestTsRange(t *testing.T) {
	testlist2{{
		valuer:  TsRangeFromTimeArray2,
		scanner: TsRangeToTimeArray2,
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
		data: []testdata{
			{
				input:  string(`["1999-01-08 04:05:06","2004-10-19 10:23:54")`),
				output: string(`["1999-01-08 04:05:06","2004-10-19 10:23:54")`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte(`["1999-01-08 04:05:06","2004-10-19 10:23:54")`),
				output: []byte(`["1999-01-08 04:05:06","2004-10-19 10:23:54")`)},
		},
	}}.execute(t, "tsrange")
}
