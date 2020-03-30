package convert

import (
	"testing"
)

func TestFloat4Array_Valuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return new(Float4ArrayFromFloat32Slice)
		},
		rows: []test_valuer_row{
			{typ: "float4arr", in: nil, want: nil},
			{typ: "float4arr", in: []float32{}, want: strptr(`{}`)},
			{typ: "float4arr", in: []float32{1, 0}, want: strptr(`{1,0}`)},
			{
				typ:  "float4arr",
				in:   []float32{3.14, 0.15},
				want: strptr(`{3.1400001,0.15000001}`)},
		},
	}, {
		valuer: func() interface{} {
			return new(Float4ArrayFromFloat64Slice)
		},
		rows: []test_valuer_row{
			{typ: "float4arr", in: nil, want: nil},
			{typ: "float4arr", in: []float64{}, want: strptr(`{}`)},
			{typ: "float4arr", in: []float64{1, 0}, want: strptr(`{1,0}`)},
			{
				typ:  "float4arr",
				in:   []float64{3.14, 0.15},
				want: strptr(`{3.1400001,0.15000001}`)},
		},
	}}.execute(t)
}

func TestFloat4Array_Scanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			s := Float4ArrayToFloat32Slice{V: new([]float32)}
			return s, s.V
		},
		rows: []test_scanner_row{
			{typ: "float4arr", in: nil, want: new([]float32)},
			{typ: "float4arr", in: `{}`, want: &[]float32{}},
			{typ: "float4arr", in: `{0.1}`, want: &[]float32{0.1}},
			{typ: "float4arr", in: `{1.9}`, want: &[]float32{1.9}},
			{
				typ:  "float4arr",
				in:   `{3.4,5.6,3.14159}`,
				want: &[]float32{3.4, 5.6, 3.14159}},
			{
				typ:  "float4arr",
				in:   `{0.0024,1.4,-89.2345,0.0}`,
				want: &[]float32{0.0024, 1.4, -89.2345, 0.0}},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			s := Float4ArrayToFloat64Slice{V: new([]float64)}
			return s, s.V
		},
		rows: []test_scanner_row{
			{typ: "float4arr", in: nil, want: new([]float64)},
			{typ: "float4arr", in: `{}`, want: &[]float64{}},
			{typ: "float4arr", in: `{0.1}`, want: &[]float64{0.1}},
			{typ: "float4arr", in: `{1.9}`, want: &[]float64{1.9}},
			{
				typ:  "float4arr",
				in:   `{3.4,5.6,3.14159}`,
				want: &[]float64{3.4000001, 5.5999999, 3.1415901}},
			{
				typ:  "float4arr",
				in:   `{0.0024,1.4,-89.2345,0.0}`,
				want: &[]float64{0.0024000001, 1.4, -89.234497, 0.0}},
		},
	}}.execute(t)
}

func TestFloat4Array_NoValuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return nil // string
		},
		rows: []test_valuer_row{
			{typ: "float4arr", in: nil, want: nil},
			{typ: "float4arr", in: `{}`, want: strptr(`{}`)},
			{typ: "float4arr", in: `{1,0}`, want: strptr(`{1,0}`)},
			{
				typ:  "float4arr",
				in:   `{3.14,0.15}`,
				want: strptr(`{3.1400001,0.15000001}`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // string
		},
		rows: []test_valuer_row{
			{typ: "float4arr", in: nil, want: nil},
			{typ: "float4arr", in: []byte(`{}`), want: strptr(`{}`)},
			{typ: "float4arr", in: []byte(`{1,0}`), want: strptr(`{1,0}`)},
			{
				typ:  "float4arr",
				in:   []byte(`{3.14,0.15}`),
				want: strptr(`{3.1400001,0.15000001}`)},
		},
	}}.execute(t)
}

func TestFloat4Array_NoScanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			v := new(string)
			return v, v
		},
		rows: []test_scanner_row{
			{typ: "float4arr", in: `{}`, want: strptr(`{}`)},
			{typ: "float4arr", in: `{0.1}`, want: strptr(`{0.1}`)},
			{typ: "float4arr", in: `{0.9}`, want: strptr(`{0.89999998}`)},
			{
				typ:  "float4arr",
				in:   `{3.4,5.6,3.14159}`,
				want: strptr(`{3.4000001,5.5999999,3.1415901}`)},
			{
				typ:  "float4arr",
				in:   `{0.0024,1.4,-89.2345,0.0}`,
				want: strptr(`{0.0024000001,1.4,-89.234497,0}`)},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			v := new([]byte)
			return v, v
		},
		rows: []test_scanner_row{
			{typ: "float4arr", in: `{}`, want: bytesptr(`{}`)},
			{typ: "float4arr", in: `{0.1}`, want: bytesptr(`{0.1}`)},
			{typ: "float4arr", in: `{0.9}`, want: bytesptr(`{0.89999998}`)},
			{
				typ:  "float4arr",
				in:   `{3.4,5.6,3.14159}`,
				want: bytesptr(`{3.4000001,5.5999999,3.1415901}`)},
			{
				typ:  "float4arr",
				in:   `{0.0024,1.4,-89.2345,0.0}`,
				want: bytesptr(`{0.0024000001,1.4,-89.234497,0}`)},
		},
	}}.execute(t)
}
