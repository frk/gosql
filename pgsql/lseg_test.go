package pgsql

import (
	"testing"
)

func TestLineSegment(t *testing.T) {
	testlist2{{
		valuer:  LsegFromFloat64Array2Array2,
		scanner: LsegToFloat64Array2Array2,
		data: []testdata{
			{
				input:  [2][2]float64{{0, 0}, {1, 1}},
				output: [2][2]float64{{0, 0}, {1, 1}}},
			{
				input:  [2][2]float64{{0.5, 0.5}, {1.5, 1.5}},
				output: [2][2]float64{{0.5, 0.5}, {1.5, 1.5}}},
		},
	}, {
		data: []testdata{
			{
				input:  string(`[(0,0),(1,1)]`),
				output: string(`[(0,0),(1,1)]`)},
			{
				input:  string(`[(0.5,0.5),(1.5,1.5)]`),
				output: string(`[(0.5,0.5),(1.5,1.5)]`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte(`[(0,0),(1,1)]`),
				output: []byte(`[(0,0),(1,1)]`)},
			{
				input:  []byte(`[(0.5,0.5),(1.5,1.5)]`),
				output: []byte(`[(0.5,0.5),(1.5,1.5)]`)},
		},
	}}.execute(t, "lseg")
}
