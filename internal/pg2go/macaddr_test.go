package pg2go

import (
	"net"
	"testing"
)

func TestMACAddr(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(MACAddrFromHardwareAddr)
		},
		scanner: func() (interface{}, interface{}) {
			v := MACAddrToHardwareAddr{Val: new(net.HardwareAddr)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  netMAC(`08:00:2b:01:02:03`),
				output: netMAC(`08:00:2b:01:02:03`)},
			{
				input:  netMAC(`00:00:5e:00:53:01`),
				output: netMAC(`00:00:5e:00:53:01`)},
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
			{input: string(`08:00:2b:01:02:03`), output: string(`08:00:2b:01:02:03`)},
			{input: string(`00:00:5e:00:53:01`), output: string(`00:00:5e:00:53:01`)},
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
			{input: []byte(`08:00:2b:01:02:03`), output: []byte(`08:00:2b:01:02:03`)},
			{input: []byte(`00:00:5e:00:53:01`), output: []byte(`00:00:5e:00:53:01`)},
		},
	}}.execute(t, "macaddr")
}
