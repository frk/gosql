package convert

import (
	"testing"
	"time"
)

func TestDateRangeArray_Valuer(t *testing.T) {
	test_valuer{{
		valuer: func() interface{} {
			return new(DateRangeArrayFromTimeArray2Slice)
		},
		rows: []test_valuer_row{
			{typ: "daterangearr", in: nil, want: nil},
			{
				typ:  "daterangearr",
				in:   [][2]time.Time{{dateval(1999, 1, 8), dateval(2001, 5, 5)}},
				want: strptr(`{"[1999-01-08,2001-05-05)"}`)},
			{
				typ: "daterangearr",
				in: [][2]time.Time{
					{{}, dateval(2001, 5, 5)},
					{dateval(1999, 1, 8), {}},
					{dateval(1999, 1, 8), dateval(2001, 5, 5)},
				},
				want: strptr(`{"(,2001-05-05)","[1999-01-08,)","[1999-01-08,2001-05-05)"}`)},
			{
				typ:  "daterangearr",
				in:   [][2]time.Time{{}, {}, {}},
				want: strptr(`{"(,)","(,)","(,)"}`)},
		},
	}}.execute(t)
}

func TestDateRangeArray_Scanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			v := DateRangeArrayToTimeArray2Slice{V: new([][2]time.Time)}
			return v, v.V
		},
		rows: []test_scanner_row{
			{typ: "daterangearr", in: nil, want: new([][2]time.Time)},
			{
				typ:  "daterangearr",
				in:   `{"[1999-01-08,2001-05-05)"}`,
				want: &[][2]time.Time{{dateval(1999, 1, 8), dateval(2001, 5, 5)}}},
			{
				typ: "daterangearr",
				in:  `{"(,2001-05-05)","[1999-01-08,)","[1999-01-08,2001-05-05)"}`,
				want: &[][2]time.Time{
					{{}, dateval(2001, 5, 5)},
					{dateval(1999, 1, 8), {}},
					{dateval(1999, 1, 8), dateval(2001, 5, 5)}}},
			{
				typ:  "daterangearr",
				in:   `{"(,)","(,)","(,)"}`,
				want: &[][2]time.Time{{}, {}, {}}},
		},
	}}.execute(t)
}
