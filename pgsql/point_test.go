package pgsql

import (
	"testing"
)

func TestPoint(t *testing.T) {
	testlist2{{
		valuer:  PointFromFloat64Array2,
		scanner: PointToFloat64Array2,
		data: []testdata{
			{input: [2]float64{}, output: [2]float64{}},
			{input: [2]float64{1, 1}, output: [2]float64{1, 1}},
			{input: [2]float64{0.5, 1.5}, output: [2]float64{0.5, 1.5}},
		},
	}, {
		data: []testdata{
			{input: string(`(0,0)`), output: string(`(0,0)`)},
			{input: string(`(1,1)`), output: string(`(1,1)`)},
			{input: string(`(0.5,1.5)`), output: string(`(0.5,1.5)`)},
		},
	}, {
		data: []testdata{
			{input: []byte(`(0,0)`), output: []byte(`(0,0)`)},
			{input: []byte(`(1,1)`), output: []byte(`(1,1)`)},
			{input: []byte(`(0.5,1.5)`), output: []byte(`(0.5,1.5)`)},
		},
	}}.execute(t, "point")
}
