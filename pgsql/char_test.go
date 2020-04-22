package pgsql

import (
	"testing"
)

func TestChar(t *testing.T) {
	testlist2{{
		valuer:  CharFromByte,
		scanner: CharToByte,
		data: []testdata{
			{input: byte('A'), output: byte('A')},
			{input: byte('b'), output: byte('b')},
			{input: byte('X'), output: byte('X')},
		},
	}, {
		valuer:  CharFromRune,
		scanner: CharToRune,
		data: []testdata{
			{input: rune('A'), output: rune('A')},
			{input: rune('z'), output: rune('z')},
			{input: rune('馬'), output: rune('馬')},
			{input: rune('駮'), output: rune('駮')},
		},
	}, {
		data: []testdata{
			{input: string("A"), output: string("A")},
			{input: string("z"), output: string("z")},
			{input: string("馬"), output: string("馬")},
			{input: string("駮"), output: string("駮")},
		},
	}, {
		data: []testdata{
			{input: []byte("A"), output: []byte("A")},
			{input: []byte("z"), output: []byte("z")},
		},
	}}.execute(t, "char")
}
