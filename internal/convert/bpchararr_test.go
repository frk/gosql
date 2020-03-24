package convert

import (
	"testing"
)

func TestBPCharArray_Valuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return new(BPCharArrayFromString)
		},
		rows: []test_valuer_row{
			{typ: "chararr", in: "", want: strptr(`{}`)},
			{typ: "chararr", in: `abc`, want: strptr(`{a,b,c}`)},
			{typ: "chararr", in: `ab'`, want: strptr(`{a,b,'}`)},
			{typ: "chararr", in: `a,'`, want: strptr(`{a,",",'}`)},
			{typ: "chararr", in: `",'`, want: strptr(`{"\"",",",'}`)},
			{typ: "chararr", in: `魔,'`, want: strptr(`{魔,",",'}`)},
			{typ: "chararr", in: `""'`, want: strptr(`{"\"","\"",'}`)},
			{typ: "chararr", in: "\t \r", want: strptr("{\"\t\",\" \",\"\r\"}")},
			{typ: "chararr", in: "\f\n", want: strptr("{\"\f\",\"\n\"}")},
			{typ: "chararr", in: "\v\\", want: strptr("{\"\v\",\"\\\\\"}")},

			// TODO handle \b and \a
			// {typ: "chararr", in: "\b\a", want: strptr("{\"\b\",\"\a\"}")},
		},
	}, {
		valuer: func() interface{} {
			return new(BPCharArrayFromByteSlice)
		},
		rows: []test_valuer_row{
			{typ: "chararr", in: nil, want: nil},
			{typ: "chararr", in: []byte{}, want: strptr(`{}`)},
			{typ: "chararr", in: []byte(`abc`), want: strptr(`{a,b,c}`)},
			{typ: "chararr", in: []byte(`ab'`), want: strptr(`{a,b,'}`)},
			{typ: "chararr", in: []byte(`a,'`), want: strptr(`{a,",",'}`)},
			{typ: "chararr", in: []byte(`",'`), want: strptr(`{"\"",",",'}`)},
			{typ: "chararr", in: []byte(`""'`), want: strptr(`{"\"","\"",'}`)},
			{typ: "chararr", in: []byte("\t \r"), want: strptr("{\"\t\",\" \",\"\r\"}")},
			{typ: "chararr", in: []byte("\f\n"), want: strptr("{\"\f\",\"\n\"}")},
			{typ: "chararr", in: []byte("\v\\"), want: strptr("{\"\v\",\"\\\\\"}")},

			// TODO handle \b and \a
			//{typ: "chararr", in: []byte("\b\a"), want: strptr("{\"\b\",\"\a\"}")},
		},
	}, {
		valuer: func() interface{} {
			return new(BPCharArrayFromRuneSlice)
		},
		rows: []test_valuer_row{
			{typ: "chararr", in: nil, want: nil},
			{typ: "chararr", in: []rune{}, want: strptr(`{}`)},
			{typ: "chararr", in: []rune(`abc`), want: strptr(`{a,b,c}`)},
			{typ: "chararr", in: []rune(`ab'`), want: strptr(`{a,b,'}`)},
			{typ: "chararr", in: []rune(`a,'`), want: strptr(`{a,",",'}`)},
			{typ: "chararr", in: []rune(`",'`), want: strptr(`{"\"",",",'}`)},
			{typ: "chararr", in: []rune(`""'`), want: strptr(`{"\"","\"",'}`)},
			{typ: "chararr", in: []rune("\t \r"), want: strptr("{\"\t\",\" \",\"\r\"}")},
			{typ: "chararr", in: []rune("\f\n"), want: strptr("{\"\f\",\"\n\"}")},
			{typ: "chararr", in: []rune("\v\\"), want: strptr("{\"\v\",\"\\\\\"}")},

			// TODO handle \b and \a
			//{typ: "chararr", in: []byte("\b\a"), want: strptr("{\"\b\",\"\a\"}")},
		},
	}, {
		valuer: func() interface{} {
			return new(BPCharArrayFromStringSlice)
		},
		rows: []test_valuer_row{
			{typ: "chararr", in: nil, want: nil},
			{typ: "chararr", in: []string{}, want: strptr(`{}`)},
			{typ: "chararr", in: []string{"a", "b", "c"}, want: strptr(`{a,b,c}`)},
			{typ: "chararr", in: []string{"a", "b", "'"}, want: strptr(`{a,b,'}`)},
			{typ: "chararr", in: []string{"a", ",", "'"}, want: strptr(`{a,",",'}`)},
			{typ: "chararr", in: []string{"\"", ",", "'"}, want: strptr(`{"\"",",",'}`)},
			{typ: "chararr", in: []string{"\"", "\"", "'"}, want: strptr(`{"\"","\"",'}`)},
			{typ: "chararr", in: []string{"\t", " ", "\r"}, want: strptr("{\"\t\",\" \",\"\r\"}")},
			{typ: "chararr", in: []string{"\f", "\n"}, want: strptr("{\"\f\",\"\n\"}")},
			{typ: "chararr", in: []string{"\v", "\\"}, want: strptr("{\"\v\",\"\\\\\"}")},

			// TODO handle \b and \a
			//{typ: "chararr", in: []byte("\b\a"), want: strptr("{\"\b\",\"\a\"}")},
		},
	}}.execute(t)
}

func TestBPCharArray_Scanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			s := BPCharArrayToString{S: new(string)}
			return s, s.S
		},
		rows: []test_scanner_row{
			{typ: "bpchararr", in: nil, want: new(string)},
			{typ: "bpchararr", in: `{}`, want: strptr(``)},
			{typ: "bpchararr", in: `{a,b,c}`, want: strptr(`abc`)},
			{typ: "bpchararr", in: `{a,b,'}`, want: strptr(`ab'`)},
			{typ: "bpchararr", in: `{a,",",'}`, want: strptr(`a,'`)},
			{typ: "bpchararr", in: `{"\"",",",'}`, want: strptr(`",'`)},
			{typ: "bpchararr", in: `{魔,",",'}`, want: strptr(`魔,'`)},
			{typ: "bpchararr", in: `{"\"","\"",'}`, want: strptr(`""'`)},
			{typ: "bpchararr", in: "{\"\t\",\" \",\"\r\"}", want: strptr("\t \r")},
			{typ: "bpchararr", in: "{\"\b\",\"\f\",\"\n\"}", want: strptr("\b\f\n")},
			{typ: "bpchararr", in: "{\"\a\",\"\v\",\"\\\\\"}", want: strptr("\a\v\\")},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			s := BPCharArrayToByteSlice{S: new([]byte)}
			return s, s.S
		},
		rows: []test_scanner_row{
			{typ: "bpchararr", in: nil, want: new([]byte)},
			{typ: "bpchararr", in: `{}`, want: &[]byte{}},
			{typ: "bpchararr", in: `{a,b,c}`, want: &[]byte{'a', 'b', 'c'}},
			{typ: "bpchararr", in: `{a,b,'}`, want: &[]byte{'a', 'b', '\''}},
			{typ: "bpchararr", in: `{a,",",'}`, want: &[]byte{'a', ',', '\''}},
			{typ: "bpchararr", in: `{"\"",",",'}`, want: &[]byte{'"', ',', '\''}},
			{typ: "bpchararr", in: `{魔,",",'}`, want: &[]byte{233, 173, 148, ',', '\''}},
			{typ: "bpchararr", in: `{"\"","\"",'}`, want: &[]byte{'"', '"', '\''}},
			{typ: "bpchararr", in: "{\"\t\",\" \",\"\r\"}", want: &[]byte{'\t', ' ', '\r'}},
			{typ: "bpchararr", in: "{\"\b\",\"\f\",\"\n\"}", want: &[]byte{'\b', '\f', '\n'}},
			{typ: "bpchararr", in: "{\"\a\",\"\v\",\"\\\\\"}", want: &[]byte{'\a', '\v', '\\'}},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			s := BPCharArrayToRuneSlice{S: new([]rune)}
			return s, s.S
		},
		rows: []test_scanner_row{
			{typ: "bpchararr", in: nil, want: new([]rune)},
			{typ: "bpchararr", in: `{}`, want: &[]rune{}},
			{typ: "bpchararr", in: `{a,b,c}`, want: &[]rune{'a', 'b', 'c'}},
			{typ: "bpchararr", in: `{a,b,'}`, want: &[]rune{'a', 'b', '\''}},
			{typ: "bpchararr", in: `{a,",",'}`, want: &[]rune{'a', ',', '\''}},
			{typ: "bpchararr", in: `{"\"",",",'}`, want: &[]rune{'"', ',', '\''}},
			{typ: "bpchararr", in: `{魔,",",'}`, want: &[]rune{'魔', ',', '\''}},
			{typ: "bpchararr", in: `{"\"","\"",'}`, want: &[]rune{'"', '"', '\''}},
			{typ: "bpchararr", in: "{\"\t\",\" \",\"\r\"}", want: &[]rune{'\t', ' ', '\r'}},
			{typ: "bpchararr", in: "{\"\b\",\"\f\",\"\n\"}", want: &[]rune{'\b', '\f', '\n'}},
			{typ: "bpchararr", in: "{\"\a\",\"\v\",\"\\\\\"}", want: &[]rune{'\a', '\v', '\\'}},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			s := BPCharArrayToStringSlice{S: new([]string)}
			return s, s.S
		},
		rows: []test_scanner_row{
			{typ: "chararr", in: nil, want: new([]string)},
			{typ: "chararr", in: `{}`, want: &[]string{}},
			{typ: "chararr", in: `{a,b,c}`, want: &[]string{"a", "b", "c"}},
			{typ: "chararr", in: `{a,b,'}`, want: &[]string{"a", "b", "'"}},
			{typ: "chararr", in: `{a,",",'}`, want: &[]string{"a", ",", "'"}},
			{typ: "chararr", in: `{"\"",",",'}`, want: &[]string{"\"", ",", "'"}},
			{typ: "chararr", in: `{魔,",",'}`, want: &[]string{"魔", ",", "'"}},
			{typ: "chararr", in: `{"\"","\"",'}`, want: &[]string{"\"", "\"", "'"}},
			{typ: "chararr", in: "{\"\t\",\" \",\"\r\"}", want: &[]string{"\t", " ", "\r"}},
			{typ: "chararr", in: "{\"\b\",\"\f\",\"\n\"}", want: &[]string{"\b", "\f", "\n"}},
			{typ: "chararr", in: "{\"\a\",\"\v\",\"\\\\\"}", want: &[]string{"\a", "\v", "\\"}},
		},
	}}.execute(t)
}
