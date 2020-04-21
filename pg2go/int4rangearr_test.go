package pg2go

import (
	"testing"
)

func TestInt4RangeArray(t *testing.T) {
	testlist2{{
		valuer:  Int4RangeArrayFromIntArray2Slice,
		scanner: Int4RangeArrayToIntArray2Slice,
		data: []testdata{
			{
				input:  [][2]int{{-2147483648, 2147483647}, {0, 21}},
				output: [][2]int{{-2147483648, 2147483647}, {0, 21}}},
		},
	}, {
		valuer:  Int4RangeArrayFromInt8Array2Slice,
		scanner: Int4RangeArrayToInt8Array2Slice,
		data: []testdata{
			{
				input:  [][2]int8{{-128, 127}, {0, 21}},
				output: [][2]int8{{-128, 127}, {0, 21}}},
		},
	}, {
		valuer:  Int4RangeArrayFromInt16Array2Slice,
		scanner: Int4RangeArrayToInt16Array2Slice,
		data: []testdata{
			{
				input:  [][2]int16{{-32768, 32767}, {0, 21}},
				output: [][2]int16{{-32768, 32767}, {0, 21}}},
		},
	}, {
		valuer:  Int4RangeArrayFromInt32Array2Slice,
		scanner: Int4RangeArrayToInt32Array2Slice,
		data: []testdata{
			{
				input:  [][2]int32{{-2147483648, 2147483647}, {0, 21}},
				output: [][2]int32{{-2147483648, 2147483647}, {0, 21}}},
		},
	}, {
		valuer:  Int4RangeArrayFromInt64Array2Slice,
		scanner: Int4RangeArrayToInt64Array2Slice,
		data: []testdata{
			{
				input:  [][2]int64{{-2147483648, 2147483647}, {0, 21}},
				output: [][2]int64{{-2147483648, 2147483647}, {0, 21}}},
		},
	}, {
		valuer:  Int4RangeArrayFromUintArray2Slice,
		scanner: Int4RangeArrayToUintArray2Slice,
		data: []testdata{
			{
				input:  [][2]uint{{0, 2147483647}, {0, 21}},
				output: [][2]uint{{0, 2147483647}, {0, 21}}},
		},
	}, {
		valuer:  Int4RangeArrayFromUint8Array2Slice,
		scanner: Int4RangeArrayToUint8Array2Slice,
		data: []testdata{
			{
				input:  [][2]uint8{{0, 255}, {0, 21}},
				output: [][2]uint8{{0, 255}, {0, 21}}},
		},
	}, {
		valuer:  Int4RangeArrayFromUint16Array2Slice,
		scanner: Int4RangeArrayToUint16Array2Slice,
		data: []testdata{
			{
				input:  [][2]uint16{{0, 65535}, {0, 21}},
				output: [][2]uint16{{0, 65535}, {0, 21}}},
		},
	}, {
		valuer:  Int4RangeArrayFromUint32Array2Slice,
		scanner: Int4RangeArrayToUint32Array2Slice,
		data: []testdata{
			{
				input:  [][2]uint32{{0, 2147483647}, {0, 21}},
				output: [][2]uint32{{0, 2147483647}, {0, 21}}},
		},
	}, {
		valuer:  Int4RangeArrayFromUint64Array2Slice,
		scanner: Int4RangeArrayToUint64Array2Slice,
		data: []testdata{
			{
				input:  [][2]uint64{{0, 2147483647}, {0, 21}},
				output: [][2]uint64{{0, 2147483647}, {0, 21}}},
		},
	}, {
		valuer:  Int4RangeArrayFromFloat32Array2Slice,
		scanner: Int4RangeArrayToFloat32Array2Slice,
		data: []testdata{
			{
				input:  [][2]float32{{-2147483648.0, 214748364.0}, {0.0, 21.0}},
				output: [][2]float32{{-2147483648.0, 214748364.0}, {0.0, 21.0}}},
		},
	}, {
		valuer:  Int4RangeArrayFromFloat64Array2Slice,
		scanner: Int4RangeArrayToFloat64Array2Slice,
		data: []testdata{
			{
				input:  [][2]float64{{-2147483648.0, 2147483647.0}, {0.0, 21.0}},
				output: [][2]float64{{-2147483648.0, 2147483647.0}, {0.0, 21.0}}},
		},
	}, {
		data: []testdata{
			{
				input:  string(`{"[-2147483648,2147483647)","[0,21)"}`),
				output: string(`{"[-2147483648,2147483647)","[0,21)"}`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte(`{"[-2147483648,2147483647)","[0,21)"}`),
				output: []byte(`{"[-2147483648,2147483647)","[0,21)"}`)},
		},
	}}.execute(t, "int4rangearr")
}
