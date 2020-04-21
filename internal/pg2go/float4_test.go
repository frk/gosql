package pg2go

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
			{input: float32(0), output: float32(0)},
			{input: float32(1), output: float32(1)},
			{input: float32(3.14), output: float32(3.14)},
			{input: float32(0.15), output: float32(0.15)},
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
			{input: float64(0), output: float64(0)},
			{input: float64(1), output: float64(1)},
			{input: float64(3.14), output: float64(3.1400001)},
			{input: float64(0.15), output: float64(0.15000001)},
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
			{input: string("0"), output: string(`0`)},
			{input: string("1"), output: string(`1`)},
			{input: string("3.14"), output: string(`3.1400001`)},
			{input: string("0.15"), output: string(`0.15000001`)},
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
			{input: []byte("0"), output: []byte(`0`)},
			{input: []byte("1"), output: []byte(`1`)},
			{input: []byte("3.14"), output: []byte(`3.1400001`)},
			{input: []byte("0.15"), output: []byte(`0.15000001`)},
		},
	}}.execute(t, "float4")
}
