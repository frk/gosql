package pgsql

import (
	"testing"
)

func TestInt4Array(t *testing.T) {
	testlist2{{
		valuer:  Int4ArrayFromIntSlice,
		scanner: Int4ArrayToIntSlice,
		data: []testdata{
			{
				input:  []int{-2147483648, 2147483647},
				output: []int{-2147483648, 2147483647}},
		},
	}, {
		valuer:  Int4ArrayFromInt8Slice,
		scanner: Int4ArrayToInt8Slice,
		data: []testdata{
			{
				input:  []int8{-128, 127},
				output: []int8{-128, 127}},
		},
	}, {
		valuer:  Int4ArrayFromInt16Slice,
		scanner: Int4ArrayToInt16Slice,
		data: []testdata{
			{
				input:  []int16{-32768, 32767},
				output: []int16{-32768, 32767}},
		},
	}, {
		valuer:  Int4ArrayFromInt32Slice,
		scanner: Int4ArrayToInt32Slice,
		data: []testdata{
			{
				input:  []int32{-2147483648, 2147483647},
				output: []int32{-2147483648, 2147483647}},
		},
	}, {
		valuer:  Int4ArrayFromInt64Slice,
		scanner: Int4ArrayToInt64Slice,
		data: []testdata{
			{
				input:  []int64{-2147483648, 2147483647},
				output: []int64{-2147483648, 2147483647}},
		},
	}, {
		valuer:  Int4ArrayFromUintSlice,
		scanner: Int4ArrayToUintSlice,
		data: []testdata{
			{
				input:  []uint{0, 2147483647},
				output: []uint{0, 2147483647}},
		},
	}, {
		valuer:  Int4ArrayFromUint8Slice,
		scanner: Int4ArrayToUint8Slice,
		data: []testdata{
			{
				input:  []uint8{0, 255},
				output: []uint8{0, 255}},
		},
	}, {
		valuer:  Int4ArrayFromUint16Slice,
		scanner: Int4ArrayToUint16Slice,
		data: []testdata{
			{
				input:  []uint16{0, 65535},
				output: []uint16{0, 65535}},
		},
	}, {
		valuer:  Int4ArrayFromUint32Slice,
		scanner: Int4ArrayToUint32Slice,
		data: []testdata{
			{
				input:  []uint32{0, 2147483647},
				output: []uint32{0, 2147483647}},
		},
	}, {
		valuer:  Int4ArrayFromUint64Slice,
		scanner: Int4ArrayToUint64Slice,
		data: []testdata{
			{
				input:  []uint64{0, 2147483647},
				output: []uint64{0, 2147483647}},
		},
	}, {
		valuer:  Int4ArrayFromFloat32Slice,
		scanner: Int4ArrayToFloat32Slice,
		data: []testdata{
			{
				input:  []float32{0, 214748364.0},
				output: []float32{0, 214748364.0}},
		},
	}, {
		valuer:  Int4ArrayFromFloat64Slice,
		scanner: Int4ArrayToFloat64Slice,
		data: []testdata{
			{
				input:  []float64{0, 2147483647.0},
				output: []float64{0, 2147483647.0}},
		},
	}, {
		data: []testdata{
			{
				input:  string("{-2147483648,2147483647}"),
				output: string(`{-2147483648,2147483647}`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte("{-2147483648,2147483647}"),
				output: []byte(`{-2147483648,2147483647}`)},
		},
	}}.execute(t, "int4arr")
}
