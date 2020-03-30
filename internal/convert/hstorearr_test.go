package convert

import (
	"database/sql"
	"testing"
)

func TestHStoreArray_ValuerAndScanner(t *testing.T) {
	A := strptr

	test_valuer_with_scanner{{
		valuer: func() interface{} {
			return new(HStoreArrayFromStringMapSlice)
		},
		scanner: func() (interface{}, interface{}) {
			v := &HStoreArrayToStringMapSlice{Val: new([]map[string]string)}
			return v, v.Val
		},
		rows: []test_valuer_with_scanner_row{
			{typ: "hstorearr", input: nil, output: new([]map[string]string)},
			{typ: "hstorearr", input: []map[string]string{}, output: &[]map[string]string{}},
			{
				typ:    "hstorearr",
				input:  []map[string]string{{"a": "1"}},
				output: &[]map[string]string{{"a": "1"}}},
			{
				typ:    "hstorearr",
				input:  []map[string]string{{"a": "1", "b": "2", "c": "3"}},
				output: &[]map[string]string{{"a": "1", "b": "2", "c": "3"}}},
			{
				typ:    "hstorearr",
				input:  []map[string]string{{"a": "1"}, {"b": "2"}, {"c": "3"}},
				output: &[]map[string]string{{"a": "1"}, {"b": "2"}, {"c": "3"}}},
			{
				typ:    "hstorearr",
				input:  []map[string]string{{"a": "1"}, nil, {"c": ""}},
				output: &[]map[string]string{{"a": "1"}, nil, {"c": ""}}},
			{
				typ:    "hstorearr",
				input:  []map[string]string{{"a": `'`}},
				output: &[]map[string]string{{"a": `'`}}},
			{
				typ:    "hstorearr",
				input:  []map[string]string{{"a": `foo" bar`}, {"b": `\`}},
				output: &[]map[string]string{{"a": `foo" bar`}, {"b": `\`}}},
			{
				typ:    "hstorearr",
				input:  []map[string]string{{"a": `foo\" bar`}, {"b": `\\`}},
				output: &[]map[string]string{{"a": `foo\" bar`}, {"b": `\\`}}},
			{
				typ:    "hstorearr",
				input:  []map[string]string{{"a": "\t"}},
				output: &[]map[string]string{{"a": "\t"}}},
			{
				typ:    "hstorearr",
				input:  []map[string]string{{"a": "\a"}, {"b": "\b"}},
				output: &[]map[string]string{{"a": "\a"}, {"b": "\b"}}},
		},
	}, {
		valuer: func() interface{} {
			return new(HStoreArrayFromStringPtrMapSlice)
		},
		scanner: func() (interface{}, interface{}) {
			v := &HStoreArrayToStringPtrMapSlice{Val: new([]map[string]*string)}
			return v, v.Val
		},
		rows: []test_valuer_with_scanner_row{
			{typ: "hstorearr", input: nil, output: new([]map[string]*string)},
			{typ: "hstorearr", input: []map[string]*string{}, output: &[]map[string]*string{}},
			{
				typ:    "hstorearr",
				input:  []map[string]*string{{"a": A("1")}},
				output: &[]map[string]*string{{"a": A("1")}}},
			{
				typ:    "hstorearr",
				input:  []map[string]*string{{"a": A("1"), "b": A("2"), "c": A("3")}},
				output: &[]map[string]*string{{"a": A("1"), "b": A("2"), "c": A("3")}}},
			{
				typ:    "hstorearr",
				input:  []map[string]*string{{"a": A("1")}, {"b": A("2")}, {"c": A("3")}},
				output: &[]map[string]*string{{"a": A("1")}, {"b": A("2")}, {"c": A("3")}}},
			{
				typ:    "hstorearr",
				input:  []map[string]*string{{"a": A("1")}, nil, {"c": nil}},
				output: &[]map[string]*string{{"a": A("1")}, nil, {"c": nil}}},
			{
				typ:    "hstorearr",
				input:  []map[string]*string{{"a": A(`'`)}},
				output: &[]map[string]*string{{"a": A(`'`)}}},
			{
				typ:    "hstorearr",
				input:  []map[string]*string{{"a": A(`foo" bar`)}, {"b": A(`\`)}},
				output: &[]map[string]*string{{"a": A(`foo" bar`)}, {"b": A(`\`)}}},
			{
				typ:    "hstorearr",
				input:  []map[string]*string{{"a": A(`foo\" bar`)}, {"b": A(`\\`)}},
				output: &[]map[string]*string{{"a": A(`foo\" bar`)}, {"b": A(`\\`)}}},
			{
				typ:    "hstorearr",
				input:  []map[string]*string{{"a": A("\t")}},
				output: &[]map[string]*string{{"a": A("\t")}}},
			{
				typ:    "hstorearr",
				input:  []map[string]*string{{"a": A("\a")}, {"b": A("\b")}},
				output: &[]map[string]*string{{"a": A("\a")}, {"b": A("\b")}}},
		},
	}, {
		valuer: func() interface{} {
			return new(HStoreArrayFromNullStringMapSlice)
		},
		scanner: func() (interface{}, interface{}) {
			v := &HStoreArrayToNullStringMapSlice{Val: new([]map[string]sql.NullString)}
			return v, v.Val
		},
		rows: []test_valuer_with_scanner_row{
			{typ: "hstorearr", input: nil, output: new([]map[string]sql.NullString)},
			{typ: "hstorearr", input: []map[string]sql.NullString{}, output: &[]map[string]sql.NullString{}},
			{
				typ:    "hstorearr",
				input:  []map[string]sql.NullString{{"a": {"1", true}}},
				output: &[]map[string]sql.NullString{{"a": {"1", true}}}},
			{
				typ:    "hstorearr",
				input:  []map[string]sql.NullString{{"a": {"1", true}, "b": {"2", true}, "c": {"3", true}}},
				output: &[]map[string]sql.NullString{{"a": {"1", true}, "b": {"2", true}, "c": {"3", true}}}},
			{
				typ:    "hstorearr",
				input:  []map[string]sql.NullString{{"a": {"1", true}}, {"b": {"2", true}}, {"c": {"3", true}}},
				output: &[]map[string]sql.NullString{{"a": {"1", true}}, {"b": {"2", true}}, {"c": {"3", true}}}},
			{
				typ:    "hstorearr",
				input:  []map[string]sql.NullString{{"a": {"1", true}}, nil, {"c": {"", false}}},
				output: &[]map[string]sql.NullString{{"a": {"1", true}}, nil, {"c": {"", false}}}},
			{
				typ:    "hstorearr",
				input:  []map[string]sql.NullString{{"a": {`'`, true}}},
				output: &[]map[string]sql.NullString{{"a": {`'`, true}}}},
			{
				typ:    "hstorearr",
				input:  []map[string]sql.NullString{{"a": {`foo" bar`, true}}, {"b": {`\`, true}}},
				output: &[]map[string]sql.NullString{{"a": {`foo" bar`, true}}, {"b": {`\`, true}}}},
			{
				typ:    "hstorearr",
				input:  []map[string]sql.NullString{{"a": {`foo\" bar`, true}}, {"b": {`\\`, true}}},
				output: &[]map[string]sql.NullString{{"a": {`foo\" bar`, true}}, {"b": {`\\`, true}}}},
			{
				typ:    "hstorearr",
				input:  []map[string]sql.NullString{{"a": {"\t", true}}},
				output: &[]map[string]sql.NullString{{"a": {"\t", true}}}},
			{
				typ:    "hstorearr",
				input:  []map[string]sql.NullString{{"a": {"\a", true}}, {"b": {"\b", true}}},
				output: &[]map[string]sql.NullString{{"a": {"\a", true}}, {"b": {"\b", true}}}},
		},
	}}.execute(t)
}

func TestHStoreArray_NoValuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return nil // string
		},
		rows: []test_valuer_row{
			{typ: "hstorearr", in: nil, want: nil},
			{typ: "hstorearr", in: "{}", want: strptr("{}")},
			{
				typ:  "hstorearr",
				in:   `{"\"a\"=>\"1\", \"b\"=>\"2\""}`,
				want: strptr(`{"\"a\"=>\"1\", \"b\"=>\"2\""}`)},
			{
				typ:  "hstorearr",
				in:   `{"\"text\"=>\"foo' \\\"bar\\\", baz \\\\quux\""}`,
				want: strptr(`{"\"text\"=>\"foo' \\\"bar\\\", baz \\\\quux\""}`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		rows: []test_valuer_row{
			{typ: "hstorearr", in: nil, want: nil},
			{typ: "hstorearr", in: "{}", want: strptr("{}")},
			{
				typ:  "hstorearr",
				in:   []byte(`{"\"a\"=>\"1\", \"b\"=>\"2\""}`),
				want: strptr(`{"\"a\"=>\"1\", \"b\"=>\"2\""}`)},
			{
				typ:  "hstorearr",
				in:   []byte(`{"\"text\"=>\"foo' \\\"bar\\\", baz \\\\quux\""}`),
				want: strptr(`{"\"text\"=>\"foo' \\\"bar\\\", baz \\\\quux\""}`)},
		},
	}}.execute(t)
}

func TestHStoreArray_NoScanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			v := new(string)
			return v, v
		},
		rows: []test_scanner_row{
			{typ: "hstorearr", in: `{}`, want: strptr(`{}`)},
			{
				typ:  "hstorearr",
				in:   `{"\"a\"=>\"1\""}`,
				want: strptr(`{"\"a\"=>\"1\""}`)},
			{
				typ:  "hstorearr",
				in:   `{"\"a\"=>\"1\"","\"b\"=>\"2\", \"c\"=>\"3\""}`,
				want: strptr(`{"\"a\"=>\"1\"","\"b\"=>\"2\", \"c\"=>\"3\""}`)},
			{
				typ:  "hstorearr",
				in:   `{"\"a\"=>\"1\", \"b\"=>NULL",NULL}`,
				want: strptr(`{"\"a\"=>\"1\", \"b\"=>NULL",NULL}`)},
			{
				typ:  "hstorearr",
				in:   `{"\"a\"=>\"1\", \"b\"=>\"\\\"\", \"c\"=>\"\\\\\\\\\""}`,
				want: strptr(`{"\"a\"=>\"1\", \"b\"=>\"\\\"\", \"c\"=>\"\\\\\\\\\""}`)},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			v := new([]byte)
			return v, v
		},
		rows: []test_scanner_row{
			{typ: "hstorearr", in: `{}`, want: bytesptr(`{}`)},
			{
				typ:  "hstorearr",
				in:   `{"\"a\"=>\"1\""}`,
				want: bytesptr(`{"\"a\"=>\"1\""}`)},
			{
				typ:  "hstorearr",
				in:   `{"\"a\"=>\"1\"","\"b\"=>\"2\", \"c\"=>\"3\""}`,
				want: bytesptr(`{"\"a\"=>\"1\"","\"b\"=>\"2\", \"c\"=>\"3\""}`)},
			{
				typ:  "hstorearr",
				in:   `{"\"a\"=>\"1\", \"b\"=>NULL",NULL}`,
				want: bytesptr(`{"\"a\"=>\"1\", \"b\"=>NULL",NULL}`)},
			{
				typ:  "hstorearr",
				in:   `{"\"a\"=>\"1\", \"b\"=>\"\\\"\", \"c\"=>\"\\\\\\\\\""}`,
				want: bytesptr(`{"\"a\"=>\"1\", \"b\"=>\"\\\"\", \"c\"=>\"\\\\\\\\\""}`)},
		},
	}}.execute(t)
}
