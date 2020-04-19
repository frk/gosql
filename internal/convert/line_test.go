package convert

import (
	"testing"
)

func TestLine(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(LineFromFloat64Array3)
		},
		scanner: func() (interface{}, interface{}) {
			v := LineToFloat64Array3{Val: new([3]float64)}
			return v, v.Val
		},
		data: []testdata{
			{input: [3]float64{1, 1, 1}, output: [3]float64{1, 1, 1}},
			{input: [3]float64{0, 5, 10}, output: [3]float64{0, 5, 10}},
			{
				input:  [3]float64{0.5, 5.5, 10.5},
				output: [3]float64{0.5, 5.5, 10.5}},
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
			{input: string(`{1,1,1}`), output: string(`{1,1,1}`)},
			{input: string(`{0,5,10}`), output: string(`{0,5,10}`)},
			{
				input:  string(`{0.5,5.5,10.5}`),
				output: string(`{0.5,5.5,10.5}`)},
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
			{input: []byte(`{1,1,1}`), output: []byte(`{1,1,1}`)},
			{input: []byte(`{0,5,10}`), output: []byte(`{0,5,10}`)},
			{
				input:  []byte(`{0.5,5.5,10.5}`),
				output: []byte(`{0.5,5.5,10.5}`)},
		},
	}}.execute(t, "line")
}
