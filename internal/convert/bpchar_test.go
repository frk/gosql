package convert

import (
	"testing"
)

func TestBPChar(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(BPCharFromByte)
		},
		scanner: func() (interface{}, interface{}) {
			s := &BPCharToByte{Val: new(byte)}
			return s, s.Val
		},
		data: []testdata{
			{input: byte('A'), output: byteptr('A')},
			{input: byte('b'), output: byteptr('b')},
			{input: byte('X'), output: byteptr('X')},
			{input: byte('z'), output: byteptr('z')},
		},
	}, {
		valuer: func() interface{} {
			return new(BPCharFromRune)
		},
		scanner: func() (interface{}, interface{}) {
			s := &BPCharToRune{Val: new(rune)}
			return s, s.Val
		},
		data: []testdata{
			{input: rune('A'), output: runeptr('A')},
			{input: rune('馬'), output: runeptr('馬')},
			{input: rune('駮'), output: runeptr('駮')},
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
			{input: string(`a`), output: strptr(`a`)},
			{input: string(`馬`), output: strptr(`馬`)},
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
			{input: []byte{'a'}, output: bytesptr(`a`)},
			{input: []byte{'Z'}, output: bytesptr(`Z`)},
		},
	}}.execute(t, "bpchar")
}
