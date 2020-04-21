package pg2go

import (
	"testing"
)

func TestLine(t *testing.T) {
	testlist2{{
		valuer:  LineFromFloat64Array3,
		scanner: LineToFloat64Array3,
		data: []testdata{
			{input: [3]float64{1, 1, 1}, output: [3]float64{1, 1, 1}},
			{input: [3]float64{0, 5, 10}, output: [3]float64{0, 5, 10}},
			{
				input:  [3]float64{0.5, 5.5, 10.5},
				output: [3]float64{0.5, 5.5, 10.5}},
		},
	}, {
		data: []testdata{
			{input: string(`{1,1,1}`), output: string(`{1,1,1}`)},
			{input: string(`{0,5,10}`), output: string(`{0,5,10}`)},
			{
				input:  string(`{0.5,5.5,10.5}`),
				output: string(`{0.5,5.5,10.5}`)},
		},
	}, {
		data: []testdata{
			{input: []byte(`{1,1,1}`), output: []byte(`{1,1,1}`)},
			{input: []byte(`{0,5,10}`), output: []byte(`{0,5,10}`)},
			{
				input:  []byte(`{0.5,5.5,10.5}`),
				output: []byte(`{0.5,5.5,10.5}`)},
		},
	}}.execute(t, "line")
}
