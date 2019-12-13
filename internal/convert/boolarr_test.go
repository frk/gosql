package convert

import (
	"database/sql"
	"testing"
)

func TestBoolArrScanners(t *testing.T) {
	test_table{{
		scnr: func() (sql.Scanner, interface{}) {
			s := BoolArr2BoolSlice{Ptr: new([]bool)}
			return s, s.Ptr
		},
		rows: []testrow{
			{typ: "boolarr", in: nil, want: new([]bool)},
			{typ: "boolarr", in: `{}`, want: &[]bool{}},
			{typ: "boolarr", in: `{f}`, want: &[]bool{false}},
			{typ: "boolarr", in: `{t}`, want: &[]bool{true}},
			{typ: "boolarr", in: `{t,f}`, want: &[]bool{true, false}},
			{typ: "boolarr", in: `{f,t,t,t,f,t,f,f}`, want: &[]bool{false, true, true, true, false, true, false, false}},
		},
	}}.execute(t)
}
