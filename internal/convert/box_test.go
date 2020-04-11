package convert

import (
	"testing"
)

func TestBox(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(BoxFromFloat64Array2Array2)
		},
		scanner: func() (interface{}, interface{}) {
			s := BoxToFloat64Array2Array2{Val: new([2][2]float64)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  nil,
				output: new([2][2]float64)},
			{
				input:  [2][2]float64{{1, 1}, {0, 0}},
				output: &[2][2]float64{{1, 1}, {0, 0}}},
			{
				input:  [2][2]float64{{0, 0}, {1, 1}},
				output: &[2][2]float64{{1, 1}, {0, 0}}},
			{
				input:  [2][2]float64{{4.5203, 0.79322}, {3.2, 5.63333}},
				output: &[2][2]float64{{4.5203, 5.63333}, {3.2, 0.79322}}},
		},
	}, {
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			d := new(string)
			return d, d
		},
		data: []testdata{
			{input: string("(0,0),(0,0)"), output: strptr(`(0,0),(0,0)`)},
			{input: string("(1,1),(0,0)"), output: strptr(`(1,1),(0,0)`)},
			{input: string("(0,0),(1,1)"), output: strptr(`(1,1),(0,0)`)},
			{
				input: string("(4.5203,0.79322),(3.2,5.63333)"),
				output: strptr(`(4.5202999999999998,5.6333299999999999),` +
					`(3.2000000000000002,0.79322000000000004)`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			d := new([]byte)
			return d, d
		},
		data: []testdata{
			{input: nil, output: new([]byte)},
			{input: []byte("(0,0),(0,0)"), output: bytesptr(`(0,0),(0,0)`)},
			{input: []byte("(1,1),(0,0)"), output: bytesptr(`(1,1),(0,0)`)},
			{input: []byte("(0,0),(1,1)"), output: bytesptr(`(1,1),(0,0)`)},
			{
				input: []byte("(4.5203,0.79322),(3.2,5.63333)"),
				output: bytesptr(`(4.5202999999999998,5.6333299999999999),` +
					`(3.2000000000000002,0.79322000000000004)`)},
		},
	}}.execute(t, "box")
}
