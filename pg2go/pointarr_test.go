package pg2go

import (
	"testing"
)

func TestPointArray(t *testing.T) {
	testlist2{{
		valuer:  PointArrayFromFloat64Array2Slice,
		scanner: PointArrayToFloat64Array2Slice,
		data: []testdata{
			{input: [][2]float64(nil), output: [][2]float64(nil)},
			{input: [][2]float64{}, output: [][2]float64{}},
			{
				input:  [][2]float64{{1, 2}},
				output: [][2]float64{{1, 2}}},
			{
				input:  [][2]float64{{0, 0.5}, {1, 2.2}, {5.5, 5}},
				output: [][2]float64{{0, 0.5}, {1, 2.2}, {5.5, 5}}},
		},
	}, {
		data: []testdata{
			{input: string(`{}`), output: string(`{}`)},
			{input: string(`{"(0,0)"}`), output: string(`{"(0,0)"}`)},
			{
				input:  string(`{"(0,0)","(1,1)"}`),
				output: string(`{"(0,0)","(1,1)"}`)},
			{
				input:  string(`{"(0,0.5)","(1,2.23)","(5.5,5)"}`),
				output: string(`{"(0,0.5)","(1,2.23)","(5.5,5)"}`)},
		},
	}, {
		data: []testdata{
			{input: []byte(`{}`), output: []byte(`{}`)},
			{input: []byte(`{"(0,0)"}`), output: []byte(`{"(0,0)"}`)},
			{
				input:  []byte(`{"(0,0)","(1,1)"}`),
				output: []byte(`{"(0,0)","(1,1)"}`)},
			{
				input:  []byte(`{"(0,0.5)","(1,2.23)","(5.5,5)"}`),
				output: []byte(`{"(0,0.5)","(1,2.23)","(5.5,5)"}`)},
		},
	}}.execute(t, "pointarr")
}
