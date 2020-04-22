package pgsql

import (
	"testing"
)

func TestXMLArray(t *testing.T) {
	B := func(s string) []byte { return []byte(s) }

	testlist2{{
		valuer:  XMLArrayFromByteSliceSlice,
		scanner: XMLArrayToByteSliceSlice,
		data: []testdata{
			{input: [][]byte(nil), output: [][]byte(nil)},
			{input: [][]byte{}, output: [][]byte{}},
			{
				input:  [][]byte{B(`<foo>bar</foo>`), B(`<bar>hello, world</bar>`)},
				output: [][]byte{B(`<foo>bar</foo>`), B(`<bar>hello, world</bar>`)}},
			{
				input:  [][]byte{B(`<foo>bar</foo>`), nil},
				output: [][]byte{B(`<foo>bar</foo>`), nil}},
		},
	}, {
		data: []testdata{
			{
				input:  string(`{}`),
				output: string(`{}`)},
			{
				input:  string(`{<foo>bar</foo>,"<bar>hello, world</bar>"}`),
				output: string(`{<foo>bar</foo>,"<bar>hello, world</bar>"}`)},
			{
				input:  string(`{<foo>bar</foo>,NULL}`),
				output: string(`{<foo>bar</foo>,NULL}`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte(`{}`),
				output: []byte(`{}`)},
			{
				input:  []byte(`{<foo>bar</foo>,"<bar>hello, world</bar>"}`),
				output: []byte(`{<foo>bar</foo>,"<bar>hello, world</bar>"}`)},
			{
				input:  []byte(`{<foo>bar</foo>,NULL}`),
				output: []byte(`{<foo>bar</foo>,NULL}`)},
		},
	}}.execute(t, "xmlarr")
}
