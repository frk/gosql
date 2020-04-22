package pgsql

import (
	"testing"
	"time"
)

func TestTimestamptzArray(t *testing.T) {
	dublin, err := time.LoadLocation("Europe/Dublin")
	if err != nil {
		t.Fatal(err)
	}
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}

	testlist2{{
		valuer:  TimestamptzArrayFromTimeSlice,
		scanner: TimestamptzArrayToTimeSlice,
		data: []testdata{
			{input: []time.Time(nil), output: []time.Time(nil)},
			{input: []time.Time{}, output: []time.Time{}},
			{
				input:  []time.Time{timestamptz(1999, 1, 8, 4, 5, 6, 0, dublin)},
				output: []time.Time{timestamptz(1999, 1, 8, 4, 5, 6, 0, dublin)}},
			{
				input: []time.Time{
					timestamptz(1999, 1, 8, 4, 5, 6, 987, dublin),
					timestamptz(2004, 10, 19, 10, 23, 54, 0, tokyo)},
				output: []time.Time{
					timestamptz(1999, 1, 8, 4, 5, 6, 987, dublin),
					timestamptz(2004, 10, 19, 10, 23, 54, 0, tokyo)}},
		},
	}, {
		data: []testdata{
			{input: string(`{}`), output: string(`{}`)},
			{
				input:  string(`{"1999-01-08 04:05:06+01"}`),
				output: string(`{"1999-01-08 04:05:06+01"}`)},
			{
				input:  string(`{"1999-01-08 04:05:06.987-05:00","1999-01-08 04:05:06+02:00"}`),
				output: string(`{"1999-01-08 10:05:06.987+01","1999-01-08 03:05:06+01"}`)},
		},
	}, {
		data: []testdata{
			{input: []byte(`{}`), output: []byte(`{}`)},
			{
				input:  []byte(`{"1999-01-08 04:05:06+01"}`),
				output: []byte(`{"1999-01-08 04:05:06+01"}`)},
			{
				input:  []byte(`{"1999-01-08 04:05:06.987-05:00","1999-01-08 04:05:06+02:00"}`),
				output: []byte(`{"1999-01-08 10:05:06.987+01","1999-01-08 03:05:06+01"}`)},
		},
	}}.execute(t, "timestamptzarr")
}
