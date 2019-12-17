package convert

import (
	"database/sql"
	"testing"
)

func TestCharArrScanners(t *testing.T) {
	test_table{{
		scnr: func() (sql.Scanner, interface{}) {
			s := CharArr2ByteSlice{Ptr: new([]byte)}
			return s, s.Ptr
		},
		rows: []testrow{
			{typ: "chararr", in: nil, want: new([]byte)},
			{typ: "chararr", in: `{}`, want: &[]byte{}},
			{typ: "chararr", in: `{a,b,c}`, want: &[]byte{'a', 'b', 'c'}},
			{typ: "chararr", in: `{a,b,'}`, want: &[]byte{'a', 'b', '\''}},
			{typ: "chararr", in: `{a,",",'}`, want: &[]byte{'a', ',', '\''}},
			{typ: "chararr", in: `{"\"",",",'}`, want: &[]byte{'"', ',', '\''}},
			{typ: "chararr", in: `{魔,",",'}`, want: &[]byte{233, 173, 148, ',', '\''}},
			{typ: "chararr", in: `{"\"","\"",'}`, want: &[]byte{'"', '"', '\''}},
			{typ: "chararr", in: "{\"\t\",\" \",\"\r\"}", want: &[]byte{'\t', ' ', '\r'}},
			{typ: "chararr", in: "{\"\b\",\"\f\",\"\n\"}", want: &[]byte{'\b', '\f', '\n'}},
			{typ: "chararr", in: "{\"\a\",\"\v\",\"\\\\\"}", want: &[]byte{'\a', '\v', '\\'}},
		},
	}, {
		scnr: func() (sql.Scanner, interface{}) {
			s := CharArr2RuneSlice{Ptr: new([]rune)}
			return s, s.Ptr
		},
		rows: []testrow{
			{typ: "chararr", in: nil, want: new([]rune)},
			{typ: "chararr", in: `{}`, want: &[]rune{}},
			{typ: "chararr", in: `{a,b,c}`, want: &[]rune{'a', 'b', 'c'}},
			{typ: "chararr", in: `{a,b,'}`, want: &[]rune{'a', 'b', '\''}},
			{typ: "chararr", in: `{a,",",'}`, want: &[]rune{'a', ',', '\''}},
			{typ: "chararr", in: `{"\"",",",'}`, want: &[]rune{'"', ',', '\''}},
			{typ: "chararr", in: `{魔,",",'}`, want: &[]rune{'魔', ',', '\''}},
			{typ: "chararr", in: `{"\"","\"",'}`, want: &[]rune{'"', '"', '\''}},
			{typ: "chararr", in: "{\"\t\",\" \",\"\r\"}", want: &[]rune{'\t', ' ', '\r'}},
			{typ: "chararr", in: "{\"\b\",\"\f\",\"\n\"}", want: &[]rune{'\b', '\f', '\n'}},
			{typ: "chararr", in: "{\"\a\",\"\v\",\"\\\\\"}", want: &[]rune{'\a', '\v', '\\'}},
		},
	}, {
		scnr: func() (sql.Scanner, interface{}) {
			s := CharArr2String{Ptr: new(string)}
			return s, s.Ptr
		},
		rows: []testrow{
			{typ: "chararr", in: nil, want: new(string)},
			{typ: "chararr", in: `{}`, want: strptr(``)},
			{typ: "chararr", in: `{a,b,c}`, want: strptr(`abc`)},
			{typ: "chararr", in: `{a,b,'}`, want: strptr(`ab'`)},
			{typ: "chararr", in: `{a,",",'}`, want: strptr(`a,'`)},
			{typ: "chararr", in: `{"\"",",",'}`, want: strptr(`",'`)},
			{typ: "chararr", in: `{魔,",",'}`, want: strptr(`魔,'`)},
			{typ: "chararr", in: `{"\"","\"",'}`, want: strptr(`""'`)},
			{typ: "chararr", in: "{\"\t\",\" \",\"\r\"}", want: strptr("\t \r")},
			{typ: "chararr", in: "{\"\b\",\"\f\",\"\n\"}", want: strptr("\b\f\n")},
			{typ: "chararr", in: "{\"\a\",\"\v\",\"\\\\\"}", want: strptr("\a\v\\")},
		},
	}}.execute(t)
}
