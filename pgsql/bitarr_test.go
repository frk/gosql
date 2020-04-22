package pgsql

import (
	"testing"
)

func TestBitArray(t *testing.T) {
	testlist2{{
		valuer:  BitArrayFromBoolSlice,
		scanner: BitArrayToBoolSlice,
		data: []testdata{
			{input: []bool(nil), output: []bool(nil)},
			{input: []bool{}, output: []bool{}},
			{input: []bool{true}, output: []bool{true}},
			{input: []bool{false}, output: []bool{false}},
			{
				input:  []bool{false, false, false, true, true, false, true, true},
				output: []bool{false, false, false, true, true, false, true, true}},
		},
	}, {
		valuer:  BitArrayFromUint8Slice,
		scanner: BitArrayToUint8Slice,
		data: []testdata{
			{input: []uint8(nil), output: []uint8(nil)},
			{input: []uint8{}, output: []uint8{}},
			{input: []uint8{1}, output: []uint8{1}},
			{input: []uint8{0}, output: []uint8{0}},
			{
				input:  []uint8{0, 0, 0, 1, 1, 0, 1, 1},
				output: []uint8{0, 0, 0, 1, 1, 0, 1, 1}},
		},
	}, {
		valuer:  BitArrayFromUintSlice,
		scanner: BitArrayToUintSlice,
		data: []testdata{
			{input: []uint(nil), output: []uint(nil)},
			{input: []uint{}, output: []uint{}},
			{input: []uint{1}, output: []uint{1}},
			{input: []uint{0}, output: []uint{0}},
			{
				input:  []uint{0, 0, 0, 1, 1, 0, 1, 1},
				output: []uint{0, 0, 0, 1, 1, 0, 1, 1}},
		},
	}, {
		data: []testdata{
			{input: string("{}"), output: string(`{}`)},
			{input: string("{1,0}"), output: string(`{1,0}`)},
			{input: string("{0,1}"), output: string(`{0,1}`)},
			{
				input:  string("{0,1,1,1,0,1,0,0}"),
				output: string(`{0,1,1,1,0,1,0,0}`)},
		},
	}, {
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
