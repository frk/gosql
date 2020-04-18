package convert

import (
	"testing"
)

func TestPointArray(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(PointArrayFromFloat64Array2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := PointArrayToFloat64Array2Slice{Val: new([][2]float64)}
			return v, v.Val
		},
		data: []testdata{
			{input: [][2]float64(nil), output: [][2]float64(nil)},
			{input: [][2]float64{}, output: [][2]float64{}},
			{
				input:  [][2]float64{{1, 2}},
				output: [][2]float64{{1, 2}}},
			{
				input:  [][2]float64{{0, 0.5}, {1, 2.2}, {5.5, 5}},
				output: [][2]float64{{0, 0.5}, {1, 2.2}, {5.5, 5}}},
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
			{input: string(`{}`), output: string(`{}`)},
			{input: string(`{"(0,0)"}`), output: string(`{"(0,0)"}`)},
			{
				input:  string(`{"(0,0)","(1,1)"}`),
				output: string(`{"(0,0)","(1,1)"}`)},
			{
				input:  string(`{"(0,0.5)","(1,2.23)","(5.5,5)"}`),
				output: string(`{"(0,0.5)","(1,2.23)","(5.5,5)"}`)},
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
			{input: []byte(`{}`), output: []byte(`{}`)},
			{input: []byte(`{"(0,0)"}`), output: []byte(`{"(0,0)"}`)},
			{
				input:  []byte(`{"(0,0)","(1,1)"}`),
				output: []byte(`{"(0,0)","(1,1)"}`)},
			{
				input:  []byte(`{"(0,0.5)","(1,2.23)","(5.5,5)"}`),
				output: []byte(`{"(0,0.5)","(1,2.23)","(5.5,5)"}`)},
		},
	}}.execute(t, "pointarr")
}
