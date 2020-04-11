package convert

import (
	"testing"
)

func TestChar(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(CharFromByte)
		},
		scanner: func() (interface{}, interface{}) {
			s := &CharToByte{Val: new(byte)}
			return s, s.Val
		},
		data: []testdata{
			{input: byte('A'), output: byteptr('A')},
			{input: byte('b'), output: byteptr('b')},
			{input: byte('X'), output: byteptr('X')},
		},
	}, {
		valuer: func() interface{} {
			return new(CharFromRune)
		},
		scanner: func() (interface{}, interface{}) {
			s := &CharToRune{Val: new(rune)}
			return s, s.Val
		},
		data: []testdata{
			{input: rune('A'), output: runeptr('A')},
			{input: rune('z'), output: runeptr('z')},
			{input: rune('馬'), output: runeptr('馬')},
			{input: rune('駮'), output: runeptr('駮')},
		},
	}, {
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			s := new(string)
			return s, s
		},
		data: []testdata{
			{input: string("A"), output: strptr("A")},
			{input: string("z"), output: strptr("z")},
			{input: string("馬"), output: strptr("馬")},
			{input: string("駮"), output: strptr("駮")},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			s := new([]byte)
			return s, s
		},
		data: []testdata{
			{input: []byte("A"), output: bytesptr("A")},
			{input: []byte("z"), output: bytesptr("z")},
		},
	}}.execute(t, "char")
}
