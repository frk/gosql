package pgsql

import (
	"testing"
	"time"
)

func TestDateArray(t *testing.T) {
	testlist2{{
		valuer:  DateArrayFromTimeSlice,
		scanner: DateArrayToTimeSlice,
		data: []testdata{
			{input: []time.Time(nil), output: []time.Time(nil)},
			{input: []time.Time{}, output: []time.Time{}},
			{
				input:  []time.Time{dateval(1999, 1, 8)},
				output: []time.Time{dateval(1999, 1, 8)}},
			{
				input:  []time.Time{dateval(1999, 1, 8), dateval(2001, 5, 5)},
				output: []time.Time{dateval(1999, 1, 8), dateval(2001, 5, 5)}},
			{
				input:  []time.Time{dateval(2020, 3, 28), dateval(2001, 5, 5)},
				output: []time.Time{dateval(2020, 3, 28), dateval(2001, 5, 5)}},
		},
	}, {
		data: []testdata{
			{input: string(`{}`), output: string(`{}`)},
			{input: string(`{1999-01-08}`), output: string(`{1999-01-08}`)},
			{input: string(`{1999-01-08,2001-05-05}`), output: string(`{1999-01-08,2001-05-05}`)},
			{input: string(`{2020-03-28,2001-05-05}`), output: string(`{2020-03-28,2001-05-05}`)},
		},
	}, {
		data: []testdata{
			{input: []byte(`{1999-01-08}`), output: []byte(`{1999-01-08}`)},
			{input: []byte(`{1999-01-08,2001-05-05}`), output: []byte(`{1999-01-08,2001-05-05}`)},
			{input: []byte(`{2020-03-28,2001-05-05}`), output: []byte(`{2020-03-28,2001-05-05}`)},
		},
	}}.execute(t, "datearr")
}
