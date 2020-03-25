package convert

import (
	"testing"
)

func TestByteaArray_Valuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return new(ByteaArrayFromStringSlice)
		},
		rows: []test_valuer_row{
			{typ: "byteaarr", in: nil, want: nil},
			{typ: "byteaarr", in: []string{}, want: strptr(`{}`)},
			{typ: "byteaarr", in: []string{"", ""}, want: strptr(`{"\\x","\\x"}`)},
			{typ: "byteaarr", in: []string{"abc", "def"}, want: strptr(`{"\\x616263","\\x646566"}`)},
			{
				typ:  "byteaarr",
				in:   []string{`\xdeadbeef`, `\xDEADBEEF`},
				want: strptr(`{"\\x5c786465616462656566","\\x5c784445414442454546"}`)},
		},
	}, {
		valuer: func() interface{} {
			return new(ByteaArrayFromByteSliceSlice)
		},
		rows: []test_valuer_row{
			{typ: "byteaarr", in: nil, want: nil},
			{typ: "byteaarr", in: [][]byte{}, want: strptr(`{}`)},
			{typ: "byteaarr", in: [][]byte{{}, {}}, want: strptr(`{"\\x","\\x"}`)},
			{
				typ:  "byteaarr",
				in:   [][]byte{[]byte("abc"), []byte("def")},
				want: strptr(`{"\\x616263","\\x646566"}`)},
			{
				typ:  "byteaarr",
				in:   [][]byte{[]byte(`\xdeadbeef`), []byte(`\xDEADBEEF`)},
				want: strptr(`{"\\x5c786465616462656566","\\x5c784445414442454546"}`)},
		},
	}}.execute(t)
}

func TestByteaArray_Scanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			s := &ByteaArrayToStringSlice{S: new([]string)}
			return s, s.S
		},
		rows: []test_scanner_row{
			{typ: "byteaarr", in: nil, want: new([]string)},
			{typ: "byteaarr", in: `{}`, want: &[]string{}},
			{typ: "byteaarr", in: `{"\\x","\\x"}`, want: &[]string{"", ""}},
			{typ: "byteaarr", in: `{"\\x616263","\\x646566"}`, want: &[]string{"abc", "def"}},
			{
				typ:  "byteaarr",
				in:   `{"\\x5c786465616462656566","\\x5c784445414442454546"}`,
				want: &[]string{`\xdeadbeef`, `\xDEADBEEF`}},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			s := &ByteaArrayToByteSliceSlice{S: new([][]byte)}
			return s, s.S
		},
		rows: []test_scanner_row{
			{typ: "byteaarr", in: nil, want: new([][]byte)},
			{typ: "byteaarr", in: `{}`, want: &[][]byte{}},
			{typ: "byteaarr", in: `{"\\x","\\x"}`, want: &[][]byte{{}, {}}},
			{
				typ:  "byteaarr",
				in:   `{"\\x616263","\\x646566"}`,
				want: &[][]byte{[]byte("abc"), []byte("def")}},
			{
				typ:  "byteaarr",
				in:   `{"\\x5c786465616462656566","\\x5c784445414442454546"}`,
				want: &[][]byte{[]byte(`\xdeadbeef`), []byte(`\xDEADBEEF`)}},
		},
	}}.execute(t)
}
