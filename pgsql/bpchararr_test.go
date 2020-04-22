package pgsql

import (
	"testing"
)

func TestBPCharArray(t *testing.T) {
	testlist2{{
		valuer:  BPCharArrayFromString,
		scanner: BPCharArrayToString,
		data: []testdata{
			{input: "", output: string(``)},
			{input: `abc`, output: string(`abc`)},
			{input: `ab'`, output: string(`ab'`)},
			{input: `a,'`, output: string(`a,'`)},
			{input: `",'`, output: string(`",'`)},
			{input: `魔,'`, output: string(`魔,'`)},
			{input: `""'`, output: string(`""'`)},
			{input: "\t \r", output: string("\t \r")},
			{input: "\f\n", output: string("\f\n")},
			{input: "\v\\", output: string("\v\\")},

			// TODO handle \b and \a
			// {input: "\a", output: string("\a")},
			// {input: "\b", output: string("\b")},
		},
	}, {
		valuer:  BPCharArrayFromByteSlice,
		scanner: BPCharArrayToByteSlice,
		data: []testdata{
			{input: []byte(nil), output: []byte(nil)},
			{input: []byte{}, output: []byte{}},
			{input: []byte(`abc`), output: []byte{'a', 'b', 'c'}},
			{input: []byte(`ab'`), output: []byte{'a', 'b', '\''}},
			{input: []byte(`a,'`), output: []byte{'a', ',', '\''}},
			{input: []byte(`",'`), output: []byte{'"', ',', '\''}},
			{input: []byte(`""'`), output: []byte{'"', '"', '\''}},
			{input: []byte("\t \r"), output: []byte{'\t', ' ', '\r'}},
			{input: []byte("\f\n"), output: []byte{'\f', '\n'}},
			{input: []byte("\v\\"), output: []byte{'\v', '\\'}},

			// TODO handle (pq: invalid byte sequence for encoding "UTF8": 0xe9 0x2c 0xad)
			//{input: []byte(`魔,'`), output: []byte{233, 173, 148, ',', '\''}},

			// TODO handle \b and \a
			// {input: []byte("\b"), output: []byte{'\b'}},
			// {input: []byte("\a"), output: []byte{'\a'}},
		},
	}, {
		valuer:  BPCharArrayFromRuneSlice,
		scanner: BPCharArrayToRuneSlice,
		data: []testdata{
			{input: []rune(nil), output: []rune(nil)},
			{input: []rune{}, output: []rune{}},
			{input: []rune(`abc`), output: []rune{'a', 'b', 'c'}},
			{input: []rune(`ab'`), output: []rune{'a', 'b', '\''}},
			{input: []rune(`a,'`), output: []rune{'a', ',', '\''}},
			{input: []rune(`",'`), output: []rune{'"', ',', '\''}},
			{input: []rune(`""'`), output: []rune{'"', '"', '\''}},
			{input: []rune("\t \r"), output: []rune{'\t', ' ', '\r'}},
			{input: []rune("\f\n"), output: []rune{'\f', '\n'}},
			{input: []rune("\v\\"), output: []rune{'\v', '\\'}},
			{input: []rune(`魔,'`), output: []rune{'魔', ',', '\''}},

			// TODO handle \b and \a
			// {input: []rune("\a"), output: []rune{'\a'}},
			// {input: []rune("\b"), output: []rune{'\b'}},
		},
	}, {
		valuer:  BPCharArrayFromStringSlice,
		scanner: BPCharArrayToStringSlice,
		data: []testdata{
			{input: []string(nil), output: []string(nil)},
			{input: []string{}, output: []string{}},
			{input: []string{"a", "b", "c"}, output: []string{"a", "b", "c"}},
			{input: []string{"a", "b", "'"}, output: []string{"a", "b", "'"}},
			{input: []string{"a", ",", "'"}, output: []string{"a", ",", "'"}},
			{input: []string{"\"", ",", "'"}, output: []string{"\"", ",", "'"}},
			{input: []string{"\"", "\"", "'"}, output: []string{"\"", "\"", "'"}},
			{input: []string{"\t", " ", "\r"}, output: []string{"\t", " ", "\r"}},
			{input: []string{"\f", "\n"}, output: []string{"\f", "\n"}},
			{input: []string{"\v", "\\"}, output: []string{"\v", "\\"}},
			{input: []string{"魔", ",", "'"}, output: []string{"魔", ",", "'"}},

			// TODO handle \b and \a
			// {input: []string{"\b"}, output: []string{"\b"}},
			// {input: []string{"\a"}, output: []string{"\a"}},
		},
	}, {
		data: []testdata{
			{input: string(`{}`), output: string(`{}`)},
			{
				input:  string(`{a,b,c}`),
				output: string(`{a,b,c}`)},
			{
				input:  string(`{a,b,'}`),
				output: string(`{a,b,'}`)},
			{
				input:  string(`{a,",",'}`),
				output: string(`{a,",",'}`)},
			{
				input:  string(`{"\"",",",'}`),
				output: string(`{"\"",",",'}`)},
			{
				input:  string(`{魔,",",'}`),
				output: string(`{魔,",",'}`)},
			{
				input:  string(`{"\"","\"",'}`),
				output: string(`{"\"","\"",'}`)},
			{
				input:  string("{\"\t\",\" \",\"\r\"}"),
				output: string("{\"\t\",\" \",\"\r\"}")},
			{
				input:  string("{\"\f\",\"\n\"}"),
				output: string("{\"\f\",\"\n\"}")},
			{
				input:  string("{\"\v\",\"\\\\\"}"),
				output: string("{\"\v\",\"\\\\\"}")},

			// TODO handle \b and \a
			// { input: string("{\"\a\"}"), output: string("{\"\a\"}")},
			// { input: string("{\"\b\"}"), output: string("{\"\b\"}")},
		},
	}, {
		data: []testdata{
			{input: []byte(`{}`), output: []byte(`{}`)},
			{
				input:  []byte(`{a,b,c}`),
				output: []byte(`{a,b,c}`)},
			{
				input:  []byte(`{a,b,'}`),
				output: []byte(`{a,b,'}`)},
			{
				input:  []byte(`{a,",",'}`),
				output: []byte(`{a,",",'}`)},
			{
				input:  []byte(`{"\"",",",'}`),
				output: []byte(`{"\"",",",'}`)},
			{
				input:  []byte(`{魔,",",'}`),
				output: []byte(`{魔,",",'}`)},
			{
				input:  []byte(`{"\"","\"",'}`),
				output: []byte(`{"\"","\"",'}`)},
			{
				input:  []byte("{\"\t\",\" \",\"\r\"}"),
				output: []byte("{\"\t\",\" \",\"\r\"}")},
			{
				input:  []byte("{\"\f\",\"\n\"}"),
				output: []byte("{\"\f\",\"\n\"}")},
			{
				input:  []byte("{\"\v\",\"\\\\\"}"),
				output: []byte("{\"\v\",\"\\\\\"}")},

			// TODO handle \b and \a
			// { input: []byte("{\"\a\"}"), output: []byte("{\"\a\"}")},
			// { input: []byte("{\"\b\"}"), output: []byte("{\"\b\"}")},
		},
	}}.execute(t, "bpchararr")
}
