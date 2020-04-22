package pgsql

import (
	"testing"
)

func TestMACAddr8(t *testing.T) {
	testlist2{{
		valuer:  MACAddr8FromHardwareAddr,
		scanner: MACAddr8ToHardwareAddr,
		data: []testdata{
			{
				input:  netMAC(`08:00:2b:01:02:03:04:05`),
				output: netMAC(`08:00:2b:01:02:03:04:05`)},
			{
				input:  netMAC(`02:00:5e:10:00:00:00:01`),
				output: netMAC(`02:00:5e:10:00:00:00:01`)},
		},
	}, {
		data: []testdata{
			{input: string(`08:00:2b:01:02:03:04:05`), output: string(`08:00:2b:01:02:03:04:05`)},
			{input: string(`02:00:5e:10:00:00:00:01`), output: string(`02:00:5e:10:00:00:00:01`)},
		},
	}, {
		data: []testdata{
			{input: []byte(`08:00:2b:01:02:03:04:05`), output: []byte(`08:00:2b:01:02:03:04:05`)},
			{input: []byte(`02:00:5e:10:00:00:00:01`), output: []byte(`02:00:5e:10:00:00:00:01`)},
		},
	}}.execute(t, "macaddr8")
}
