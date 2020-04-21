package pg2go

import (
	"testing"
)

func TestBPChar(t *testing.T) {
	testlist2{{
		valuer:  BPCharFromByte,
		scanner: BPCharToByte,
		data: []testdata{
			{input: byte('A'), output: byte('A')},
			{input: byte('b'), output: byte('b')},
			{input: byte('X'), output: byte('X')},
			{input: byte('z'), output: byte('z')},
		},
	}, {
		valuer:  BPCharFromRune,
		scanner: BPCharToRune,
		data: []testdata{
			{input: rune('A'), output: rune('A')},
			{input: rune('馬'), output: rune('馬')},
			{input: rune('駮'), output: rune('駮')},
		},
	}, {
		data: []testdata{
			{input: string(`a`), output: string(`a`)},
			{input: string(`馬`), output: string(`馬`)},
		},
	}, {
		data: []testdata{
			{input: []byte{'a'}, output: []byte(`a`)},
			{input: []byte{'Z'}, output: []byte(`Z`)},
		},
	}}.execute(t, "bpchar")
}
