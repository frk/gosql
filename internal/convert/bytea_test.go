package convert

import (
	"testing"
)

func TestByteaWithoutScanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			return nil, new([]byte)
		},
		rows: []test_scanner_row{
			{typ: "bytea", in: nil, want: new([]byte)},
			{typ: "bytea", in: ``, want: &[]byte{}},
			{typ: "bytea", in: `\xdeadbeef`, want: bytesptr(`\xdeadbeef`)},
			{typ: "bytea", in: `\xDEADBEEF`, want: bytesptr(`\xDEADBEEF`)},
			{typ: "bytea", in: "\xfffefdfcfbfaf9f8f7f6f5f4f3f2f1f", want: bytesptr("\xfffefdfcfbfaf9f8f7f6f5f4f3f2f1f")},
			{typ: "bytea", in: `\xfffefdfcfbfaf9f8f7f6f5f4f3f2f1f`, want: bytesptr(`\xfffefdfcfbfaf9f8f7f6f5f4f3f2f1f`)},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			return nil, new(string)
		},
		rows: []test_scanner_row{
			// NOTE: NULL is not supported by string, if the source column
			// is NULLable one should use sql.NullString or *string instead.
			// {typ: "bytea", in: nil, want: new(string)},

			{typ: "bytea", in: ``, want: strptr(``)},
			{typ: "bytea", in: `\xdeadbeef`, want: strptr(`\xdeadbeef`)},
			{typ: "bytea", in: `\xDEADBEEF`, want: strptr(`\xDEADBEEF`)},
			{typ: "bytea", in: `\xfffefdfcfbfaf9f8f7f6f5f4f3f2f1f`, want: strptr(`\xfffefdfcfbfaf9f8f7f6f5f4f3f2f1f`)},
		},
	}}.execute(t)
}
