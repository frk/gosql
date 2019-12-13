package convert

import (
	"database/sql"
	"testing"
)

func TestBitArrScanners(t *testing.T) {
	test_table{{
		scnr: func() (sql.Scanner, interface{}) {
			s := BitArr2BoolSlice{Ptr: new([]bool)}
			return s, s.Ptr
		},
		rows: []testrow{
			{typ: "bitarr", in: nil, want: new([]bool)},
			{typ: "bitarr", in: `{}`, want: &[]bool{}},
			{typ: "bitarr", in: `{0}`, want: &[]bool{false}},
			{typ: "bitarr", in: `{1}`, want: &[]bool{true}},
			{typ: "bitarr", in: `{1,0}`, want: &[]bool{true, false}},
			{typ: "bitarr", in: `{0,1,1,1,0,1,0,0}`, want: &[]bool{false, true, true, true, false, true, false, false}},
		},
	}, {
		scnr: func() (sql.Scanner, interface{}) {
			s := BitArr2Uint8Slice{Ptr: new([]uint8)}
			return s, s.Ptr
		},
		rows: []testrow{
			{typ: "bitarr", in: nil, want: new([]uint8)},
			{typ: "bitarr", in: `{}`, want: &[]uint8{}},
			{typ: "bitarr", in: `{0}`, want: &[]uint8{0}},
			{typ: "bitarr", in: `{1}`, want: &[]uint8{1}},
			{typ: "bitarr", in: `{1,0}`, want: &[]uint8{1, 0}},
			{typ: "bitarr", in: `{0,1,1,1,0,1,0,0}`, want: &[]uint8{0, 1, 1, 1, 0, 1, 0, 0}},
		},
	}, {
		scnr: func() (sql.Scanner, interface{}) {
			s := BitArr2UintSlice{Ptr: new([]uint)}
			return s, s.Ptr
		},
		rows: []testrow{
			{typ: "bitarr", in: nil, want: new([]uint)},
			{typ: "bitarr", in: `{}`, want: &[]uint{}},
			{typ: "bitarr", in: `{0}`, want: &[]uint{0}},
			{typ: "bitarr", in: `{1}`, want: &[]uint{1}},
			{typ: "bitarr", in: `{1,0}`, want: &[]uint{1, 0}},
			{typ: "bitarr", in: `{0,1,1,1,0,1,0,0}`, want: &[]uint{0, 1, 1, 1, 0, 1, 0, 0}},
		},
	}}.execute(t)
}
