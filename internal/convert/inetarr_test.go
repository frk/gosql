package convert

import (
	"net"
	"testing"
)

func TestInetArray(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(InetArrayFromIPNetSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := InetArrayToIPNetSlice{Val: new([]net.IPNet)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  nil,
				output: []net.IPNet(nil)},
			{
				input:  netCIDRSlice("192.168.100.128/25"),
				output: netCIDRSlice("192.168.100.128/25")},
			{
				input:  netCIDRSlice("192.168.100.128/25", "128.1.0.0/16"),
				output: netCIDRSlice("192.168.100.128/25", "128.1.0.0/16")},
			{
				input:  netCIDRSlice("2001:4f8:3:ba::/64", "128.1.0.0/16"),
				output: netCIDRSlice("2001:4f8:3:ba::/64", "128.1.0.0/16")},
			{
				input:  netCIDRSlice("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128", "2001:4f8:3:ba::/64"),
				output: []net.IPNet{{IP: netIP("2001:4f8:3:ba:2e0:81ff:fe22:d1f1")}, netCIDR("2001:4f8:3:ba::/64")}},
		},
	}, {
		valuer: func() interface{} {
			return new(InetArrayFromIPSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := InetArrayToIPSlice{Val: new([]net.IP)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  nil,
				output: []net.IP(nil)},
			{
				input:  netIPSlice("192.168.100.128"),
				output: netIPSlice("192.168.100.128")},
			{
				input:  netIPSlice("192.168.100.128", "128.1.0.0"),
				output: netIPSlice("192.168.100.128", "128.1.0.0")},
			{
				input:  netIPSlice("2001:4f8:3:ba::", "128.1.0.0"),
				output: netIPSlice("2001:4f8:3:ba::", "128.1.0.0")},
			{
				input:  netIPSlice("2001:4f8:3:ba:2e0:81ff:fe22:d1f1", "2001:4f8:3:ba::"),
				output: netIPSlice("2001:4f8:3:ba:2e0:81ff:fe22:d1f1", "2001:4f8:3:ba::")},
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
				input:  string("{192.168.100.128}"),
				output: string(`{192.168.100.128}`)},
			{
				input:  string("{192.168.100.128,128.1.0.0}"),
				output: string(`{192.168.100.128,128.1.0.0}`)},
			{
				input:  string("{2001:4f8:3:ba::,128.1.0.0}"),
				output: string(`{2001:4f8:3:ba::,128.1.0.0}`)},
			{
				input:  string("{2001:4f8:3:ba:2e0:81ff:fe22:d1f1,2001:4f8:3:ba::}"),
				output: string(`{2001:4f8:3:ba:2e0:81ff:fe22:d1f1,2001:4f8:3:ba::}`)},
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
				input:  []byte("{192.168.100.128}"),
				output: []byte(`{192.168.100.128}`)},
			{
				input:  []byte("{192.168.100.128,128.1.0.0}"),
				output: []byte(`{192.168.100.128,128.1.0.0}`)},
			{
				input:  []byte("{2001:4f8:3:ba::,128.1.0.0}"),
				output: []byte(`{2001:4f8:3:ba::,128.1.0.0}`)},
			{
				input:  []byte("{2001:4f8:3:ba:2e0:81ff:fe22:d1f1,2001:4f8:3:ba::}"),
				output: []byte(`{2001:4f8:3:ba:2e0:81ff:fe22:d1f1,2001:4f8:3:ba::}`)},
		},
	}}.execute(t, "inetarr")
}
