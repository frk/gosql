package pg2go

import (
	"testing"
)

func TestMACAddr(t *testing.T) {
	testlist2{{
		valuer:  MACAddrFromHardwareAddr,
		scanner: MACAddrToHardwareAddr,
		data: []testdata{
			{
				input:  netMAC(`08:00:2b:01:02:03`),
				output: netMAC(`08:00:2b:01:02:03`)},
			{
				input:  netMAC(`00:00:5e:00:53:01`),
				output: netMAC(`00:00:5e:00:53:01`)},
		},
	}, {
		data: []testdata{
			{input: string(`08:00:2b:01:02:03`), output: string(`08:00:2b:01:02:03`)},
			{input: string(`00:00:5e:00:53:01`), output: string(`00:00:5e:00:53:01`)},
		},
	}, {
		data: []testdata{
			{input: []byte(`08:00:2b:01:02:03`), output: []byte(`08:00:2b:01:02:03`)},
			{input: []byte(`00:00:5e:00:53:01`), output: []byte(`00:00:5e:00:53:01`)},
		},
	}}.execute(t, "macaddr")
}
