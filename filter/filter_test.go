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

type myString string

func convertToMyString(v any) (any, error) {
	if s, ok := v.(string); ok {
		return myString(s), nil
	}
	return v, nil
}

var test_colmap = map[string]Column{
	"a": {Name: "col_a", IsNULLable: true},
	"b": {Name: "col_b"},
	"c": {Name: "col_c"},
	"d": {Name: "col_d"},
	"e": {Name: "col_e", ConvertValue: convertToMyString},
	"f": {Name: "col_f"},
	"g": {Name: "col_g"},
	"h": {Name: "col_h"},
	"i": {Name: "col_i"},
	"j": {Name: "col_j"},
}

func TestFilter(t *testing.T) {
	type result struct {
		where  string
		params []interface{}
	}

	tests := []struct {
		name string
		run  func(*Constructor) error
		pos  int
		err  error
		want result
	}{{
		name: "test_fql_1",
		run: func(c *Constructor) error {
			return c.UnmarshalFQL("a:true")
		},
		want: result{where: ` WHERE col_a IS TRUE`},
	}, {
		name: "test_fql_2",
		run: func(c *Constructor) error {
			return c.UnmarshalFQL("b:!true")
		},
		want: result{where: ` WHERE col_b IS NOT TRUE`},
	}, {
		name: "test_fql_3",
		run: func(c *Constructor) error {
			return c.UnmarshalFQL("c:!false;d:false")
		},
		want: result{where: ` WHERE col_c IS NOT FALSE AND col_d IS FALSE`},
	}, {
		name: "test_fql_4",
		run: func(c *Constructor) error {
			return c.UnmarshalFQL("a:null;b:!null")
		},
		want: result{where: ` WHERE col_a IS NULL AND col_b IS NOT NULL`},
	}, {
		name: "test_fql_5",
		run: func(c *Constructor) error {
			return c.UnmarshalFQL("a:>18")
		},
		want: result{where: ` WHERE col_a > $1`, params: []interface{}{int64(18)}},
	}, {
		name: "test_fql_6",
		run: func(c *Constructor) error {
			return c.UnmarshalFQL("b:>=-90")
		},
		want: result{where: ` WHERE col_b >= $1`, params: []interface{}{int64(-90)}},
	}, {
		name: "test_fql_7",
		run: func(c *Constructor) error {
			return c.UnmarshalFQL("c:<-45.004532")
		},
		want: result{where: ` WHERE col_c < $1`, params: []interface{}{float64(-45.004532)}},
	}, {
		name: "test_fql_8",
		run: func(c *Constructor) error {
			return c.UnmarshalFQL("d:<=d1257894000")
		},
		want: result{where: ` WHERE col_d <= $1`, params: []interface{}{time.Unix(1257894000, 0)}},
	}, {
		name: "test_fql_9",
		run: func(c *Constructor) error {
			return c.UnmarshalFQL(`e:!"John Doe"`)
		},
		want: result{where: ` WHERE col_e <> $1`, params: []interface{}{myString("John Doe")}},
	}, {
		name: "test_fql_10",
		run: func(c *Constructor) error {
			return c.UnmarshalFQL("(a:123)")
		},
		want: result{where: ` WHERE (col_a = $1)`, params: []interface{}{int64(123)}},
	}, {
		name: "test_fql_11",
		run: func(c *Constructor) error {
			return c.UnmarshalFQL("(a:123),b:456")
		},
		want: result{where: ` WHERE (col_a = $1) OR col_b = $2`, params: i64s(123, 456)},
	}, {
		name: "test_fql_12",
		run: func(c *Constructor) error {
			return c.UnmarshalFQL("a:123;(b:456,c:789)")
		},
		want: result{where: ` WHERE col_a = $1 AND (col_b = $2 OR col_c = $3)`,
			params: i64s(123, 456, 789)},
	}, {
		name: "test_fql_13",
		run: func(c *Constructor) error {
			return c.UnmarshalFQL("a:123;(b:456,c:789)")
		},
		pos: 24,
		want: result{where: ` WHERE col_a = $25 AND (col_b = $26 OR col_c = $27)`,
			params: i64s(123, 456, 789)},
	}, {
		name: "test_fql_14",
		run: func(c *Constructor) error {
			return c.UnmarshalFQL("a:1;(b:2,(c:3;d:4;(e:5,f:6)),(g:7;h:8));(i:9,j:0)")
		},
		want: result{where: ` WHERE col_a = $1 AND (col_b = $2 OR ` +
			`(col_c = $3 AND col_d = $4 AND (col_e = $5 OR col_f = $6))` +
			` OR (col_g = $7 AND col_h = $8)) AND (col_i = $9 OR col_j = $10)`,
			params: i64s(1, 2, 3, 4, 5, 6, 7, 8, 9, 0)},
	}, {
		name: "test_fql_15",
		run: func(c *Constructor) error {
			return c.UnmarshalFQL("a:1;(b:2,(c:3;d:4;(e:5,f:6)),(g:7;h:8));(i:9,j:0)")
		},
		pos: 10,
		want: result{where: ` WHERE col_a = $11 AND (col_b = $12 OR ` +
			`(col_c = $13 AND col_d = $14 AND (col_e = $15 OR col_f = $16))` +
			` OR (col_g = $17 AND col_h = $18)) AND (col_i = $19 OR col_j = $20)`,
			params: i64s(1, 2, 3, 4, 5, 6, 7, 8, 9, 0)},
	}, {
		name: "test_fql_error_1",
		run: func(c *Constructor) error {
			return c.UnmarshalFQL("foo:1")
		},
		err:  UnknownColumnKeyError{Key: "foo"},
		want: result{},
	}, {
		name: "test_fql_error_2",
		run: func(c *Constructor) error {
			return c.UnmarshalFQL("a:123;(b:456,x:789)")
		},
		err:  UnknownColumnKeyError{Key: "x"},
		want: result{},
	}, {
		name: "test_sort_1",
		run: func(c *Constructor) error {
			return c.UnmarshalSort("-a")
		},
		want: result{where: ` ORDER BY col_a DESC NULLS LAST`},
	}, {
		name: "test_sort_2",
		run: func(c *Constructor) error {
			return c.UnmarshalSort("a,b,-c")
		},
		want: result{where: ` ORDER BY col_a ASC NULLS LAST, col_b ASC, col_c DESC`},
	}, {
		name: "test_sort_3",
		run: func(c *Constructor) error {
			return c.UnmarshalSort("-c,b,a")
		},
		want: result{where: ` ORDER BY col_c DESC, col_b ASC, col_a ASC NULLS LAST`},
	}, {
		name: "test_sort_error_1",
		run: func(c *Constructor) error {
			return c.UnmarshalSort("-c,b,x")
		},
		err: UnknownColumnKeyError{Key: "x"},
	}, {
		name: "test_sort_error_2",
		run: func(c *Constructor) error {
			return c.UnmarshalSort("foo")
		},
		err: UnknownColumnKeyError{Key: "foo"},
	}, {
		name: "test_text_search_1",
		run: func(c *Constructor) error {
			c.TextSearch("foo bar")
			return nil
		},
		want: result{where: ` WHERE t."tsvec" @@ to_tsquery('simple', $1)`,
			params: []interface{}{"foo:* & bar:*"}},
	}, {
		name: "test_text_search_2",
		run: func(c *Constructor) error {
			c.TextSearch("foo")
			c.Or(nil)
			c.TextSearch("bar")
			return nil
		},
		want: result{where: ` WHERE t."tsvec" @@ to_tsquery('simple', $1)` +
			` OR t."tsvec" @@ to_tsquery('simple', $2)`,
			params: []interface{}{"foo:*", "bar:*"}},
	}, {
		name: "test_text_search_3",
		run: func(c *Constructor) error {
			c.TextSearch("foo bar")
			c.UnmarshalFQL("(a:>123;b:!null)")
			return nil
		},
		want: result{where: ` WHERE t."tsvec" @@ to_tsquery('simple', $1) AND (col_a > $2 AND col_b IS NOT NULL)`,
			params: []interface{}{"foo:* & bar:*", int64(123)}},
	}, {
		name: "test_and",
		run: func(c *Constructor) error {
			c.Col("col_a", "=", 123)
			c.And(func() {
				c.Col("col_b", "=", 123)
				c.And(func() {
					c.Col("col_c", "=", 123)
				})
				c.Col("col_d", "=", 123)
			})
			c.Col("col_e", "=", 123)
			return nil
		},
		want: result{where: ` WHERE col_a = $1 AND (col_b = $2 AND (col_c = $3) AND col_d = $4) AND col_e = $5`,
			params: []interface{}{123, 123, 123, 123, 123}},
	}, {
		name: "test_and_2",
		run: func(c *Constructor) error {
			c.And(func() {
				c.Col("col_a", "=", 123)
				c.Or(nil)
				c.Col("col_a", "=", 0)
			})
			c.And(func() {
				c.Col("col_b", "=", 876)
				c.Or(nil)
				c.Col("col_b", "=", 0)
			})
			return nil
		},
		want: result{where: ` WHERE (col_a = $1 OR col_a = $2) AND (col_b = $3 OR col_b = $4)`,
			params: []interface{}{123, 0, 876, 0}},
	}, {
		name: "test_or",
		run: func(c *Constructor) error {
			c.Col("col_a", "=", 123)
			c.Or(func() {
				c.Col("col_b", "=", 123)
				c.Or(func() {
					c.Col("col_c", "=", 123)
					c.Col("col_d", "=", 123)
				})
			})
			c.Col("col_e", "=", 123)
			return nil
		},
		want: result{where: ` WHERE col_a = $1 OR (col_b = $2 OR (col_c = $3 AND col_d = $4)) AND col_e = $5`,
			params: []interface{}{123, 123, 123, 123, 123}},
	}, {
		name: "test_order_by_1",
		run: func(c *Constructor) error {
			c.OrderBy("col_a", true, true)
			return nil
		},
		want: result{where: ` ORDER BY col_a DESC NULLS FIRST`},
	}, {
		name: "test_order_by_2",
		run: func(c *Constructor) error {
			c.OrderBy("col_b", false, true)
			return nil
		},
		want: result{where: ` ORDER BY col_b ASC NULLS FIRST`},
	}, {
		name: "test_order_by_3",
		run: func(c *Constructor) error {
			c.OrderBy("col_c", true, false)
			return nil
		},
		want: result{where: ` ORDER BY col_c DESC NULLS LAST`},
	}, {
		name: "test_order_by_4",
		run: func(c *Constructor) error {
			c.OrderBy("col_d", false, false)
			return nil
		},
		want: result{where: ` ORDER BY col_d ASC NULLS LAST`},
	}, {
		name: "test_order_by_5",
		run: func(c *Constructor) error {
			c.OrderBy("col_a", true, true)
			c.OrderBy("col_b", false, true)
			c.OrderBy("col_c", true, false)
			c.OrderBy("col_d", false, false)
			return nil
		},
		want: result{where: ` ORDER BY col_a DESC NULLS FIRST,` +
			` col_b ASC NULLS FIRST,` +
			` col_c DESC NULLS LAST,` +
			` col_d ASC NULLS LAST`},
	}, {
		name: "test_full",
		run: func(c *Constructor) error {
			if err := c.UnmarshalFQL("a:>123;b:!null"); err != nil {
				return err
			}
			if err := c.UnmarshalSort("-c"); err != nil {
				return err
			}
			c.Col("col_d", "<=", 123)
			c.Or(func() {
				c.Col("col_e", ">", 10)
				c.Col("col_e", "<", 20)
			})
			c.Or(nil)
			c.TextSearch("foo bar baz")
			c.Limit(5)
			c.Offset(10)
			return nil
		},
		want: result{where: ` WHERE col_a > $1 AND col_b IS NOT NULL` +
			` AND col_d <= $2 OR (col_e > $3 AND col_e < $4)` +
			` OR t."tsvec" @@ to_tsquery('simple', $5)` +
			` ORDER BY col_c DESC` +
			` LIMIT 5` +
			` OFFSET 10`,
			params: []interface{}{int64(123), 123, 10, 20, "foo:* & bar:* & baz:*"}},
	}, {
		name: "test_=any",
		run: func(c *Constructor) error {
			if err := c.UnmarshalFQL("a:>123;b:!null"); err != nil {
				return err
			}
			if err := c.UnmarshalSort("-c"); err != nil {
				return err
			}
			c.And(func() {
				c.Col("kyb_assignee_id", "= ANY", []int{123, 3543, 920348, 28})
				c.Or(nil)
				c.Col("sponsor_assignee_id", "= ANY", []int{123, 3543, 920348, 28})
			})
			c.Limit(5)
			c.Offset(10)
			return nil
		},
		want: result{where: ` WHERE col_a > $1 AND col_b IS NOT NULL ` +
			`AND (kyb_assignee_id = ANY ($2::int4[]) ` +
			`OR sponsor_assignee_id = ANY ($3::int4[])) ` +
			`ORDER BY col_c DESC LIMIT 5 OFFSET 10`,
			params: []interface{}{int64(123), []int{123, 3543, 920348, 28}, []int{123, 3543, 920348, 28}}},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Constructor{colmap: test_colmap, tscol: `t."tsvec"`, strict: true}
			err := tt.run(c)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Error(e)
			} else if err == nil {
				where, params := c.Filter().ToSQL(tt.pos)
				got := result{where, params}
				if e := compare.Compare(got, tt.want); e != nil {
					t.Error(e)
				}
			}
		})
	}
}

func Test_toTSQuery(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"foo", "foo:*"},
		{"foo bar", "foo:* & bar:*"},
		{" foo   bar   ", "foo:* & bar:*"},
		{" foo   bar   baz", "foo:* & bar:* & baz:*"},
	}

	for _, tt := range tests {
		got := toTSQuery(tt.input)
		if got != tt.want {
			t.Errorf("got=%q; want=%q", got, tt.want)
		}
	}
}
