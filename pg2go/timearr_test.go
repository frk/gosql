package pg2go

import (
	"testing"
	"time"
)

func TestTimeArray(t *testing.T) {
	testlist2{{
		valuer:  TimeArrayFromTimeSlice,
		scanner: TimeArrayToTimeSlice,
		data: []testdata{
			{input: []time.Time(nil), output: []time.Time(nil)},
			{input: []time.Time{}, output: []time.Time{}},
			{
				input:  []time.Time{timeval(4, 5, 6, 789)},
				output: []time.Time{timeval(4, 5, 6, 789)}},
			{
				input:  []time.Time{timeval(21, 5, 33, 0), timeval(4, 5, 6, 789)},
				output: []time.Time{timeval(21, 5, 33, 0), timeval(4, 5, 6, 789)}},
		},
	}, {
		data: []testdata{
			{input: string("{}"), output: string("{}")},
			{
				input:  string("{04:05:06.789}"),
				output: string("{04:05:06.789}")},
			{
				input:  string("{21:05:33,04:05:06.789}"),
				output: string("{21:05:33,04:05:06.789}")},
		},
	}, {
		data: []testdata{
			{input: []byte("{}"), output: []byte("{}")},
			{
				input:  []byte("{04:05:06.789}"),
				output: []byte("{04:05:06.789}")},
			{
				input:  []byte("{21:05:33,04:05:06.789}"),
				output: []byte("{21:05:33,04:05:06.789}")},
		},
	}}.execute(t, "timearr")
}
