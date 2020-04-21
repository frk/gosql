package pg2go

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
			{input: bool(true), output: bool(true)},
			{input: bool(false), output: bool(false)},
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
			{input: uint8(1), output: uint8(1)},
			{input: uint8(0), output: uint8(0)},
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
			{input: uint(0), output: uint(0)},
			{input: uint(1), output: uint(1)},
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
			{input: string("0"), output: string(`0`)},
			{input: string("1"), output: string(`1`)},
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
			{input: []byte(`0`), output: []byte(`0`)},
			{input: []byte(`1`), output: []byte(`1`)},
		},
	}}.execute(t, "bit")
}
