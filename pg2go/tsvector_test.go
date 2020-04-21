package pg2go

import (
	"testing"
)

func TestTSVector(t *testing.T) {
	B := func(v string) []byte { return []byte(v) }

	testlist2{{
		valuer:  TSVectorFromStringSlice,
		scanner: TSVectorToStringSlice,
		data: []testdata{
			{input: []string(nil), output: []string(nil)},
			{input: []string{}, output: []string{}},
			{
				input:  []string{"'cat'", "'fat'", "'rat'"},
				output: []string{"'cat'", "'fat'", "'rat'"}},
			{
				input:  []string{"'cat':2", "'fat':7,11", "'rat'"},
				output: []string{"'cat':2", "'fat':7,11", "'rat'"}},
			{
				input:  []string{"'cat':2", "'fat':7,11", "'rat':2B,4C"},
				output: []string{"'cat':2", "'fat':7,11", "'rat':2B,4C"}},
		},
	}, {
		valuer:  TSVectorFromByteSliceSlice,
		scanner: TSVectorToByteSliceSlice,
		data: []testdata{
			{input: [][]byte(nil), output: [][]byte(nil)},
			{input: [][]byte{}, output: [][]byte{}},
			{
				input:  [][]byte{B("'cat'"), B("'fat'"), B("'rat'")},
				output: [][]byte{B("'cat'"), B("'fat'"), B("'rat'")}},
			{
				input:  [][]byte{B("'cat':2"), B("'fat':7,11"), B("'rat'")},
				output: [][]byte{B("'cat':2"), B("'fat':7,11"), B("'rat'")}},
			{
				input:  [][]byte{B("'cat':2"), B("'fat':7,11"), B("'rat':2B,4C")},
				output: [][]byte{B("'cat':2"), B("'fat':7,11"), B("'rat':2B,4C")}},
		},
	}, {
		data: []testdata{
			{
				input:  string(""),
				output: string("")},
			{
				input:  string("'cat' 'fat' 'rat'"),
				output: string("'cat' 'fat' 'rat'")},
			{
				input:  string("'cat':2 'fat':7,11 'rat'"),
				output: string("'cat':2 'fat':7,11 'rat'")},
			{
				input:  string("'cat':2 'fat':7,11 'rat':2B,4C"),
				output: string("'cat':2 'fat':7,11 'rat':2B,4C")},
		},
	}, {
		data: []testdata{
			{
				input:  []byte(""),
				output: []byte("")},
			{
				input:  []byte("'cat' 'fat' 'rat'"),
				output: []byte("'cat' 'fat' 'rat'")},
			{
				input:  []byte("'cat':2 'fat':7,11 'rat'"),
				output: []byte("'cat':2 'fat':7,11 'rat'")},
			{
				input:  []byte("'cat':2 'fat':7,11 'rat':2B,4C"),
				output: []byte("'cat':2 'fat':7,11 'rat':2B,4C")},
		},
	}}.execute(t, "tsvector")
}
