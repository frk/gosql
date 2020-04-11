package convert

import (
	"testing"
	"time"
)

func TestDateRange(t *testing.T) {
	testlist{{
		valuer: func() interface{} {
			return new(DateRangeFromTimeArray2)
		},
		scanner: func() (interface{}, interface{}) {
			v := DateRangeToTimeArray2{Val: new([2]time.Time)}
			return v, v.Val
		},
		data: []testdata{
			{input: nil, output: new([2]time.Time)},
			{
				input:  [2]time.Time{dateval(1999, 1, 8), dateval(2001, 5, 5)},
				output: &[2]time.Time{dateval(1999, 1, 8), dateval(2001, 5, 5)}},
			{
				input:  [2]time.Time{{}, dateval(2001, 5, 5)},
				output: &[2]time.Time{{}, dateval(2001, 5, 5)}},
			{
				input:  [2]time.Time{dateval(1999, 1, 8), {}},
				output: &[2]time.Time{dateval(1999, 1, 8), {}}},
			{
				input:  [2]time.Time{{}, {}},
				output: &[2]time.Time{{}, {}}},
		},
	}, {
		valuer: func() interface{} {
			return nil // string
		},
		scanner: func() (interface{}, interface{}) {
			s := new(string)
			return s, s
		},
		data: []testdata{
			{
				input:  string("[1999-01-08,2001-05-05)"),
				output: strptr(`[1999-01-08,2001-05-05)`)},
			{
				input:  string("[1999-01-08,2001-05-05]"),
				output: strptr(`[1999-01-08,2001-05-06)`)},
			{
				input:  string("(1999-01-08,2001-05-05]"),
				output: strptr(`[1999-01-09,2001-05-06)`)},
			{
				input:  string("[,2001-05-05)"),
				output: strptr(`(,2001-05-05)`)},
			{
				input:  string("[1999-01-08,]"),
				output: strptr(`[1999-01-08,)`)},
			{
				input:  string("[,]"),
				output: strptr(`(,)`)},
		},
	}, {
		valuer: func() interface{} {
			return nil // []byte
		},
		scanner: func() (interface{}, interface{}) {
			s := new([]byte)
			return s, s
		},
		data: []testdata{
			{
				input:  []byte("[1999-01-08,2001-05-05)"),
				output: bytesptr(`[1999-01-08,2001-05-05)`)},
			{
				input:  []byte("[1999-01-08,2001-05-05]"),
				output: bytesptr(`[1999-01-08,2001-05-06)`)},
			{
				input:  []byte("(1999-01-08,2001-05-05]"),
				output: bytesptr(`[1999-01-09,2001-05-06)`)},
			{
				input:  []byte("[,2001-05-05)"),
				output: bytesptr(`(,2001-05-05)`)},
			{
				input:  []byte("[1999-01-08,]"),
				output: bytesptr(`[1999-01-08,)`)},
			{
				input:  []byte("[,]"),
				output: bytesptr(`(,)`)},
		},
	}}.execute(t, "daterange")
}
