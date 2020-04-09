package convert

import (
	"testing"
)

func TestInt2Vector(t *testing.T) {
	test_valuer_with_scanner{{
		valuer: func() interface{} {
			return new(Int2VectorFromIntSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorToIntSlice{Val: new([]int)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vector",
				input:  []int{-32768, 32767},
				output: &[]int{-32768, 32767}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorFromInt8Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorToInt8Slice{Val: new([]int8)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vector",
				input:  []int8{-128, 127},
				output: &[]int8{-128, 127}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorFromInt16Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorToInt16Slice{Val: new([]int16)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vector",
				input:  []int16{-32768, 32767},
				output: &[]int16{-32768, 32767}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorFromInt32Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorToInt32Slice{Val: new([]int32)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vector",
				input:  []int32{-32768, 32767},
				output: &[]int32{-32768, 32767}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorFromInt64Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorToInt64Slice{Val: new([]int64)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vector",
				input:  []int64{-32768, 32767},
				output: &[]int64{-32768, 32767}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorFromUintSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorToUintSlice{Val: new([]uint)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vector",
				input:  []uint{0, 32767},
				output: &[]uint{0, 32767}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorFromUint8Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorToUint8Slice{Val: new([]uint8)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vector",
				input:  []uint8{0, 255},
				output: &[]uint8{0, 255}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorFromUint16Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorToUint16Slice{Val: new([]uint16)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vector",
				input:  []uint16{0, 32767},
				output: &[]uint16{0, 32767}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorFromUint32Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorToUint32Slice{Val: new([]uint32)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vector",
				input:  []uint32{0, 32767},
				output: &[]uint32{0, 32767}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorFromUint64Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorToUint64Slice{Val: new([]uint64)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vector",
				input:  []uint64{0, 32767},
				output: &[]uint64{0, 32767}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorFromFloat32Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorToFloat32Slice{Val: new([]float32)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vector",
				input:  []float32{0, 32767.0},
				output: &[]float32{0, 32767.0}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int2VectorFromFloat64Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2VectorToFloat64Slice{Val: new([]float64)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vector",
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
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2vector",
				input:  "-32768 32767",
				output: strptr(`-32768 32767`)},
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
				typ:    "int2vector",
				input:  []byte("-32768 32767"),
				output: bytesptr(`-32768 32767`)},
		},
	}}.execute(t)
}
