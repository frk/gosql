package pg2go

import (
	"testing"
)

func TestTSVectorArray(t *testing.T) {
	B := func(v string) []byte { return []byte(v) }

	testlist2{{
		valuer:  TSVectorArrayFromStringSliceSlice,
		scanner: TSVectorArrayToStringSliceSlice,
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
		valuer:  TSVectorArrayFromByteSliceSliceSlice,
		scanner: TSVectorArrayToByteSliceSliceSlice,
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
		data: []testdata{
			{
				input:  string("{}"),
				output: string("{}")},
			{
				input:  string(`{"'cat' 'fat'","'rat'"}`),
				output: string(`{"'cat' 'fat'",'rat'}`)},
		},
	}, {
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
