package pg2go

import (
	"testing"
)

func TestNumRange(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(NumRangeFromIntArray2)
		},
		scanner: func() (interface{}, interface{}) {
			s := NumRangeToIntArray2{Val: new([2]int)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [2]int{-9223372036854775808, 9223372036854775807},
				output: [2]int{-9223372036854775808, 9223372036854775807}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumRangeFromInt8Array2)
		},
		scanner: func() (interface{}, interface{}) {
			s := NumRangeToInt8Array2{Val: new([2]int8)}
			return s, s.Val
		},
		data: []testdata{
			{input: [2]int8{-128, 127}, output: [2]int8{-128, 127}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumRangeFromInt16Array2)
		},
		scanner: func() (interface{}, interface{}) {
			s := NumRangeToInt16Array2{Val: new([2]int16)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [2]int16{-32768, 32767},
				output: [2]int16{-32768, 32767}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumRangeFromInt32Array2)
		},
		scanner: func() (interface{}, interface{}) {
			s := NumRangeToInt32Array2{Val: new([2]int32)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [2]int32{-2147483648, 2147483647},
				output: [2]int32{-2147483648, 2147483647}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumRangeFromInt64Array2)
		},
		scanner: func() (interface{}, interface{}) {
			s := NumRangeToInt64Array2{Val: new([2]int64)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [2]int64{-9223372036854775808, 9223372036854775807},
				output: [2]int64{-9223372036854775808, 9223372036854775807}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumRangeFromUintArray2)
		},
		scanner: func() (interface{}, interface{}) {
			s := NumRangeToUintArray2{Val: new([2]uint)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [2]uint{0, 9223372036854775807},
				output: [2]uint{0, 9223372036854775807}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumRangeFromUint8Array2)
		},
		scanner: func() (interface{}, interface{}) {
			s := NumRangeToUint8Array2{Val: new([2]uint8)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [2]uint8{0, 255},
				output: [2]uint8{0, 255}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumRangeFromUint16Array2)
		},
		scanner: func() (interface{}, interface{}) {
			s := NumRangeToUint16Array2{Val: new([2]uint16)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [2]uint16{0, 65535},
				output: [2]uint16{0, 65535}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumRangeFromUint32Array2)
		},
		scanner: func() (interface{}, interface{}) {
			s := NumRangeToUint32Array2{Val: new([2]uint32)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [2]uint32{0, 4294967295},
				output: [2]uint32{0, 4294967295}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumRangeFromUint64Array2)
		},
		scanner: func() (interface{}, interface{}) {
			s := NumRangeToUint64Array2{Val: new([2]uint64)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [2]uint64{0, 9223372036854775807},
				output: [2]uint64{0, 9223372036854775807}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumRangeFromFloat32Array2)
		},
		scanner: func() (interface{}, interface{}) {
			s := NumRangeToFloat32Array2{Val: new([2]float32)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [2]float32{-9223372036854775808.0, 922337203685477580.0},
				output: [2]float32{-9223372036854775808.0, 922337203685477580.0}},
		},
	}, {
		valuer: func() interface{} {
			return new(NumRangeFromFloat64Array2)
		},
		scanner: func() (interface{}, interface{}) {
			s := NumRangeToFloat64Array2{Val: new([2]float64)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  [2]float64{-9223372036854775808.0, 922337203685477580.0},
				output: [2]float64{-9223372036854775808.0, 922337203685477580.0}},
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
				input:  string("[-9223372036854775808,9223372036854775807)"),
				output: string(`[-9223372036854775808,9223372036854775807)`)},
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
				input:  []byte("[-9223372036854775808,9223372036854775807)"),
				output: []byte(`[-9223372036854775808,9223372036854775807)`)},
		},
	}}.execute(t, "numrange")
}
