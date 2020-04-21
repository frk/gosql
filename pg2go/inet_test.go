package pg2go

import (
	"net"
	"testing"
)

func TestInet(t *testing.T) {
	testlist2{{
		valuer:  InetFromIPNet,
		scanner: InetToIPNet,
		data: []testdata{
			{
				input:  net.IPNet{},
				output: net.IPNet{}},
			{
				input:  netCIDR("192.168.100.128/25"),
				output: netCIDR("192.168.100.128/25")},
			{
				input:  netCIDR("128.1.0.0/16"),
				output: netCIDR("128.1.0.0/16")},
			{
				input:  netCIDR("2001:4f8:3:ba::/64"),
				output: netCIDR("2001:4f8:3:ba::/64")},
			{
				input:  netCIDR("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/64"),
				output: netCIDR("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/64")},
			{
				input:  netCIDR("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128"),
				output: net.IPNet{IP: netIP("2001:4f8:3:ba:2e0:81ff:fe22:d1f1")}},
		},
	}, {
		valuer:  InetFromIP,
		scanner: InetToIP,
		data: []testdata{
			{
				input:  net.IP(nil),
				output: net.IP(nil)},
			{
				input:  netIP("192.168.100.128"),
				output: netIP("192.168.100.128")},
			{
				input:  netIP("128.1.0.0"),
				output: netIP("128.1.0.0")},
			{
				input:  netIP("2001:4f8:3:ba:2e0:81ff:fe22:d1f1"),
				output: netIP("2001:4f8:3:ba:2e0:81ff:fe22:d1f1")},
		},
	}, {
		data: []testdata{
			{
				input:  string("192.168.100.128/25"),
				output: string(`192.168.100.128/25`)},
			{
				input:  string("128.1.0.0/16"),
				output: string(`128.1.0.0/16`)},
			{
				input:  string("2001:4f8:3:ba::/64"),
				output: string(`2001:4f8:3:ba::/64`)},
			{
				input:  string("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128"),
				output: string(`2001:4f8:3:ba:2e0:81ff:fe22:d1f1`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte("192.168.100.128/25"),
				output: []byte(`192.168.100.128/25`)},
			{
				input:  []byte("128.1.0.0/16"),
				output: []byte(`128.1.0.0/16`)},
			{
				input:  []byte("2001:4f8:3:ba::/64"),
				output: []byte(`2001:4f8:3:ba::/64`)},
			{
				input:  []byte("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128"),
				output: []byte(`2001:4f8:3:ba:2e0:81ff:fe22:d1f1`)},
		},
	}}.execute(t, "inet")
}
