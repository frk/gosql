package pg2go

import (
	"testing"
)

func TestInt8Range(t *testing.T) {
	testlist2{{
		valuer:  Int8RangeFromIntArray2,
		scanner: Int8RangeToIntArray2,
		data: []testdata{
			{
				input:  [2]int{-9223372036854775808, 9223372036854775807},
				output: [2]int{-9223372036854775808, 9223372036854775807}},
		},
	}, {
		valuer:  Int8RangeFromInt8Array2,
		scanner: Int8RangeToInt8Array2,
		data: []testdata{
			{input: [2]int8{-128, 127}, output: [2]int8{-128, 127}},
		},
	}, {
		valuer:  Int8RangeFromInt16Array2,
		scanner: Int8RangeToInt16Array2,
		data: []testdata{
			{
				input:  [2]int16{-32768, 32767},
				output: [2]int16{-32768, 32767}},
		},
	}, {
		valuer:  Int8RangeFromInt32Array2,
		scanner: Int8RangeToInt32Array2,
		data: []testdata{
			{
				input:  [2]int32{-2147483648, 2147483647},
				output: [2]int32{-2147483648, 2147483647}},
		},
	}, {
		valuer:  Int8RangeFromInt64Array2,
		scanner: Int8RangeToInt64Array2,
		data: []testdata{
			{
				input:  [2]int64{-9223372036854775808, 9223372036854775807},
				output: [2]int64{-9223372036854775808, 9223372036854775807}},
		},
	}, {
		valuer:  Int8RangeFromUintArray2,
		scanner: Int8RangeToUintArray2,
		data: []testdata{
			{
				input:  [2]uint{0, 9223372036854775807},
				output: [2]uint{0, 9223372036854775807}},
		},
	}, {
		valuer:  Int8RangeFromUint8Array2,
		scanner: Int8RangeToUint8Array2,
		data: []testdata{
			{
				input:  [2]uint8{0, 255},
				output: [2]uint8{0, 255}},
		},
	}, {
		valuer:  Int8RangeFromUint16Array2,
		scanner: Int8RangeToUint16Array2,
		data: []testdata{
			{
				input:  [2]uint16{0, 65535},
				output: [2]uint16{0, 65535}},
		},
	}, {
		valuer:  Int8RangeFromUint32Array2,
		scanner: Int8RangeToUint32Array2,
		data: []testdata{
			{
				input:  [2]uint32{0, 4294967295},
				output: [2]uint32{0, 4294967295}},
		},
	}, {
		valuer:  Int8RangeFromUint64Array2,
		scanner: Int8RangeToUint64Array2,
		data: []testdata{
			{
				input:  [2]uint64{0, 9223372036854775807},
				output: [2]uint64{0, 9223372036854775807}},
		},
	}, {
		valuer:  Int8RangeFromFloat32Array2,
		scanner: Int8RangeToFloat32Array2,
		data: []testdata{
			{
				input:  [2]float32{-9223372036854775808.0, 922337203685477580.0},
				output: [2]float32{-9223372036854775808.0, 922337203685477580.0}},
		},
	}, {
		valuer:  Int8RangeFromFloat64Array2,
		scanner: Int8RangeToFloat64Array2,
		data: []testdata{
			{
				input:  [2]float64{-9223372036854775808.0, 922337203685477580.0},
				output: [2]float64{-9223372036854775808.0, 922337203685477580.0}},
		},
	}, {
		data: []testdata{
			{
				input:  string("[-9223372036854775808,9223372036854775807)"),
				output: string(`[-9223372036854775808,9223372036854775807)`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte("[-9223372036854775808,9223372036854775807)"),
				output: []byte(`[-9223372036854775808,9223372036854775807)`)},
		},
	}}.execute(t, "int8range")
}
