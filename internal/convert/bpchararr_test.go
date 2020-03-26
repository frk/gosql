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
			{typ: "bpchararr", in: "", want: strptr(`{}`)},
			{typ: "bpchararr", in: `abc`, want: strptr(`{a,b,c}`)},
			{typ: "bpchararr", in: `ab'`, want: strptr(`{a,b,'}`)},
			{typ: "bpchararr", in: `a,'`, want: strptr(`{a,",",'}`)},
			{typ: "bpchararr", in: `",'`, want: strptr(`{"\"",",",'}`)},
			{typ: "bpchararr", in: `魔,'`, want: strptr(`{魔,",",'}`)},
			{typ: "bpchararr", in: `""'`, want: strptr(`{"\"","\"",'}`)},
			{typ: "bpchararr", in: "\t \r", want: strptr("{\"\t\",\" \",\"\r\"}")},
			{typ: "bpchararr", in: "\f\n", want: strptr("{\"\f\",\"\n\"}")},
			{typ: "bpchararr", in: "\v\\", want: strptr("{\"\v\",\"\\\\\"}")},

			// TODO handle \b and \a
			// {typ: "bpchararr", in: "\b\a", want: strptr("{\"\b\",\"\a\"}")},
		},
	}, {
		valuer: func() interface{} {
			return new(BPCharArrayFromByteSlice)
		},
		rows: []test_valuer_row{
			{typ: "bpchararr", in: nil, want: nil},
			{typ: "bpchararr", in: []byte{}, want: strptr(`{}`)},
			{typ: "bpchararr", in: []byte(`abc`), want: strptr(`{a,b,c}`)},
			{typ: "bpchararr", in: []byte(`ab'`), want: strptr(`{a,b,'}`)},
			{typ: "bpchararr", in: []byte(`a,'`), want: strptr(`{a,",",'}`)},
			{typ: "bpchararr", in: []byte(`",'`), want: strptr(`{"\"",",",'}`)},
			{typ: "bpchararr", in: []byte(`""'`), want: strptr(`{"\"","\"",'}`)},
			{typ: "bpchararr", in: []byte("\t \r"), want: strptr("{\"\t\",\" \",\"\r\"}")},
			{typ: "bpchararr", in: []byte("\f\n"), want: strptr("{\"\f\",\"\n\"}")},
			{typ: "bpchararr", in: []byte("\v\\"), want: strptr("{\"\v\",\"\\\\\"}")},

			// TODO handle \b and \a
			//{typ: "bpchararr", in: []byte("\b\a"), want: strptr("{\"\b\",\"\a\"}")},
		},
	}, {
		valuer: func() interface{} {
			return new(BPCharArrayFromRuneSlice)
		},
		rows: []test_valuer_row{
			{typ: "bpchararr", in: nil, want: nil},
			{typ: "bpchararr", in: []rune{}, want: strptr(`{}`)},
			{typ: "bpchararr", in: []rune(`abc`), want: strptr(`{a,b,c}`)},
			{typ: "bpchararr", in: []rune(`ab'`), want: strptr(`{a,b,'}`)},
			{typ: "bpchararr", in: []rune(`a,'`), want: strptr(`{a,",",'}`)},
			{typ: "bpchararr", in: []rune(`",'`), want: strptr(`{"\"",",",'}`)},
			{typ: "bpchararr", in: []rune(`""'`), want: strptr(`{"\"","\"",'}`)},
			{typ: "bpchararr", in: []rune("\t \r"), want: strptr("{\"\t\",\" \",\"\r\"}")},
			{typ: "bpchararr", in: []rune("\f\n"), want: strptr("{\"\f\",\"\n\"}")},
			{typ: "bpchararr", in: []rune("\v\\"), want: strptr("{\"\v\",\"\\\\\"}")},

			// TODO handle \b and \a
			//{typ: "bpchararr", in: []byte("\b\a"), want: strptr("{\"\b\",\"\a\"}")},
		},
	}, {
		valuer: func() interface{} {
			return new(BPCharArrayFromStringSlice)
		},
		rows: []test_valuer_row{
			{typ: "bpchararr", in: nil, want: nil},
			{typ: "bpchararr", in: []string{}, want: strptr(`{}`)},
			{typ: "bpchararr", in: []string{"a", "b", "c"}, want: strptr(`{a,b,c}`)},
			{typ: "bpchararr", in: []string{"a", "b", "'"}, want: strptr(`{a,b,'}`)},
			{typ: "bpchararr", in: []string{"a", ",", "'"}, want: strptr(`{a,",",'}`)},
			{typ: "bpchararr", in: []string{"\"", ",", "'"}, want: strptr(`{"\"",",",'}`)},
			{typ: "bpchararr", in: []string{"\"", "\"", "'"}, want: strptr(`{"\"","\"",'}`)},
			{typ: "bpchararr", in: []string{"\t", " ", "\r"}, want: strptr("{\"\t\",\" \",\"\r\"}")},
			{typ: "bpchararr", in: []string{"\f", "\n"}, want: strptr("{\"\f\",\"\n\"}")},
			{typ: "bpchararr", in: []string{"\v", "\\"}, want: strptr("{\"\v\",\"\\\\\"}")},

			// TODO handle \b and \a
			//{typ: "bpchararr", in: []byte("\b\a"), want: strptr("{\"\b\",\"\a\"}")},
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
			{typ: "bpchararr", in: nil, want: new([]string)},
			{typ: "bpchararr", in: `{}`, want: &[]string{}},
			{typ: "bpchararr", in: `{a,b,c}`, want: &[]string{"a", "b", "c"}},
			{typ: "bpchararr", in: `{a,b,'}`, want: &[]string{"a", "b", "'"}},
			{typ: "bpchararr", in: `{a,",",'}`, want: &[]string{"a", ",", "'"}},
			{typ: "bpchararr", in: `{"\"",",",'}`, want: &[]string{"\"", ",", "'"}},
			{typ: "bpchararr", in: `{魔,",",'}`, want: &[]string{"魔", ",", "'"}},
			{typ: "bpchararr", in: `{"\"","\"",'}`, want: &[]string{"\"", "\"", "'"}},
			{typ: "bpchararr", in: "{\"\t\",\" \",\"\r\"}", want: &[]string{"\t", " ", "\r"}},
			{typ: "bpchararr", in: "{\"\b\",\"\f\",\"\n\"}", want: &[]string{"\b", "\f", "\n"}},
			{typ: "bpchararr", in: "{\"\a\",\"\v\",\"\\\\\"}", want: &[]string{"\a", "\v", "\\"}},
		},
	}}.execute(t)
}
