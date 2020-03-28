package convert

import (
	"net"
	"testing"
)

func TestCIDRArray_Valuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return new(CIDRArrayFromIPNetSlice)
		},
		rows: []test_valuer_row{
			{typ: "cidrarr", in: nil, want: nil},
			{typ: "cidrarr", in: cidrIPNetSlice("192.168.100.128/25"), want: strptr(`{192.168.100.128/25}`)},
			{
				typ:  "cidrarr",
				in:   cidrIPNetSlice("192.168.100.128/25", "128.1.0.0/16"),
				want: strptr(`{192.168.100.128/25,128.1.0.0/16}`)},
			{
				typ:  "cidrarr",
				in:   cidrIPNetSlice("2001:4f8:3:ba::/64", "128.1.0.0/16"),
				want: strptr(`{2001:4f8:3:ba::/64,128.1.0.0/16}`)},
			{
				typ:  "cidrarr",
				in:   cidrIPNetSlice("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128", "2001:4f8:3:ba::/64"),
				want: strptr(`{2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128,2001:4f8:3:ba::/64}`)},
		},
	}}.execute(t)
}

func TestCIDRArray_NoValuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return nil // string
		},
		rows: []test_valuer_row{
			{typ: "cidrarr", in: nil, want: nil},
			{typ: "cidrarr", in: "{192.168.100.128/25}", want: strptr(`{192.168.100.128/25}`)},
			{
				typ:  "cidrarr",
				in:   "{192.168.100.128/25,128.1.0.0/16}",
				want: strptr(`{192.168.100.128/25,128.1.0.0/16}`)},
			{
				typ:  "cidrarr",
				in:   "{2001:4f8:3:ba::/64,128.1.0.0/16}",
				want: strptr(`{2001:4f8:3:ba::/64,128.1.0.0/16}`)},
			{
				typ:  "cidrarr",
				in:   "{2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128,2001:4f8:3:ba::/64}",
				want: strptr(`{2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128,2001:4f8:3:ba::/64}`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		rows: []test_valuer_row{
			{typ: "cidrarr", in: nil, want: nil},
			{typ: "cidrarr", in: []byte("{192.168.100.128/25}"), want: strptr(`{192.168.100.128/25}`)},
			{
				typ:  "cidrarr",
				in:   []byte("{192.168.100.128/25,128.1.0.0/16}"),
				want: strptr(`{192.168.100.128/25,128.1.0.0/16}`)},
			{
				typ:  "cidrarr",
				in:   []byte("{2001:4f8:3:ba::/64,128.1.0.0/16}"),
				want: strptr(`{2001:4f8:3:ba::/64,128.1.0.0/16}`)},
			{
				typ:  "cidrarr",
				in:   []byte("{2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128,2001:4f8:3:ba::/64}"),
				want: strptr(`{2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128,2001:4f8:3:ba::/64}`)},
		},
	}}.execute(t)
}

func TestCIDRArray_Scanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			s := CIDRArrayToIPNetSlice{V: new([]net.IPNet)}
			return s, s.V
		},
		rows: []test_scanner_row{
			{typ: "cidrarr", in: nil, want: new([]net.IPNet)},
			{typ: "cidrarr", in: `{192.168.100.128/25}`, want: cidrIPNetSlicep("192.168.100.128/25")},
			{
				typ:  "cidrarr",
				in:   `{192.168.100.128/25,128.1.0.0/16}`,
				want: cidrIPNetSlicep("192.168.100.128/25", "128.1.0.0/16")},
			{
				typ:  "cidrarr",
				in:   `{2001:4f8:3:ba::/64,128.1.0.0/16}`,
				want: cidrIPNetSlicep("2001:4f8:3:ba::/64", "128.1.0.0/16")},
			{
				typ:  "cidrarr",
				in:   `{2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128,2001:4f8:3:ba::/64}`,
				want: cidrIPNetSlicep("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128", "2001:4f8:3:ba::/64")},
		},
	}}.execute(t)
}

func TestCIDRArray_NoScanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			d := new(string)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "cidrarr", in: `{192.168.100.128/25}`, want: strptr("{192.168.100.128/25}")},
			{
				typ:  "cidrarr",
				in:   `{192.168.100.128/25,128.1.0.0/16}`,
				want: strptr("{192.168.100.128/25,128.1.0.0/16}")},
			{
				typ:  "cidrarr",
				in:   `{2001:4f8:3:ba::/64,128.1.0.0/16}`,
				want: strptr("{2001:4f8:3:ba::/64,128.1.0.0/16}")},
			{
				typ:  "cidrarr",
				in:   `{2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128,2001:4f8:3:ba::/64}`,
				want: strptr("{2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128,2001:4f8:3:ba::/64}")},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			d := new([]byte)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "cidrarr", in: `{192.168.100.128/25}`, want: bytesptr("{192.168.100.128/25}")},
			{
				typ:  "cidrarr",
				in:   `{192.168.100.128/25,128.1.0.0/16}`,
				want: bytesptr("{192.168.100.128/25,128.1.0.0/16}")},
			{
				typ:  "cidrarr",
				in:   `{2001:4f8:3:ba::/64,128.1.0.0/16}`,
				want: bytesptr("{2001:4f8:3:ba::/64,128.1.0.0/16}")},
			{
				typ:  "cidrarr",
				in:   `{2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128,2001:4f8:3:ba::/64}`,
				want: bytesptr("{2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128,2001:4f8:3:ba::/64}")},
		},
	}}.execute(t)
}
