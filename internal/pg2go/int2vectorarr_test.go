package pg2go

import (
	"testing"
)

func TestInt2VectorArray(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(Int2VectorArrayFromIntSliceSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorArrayToIntSliceSlice{Val: new([][]int)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [][]int{{-32768, 32767}, {0, 1, 2, 3}},
				output: [][]int{{-32768, 32767}, {0, 1, 2, 3}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorArrayFromInt8SliceSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorArrayToInt8SliceSlice{Val: new([][]int8)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [][]int8{{-128, 127}, {0, 1, 2, 3}},
				output: [][]int8{{-128, 127}, {0, 1, 2, 3}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorArrayFromInt16SliceSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorArrayToInt16SliceSlice{Val: new([][]int16)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [][]int16{{-32768, 32767}, {0, 1, 2, 3}},
				output: [][]int16{{-32768, 32767}, {0, 1, 2, 3}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorArrayFromInt32SliceSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorArrayToInt32SliceSlice{Val: new([][]int32)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [][]int32{{-32768, 32767}, {0, 1, 2, 3}},
				output: [][]int32{{-32768, 32767}, {0, 1, 2, 3}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorArrayFromInt64SliceSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorArrayToInt64SliceSlice{Val: new([][]int64)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [][]int64{{-32768, 32767}, {0, 1, 2, 3}},
				output: [][]int64{{-32768, 32767}, {0, 1, 2, 3}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorArrayFromUintSliceSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorArrayToUintSliceSlice{Val: new([][]uint)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [][]uint{{0, 32767}, {0, 1, 2, 3}},
				output: [][]uint{{0, 32767}, {0, 1, 2, 3}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorArrayFromUint8SliceSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorArrayToUint8SliceSlice{Val: new([][]uint8)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [][]uint8{{0, 255}, {0, 1, 2, 3}},
				output: [][]uint8{{0, 255}, {0, 1, 2, 3}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorArrayFromUint16SliceSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorArrayToUint16SliceSlice{Val: new([][]uint16)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [][]uint16{{0, 32767}, {0, 1, 2, 3}},
				output: [][]uint16{{0, 32767}, {0, 1, 2, 3}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorArrayFromUint32SliceSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorArrayToUint32SliceSlice{Val: new([][]uint32)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [][]uint32{{0, 32767}, {0, 1, 2, 3}},
				output: [][]uint32{{0, 32767}, {0, 1, 2, 3}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorArrayFromUint64SliceSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorArrayToUint64SliceSlice{Val: new([][]uint64)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [][]uint64{{0, 32767}, {0, 1, 2, 3}},
				output: [][]uint64{{0, 32767}, {0, 1, 2, 3}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorArrayFromFloat32SliceSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorArrayToFloat32SliceSlice{Val: new([][]float32)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [][]float32{{-32768.0, 32767.0}, {0.0, 1.0, 2.0, 3.0}},
				output: [][]float32{{-32768.0, 32767.0}, {0.0, 1.0, 2.0, 3.0}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorArrayFromFloat64SliceSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorArrayToFloat64SliceSlice{Val: new([][]float64)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [][]float64{{-32768.0, 32767.0}, {0.0, 1.0, 2.0, 3.0}},
				output: [][]float64{{-32768.0, 32767.0}, {0.0, 1.0, 2.0, 3.0}}},
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
				input:  string(`{"-32768 32767","0 1 2 3"}`),
				output: string(`{"-32768 32767","0 1 2 3"}`)},
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
				input:  []byte(`{"-32768 32767","0 1 2 3"}`),
				output: []byte(`{"-32768 32767","0 1 2 3"}`)},
		},
	}}.execute(t, "int2vectorarr")
}
