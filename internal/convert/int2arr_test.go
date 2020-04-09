package convert

import (
	"testing"
)

func TestInt2Array(t *testing.T) {
	test_valuer_with_scanner{{
		valuer: func() interface{} {
			return new(Int2ArrayFromIntSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int2ArrayToIntSlice{Val: new([]int)}
			return s, s.Val
		},
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2arr",
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
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2arr",
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
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2arr",
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
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2arr",
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
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2arr",
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
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2arr",
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
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2arr",
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
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2arr",
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
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2arr",
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
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2arr",
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
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2arr",
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
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2arr",
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
				typ:    "int2arr",
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
		rows: []test_valuer_with_scanner_row{
			{
				typ:    "int2arr",
				input:  []byte("{-32768,32767}"),
				output: bytesptr(`{-32768,32767}`)},
		},
	}}.execute(t)
}
