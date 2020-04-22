package pgsql

import (
	"testing"
	"time"
)

func TestTSRangeArray(t *testing.T) {
	testlist2{{
		valuer:  TsRangeArrayFromTimeArray2Slice,
		scanner: TsRangeArrayToTimeArray2Slice,
		data: []testdata{
			{input: [][2]time.Time(nil), output: [][2]time.Time(nil)},
			{input: [][2]time.Time{}, output: [][2]time.Time{}},
			{
				input: [][2]time.Time{{
					timestamp(1999, 1, 8, 21, 5, 33, 0),
					timestamp(2004, 10, 19, 10, 23, 54, 789),
				}},
				output: [][2]time.Time{{
					timestamp(1999, 1, 8, 21, 5, 33, 0),
					timestamp(2004, 10, 19, 10, 23, 54, 789),
				}}},
			{
				input: [][2]time.Time{{
					timestamp(1999, 1, 8, 21, 5, 33, 0),
					timestamp(2004, 10, 19, 10, 23, 54, 789),
				}, {
					timestamp(2004, 10, 19, 10, 23, 54, 789),
					timestamp(2019, 1, 8, 21, 5, 33, 0),
				}},
				output: [][2]time.Time{{
					timestamp(1999, 1, 8, 21, 5, 33, 0),
					timestamp(2004, 10, 19, 10, 23, 54, 789),
				}, {
					timestamp(2004, 10, 19, 10, 23, 54, 789),
					timestamp(2019, 1, 8, 21, 5, 33, 0),
				}}},
		},
	}, {
		data: []testdata{
			{input: string(`{}`), output: string(`{}`)},
			{
				input:  string(`{"[\"1999-01-08 04:05:06\",\"2004-10-19 10:23:54\")"}`),
				output: string(`{"[\"1999-01-08 04:05:06\",\"2004-10-19 10:23:54\")"}`)},
			{
				input: string(`{"[\"1999-01-08 04:05:06\",\"2004-10-19 10:23:54\")",` +
					`"(\"1999-01-08 04:05:06\",\"2004-10-19 10:23:54\"]"}`),
				output: string(`{"[\"1999-01-08 04:05:06\",\"2004-10-19 10:23:54\")",` +
					`"(\"1999-01-08 04:05:06\",\"2004-10-19 10:23:54\"]"}`)},
			{
				input:  string(`{"(,)",NULL}`),
				output: string(`{"(,)",NULL}`)},
		},
	}, {
		data: []testdata{
			{input: []byte(`{}`), output: []byte(`{}`)},
			{
				input:  []byte(`{"[\"1999-01-08 04:05:06\",\"2004-10-19 10:23:54\")"}`),
				output: []byte(`{"[\"1999-01-08 04:05:06\",\"2004-10-19 10:23:54\")"}`)},
			{
				input: []byte(`{"[\"1999-01-08 04:05:06\",\"2004-10-19 10:23:54\")",` +
					`"(\"1999-01-08 04:05:06\",\"2004-10-19 10:23:54\"]"}`),
				output: []byte(`{"[\"1999-01-08 04:05:06\",\"2004-10-19 10:23:54\")",` +
					`"(\"1999-01-08 04:05:06\",\"2004-10-19 10:23:54\"]"}`)},
			{
				input:  []byte(`{"(,)",NULL}`),
				output: []byte(`{"(,)",NULL}`)},
		},
	}}.execute(t, "tsrangearr")
}
