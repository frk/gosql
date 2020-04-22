package pgsql

import (
	"testing"
)

func TestNumericArray(t *testing.T) {
	testlist2{{
		valuer:  NumericArrayFromIntSlice,
		scanner: NumericArrayToIntSlice,
		data: []testdata{
			{
				input:  []int{-9223372036854775808, 9223372036854775807},
				output: []int{-9223372036854775808, 9223372036854775807}},
		},
	}, {
		valuer:  NumericArrayFromInt8Slice,
		scanner: NumericArrayToInt8Slice,
		data: []testdata{
			{
				input:  []int8{-128, 127},
				output: []int8{-128, 127}},
		},
	}, {
		valuer:  NumericArrayFromInt16Slice,
		scanner: NumericArrayToInt16Slice,
		data: []testdata{
			{
				input:  []int16{-32768, 32767},
				output: []int16{-32768, 32767}},
		},
	}, {
		valuer:  NumericArrayFromInt32Slice,
		scanner: NumericArrayToInt32Slice,
		data: []testdata{
			{
				input:  []int32{-2147483648, 2147483647},
				output: []int32{-2147483648, 2147483647}},
		},
	}, {
		valuer:  NumericArrayFromInt64Slice,
		scanner: NumericArrayToInt64Slice,
		data: []testdata{
			{
				input:  []int64{-9223372036854775808, 9223372036854775807},
				output: []int64{-9223372036854775808, 9223372036854775807}},
		},
	}, {
		valuer:  NumericArrayFromUintSlice,
		scanner: NumericArrayToUintSlice,
		data: []testdata{
			{
				input:  []uint{0, 2147483647},
				output: []uint{0, 2147483647}},
		},
	}, {
		valuer:  NumericArrayFromUint8Slice,
		scanner: NumericArrayToUint8Slice,
		data: []testdata{
			{
				input:  []uint8{0, 255},
				output: []uint8{0, 255}},
		},
	}, {
		valuer:  NumericArrayFromUint16Slice,
		scanner: NumericArrayToUint16Slice,
		data: []testdata{
			{
				input:  []uint16{0, 65535},
				output: []uint16{0, 65535}},
		},
	}, {
		valuer:  NumericArrayFromUint32Slice,
		scanner: NumericArrayToUint32Slice,
		data: []testdata{
			{
				input:  []uint32{0, 4294967295},
				output: []uint32{0, 4294967295}},
		},
	}, {
		valuer:  NumericArrayFromUint64Slice,
		scanner: NumericArrayToUint64Slice,
		data: []testdata{
			{
				input:  []uint64{0, 9223372036854775807},
				output: []uint64{0, 9223372036854775807}},
		},
	}, {
		valuer:  NumericArrayFromFloat32Slice,
		scanner: NumericArrayToFloat32Slice,
		data: []testdata{
			{
				input:  []float32{0, 2147483647.0},
				output: []float32{0, 2147483647.0}},
		},
	}, {
		valuer:  NumericArrayFromFloat64Slice,
		scanner: NumericArrayToFloat64Slice,
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
	}}.execute(t, "numericarr")
}
