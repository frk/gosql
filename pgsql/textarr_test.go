package pgsql

import (
	"testing"
)

func TestTextArray(t *testing.T) {
	B := func(s string) []byte { return []byte(s) }

	testlist2{{
		valuer:  TextArrayFromStringSlice,
		scanner: TextArrayToStringSlice,
		data: []testdata{
			{input: []string(nil), output: []string(nil)},
			{
				input:  []string{"foo", "bar"},
				output: []string{"foo", "bar"}},
			{
				input:  []string{"foo,bar", "foo bar"},
				output: []string{"foo,bar", "foo bar"}},
			{
				input:  []string{"foo", "", "NULL"},
				output: []string{"foo", "", "NULL"}},
			{
				input:  []string{"foo\""},
				output: []string{"foo\""}},
			{
				input:  []string{"foo\\\\\""},
				output: []string{"foo\\\\\""}},
		},
	}, {
		valuer:  TextArrayFromByteSliceSlice,
		scanner: TextArrayToByteSliceSlice,
		data: []testdata{
			{input: [][]byte(nil), output: [][]byte(nil)},
			{
				input:  [][]byte{B("foo"), B("bar")},
				output: [][]byte{B("foo"), B("bar")}},
			{
				input:  [][]byte{B("foo,bar"), B("foo bar")},
				output: [][]byte{B("foo,bar"), B("foo bar")}},
			{
				input:  [][]byte{B("foo"), nil, B("NULL")},
				output: [][]byte{B("foo"), nil, B("NULL")}},
			{
				input:  [][]byte{B("foo\"")},
				output: [][]byte{B("foo\"")}},
			{
				input:  [][]byte{B("foo\\\\\"")},
				output: [][]byte{B("foo\\\\\"")}},
		},
	}, {
		data: []testdata{
			{input: string(`{}`), output: string(`{}`)},
			{input: string(`{foo,bar}`), output: string(`{foo,bar}`)},
			{
				input:  string(`{"foo,bar","foo bar"}`),
				output: string(`{"foo,bar","foo bar"}`)},
		},
	}, {
		data: []testdata{
			{input: []byte(`{}`), output: []byte(`{}`)},
			{input: []byte(`{foo,bar}`), output: []byte(`{foo,bar}`)},
			{
				input:  []byte(`{"foo,bar","foo bar"}`),
				output: []byte(`{"foo,bar","foo bar"}`)},
		},
	}}.execute(t, "textarr")
}
