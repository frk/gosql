package pgsql

import (
	"testing"
	"time"
)

func TestTstzRangeArray(t *testing.T) {
	dublin, err := time.LoadLocation("Europe/Dublin")
	if err != nil {
		t.Fatal(err)
	}
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}

	testlist2{{
		valuer:  TstzRangeArrayFromTimeArray2Slice,
		scanner: TstzRangeArrayToTimeArray2Slice,
		data: []testdata{
			{input: [][2]time.Time(nil), output: [][2]time.Time(nil)},
			{input: [][2]time.Time{}, output: [][2]time.Time{}},
			{
				input: [][2]time.Time{{
					timestamptz(1999, 1, 8, 21, 5, 33, 0, dublin),
					timestamptz(2004, 10, 19, 10, 23, 54, 789, tokyo),
				}},
				output: [][2]time.Time{{
					timestamptz(1999, 1, 8, 21, 5, 33, 0, dublin),
					timestamptz(2004, 10, 19, 10, 23, 54, 789, tokyo),
				}}},
			{
				input: [][2]time.Time{{
					timestamptz(1999, 1, 8, 21, 5, 33, 0, time.UTC),
					timestamptz(2004, 10, 19, 10, 23, 54, 789, time.UTC),
				}, {
					timestamptz(2004, 10, 19, 10, 23, 54, 789, dublin),
					timestamptz(2019, 1, 8, 21, 5, 33, 0, dublin),
				}},
				output: [][2]time.Time{{
					timestamptz(1999, 1, 8, 21, 5, 33, 0, time.UTC),
					timestamptz(2004, 10, 19, 10, 23, 54, 789, time.UTC),
				}, {
					timestamptz(2004, 10, 19, 10, 23, 54, 789, dublin),
					timestamptz(2019, 1, 8, 21, 5, 33, 0, dublin),
				}}},
		},
	}, {
		data: []testdata{
			{input: string(`{}`), output: string(`{}`)},
			{
				input:  string(`{"[\"1999-01-08 04:05:06+03\",\"2004-10-19 10:23:54-02\")"}`),
				output: string(`{"[\"1999-01-08 02:05:06+01\",\"2004-10-19 14:23:54+02\")"}`)},
			{
				input: string(`{"[\"1999-01-08 04:05:06-02\",\"2004-10-19 10:23:54-02\")",` +
					`"(\"1999-01-08 04:05:06-05\",\"2004-10-19 10:23:54-05\"]"}`),
				output: string(`{"[\"1999-01-08 07:05:06+01\",\"2004-10-19 14:23:54+02\")",` +
					`"(\"1999-01-08 10:05:06+01\",\"2004-10-19 17:23:54+02\"]"}`)},
			{
				input:  string(`{"(,)",NULL}`),
				output: string(`{"(,)",NULL}`)},
		},
	}, {
		data: []testdata{
			{input: []byte(`{}`), output: []byte(`{}`)},
			{
				input:  []byte(`{"[\"1999-01-08 04:05:06+03\",\"2004-10-19 10:23:54-02\")"}`),
				output: []byte(`{"[\"1999-01-08 02:05:06+01\",\"2004-10-19 14:23:54+02\")"}`)},
			{
				input: []byte(`{"[\"1999-01-08 04:05:06-02\",\"2004-10-19 10:23:54-02\")",` +
					`"(\"1999-01-08 04:05:06-05\",\"2004-10-19 10:23:54-05\"]"}`),
				output: []byte(`{"[\"1999-01-08 07:05:06+01\",\"2004-10-19 14:23:54+02\")",` +
					`"(\"1999-01-08 10:05:06+01\",\"2004-10-19 17:23:54+02\"]"}`)},
			{
				input:  []byte(`{"(,)",NULL}`),
				output: []byte(`{"(,)",NULL}`)},
		},
	}}.execute(t, "tstzrangearr")
}
