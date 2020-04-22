package pgsql

import (
	"testing"
)

func TestInt2VectorArray(t *testing.T) {
	testlist2{{
		valuer:  Int2VectorArrayFromIntSliceSlice,
		scanner: Int2VectorArrayToIntSliceSlice,
		data: []testdata{
			{
				input:  [][]int{{-32768, 32767}, {0, 1, 2, 3}},
				output: [][]int{{-32768, 32767}, {0, 1, 2, 3}}},
		},
	}, {
		valuer:  Int2VectorArrayFromInt8SliceSlice,
		scanner: Int2VectorArrayToInt8SliceSlice,
		data: []testdata{
			{
				input:  [][]int8{{-128, 127}, {0, 1, 2, 3}},
				output: [][]int8{{-128, 127}, {0, 1, 2, 3}}},
		},
	}, {
		valuer:  Int2VectorArrayFromInt16SliceSlice,
		scanner: Int2VectorArrayToInt16SliceSlice,
		data: []testdata{
			{
				input:  [][]int16{{-32768, 32767}, {0, 1, 2, 3}},
				output: [][]int16{{-32768, 32767}, {0, 1, 2, 3}}},
		},
	}, {
		valuer:  Int2VectorArrayFromInt32SliceSlice,
		scanner: Int2VectorArrayToInt32SliceSlice,
		data: []testdata{
			{
				input:  [][]int32{{-32768, 32767}, {0, 1, 2, 3}},
				output: [][]int32{{-32768, 32767}, {0, 1, 2, 3}}},
		},
	}, {
		valuer:  Int2VectorArrayFromInt64SliceSlice,
		scanner: Int2VectorArrayToInt64SliceSlice,
		data: []testdata{
			{
				input:  [][]int64{{-32768, 32767}, {0, 1, 2, 3}},
				output: [][]int64{{-32768, 32767}, {0, 1, 2, 3}}},
		},
	}, {
		valuer:  Int2VectorArrayFromUintSliceSlice,
		scanner: Int2VectorArrayToUintSliceSlice,
		data: []testdata{
			{
				input:  [][]uint{{0, 32767}, {0, 1, 2, 3}},
				output: [][]uint{{0, 32767}, {0, 1, 2, 3}}},
		},
	}, {
		valuer:  Int2VectorArrayFromUint8SliceSlice,
		scanner: Int2VectorArrayToUint8SliceSlice,
		data: []testdata{
			{
				input:  [][]uint8{{0, 255}, {0, 1, 2, 3}},
				output: [][]uint8{{0, 255}, {0, 1, 2, 3}}},
		},
	}, {
		valuer:  Int2VectorArrayFromUint16SliceSlice,
		scanner: Int2VectorArrayToUint16SliceSlice,
		data: []testdata{
			{
				input:  [][]uint16{{0, 32767}, {0, 1, 2, 3}},
				output: [][]uint16{{0, 32767}, {0, 1, 2, 3}}},
		},
	}, {
		valuer:  Int2VectorArrayFromUint32SliceSlice,
		scanner: Int2VectorArrayToUint32SliceSlice,
		data: []testdata{
			{
				input:  [][]uint32{{0, 32767}, {0, 1, 2, 3}},
				output: [][]uint32{{0, 32767}, {0, 1, 2, 3}}},
		},
	}, {
		valuer:  Int2VectorArrayFromUint64SliceSlice,
		scanner: Int2VectorArrayToUint64SliceSlice,
		data: []testdata{
			{
				input:  [][]uint64{{0, 32767}, {0, 1, 2, 3}},
				output: [][]uint64{{0, 32767}, {0, 1, 2, 3}}},
		},
	}, {
		valuer:  Int2VectorArrayFromFloat32SliceSlice,
		scanner: Int2VectorArrayToFloat32SliceSlice,
		data: []testdata{
			{
				input:  [][]float32{{-32768.0, 32767.0}, {0.0, 1.0, 2.0, 3.0}},
				output: [][]float32{{-32768.0, 32767.0}, {0.0, 1.0, 2.0, 3.0}}},
		},
	}, {
		valuer:  Int2VectorArrayFromFloat64SliceSlice,
		scanner: Int2VectorArrayToFloat64SliceSlice,
		data: []testdata{
			{
				input:  [][]float64{{-32768.0, 32767.0}, {0.0, 1.0, 2.0, 3.0}},
				output: [][]float64{{-32768.0, 32767.0}, {0.0, 1.0, 2.0, 3.0}}},
		},
	}, {
		data: []testdata{
			{
				input:  string(`{"-32768 32767","0 1 2 3"}`),
				output: string(`{"-32768 32767","0 1 2 3"}`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte(`{"-32768 32767","0 1 2 3"}`),
				output: []byte(`{"-32768 32767","0 1 2 3"}`)},
		},
	}}.execute(t, "int2vectorarr")
}
