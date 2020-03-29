package convert

import (
	"testing"
	"time"
)

func TestDateRange_Valuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return new(DateRangeFromTimeArray2)
		},
		rows: []test_valuer_row{
			{typ: "daterange", in: nil, want: strptr(`(,)`)},
			{
				typ:  "daterange",
				in:   [2]time.Time{dateval(1999, 1, 8), dateval(2001, 5, 5)},
				want: strptr(`[1999-01-08,2001-05-05)`)},
			{
				typ:  "daterange",
				in:   [2]time.Time{{}, dateval(2001, 5, 5)},
				want: strptr(`(,2001-05-05)`)},
			{
				typ:  "daterange",
				in:   [2]time.Time{dateval(1999, 1, 8), {}},
				want: strptr(`[1999-01-08,)`)},
			{
				typ:  "daterange",
				in:   [2]time.Time{{}, {}},
				want: strptr(`(,)`)},
		},
	}}.execute(t)
}

func TestDateRange_Scanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			v := DateRangeToTimeArray2{V: new([2]time.Time)}
			return v, v.V
		},
		rows: []test_scanner_row{
			{typ: "daterange", in: nil, want: new([2]time.Time)},
			{
				typ:  "daterange",
				in:   `[1999-01-08,2001-05-05)`,
				want: &[2]time.Time{dateval(1999, 1, 8), dateval(2001, 5, 5)}},
			{
				typ:  "daterange",
				in:   `[1999-01-08,2001-05-05)`,
				want: &[2]time.Time{dateval(1999, 1, 8), dateval(2001, 5, 5)}},
			{
				typ:  "daterange",
				in:   `[1999-01-08,2001-05-05]`,
				want: &[2]time.Time{dateval(1999, 1, 8), dateval(2001, 5, 6)}},
			{
				typ:  "daterange",
				in:   `(1999-01-08,2001-05-05]`,
				want: &[2]time.Time{dateval(1999, 1, 9), dateval(2001, 5, 6)}},
			{
				typ:  "daterange",
				in:   "[,2001-05-05]",
				want: &[2]time.Time{{}, dateval(2001, 5, 6)}},
			{
				typ:  "daterange",
				in:   "[1999-01-08,]",
				want: &[2]time.Time{dateval(1999, 1, 8), {}}},
			{
				typ:  "daterange",
				in:   "[,]",
				want: &[2]time.Time{{}, {}}},
		},
	}}.execute(t)
}

func TestDateRange_NoValuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return nil // string
		},
		rows: []test_valuer_row{
			{typ: "daterange", in: nil, want: nil},
			{
				typ:  "daterange",
				in:   "[1999-01-08,2001-05-05)",
				want: strptr(`[1999-01-08,2001-05-05)`)},
			{
				typ:  "daterange",
				in:   "[1999-01-08,2001-05-05]",
				want: strptr(`[1999-01-08,2001-05-06)`)},
			{
				typ:  "daterange",
				in:   "(1999-01-08,2001-05-05]",
				want: strptr(`[1999-01-09,2001-05-06)`)},
			{
				typ:  "daterange",
				in:   "[,2001-05-05)",
				want: strptr(`(,2001-05-05)`)},
			{
				typ:  "daterange",
				in:   "[1999-01-08,]",
				want: strptr(`[1999-01-08,)`)},
			{
				typ:  "daterange",
				in:   "[,]",
				want: strptr(`(,)`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		rows: []test_valuer_row{
			{typ: "daterange", in: nil, want: nil},
			{
				typ:  "daterange",
				in:   []byte("[1999-01-08,2001-05-05)"),
				want: strptr(`[1999-01-08,2001-05-05)`)},
			{
				typ:  "daterange",
				in:   []byte("[1999-01-08,2001-05-05]"),
				want: strptr(`[1999-01-08,2001-05-06)`)},
			{
				typ:  "daterange",
				in:   []byte("(1999-01-08,2001-05-05]"),
				want: strptr(`[1999-01-09,2001-05-06)`)},
			{
				typ:  "daterange",
				in:   []byte("[,2001-05-05)"),
				want: strptr(`(,2001-05-05)`)},
			{
				typ:  "daterange",
				in:   []byte("[1999-01-08,]"),
				want: strptr(`[1999-01-08,)`)},
			{
				typ:  "daterange",
				in:   []byte("[,]"),
				want: strptr(`(,)`)},
		},
	}}.execute(t)
}

func TestDateRange_NoScanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			d := new(string)
			return d, d
		},
		rows: []test_scanner_row{
			{
				typ:  "daterange",
				in:   "[1999-01-08,2001-05-05)",
				want: strptr(`[1999-01-08,2001-05-05)`)},
			{
				typ:  "daterange",
				in:   "[1999-01-08,2001-05-05]",
				want: strptr(`[1999-01-08,2001-05-06)`)},
			{
				typ:  "daterange",
				in:   "(1999-01-08,2001-05-05]",
				want: strptr(`[1999-01-09,2001-05-06)`)},
			{
				typ:  "daterange",
				in:   "[,2001-05-05)",
				want: strptr(`(,2001-05-05)`)},
			{
				typ:  "daterange",
				in:   "[1999-01-08,]",
				want: strptr(`[1999-01-08,)`)},
			{
				typ:  "daterange",
				in:   "[,]",
				want: strptr(`(,)`)},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			d := new([]byte)
			return d, d
		},
		rows: []test_scanner_row{
			{
				typ:  "daterange",
				in:   "[1999-01-08,2001-05-05)",
				want: bytesptr(`[1999-01-08,2001-05-05)`)},
			{
				typ:  "daterange",
				in:   "[1999-01-08,2001-05-05]",
				want: bytesptr(`[1999-01-08,2001-05-06)`)},
			{
				typ:  "daterange",
				in:   "(1999-01-08,2001-05-05]",
				want: bytesptr(`[1999-01-09,2001-05-06)`)},
			{
				typ:  "daterange",
				in:   "[,2001-05-05)",
				want: bytesptr(`(,2001-05-05)`)},
			{
				typ:  "daterange",
				in:   "[1999-01-08,]",
				want: bytesptr(`[1999-01-08,)`)},
			{
				typ:  "daterange",
				in:   "[,]",
				want: bytesptr(`(,)`)},
		},
	}}.execute(t)
}
