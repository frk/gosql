package convert

import (
	"testing"
)

func TestBitArray_Valuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return new(BitArrayFromBoolSlice)
		},
		rows: []test_valuer_row{
			{typ: "bitarr", in: nil, want: nil},
			{typ: "bitarr", in: []bool{}, want: strptr(`{}`)},
			{typ: "bitarr", in: []bool{true}, want: strptr(`{1}`)},
			{typ: "bitarr", in: []bool{false}, want: strptr(`{0}`)},
			{
				typ:  "bitarr",
				in:   []bool{false, false, false, true, true, false, true, true},
				want: strptr(`{0,0,0,1,1,0,1,1}`)},
		},
	}, {
		valuer: func() interface{} {
			return new(BitArrayFromUint8Slice)
		},
		rows: []test_valuer_row{
			{typ: "bitarr", in: nil, want: nil},
			{typ: "bitarr", in: []uint8{}, want: strptr(`{}`)},
			{typ: "bitarr", in: []uint8{1}, want: strptr(`{1}`)},
			{typ: "bitarr", in: []uint8{0}, want: strptr(`{0}`)},
			{
				typ:  "bitarr",
				in:   []uint8{0, 0, 0, 1, 1, 0, 1, 1},
				want: strptr(`{0,0,0,1,1,0,1,1}`)},
		},
	}, {
		valuer: func() interface{} {
			return new(BitArrayFromUintSlice)
		},
		rows: []test_valuer_row{
			{typ: "bitarr", in: nil, want: nil},
			{typ: "bitarr", in: []uint{}, want: strptr(`{}`)},
			{typ: "bitarr", in: []uint{1}, want: strptr(`{1}`)},
			{typ: "bitarr", in: []uint{0}, want: strptr(`{0}`)},
			{
				typ:  "bitarr",
				in:   []uint{0, 0, 0, 1, 1, 0, 1, 1},
				want: strptr(`{0,0,0,1,1,0,1,1}`)},
		},
	}}.execute(t)
}

func TestBitArray_Scanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			s := BitArrayToBoolSlice{S: new([]bool)}
			return s, s.S
		},
		rows: []test_scanner_row{
			{typ: "bitarr", in: nil, want: new([]bool)},
			{typ: "bitarr", in: `{}`, want: &[]bool{}},
			{typ: "bitarr", in: `{0}`, want: &[]bool{false}},
			{typ: "bitarr", in: `{1}`, want: &[]bool{true}},
			{typ: "bitarr", in: `{1,0}`, want: &[]bool{true, false}},
			{
				typ:  "bitarr",
				in:   `{0,1,1,1,0,1,0,0}`,
				want: &[]bool{false, true, true, true, false, true, false, false}},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			s := BitArrayToUint8Slice{S: new([]uint8)}
			return s, s.S
		},
		rows: []test_scanner_row{
			{typ: "bitarr", in: nil, want: new([]uint8)},
			{typ: "bitarr", in: `{}`, want: &[]uint8{}},
			{typ: "bitarr", in: `{0}`, want: &[]uint8{0}},
			{typ: "bitarr", in: `{1}`, want: &[]uint8{1}},
			{typ: "bitarr", in: `{1,0}`, want: &[]uint8{1, 0}},
			{
				typ:  "bitarr",
				in:   `{0,1,1,1,0,1,0,0}`,
				want: &[]uint8{0, 1, 1, 1, 0, 1, 0, 0}},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			s := BitArrayToUintSlice{S: new([]uint)}
			return s, s.S
		},
		rows: []test_scanner_row{
			{typ: "bitarr", in: nil, want: new([]uint)},
			{typ: "bitarr", in: `{}`, want: &[]uint{}},
			{typ: "bitarr", in: `{0}`, want: &[]uint{0}},
			{typ: "bitarr", in: `{1}`, want: &[]uint{1}},
			{typ: "bitarr", in: `{1,0}`, want: &[]uint{1, 0}},
			{
				typ:  "bitarr",
				in:   `{0,1,1,1,0,1,0,0}`,
				want: &[]uint{0, 1, 1, 1, 0, 1, 0, 0}},
		},
	}}.execute(t)
}

func TestBitArray_NoValuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return nil // string
		},
		rows: []test_valuer_row{
			{typ: "bitarr", in: nil, want: nil},
			{typ: "bitarr", in: "{}", want: strptr(`{}`)},
			{typ: "bitarr", in: "{1,0}", want: strptr(`{1,0}`)},
			{typ: "bitarr", in: "{0,1}", want: strptr(`{0,1}`)},
			{typ: "bitarr", in: "{0,1,1,1,0,1,0,0}", want: strptr(`{0,1,1,1,0,1,0,0}`)},
		},
	}, {
		valuer: func() interface{} {
			return nil
		},
		rows: []test_valuer_row{
			{typ: "bitarr", in: nil, want: nil},
			{typ: "bitarr", in: []byte("{}"), want: strptr(`{}`)},
			{typ: "bitarr", in: []byte("{1,0}"), want: strptr(`{1,0}`)},
			{typ: "bitarr", in: []byte("{0,1}"), want: strptr(`{0,1}`)},
			{typ: "bitarr", in: []byte("{0,1,1,1,0,1,0,0}"), want: strptr(`{0,1,1,1,0,1,0,0}`)},
		},
	}}.execute(t)
}

func TestBitArray_NoScanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			d := new(string)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "bitarr", in: `{}`, want: strptr("{}")},
			{typ: "bitarr", in: `{1,0}`, want: strptr("{1,0}")},
			{typ: "bitarr", in: `{0,1}`, want: strptr("{0,1}")},
			{typ: "bitarr", in: `{0,1,1,1,0,1,0,0}`, want: strptr("{0,1,1,1,0,1,0,0}")},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			d := new([]byte)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "bitarr", in: `{}`, want: bytesptr("{}")},
			{typ: "bitarr", in: `{1,0}`, want: bytesptr("{1,0}")},
			{typ: "bitarr", in: `{0,1}`, want: bytesptr("{0,1}")},
			{typ: "bitarr", in: `{0,1,1,1,0,1,0,0}`, want: bytesptr("{0,1,1,1,0,1,0,0}")},
		},
	}}.execute(t)
}
