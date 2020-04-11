package convert

import (
	"testing"
)

func TestXMLArray(t *testing.T) {
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
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			s := new([]byte)
			return s, s
		},
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
