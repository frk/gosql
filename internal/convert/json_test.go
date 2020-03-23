package convert

import (
	"testing"
)

func TestJSON_Scanner(t *testing.T) {
	type data struct {
		Foo []interface{} `json:"foo"`
	}

	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			s := JSON{V: new(data)}
			return s, s.V
		},
		rows: []test_scanner_row{
			{typ: "json", in: nil, want: new(data)},
			{typ: "json", in: `{}`, want: new(data)},
			{typ: "json", in: `{"foo":["bar", "baz", 123]}`, want: &data{[]interface{}{"bar", "baz", float64(123)}}},
			{typ: "json", in: `{"foo":[123, "baz", "bar"]}`, want: &data{[]interface{}{float64(123), "baz", "bar"}}},
		},
	}}.execute(t)
}

func TestJSON_Valuer(t *testing.T) {
	type data struct {
		Foo []interface{} `json:"foo"`
		Bar string        `json:"bar"`
	}

	test_valuer{{
		valuer: func() interface{} {
			return new(JSON)
		},
		rows: []test_valuer_row{
			{typ: "json", in: nil, want: nil},
			{typ: "json", in: data{[]interface{}{1, 8}, "abcdef"}, want: strptr(`{"foo":[1,8],"bar":"abcdef"}`)},
		},
	}}.execute(t)
}

func TestJSON_NoScanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			return nil, new([]byte)
		},
		rows: []test_scanner_row{
			{typ: "json", in: `{"foo":[123, "baz", "bar"]}`, want: bytesptr(`{"foo":[123, "baz", "bar"]}`)},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			return nil, new(string)
		},
		rows: []test_scanner_row{
			{typ: "json", in: `{"foo":[123, "baz", "bar"]}`, want: strptr(`{"foo":[123, "baz", "bar"]}`)},
		},
	}}.execute(t)
}

func TestJSON_NoValuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return nil
		},
		rows: []test_valuer_row{
			{typ: "json", in: nil, want: nil},
			{typ: "json", in: `{"foo":[1,8],"bar":"abcdef"}`, want: strptr(`{"foo":[1,8],"bar":"abcdef"}`)},
			{typ: "json", in: []byte(`{"foo":[1,8],"bar":"abcdef"}`), want: strptr(`{"foo":[1,8],"bar":"abcdef"}`)},
		},
	}}.execute(t)
}
