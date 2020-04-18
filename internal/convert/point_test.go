package convert

import (
	"testing"
)

func TestPoint(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(PointFromFloat64Array2)
		},
		scanner: func() (interface{}, interface{}) {
			v := PointToFloat64Array2{Val: new([2]float64)}
			return v, v.Val
		},
		data: []testdata{
			{input: [2]float64{}, output: [2]float64{}},
			{input: [2]float64{1, 1}, output: [2]float64{1, 1}},
			{input: [2]float64{0.5, 1.5}, output: [2]float64{0.5, 1.5}},
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
			{input: string(`(0,0)`), output: string(`(0,0)`)},
			{input: string(`(1,1)`), output: string(`(1,1)`)},
			{input: string(`(0.5,1.5)`), output: string(`(0.5,1.5)`)},
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
			{input: []byte(`(0,0)`), output: []byte(`(0,0)`)},
			{input: []byte(`(1,1)`), output: []byte(`(1,1)`)},
			{input: []byte(`(0.5,1.5)`), output: []byte(`(0.5,1.5)`)},
		},
	}}.execute(t, "point")
}
