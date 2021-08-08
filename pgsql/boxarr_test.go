package pgsql

import (
	"testing"
)

func TestBoxArray(t *testing.T) {
	testlist2{{
		valuer:  BoxArrayFromFloat64Array2Array2Slice,
		scanner: BoxArrayToFloat64Array2Array2Slice,
		data: []testdata{
			{
				input:  [][2][2]float64(nil),
				output: [][2][2]float64(nil)},
			{
				input:  [][2][2]float64{},
				output: [][2][2]float64{}},
			{
				input:  [][2][2]float64{{{1, 1}, {0, 0}}},
				output: [][2][2]float64{{{1, 1}, {0, 0}}}},
			{
				input:  [][2][2]float64{{{1, 1}, {0, 0}}, {{0, 0}, {1, 1}}},
				output: [][2][2]float64{{{1, 1}, {0, 0}}, {{1, 1}, {0, 0}}}},
			{
				input: [][2][2]float64{
					{{4.5, 0.79}, {3.2, 5.63}},
					{{43.57, 2.12}, {4.99, 0.22}},
					{{0.22, 0.31}, {4, 1.07}},
				},
				output: [][2][2]float64{
					{{4.5, 5.63}, {3.2, 0.79}},
					{{43.57, 2.12}, {4.99, 0.22}},
					{{4, 1.07}, {0.22, 0.31}},
				},
			},
		},
	}, {
		data: []testdata{
			{
				input:  string("{}"),
				output: string(`{}`)},
			{
				input:  string("{(0,0),(0,0)}"),
				output: string(`{(0,0),(0,0)}`)},
			{
				input:  string("{(1,1),(0,0);(0,0),(1,1)}"),
				output: string(`{(1,1),(0,0);(1,1),(0,0)}`)},
			{
				input:  string(`{(4.5,0.79),(3.2,5.63);(43.57,2.12),(4.99,0.22);(0.22,0.31),(4,1.07)}`),
				output: string(`{(4.5,5.63),(3.2,0.79);(43.57,2.12),(4.99,0.22);(4,1.07),(0.22,0.31)}`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte("{}"),
				output: []byte(`{}`)},
			{
				input:  []byte("{(0,0),(0,0)}"),
				output: []byte(`{(0,0),(0,0)}`)},
			{
				input:  []byte("{(1,1),(0,0);(0,0),(1,1)}"),
				output: []byte(`{(1,1),(0,0);(1,1),(0,0)}`)},
			{
				input:  []byte(`{(4.5,0.79),(3.2,5.63);(43.57,2.12),(4.99,0.22);(0.22,0.31),(4,1.07)}`),
				output: []byte(`{(4.5,5.63),(3.2,0.79);(43.57,2.12),(4.99,0.22);(4,1.07),(0.22,0.31)}`)},
		},
	}}.execute(t, "boxarr")
}
