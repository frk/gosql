package pg2go

import (
	"testing"
)

func TestNumRangeArray(t *testing.T) {
	testlist2{{
		valuer:  NumRangeArrayFromIntArray2Slice,
		scanner: NumRangeArrayToIntArray2Slice,
		data: []testdata{
			{
				input:  [][2]int{{-9223372036854775808, 9223372036854775807}, {0, 21}},
				output: [][2]int{{-9223372036854775808, 9223372036854775807}, {0, 21}}},
		},
	}, {
		valuer:  NumRangeArrayFromInt8Array2Slice,
		scanner: NumRangeArrayToInt8Array2Slice,
		data: []testdata{
			{
				input:  [][2]int8{{-128, 127}, {0, 21}},
				output: [][2]int8{{-128, 127}, {0, 21}}},
		},
	}, {
		valuer:  NumRangeArrayFromInt16Array2Slice,
		scanner: NumRangeArrayToInt16Array2Slice,
		data: []testdata{
			{
				input:  [][2]int16{{-32768, 32767}, {0, 21}},
				output: [][2]int16{{-32768, 32767}, {0, 21}}},
		},
	}, {
		valuer:  NumRangeArrayFromInt32Array2Slice,
		scanner: NumRangeArrayToInt32Array2Slice,
		data: []testdata{
			{
				input:  [][2]int32{{-2147483648, 2147483647}, {0, 21}},
				output: [][2]int32{{-2147483648, 2147483647}, {0, 21}}},
		},
	}, {
		valuer:  NumRangeArrayFromInt64Array2Slice,
		scanner: NumRangeArrayToInt64Array2Slice,
		data: []testdata{
			{
				input:  [][2]int64{{-9223372036854775808, 9223372036854775807}, {0, 21}},
				output: [][2]int64{{-9223372036854775808, 9223372036854775807}, {0, 21}}},
		},
	}, {
		valuer:  NumRangeArrayFromUintArray2Slice,
		scanner: NumRangeArrayToUintArray2Slice,
		data: []testdata{
			{
				input:  [][2]uint{{0, 9223372036854775807}, {0, 21}},
				output: [][2]uint{{0, 9223372036854775807}, {0, 21}}},
		},
	}, {
		valuer:  NumRangeArrayFromUint8Array2Slice,
		scanner: NumRangeArrayToUint8Array2Slice,
		data: []testdata{
			{
				input:  [][2]uint8{{0, 255}, {0, 21}},
				output: [][2]uint8{{0, 255}, {0, 21}}},
		},
	}, {
		valuer:  NumRangeArrayFromUint16Array2Slice,
		scanner: NumRangeArrayToUint16Array2Slice,
		data: []testdata{
			{
				input:  [][2]uint16{{0, 65535}, {0, 21}},
				output: [][2]uint16{{0, 65535}, {0, 21}}},
		},
	}, {
		valuer:  NumRangeArrayFromUint32Array2Slice,
		scanner: NumRangeArrayToUint32Array2Slice,
		data: []testdata{
			{
				input:  [][2]uint32{{0, 4294967295}, {0, 21}},
				output: [][2]uint32{{0, 4294967295}, {0, 21}}},
		},
	}, {
		valuer:  NumRangeArrayFromUint64Array2Slice,
		scanner: NumRangeArrayToUint64Array2Slice,
		data: []testdata{
			{
				input:  [][2]uint64{{0, 9223372036854775807}, {0, 21}},
				output: [][2]uint64{{0, 9223372036854775807}, {0, 21}}},
		},
	}, {
		valuer:  NumRangeArrayFromFloat32Array2Slice,
		scanner: NumRangeArrayToFloat32Array2Slice,
		data: []testdata{
			{
				input:  [][2]float32{{-2147483648.0, 214748364.0}, {0.0, 21.0}},
				output: [][2]float32{{-2147483648.0, 214748364.0}, {0.0, 21.0}}},
		},
	}, {
		valuer:  NumRangeArrayFromFloat64Array2Slice,
		scanner: NumRangeArrayToFloat64Array2Slice,
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
	}}.execute(t, "numrangearr")
}
