package pgsql

import (
	"testing"
)

func TestTSQueryArray(t *testing.T) {
	B := func(s string) []byte { return []byte(s) }

	testlist2{{
		valuer:  TSQueryArrayFromStringSlice,
		scanner: TSQueryArrayToStringSlice,
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
		valuer:  TSQueryArrayFromByteSliceSlice,
		scanner: TSQueryArrayToByteSliceSlice,
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
