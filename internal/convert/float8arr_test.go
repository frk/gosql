package convert

import (
	"testing"
)

func TestFloat8Array_Valuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return new(Float8ArrayFromFloat32Slice)
		},
		rows: []test_valuer_row{
			{typ: "float8arr", in: nil, want: nil},
			{typ: "float8arr", in: []float32{}, want: strptr(`{}`)},
			{typ: "float8arr", in: []float32{1, 0}, want: strptr(`{1,0}`)},
			{
				typ:  "float8arr",
				in:   []float32{3.14, 0.15},
				want: strptr(`{3.1400001049041748,0.15000000596046448}`)},
		},
	}, {
		valuer: func() interface{} {
			return new(Float8ArrayFromFloat64Slice)
		},
		rows: []test_valuer_row{
			{typ: "float8arr", in: nil, want: nil},
			{typ: "float8arr", in: []float64{}, want: strptr(`{}`)},
			{typ: "float8arr", in: []float64{1, 0}, want: strptr(`{1,0}`)},
			{
				typ:  "float8arr",
				in:   []float64{3.14, 0.15},
				want: strptr(`{3.1400000000000001,0.14999999999999999}`)},
		},
	}}.execute(t)
}

func TestFloat8Array_Scanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			s := Float8ArrayToFloat32Slice{V: new([]float32)}
			return s, s.V
		},
		rows: []test_scanner_row{
			{typ: "float8arr", in: nil, want: new([]float32)},
			{typ: "float8arr", in: `{}`, want: &[]float32{}},
			{typ: "float8arr", in: `{0.1}`, want: &[]float32{0.1}},
			{typ: "float8arr", in: `{1.9}`, want: &[]float32{1.9}},
			{
				typ:  "float8arr",
				in:   `{3.4,5.6,3.14159}`,
				want: &[]float32{3.4, 5.6, 3.14159}},
			{
				typ:  "float8arr",
				in:   `{0.0024,1.4,-89.2345,0.0}`,
				want: &[]float32{0.0024, 1.4, -89.2345, 0.0}},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			s := FloatArr2Float64Slice{Ptr: new([]float64)}
			return s, s.Ptr
		},
		rows: []test_scanner_row{
			{typ: "float8arr", in: nil, want: new([]float64)},
			{typ: "float8arr", in: `{}`, want: &[]float64{}},
			{typ: "float8arr", in: `{0.1}`, want: &[]float64{0.1}},
			{typ: "float8arr", in: `{1.9}`, want: &[]float64{1.9}},
			{
				typ:  "float8arr",
				in:   `{3.4,5.6,3.14159}`,
				want: &[]float64{3.4, 5.6, 3.14159}},
			{
				typ:  "float8arr",
				in:   `{0.0024,1.4,-89.2345,0.0}`,
				want: &[]float64{0.0024, 1.4, -89.2345, 0.0}},
		},
	}}.execute(t)
}

func TestFloat8Array_NoValuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return nil // string
		},
		rows: []test_valuer_row{
			{typ: "float8arr", in: nil, want: nil},
			{typ: "float8arr", in: `{}`, want: strptr(`{}`)},
			{typ: "float8arr", in: `{1,0}`, want: strptr(`{1,0}`)},
			{
				typ:  "float8arr",
				in:   `{3.14,0.15}`,
				want: strptr(`{3.1400000000000001,0.14999999999999999}`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // string
		},
		rows: []test_valuer_row{
			{typ: "float8arr", in: nil, want: nil},
			{typ: "float8arr", in: []byte(`{}`), want: strptr(`{}`)},
			{typ: "float8arr", in: []byte(`{1,0}`), want: strptr(`{1,0}`)},
			{
				typ:  "float8arr",
				in:   []byte(`{3.14,0.15}`),
				want: strptr(`{3.1400000000000001,0.14999999999999999}`)},
		},
	}}.execute(t)
}

func TestFloat8Array_NoScanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			v := new(string)
			return v, v
		},
		rows: []test_scanner_row{
			{typ: "float8arr", in: `{}`, want: strptr(`{}`)},
			{typ: "float8arr", in: `{0.1}`, want: strptr(`{0.10000000000000001}`)},
			{typ: "float8arr", in: `{0.9}`, want: strptr(`{0.90000000000000002}`)},
			{
				typ:  "float8arr",
				in:   `{3.4,5.6,3.14159}`,
				want: strptr(`{3.3999999999999999,5.5999999999999996,3.1415899999999999}`)},
			{
				typ:  "float8arr",
				in:   `{0.0024,1.4,-89.2345,0.0}`,
				want: strptr(`{0.0023999999999999998,1.3999999999999999,-89.234499999999997,0}`)},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			v := new([]byte)
			return v, v
		},
		rows: []test_scanner_row{
			{typ: "float8arr", in: `{}`, want: bytesptr(`{}`)},
			{typ: "float8arr", in: `{0.1}`, want: bytesptr(`{0.10000000000000001}`)},
			{typ: "float8arr", in: `{0.9}`, want: bytesptr(`{0.90000000000000002}`)},
			{
				typ:  "float8arr",
				in:   `{3.4,5.6,3.14159}`,
				want: bytesptr(`{3.3999999999999999,5.5999999999999996,3.1415899999999999}`)},
			{
				typ:  "float8arr",
				in:   `{0.0024,1.4,-89.2345,0.0}`,
				want: bytesptr(`{0.0023999999999999998,1.3999999999999999,-89.234499999999997,0}`)},
		},
	}}.execute(t)
}
