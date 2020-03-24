package convert

import (
	"testing"
)

func TestBPChar_Valuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return new(BPCharFromByte)
		},
		rows: []test_valuer_row{
			{typ: "bpchar", in: byte('A'), want: strptr(`A`)},
			{typ: "bpchar", in: byte('b'), want: strptr(`b`)},
			{typ: "bpchar", in: byte('X'), want: strptr(`X`)},
		},
	}, {
		valuer: func() interface{} {
			return new(BPCharFromRune)
		},
		rows: []test_valuer_row{
			{typ: "bpchar", in: rune('A'), want: strptr(`A`)},
			{typ: "bpchar", in: rune('馬'), want: strptr(`馬`)},
			{typ: "bpchar", in: rune('駮'), want: strptr(`駮`)},
		},
	}}.execute(t)
}

func TestBPChar_Scanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			s := &BPCharToByte{B: new(byte)}
			return s, s.B
		},
		rows: []test_scanner_row{
			{typ: "bpchar", in: `a`, want: byteptr('a')},
			{typ: "bpchar", in: `Z`, want: byteptr('Z')},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			s := &BPCharToRune{R: new(rune)}
			return s, s.R
		},
		rows: []test_scanner_row{
			{typ: "bpchar", in: `a`, want: runeptr('a')},
			{typ: "bpchar", in: `Z`, want: runeptr('Z')},
			{typ: "bpchar", in: `馬`, want: runeptr('馬')},
			{typ: "bpchar", in: `駮`, want: runeptr('駮')},
		},
	}}.execute(t)
}

func TestBPChar_NoValuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return nil
		},
		rows: []test_valuer_row{
			{typ: "bpchar", in: `a`, want: strptr(`a`)},
			{typ: "bpchar", in: `馬`, want: strptr(`馬`)},
			{typ: "bpchar", in: []byte{'a'}, want: strptr(`a`)},
			{typ: "bpchar", in: []byte{'Z'}, want: strptr(`Z`)},
		},
	}}.execute(t)
}

func TestBPChar_NoScanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			d := new([]byte)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "bpchar", in: `a`, want: bytesptr(`a`)},
			{typ: "bpchar", in: `Z`, want: bytesptr(`Z`)},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			d := new(string)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "bpchar", in: `a`, want: strptr(`a`)},
			{typ: "bpchar", in: `馬`, want: strptr(`馬`)},
		},
	}}.execute(t)
}
