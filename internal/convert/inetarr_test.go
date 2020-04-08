package convert

import (
	"net"
	"testing"
)

func TestInetArray_ValuerAndScanner(t *testing.T) {
	test_valuer_with_scanner{{
		valuer: func() interface{} {
			return new(InetArrayFromIPNetSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := InetArrayToIPNetSlice{Val: new([]net.IPNet)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "inetarr",
				input:  nil,
				output: new([]net.IPNet)},
			{
				typ:    "inetarr",
				input:  netCIDRSlice("192.168.100.128/25"),
				output: netCIDRSliceptr("192.168.100.128/25")},
			{
				typ:    "inetarr",
				input:  netCIDRSlice("192.168.100.128/25", "128.1.0.0/16"),
				output: netCIDRSliceptr("192.168.100.128/25", "128.1.0.0/16")},
			{
				typ:    "inetarr",
				input:  netCIDRSlice("2001:4f8:3:ba::/64", "128.1.0.0/16"),
				output: netCIDRSliceptr("2001:4f8:3:ba::/64", "128.1.0.0/16")},
			{
				typ:    "inetarr",
				input:  netCIDRSlice("2001:4f8:3:ba:2e0:81ff:fe22:d1f1/128", "2001:4f8:3:ba::/64"),
				output: &[]net.IPNet{{IP: netIP("2001:4f8:3:ba:2e0:81ff:fe22:d1f1")}, netCIDR("2001:4f8:3:ba::/64")}},
		},
	}, {
		valuer: func() interface{} {
			return new(InetArrayFromIPSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := InetArrayToIPSlice{Val: new([]net.IP)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "inetarr",
				input:  nil,
				output: new([]net.IP)},
			{
				typ:    "inetarr",
				input:  netIPSlice("192.168.100.128"),
				output: netIPSliceptr("192.168.100.128")},
			{
				typ:    "inetarr",
				input:  netIPSlice("192.168.100.128", "128.1.0.0"),
				output: netIPSliceptr("192.168.100.128", "128.1.0.0")},
			{
				typ:    "inetarr",
				input:  netIPSlice("2001:4f8:3:ba::", "128.1.0.0"),
				output: netIPSliceptr("2001:4f8:3:ba::", "128.1.0.0")},
			{
				typ:    "inetarr",
				input:  netIPSlice("2001:4f8:3:ba:2e0:81ff:fe22:d1f1", "2001:4f8:3:ba::"),
				output: netIPSliceptr("2001:4f8:3:ba:2e0:81ff:fe22:d1f1", "2001:4f8:3:ba::")},
		},
	}}.execute(t)
}

func TestInetArray_NoValuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return nil // string
		},
		rows: []test_valuer_row{
			{
				typ:  "inetarr",
				in:   nil,
				want: nil},
			{
				typ:  "inetarr",
				in:   "{192.168.100.128}",
				want: strptr(`{192.168.100.128}`)},
			{
				typ:  "inetarr",
				in:   "{192.168.100.128,128.1.0.0}",
				want: strptr(`{192.168.100.128,128.1.0.0}`)},
			{
				typ:  "inetarr",
				in:   "{2001:4f8:3:ba::,128.1.0.0}",
				want: strptr(`{2001:4f8:3:ba::,128.1.0.0}`)},
			{
				typ:  "inetarr",
				in:   "{2001:4f8:3:ba:2e0:81ff:fe22:d1f1,2001:4f8:3:ba::}",
				want: strptr(`{2001:4f8:3:ba:2e0:81ff:fe22:d1f1,2001:4f8:3:ba::}`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		rows: []test_valuer_row{
			{
				typ:  "inetarr",
				in:   nil,
				want: nil},
			{
				typ:  "inetarr",
				in:   []byte("{192.168.100.128}"),
				want: strptr(`{192.168.100.128}`)},
			{
				typ:  "inetarr",
				in:   []byte("{192.168.100.128,128.1.0.0}"),
				want: strptr(`{192.168.100.128,128.1.0.0}`)},
			{
				typ:  "inetarr",
				in:   []byte("{2001:4f8:3:ba::,128.1.0.0}"),
				want: strptr(`{2001:4f8:3:ba::,128.1.0.0}`)},
			{
				typ:  "inetarr",
				in:   []byte("{2001:4f8:3:ba:2e0:81ff:fe22:d1f1,2001:4f8:3:ba::}"),
				want: strptr(`{2001:4f8:3:ba:2e0:81ff:fe22:d1f1,2001:4f8:3:ba::}`)},
		},
	}}.execute(t)
}

func TestInetArray_NoScanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			d := new(string)
			return d, d
		},
		rows: []test_scanner_row{
			{
				typ:  "inetarr",
				in:   `{192.168.100.128}`,
				want: strptr("{192.168.100.128}")},
			{
				typ:  "inetarr",
				in:   `{192.168.100.128,128.1.0.0}`,
				want: strptr("{192.168.100.128,128.1.0.0}")},
			{
				typ:  "inetarr",
				in:   `{2001:4f8:3:ba::,128.1.0.0}`,
				want: strptr("{2001:4f8:3:ba::,128.1.0.0}")},
			{
				typ:  "inetarr",
				in:   `{2001:4f8:3:ba:2e0:81ff:fe22:d1f1,2001:4f8:3:ba::}`,
				want: strptr("{2001:4f8:3:ba:2e0:81ff:fe22:d1f1,2001:4f8:3:ba::}")},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			d := new([]byte)
			return d, d
		},
		rows: []test_scanner_row{
			{
				typ:  "inetarr",
				in:   `{192.168.100.128}`,
				want: bytesptr("{192.168.100.128}")},
			{
				typ:  "inetarr",
				in:   `{192.168.100.128,128.1.0.0}`,
				want: bytesptr("{192.168.100.128,128.1.0.0}")},
			{
				typ:  "inetarr",
				in:   `{2001:4f8:3:ba::,128.1.0.0}`,
				want: bytesptr("{2001:4f8:3:ba::,128.1.0.0}")},
			{
				typ:  "inetarr",
				in:   `{2001:4f8:3:ba:2e0:81ff:fe22:d1f1,2001:4f8:3:ba::}`,
				want: bytesptr("{2001:4f8:3:ba:2e0:81ff:fe22:d1f1,2001:4f8:3:ba::}")},
		},
	}}.execute(t)
}
