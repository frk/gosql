package convert

import (
	"testing"
)

func TestFloat8_NoValuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return nil // float32
		},
		rows: []test_valuer_row{
			{typ: "float8", in: nil, want: nil},
			{typ: "float8", in: float32(0), want: strptr(`0`)},
			{typ: "float8", in: float32(1), want: strptr(`1`)},
			{typ: "float8", in: float32(3.14), want: strptr(`3.140000104904175`)},
			{typ: "float8", in: float32(0.15), want: strptr(`0.15000000596046448`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // float64
		},
		rows: []test_valuer_row{
			{typ: "float8", in: nil, want: nil},
			{typ: "float8", in: float64(0), want: strptr(`0`)},
			{typ: "float8", in: float64(1), want: strptr(`1`)},
			{typ: "float8", in: float64(3.14), want: strptr(`3.14`)},
			{typ: "float8", in: float64(0.15), want: strptr(`0.15`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // string
		},
		rows: []test_valuer_row{
			{typ: "float8", in: nil, want: nil},
			{typ: "float8", in: "0", want: strptr(`0`)},
			{typ: "float8", in: "1", want: strptr(`1`)},
			{typ: "float8", in: "3.14", want: strptr(`3.14`)},
			{typ: "float8", in: "0.15", want: strptr(`0.15`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		rows: []test_valuer_row{
			{typ: "float8", in: nil, want: nil},
			{typ: "float8", in: []byte("0"), want: strptr(`0`)},
			{typ: "float8", in: []byte("1"), want: strptr(`1`)},
			{typ: "float8", in: []byte("3.14"), want: strptr(`3.14`)},
			{typ: "float8", in: []byte("0.15"), want: strptr(`0.15`)},
		},
	}}.execute(t)
}

func TestFloat8_NoScanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			d := new(float32)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "float8", in: `0`, want: f32ptr(0)},
			{typ: "float8", in: `1`, want: f32ptr(1)},
			{typ: "float8", in: `3.14`, want: f32ptr(3.14)},
			{typ: "float8", in: `0.15`, want: f32ptr(0.15)},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			d := new(float64)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "float8", in: `0`, want: f64ptr(0)},
			{typ: "float8", in: `1`, want: f64ptr(1)},
			{typ: "float8", in: `3.14`, want: f64ptr(3.14)},
			{typ: "float8", in: `0.15`, want: f64ptr(0.15)},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			d := new(string)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "float8", in: `0`, want: strptr(`0`)},
			{typ: "float8", in: `1`, want: strptr(`1`)},
			{typ: "float8", in: `3.14`, want: strptr(`3.14`)},
			{typ: "float8", in: `0.15`, want: strptr(`0.15`)},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			d := new([]byte)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "float8", in: `0`, want: bytesptr(`0`)},
			{typ: "float8", in: `1`, want: bytesptr(`1`)},
			{typ: "float8", in: `3.14`, want: bytesptr(`3.14`)},
			{typ: "float8", in: `0.15`, want: bytesptr(`0.15`)},
		},
	}}.execute(t)
}
