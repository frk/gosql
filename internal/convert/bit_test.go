package convert

import (
	"testing"
)

func TestBit(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(BitFromBool)
		},
		scanner: func() (interface{}, interface{}) {
			d := new(bool)
			return d, d
		},
		data: []testdata{
			{input: bool(true), output: boolptr(true)},
			{input: bool(false), output: boolptr(false)},
		},
	}, {
		valuer: func() interface{} {
			return nil // uint8
		},
		scanner: func() (interface{}, interface{}) {
			d := new(uint8)
			return d, d
		},
		data: []testdata{
			{input: uint8(1), output: u8ptr(1)},
			{input: uint8(0), output: u8ptr(0)},
		},
	}, {
		valuer: func() interface{} {
			return nil // uint
		},
		scanner: func() (interface{}, interface{}) {
			d := new(uint)
			return d, d
		},
		data: []testdata{
			{input: uint(0), output: uptr(0)},
			{input: uint(1), output: uptr(1)},
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
			{input: string("0"), output: strptr(`0`)},
			{input: string("1"), output: strptr(`1`)},
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
			{input: []byte(`0`), output: bytesptr(`0`)},
			{input: []byte(`1`), output: bytesptr(`1`)},
		},
	}}.execute(t, "bit")
}
