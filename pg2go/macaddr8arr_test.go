package pg2go

import (
	"net"
	"testing"
)

func TestMACAddr8Array(t *testing.T) {
	testlist2{{
		valuer:  MACAddr8ArrayFromHardwareAddrSlice,
		scanner: MACAddr8ArrayToHardwareAddrSlice,
		data: []testdata{
			{input: []net.HardwareAddr(nil), output: []net.HardwareAddr(nil)},
			{input: []net.HardwareAddr{}, output: []net.HardwareAddr{}},
			{
				input:  netMACSlice("08:00:2b:01:02:03:04:05"),
				output: netMACSlice("08:00:2b:01:02:03:04:05")},
			{
				input:  netMACSlice("08:00:2b:01:02:03:04:05", "02:00:5e:10:00:00:00:01"),
				output: netMACSlice("08:00:2b:01:02:03:04:05", "02:00:5e:10:00:00:00:01")},
		},
	}, {
		data: []testdata{
			{
				input:  string("{08:00:2b:01:02:03:04:05}"),
				output: string(`{08:00:2b:01:02:03:04:05}`)},
			{
				input:  string("{08:00:2b:01:02:03:04:05,02:00:5e:10:00:00:00:01}"),
				output: string(`{08:00:2b:01:02:03:04:05,02:00:5e:10:00:00:00:01}`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte("{08:00:2b:01:02:03:04:05}"),
				output: []byte(`{08:00:2b:01:02:03:04:05}`)},
			{
				input:  []byte("{08:00:2b:01:02:03:04:05,02:00:5e:10:00:00:00:01}"),
				output: []byte(`{08:00:2b:01:02:03:04:05,02:00:5e:10:00:00:00:01}`)},
		},
	}}.execute(t, "macaddr8arr")
}
