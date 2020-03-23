package convert

import (
	"testing"
)

func TestBoolArrScanners(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			s := BoolArr2BoolSlice{Ptr: new([]bool)}
			return s, s.Ptr
		},
		rows: []test_scanner_row{
			{typ: "boolarr", in: nil, want: new([]bool)},
			{typ: "boolarr", in: `{}`, want: &[]bool{}},
			{typ: "boolarr", in: `{f}`, want: &[]bool{false}},
			{typ: "boolarr", in: `{t}`, want: &[]bool{true}},
			{typ: "boolarr", in: `{t,f}`, want: &[]bool{true, false}},
			{typ: "boolarr", in: `{f,t,t,t,f,t,f,f}`, want: &[]bool{false, true, true, true, false, true, false, false}},
		},
	}}.execute(t)
}
