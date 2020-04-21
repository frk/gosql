package pg2go

import (
	"testing"
)

func TestInt2Array(t *testing.T) {
	testlist2{{
		valuer:  Int2ArrayFromIntSlice,
		scanner: Int2ArrayToIntSlice,
		data: []testdata{
			{
				input:  []int{-32768, 32767},
				output: []int{-32768, 32767}},
		},
	}, {
		valuer:  Int2ArrayFromInt8Slice,
		scanner: Int2ArrayToInt8Slice,
		data: []testdata{
			{
				input:  []int8{-128, 127},
				output: []int8{-128, 127}},
		},
	}, {
		valuer:  Int2ArrayFromInt16Slice,
		scanner: Int2ArrayToInt16Slice,
		data: []testdata{
			{
				input:  []int16{-32768, 32767},
				output: []int16{-32768, 32767}},
		},
	}, {
		valuer:  Int2ArrayFromInt32Slice,
		scanner: Int2ArrayToInt32Slice,
		data: []testdata{
			{
				input:  []int32{-32768, 32767},
				output: []int32{-32768, 32767}},
		},
	}, {
		valuer:  Int2ArrayFromInt64Slice,
		scanner: Int2ArrayToInt64Slice,
		data: []testdata{
			{
				input:  []int64{-32768, 32767},
				output: []int64{-32768, 32767}},
		},
	}, {
		valuer:  Int2ArrayFromUintSlice,
		scanner: Int2ArrayToUintSlice,
		data: []testdata{
			{
				input:  []uint{0, 32767},
				output: []uint{0, 32767}},
		},
	}, {
		valuer:  Int2ArrayFromUint8Slice,
		scanner: Int2ArrayToUint8Slice,
		data: []testdata{
			{
				input:  []uint8{0, 255},
				output: []uint8{0, 255}},
		},
	}, {
		valuer:  Int2ArrayFromUint16Slice,
		scanner: Int2ArrayToUint16Slice,
		data: []testdata{
			{
				input:  []uint16{0, 32767},
				output: []uint16{0, 32767}},
		},
	}, {
		valuer:  Int2ArrayFromUint32Slice,
		scanner: Int2ArrayToUint32Slice,
		data: []testdata{
			{
				input:  []uint32{0, 32767},
				output: []uint32{0, 32767}},
		},
	}, {
		valuer:  Int2ArrayFromUint64Slice,
		scanner: Int2ArrayToUint64Slice,
		data: []testdata{
			{
				input:  []uint64{0, 32767},
				output: []uint64{0, 32767}},
		},
	}, {
		valuer:  Int2ArrayFromFloat32Slice,
		scanner: Int2ArrayToFloat32Slice,
		data: []testdata{
			{
				input:  []float32{0, 32767.0},
				output: []float32{0, 32767.0}},
		},
	}, {
		valuer:  Int2ArrayFromFloat64Slice,
		scanner: Int2ArrayToFloat64Slice,
		data: []testdata{
			{
				input:  []float64{0, 32767.0},
				output: []float64{0, 32767.0}},
		},
	}, {
		data: []testdata{
			{
				input:  string("{-32768,32767}"),
				output: string(`{-32768,32767}`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte("{-32768,32767}"),
				output: []byte(`{-32768,32767}`)},
		},
	}}.execute(t, "int2arr")
}
