package convert

import (
	"testing"
)

func TestLsegArray(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(LsegArrayFromFloat64Array2Array2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			v := LsegArrayToFloat64Array2Array2Slice{Val: new([][2][2]float64)}
			return v, v.Val
		},
		data: []testdata{
			{input: [][2][2]float64(nil), output: [][2][2]float64(nil)},
			{input: [][2][2]float64{}, output: [][2][2]float64{}},
			{
				input:  [][2][2]float64{{{0, 0}, {1, 1}}},
				output: [][2][2]float64{{{0, 0}, {1, 1}}}},
			{
				input:  [][2][2]float64{{{1, 2}, {3, 4}}, {{1.5, 2.5}, {3.5, 4.5}}},
				output: [][2][2]float64{{{1, 2}, {3, 4}}, {{1.5, 2.5}, {3.5, 4.5}}}},
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
			{
				input:  string(`{"[(0,0),(1,1)]"}`),
				output: string(`{"[(0,0),(1,1)]"}`)},
			{
				input:  string(`{"[(1,2),(3,4)]","[(1.5,2.5),(3.5,4.5)]"}`),
				output: string(`{"[(1,2),(3,4)]","[(1.5,2.5),(3.5,4.5)]"}`)},
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
			{
				input:  []byte(`{"[(0,0),(1,1)]"}`),
				output: []byte(`{"[(0,0),(1,1)]"}`)},
			{
				input:  []byte(`{"[(1,2),(3,4)]","[(1.5,2.5),(3.5,4.5)]"}`),
				output: []byte(`{"[(1,2),(3,4)]","[(1.5,2.5),(3.5,4.5)]"}`)},
		},
	}}.execute(t, "lsegarr")
}
