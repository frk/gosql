package pgsql

import (
	"testing"
)

func TestInt8Array(t *testing.T) {
	testlist2{{
		valuer:  Int8ArrayFromIntSlice,
		scanner: Int8ArrayToIntSlice,
		data: []testdata{
			{
				input:  []int{-9223372036854775808, 9223372036854775807},
				output: []int{-9223372036854775808, 9223372036854775807}},
		},
	}, {
		valuer:  Int8ArrayFromInt8Slice,
		scanner: Int8ArrayToInt8Slice,
		data: []testdata{
			{
				input:  []int8{-128, 127},
				output: []int8{-128, 127}},
		},
	}, {
		valuer:  Int8ArrayFromInt16Slice,
		scanner: Int8ArrayToInt16Slice,
		data: []testdata{
			{
				input:  []int16{-32768, 32767},
				output: []int16{-32768, 32767}},
		},
	}, {
		valuer:  Int8ArrayFromInt32Slice,
		scanner: Int8ArrayToInt32Slice,
		data: []testdata{
			{
				input:  []int32{-2147483648, 2147483647},
				output: []int32{-2147483648, 2147483647}},
		},
	}, {
		valuer:  Int8ArrayFromInt64Slice,
		scanner: Int8ArrayToInt64Slice,
		data: []testdata{
			{
				input:  []int64{-9223372036854775808, 9223372036854775807},
				output: []int64{-9223372036854775808, 9223372036854775807}},
		},
	}, {
		valuer:  Int8ArrayFromUintSlice,
		scanner: Int8ArrayToUintSlice,
		data: []testdata{
			{
				input:  []uint{0, 2147483647},
				output: []uint{0, 2147483647}},
		},
	}, {
		valuer:  Int8ArrayFromUint8Slice,
		scanner: Int8ArrayToUint8Slice,
		data: []testdata{
			{
				input:  []uint8{0, 255},
				output: []uint8{0, 255}},
		},
	}, {
		valuer:  Int8ArrayFromUint16Slice,
		scanner: Int8ArrayToUint16Slice,
		data: []testdata{
			{
				input:  []uint16{0, 65535},
				output: []uint16{0, 65535}},
		},
	}, {
		valuer:  Int8ArrayFromUint32Slice,
		scanner: Int8ArrayToUint32Slice,
		data: []testdata{
			{
				input:  []uint32{0, 4294967295},
				output: []uint32{0, 4294967295}},
		},
	}, {
		valuer:  Int8ArrayFromUint64Slice,
		scanner: Int8ArrayToUint64Slice,
		data: []testdata{
			{
				input:  []uint64{0, 9223372036854775807},
				output: []uint64{0, 9223372036854775807}},
		},
	}, {
		valuer:  Int8ArrayFromFloat32Slice,
		scanner: Int8ArrayToFloat32Slice,
		data: []testdata{
			{
				input:  []float32{0, 2147483647.0},
				output: []float32{0, 2147483647.0}},
		},
	}, {
		valuer:  Int8ArrayFromFloat64Slice,
		scanner: Int8ArrayToFloat64Slice,
		data: []testdata{
			{
				input:  []float64{0, 922337203685477580.0},
				output: []float64{0, 922337203685477580.0}},
		},
	}, {
		data: []testdata{
			{
				input:  string("{-9223372036854775808,9223372036854775807}"),
				output: string(`{-9223372036854775808,9223372036854775807}`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte("{-9223372036854775808,9223372036854775807}"),
				output: []byte(`{-9223372036854775808,9223372036854775807}`)},
		},
	}}.execute(t, "int8arr")
}
