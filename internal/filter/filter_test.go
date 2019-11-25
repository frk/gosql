package filter

import (
	"testing"
	"time"

	"github.com/frk/compare"
)

// helper
func i64s(ii ...int64) (out []interface{}) {
	for _, i := range ii {
		out = append(out, i)
	}
	return out
}

var test_colmap = map[string]string{
	"a": "col_a",
	"b": "col_b",
	"c": "col_c",
	"d": "col_d",
	"e": "col_e",
	"f": "col_f",
	"g": "col_g",
	"h": "col_h",
	"i": "col_i",
	"j": "col_j",
}

func TestParseFQL(t *testing.T) {
	type result struct {
		clause string
		params []interface{}
	}

	tests := []struct {
		fql  string
		want result
		err  error
	}{{
		fql:  "a:true",
		want: result{clause: ` WHERE col_a IS TRUE`},
	}, {
		fql:  `b:!true`,
		want: result{clause: ` WHERE col_b IS NOT TRUE`},
	}, {
		fql:  `c:!false;d:false`,
		want: result{clause: ` WHERE col_c IS NOT FALSE AND col_d IS FALSE`},
	}, {
		fql:  `a:null;b:!null`,
		want: result{clause: ` WHERE col_a IS NULL AND col_b IS NOT NULL`},
	}, {
		fql:  `a:>18`,
		want: result{clause: ` WHERE col_a > $1`, params: []interface{}{int64(18)}},
	}, {
		fql:  `b:>=-90`,
		want: result{clause: ` WHERE col_b >= $1`, params: []interface{}{int64(-90)}},
	}, {
		fql:  `c:<-45.004532`,
		want: result{clause: ` WHERE col_c < $1`, params: []interface{}{float64(-45.004532)}},
	}, {
		fql:  `d:<=d1257894000`,
		want: result{clause: ` WHERE col_d <= $1`, params: []interface{}{time.Unix(1257894000, 0)}},
	}, {
		fql:  `e:!"John Doe"`,
		want: result{clause: ` WHERE col_e <> $1`, params: []interface{}{"John Doe"}},
	}, {
		fql:  `(a:123)`,
		want: result{clause: ` WHERE (col_a = $1)`, params: []interface{}{int64(123)}},
	}, {
		fql: `(a:123),b:456`,
		want: result{clause: ` WHERE (col_a = $1) OR col_b = $2`,
			params: i64s(123, 456)},
	}, {
		fql: `a:123;(b:456,c:789)`,
		want: result{clause: ` WHERE col_a = $1 AND (col_b = $2 OR col_c = $3)`,
			params: i64s(123, 456, 789)},
	}, {
		fql: `a:1;(b:2,(c:3;d:4;(e:5,f:6)),(g:7;h:8));(i:9,j:0)`,
		want: result{clause: ` WHERE col_a = $1 AND (col_b = $2 OR ` +
			`(col_c = $3 AND col_d = $4 AND (col_e = $5 OR col_f = $6))` +
			` OR (col_g = $7 AND col_h = $8)) AND (col_i = $9 OR col_j = $10)`,
			params: i64s(1, 2, 3, 4, 5, 6, 7, 8, 9, 0)},
	}}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			f := new(Filter)
			err := f.ParseFQL(tt.fql, test_colmap)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Error(e)
			} else {
				clause, params := f.ToSQL(), f.Params()
				got := result{clause, params}
				if e := compare.Compare(got, tt.want); e != nil {
					t.Error(e)
				}
			}
		})
	}
}
