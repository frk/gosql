package convert

import (
	"testing"
)

func TestBoolArray_Valuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return new(BoolArrayFromBoolSlice)
		},
		rows: []test_valuer_row{
			{typ: "boolarr", in: nil, want: nil},
			{typ: "boolarr", in: []bool{}, want: strptr(`{}`)},
			{typ: "boolarr", in: []bool{true}, want: strptr(`{t}`)},
			{typ: "boolarr", in: []bool{false}, want: strptr(`{f}`)},
			{
				typ:  "boolarr",
				in:   []bool{false, false, false, true, true, false, true, true},
				want: strptr(`{f,f,f,t,t,f,t,t}`)},
		},
	}}.execute(t)
}

func TestBoolArray_Scanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			s := BoolArrayToBoolSlice{S: new([]bool)}
			return s, s.S
		},
		rows: []test_scanner_row{
			{typ: "boolarr", in: nil, want: new([]bool)},
			{typ: "boolarr", in: `{}`, want: &[]bool{}},
			{typ: "boolarr", in: `{f}`, want: &[]bool{false}},
			{typ: "boolarr", in: `{t}`, want: &[]bool{true}},
			{typ: "boolarr", in: `{t,f}`, want: &[]bool{true, false}},
			{
				typ:  "boolarr",
				in:   `{f,t,t,t,f,t,f,f}`,
				want: &[]bool{false, true, true, true, false, true, false, false}},
		},
	}}.execute(t)
}

func TestBoolArray_NoValuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return nil // string
		},
		rows: []test_valuer_row{
			{typ: "boolarr", in: nil, want: nil},
			{typ: "boolarr", in: "{}", want: strptr(`{}`)},
			{typ: "boolarr", in: "{true,false}", want: strptr(`{t,f}`)},
			{typ: "boolarr", in: "{t,f,f,f,t,t}", want: strptr(`{t,f,f,f,t,t}`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		rows: []test_valuer_row{
			{typ: "boolarr", in: nil, want: nil},
			{typ: "boolarr", in: []byte("{}"), want: strptr(`{}`)},
			{typ: "boolarr", in: []byte("{true,false}"), want: strptr(`{t,f}`)},
			{typ: "boolarr", in: []byte("{t,f,f,f,t,t}"), want: strptr(`{t,f,f,f,t,t}`)},
		},
	}}.execute(t)
}

func TestBoolArray_NoScanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			d := new(string)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "boolarr", in: `{}`, want: strptr("{}")},
			{typ: "boolarr", in: `{true,false}`, want: strptr("{t,f}")},
			{typ: "boolarr", in: `{t,f,f,f,t,t}`, want: strptr("{t,f,f,f,t,t}")},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			d := new([]byte)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "boolarr", in: `{}`, want: bytesptr("{}")},
			{typ: "boolarr", in: `{true,false}`, want: bytesptr("{t,f}")},
			{typ: "boolarr", in: `{t,f,f,f,t,t}`, want: bytesptr("{t,f,f,f,t,t}")},
		},
	}}.execute(t)
}
