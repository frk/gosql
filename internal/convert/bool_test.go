package convert

import (
	"testing"
)

func TestBool_NoValuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return nil
		},
		rows: []test_valuer_row{
			{typ: "bool", in: true, want: strptr(`true`)},
			{typ: "bool", in: false, want: strptr(`false`)},
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
	}}.execute(t)
}
