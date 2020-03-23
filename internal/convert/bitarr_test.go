package convert

import (
	"testing"
)

func TestBitArr_Scanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			s := BitArrToBoolSlice{Ptr: new([]bool)}
			return s, s.Ptr
		},
		rows: []test_scanner_row{
			{typ: "bitarr", in: nil, want: new([]bool)},
			{typ: "bitarr", in: `{}`, want: &[]bool{}},
			{typ: "bitarr", in: `{0}`, want: &[]bool{false}},
			{typ: "bitarr", in: `{1}`, want: &[]bool{true}},
			{typ: "bitarr", in: `{1,0}`, want: &[]bool{true, false}},
			{typ: "bitarr", in: `{0,1,1,1,0,1,0,0}`, want: &[]bool{false, true, true, true, false, true, false, false}},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			s := BitArrToUint8Slice{Ptr: new([]uint8)}
			return s, s.Ptr
		},
		rows: []test_scanner_row{
			{typ: "bitarr", in: nil, want: new([]uint8)},
			{typ: "bitarr", in: `{}`, want: &[]uint8{}},
			{typ: "bitarr", in: `{0}`, want: &[]uint8{0}},
			{typ: "bitarr", in: `{1}`, want: &[]uint8{1}},
			{typ: "bitarr", in: `{1,0}`, want: &[]uint8{1, 0}},
			{typ: "bitarr", in: `{0,1,1,1,0,1,0,0}`, want: &[]uint8{0, 1, 1, 1, 0, 1, 0, 0}},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			s := BitArrToUintSlice{Ptr: new([]uint)}
			return s, s.Ptr
		},
		rows: []test_scanner_row{
			{typ: "bitarr", in: nil, want: new([]uint)},
			{typ: "bitarr", in: `{}`, want: &[]uint{}},
			{typ: "bitarr", in: `{0}`, want: &[]uint{0}},
			{typ: "bitarr", in: `{1}`, want: &[]uint{1}},
			{typ: "bitarr", in: `{1,0}`, want: &[]uint{1, 0}},
			{typ: "bitarr", in: `{0,1,1,1,0,1,0,0}`, want: &[]uint{0, 1, 1, 1, 0, 1, 0, 0}},
		},
	}}.execute(t)
}

func TestBitArr_NoScanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			d := new([]byte)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "bitarr", in: nil, want: new([]byte)},
			{typ: "bitarr", in: `{}`, want: bytesptr(`{}`)},
			{typ: "bitarr", in: `{0}`, want: bytesptr(`{0}`)},
			{typ: "bitarr", in: `{1}`, want: bytesptr(`{1}`)},
			{typ: "bitarr", in: `{1,0}`, want: bytesptr(`{1,0}`)},
			{typ: "bitarr", in: `{0,1,1,1,0,1,0,0}`, want: bytesptr(`{0,1,1,1,0,1,0,0}`)},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			d := new(string)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "bitarr", in: `{}`, want: strptr(`{}`)},
			{typ: "bitarr", in: `{0}`, want: strptr(`{0}`)},
			{typ: "bitarr", in: `{1}`, want: strptr(`{1}`)},
			{typ: "bitarr", in: `{1,0}`, want: strptr(`{1,0}`)},
			{typ: "bitarr", in: `{0,1,1,1,0,1,0,0}`, want: strptr(`{0,1,1,1,0,1,0,0}`)},
		},
	}}.execute(t)
}
