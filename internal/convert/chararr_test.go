package convert

import (
	"testing"
)

func TestCharArray(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(CharArrayFromString)
		},
		scanner: func() (interface{}, interface{}) {
			s := CharArrayToString{Val: new(string)}
			return s, s.Val
		},
		data: []testdata{
			{input: "", output: strptr(``)},
			{input: `abc`, output: strptr(`abc`)},
			{input: `ab'`, output: strptr(`ab'`)},
			{input: `a,'`, output: strptr(`a,'`)},
			{input: `",'`, output: strptr(`",'`)},
			{input: `魔,'`, output: strptr(`魔,'`)},
			{input: `""'`, output: strptr(`""'`)},
			{input: "\t \r", output: strptr("\t \r")},
			{input: "\f\n", output: strptr("\f\n")},
			{input: "\v\\", output: strptr("\v\\")},

			// TODO handle \b and \a
			// {input: "\a", output: strptr("\a")},
			// {input: "\b", output: strptr("\b")},
		},
	}, {
		valuer: func() interface{} {
			return new(CharArrayFromByteSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := CharArrayToByteSlice{Val: new([]byte)}
			return s, s.Val
		},
		data: []testdata{
			{input: nil, output: new([]byte)},
			{input: []byte{}, output: &[]byte{}},
			{input: []byte(`abc`), output: &[]byte{'a', 'b', 'c'}},
			{input: []byte(`ab'`), output: &[]byte{'a', 'b', '\''}},
			{input: []byte(`a,'`), output: &[]byte{'a', ',', '\''}},
			{input: []byte(`",'`), output: &[]byte{'"', ',', '\''}},
			{input: []byte(`""'`), output: &[]byte{'"', '"', '\''}},
			{input: []byte("\t \r"), output: &[]byte{'\t', ' ', '\r'}},
			{input: []byte("\f\n"), output: &[]byte{'\f', '\n'}},
			{input: []byte("\v\\"), output: &[]byte{'\v', '\\'}},

			// TODO handle (pq: invalid byte sequence for encoding "UTF8": 0xe9 0x2c 0xad)
			//{input: []byte(`魔,'`), output: &[]byte{233, 173, 148, ',', '\''}},

			// TODO handle \b and \a
			// {input: []byte("\b"), output: &[]byte{'\b'}},
			// {input: []byte("\a"), output: &[]byte{'\a'}},
		},
	}, {
		valuer: func() interface{} {
			return new(CharArrayFromRuneSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := CharArrayToRuneSlice{Val: new([]rune)}
			return s, s.Val
		},
		data: []testdata{
			{input: nil, output: new([]rune)},
			{input: []rune{}, output: &[]rune{}},
			{input: []rune(`abc`), output: &[]rune{'a', 'b', 'c'}},
			{input: []rune(`ab'`), output: &[]rune{'a', 'b', '\''}},
			{input: []rune(`a,'`), output: &[]rune{'a', ',', '\''}},
			{input: []rune(`",'`), output: &[]rune{'"', ',', '\''}},
			{input: []rune(`""'`), output: &[]rune{'"', '"', '\''}},
			{input: []rune("\t \r"), output: &[]rune{'\t', ' ', '\r'}},
			{input: []rune("\f\n"), output: &[]rune{'\f', '\n'}},
			{input: []rune("\v\\"), output: &[]rune{'\v', '\\'}},
			{input: []rune(`魔,'`), output: &[]rune{'魔', ',', '\''}},

			// TODO handle \b and \a
			// {input: []rune("\a"), output: &[]rune{'\a'}},
			// {input: []rune("\b"), output: &[]rune{'\b'}},
		},
	}, {
		valuer: func() interface{} {
			return new(CharArrayFromStringSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := CharArrayToStringSlice{Val: new([]string)}
			return s, s.Val
		},
		data: []testdata{
			{input: nil, output: new([]string)},
			{input: []string{}, output: &[]string{}},
			{input: []string{"a", "b", "c"}, output: &[]string{"a", "b", "c"}},
			{input: []string{"a", "b", "'"}, output: &[]string{"a", "b", "'"}},
			{input: []string{"a", ",", "'"}, output: &[]string{"a", ",", "'"}},
			{input: []string{"\"", ",", "'"}, output: &[]string{"\"", ",", "'"}},
			{input: []string{"\"", "\"", "'"}, output: &[]string{"\"", "\"", "'"}},
			{input: []string{"\t", " ", "\r"}, output: &[]string{"\t", " ", "\r"}},
			{input: []string{"\f", "\n"}, output: &[]string{"\f", "\n"}},
			{input: []string{"\v", "\\"}, output: &[]string{"\v", "\\"}},
			{input: []string{"魔", ",", "'"}, output: &[]string{"魔", ",", "'"}},

			// TODO handle \b and \a
			// {input: []string{"\b"}, output: &[]string{"\b"}},
			// {input: []string{"\a"}, output: &[]string{"\a"}},
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
			{input: string(`{a,b,c}`), output: strptr(`{a,b,c}`)},
			{input: string(`{a,b,'}`), output: strptr(`{a,b,'}`)},
			{input: string(`{a,",",'}`), output: strptr(`{a,",",'}`)},
			{input: string(`{"\"",",",'}`), output: strptr(`{"\"",",",'}`)},
			{input: string(`{魔,",",'}`), output: strptr(`{魔,",",'}`)},
			{input: string(`{"\"","\"",'}`), output: strptr(`{"\"","\"",'}`)},
			{
				input:  string("{\"\t\",\" \",\"\r\"}"),
				output: strptr("{\"\t\",\" \",\"\r\"}")},
			{
				input:  string("{\"\f\",\"\n\"}"),
				output: strptr("{\"\f\",\"\n\"}")},
			{
				input:  string("{\"\v\",\"\\\\\"}"),
				output: strptr("{\"\v\",\"\\\\\"}")},

			// TODO handle \b and \a
			// {input: string("{\"\a\"}"), output: strptr("{\"\a\"}")},
			// {input: string("{\"\b\"}"), output: strptr("{\"\b\"}")},
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
			{input: []byte(`{a,b,c}`), output: bytesptr(`{a,b,c}`)},
			{input: []byte(`{a,b,'}`), output: bytesptr(`{a,b,'}`)},
			{input: []byte(`{a,",",'}`), output: bytesptr(`{a,",",'}`)},
			{input: []byte(`{"\"",",",'}`), output: bytesptr(`{"\"",",",'}`)},
			{input: []byte(`{魔,",",'}`), output: bytesptr(`{魔,",",'}`)},
			{input: []byte(`{"\"","\"",'}`), output: bytesptr(`{"\"","\"",'}`)},
			{
				input:  []byte("{\"\t\",\" \",\"\r\"}"),
				output: bytesptr("{\"\t\",\" \",\"\r\"}")},
			{
				input:  []byte("{\"\f\",\"\n\"}"),
				output: bytesptr("{\"\f\",\"\n\"}")},
			{
				input:  []byte("{\"\v\",\"\\\\\"}"),
				output: bytesptr("{\"\v\",\"\\\\\"}")},

			// TODO handle \b and \a
			// {input: []byte("{\"\a\"}"), output: bytesptr("{\"\a\"}")},
			// {input: []byte("{\"\b\"}"), output: bytesptr("{\"\b\"}")},
		},
	}}.execute(t, "chararr")
}
