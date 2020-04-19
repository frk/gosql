package convert

import (
	"testing"
)

func TestNumRangeArray(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(NumRangeArrayFromIntArray2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := NumRangeArrayToIntArray2Slice{Val: new([][2]int)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  [][2]int{{-9223372036854775808, 9223372036854775807}, {0, 21}},
				output: [][2]int{{-9223372036854775808, 9223372036854775807}, {0, 21}}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumRangeArrayFromInt8Array2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := NumRangeArrayToInt8Array2Slice{Val: new([][2]int8)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  [][2]int8{{-128, 127}, {0, 21}},
				output: [][2]int8{{-128, 127}, {0, 21}}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumRangeArrayFromInt16Array2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := NumRangeArrayToInt16Array2Slice{Val: new([][2]int16)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  [][2]int16{{-32768, 32767}, {0, 21}},
				output: [][2]int16{{-32768, 32767}, {0, 21}}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumRangeArrayFromInt32Array2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := NumRangeArrayToInt32Array2Slice{Val: new([][2]int32)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  [][2]int32{{-2147483648, 2147483647}, {0, 21}},
				output: [][2]int32{{-2147483648, 2147483647}, {0, 21}}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumRangeArrayFromInt64Array2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := NumRangeArrayToInt64Array2Slice{Val: new([][2]int64)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  [][2]int64{{-9223372036854775808, 9223372036854775807}, {0, 21}},
				output: [][2]int64{{-9223372036854775808, 9223372036854775807}, {0, 21}}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumRangeArrayFromUintArray2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := NumRangeArrayToUintArray2Slice{Val: new([][2]uint)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  [][2]uint{{0, 9223372036854775807}, {0, 21}},
				output: [][2]uint{{0, 9223372036854775807}, {0, 21}}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumRangeArrayFromUint8Array2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := NumRangeArrayToUint8Array2Slice{Val: new([][2]uint8)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  [][2]uint8{{0, 255}, {0, 21}},
				output: [][2]uint8{{0, 255}, {0, 21}}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumRangeArrayFromUint16Array2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := NumRangeArrayToUint16Array2Slice{Val: new([][2]uint16)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  [][2]uint16{{0, 65535}, {0, 21}},
				output: [][2]uint16{{0, 65535}, {0, 21}}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumRangeArrayFromUint32Array2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := NumRangeArrayToUint32Array2Slice{Val: new([][2]uint32)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  [][2]uint32{{0, 4294967295}, {0, 21}},
				output: [][2]uint32{{0, 4294967295}, {0, 21}}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumRangeArrayFromUint64Array2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := NumRangeArrayToUint64Array2Slice{Val: new([][2]uint64)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  [][2]uint64{{0, 9223372036854775807}, {0, 21}},
				output: [][2]uint64{{0, 9223372036854775807}, {0, 21}}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumRangeArrayFromFloat32Array2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := NumRangeArrayToFloat32Array2Slice{Val: new([][2]float32)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  [][2]float32{{-2147483648.0, 214748364.0}, {0.0, 21.0}},
				output: [][2]float32{{-2147483648.0, 214748364.0}, {0.0, 21.0}}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumRangeArrayFromFloat64Array2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := NumRangeArrayToFloat64Array2Slice{Val: new([][2]float64)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  [][2]float64{{-9223372036854775808.0, 922337203685477580.0}, {0.0, 21.0}},
				output: [][2]float64{{-9223372036854775808.0, 922337203685477580.0}, {0.0, 21.0}}},
		},
	}, {
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			v := new(string)
			return v, v
		},
		data: []testdata{
			{
				input:  string(`{"[-9223372036854775808,9223372036854775807)","[0,21)"}`),
				output: string(`{"[-9223372036854775808,9223372036854775807)","[0,21)"}`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			v := new([]byte)
			return v, v
		},
		data: []testdata{
			{
				input:  []byte(`{"[-9223372036854775808,9223372036854775807)","[0,21)"}`),
				output: []byte(`{"[-9223372036854775808,9223372036854775807)","[0,21)"}`)},
		},
	}}.execute(t, "numrangearr")
}
