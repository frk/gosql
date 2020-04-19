package convert

import (
	"testing"
)

func TestJSONArray(t *testing.T) {
	B := func(s string) []byte { return []byte(s) }

	testlist{{
		valuer: func() interface{} {
			return new(JSONArrayFromByteSliceSlice)
		},
		scanner: func() (interface{}, interface{}) {
			v := JSONArrayToByteSliceSlice{Val: new([][]byte)}
			return v, v.Val
		},
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
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			s := new(string)
			return s, s
		},
		data: []testdata{
			{
				input:  string(`{"{\"foo\":[\"bar\", \"baz\", 123]}","[\"foo\", 123]"}`),
				output: string(`{"{\"foo\":[\"bar\", \"baz\", 123]}","[\"foo\", 123]"}`)},
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
			{
				input:  []byte(`{"{\"foo\":[\"bar\", \"baz\", 123]}","[\"foo\", 123]"}`),
				output: []byte(`{"{\"foo\":[\"bar\", \"baz\", 123]}","[\"foo\", 123]"}`)},
		},
	}}.execute(t, "jsonarr")
}

func TestJSONBArray(t *testing.T) {
	B := func(s string) []byte { return []byte(s) }

	testlist{{
		valuer: func() interface{} {
			return new(JSONArrayFromByteSliceSlice)
		},
		scanner: func() (interface{}, interface{}) {
			v := JSONArrayToByteSliceSlice{Val: new([][]byte)}
			return v, v.Val
		},
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
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			s := new(string)
			return s, s
		},
		data: []testdata{
			{
				input:  string(`{"{\"foo\":[\"bar\",\"baz\",123]}","[\"foo\",123]"}`),
				output: string(`{"{\"foo\": [\"bar\", \"baz\", 123]}","[\"foo\", 123]"}`)},
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
			{
				input:  []byte(`{"{\"foo\":[\"bar\",\"baz\",123]}","[\"foo\",123]"}`),
				output: []byte(`{"{\"foo\": [\"bar\", \"baz\", 123]}","[\"foo\", 123]"}`)},
		},
	}}.execute(t, "jsonbarr")
}
