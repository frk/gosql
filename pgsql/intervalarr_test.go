package pgsql

import (
	"testing"
)

func TestIntervalArray(t *testing.T) {
	testlist2{{
		data: []testdata{
			{
				input:  string(`{"1 day","-5 years -4 mons -00:34:00"}`),
				output: string(`{"1 day","-5 years -4 mons -00:34:00"}`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte(`{"1 day","-5 years -4 mons -00:34:00"}`),
				output: []byte(`{"1 day","-5 years -4 mons -00:34:00"}`)},
		},
	}}.execute(t, "intervalarr")
}
