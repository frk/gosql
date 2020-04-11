package convert

import (
	"testing"
)

func TestBitArray(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(BitArrayFromBoolSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := BitArrayToBoolSlice{Val: new([]bool)}
			return s, s.Val
		},
		data: []testdata{
			{input: nil, output: []bool(nil)},
			{input: []bool{}, output: []bool{}},
			{input: []bool{true}, output: []bool{true}},
			{input: []bool{false}, output: []bool{false}},
			{
				input:  []bool{false, false, false, true, true, false, true, true},
				output: []bool{false, false, false, true, true, false, true, true}},
		},
	}, {
		valuer: func() interface{} {
			return new(BitArrayFromUint8Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := BitArrayToUint8Slice{Val: new([]uint8)}
			return s, s.Val
		},
		data: []testdata{
			{input: nil, output: []uint8(nil)},
			{input: []uint8{}, output: []uint8{}},
			{input: []uint8{1}, output: []uint8{1}},
			{input: []uint8{0}, output: []uint8{0}},
			{
				input:  []uint8{0, 0, 0, 1, 1, 0, 1, 1},
				output: []uint8{0, 0, 0, 1, 1, 0, 1, 1}},
		},
	}, {
		valuer: func() interface{} {
			return new(BitArrayFromUintSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := BitArrayToUintSlice{Val: new([]uint)}
			return s, s.Val
		},
		data: []testdata{
			{input: nil, output: []uint(nil)},
			{input: []uint{}, output: []uint{}},
			{input: []uint{1}, output: []uint{1}},
			{input: []uint{0}, output: []uint{0}},
			{
				input:  []uint{0, 0, 0, 1, 1, 0, 1, 1},
				output: []uint{0, 0, 0, 1, 1, 0, 1, 1}},
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
			{input: string("{}"), output: string(`{}`)},
			{input: string("{1,0}"), output: string(`{1,0}`)},
			{input: string("{0,1}"), output: string(`{0,1}`)},
			{
				input:  string("{0,1,1,1,0,1,0,0}"),
				output: string(`{0,1,1,1,0,1,0,0}`)},
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
			{input: nil, output: []byte(nil)},
			{input: []byte("{}"), output: []byte(`{}`)},
			{input: []byte("{1,0}"), output: []byte(`{1,0}`)},
			{input: []byte("{0,1}"), output: []byte(`{0,1}`)},
			{
				input:  []byte("{0,1,1,1,0,1,0,0}"),
				output: []byte(`{0,1,1,1,0,1,0,0}`)},
		},
	}}.execute(t, "bitarr")
}
