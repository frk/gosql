package pgsql

import (
	"testing"
)

func TestBit(t *testing.T) {
	testlist2{{
		valuer: BitFromBool,
		data: []testdata{
			{input: bool(true), output: bool(true)},
			{input: bool(false), output: bool(false)},
		},
	}, {
		data: []testdata{
			{input: uint8(1), output: uint8(1)},
			{input: uint8(0), output: uint8(0)},
		},
	}, {
		data: []testdata{
			{input: uint(0), output: uint(0)},
			{input: uint(1), output: uint(1)},
		},
	}, {
		data: []testdata{
			{input: string("0"), output: string(`0`)},
			{input: string("1"), output: string(`1`)},
		},
	}, {
		data: []testdata{
			{input: []byte(`0`), output: []byte(`0`)},
			{input: []byte(`1`), output: []byte(`1`)},
		},
	}}.execute(t, "bit")
}
