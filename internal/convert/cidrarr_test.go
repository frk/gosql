package convert

import (
	"net"
	"testing"
)

func TestCIDRArray(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(CIDRArrayFromIPNetSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := CIDRArrayToIPNetSlice{Val: new([]net.IPNet)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  nil,
				output: new([]net.IPNet)},
			{
				input:  netCIDRSlice("192.168.100.128/25"),
				output: netCIDRSliceptr("192.168.100.128/25")},
			{
				input:  netCIDRSlice("192.168.100.128/25", "128.1.0.0/16"),
				output: netCIDRSliceptr("192.168.100.128/25", "128.1.0.0/16")},
			{
				input:  netCIDRSlice("2001:4f8:3:ba::/64", "128.1.0.0/16"),
				output: netCIDRSliceptr("2001:4f8:3:ba::/64", "128.1.0.0/16")},
			{
				input:  netCIDRSlice("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128", "2001:4f8:3:ba::/64"),
				output: netCIDRSliceptr("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128", "2001:4f8:3:ba::/64")},
		},
	}, {
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			d := new(string)
			return d, d
		},
		data: []testdata{
			{
				input:  string("{192.168.100.128/25}"),
				output: strptr(`{192.168.100.128/25}`)},
			{
				input:  string("{192.168.100.128/25,128.1.0.0/16}"),
				output: strptr(`{192.168.100.128/25,128.1.0.0/16}`)},
			{
				input:  string("{2001:4f8:3:ba::/64,128.1.0.0/16}"),
				output: strptr(`{2001:4f8:3:ba::/64,128.1.0.0/16}`)},
			{
				input:  string("{2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128,2001:4f8:3:ba::/64}"),
				output: strptr(`{2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128,2001:4f8:3:ba::/64}`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			d := new([]byte)
			return d, d
		},
		data: []testdata{
			{
				input:  nil,
				output: new([]byte)},
			{
				input:  []byte("{192.168.100.128/25}"),
				output: bytesptr(`{192.168.100.128/25}`)},
			{
				input:  []byte("{192.168.100.128/25,128.1.0.0/16}"),
				output: bytesptr(`{192.168.100.128/25,128.1.0.0/16}`)},
			{
				input:  []byte("{2001:4f8:3:ba::/64,128.1.0.0/16}"),
				output: bytesptr(`{2001:4f8:3:ba::/64,128.1.0.0/16}`)},
			{
				input:  []byte("{2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128,2001:4f8:3:ba::/64}"),
				output: bytesptr(`{2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128,2001:4f8:3:ba::/64}`)},
		},
	}}.execute(t, "cidrarr")
}
