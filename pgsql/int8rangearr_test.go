package pgsql

import (
	"testing"
)

func TestInt8RangeArray(t *testing.T) {
	testlist2{{
		valuer:  Int8RangeArrayFromIntArray2Slice,
		scanner: Int8RangeArrayToIntArray2Slice,
		data: []testdata{
			{
				input:  [][2]int{{-9223372036854775808, 9223372036854775807}, {0, 21}},
				output: [][2]int{{-9223372036854775808, 9223372036854775807}, {0, 21}}},
		},
	}, {
		valuer:  Int8RangeArrayFromInt8Array2Slice,
		scanner: Int8RangeArrayToInt8Array2Slice,
		data: []testdata{
			{
				input:  [][2]int8{{-128, 127}, {0, 21}},
				output: [][2]int8{{-128, 127}, {0, 21}}},
		},
	}, {
		valuer:  Int8RangeArrayFromInt16Array2Slice,
		scanner: Int8RangeArrayToInt16Array2Slice,
		data: []testdata{
			{
				input:  [][2]int16{{-32768, 32767}, {0, 21}},
				output: [][2]int16{{-32768, 32767}, {0, 21}}},
		},
	}, {
		valuer:  Int8RangeArrayFromInt32Array2Slice,
		scanner: Int8RangeArrayToInt32Array2Slice,
		data: []testdata{
			{
				input:  [][2]int32{{-2147483648, 2147483647}, {0, 21}},
				output: [][2]int32{{-2147483648, 2147483647}, {0, 21}}},
		},
	}, {
		valuer:  Int8RangeArrayFromInt64Array2Slice,
		scanner: Int8RangeArrayToInt64Array2Slice,
		data: []testdata{
			{
				input:  [][2]int64{{-9223372036854775808, 9223372036854775807}, {0, 21}},
				output: [][2]int64{{-9223372036854775808, 9223372036854775807}, {0, 21}}},
		},
	}, {
		valuer:  Int8RangeArrayFromUintArray2Slice,
		scanner: Int8RangeArrayToUintArray2Slice,
		data: []testdata{
			{
				input:  [][2]uint{{0, 9223372036854775807}, {0, 21}},
				output: [][2]uint{{0, 9223372036854775807}, {0, 21}}},
		},
	}, {
		valuer:  Int8RangeArrayFromUint8Array2Slice,
		scanner: Int8RangeArrayToUint8Array2Slice,
		data: []testdata{
			{
				input:  [][2]uint8{{0, 255}, {0, 21}},
				output: [][2]uint8{{0, 255}, {0, 21}}},
		},
	}, {
		valuer:  Int8RangeArrayFromUint16Array2Slice,
		scanner: Int8RangeArrayToUint16Array2Slice,
		data: []testdata{
			{
				input:  [][2]uint16{{0, 65535}, {0, 21}},
				output: [][2]uint16{{0, 65535}, {0, 21}}},
		},
	}, {
		valuer:  Int8RangeArrayFromUint32Array2Slice,
		scanner: Int8RangeArrayToUint32Array2Slice,
		data: []testdata{
			{
				input:  [][2]uint32{{0, 4294967295}, {0, 21}},
				output: [][2]uint32{{0, 4294967295}, {0, 21}}},
		},
	}, {
		valuer:  Int8RangeArrayFromUint64Array2Slice,
		scanner: Int8RangeArrayToUint64Array2Slice,
		data: []testdata{
			{
				input:  [][2]uint64{{0, 9223372036854775807}, {0, 21}},
				output: [][2]uint64{{0, 9223372036854775807}, {0, 21}}},
		},
	}, {
		valuer:  Int8RangeArrayFromFloat32Array2Slice,
		scanner: Int8RangeArrayToFloat32Array2Slice,
		data: []testdata{
			{
				input:  [][2]float32{{-2147483648.0, 214748364.0}, {0.0, 21.0}},
				output: [][2]float32{{-2147483648.0, 214748364.0}, {0.0, 21.0}}},
		},
	}, {
		valuer:  Int8RangeArrayFromFloat64Array2Slice,
		scanner: Int8RangeArrayToFloat64Array2Slice,
		data: []testdata{
			{
				input:  [][2]float64{{-9223372036854775808.0, 922337203685477580.0}, {0.0, 21.0}},
				output: [][2]float64{{-9223372036854775808.0, 922337203685477580.0}, {0.0, 21.0}}},
		},
	}, {
		data: []testdata{
			{
				input:  string(`{"[-9223372036854775808,9223372036854775807)","[0,21)"}`),
				output: string(`{"[-9223372036854775808,9223372036854775807)","[0,21)"}`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte(`{"[-9223372036854775808,9223372036854775807)","[0,21)"}`),
				output: []byte(`{"[-9223372036854775808,9223372036854775807)","[0,21)"}`)},
		},
	}}.execute(t, "int8rangearr")
}
