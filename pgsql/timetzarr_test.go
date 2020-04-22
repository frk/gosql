package pgsql

import (
	"testing"
	"time"
)

func TestTimetzArray(t *testing.T) {
	dublin, err := time.LoadLocation("Europe/Dublin")
	if err != nil {
		t.Fatal(err)
	}

	testlist2{{
		valuer:  TimetzArrayFromTimeSlice,
		scanner: TimetzArrayToTimeSlice,
		data: []testdata{
			{input: []time.Time(nil), output: []time.Time(nil)},
			{input: []time.Time{}, output: []time.Time{}},
			{
				input:  []time.Time{timetzval(4, 5, 6, 789, dublin)},
				output: []time.Time{timetzval(4, 5, 6, 789, dublin)}},
			{
				input: []time.Time{
					timetzval(21, 5, 33, 0, dublin),
					timetzval(4, 5, 6, 789, dublin)},
				output: []time.Time{
					timetzval(21, 5, 33, 0, dublin),
					timetzval(4, 5, 6, 789, dublin)}},
		},
	}, {
		data: []testdata{
			{input: string("{}"), output: string("{}")},
			{
				input:  string("{04:05:06.789-08}"),
				output: string("{04:05:06.789-08}")},
			{
				input:  string("{21:05:33-08,04:05:06.789+04}"),
				output: string("{21:05:33-08,04:05:06.789+04}")},
		},
	}, {
		data: []testdata{
			{input: []byte("{}"), output: []byte("{}")},
			{
				input:  []byte("{04:05:06.789-08}"),
				output: []byte("{04:05:06.789-08}")},
			{
				input:  []byte("{21:05:33-08,04:05:06.789+04}"),
				output: []byte("{21:05:33-08,04:05:06.789+04}")},
		},
	}}.execute(t, "timetzarr")
}
