package convert

import (
	"testing"
)

func TestPolygon(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(PolygonFromFloat64Array2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := PolygonToFloat64Array2Slice{Val: new([][2]float64)}
			return v, v.Val
		},
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
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			v := new(string)
			return v, v
		},
		data: []testdata{
			{input: string(`((0,0))`), output: string(`((0,0))`)},
			{
				input:  string(`((0,0),(1,1))`),
				output: string(`((0,0),(1,1))`)},
			{
				input:  string(`((0,0.5),(1.5,1),(2.23,0.623))`),
				output: string(`((0,0.5),(1.5,1),(2.23,0.623))`)},
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
			{input: []byte(`((0,0))`), output: []byte(`((0,0))`)},
			{
				input:  []byte(`((0,0),(1,1))`),
				output: []byte(`((0,0),(1,1))`)},
			{
				input:  []byte(`((0,0.5),(1.5,1),(2.23,0.623))`),
				output: []byte(`((0,0.5),(1.5,1),(2.23,0.623))`)},
		},
	}}.execute(t, "polygon")
}
