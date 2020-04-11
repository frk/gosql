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
			{input: byte('A'), output: byte('A')},
			{input: byte('b'), output: byte('b')},
			{input: byte('X'), output: byte('X')},
			{input: byte('z'), output: byte('z')},
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
			{input: rune('A'), output: rune('A')},
			{input: rune('馬'), output: rune('馬')},
			{input: rune('駮'), output: rune('駮')},
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
			{input: string(`a`), output: string(`a`)},
			{input: string(`馬`), output: string(`馬`)},
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
			{input: []byte{'a'}, output: []byte(`a`)},
			{input: []byte{'Z'}, output: []byte(`Z`)},
		},
	}}.execute(t, "bpchar")
}
