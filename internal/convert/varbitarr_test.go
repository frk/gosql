package convert

import (
	"testing"
)

func TestVarBitArrScanners(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			s := VarBitArr2StringSlice{Ptr: new([]string)}
			return s, s.Ptr
		},
		rows: []test_scanner_row{
			{typ: "varbitarr", in: nil, want: new([]string)},
			{typ: "varbitarr", in: `{}`, want: &[]string{}},
			{typ: "varbitarr", in: `{101010,01,10,1111100}`, want: &[]string{"101010", "01", "10", "1111100"}},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			s := VarBitArr2Int64Slice{Ptr: new([]int64)}
			return s, s.Ptr
		},
		rows: []test_scanner_row{
			{typ: "varbitarr", in: nil, want: new([]int64)},
			{typ: "varbitarr", in: `{}`, want: &[]int64{}},
			{typ: "varbitarr", in: `{101010,01,10,1111100}`, want: &[]int64{42, 1, 2, 124}},
		},
	}}.execute(t)
}
