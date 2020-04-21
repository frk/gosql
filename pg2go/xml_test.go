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

	testlist2{{
		valuer:  XML,
		scanner: XML,
		data: []testdata{
			{
				input:  data{},
				output: data{}},
			{
				input:  data{Foo: "some value"},
				output: data{Foo: "some value"}},
		},
	}, {
		data: []testdata{
			{
				input:  string("<document><foo>some value</foo></document>"),
				output: string("<document><foo>some value</foo></document>")},
		},
	}, {
		data: []testdata{
			{
				input:  []byte("<document><foo>some value</foo></document>"),
				output: []byte("<document><foo>some value</foo></document>")},
		},
	}}.execute(t, "xml")
}
