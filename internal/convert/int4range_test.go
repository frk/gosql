package convert

import (
	"testing"
)

func TestInt4Range(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(Int4RangeFromIntArray2)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int4RangeToIntArray2{Val: new([2]int)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [2]int{-2147483648, 2147483647},
				output: [2]int{-2147483648, 2147483647}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int4RangeFromInt8Array2)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int4RangeToInt8Array2{Val: new([2]int8)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [2]int8{-128, 127},
				output: [2]int8{-128, 127}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int4RangeFromInt16Array2)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int4RangeToInt16Array2{Val: new([2]int16)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [2]int16{-32768, 32767},
				output: [2]int16{-32768, 32767}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int4RangeFromInt32Array2)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int4RangeToInt32Array2{Val: new([2]int32)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [2]int32{-2147483648, 2147483647},
				output: [2]int32{-2147483648, 2147483647}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int4RangeFromInt64Array2)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int4RangeToInt64Array2{Val: new([2]int64)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [2]int64{-2147483648, 2147483647},
				output: [2]int64{-2147483648, 2147483647}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int4RangeFromUintArray2)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int4RangeToUintArray2{Val: new([2]uint)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [2]uint{0, 2147483647},
				output: [2]uint{0, 2147483647}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int4RangeFromUint8Array2)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int4RangeToUint8Array2{Val: new([2]uint8)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [2]uint8{0, 255},
				output: [2]uint8{0, 255}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int4RangeFromUint16Array2)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int4RangeToUint16Array2{Val: new([2]uint16)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [2]uint16{0, 65535},
				output: [2]uint16{0, 65535}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int4RangeFromUint32Array2)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int4RangeToUint32Array2{Val: new([2]uint32)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [2]uint32{0, 2147483647},
				output: [2]uint32{0, 2147483647}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int4RangeFromUint64Array2)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int4RangeToUint64Array2{Val: new([2]uint64)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [2]uint64{0, 2147483647},
				output: [2]uint64{0, 2147483647}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int4RangeFromFloat32Array2)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int4RangeToFloat32Array2{Val: new([2]float32)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [2]float32{-2147483648.0, 214748364.0},
				output: [2]float32{-2147483648.0, 214748364.0}},
		},
	}, {
		valuer: func() interface{} {
			return new(Int4RangeFromFloat64Array2)
		},
		scanner: func() (interface{}, interface{}) {
			s := Int4RangeToFloat64Array2{Val: new([2]float64)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [2]float64{-2147483648.0, 2147483647.0},
				output: [2]float64{-2147483648.0, 2147483647.0}},
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
				input:  string("[-2147483648,2147483647)"),
				output: string(`[-2147483648,2147483647)`)},
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
				input:  []byte("[-2147483648,2147483647)"),
				output: []byte(`[-2147483648,2147483647)`)},
		},
	}}.execute(t, "int4range")
}
