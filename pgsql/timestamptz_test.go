package pgsql

import (
	"testing"
	"time"
)

func TestTimestamptz(t *testing.T) {
	dublin, err := time.LoadLocation("Europe/Dublin")
	if err != nil {
		t.Fatal(err)
	}
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}

	testlist2{{
		data: []testdata{
			{
				input:  timestamptz(2004, 10, 19, 10, 23, 54, 987, dublin),
				output: timestamptz(2004, 10, 19, 10, 23, 54, 987, dublin)},
			{
				input:  timestamptz(1999, 1, 8, 21, 5, 33, 0, tokyo),
				output: timestamptz(1999, 1, 8, 21, 5, 33, 0, tokyo)},
		},
	}, {
		data: []testdata{
			{
				input:  string("1999-01-08 04:05:06+01"),
				output: string("1999-01-08T04:05:06+01:00")},
			{
				input:  string("1999-01-08T04:05:06+02:00"),
				output: string("1999-01-08T03:05:06+01:00")},
			{
				input:  string("1999-01-08T04:05:06.987-05:00"),
				output: string("1999-01-08T10:05:06.987+01:00")},
		},
	}, {
		data: []testdata{
			{
				input:  []byte("1999-01-08 04:05:06+01"),
				output: []byte("1999-01-08T04:05:06+01:00")},
			{
				input:  []byte("1999-01-08T04:05:06+02:00"),
				output: []byte("1999-01-08T03:05:06+01:00")},
			{
				input:  []byte("1999-01-08T04:05:06.987-05:00"),
				output: []byte("1999-01-08T10:05:06.987+01:00")},
		},
	}}.execute(t, "timestamptz")
}
