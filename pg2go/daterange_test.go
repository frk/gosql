package pg2go

import (
	"testing"
	"time"
)

func TestDateRange(t *testing.T) {
	testlist2{{
		valuer:  DateRangeFromTimeArray2,
		scanner: DateRangeToTimeArray2,
		data: []testdata{
			{input: [2]time.Time{}, output: [2]time.Time{}},
			{
				input:  [2]time.Time{dateval(1999, 1, 8), dateval(2001, 5, 5)},
				output: [2]time.Time{dateval(1999, 1, 8), dateval(2001, 5, 5)}},
			{
				input:  [2]time.Time{{}, dateval(2001, 5, 5)},
				output: [2]time.Time{{}, dateval(2001, 5, 5)}},
			{
				input:  [2]time.Time{dateval(1999, 1, 8), {}},
				output: [2]time.Time{dateval(1999, 1, 8), {}}},
			{
				input:  [2]time.Time{{}, {}},
				output: [2]time.Time{{}, {}}},
		},
	}, {
		data: []testdata{
			{
				input:  string("[1999-01-08,2001-05-05)"),
				output: string(`[1999-01-08,2001-05-05)`)},
			{
				input:  string("[1999-01-08,2001-05-05]"),
				output: string(`[1999-01-08,2001-05-06)`)},
			{
				input:  string("(1999-01-08,2001-05-05]"),
				output: string(`[1999-01-09,2001-05-06)`)},
			{
				input:  string("[,2001-05-05)"),
				output: string(`(,2001-05-05)`)},
			{
				input:  string("[1999-01-08,]"),
				output: string(`[1999-01-08,)`)},
			{
				input:  string("[,]"),
				output: string(`(,)`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte("[1999-01-08,2001-05-05)"),
				output: []byte(`[1999-01-08,2001-05-05)`)},
			{
				input:  []byte("[1999-01-08,2001-05-05]"),
				output: []byte(`[1999-01-08,2001-05-06)`)},
			{
				input:  []byte("(1999-01-08,2001-05-05]"),
				output: []byte(`[1999-01-09,2001-05-06)`)},
			{
				input:  []byte("[,2001-05-05)"),
				output: []byte(`(,2001-05-05)`)},
			{
				input:  []byte("[1999-01-08,]"),
				output: []byte(`[1999-01-08,)`)},
			{
				input:  []byte("[,]"),
				output: []byte(`(,)`)},
		},
	}}.execute(t, "daterange")
}
