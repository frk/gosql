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
			{input: string(``), output: strptr(``)},
			{input: string(`abc`), output: strptr(`abc`)},
			{
				input:  string(`\xdeadbeef`),
				output: strptr(`\xdeadbeef`)},
			{
				input:  string(`\xDEADBEEF`),
				output: strptr(`\xDEADBEEF`)},
			{
				input:  string("\xfffefdfcfbfaf9f8f7f6f5f4f3f2f1f"),
				output: strptr("\xfffefdfcfbfaf9f8f7f6f5f4f3f2f1f")},
			{
				input:  string(`\xfffefdfcfbfaf9f8f7f6f5f4f3f2f1f`),
				output: strptr(`\xfffefdfcfbfaf9f8f7f6f5f4f3f2f1f`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			return nil, new([]byte)
		},
		data: []testdata{
			{input: nil, output: new([]byte)},
			{input: []byte(``), output: bytesptr(``)},
			{
				input:  []byte(`\xdeadbeef`),
				output: bytesptr(`\xdeadbeef`)},
			{
				input:  []byte(`\xDEADBEEF`),
				output: bytesptr(`\xDEADBEEF`)},
			{
				input:  []byte("\xfffefdfcfbfaf9f8f7f6f5f4f3f2f1f"),
				output: bytesptr("\xfffefdfcfbfaf9f8f7f6f5f4f3f2f1f")},
			{
				input:  []byte(`\xfffefdfcfbfaf9f8f7f6f5f4f3f2f1f`),
				output: bytesptr(`\xfffefdfcfbfaf9f8f7f6f5f4f3f2f1f`)},
		},
	}}.execute(t, "bytea")
}
