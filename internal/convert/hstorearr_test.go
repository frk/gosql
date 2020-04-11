package convert

import (
	"database/sql"
	"testing"
)

func TestHStoreArray(t *testing.T) {
	A := strptr

	testlist{{
		valuer: func() interface{} {
			return new(HStoreArrayFromStringMapSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := HStoreArrayToStringMapSlice{Val: new([]map[string]string)}
			return s, s.Val
		},
		data: []testdata{
			{input: nil, output: []map[string]string(nil)},
			{input: []map[string]string{}, output: []map[string]string{}},
			{
				input:  []map[string]string{{"a": "1"}},
				output: []map[string]string{{"a": "1"}}},
			{
				input:  []map[string]string{{"a": "1", "b": "2", "c": "3"}},
				output: []map[string]string{{"a": "1", "b": "2", "c": "3"}}},
			{
				input:  []map[string]string{{"a": "1"}, {"b": "2"}, {"c": "3"}},
				output: []map[string]string{{"a": "1"}, {"b": "2"}, {"c": "3"}}},
			{
				input:  []map[string]string{{"a": "1"}, nil, {"c": ""}},
				output: []map[string]string{{"a": "1"}, nil, {"c": ""}}},
			{
				input:  []map[string]string{{"a": `'`}},
				output: []map[string]string{{"a": `'`}}},
			{
				input:  []map[string]string{{"a": `foo" bar`}, {"b": `\`}},
				output: []map[string]string{{"a": `foo" bar`}, {"b": `\`}}},
			{
				input:  []map[string]string{{"a": `foo\" bar`}, {"b": `\\`}},
				output: []map[string]string{{"a": `foo\" bar`}, {"b": `\\`}}},
			{
				input:  []map[string]string{{"a": "\t"}},
				output: []map[string]string{{"a": "\t"}}},
			{
				input:  []map[string]string{{"a": "\a"}, {"b": "\b"}},
				output: []map[string]string{{"a": "\a"}, {"b": "\b"}}},
		},
	}, {
		valuer: func() interface{} {
			return new(HStoreArrayFromStringPtrMapSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := HStoreArrayToStringPtrMapSlice{Val: new([]map[string]*string)}
			return s, s.Val
		},
		data: []testdata{
			{input: nil, output: []map[string]*string(nil)},
			{input: []map[string]*string{}, output: []map[string]*string{}},
			{
				input:  []map[string]*string{{"a": A("1")}},
				output: []map[string]*string{{"a": A("1")}}},
			{
				input:  []map[string]*string{{"a": A("1"), "b": A("2"), "c": A("3")}},
				output: []map[string]*string{{"a": A("1"), "b": A("2"), "c": A("3")}}},
			{
				input:  []map[string]*string{{"a": A("1")}, {"b": A("2")}, {"c": A("3")}},
				output: []map[string]*string{{"a": A("1")}, {"b": A("2")}, {"c": A("3")}}},
			{
				input:  []map[string]*string{{"a": A("1")}, nil, {"c": nil}},
				output: []map[string]*string{{"a": A("1")}, nil, {"c": nil}}},
			{
				input:  []map[string]*string{{"a": A(`'`)}},
				output: []map[string]*string{{"a": A(`'`)}}},
			{
				input:  []map[string]*string{{"a": A(`foo" bar`)}, {"b": A(`\`)}},
				output: []map[string]*string{{"a": A(`foo" bar`)}, {"b": A(`\`)}}},
			{
				input:  []map[string]*string{{"a": A(`foo\" bar`)}, {"b": A(`\\`)}},
				output: []map[string]*string{{"a": A(`foo\" bar`)}, {"b": A(`\\`)}}},
			{
				input:  []map[string]*string{{"a": A("\t")}},
				output: []map[string]*string{{"a": A("\t")}}},
			{
				input:  []map[string]*string{{"a": A("\a")}, {"b": A("\b")}},
				output: []map[string]*string{{"a": A("\a")}, {"b": A("\b")}}},
		},
	}, {
		valuer: func() interface{} {
			return new(HStoreArrayFromNullStringMapSlice)
		},
		scanner: func() (interface{}, interface{}) {
			s := HStoreArrayToNullStringMapSlice{Val: new([]map[string]sql.NullString)}
			return s, s.Val
		},
		data: []testdata{
			{input: nil, output: []map[string]sql.NullString(nil)},
			{input: []map[string]sql.NullString{}, output: []map[string]sql.NullString{}},
			{
				input:  []map[string]sql.NullString{{"a": {"1", true}}},
				output: []map[string]sql.NullString{{"a": {"1", true}}}},
			{
				input:  []map[string]sql.NullString{{"a": {"1", true}, "b": {"2", true}, "c": {"3", true}}},
				output: []map[string]sql.NullString{{"a": {"1", true}, "b": {"2", true}, "c": {"3", true}}}},
			{
				input:  []map[string]sql.NullString{{"a": {"1", true}}, {"b": {"2", true}}, {"c": {"3", true}}},
				output: []map[string]sql.NullString{{"a": {"1", true}}, {"b": {"2", true}}, {"c": {"3", true}}}},
			{
				input:  []map[string]sql.NullString{{"a": {"1", true}}, nil, {"c": {"", false}}},
				output: []map[string]sql.NullString{{"a": {"1", true}}, nil, {"c": {"", false}}}},
			{
				input:  []map[string]sql.NullString{{"a": {`'`, true}}},
				output: []map[string]sql.NullString{{"a": {`'`, true}}}},
			{
				input:  []map[string]sql.NullString{{"a": {`foo" bar`, true}}, {"b": {`\`, true}}},
				output: []map[string]sql.NullString{{"a": {`foo" bar`, true}}, {"b": {`\`, true}}}},
			{
				input:  []map[string]sql.NullString{{"a": {`foo\" bar`, true}}, {"b": {`\\`, true}}},
				output: []map[string]sql.NullString{{"a": {`foo\" bar`, true}}, {"b": {`\\`, true}}}},
			{
				input:  []map[string]sql.NullString{{"a": {"\t", true}}},
				output: []map[string]sql.NullString{{"a": {"\t", true}}}},
			{
				input:  []map[string]sql.NullString{{"a": {"\a", true}}, {"b": {"\b", true}}},
				output: []map[string]sql.NullString{{"a": {"\a", true}}, {"b": {"\b", true}}}},
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
			{input: string("{}"), output: string("{}")},
			{
				input:  string(`{"\"a\"=>\"1\", \"b\"=>\"2\""}`),
				output: string(`{"\"a\"=>\"1\", \"b\"=>\"2\""}`)},
			{
				input:  string(`{"\"text\"=>\"foo' \\\"bar\\\", baz \\\\quux\""}`),
				output: string(`{"\"text\"=>\"foo' \\\"bar\\\", baz \\\\quux\""}`)},
			{
				input:  string(`{"\"a\"=>\"1\", \"b\"=>NULL",NULL}`),
				output: string(`{"\"a\"=>\"1\", \"b\"=>NULL",NULL}`)},
			{
				input:  string(`{"\"a\"=>\"1\", \"b\"=>\"\\\"\", \"c\"=>\"\\\\\\\\\""}`),
				output: string(`{"\"a\"=>\"1\", \"b\"=>\"\\\"\", \"c\"=>\"\\\\\\\\\""}`)},
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
			{input: []byte("{}"), output: []byte("{}")},
			{
				input:  []byte(`{"\"a\"=>\"1\", \"b\"=>\"2\""}`),
				output: []byte(`{"\"a\"=>\"1\", \"b\"=>\"2\""}`)},
			{
				input:  []byte(`{"\"text\"=>\"foo' \\\"bar\\\", baz \\\\quux\""}`),
				output: []byte(`{"\"text\"=>\"foo' \\\"bar\\\", baz \\\\quux\""}`)},
			{
				input:  []byte(`{"\"a\"=>\"1\", \"b\"=>NULL",NULL}`),
				output: []byte(`{"\"a\"=>\"1\", \"b\"=>NULL",NULL}`)},
			{
				input:  []byte(`{"\"a\"=>\"1\", \"b\"=>\"\\\"\", \"c\"=>\"\\\\\\\\\""}`),
				output: []byte(`{"\"a\"=>\"1\", \"b\"=>\"\\\"\", \"c\"=>\"\\\\\\\\\""}`)},
		},
	}}.execute(t, "hstorearr")
}
