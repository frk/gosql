package convert

import (
	"net"
	"testing"
)

func TestCIDR_Valuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return new(CIDRFromIPNet)
		},
		rows: []test_valuer_row{
			{typ: "cidr", in: nil, want: nil},
			{typ: "cidr", in: cidrIPNet("192.168.100.128/25"), want: strptr(`192.168.100.128/25`)},
			{typ: "cidr", in: cidrIPNet("128.1.0.0/16"), want: strptr(`128.1.0.0/16`)},
			{typ: "cidr", in: cidrIPNet("2001:4f8:3:ba::/64"), want: strptr(`2001:4f8:3:ba::/64`)},
			{
				typ:  "cidr",
				in:   cidrIPNet("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128"),
				want: strptr(`2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128`)},
		},
	}}.execute(t)
}

func TestCIDR_NoValuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return nil // string
		},
		rows: []test_valuer_row{
			{typ: "cidr", in: nil, want: nil},
			{typ: "cidr", in: "192.168.100.128/25", want: strptr(`192.168.100.128/25`)},
			{typ: "cidr", in: "128.1.0.0/16", want: strptr(`128.1.0.0/16`)},
			{typ: "cidr", in: "2001:4f8:3:ba::/64", want: strptr(`2001:4f8:3:ba::/64`)},
			{
				typ:  "cidr",
				in:   "2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128",
				want: strptr(`2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		rows: []test_valuer_row{
			{typ: "cidr", in: nil, want: nil},
			{typ: "cidr", in: []byte("192.168.100.128/25"), want: strptr(`192.168.100.128/25`)},
			{typ: "cidr", in: []byte("128.1.0.0/16"), want: strptr(`128.1.0.0/16`)},
			{typ: "cidr", in: []byte("2001:4f8:3:ba::/64"), want: strptr(`2001:4f8:3:ba::/64`)},
			{
				typ:  "cidr",
				in:   []byte("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128"),
				want: strptr(`2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128`)},
		},
	}}.execute(t)
}

func TestCIDR_Scanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			s := CIDRToIPNet{V: new(net.IPNet)}
			return s, s.V
		},
		rows: []test_scanner_row{
			{typ: "cidr", in: nil, want: new(net.IPNet)},
			{typ: "cidr", in: `192.168.100.128/25`, want: cidrIPNetp(`192.168.100.128/25`)},
			{typ: "cidr", in: `128.1.0.0/16`, want: cidrIPNetp(`128.1.0.0/16`)},
			{typ: "cidr", in: `2001:4f8:3:ba::/64`, want: cidrIPNetp(`2001:4f8:3:ba::/64`)},
			{
				typ:  "cidr",
				in:   `2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128`,
				want: cidrIPNetp(`2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128`)},
		},
	}}.execute(t)
}

func TestCIDR_NoScanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			d := new(string)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "cidr", in: `192.168.100.128/25`, want: strptr("192.168.100.128/25")},
			{typ: "cidr", in: `128.1.0.0/16`, want: strptr("128.1.0.0/16")},
			{typ: "cidr", in: `2001:4f8:3:ba::/64`, want: strptr("2001:4f8:3:ba::/64")},
			{
				typ:  "cidr",
				in:   `2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128`,
				want: strptr("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128")},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			d := new([]byte)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "cidr", in: `192.168.100.128/25`, want: bytesptr("192.168.100.128/25")},
			{typ: "cidr", in: `128.1.0.0/16`, want: bytesptr("128.1.0.0/16")},
			{typ: "cidr", in: `2001:4f8:3:ba::/64`, want: bytesptr("2001:4f8:3:ba::/64")},
			{
				typ:  "cidr",
				in:   `2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128`,
				want: bytesptr("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128")},
		},
	}}.execute(t)
}
