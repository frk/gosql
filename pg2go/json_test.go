package pg2go

import (
	"testing"
)

func TestJSON(t *testing.T) {
	type data struct {
		Foo []interface{} `json:"foo"`
	}

	testlist2{{
		valuer:  JSON,
		scanner: JSON,
		data: []testdata{
			{
				input:  data{},
				output: data{}},
			{
				input:  data{[]interface{}{"bar", "baz", float64(123)}},
				output: data{[]interface{}{"bar", "baz", float64(123)}}},
			{
				input:  data{[]interface{}{float64(123), "baz", "bar"}},
				output: data{[]interface{}{float64(123), "baz", "bar"}}},
		},
	}, {
		data: []testdata{
			{
				input:  string(`{}`),
				output: string(`{}`)},
			{
				input:  string(`{"foo":["bar", "baz", 123]}`),
				output: string(`{"foo":["bar", "baz", 123]}`)},
			{
				input:  string(`{"foo":[123, "baz", "bar"]}`),
				output: string(`{"foo":[123, "baz", "bar"]}`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte(`{}`),
				output: []byte(`{}`)},
			{
				input:  []byte(`{"foo":["bar", "baz", 123]}`),
				output: []byte(`{"foo":["bar", "baz", 123]}`)},
			{
				input:  []byte(`{"foo":[123, "baz", "bar"]}`),
				output: []byte(`{"foo":[123, "baz", "bar"]}`)},
		},
	}}.execute(t, "json")
}

func TestJSONB(t *testing.T) {
	type data struct {
		Foo []interface{} `json:"foo"`
	}

	testlist2{{
		valuer:  JSON,
		scanner: JSON,
		data: []testdata{
			{
				input:  data{},
				output: data{}},
			{
				input:  data{[]interface{}{"bar", "baz", float64(123)}},
				output: data{[]interface{}{"bar", "baz", float64(123)}}},
			{
				input:  data{[]interface{}{float64(123), "baz", "bar"}},
				output: data{[]interface{}{float64(123), "baz", "bar"}}},
		},
	}, {
		data: []testdata{
			{
				input:  string(`{}`),
				output: string(`{}`)},
			{
				input:  string(`{"foo": ["bar", "baz", 123]}`),
				output: string(`{"foo": ["bar", "baz", 123]}`)},
			{
				input:  string(`{"foo": [123, "baz", "bar"]}`),
				output: string(`{"foo": [123, "baz", "bar"]}`)},
		},
	}, {
		data: []testdata{
			{
				input:  []byte(`{}`),
				output: []byte(`{}`)},
			{
				input:  []byte(`{"foo": ["bar", "baz", 123]}`),
				output: []byte(`{"foo": ["bar", "baz", 123]}`)},
			{
				input:  []byte(`{"foo": [123, "baz", "bar"]}`),
				output: []byte(`{"foo": [123, "baz", "bar"]}`)},
		},
	}}.execute(t, "jsonb")
}
