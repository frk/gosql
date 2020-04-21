package pg2go

import (
	"testing"
)

func TestInt4Range(t *testing.T) {
	testlist2{{
		valuer:  Int4RangeFromIntArray2,
		scanner: Int4RangeToIntArray2,
		data: []testdata{
			{
				input:  [2]int{-2147483648, 2147483647},
				output: [2]int{-2147483648, 2147483647}},
		},
	}, {
		valuer:  Int4RangeFromInt8Array2,
		scanner: Int4RangeToInt8Array2,
		data: []testdata{
			{
				input:  [2]int8{-128, 127},
				output: [2]int8{-128, 127}},
		},
	}, {
		valuer:  Int4RangeFromInt16Array2,
		scanner: Int4RangeToInt16Array2,
		data: []testdata{
			{
				input:  [2]int16{-32768, 32767},
				output: [2]int16{-32768, 32767}},
		},
	}, {
		valuer:  Int4RangeFromInt32Array2,
		scanner: Int4RangeToInt32Array2,
		data: []testdata{
			{
				input:  [2]int32{-2147483648, 2147483647},
				output: [2]int32{-2147483648, 2147483647}},
		},
	}, {
		valuer:  Int4RangeFromInt64Array2,
		scanner: Int4RangeToInt64Array2,
		data: []testdata{
			{
				input:  [2]int64{-2147483648, 2147483647},
				output: [2]int64{-2147483648, 2147483647}},
		},
	}, {
		valuer:  Int4RangeFromUintArray2,
		scanner: Int4RangeToUintArray2,
		data: []testdata{
			{
				input:  [2]uint{0, 2147483647},
				output: [2]uint{0, 2147483647}},
		},
	}, {
		valuer:  Int4RangeFromUint8Array2,
		scanner: Int4RangeToUint8Array2,
		data: []testdata{
			{
				input:  [2]uint8{0, 255},
				output: [2]uint8{0, 255}},
		},
	}, {
		valuer:  Int4RangeFromUint16Array2,
		scanner: Int4RangeToUint16Array2,
		data: []testdata{
			{
				input:  [2]uint16{0, 65535},
				output: [2]uint16{0, 65535}},
		},
	}, {
		valuer:  Int4RangeFromUint32Array2,
		scanner: Int4RangeToUint32Array2,
		data: []testdata{
			{
				input:  [2]uint32{0, 2147483647},
				output: [2]uint32{0, 2147483647}},
		},
	}, {
		valuer:  Int4RangeFromUint64Array2,
		scanner: Int4RangeToUint64Array2,
		data: []testdata{
			{
				input:  [2]uint64{0, 2147483647},
				output: [2]uint64{0, 2147483647}},
		},
	}, {
		valuer:  Int4RangeFromFloat32Array2,
		scanner: Int4RangeToFloat32Array2,
		data: []testdata{
			{
				input:  [2]float32{-2147483648.0, 214748364.0},
				output: [2]float32{-2147483648.0, 214748364.0}},
		},
	}, {
		valuer:  Int4RangeFromFloat64Array2,
		scanner: Int4RangeToFloat64Array2,
		data: []testdata{
			{
				input:  [2]float64{-2147483648.0, 2147483647.0},
				output: [2]float64{-2147483648.0, 2147483647.0}},
		},
	}, {
		data: []testdata{
			{
				input:  string("[-2147483648,2147483647)"),
				output: string(`[-2147483648,2147483647)`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte("[-2147483648,2147483647)"),
				output: []byte(`[-2147483648,2147483647)`)},
		},
	}}.execute(t, "int4range")
}
