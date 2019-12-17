package convert

import (
	"database/sql"
	"testing"
)

func TestByteaWithoutScanner(t *testing.T) {
	test_table{{
		dest: func() interface{} {
			return new([]byte)
		},
		rows: []testrow{
			{typ: "bytea", in: nil, want: new([]byte)},
			{typ: "bytea", in: ``, want: &[]byte{}},
			{typ: "bytea", in: `\xdeadbeef`, want: bytesptr(`\xdeadbeef`)},
			{typ: "bytea", in: `\xDEADBEEF`, want: bytesptr(`\xDEADBEEF`)},
			{typ: "bytea", in: "\xfffefdfcfbfaf9f8f7f6f5f4f3f2f1f", want: bytesptr(`\xfffefdfcfbfaf9f8f7f6f5f4f3f2f1f`)},
		},
	}, {
		dest: func() interface{} {
			return new(string)
		},
		rows: []testrow{
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
