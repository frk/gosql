package pgsql

import (
	"testing"
	"time"
)

func TestDateRangeArray(t *testing.T) {
	testlist2{{
		valuer:  DateRangeArrayFromTimeArray2Slice,
		scanner: DateRangeArrayToTimeArray2Slice,
		data: []testdata{
			{input: [][2]time.Time(nil), output: [][2]time.Time(nil)},
			{
				input:  [][2]time.Time{{}, {}, {}},
				output: [][2]time.Time{{}, {}, {}}},
			{
				input:  [][2]time.Time{{dateval(1999, 1, 8), dateval(2001, 5, 5)}},
				output: [][2]time.Time{{dateval(1999, 1, 8), dateval(2001, 5, 5)}}},
			{
				input: [][2]time.Time{
					{{}, dateval(2001, 5, 5)},
					{dateval(1999, 1, 8), {}},
					{dateval(1999, 1, 8), dateval(2001, 5, 5)},
				},
				output: [][2]time.Time{
					{{}, dateval(2001, 5, 5)},
					{dateval(1999, 1, 8), {}},
					{dateval(1999, 1, 8), dateval(2001, 5, 5)}}},
		},
	}, {
		data: []testdata{
			{
				input:  string(`{"(,)","(,)","(,)"}`),
				output: string(`{"(,)","(,)","(,)"}`)},
			{
				input:  string(`{"[1999-01-08,2001-05-05)"}`),
				output: string(`{"[1999-01-08,2001-05-05)"}`)},
			{
				input:  string(`{"(,2001-05-05)","[1999-01-08,)","[1999-01-08,2001-05-05)"}`),
				output: string(`{"(,2001-05-05)","[1999-01-08,)","[1999-01-08,2001-05-05)"}`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte(`{"(,)","(,)","(,)"}`),
				output: []byte(`{"(,)","(,)","(,)"}`)},
			{
				input:  []byte(`{"[1999-01-08,2001-05-05)"}`),
				output: []byte(`{"[1999-01-08,2001-05-05)"}`)},
			{
				input:  []byte(`{"(,2001-05-05)","[1999-01-08,)","[1999-01-08,2001-05-05)"}`),
				output: []byte(`{"(,2001-05-05)","[1999-01-08,)","[1999-01-08,2001-05-05)"}`)},
		},
	}}.execute(t, "daterangearr")
}
