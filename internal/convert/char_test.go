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
			{input: byte('A'), output: byte('A')},
			{input: byte('b'), output: byte('b')},
			{input: byte('X'), output: byte('X')},
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
			{input: rune('A'), output: rune('A')},
			{input: rune('z'), output: rune('z')},
			{input: rune('馬'), output: rune('馬')},
			{input: rune('駮'), output: rune('駮')},
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
			{input: string("A"), output: string("A")},
			{input: string("z"), output: string("z")},
			{input: string("馬"), output: string("馬")},
			{input: string("駮"), output: string("駮")},
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
			{input: []byte("A"), output: []byte("A")},
			{input: []byte("z"), output: []byte("z")},
		},
	}}.execute(t, "char")
}
