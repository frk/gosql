package convert

import (
	"net"
	"testing"
)

func TestMACAddrArray(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(MACAddrArrayFromHardwareAddrSlice)
		},
		scanner: func() (interface{}, interface{}) {
			v := MACAddrArrayToHardwareAddrSlice{Val: new([]net.HardwareAddr)}
			return v, v.Val
		},
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
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			v := new(string)
			return v, v
		},
		data: []testdata{
			{
				input:  string("{08:00:2b:01:02:03}"),
				output: string(`{08:00:2b:01:02:03}`)},
			{
				input:  string("{08:00:2b:01:02:03,00:00:5e:00:53:01}"),
				output: string(`{08:00:2b:01:02:03,00:00:5e:00:53:01}`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			v := new([]byte)
			return v, v
		},
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
