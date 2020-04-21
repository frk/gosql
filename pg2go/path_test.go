package pg2go

import (
	"testing"
)

func TestPath(t *testing.T) {
	testlist2{{
		valuer:  PathFromFloat64Array2Slice,
		scanner: PathToFloat64Array2Slice,
		data: []testdata{
			{input: [][2]float64(nil), output: [][2]float64(nil)},
			{
				input:  [][2]float64{{0, 0}},
				output: [][2]float64{{0, 0}}},
			{
				input:  [][2]float64{{0, 0}, {1, 1}, {2, 0}},
				output: [][2]float64{{0, 0}, {1, 1}, {2, 0}}},
			{
				input:  [][2]float64{{0, 0.5}, {1.5, 1}, {2.4, 0.4}},
				output: [][2]float64{{0, 0.5}, {1.5, 1}, {2.4, 0.4}}},
		},
	}, {
		data: []testdata{
			{input: string(`((0,0))`), output: string(`((0,0))`)},
			{
				input:  string(`((0,0),(1,1))`),
				output: string(`((0,0),(1,1))`)},
			{
				input:  string(`[(0,0.5),(1.5,1),(2.23,0.623)]`),
				output: string(`[(0,0.5),(1.5,1),(2.23,0.623)]`)},
		},
	}, {
		data: []testdata{
			{input: []byte(`((0,0))`), output: []byte(`((0,0))`)},
			{
				input:  []byte(`((0,0),(1,1))`),
				output: []byte(`((0,0),(1,1))`)},
			{
				input:  []byte(`[(0,0.5),(1.5,1),(2.23,0.623)]`),
				output: []byte(`[(0,0.5),(1.5,1),(2.23,0.623)]`)},
		},
	}}.execute(t, "path")
}
