package convert

import (
	"net"
	"testing"
)

func TestInet(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(InetFromIPNet)
		},
		scanner: func() (interface{}, interface{}) {
			s := InetToIPNet{Val: new(net.IPNet)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  nil,
				output: new(net.IPNet)},
			{
				input:  netCIDR("192.168.100.128/25"),
				output: netCIDRptr("192.168.100.128/25")},
			{
				input:  netCIDR("128.1.0.0/16"),
				output: netCIDRptr("128.1.0.0/16")},
			{
				input:  netCIDR("2001:4f8:3:ba::/64"),
				output: netCIDRptr("2001:4f8:3:ba::/64")},
			{
				input:  netCIDR("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/64"),
				output: netCIDRptr("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/64")},
			{
				input:  netCIDR("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128"),
				output: &net.IPNet{IP: netIP("2001:4f8:3:ba:2e0:81ff:fe22:d1f1")}},
		},
	}, {
		valuer: func() interface{} {
			return new(InetFromIP)
		},
		scanner: func() (interface{}, interface{}) {
			s := InetToIP{Val: new(net.IP)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  nil,
				output: new(net.IP)},
			{
				input:  netIP("192.168.100.128"),
				output: netIPptr("192.168.100.128")},
			{
				input:  netIP("128.1.0.0"),
				output: netIPptr("128.1.0.0")},
			{
				input:  netIP("2001:4f8:3:ba:2e0:81ff:fe22:d1f1"),
				output: netIPptr("2001:4f8:3:ba:2e0:81ff:fe22:d1f1")},
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
				input:  string("192.168.100.128/25"),
				output: strptr(`192.168.100.128/25`)},
			{
				input:  string("128.1.0.0/16"),
				output: strptr(`128.1.0.0/16`)},
			{
				input:  string("2001:4f8:3:ba::/64"),
				output: strptr(`2001:4f8:3:ba::/64`)},
			{
				input:  string("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128"),
				output: strptr(`2001:4f8:3:ba:2e0:81ff:fe22:d1f1`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			s := new([]byte)
			return s, s
		},
		data: []testdata{
			{
				input:  []byte("192.168.100.128/25"),
				output: bytesptr(`192.168.100.128/25`)},
			{
				input:  []byte("128.1.0.0/16"),
				output: bytesptr(`128.1.0.0/16`)},
			{
				input:  []byte("2001:4f8:3:ba::/64"),
				output: bytesptr(`2001:4f8:3:ba::/64`)},
			{
				input:  []byte("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128"),
				output: bytesptr(`2001:4f8:3:ba:2e0:81ff:fe22:d1f1`)},
		},
	}}.execute(t, "inet")
}
