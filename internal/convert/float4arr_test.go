package convert

import (
	"testing"
)

func TestFloat4Array(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(Float4ArrayFromFloat32Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Float4ArrayToFloat32Slice{Val: new([]float32)}
			return s, s.Val
		},
		data: []testdata{
			{input: nil, output: []float32(nil)},
			{input: []float32{}, output: []float32{}},
			{input: []float32{1, 0}, output: []float32{1, 0}},
			{
				input:  []float32{3.14, 0.15},
				output: []float32{3.14, 0.15}},
			{
				input:  []float32{3.4, 5.6, 3.14159},
				output: []float32{3.4, 5.6, 3.14159}},
			{
				input:  []float32{0.0024, 1.4, -89.2345, 0.0},
				output: []float32{0.0024, 1.4, -89.2345, 0.0}},
		},
	}, {
		valuer: func() interface{} {
			return new(Float4ArrayFromFloat64Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := Float4ArrayToFloat64Slice{Val: new([]float64)}
			return s, s.Val
		},
		data: []testdata{
			{input: nil, output: []float64(nil)},
			{input: []float64{}, output: []float64{}},
			{input: []float64{1, 0}, output: []float64{1, 0}},
			{
				input:  []float64{3.14, 0.15},
				output: []float64{3.1400001, 0.15000001}},
			{
				input:  []float64{3.4, 5.6, 3.14159},
				output: []float64{3.4000001, 5.5999999, 3.1415901}},
			{
				input:  []float64{0.0024, 1.4, -89.2345, 0.0},
				output: []float64{0.0024000001, 1.4, -89.234497, 0.0}},
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
			{input: string(`{}`), output: string(`{}`)},
			{input: string(`{1,0}`), output: string(`{1,0}`)},
			{
				input:  string(`{3.14,0.15}`),
				output: string(`{3.1400001,0.15000001}`)},
			{
				input:  string(`{3.4,5.6,3.14159}`),
				output: string(`{3.4000001,5.5999999,3.1415901}`)},
			{
				input:  string(`{0.0024,1.4,-89.2345,0.0}`),
				output: string(`{0.0024000001,1.4,-89.234497,0}`)},
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
			{input: []byte(`{}`), output: []byte(`{}`)},
			{input: []byte(`{1,0}`), output: []byte(`{1,0}`)},
			{
				input:  []byte(`{3.14,0.15}`),
				output: []byte(`{3.1400001,0.15000001}`)},
			{
				input:  []byte(`{3.4,5.6,3.14159}`),
				output: []byte(`{3.4000001,5.5999999,3.1415901}`)},
			{
				input:  []byte(`{0.0024,1.4,-89.2345,0.0}`),
				output: []byte(`{0.0024000001,1.4,-89.234497,0}`)},
		},
	}}.execute(t, "float4arr")
}
