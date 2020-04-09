package convert

import (
	"testing"
)

func TestInt2VectorArray(t *testing.T) {
	test_valuer_with_scanner{{
		valuer: func() interface{} {
			return new(Int2VectorArrayFromIntSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorArrayToIntSlice{Val: new([][]int)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vectorarr",
				input:  [][]int{{-32768, 32767}, {0, 1, 2, 3}},
				output: &[][]int{{-32768, 32767}, {0, 1, 2, 3}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorArrayFromInt8Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorArrayToInt8Slice{Val: new([][]int8)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vectorarr",
				input:  [][]int8{{-128, 127}, {0, 1, 2, 3}},
				output: &[][]int8{{-128, 127}, {0, 1, 2, 3}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorArrayFromInt16Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorArrayToInt16Slice{Val: new([][]int16)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vectorarr",
				input:  [][]int16{{-32768, 32767}, {0, 1, 2, 3}},
				output: &[][]int16{{-32768, 32767}, {0, 1, 2, 3}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorArrayFromInt32Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorArrayToInt32Slice{Val: new([][]int32)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vectorarr",
				input:  [][]int32{{-32768, 32767}, {0, 1, 2, 3}},
				output: &[][]int32{{-32768, 32767}, {0, 1, 2, 3}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorArrayFromInt64Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorArrayToInt64Slice{Val: new([][]int64)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vectorarr",
				input:  [][]int64{{-32768, 32767}, {0, 1, 2, 3}},
				output: &[][]int64{{-32768, 32767}, {0, 1, 2, 3}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorArrayFromUintSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorArrayToUintSlice{Val: new([][]uint)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vectorarr",
				input:  [][]uint{{0, 32767}, {0, 1, 2, 3}},
				output: &[][]uint{{0, 32767}, {0, 1, 2, 3}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorArrayFromUint8Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorArrayToUint8Slice{Val: new([][]uint8)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vectorarr",
				input:  [][]uint8{{0, 255}, {0, 1, 2, 3}},
				output: &[][]uint8{{0, 255}, {0, 1, 2, 3}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorArrayFromUint16Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorArrayToUint16Slice{Val: new([][]uint16)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vectorarr",
				input:  [][]uint16{{0, 32767}, {0, 1, 2, 3}},
				output: &[][]uint16{{0, 32767}, {0, 1, 2, 3}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorArrayFromUint32Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorArrayToUint32Slice{Val: new([][]uint32)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vectorarr",
				input:  [][]uint32{{0, 32767}, {0, 1, 2, 3}},
				output: &[][]uint32{{0, 32767}, {0, 1, 2, 3}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorArrayFromUint64Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorArrayToUint64Slice{Val: new([][]uint64)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vectorarr",
				input:  [][]uint64{{0, 32767}, {0, 1, 2, 3}},
				output: &[][]uint64{{0, 32767}, {0, 1, 2, 3}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorArrayFromFloat32Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorArrayToFloat32Slice{Val: new([][]float32)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vectorarr",
				input:  [][]float32{{-32768.0, 32767.0}, {0.0, 1.0, 2.0, 3.0}},
				output: &[][]float32{{-32768.0, 32767.0}, {0.0, 1.0, 2.0, 3.0}}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorArrayFromFloat64Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorArrayToFloat64Slice{Val: new([][]float64)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vectorarr",
				input:  [][]float64{{-32768.0, 32767.0}, {0.0, 1.0, 2.0, 3.0}},
				output: &[][]float64{{-32768.0, 32767.0}, {0.0, 1.0, 2.0, 3.0}}},
		},
	}, {
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			s := new(string)
			return s, s
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vectorarr",
				input:  `{"-32768 32767","0 1 2 3"}`,
				output: strptr(`{"-32768 32767","0 1 2 3"}`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			s := new([]byte)
			return s, s
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vectorarr",
				input:  []byte(`{"-32768 32767","0 1 2 3"}`),
				output: bytesptr(`{"-32768 32767","0 1 2 3"}`)},
		},
	}}.execute(t)
}
