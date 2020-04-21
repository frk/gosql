package pg2go

import (
	"testing"
	"time"
)

func TestTimetz(t *testing.T) {
	dublin, err := time.LoadLocation("Europe/Dublin")
	if err != nil {
		t.Fatal(err)
	}

	testlist2{{
		data: []testdata{
			{
				input:  timetzval(21, 5, 33, 0, dublin),
				output: timetzval(21, 5, 33, 0, dublin)},
			{
				input:  timetzval(4, 5, 6, 789, time.UTC),
				output: timetzval(4, 5, 6, 789, time.UTC)},
		},
	}, {
		scanner: TimetzToString,
		data: []testdata{
			{
				input:  string("21:05:33+01"),
				output: string("21:05:33+01:00")},
			{
				input:  string("04:05:06.789-08"),
				output: string("04:05:06.789-08:00")},
		},
	}, {
		scanner: TimetzToByteSlice,
		data: []testdata{
			{
				input:  []byte("21:05:33+01"),
				output: []byte("21:05:33+01:00")},
			{
				input:  []byte("04:05:06.789-08"),
				output: []byte("04:05:06.789-08:00")},
		},
	}}.execute(t, "timetz")
}
