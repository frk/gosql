package convert

import (
	"testing"
)

func TestVarBitArray(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return nil // TODO
		},
		scanner: func() (interface{}, interface{}) {
			s := VarBitArr2StringSlice{Val: new([]string)}
			return s, s.Val
		},
		data: []testdata{
			{input: nil, output: new([]string)},
			{input: string(`{}`), output: &[]string{}},
			{
				input:  string(`{101010,01,10,1111100}`),
				output: &[]string{"101010", "01", "10", "1111100"}},
		},
	}, {
		valuer: func() interface{} {
			return nil // TODO
		},
		scanner: func() (interface{}, interface{}) {
			s := VarBitArr2Int64Slice{Val: new([]int64)}
			return s, s.Val
		},
		data: []testdata{
			{input: nil, output: new([]int64)},
			{input: string(`{}`), output: &[]int64{}},
			{
				input:  string(`{101010,01,10,1111100}`),
				output: &[]int64{42, 1, 2, 124}},
		},
	}, {
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			s := new(string)
			return s, s
		},
		data: []testdata{
			{input: string(`{}`), output: strptr(`{}`)},
			{
				input:  string(`{101010,01,10,1111100}`),
				output: strptr(`{101010,01,10,1111100}`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			s := new([]byte)
			return s, s
		},
		data: []testdata{
			{input: []byte(`{}`), output: bytesptr(`{}`)},
			{
				input:  []byte(`{101010,01,10,1111100}`),
				output: bytesptr(`{101010,01,10,1111100}`)},
		},
	}}.execute(t, "varbitarr")
}
