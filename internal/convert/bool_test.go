package convert

import (
	"testing"
)

func TestBool_NoValuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return nil // bool
		},
		rows: []test_valuer_row{
			{typ: "bool", in: nil, want: nil},
			{typ: "bool", in: true, want: strptr(`true`)},
			{typ: "bool", in: false, want: strptr(`false`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // string
		},
		rows: []test_valuer_row{
			{typ: "bool", in: nil, want: nil},
			{typ: "bool", in: "true", want: strptr(`true`)},
			{typ: "bool", in: "false", want: strptr(`false`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		rows: []test_valuer_row{
			{typ: "bool", in: nil, want: nil},
			{typ: "bool", in: []byte("true"), want: strptr(`true`)},
			{typ: "bool", in: []byte("false"), want: strptr(`false`)},
		},
	}}.execute(t)
}

func TestBool_NoScanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			d := new(bool)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "bool", in: `true`, want: boolptr(true)},
			{typ: "bool", in: `false`, want: boolptr(false)},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			d := new(string)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "bool", in: `true`, want: strptr("true")},
			{typ: "bool", in: `false`, want: strptr("false")},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			d := new([]byte)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "bool", in: `true`, want: bytesptr("true")},
			{typ: "bool", in: `false`, want: bytesptr("false")},
		},
	}}.execute(t)
}
