package pgsql

import (
	"testing"
)

func TestJSONArray(t *testing.T) {
	B := func(s string) []byte { return []byte(s) }

	testlist2{{
		valuer:  JSONArrayFromByteSliceSlice,
		scanner: JSONArrayToByteSliceSlice,
		data: []testdata{
			{input: [][]byte(nil), output: [][]byte(nil)},
			{input: [][]byte{}, output: [][]byte{}},
			{
				input: [][]byte{
					B(`{"foo":["bar", "baz", 123]}`),
					B(`["foo", 123]`)},
				output: [][]byte{
					B(`{"foo":["bar", "baz", 123]}`),
					B(`["foo", 123]`)}},
			{
				input:  [][]byte{B(`{"foo":["bar", "baz", 123]}`), nil},
				output: [][]byte{B(`{"foo":["bar", "baz", 123]}`), nil}},
		},
	}, {
		data: []testdata{
			{
				input:  string(`{"{\"foo\":[\"bar\", \"baz\", 123]}","[\"foo\", 123]"}`),
				output: string(`{"{\"foo\":[\"bar\", \"baz\", 123]}","[\"foo\", 123]"}`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte(`{"{\"foo\":[\"bar\", \"baz\", 123]}","[\"foo\", 123]"}`),
				output: []byte(`{"{\"foo\":[\"bar\", \"baz\", 123]}","[\"foo\", 123]"}`)},
		},
	}}.execute(t, "jsonarr")
}

func TestJSONBArray(t *testing.T) {
	B := func(s string) []byte { return []byte(s) }

	testlist2{{
		valuer:  JSONArrayFromByteSliceSlice,
		scanner: JSONArrayToByteSliceSlice,
		data: []testdata{
			{input: [][]byte(nil), output: [][]byte(nil)},
			{input: [][]byte{}, output: [][]byte{}},
			{
				input: [][]byte{
					B(`{"foo": ["bar", "baz", 123]}`),
					B(`["foo", 123]`)},
				output: [][]byte{
					B(`{"foo": ["bar", "baz", 123]}`),
					B(`["foo", 123]`)}},
			{
				input:  [][]byte{B(`{"foo": ["bar", "baz", 123]}`), nil},
				output: [][]byte{B(`{"foo": ["bar", "baz", 123]}`), nil}},
		},
	}, {
		data: []testdata{
			{
				input:  string(`{"{\"foo\":[\"bar\",\"baz\",123]}","[\"foo\",123]"}`),
				output: string(`{"{\"foo\": [\"bar\", \"baz\", 123]}","[\"foo\", 123]"}`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte(`{"{\"foo\":[\"bar\",\"baz\",123]}","[\"foo\",123]"}`),
				output: []byte(`{"{\"foo\": [\"bar\", \"baz\", 123]}","[\"foo\", 123]"}`)},
		},
	}}.execute(t, "jsonbarr")
}
