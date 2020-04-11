package convert

import (
	"testing"
)

func TestInt4Array(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(Int4ArrayFromIntSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int4ArrayToIntSlice{Val: new([]int)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []int{-2147483648, 2147483647},
				output: []int{-2147483648, 2147483647}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int4ArrayFromInt8Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int4ArrayToInt8Slice{Val: new([]int8)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []int8{-128, 127},
				output: []int8{-128, 127}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int4ArrayFromInt16Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int4ArrayToInt16Slice{Val: new([]int16)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []int16{-32768, 32767},
				output: []int16{-32768, 32767}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int4ArrayFromInt32Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int4ArrayToInt32Slice{Val: new([]int32)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []int32{-2147483648, 2147483647},
				output: []int32{-2147483648, 2147483647}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int4ArrayFromInt64Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int4ArrayToInt64Slice{Val: new([]int64)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []int64{-2147483648, 2147483647},
				output: []int64{-2147483648, 2147483647}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int4ArrayFromUintSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int4ArrayToUintSlice{Val: new([]uint)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []uint{0, 2147483647},
				output: []uint{0, 2147483647}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int4ArrayFromUint8Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int4ArrayToUint8Slice{Val: new([]uint8)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []uint8{0, 255},
				output: []uint8{0, 255}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int4ArrayFromUint16Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int4ArrayToUint16Slice{Val: new([]uint16)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []uint16{0, 65535},
				output: []uint16{0, 65535}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int4ArrayFromUint32Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int4ArrayToUint32Slice{Val: new([]uint32)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []uint32{0, 2147483647},
				output: []uint32{0, 2147483647}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int4ArrayFromUint64Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int4ArrayToUint64Slice{Val: new([]uint64)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []uint64{0, 2147483647},
				output: []uint64{0, 2147483647}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int4ArrayFromFloat32Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int4ArrayToFloat32Slice{Val: new([]float32)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []float32{0, 214748364.0},
				output: []float32{0, 214748364.0}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int4ArrayFromFloat64Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int4ArrayToFloat64Slice{Val: new([]float64)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []float64{0, 2147483647.0},
				output: []float64{0, 2147483647.0}},
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
				input:  string("{-2147483648,2147483647}"),
				output: string(`{-2147483648,2147483647}`)},
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
				input:  []byte("{-2147483648,2147483647}"),
				output: []byte(`{-2147483648,2147483647}`)},
		},
	}}.execute(t, "int4arr")
}
