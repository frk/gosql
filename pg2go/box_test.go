package pg2go

import (
	"testing"
)

func TestBox(t *testing.T) {
	testlist2{{
		valuer:  BoxFromFloat64Array2Array2,
		scanner: BoxToFloat64Array2Array2,
		data: []testdata{
			{
				input:  [2][2]float64{},
				output: [2][2]float64{}},
			{
				input:  [2][2]float64{{1, 1}, {0, 0}},
				output: [2][2]float64{{1, 1}, {0, 0}}},
			{
				input:  [2][2]float64{{0, 0}, {1, 1}},
				output: [2][2]float64{{1, 1}, {0, 0}}},
			{
				input:  [2][2]float64{{4.5203, 0.79322}, {3.2, 5.63333}},
				output: [2][2]float64{{4.5203, 5.63333}, {3.2, 0.79322}}},
		},
	}, {
		data: []testdata{
			{input: string("(0,0),(0,0)"), output: string(`(0,0),(0,0)`)},
			{input: string("(1,1),(0,0)"), output: string(`(1,1),(0,0)`)},
			{input: string("(0,0),(1,1)"), output: string(`(1,1),(0,0)`)},
			{
				input: string("(4.5203,0.79322),(3.2,5.63333)"),
				output: string(`(4.5202999999999998,5.6333299999999999),` +
					`(3.2000000000000002,0.79322000000000004)`)},
		},
	}, {
		data: []testdata{
			{input: []byte("(0,0),(0,0)"), output: []byte(`(0,0),(0,0)`)},
			{input: []byte("(1,1),(0,0)"), output: []byte(`(1,1),(0,0)`)},
			{input: []byte("(0,0),(1,1)"), output: []byte(`(1,1),(0,0)`)},
			{
				input: []byte("(4.5203,0.79322),(3.2,5.63333)"),
				output: []byte(`(4.5202999999999998,5.6333299999999999),` +
					`(3.2000000000000002,0.79322000000000004)`)},
		},
	}}.execute(t, "box")
}
