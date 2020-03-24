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
