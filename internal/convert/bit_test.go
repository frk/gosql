package convert

import (
	"testing"
)

func TestBit_Valuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return new(BitFromBool)
		},
		rows: []test_valuer_row{
			{typ: "bit", in: bool(true), want: strptr(`1`)},
			{typ: "bit", in: bool(false), want: strptr(`0`)},
		},
	}}.execute(t)
}

func TestBit_NoScanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			d := new(byte)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "bit", in: `0`, want: byteptr(0)},
			{typ: "bit", in: `1`, want: byteptr(1)},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			d := new(bool)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "bit", in: `1`, want: boolptr(true)},
			{typ: "bit", in: `0`, want: boolptr(false)},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			d := new(string)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "bit", in: `1`, want: strptr(`1`)},
			{typ: "bit", in: `0`, want: strptr(`0`)},
		},
	}}.execute(t)
}

func TestBit_NoValuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return nil
		},
		rows: []test_valuer_row{
			{typ: "bit", in: byte(1), want: strptr(`1`)},
			{typ: "bit", in: byte(0), want: strptr(`0`)},
			{typ: "bit", in: string(`1`), want: strptr(`1`)},
			{typ: "bit", in: string(`0`), want: strptr(`0`)},
		},
	}}.execute(t)
}
