package pg2go

import (
	"net"
	"testing"
)

func TestMACAddrArray(t *testing.T) {
	testlist2{{
		valuer:  MACAddrArrayFromHardwareAddrSlice,
		scanner: MACAddrArrayToHardwareAddrSlice,
		data: []testdata{
			{input: []net.HardwareAddr(nil), output: []net.HardwareAddr(nil)},
			{input: []net.HardwareAddr{}, output: []net.HardwareAddr{}},
			{
				input:  netMACSlice("08:00:2b:01:02:03"),
				output: netMACSlice("08:00:2b:01:02:03")},
			{
				input:  netMACSlice("08:00:2b:01:02:03", "00:00:5e:00:53:01"),
				output: netMACSlice("08:00:2b:01:02:03", "00:00:5e:00:53:01")},
		},
	}, {
		data: []testdata{
			{
				input:  string("{08:00:2b:01:02:03}"),
				output: string(`{08:00:2b:01:02:03}`)},
			{
				input:  string("{08:00:2b:01:02:03,00:00:5e:00:53:01}"),
				output: string(`{08:00:2b:01:02:03,00:00:5e:00:53:01}`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte("{08:00:2b:01:02:03}"),
				output: []byte(`{08:00:2b:01:02:03}`)},
			{
				input:  []byte("{08:00:2b:01:02:03,00:00:5e:00:53:01}"),
				output: []byte(`{08:00:2b:01:02:03,00:00:5e:00:53:01}`)},
		},
	}}.execute(t, "macaddrarr")
}
