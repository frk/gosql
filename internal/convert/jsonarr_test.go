package convert

import (
	"testing"
)

func TestJSONArray(t *testing.T) {
	testlist{{
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
	testlist{{
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
