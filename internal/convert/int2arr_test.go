package convert

import (
	"testing"
)

func TestInt2Array(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(Int2ArrayFromIntSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2ArrayToIntSlice{Val: new([]int)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []int{-32768, 32767},
				output: &[]int{-32768, 32767}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2ArrayFromInt8Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2ArrayToInt8Slice{Val: new([]int8)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []int8{-128, 127},
				output: &[]int8{-128, 127}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2ArrayFromInt16Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2ArrayToInt16Slice{Val: new([]int16)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []int16{-32768, 32767},
				output: &[]int16{-32768, 32767}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2ArrayFromInt32Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2ArrayToInt32Slice{Val: new([]int32)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []int32{-32768, 32767},
				output: &[]int32{-32768, 32767}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2ArrayFromInt64Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2ArrayToInt64Slice{Val: new([]int64)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []int64{-32768, 32767},
				output: &[]int64{-32768, 32767}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2ArrayFromUintSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2ArrayToUintSlice{Val: new([]uint)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []uint{0, 32767},
				output: &[]uint{0, 32767}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2ArrayFromUint8Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2ArrayToUint8Slice{Val: new([]uint8)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []uint8{0, 255},
				output: &[]uint8{0, 255}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2ArrayFromUint16Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2ArrayToUint16Slice{Val: new([]uint16)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []uint16{0, 32767},
				output: &[]uint16{0, 32767}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2ArrayFromUint32Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2ArrayToUint32Slice{Val: new([]uint32)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []uint32{0, 32767},
				output: &[]uint32{0, 32767}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2ArrayFromUint64Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2ArrayToUint64Slice{Val: new([]uint64)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []uint64{0, 32767},
				output: &[]uint64{0, 32767}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2ArrayFromFloat32Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2ArrayToFloat32Slice{Val: new([]float32)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []float32{0, 32767.0},
				output: &[]float32{0, 32767.0}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2ArrayFromFloat64Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2ArrayToFloat64Slice{Val: new([]float64)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  []float64{0, 32767.0},
				output: &[]float64{0, 32767.0}},
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
				input:  "{-32768,32767}",
				output: strptr(`{-32768,32767}`)},
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
				input:  []byte("{-32768,32767}"),
				output: bytesptr(`{-32768,32767}`)},
		},
	}}.execute(t, "int2arr")
}
