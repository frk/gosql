package convert

import (
	"testing"
	"time"
)

func TestDate_Scanner(t *testing.T) {
	test_scanner{{
		scanner: func() (interface{}, interface{}) {
			d := &DateToTime{V: new(time.Time)}
			return d, d.V
		},
		rows: []test_scanner_row{
			{typ: "date", in: `1999-01-08`, want: dateptr(1999, 1, 8)},
			{typ: "date", in: `2001-05-05`, want: dateptr(2001, 5, 5)},
			{typ: "date", in: `2020-03-28`, want: dateptr(2020, 3, 28)},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			d := &DateToString{V: new(string)}
			return d, d.V
		},
		rows: []test_scanner_row{
			{typ: "date", in: nil, want: new(string)},
			{typ: "date", in: `1999-01-08`, want: strptr(`1999-01-08`)},
			{typ: "date", in: `2001-05-05`, want: strptr(`2001-05-05`)},
			{typ: "date", in: `2020-03-28`, want: strptr(`2020-03-28`)},
		},
	}, {
		scanner: func() (interface{}, interface{}) {
			d := &DateToByteSlice{V: new([]byte)}
			return d, d.V
		},
		rows: []test_scanner_row{
			{typ: "date", in: nil, want: new([]byte)},
			{typ: "date", in: `1999-01-08`, want: bytesptr(`1999-01-08`)},
			{typ: "date", in: `2001-05-05`, want: bytesptr(`2001-05-05`)},
			{typ: "date", in: `2020-03-28`, want: bytesptr(`2020-03-28`)},
		},
	}}.execute(t)
}

func TestDate_NoValuer(t *testing.T) {
	// NOTE(mkopriva): the "T00:00:00Z" extra suffix is what lib/pq appends
	// when decoding the postgres' date text into plain string/[]byte, that's
	// also the reason why custom scanners, see above, are provided for these
	// two builtin types, with that in mind, what is actually stored however
	// is just the provided date string without the suffix.
	test_valuer{{
		valuer: func() interface{} {
			return nil // time.Time
		},
		rows: []test_valuer_row{
			{typ: "date", in: nil, want: nil},
			{typ: "date", in: dateptr(1999, 1, 8), want: strptr(`1999-01-08T00:00:00Z`)},
			{typ: "date", in: dateptr(2001, 5, 5), want: strptr(`2001-05-05T00:00:00Z`)},
			{typ: "date", in: dateptr(2020, 3, 28), want: strptr(`2020-03-28T00:00:00Z`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // string
		},
		rows: []test_valuer_row{
			{typ: "date", in: nil, want: nil},
			{typ: "date", in: `1999-01-08`, want: strptr(`1999-01-08T00:00:00Z`)},
			{typ: "date", in: `2001-05-05`, want: strptr(`2001-05-05T00:00:00Z`)},
			{typ: "date", in: `2020-03-28`, want: strptr(`2020-03-28T00:00:00Z`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		rows: []test_valuer_row{
			{typ: "date", in: nil, want: nil},
			{typ: "date", in: []byte(`1999-01-08`), want: strptr(`1999-01-08T00:00:00Z`)},
			{typ: "date", in: []byte(`2001-05-05`), want: strptr(`2001-05-05T00:00:00Z`)},
			{typ: "date", in: []byte(`2020-03-28`), want: strptr(`2020-03-28T00:00:00Z`)},
		},
	}}.execute(t)
}
