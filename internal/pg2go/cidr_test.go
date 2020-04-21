package pg2go

import (
	"net"
	"testing"
)

func TestCIDR(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(CIDRFromIPNet)
		},
		scanner: func() (interface{}, interface{}) {
			s := CIDRToIPNet{Val: new(net.IPNet)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  nil,
				output: net.IPNet{}},
			{
				input:  netCIDR("192.168.100.128/25"),
				output: netCIDR(`192.168.100.128/25`)},
			{
				input:  netCIDR("128.1.0.0/16"),
				output: netCIDR(`128.1.0.0/16`)},
			{
				input:  netCIDR("2001:4f8:3:ba::/64"),
				output: netCIDR(`2001:4f8:3:ba::/64`)},
			{
				input:  netCIDR("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128"),
				output: netCIDR(`2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128`)},
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
				output: string(`192.168.100.128/25`)},
			{
				input:  string("128.1.0.0/16"),
				output: string(`128.1.0.0/16`)},
			{
				input:  string("2001:4f8:3:ba::/64"),
				output: string(`2001:4f8:3:ba::/64`)},
			{
				input:  string("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128"),
				output: string(`2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128`)},
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
				output: []byte(nil)},
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
				output: []byte(`2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128`)},
		},
	}}.execute(t, "cidr")
}
