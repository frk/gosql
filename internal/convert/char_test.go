package convert

import (
	"testing"
)

func TestChar_Valuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return new(CharFromByte)
		},
		rows: []test_valuer_row{
			{typ: "char", in: byte('A'), want: strptr(`A`)},
			{typ: "char", in: byte('b'), want: strptr(`b`)},
			{typ: "char", in: byte('X'), want: strptr(`X`)},
		},
	}, {
		valuer: func() interface{} {
			return new(CharFromRune)
		},
		rows: []test_valuer_row{
			{typ: "char", in: rune('A'), want: strptr(`A`)},
			{typ: "char", in: rune('馬'), want: strptr(`馬`)},
			{typ: "char", in: rune('駮'), want: strptr(`駮`)},
		},
	}}.execute(t)
}

func TestChar_Scanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			s := &CharToByte{B: new(byte)}
			return s, s.B
		},
		rows: []test_scanner_row{
			{typ: "char", in: `a`, want: byteptr('a')},
			{typ: "char", in: `Z`, want: byteptr('Z')},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			s := &CharToRune{R: new(rune)}
			return s, s.R
		},
		rows: []test_scanner_row{
			{typ: "char", in: `a`, want: runeptr('a')},
			{typ: "char", in: `Z`, want: runeptr('Z')},
			{typ: "char", in: `馬`, want: runeptr('馬')},
			{typ: "char", in: `駮`, want: runeptr('駮')},
		},
	}}.execute(t)
}

func TestChar_NoValuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return nil
		},
		rows: []test_valuer_row{
			{typ: "char", in: `a`, want: strptr(`a`)},
			{typ: "char", in: `馬`, want: strptr(`馬`)},
			{typ: "char", in: []byte{'a'}, want: strptr(`a`)},
			{typ: "char", in: []byte{'Z'}, want: strptr(`Z`)},
		},
	}}.execute(t)
}

func TestChar_NoScanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			d := new([]byte)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "char", in: `a`, want: bytesptr(`a`)},
			{typ: "char", in: `Z`, want: bytesptr(`Z`)},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			d := new(string)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "char", in: `a`, want: strptr(`a`)},
			{typ: "char", in: `馬`, want: strptr(`馬`)},
		},
	}}.execute(t)
}
