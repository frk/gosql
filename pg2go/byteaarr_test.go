package pg2go

import (
	"testing"
)

func TestByteaArray(t *testing.T) {
	testlist2{{
		valuer:  ByteaArrayFromStringSlice,
		scanner: ByteaArrayToStringSlice,
		data: []testdata{
			{input: []string(nil), output: []string(nil)},
			{input: []string{}, output: []string{}},
			{
				input:  []string{"", ""},
				output: []string{"", ""}},
			{
				input:  []string{"abc", "def"},
				output: []string{"abc", "def"}},
			{
				input:  []string{`\xdeadbeef`, `\xDEADBEEF`},
				output: []string{`\xdeadbeef`, `\xDEADBEEF`}},
		},
	}, {
		valuer:  ByteaArrayFromByteSliceSlice,
		scanner: ByteaArrayToByteSliceSlice,
		data: []testdata{
			{input: [][]byte(nil), output: [][]byte(nil)},
			{input: [][]byte{}, output: [][]byte{}},
			{
				input:  [][]byte{{}, {}},
				output: [][]byte{{}, {}}},
			{
				input:  [][]byte{[]byte("abc"), []byte("def")},
				output: [][]byte{[]byte("abc"), []byte("def")}},
			{
				input:  [][]byte{[]byte(`\xdeadbeef`), []byte(`\xDEADBEEF`)},
				output: [][]byte{[]byte(`\xdeadbeef`), []byte(`\xDEADBEEF`)}},
		},
	}, {
		data: []testdata{
			{input: string(`{}`), output: string(`{}`)},
			{
				input:  string(`{"\\x","\\x"}`),
				output: string(`{"\\x","\\x"}`)},
			{
				input:  string(`{"\\x616263","\\x646566"}`),
				output: string(`{"\\x616263","\\x646566"}`)},
			{
				input:  string(`{"\\x5c786465616462656566","\\x5c784445414442454546"}`),
				output: string(`{"\\x5c786465616462656566","\\x5c784445414442454546"}`)},
		},
	}, {
		data: []testdata{
			{input: []byte(`{}`), output: []byte(`{}`)},
			{
				input:  []byte(`{"\\x","\\x"}`),
				output: []byte(`{"\\x","\\x"}`)},
			{
				input:  []byte(`{"\\x616263","\\x646566"}`),
				output: []byte(`{"\\x616263","\\x646566"}`)},
			{
				input:  []byte(`{"\\x5c786465616462656566","\\x5c784445414442454546"}`),
				output: []byte(`{"\\x5c786465616462656566","\\x5c784445414442454546"}`)},
		},
	}}.execute(t, "byteaarr")
}
