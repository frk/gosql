package pg2go

import (
	"testing"
)

func TestInt8Array(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(Int8ArrayFromIntSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int8ArrayToIntSlice{Val: new([]int)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []int{-9223372036854775808, 9223372036854775807},
				output: []int{-9223372036854775808, 9223372036854775807}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int8ArrayFromInt8Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int8ArrayToInt8Slice{Val: new([]int8)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []int8{-128, 127},
				output: []int8{-128, 127}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int8ArrayFromInt16Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int8ArrayToInt16Slice{Val: new([]int16)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []int16{-32768, 32767},
				output: []int16{-32768, 32767}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int8ArrayFromInt32Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int8ArrayToInt32Slice{Val: new([]int32)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []int32{-2147483648, 2147483647},
				output: []int32{-2147483648, 2147483647}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int8ArrayFromInt64Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int8ArrayToInt64Slice{Val: new([]int64)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []int64{-9223372036854775808, 9223372036854775807},
				output: []int64{-9223372036854775808, 9223372036854775807}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int8ArrayFromUintSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int8ArrayToUintSlice{Val: new([]uint)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []uint{0, 2147483647},
				output: []uint{0, 2147483647}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int8ArrayFromUint8Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int8ArrayToUint8Slice{Val: new([]uint8)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []uint8{0, 255},
				output: []uint8{0, 255}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int8ArrayFromUint16Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int8ArrayToUint16Slice{Val: new([]uint16)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []uint16{0, 65535},
				output: []uint16{0, 65535}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int8ArrayFromUint32Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int8ArrayToUint32Slice{Val: new([]uint32)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []uint32{0, 4294967295},
				output: []uint32{0, 4294967295}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int8ArrayFromUint64Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int8ArrayToUint64Slice{Val: new([]uint64)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []uint64{0, 9223372036854775807},
				output: []uint64{0, 9223372036854775807}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int8ArrayFromFloat32Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int8ArrayToFloat32Slice{Val: new([]float32)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []float32{0, 2147483647.0},
				output: []float32{0, 2147483647.0}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int8ArrayFromFloat64Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int8ArrayToFloat64Slice{Val: new([]float64)}
			return s, s.Val
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
			s := new(string)
			return s, s
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
			s := new([]byte)
			return s, s
		},
		data: []testdata{
			{
				input:  []byte("{-9223372036854775808,9223372036854775807}"),
				output: []byte(`{-9223372036854775808,9223372036854775807}`)},
		},
	}}.execute(t, "int8arr")
}
