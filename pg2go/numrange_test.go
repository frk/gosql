package pg2go

import (
	"testing"
)

func TestNumRange(t *testing.T) {
	testlist2{{
		valuer:  NumRangeFromIntArray2,
		scanner: NumRangeToIntArray2,
		data: []testdata{
			{
				input:  [2]int{-9223372036854775808, 9223372036854775807},
				output: [2]int{-9223372036854775808, 9223372036854775807}},
		},
	}, {
		valuer:  NumRangeFromInt8Array2,
		scanner: NumRangeToInt8Array2,
		data: []testdata{
			{input: [2]int8{-128, 127}, output: [2]int8{-128, 127}},
		},
	}, {
		valuer:  NumRangeFromInt16Array2,
		scanner: NumRangeToInt16Array2,
		data: []testdata{
			{
				input:  [2]int16{-32768, 32767},
				output: [2]int16{-32768, 32767}},
		},
	}, {
		valuer:  NumRangeFromInt32Array2,
		scanner: NumRangeToInt32Array2,
		data: []testdata{
			{
				input:  [2]int32{-2147483648, 2147483647},
				output: [2]int32{-2147483648, 2147483647}},
		},
	}, {
		valuer:  NumRangeFromInt64Array2,
		scanner: NumRangeToInt64Array2,
		data: []testdata{
			{
				input:  [2]int64{-9223372036854775808, 9223372036854775807},
				output: [2]int64{-9223372036854775808, 9223372036854775807}},
		},
	}, {
		valuer:  NumRangeFromUintArray2,
		scanner: NumRangeToUintArray2,
		data: []testdata{
			{
				input:  [2]uint{0, 9223372036854775807},
				output: [2]uint{0, 9223372036854775807}},
		},
	}, {
		valuer:  NumRangeFromUint8Array2,
		scanner: NumRangeToUint8Array2,
		data: []testdata{
			{
				input:  [2]uint8{0, 255},
				output: [2]uint8{0, 255}},
		},
	}, {
		valuer:  NumRangeFromUint16Array2,
		scanner: NumRangeToUint16Array2,
		data: []testdata{
			{
				input:  [2]uint16{0, 65535},
				output: [2]uint16{0, 65535}},
		},
	}, {
		valuer:  NumRangeFromUint32Array2,
		scanner: NumRangeToUint32Array2,
		data: []testdata{
			{
				input:  [2]uint32{0, 4294967295},
				output: [2]uint32{0, 4294967295}},
		},
	}, {
		valuer:  NumRangeFromUint64Array2,
		scanner: NumRangeToUint64Array2,
		data: []testdata{
			{
				input:  [2]uint64{0, 9223372036854775807},
				output: [2]uint64{0, 9223372036854775807}},
		},
	}, {
		valuer:  NumRangeFromFloat32Array2,
		scanner: NumRangeToFloat32Array2,
		data: []testdata{
			{
				input:  [2]float32{-9223372036854775808.0, 922337203685477580.0},
				output: [2]float32{-9223372036854775808.0, 922337203685477580.0}},
		},
	}, {
		valuer:  NumRangeFromFloat64Array2,
		scanner: NumRangeToFloat64Array2,
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
	}}.execute(t, "numrange")
}
