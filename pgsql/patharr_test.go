package pgsql

import (
	"testing"
)

func TestPathArray(t *testing.T) {
	testlist2{{
		valuer:  PathArrayFromFloat64Array2SliceSlice,
		scanner: PathArrayToFloat64Array2SliceSlice,
		data: []testdata{
			{input: [][][2]float64(nil), output: [][][2]float64(nil)},
			{input: [][][2]float64{}, output: [][][2]float64{}},
			{
				input:  [][][2]float64{{{0, 0}}},
				output: [][][2]float64{{{0, 0}}}},
			{
				input: [][][2]float64{
					{{1, 1}, {2, 2}, {3, 3}},
					{{1.5, 1.5}, {2.5, 2.5}, {3.5, 3.5}}},
				output: [][][2]float64{
					{{1, 1}, {2, 2}, {3, 3}},
					{{1.5, 1.5}, {2.5, 2.5}, {3.5, 3.5}}}},
			{
				input: [][][2]float64{
					{{1, 1}, {3, 3}},
					[][2]float64(nil),
					{{2.5, 2.5}, {3.5, 3.5}},
					[][2]float64(nil)},
				output: [][][2]float64{
					{{1, 1}, {3, 3}},
					[][2]float64(nil),
					{{2.5, 2.5}, {3.5, 3.5}},
					[][2]float64(nil)}},
		},
	}, {
		data: []testdata{
			{input: string(`{}`), output: string(`{}`)},
			{
				input:  string(`{"((0,0))"}`),
				output: string(`{"((0,0))"}`)},
			{
				input:  string(`{"((1,1),(2,2),(3,3))","[(1.5,1.5),(2.5,2.5),(3.5,3.5)]"}`),
				output: string(`{"((1,1),(2,2),(3,3))","[(1.5,1.5),(2.5,2.5),(3.5,3.5)]"}`)},
			{
				input:  string(`{"((1,1),(3,3))",NULL,"[(2.5,2.5),(3.5,3.5)]"}`),
				output: string(`{"((1,1),(3,3))",NULL,"[(2.5,2.5),(3.5,3.5)]"}`)},
		},
	}, {
		data: []testdata{
			{input: []byte(`{}`), output: []byte(`{}`)},
			{
				input:  []byte(`{"[(0,0)]"}`),
				output: []byte(`{"[(0,0)]"}`)},
			{
				input:  []byte(`{"[(1,1),(2,2),(3,3)]","((1.5,1.5),(2.5,2.5),(3.5,3.5))"}`),
				output: []byte(`{"[(1,1),(2,2),(3,3)]","((1.5,1.5),(2.5,2.5),(3.5,3.5))"}`)},
			{
				input:  []byte(`{"[(1,1),(3,3)]",NULL,"((2.5,2.5),(3.5,3.5))"}`),
				output: []byte(`{"[(1,1),(3,3)]",NULL,"((2.5,2.5),(3.5,3.5))"}`)},
		},
	}}.execute(t, "patharr")
}
