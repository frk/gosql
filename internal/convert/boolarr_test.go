package convert

import (
	"testing"
)

func TestBoolArray(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(BoolArrayFromBoolSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := BoolArrayToBoolSlice{Val: new([]bool)}
			return s, s.Val
		},
		data: []testdata{
			{input: nil, output: new([]bool)},
			{input: []bool{}, output: &[]bool{}},
			{input: []bool{true}, output: &[]bool{true}},
			{input: []bool{false}, output: &[]bool{false}},
			{
				input:  []bool{false, false, false, true, true, false, true, true},
				output: &[]bool{false, false, false, true, true, false, true, true}},
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
			{input: string("{}"), output: strptr(`{}`)},
			{
				input:  string("{true,false}"),
				output: strptr(`{t,f}`)},
			{
				input:  string("{t,f,f,f,t,t}"),
				output: strptr(`{t,f,f,f,t,t}`)},
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
			{input: []byte("{}"), output: bytesptr(`{}`)},
			{
				input:  []byte("{true,false}"),
				output: bytesptr(`{t,f}`)},
			{
				input:  []byte("{t,f,f,f,t,t}"),
				output: bytesptr(`{t,f,f,f,t,t}`)},
		},
	}}.execute(t, "boolarr")
}
