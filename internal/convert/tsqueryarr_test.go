package convert

import (
	"testing"
)

func TestTSQueryArray(t *testing.T) {
	B := func(s string) []byte { return []byte(s) }

	testlist{{
		valuer: func() interface{} {
			return new(TSQueryArrayFromStringSlice)
		},
		scanner: func() (interface{}, interface{}) {
			v := TSQueryArrayToStringSlice{Val: new([]string)}
			return v, v.Val
		},
		data: []testdata{
			{input: []string(nil), output: []string(nil)},
			{input: []string{}, output: []string{}},
			{
				input:  []string{`'rat'`},
				output: []string{`'rat'`}},
			{
				input:  []string{`'fat' & 'rat'`, `'bat':AB`},
				output: []string{`'fat' & 'rat'`, `'bat':AB`}},
			{
				input:  []string{`!'fat':* & !'rat'`, `'bat' | 'bad'`},
				output: []string{`!'fat':* & !'rat'`, `'bat' | 'bad'`}},
		},
	}, {
		valuer: func() interface{} {
			return new(TSQueryArrayFromByteSliceSlice)
		},
		scanner: func() (interface{}, interface{}) {
			v := TSQueryArrayToByteSliceSlice{Val: new([][]byte)}
			return v, v.Val
		},
		data: []testdata{
			{input: [][]byte(nil), output: [][]byte(nil)},
			{input: [][]byte{}, output: [][]byte{}},
			{
				input:  [][]byte{B(`'rat'`)},
				output: [][]byte{B(`'rat'`)}},
			{
				input:  [][]byte{B(`'fat' & 'rat'`), B(`'bat':AB`)},
				output: [][]byte{B(`'fat' & 'rat'`), B(`'bat':AB`)}},
			{
				input:  [][]byte{B(`!'fat':* & !'rat'`), B(`'bat' | 'bad'`)},
				output: [][]byte{B(`!'fat':* & !'rat'`), B(`'bat' | 'bad'`)}},
		},
	}, {
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			v := new(string)
			return v, v
		},
		data: []testdata{
			{input: string(`{}`), output: string(`{}`)},
			{input: string(`{'rat'}`), output: string(`{'rat'}`)},
			{
				input:  string(`{"'fat' & 'rat'",'bat':AB}`),
				output: string(`{"'fat' & 'rat'",'bat':AB}`)},
			{
				input:  string(`{"!'fat':* & !'rat'","'bat' | 'bad'"}`),
				output: string(`{"!'fat':* & !'rat'","'bat' | 'bad'"}`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			v := new([]byte)
			return v, v
		},
		data: []testdata{
			{input: []byte(`{}`), output: []byte(`{}`)},
			{input: []byte(`{'rat'}`), output: []byte(`{'rat'}`)},
			{
				input:  []byte(`{"'fat' & 'rat'",'bat':AB}`),
				output: []byte(`{"'fat' & 'rat'",'bat':AB}`)},
			{
				input:  []byte(`{"!'fat':* & !'rat'","'bat' | 'bad'"}`),
				output: []byte(`{"!'fat':* & !'rat'","'bat' | 'bad'"}`)},
		},
	}}.execute(t, "tsqueryarr")
}
