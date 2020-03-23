package convert

import (
	"testing"
)

func TestBPCharArrScanners(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			s := BPCharArr2ByteSlice{Ptr: new([]byte)}
			return s, s.Ptr
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
			s := BPCharArr2RuneSlice{Ptr: new([]rune)}
			return s, s.Ptr
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
			s := BPCharArr2String{Ptr: new(string)}
			return s, s.Ptr
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
	}}.execute(t)
}
