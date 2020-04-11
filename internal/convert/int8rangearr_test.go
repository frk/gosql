package convert

import (
	"testing"
)

func TestInt8RangeArray(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(Int8RangeArrayFromIntArray2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int8RangeArrayToIntArray2Slice{Val: new([][2]int)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [][2]int{{-9223372036854775808, 9223372036854775807}, {0, 21}},
				output: &[][2]int{{-9223372036854775808, 9223372036854775807}, {0, 21}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int8RangeArrayFromInt8Array2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int8RangeArrayToInt8Array2Slice{Val: new([][2]int8)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [][2]int8{{-128, 127}, {0, 21}},
				output: &[][2]int8{{-128, 127}, {0, 21}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int8RangeArrayFromInt16Array2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int8RangeArrayToInt16Array2Slice{Val: new([][2]int16)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [][2]int16{{-32768, 32767}, {0, 21}},
				output: &[][2]int16{{-32768, 32767}, {0, 21}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int8RangeArrayFromInt32Array2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int8RangeArrayToInt32Array2Slice{Val: new([][2]int32)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [][2]int32{{-2147483648, 2147483647}, {0, 21}},
				output: &[][2]int32{{-2147483648, 2147483647}, {0, 21}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int8RangeArrayFromInt64Array2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int8RangeArrayToInt64Array2Slice{Val: new([][2]int64)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [][2]int64{{-9223372036854775808, 9223372036854775807}, {0, 21}},
				output: &[][2]int64{{-9223372036854775808, 9223372036854775807}, {0, 21}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int8RangeArrayFromUintArray2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int8RangeArrayToUintArray2Slice{Val: new([][2]uint)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [][2]uint{{0, 9223372036854775807}, {0, 21}},
				output: &[][2]uint{{0, 9223372036854775807}, {0, 21}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int8RangeArrayFromUint8Array2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int8RangeArrayToUint8Array2Slice{Val: new([][2]uint8)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [][2]uint8{{0, 255}, {0, 21}},
				output: &[][2]uint8{{0, 255}, {0, 21}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int8RangeArrayFromUint16Array2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int8RangeArrayToUint16Array2Slice{Val: new([][2]uint16)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [][2]uint16{{0, 65535}, {0, 21}},
				output: &[][2]uint16{{0, 65535}, {0, 21}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int8RangeArrayFromUint32Array2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int8RangeArrayToUint32Array2Slice{Val: new([][2]uint32)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [][2]uint32{{0, 4294967295}, {0, 21}},
				output: &[][2]uint32{{0, 4294967295}, {0, 21}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int8RangeArrayFromUint64Array2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int8RangeArrayToUint64Array2Slice{Val: new([][2]uint64)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [][2]uint64{{0, 9223372036854775807}, {0, 21}},
				output: &[][2]uint64{{0, 9223372036854775807}, {0, 21}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int8RangeArrayFromFloat32Array2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int8RangeArrayToFloat32Array2Slice{Val: new([][2]float32)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [][2]float32{{-2147483648.0, 214748364.0}, {0.0, 21.0}},
				output: &[][2]float32{{-2147483648.0, 214748364.0}, {0.0, 21.0}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int8RangeArrayFromFloat64Array2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int8RangeArrayToFloat64Array2Slice{Val: new([][2]float64)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [][2]float64{{-9223372036854775808.0, 922337203685477580.0}, {0.0, 21.0}},
				output: &[][2]float64{{-9223372036854775808.0, 922337203685477580.0}, {0.0, 21.0}}},
		},
	}, {
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			s := new(string)
			return s, s
		},
		data: []testdata{
			{
				input:  `{"[-9223372036854775808,9223372036854775807)","[0,21)"}`,
				output: strptr(`{"[-9223372036854775808,9223372036854775807)","[0,21)"}`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			s := new([]byte)
			return s, s
		},
		data: []testdata{
			{
				input:  []byte(`{"[-9223372036854775808,9223372036854775807)","[0,21)"}`),
				output: bytesptr(`{"[-9223372036854775808,9223372036854775807)","[0,21)"}`)},
		},
	}}.execute(t, "int8rangearr")
}
