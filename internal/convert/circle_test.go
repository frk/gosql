package convert

import (
	"testing"
)

func TestCircle_NoValuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return nil // string
		},
		rows: []test_valuer_row{
			{typ: "circle", in: nil, want: nil},
			{typ: "circle", in: "<(0,0),3.5>", want: strptr("<(0,0),3.5>")},
			{typ: "circle", in: "<(0.5,1),5>", want: strptr("<(0.5,1),5>")},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		rows: []test_valuer_row{
			{typ: "circle", in: nil, want: nil},
			{typ: "circle", in: []byte("<(0,0),3.5>"), want: strptr("<(0,0),3.5>")},
			{typ: "circle", in: []byte("<(0.5,1),5>"), want: strptr("<(0.5,1),5>")},
		},
	}}.execute(t)
}

func TestCircle_NoScanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			d := new(string)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "circle", in: "<(0,0),3.5>", want: strptr("<(0,0),3.5>")},
			{typ: "circle", in: "<(0.5,1),5>", want: strptr("<(0.5,1),5>")},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			d := new([]byte)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "circle", in: "<(0,0),3.5>", want: bytesptr("<(0,0),3.5>")},
			{typ: "circle", in: "<(0.5,1),5>", want: bytesptr("<(0.5,1),5>")},
		},
	}}.execute(t)
}
