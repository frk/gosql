package convert

import (
	"testing"
)

func TestByteaArray(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(ByteaArrayFromStringSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := &ByteaArrayToStringSlice{Val: new([]string)}
			return s, s.Val
		},
		data: []testdata{
			{input: nil, output: new([]string)},
			{input: []string{}, output: &[]string{}},
			{
				input:  []string{"", ""},
				output: &[]string{"", ""}},
			{
				input:  []string{"abc", "def"},
				output: &[]string{"abc", "def"}},
			{
				input:  []string{`\xdeadbeef`, `\xDEADBEEF`},
				output: &[]string{`\xdeadbeef`, `\xDEADBEEF`}},
		},
	}, {
		valuer: func() interface{} {
			return new(ByteaArrayFromByteSliceSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := &ByteaArrayToByteSliceSlice{Val: new([][]byte)}
			return s, s.Val
		},
		data: []testdata{
			{input: nil, output: new([][]byte)},
			{input: [][]byte{}, output: &[][]byte{}},
			{
				input:  [][]byte{{}, {}},
				output: &[][]byte{{}, {}}},
			{
				input:  [][]byte{[]byte("abc"), []byte("def")},
				output: &[][]byte{[]byte("abc"), []byte("def")}},
			{
				input:  [][]byte{[]byte(`\xdeadbeef`), []byte(`\xDEADBEEF`)},
				output: &[][]byte{[]byte(`\xdeadbeef`), []byte(`\xDEADBEEF`)}},
		},
	}, {
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			s := new(string)
			return s, s
		},
		data: []testdata{
			{input: string(`{}`), output: strptr(`{}`)},
			{
				input:  string(`{"\\x","\\x"}`),
				output: strptr(`{"\\x","\\x"}`)},
			{
				input:  string(`{"\\x616263","\\x646566"}`),
				output: strptr(`{"\\x616263","\\x646566"}`)},
			{
				input:  string(`{"\\x5c786465616462656566","\\x5c784445414442454546"}`),
				output: strptr(`{"\\x5c786465616462656566","\\x5c784445414442454546"}`)},
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
			{input: []byte(`{}`), output: bytesptr(`{}`)},
			{
				input:  []byte(`{"\\x","\\x"}`),
				output: bytesptr(`{"\\x","\\x"}`)},
			{
				input:  []byte(`{"\\x616263","\\x646566"}`),
				output: bytesptr(`{"\\x616263","\\x646566"}`)},
			{
				input:  []byte(`{"\\x5c786465616462656566","\\x5c784445414442454546"}`),
				output: bytesptr(`{"\\x5c786465616462656566","\\x5c784445414442454546"}`)},
		},
	}}.execute(t, "byteaarr")
}
