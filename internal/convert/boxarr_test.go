package convert

import (
	"testing"
)

func TestBoxArray_Valuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return new(BoxArrayFromFloat64Array2Array2Slice)
		},
		rows: []test_valuer_row{
			{typ: "boxarr", in: nil, want: nil},
			{typ: "boxarr", in: [][2][2]float64{}, want: strptr(`{}`)},
			{typ: "boxarr", in: [][2][2]float64{{{1, 1}, {0, 0}}}, want: strptr(`{(1,1),(0,0)}`)},
			{
				typ:  "boxarr",
				in:   [][2][2]float64{{{1, 1}, {0, 0}}, {{0, 0}, {1, 1}}},
				want: strptr(`{(1,1),(0,0);(1,1),(0,0)}`),
			}, {
				typ: "boxarr",
				in: [][2][2]float64{
					{{4.5, 0.79}, {3.2, 5.63}},
					{{43.57, 2.12}, {4.99, 0.22}},
					{{0.22, 0.31}, {4, 1.07}},
				},
				// TODO(mkopriva) figure out whether postgres can be made to return
				// a string exactly matching the input Go input ...
				want: strptr(`{(4.5,5.6299999999999999),(3.2000000000000002,0.79000000000000004);` +
					`(43.57,2.1200000000000001),(4.9900000000000002,0.22);` +
					`(4,1.0700000000000001),(0.22,0.31)}`),
			},
		},
	}}.execute(t)
}

func TestBoxArray_Scanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			s := BoxArrayToFloat64Array2Array2Slice{S: new([][2][2]float64)}
			return s, s.S
		},
		rows: []test_scanner_row{
			{typ: "boxarr", in: nil, want: new([][2][2]float64)},
			{typ: "boxarr", in: `{}`, want: &[][2][2]float64{}},
			{typ: "boxarr", in: `{(1,1),(0,0)}`, want: &[][2][2]float64{{{1, 1}, {0, 0}}}},
			{
				typ:  "boxarr",
				in:   `{(1,1),(0,0);(0,0),(1,1)}`,
				want: &[][2][2]float64{{{1, 1}, {0, 0}}, {{1, 1}, {0, 0}}},
			}, {
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
