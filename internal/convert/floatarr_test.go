package convert

import (
	"testing"
)

func TestFloatArrScanners(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			s := FloatArr2Float32Slice{Ptr: new([]float32)}
			return s, s.Ptr
		},
		rows: []test_scanner_row{
			{typ: "float4arr", in: nil, want: new([]float32)},
			{typ: "float8arr", in: `{}`, want: &[]float32{}},
			{typ: "float4arr", in: `{0.1}`, want: &[]float32{0.1}},
			{typ: "float8arr", in: `{1.9}`, want: &[]float32{1.9}},
			{typ: "float4arr", in: `{3.4,5.6,3.14159}`, want: &[]float32{3.4, 5.6, 3.14159}},
			{typ: "float8arr", in: `{0.0024,1.4,-89.2345,0.0}`, want: &[]float32{0.0024, 1.4, -89.2345, 0.0}},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			s := FloatArr2Float64Slice{Ptr: new([]float64)}
			return s, s.Ptr
		},
		rows: []test_scanner_row{
			{typ: "float4arr", in: nil, want: new([]float64)},
			{typ: "float8arr", in: `{}`, want: &[]float64{}},
			{typ: "float4arr", in: `{0.1}`, want: &[]float64{0.1}},
			{typ: "float8arr", in: `{1.9}`, want: &[]float64{1.9}},
			{typ: "float4arr", in: `{3.4,5.6,3.14159}`, want: &[]float64{3.4000001, 5.5999999, 3.1415901}},
			{typ: "float8arr", in: `{0.0024,1.4,-89.2345,0.0}`, want: &[]float64{0.0024, 1.4, -89.2345, 0.0}},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			s := FloatArr2Int8Slice{Ptr: new([]int8)}
			return s, s.Ptr
		},
		rows: []test_scanner_row{
			{typ: "float4arr", in: nil, want: new([]int8)},
			{typ: "float8arr", in: `{}`, want: &[]int8{}},
			{typ: "float4arr", in: `{0.1}`, want: &[]int8{0}},
			{typ: "float8arr", in: `{1.9}`, want: &[]int8{1}},
			{typ: "float4arr", in: `{3.4,5.6,3.14159}`, want: &[]int8{3, 5, 3}},
			{typ: "float8arr", in: `{0.0024,1.4,-89.2345,0.0}`, want: &[]int8{0, 1, -89, 0}},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			s := FloatArr2Int16Slice{Ptr: new([]int16)}
			return s, s.Ptr
		},
		rows: []test_scanner_row{
			{typ: "float4arr", in: nil, want: new([]int16)},
			{typ: "float8arr", in: `{}`, want: &[]int16{}},
			{typ: "float4arr", in: `{0.1}`, want: &[]int16{0}},
			{typ: "float8arr", in: `{1.9}`, want: &[]int16{1}},
			{typ: "float4arr", in: `{3.4,5.6,3.14159}`, want: &[]int16{3, 5, 3}},
			{typ: "float8arr", in: `{0.0024,1.4,-89.2345,0.0}`, want: &[]int16{0, 1, -89, 0}},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			s := FloatArr2Int32Slice{Ptr: new([]int32)}
			return s, s.Ptr
		},
		rows: []test_scanner_row{
			{typ: "float4arr", in: nil, want: new([]int32)},
			{typ: "float8arr", in: `{}`, want: &[]int32{}},
			{typ: "float4arr", in: `{0.1}`, want: &[]int32{0}},
			{typ: "float8arr", in: `{1.9}`, want: &[]int32{1}},
			{typ: "float4arr", in: `{3.4,5.6,3.14159}`, want: &[]int32{3, 5, 3}},
			{typ: "float8arr", in: `{0.0024,1.4,-89.2345,0.0}`, want: &[]int32{0, 1, -89, 0}},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			s := FloatArr2Int64Slice{Ptr: new([]int64)}
			return s, s.Ptr
		},
		rows: []test_scanner_row{
			{typ: "float4arr", in: nil, want: new([]int64)},
			{typ: "float8arr", in: `{}`, want: &[]int64{}},
			{typ: "float4arr", in: `{0.1}`, want: &[]int64{0}},
			{typ: "float8arr", in: `{1.9}`, want: &[]int64{1}},
			{typ: "float4arr", in: `{3.4,5.6,3.14159}`, want: &[]int64{3, 5, 3}},
			{typ: "float8arr", in: `{0.0024,1.4,-89.2345,0.0}`, want: &[]int64{0, 1, -89, 0}},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			s := FloatArr2IntSlice{Ptr: new([]int)}
			return s, s.Ptr
		},
		rows: []test_scanner_row{
			{typ: "float4arr", in: nil, want: new([]int)},
			{typ: "float8arr", in: `{}`, want: &[]int{}},
			{typ: "float4arr", in: `{0.1}`, want: &[]int{0}},
			{typ: "float8arr", in: `{1.9}`, want: &[]int{1}},
			{typ: "float4arr", in: `{3.4,5.6,3.14159}`, want: &[]int{3, 5, 3}},
			{typ: "float8arr", in: `{0.0024,1.4,-89.2345,0.0}`, want: &[]int{0, 1, -89, 0}},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			s := FloatArr2Uint8Slice{Ptr: new([]uint8)}
			return s, s.Ptr
		},
		rows: []test_scanner_row{
			{typ: "float4arr", in: nil, want: new([]uint8)},
			{typ: "float8arr", in: `{}`, want: &[]uint8{}},
			{typ: "float4arr", in: `{0.1}`, want: &[]uint8{0}},
			{typ: "float8arr", in: `{1.9}`, want: &[]uint8{1}},
			{typ: "float4arr", in: `{3.4,5.6,3.14159}`, want: &[]uint8{3, 5, 3}},
			{typ: "float8arr", in: `{0.0024,1.4,-89.2345,0.0}`, want: &[]uint8{0, 1, 167, 0}},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			s := FloatArr2Uint16Slice{Ptr: new([]uint16)}
			return s, s.Ptr
		},
		rows: []test_scanner_row{
			{typ: "float4arr", in: nil, want: new([]uint16)},
			{typ: "float8arr", in: `{}`, want: &[]uint16{}},
			{typ: "float4arr", in: `{0.1}`, want: &[]uint16{0}},
			{typ: "float8arr", in: `{1.9}`, want: &[]uint16{1}},
			{typ: "float4arr", in: `{3.4,5.6,3.14159}`, want: &[]uint16{3, 5, 3}},
			{typ: "float8arr", in: `{0.0024,1.4,-89.2345,0.0}`, want: &[]uint16{0, 1, 65447, 0}},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			s := FloatArr2Uint32Slice{Ptr: new([]uint32)}
			return s, s.Ptr
		},
		rows: []test_scanner_row{
			{typ: "float4arr", in: nil, want: new([]uint32)},
			{typ: "float8arr", in: `{}`, want: &[]uint32{}},
			{typ: "float4arr", in: `{0.1}`, want: &[]uint32{0}},
			{typ: "float8arr", in: `{1.9}`, want: &[]uint32{1}},
			{typ: "float4arr", in: `{3.4,5.6,3.14159}`, want: &[]uint32{3, 5, 3}},
			{typ: "float8arr", in: `{0.0024,1.4,-89.2345,0.0}`, want: &[]uint32{0, 1, 4294967207, 0}},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			s := FloatArr2Uint64Slice{Ptr: new([]uint64)}
			return s, s.Ptr
		},
		rows: []test_scanner_row{
			{typ: "float4arr", in: nil, want: new([]uint64)},
			{typ: "float8arr", in: `{}`, want: &[]uint64{}},
			{typ: "float4arr", in: `{0.1}`, want: &[]uint64{0}},
			{typ: "float8arr", in: `{1.9}`, want: &[]uint64{1}},
			{typ: "float4arr", in: `{3.4,5.6,3.14159}`, want: &[]uint64{3, 5, 3}},
			{typ: "float8arr", in: `{0.0024,1.4,-89.2345,0.0}`, want: &[]uint64{0, 1, 18446744073709551527, 0}},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			s := FloatArr2UintSlice{Ptr: new([]uint)}
			return s, s.Ptr
		},
		rows: []test_scanner_row{
			{typ: "float4arr", in: nil, want: new([]uint)},
			{typ: "float8arr", in: `{}`, want: &[]uint{}},
			{typ: "float4arr", in: `{0.1}`, want: &[]uint{0}},
			{typ: "float8arr", in: `{1.9}`, want: &[]uint{1}},
			{typ: "float4arr", in: `{3.4,5.6,3.14159}`, want: &[]uint{3, 5, 3}},
			{typ: "float8arr", in: `{0.0024,1.4,-89.2345,0.0}`, want: &[]uint{0, 1, 18446744073709551527, 0}},
		},
	}}.execute(t)
}
