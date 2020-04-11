package convert

import (
	"testing"
)

func TestBytea(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			return nil, new(string)
		},
		data: []testdata{
			{input: string(``), output: string(``)},
			{input: string(`abc`), output: string(`abc`)},
			{
				input:  string(`\xdeadbeef`),
				output: string(`\xdeadbeef`)},
			{
				input:  string(`\xDEADBEEF`),
				output: string(`\xDEADBEEF`)},
			{
				input:  string("\xfffefdfcfbfaf9f8f7f6f5f4f3f2f1f"),
				output: string("\xfffefdfcfbfaf9f8f7f6f5f4f3f2f1f")},
			{
				input:  string(`\xfffefdfcfbfaf9f8f7f6f5f4f3f2f1f`),
				output: string(`\xfffefdfcfbfaf9f8f7f6f5f4f3f2f1f`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			return nil, new([]byte)
		},
		data: []testdata{
			{input: nil, output: []byte(nil)},
			{input: []byte(``), output: []byte(``)},
			{
				input:  []byte(`\xdeadbeef`),
				output: []byte(`\xdeadbeef`)},
			{
				input:  []byte(`\xDEADBEEF`),
				output: []byte(`\xDEADBEEF`)},
			{
				input:  []byte("\xfffefdfcfbfaf9f8f7f6f5f4f3f2f1f"),
				output: []byte("\xfffefdfcfbfaf9f8f7f6f5f4f3f2f1f")},
			{
				input:  []byte(`\xfffefdfcfbfaf9f8f7f6f5f4f3f2f1f`),
				output: []byte(`\xfffefdfcfbfaf9f8f7f6f5f4f3f2f1f`)},
		},
	}}.execute(t, "bytea")
}
