package convert

import (
	"database/sql"
	"testing"
)

func TestBoxArr2Float64a2a2Slice(t *testing.T) {
	test_table{{
		scnr: func() (sql.Scanner, interface{}) {
			s := BoxArr2Float64a2a2Slice{Ptr: new([][2][2]float64)}
			return s, s.Ptr
		},
		rows: []testrow{
			{typ: "boxarr", in: nil, want: new([][2][2]float64)},
			{typ: "boxarr", in: `{}`, want: &[][2][2]float64{}},
			{typ: "boxarr", in: `{(1,1),(0,0)}`, want: &[][2][2]float64{{{1, 1}, {0, 0}}}},
			{typ: "boxarr", in: `{(1,1),(0,0);(0,0),(1,1)}`, want: &[][2][2]float64{{{1, 1}, {0, 0}}, {{1, 1}, {0, 0}}}},
			{
				typ: "boxarr",
				in:  `{(4.5,0.79),(3.2,5.63);(43.57,2.12),(4.99,0.22);(0.22,0.31),(4,1.07)}`,
				want: &[][2][2]float64{
					{{4.5, 5.63}, {3.2, 0.79}},
					{{43.57, 2.12}, {4.99, 0.22}},
					{{4, 1.07}, {0.22, 0.31}},
				},
			},
		},
	}}.execute(t)
}
