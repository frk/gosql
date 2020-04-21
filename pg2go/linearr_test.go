package pg2go

import (
	"testing"
)

func TestLineArray(t *testing.T) {
	testlist2{{
		valuer:  LineArrayFromFloat64Array3Slice,
		scanner: LineArrayToFloat64Array3Slice,
		data: []testdata{
			{input: [][3]float64(nil), output: [][3]float64(nil)},
			{input: [][3]float64{}, output: [][3]float64{}},
			{input: [][3]float64{{1, 2, 3}}, output: [][3]float64{{1, 2, 3}}},
			{
				input:  [][3]float64{{1, 2, 3}, {4, 5, 6}},
				output: [][3]float64{{1, 2, 3}, {4, 5, 6}}},
		},
	}, {
		data: []testdata{
			{input: string(`{}`), output: string(`{}`)},
			{input: string(`{"{1,2,3}"}`), output: string(`{"{1,2,3}"}`)},
			{
				input:  string(`{"{1,2,3}","{4,5,6}"}`),
				output: string(`{"{1,2,3}","{4,5,6}"}`)},
			{
				input:  string(`{NULL,"{4,5,6}",NULL}`),
				output: string(`{NULL,"{4,5,6}",NULL}`)},
		},
	}, {
		data: []testdata{
			{input: []byte(`{}`), output: []byte(`{}`)},
			{input: []byte(`{"{1,2,3}"}`), output: []byte(`{"{1,2,3}"}`)},
			{
				input:  []byte(`{"{1,2,3}","{4,5,6}"}`),
				output: []byte(`{"{1,2,3}","{4,5,6}"}`)},
			{
				input:  []byte(`{NULL,"{4,5,6}",NULL}`),
				output: []byte(`{NULL,"{4,5,6}",NULL}`)},
		},
	}}.execute(t, "linearr")
}
