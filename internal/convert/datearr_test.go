package convert

import (
	"testing"
	"time"
)

func TestDateArray_Valuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return new(DateArrayFromTimeSlice)
		},
		rows: []test_valuer_row{
			{typ: "datearr", in: nil, want: nil},
			{typ: "datearr", in: []time.Time{}, want: strptr(`{}`)},
			{
				typ:  "datearr",
				in:   []time.Time{dateval(1999, 1, 8)},
				want: strptr(`{1999-01-08}`)},
			{
				typ:  "datearr",
				in:   []time.Time{dateval(1999, 1, 8), dateval(2001, 5, 5)},
				want: strptr(`{1999-01-08,2001-05-05}`)},
			{
				typ:  "datearr",
				in:   []time.Time{dateval(2020, 3, 28), dateval(2001, 5, 5)},
				want: strptr(`{2020-03-28,2001-05-05}`)},
		},
	}}.execute(t)
}

func TestDateArray_Scanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			d := &DateArrayToTimeSlice{V: new([]time.Time)}
			return d, d.V
		},
		rows: []test_scanner_row{
			{typ: "datearr", in: nil, want: new([]time.Time)},
			{typ: "datearr", in: `{}`, want: &[]time.Time{}},
			{
				typ:  "datearr",
				in:   `{1999-01-08}`,
				want: &[]time.Time{dateval(1999, 1, 8)}},
			{
				typ:  "datearr",
				in:   `{1999-01-08,2001-05-05}`,
				want: &[]time.Time{dateval(1999, 1, 8), dateval(2001, 5, 5)}},
			{
				typ:  "datearr",
				in:   `{2020-03-28,2001-05-05}`,
				want: &[]time.Time{dateval(2020, 3, 28), dateval(2001, 5, 5)}},
		},
	}}.execute(t)
}

func TestDateArray_NoValuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return nil // string
		},
		rows: []test_valuer_row{
			{typ: "datearr", in: nil, want: nil},
			{typ: "datearr", in: `{}`, want: strptr(`{}`)},
			{typ: "datearr", in: `{1999-01-08}`, want: strptr(`{1999-01-08}`)},
			{typ: "datearr", in: `{1999-01-08,2001-05-05}`, want: strptr(`{1999-01-08,2001-05-05}`)},
			{typ: "datearr", in: `{2020-03-28,2001-05-05}`, want: strptr(`{2020-03-28,2001-05-05}`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		rows: []test_valuer_row{
			{typ: "datearr", in: nil, want: nil},
			{typ: "datearr", in: []byte(`{1999-01-08}`), want: strptr(`{1999-01-08}`)},
			{typ: "datearr", in: []byte(`{1999-01-08,2001-05-05}`), want: strptr(`{1999-01-08,2001-05-05}`)},
			{typ: "datearr", in: []byte(`{2020-03-28,2001-05-05}`), want: strptr(`{2020-03-28,2001-05-05}`)},
		},
	}}.execute(t)
}

func TestDateArray_NoScanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			d := new(string)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "datearr", in: `{}`, want: strptr(`{}`)},
			{typ: "datearr", in: `{1999-01-08}`, want: strptr(`{1999-01-08}`)},
			{typ: "datearr", in: `{1999-01-08,2001-05-05}`, want: strptr(`{1999-01-08,2001-05-05}`)},
			{typ: "datearr", in: `{2020-03-28,2001-05-05}`, want: strptr(`{2020-03-28,2001-05-05}`)},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			d := new([]byte)
			return d, d
		},
		rows: []test_scanner_row{
			{typ: "datearr", in: `{}`, want: bytesptr(`{}`)},
			{typ: "datearr", in: `{1999-01-08}`, want: bytesptr(`{1999-01-08}`)},
			{typ: "datearr", in: `{1999-01-08,2001-05-05}`, want: bytesptr(`{1999-01-08,2001-05-05}`)},
			{typ: "datearr", in: `{2020-03-28,2001-05-05}`, want: bytesptr(`{2020-03-28,2001-05-05}`)},
		},
	}}.execute(t)
}
