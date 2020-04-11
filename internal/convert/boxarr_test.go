package convert

import (
	"testing"
)

func TestBoxArray(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(BoxArrayFromFloat64Array2Array2Slice)
		},
		scanner: func() (interface{}, interface{}) {
			s := BoxArrayToFloat64Array2Array2Slice{Val: new([][2][2]float64)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  nil,
				output: [][2][2]float64(nil)},
			{
				input:  [][2][2]float64{},
				output: [][2][2]float64{}},
			{
				input:  [][2][2]float64{{{1, 1}, {0, 0}}},
				output: [][2][2]float64{{{1, 1}, {0, 0}}}},
			{
				input:  [][2][2]float64{{{1, 1}, {0, 0}}, {{0, 0}, {1, 1}}},
				output: [][2][2]float64{{{1, 1}, {0, 0}}, {{1, 1}, {0, 0}}}},
			{
				input: [][2][2]float64{
					{{4.5, 0.79}, {3.2, 5.63}},
					{{43.57, 2.12}, {4.99, 0.22}},
					{{0.22, 0.31}, {4, 1.07}},
				},
				output: [][2][2]float64{
					{{4.5, 5.63}, {3.2, 0.79}},
					{{43.57, 2.12}, {4.99, 0.22}},
					{{4, 1.07}, {0.22, 0.31}},
				},
			},
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
			{
				input:  string("{}"),
				output: string(`{}`)},
			{
				input:  string("{(0,0),(0,0)}"),
				output: string(`{(0,0),(0,0)}`)},
			{
				input:  string("{(1,1),(0,0);(0,0),(1,1)}"),
				output: string(`{(1,1),(0,0);(1,1),(0,0)}`)},
			{
				input: string(`{(4.5,0.79),(3.2,5.63);(43.57,2.12),(4.99,0.22);(0.22,0.31),(4,1.07)}`),
				output: string(`{(4.5,5.6299999999999999),(3.2000000000000002,0.79000000000000004);` +
					`(43.57,2.1200000000000001),(4.9900000000000002,0.22);` +
					`(4,1.0700000000000001),(0.22,0.31)}`)},
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
			{
				input:  nil,
				output: []byte(nil)},
			{
				input:  []byte("{}"),
				output: []byte(`{}`)},
			{
				input:  []byte("{(0,0),(0,0)}"),
				output: []byte(`{(0,0),(0,0)}`)},
			{
				input:  []byte("{(1,1),(0,0);(0,0),(1,1)}"),
				output: []byte(`{(1,1),(0,0);(1,1),(0,0)}`)},
			{
				input: []byte(`{(4.5,0.79),(3.2,5.63);(43.57,2.12),(4.99,0.22);(0.22,0.31),(4,1.07)}`),
				output: []byte(`{(4.5,5.6299999999999999),(3.2000000000000002,0.79000000000000004);` +
					`(43.57,2.1200000000000001),(4.9900000000000002,0.22);` +
					`(4,1.0700000000000001),(0.22,0.31)}`)},
		},
	}}.execute(t, "boxarr")
}
