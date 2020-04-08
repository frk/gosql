package convert

import (
	"net"
	"testing"
)

func TestInet_ValuerAndScanner(t *testing.T) {
	test_valuer_with_scanner{{
		valuer: func() interface{} {
			return new(InetFromIPNet)
		},
		scanner: func() (interface{}, interface{}) {
			s := InetToIPNet{Val: new(net.IPNet)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "inet",
				input:  nil,
				output: new(net.IPNet)},
			{
				typ:    "inet",
				input:  netCIDR("192.168.100.128/25"),
				output: netCIDRptr("192.168.100.128/25")},
			{
				typ:    "inet",
				input:  netCIDR("128.1.0.0/16"),
				output: netCIDRptr("128.1.0.0/16")},
			{
				typ:    "inet",
				input:  netCIDR("2001:4f8:3:ba::/64"),
				output: netCIDRptr("2001:4f8:3:ba::/64")},
			{
				typ:    "inet",
				input:  netCIDR("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/64"),
				output: netCIDRptr("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/64")},
			{
				typ:    "inet",
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
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "inet",
				input:  nil,
				output: new(net.IP)},
			{
				typ:    "inet",
				input:  netIP("192.168.100.128"),
				output: netIPptr("192.168.100.128")},
			{
				typ:    "inet",
				input:  netIP("128.1.0.0"),
				output: netIPptr("128.1.0.0")},
			{
				typ:    "inet",
				input:  netIP("2001:4f8:3:ba:2e0:81ff:fe22:d1f1"),
				output: netIPptr("2001:4f8:3:ba:2e0:81ff:fe22:d1f1")},
		},
	}}.execute(t)
}

func TestInet_NoValuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return nil // string
		},
		rows: []test_valuer_row{
			{
				typ:  "inet",
				in:   nil,
				want: nil},
			{
				typ:  "inet",
				in:   "192.168.100.128/25",
				want: strptr(`192.168.100.128/25`)},
			{
				typ:  "inet",
				in:   "128.1.0.0/16",
				want: strptr(`128.1.0.0/16`)},
			{
				typ:  "inet",
				in:   "2001:4f8:3:ba::/64",
				want: strptr(`2001:4f8:3:ba::/64`)},
			{
				typ:  "inet",
				in:   "2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128",
				want: strptr(`2001:4f8:3:ba:2e0:81ff:fe22:d1f1`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		rows: []test_valuer_row{
			{
				typ:  "inet",
				in:   nil,
				want: nil},
			{
				typ:  "inet",
				in:   []byte("192.168.100.128/25"),
				want: strptr(`192.168.100.128/25`)},
			{
				typ:  "inet",
				in:   []byte("128.1.0.0/16"),
				want: strptr(`128.1.0.0/16`)},
			{
				typ:  "inet",
				in:   []byte("2001:4f8:3:ba::/64"),
				want: strptr(`2001:4f8:3:ba::/64`)},
			{
				typ:  "inet",
				in:   []byte("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128"),
				want: strptr(`2001:4f8:3:ba:2e0:81ff:fe22:d1f1`)},
		},
	}}.execute(t)
}

func TestInet_NoScanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			d := new(string)
			return d, d
		},
		rows: []test_scanner_row{
			{
				typ:  "inet",
				in:   `192.168.100.128/25`,
				want: strptr("192.168.100.128/25")},
			{
				typ:  "inet",
				in:   `128.1.0.0/16`,
				want: strptr("128.1.0.0/16")},
			{
				typ:  "inet",
				in:   `2001:4f8:3:ba::/64`,
				want: strptr("2001:4f8:3:ba::/64")},
			{
				typ:  "inet",
				in:   `2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128`,
				want: strptr("2001:4f8:3:ba:2e0:81ff:fe22:d1f1")},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			d := new([]byte)
			return d, d
		},
		rows: []test_scanner_row{
			{
				typ:  "inet",
				in:   `192.168.100.128/25`,
				want: bytesptr("192.168.100.128/25")},
			{
				typ:  "inet",
				in:   `128.1.0.0/16`,
				want: bytesptr("128.1.0.0/16")},
			{
				typ:  "inet",
				in:   `2001:4f8:3:ba::/64`,
				want: bytesptr("2001:4f8:3:ba::/64")},
			{
				typ:  "inet",
				in:   `2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128`,
				want: bytesptr("2001:4f8:3:ba:2e0:81ff:fe22:d1f1")},
		},
	}}.execute(t)
}
