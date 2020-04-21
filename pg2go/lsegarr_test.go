package pg2go

import (
	"testing"
)

func TestLsegArray(t *testing.T) {
	testlist2{{
		valuer:  LsegArrayFromFloat64Array2Array2Slice,
		scanner: LsegArrayToFloat64Array2Array2Slice,
		data: []testdata{
			{input: [][2][2]float64(nil), output: [][2][2]float64(nil)},
			{input: [][2][2]float64{}, output: [][2][2]float64{}},
			{
				input:  [][2][2]float64{{{0, 0}, {1, 1}}},
				output: [][2][2]float64{{{0, 0}, {1, 1}}}},
			{
				input:  [][2][2]float64{{{1, 2}, {3, 4}}, {{1.5, 2.5}, {3.5, 4.5}}},
				output: [][2][2]float64{{{1, 2}, {3, 4}}, {{1.5, 2.5}, {3.5, 4.5}}}},
		},
	}, {
		data: []testdata{
			{input: string(`{}`), output: string(`{}`)},
			{
				input:  string(`{"[(0,0),(1,1)]"}`),
				output: string(`{"[(0,0),(1,1)]"}`)},
			{
				input:  string(`{"[(1,2),(3,4)]","[(1.5,2.5),(3.5,4.5)]"}`),
				output: string(`{"[(1,2),(3,4)]","[(1.5,2.5),(3.5,4.5)]"}`)},
		},
	}, {
		data: []testdata{
			{input: []byte(`{}`), output: []byte(`{}`)},
			{
				input:  []byte(`{"[(0,0),(1,1)]"}`),
				output: []byte(`{"[(0,0),(1,1)]"}`)},
			{
				input:  []byte(`{"[(1,2),(3,4)]","[(1.5,2.5),(3.5,4.5)]"}`),
				output: []byte(`{"[(1,2),(3,4)]","[(1.5,2.5),(3.5,4.5)]"}`)},
		},
	}}.execute(t, "lsegarr")
}
