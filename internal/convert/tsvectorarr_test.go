package convert

import (
	"testing"
)

func TestTSVectorArray(t *testing.T) {
	B := func(v string) []byte { return []byte(v) }

	testlist{{
		valuer: func() interface{} {
			return new(TSVectorArrayFromStringSliceSlice)
		},
		scanner: func() (interface{}, interface{}) {
			v := TSVectorArrayToStringSliceSlice{Val: new([][]string)}
			return v, v.Val
		},
		data: []testdata{
			{input: [][]string(nil), output: [][]string(nil)},
			{input: [][]string{}, output: [][]string{}},
			{
				input:  [][]string{{"'cat'", "'fat'"}, {"'bat'", "'rat'"}},
				output: [][]string{{"'cat'", "'fat'"}, {"'bat'", "'rat'"}}},
			{
				input:  [][]string{{"'cat'", "'fat'"}, {}, nil},
				output: [][]string{{"'cat'", "'fat'"}, {}, nil}},
		},
	}, {
		valuer: func() interface{} {
			return new(TSVectorArrayFromByteSliceSliceSlice)
		},
		scanner: func() (interface{}, interface{}) {
			v := TSVectorArrayToByteSliceSliceSlice{Val: new([][][]byte)}
			return v, v.Val
		},
		data: []testdata{
			{input: [][][]byte(nil), output: [][][]byte(nil)},
			{input: [][][]byte{}, output: [][][]byte{}},
			{
				input:  [][][]byte{{B("'cat'"), B("'fat'")}, {B("'bat'"), B("'rat'")}},
				output: [][][]byte{{B("'cat'"), B("'fat'")}, {B("'bat'"), B("'rat'")}}},
			{
				input:  [][][]byte{{B("'cat'"), B("'fat'")}, {}, nil},
				output: [][][]byte{{B("'cat'"), B("'fat'")}, {}, nil}},
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
			{
				input:  string("{}"),
				output: string("{}")},
			{
				input:  string(`{"'cat' 'fat'","'rat'"}`),
				output: string(`{"'cat' 'fat'",'rat'}`)},
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
			{
				input:  []byte("{}"),
				output: []byte("{}")},
			{
				input:  []byte(`{"'cat' 'fat'","'rat'"}`),
				output: []byte(`{"'cat' 'fat'",'rat'}`)},
		},
	}}.execute(t, "tsvectorarr")
}
