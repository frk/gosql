package convert

import (
	"testing"
)

func TestFloat4(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return nil // float32
		},
		scanner: func() (interface{}, interface{}) {
			s := new(float32)
			return s, s
		},
		data: []testdata{
			{input: float32(0), output: f32ptr(0)},
			{input: float32(1), output: f32ptr(1)},
			{input: float32(3.14), output: f32ptr(3.14)},
			{input: float32(0.15), output: f32ptr(0.15)},
		},
	}, {
		valuer: func() interface{} {
			return nil // float64
		},
		scanner: func() (interface{}, interface{}) {
			s := new(float64)
			return s, s
		},
		data: []testdata{
			{input: float64(0), output: f64ptr(0)},
			{input: float64(1), output: f64ptr(1)},
			{input: float64(3.14), output: f64ptr(3.1400001)},
			{input: float64(0.15), output: f64ptr(0.15000001)},
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
			{input: string("0"), output: strptr(`0`)},
			{input: string("1"), output: strptr(`1`)},
			{input: string("3.14"), output: strptr(`3.1400001`)},
			{input: string("0.15"), output: strptr(`0.15000001`)},
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
			{input: []byte("0"), output: bytesptr(`0`)},
			{input: []byte("1"), output: bytesptr(`1`)},
			{input: []byte("3.14"), output: bytesptr(`3.1400001`)},
			{input: []byte("0.15"), output: bytesptr(`0.15000001`)},
		},
	}}.execute(t, "float4")
}
