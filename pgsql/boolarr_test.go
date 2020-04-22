package pgsql

import (
	"testing"
)

func TestBoolArray(t *testing.T) {
	testlist2{{
		valuer:  BoolArrayFromBoolSlice,
		scanner: BoolArrayToBoolSlice,
		data: []testdata{
			{input: []bool(nil), output: []bool(nil)},
			{input: []bool{}, output: []bool{}},
			{input: []bool{true}, output: []bool{true}},
			{input: []bool{false}, output: []bool{false}},
			{
				input:  []bool{false, false, false, true, true, false, true, true},
				output: []bool{false, false, false, true, true, false, true, true}},
		},
	}, {
		data: []testdata{
			{input: string("{}"), output: string(`{}`)},
			{
				input:  string("{true,false}"),
				output: string(`{t,f}`)},
			{
				input:  string("{t,f,f,f,t,t}"),
				output: string(`{t,f,f,f,t,t}`)},
		},
	}, {
		data: []testdata{
			{input: []byte("{}"), output: []byte(`{}`)},
			{
				input:  []byte("{true,false}"),
				output: []byte(`{t,f}`)},
			{
				input:  []byte("{t,f,f,f,t,t}"),
				output: []byte(`{t,f,f,f,t,t}`)},
		},
	}}.execute(t, "boolarr")
}
