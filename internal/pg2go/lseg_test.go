package pg2go

import (
	"testing"
)

func TestLineSegment(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(LsegFromFloat64Array2Array2)
		},
		scanner: func() (interface{}, interface{}) {
			v := LsegToFloat64Array2Array2{Val: new([2][2]float64)}
			return v, v.Val
		},
		data: []testdata{
			{
				input:  [2][2]float64{{0, 0}, {1, 1}},
				output: [2][2]float64{{0, 0}, {1, 1}}},
			{
				input:  [2][2]float64{{0.5, 0.5}, {1.5, 1.5}},
				output: [2][2]float64{{0.5, 0.5}, {1.5, 1.5}}},
		},
	}, {
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			v := new(string)
			return v, v
		},
		data: []testdata{
			{
				input:  string(`[(0,0),(1,1)]`),
				output: string(`[(0,0),(1,1)]`)},
			{
				input:  string(`[(0.5,0.5),(1.5,1.5)]`),
				output: string(`[(0.5,0.5),(1.5,1.5)]`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			v := new([]byte)
			return v, v
		},
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
