package pg2go

import (
	"encoding/xml"
	"testing"
)

func TestXML(t *testing.T) {
	type data struct {
		xml.Name `xml:"document"`
		Foo      string `xml:"foo"`
	}

	testlist{{
		valuer: func() interface{} {
			return new(XML)
		},
		scanner: func() (interface{}, interface{}) {
			s := XML{Val: new(data)}
			return s, s.Val
		},
		data: []testdata{
			{
				input:  data{},
				output: data{}},
			{
				input:  data{Foo: "some value"},
				output: data{Foo: "some value"}},
		},
	}, {
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			d := new(string)
			return d, d
		},
		data: []testdata{
			{
				input:  string("<document><foo>some value</foo></document>"),
				output: string("<document><foo>some value</foo></document>")},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			d := new([]byte)
			return d, d
		},
		data: []testdata{
			{
				input:  []byte("<document><foo>some value</foo></document>"),
				output: []byte("<document><foo>some value</foo></document>")},
		},
	}}.execute(t, "xml")
}
