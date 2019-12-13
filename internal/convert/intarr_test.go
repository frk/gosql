package convert

import (
	"database/sql"
	"testing"
)

func TestIntArrScanners(t *testing.T) {
	test_table{{
		scnr: func() (sql.Scanner, interface{}) {
			s := IntArr2IntSlice{Ptr: new([]int)}
			return s, s.Ptr
		},
		rows: []testrow{
			{typ: "int2arr", in: nil, want: new([]int)},
			{typ: "int2arr", in: `{}`, want: &[]int{}},
			{typ: "int2arr", in: `{42,32767,-10,0}`, want: &[]int{42, 32767, -10, 0}},
			{typ: "int4arr", in: nil, want: new([]int)},
			{typ: "int4arr", in: `{}`, want: &[]int{}},
			{typ: "int4arr", in: `{42,2147483647,-10,0}`, want: &[]int{42, 2147483647, -10, 0}},
			{typ: "int8arr", in: nil, want: new([]int)},
			{typ: "int8arr", in: `{}`, want: &[]int{}},
			{typ: "int8arr", in: `{42,9223372036854775807,-10,0}`, want: &[]int{42, 9223372036854775807, -10, 0}},
		},
	}, {
		scnr: func() (sql.Scanner, interface{}) {
			s := IntArr2Int8Slice{Ptr: new([]int8)}
			return s, s.Ptr
		},
		rows: []testrow{
			{typ: "int2arr", in: nil, want: new([]int8)},
			{typ: "int2arr", in: `{}`, want: &[]int8{}},
			{typ: "int2arr", in: `{42,127,-10,0}`, want: &[]int8{42, 127, -10, 0}},
			{typ: "int4arr", in: nil, want: new([]int8)},
			{typ: "int4arr", in: `{}`, want: &[]int8{}},
			{typ: "int4arr", in: `{42,127,-10,0}`, want: &[]int8{42, 127, -10, 0}},
			{typ: "int8arr", in: nil, want: new([]int8)},
			{typ: "int8arr", in: `{}`, want: &[]int8{}},
			{typ: "int8arr", in: `{42,128,-10,0}`, want: &[]int8{42, -128, -10, 0}}, // note the integer overflow
		},
	}, {
		scnr: func() (sql.Scanner, interface{}) {
			s := IntArr2Int16Slice{Ptr: new([]int16)}
			return s, s.Ptr
		},
		rows: []testrow{
			{typ: "int2arr", in: nil, want: new([]int16)},
			{typ: "int2arr", in: `{}`, want: &[]int16{}},
			{typ: "int2arr", in: `{42,32767,-10,0}`, want: &[]int16{42, 32767, -10, 0}},
			{typ: "int4arr", in: nil, want: new([]int16)},
			{typ: "int4arr", in: `{}`, want: &[]int16{}},
			{typ: "int4arr", in: `{42,32768,-10,0}`, want: &[]int16{42, -32768, -10, 0}},
			{typ: "int8arr", in: nil, want: new([]int16)},
			{typ: "int8arr", in: `{}`, want: &[]int16{}},
			{typ: "int8arr", in: `{42,32767,-10,0}`, want: &[]int16{42, 32767, -10, 0}},
		},
	}, {
		scnr: func() (sql.Scanner, interface{}) {
			s := IntArr2Int32Slice{Ptr: new([]int32)}
			return s, s.Ptr
		},
		rows: []testrow{
			{typ: "int2arr", in: nil, want: new([]int32)},
			{typ: "int2arr", in: `{}`, want: &[]int32{}},
			{typ: "int2arr", in: `{42,127,-10,0}`, want: &[]int32{42, 127, -10, 0}},
			{typ: "int4arr", in: nil, want: new([]int32)},
			{typ: "int4arr", in: `{}`, want: &[]int32{}},
			{typ: "int4arr", in: `{42,2147483647,-10,0}`, want: &[]int32{42, 2147483647, -10, 0}},
			{typ: "int8arr", in: nil, want: new([]int32)},
			{typ: "int8arr", in: `{}`, want: &[]int32{}},
			{typ: "int8arr", in: `{42,2147483648,-10,0}`, want: &[]int32{42, -2147483648, -10, 0}},
		},
	}, {
		scnr: func() (sql.Scanner, interface{}) {
			s := IntArr2Int64Slice{Ptr: new([]int64)}
			return s, s.Ptr
		},
		rows: []testrow{
			{typ: "int2arr", in: nil, want: new([]int64)},
			{typ: "int2arr", in: `{}`, want: &[]int64{}},
			{typ: "int2arr", in: `{42,127,-10,0}`, want: &[]int64{42, 127, -10, 0}},
			{typ: "int4arr", in: nil, want: new([]int64)},
			{typ: "int4arr", in: `{}`, want: &[]int64{}},
			{typ: "int4arr", in: `{42,127,-10,0}`, want: &[]int64{42, 127, -10, 0}},
			{typ: "int8arr", in: nil, want: new([]int64)},
			{typ: "int8arr", in: `{}`, want: &[]int64{}},
			{typ: "int8arr", in: `{42,9223372036854775807,-10,0}`, want: &[]int64{42, 9223372036854775807, -10, 0}},
		},
	}}.execute(t)
}
