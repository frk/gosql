package pgsql

import (
	"testing"
)

func TestInt2Vector(t *testing.T) {
	testlist2{{
		valuer:  Int2VectorFromIntSlice,
		scanner: Int2VectorToIntSlice,
		data: []testdata{
			{
				input:  []int{-32768, 32767},
				output: []int{-32768, 32767}},
		},
	}, {
		valuer:  Int2VectorFromInt8Slice,
		scanner: Int2VectorToInt8Slice,
		data: []testdata{
			{
				input:  []int8{-128, 127},
				output: []int8{-128, 127}},
		},
	}, {
		valuer:  Int2VectorFromInt16Slice,
		scanner: Int2VectorToInt16Slice,
		data: []testdata{
			{
				input:  []int16{-32768, 32767},
				output: []int16{-32768, 32767}},
		},
	}, {
		valuer:  Int2VectorFromInt32Slice,
		scanner: Int2VectorToInt32Slice,
		data: []testdata{
			{
				input:  []int32{-32768, 32767},
				output: []int32{-32768, 32767}},
		},
	}, {
		valuer:  Int2VectorFromInt64Slice,
		scanner: Int2VectorToInt64Slice,
		data: []testdata{
			{
				input:  []int64{-32768, 32767},
				output: []int64{-32768, 32767}},
		},
	}, {
		valuer:  Int2VectorFromUintSlice,
		scanner: Int2VectorToUintSlice,
		data: []testdata{
			{
				input:  []uint{0, 32767},
				output: []uint{0, 32767}},
		},
	}, {
		valuer:  Int2VectorFromUint8Slice,
		scanner: Int2VectorToUint8Slice,
		data: []testdata{
			{
				input:  []uint8{0, 255},
				output: []uint8{0, 255}},
		},
	}, {
		valuer:  Int2VectorFromUint16Slice,
		scanner: Int2VectorToUint16Slice,
		data: []testdata{
			{
				input:  []uint16{0, 32767},
				output: []uint16{0, 32767}},
		},
	}, {
		valuer:  Int2VectorFromUint32Slice,
		scanner: Int2VectorToUint32Slice,
		data: []testdata{
			{
				input:  []uint32{0, 32767},
				output: []uint32{0, 32767}},
		},
	}, {
		valuer:  Int2VectorFromUint64Slice,
		scanner: Int2VectorToUint64Slice,
		data: []testdata{
			{
				input:  []uint64{0, 32767},
				output: []uint64{0, 32767}},
		},
	}, {
		valuer:  Int2VectorFromFloat32Slice,
		scanner: Int2VectorToFloat32Slice,
		data: []testdata{
			{
				input:  []float32{0, 32767.0},
				output: []float32{0, 32767.0}},
		},
	}, {
		valuer:  Int2VectorFromFloat64Slice,
		scanner: Int2VectorToFloat64Slice,
		data: []testdata{
			{
				input:  []float64{0, 32767.0},
				output: []float64{0, 32767.0}},
		},
	}, {
		data: []testdata{
			{
				input:  string("-32768 32767"),
				output: string(`-32768 32767`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte("-32768 32767"),
				output: []byte(`-32768 32767`)},
		},
	}}.execute(t, "int2vector")
}
