package convert

import (
	"testing"
)

func TestNumericArray(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(NumericArrayFromIntSlice)
		},
		scanner: func() (interface{}, interface{}) {
			v := NumericArrayToIntSlice{Val: new([]int)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  []int{-9223372036854775808, 9223372036854775807},
				output: []int{-9223372036854775808, 9223372036854775807}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumericArrayFromInt8Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := NumericArrayToInt8Slice{Val: new([]int8)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  []int8{-128, 127},
				output: []int8{-128, 127}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumericArrayFromInt16Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := NumericArrayToInt16Slice{Val: new([]int16)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  []int16{-32768, 32767},
				output: []int16{-32768, 32767}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumericArrayFromInt32Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := NumericArrayToInt32Slice{Val: new([]int32)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  []int32{-2147483648, 2147483647},
				output: []int32{-2147483648, 2147483647}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumericArrayFromInt64Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := NumericArrayToInt64Slice{Val: new([]int64)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  []int64{-9223372036854775808, 9223372036854775807},
				output: []int64{-9223372036854775808, 9223372036854775807}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumericArrayFromUintSlice)
		},
		scanner: func() (interface{}, interface{}) {
			v := NumericArrayToUintSlice{Val: new([]uint)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  []uint{0, 2147483647},
				output: []uint{0, 2147483647}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumericArrayFromUint8Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := NumericArrayToUint8Slice{Val: new([]uint8)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  []uint8{0, 255},
				output: []uint8{0, 255}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumericArrayFromUint16Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := NumericArrayToUint16Slice{Val: new([]uint16)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  []uint16{0, 65535},
				output: []uint16{0, 65535}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumericArrayFromUint32Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := NumericArrayToUint32Slice{Val: new([]uint32)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  []uint32{0, 4294967295},
				output: []uint32{0, 4294967295}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumericArrayFromUint64Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := NumericArrayToUint64Slice{Val: new([]uint64)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  []uint64{0, 9223372036854775807},
				output: []uint64{0, 9223372036854775807}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumericArrayFromFloat32Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := NumericArrayToFloat32Slice{Val: new([]float32)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  []float32{0, 2147483647.0},
				output: []float32{0, 2147483647.0}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumericArrayFromFloat64Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := NumericArrayToFloat64Slice{Val: new([]float64)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  []float64{0, 922337203685477580.0},
				output: []float64{0, 922337203685477580.0}},
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
				input:  string("{-9223372036854775808,9223372036854775807}"),
				output: string(`{-9223372036854775808,9223372036854775807}`)},
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
				input:  []byte("{-9223372036854775808,9223372036854775807}"),
				output: []byte(`{-9223372036854775808,9223372036854775807}`)},
		},
	}}.execute(t, "numericarr")
}
