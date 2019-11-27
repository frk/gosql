package fql

import (
	"testing"
	"time"

	"github.com/frk/compare"
)

func TestTokenizer(t *testing.T) {
	type list []interface{} // for less typing

	tests := []struct {
		test string
		fql  string
		want list
		err  error
	}{
		{
			test: "PASS: empty fql returns EOF immediately",
			fql:  "", want: nil,
		},

		// Test that a rule's value is parsed as expected.
		{
			test: "PASS: parse 'true' value",
			fql:  "test_key:true",
			want: list{RULE, &Rule{Key: "test_key", Cmp: CmpEq, Val: true}},
		}, {
			test: "PASS: parse 'false' value",
			fql:  "test_key:false",
			want: list{RULE, &Rule{Key: "test_key", Cmp: CmpEq, Val: false}},
		}, {
			test: "PASS: parse 'false' value",
			fql:  "test_key:>18",
			want: list{RULE, &Rule{Key: "test_key", Cmp: CmpGt, Val: int64(18)}},
		}, {
			test: "PASS: parse integer value",
			fql:  "test_key:>=90",
			want: list{RULE, &Rule{Key: "test_key", Cmp: CmpGe, Val: int64(90)}},
		}, {
			test: "PASS: parse negative float value",
			fql:  "test_key:<-45.004532",
			want: list{RULE, &Rule{Key: "test_key", Cmp: CmpLt, Val: float64(-45.004532)}},
		}, {
			test: "PASS: parse duration value",
			fql:  "test_key:<=d1257894000",
			want: list{RULE, &Rule{Key: "test_key", Cmp: CmpLe, Val: time.Unix(1257894000, 0)}},
		}, {
			test: "PASS: parse text value",
			fql:  `test_key:!"John Doe"`,
			want: list{RULE, &Rule{Key: "test_key", Cmp: CmpNe, Val: "John Doe"}},
		}, {
			test: "PASS: parse text value with escaped quotes",
			fql:  `test_key:"They said \"Hi!\"."`,
			want: list{RULE, &Rule{Key: "test_key", Cmp: CmpEq, Val: "They said \"Hi!\"."}},
		}, {
			test: "PASS: parse text value with escaped escapes",
			fql:  `test_key:"They said \\\\\"Hi!\"."`,
			want: list{RULE, &Rule{Key: "test_key", Cmp: CmpEq, Val: "They said \\\\\"Hi!\"."}},
		}, {
			test: "PASS: parse null",
			fql:  `test_key:null`,
			want: list{RULE, &Rule{Key: "test_key", Cmp: CmpEq, Val: Null}},
		}, {
			test: "PASS: parse null 2",
			fql:  `test_key:!null`,
			want: list{RULE, &Rule{Key: "test_key", Cmp: CmpNe, Val: Null}},
		}, {
			test: "PASS: parse integer 0",
			fql:  `test_key:!0`,
			want: list{RULE, &Rule{Key: "test_key", Cmp: CmpNe, Val: int64(0)}},
		},

		// Test that logical operators are tokenized accurately.
		{
			test: "PASS: tokenize logical AND",
			fql:  "key_a:true;key_b:false",
			want: list{
				RULE, &Rule{Key: "key_a", Cmp: CmpEq, Val: true},
				AND,
				RULE, &Rule{Key: "key_b", Cmp: CmpEq, Val: false},
			},
		}, {
			test: "PASS: tokenize logical OR and AND",
			fql:  "key_a:true,key_b:false;key_c:>d1257894000",
			want: list{
				RULE, &Rule{Key: "key_a", Cmp: CmpEq, Val: true},
				OR,
				RULE, &Rule{Key: "key_b", Cmp: CmpEq, Val: false},
				AND,
				RULE, &Rule{Key: "key_c", Cmp: CmpGt, Val: time.Unix(1257894000, 0)},
			},
		},

		// Test tokenization of grouped rules.
		{
			test: "PASS: tokenize simple group",
			fql:  "(test_key:123)",
			want: list{
				LPAREN,
				RULE, &Rule{Key: "test_key", Cmp: CmpEq, Val: int64(123)},
				RPAREN,
			},
		}, {
			test: "PASS: tokenize simple group together with non-group",
			fql:  "(key_a:123),key_b:456",
			want: list{
				LPAREN,
				RULE, &Rule{Key: "key_a", Cmp: CmpEq, Val: int64(123)},
				RPAREN,
				OR,
				RULE, &Rule{Key: "key_b", Cmp: CmpEq, Val: int64(456)},
			},
		}, {
			test: "PASS: tokenize group with multiple rules",
			fql:  "key_a:123;(key_b:456,key_c:789)",
			want: list{
				RULE, &Rule{Key: "key_a", Cmp: CmpEq, Val: int64(123)},
				AND,
				LPAREN,
				RULE, &Rule{Key: "key_b", Cmp: CmpEq, Val: int64(456)},
				OR,
				RULE, &Rule{Key: "key_c", Cmp: CmpEq, Val: int64(789)},
				RPAREN,
			},
		}, {
			test: "PASS: tokenize nested groups",
			fql:  "a:1;(b:2,(c:3;d:4;(e:5,f:6)),(g:7;h:8));(i:9,j:0)",
			want: list{
				RULE, &Rule{Key: "a", Cmp: CmpEq, Val: int64(1)},
				AND,
				LPAREN,
				RULE, &Rule{Key: "b", Cmp: CmpEq, Val: int64(2)},
				OR,
				LPAREN,
				RULE, &Rule{Key: "c", Cmp: CmpEq, Val: int64(3)},
				AND,
				RULE, &Rule{Key: "d", Cmp: CmpEq, Val: int64(4)},
				AND,
				LPAREN,
				RULE, &Rule{Key: "e", Cmp: CmpEq, Val: int64(5)},
				OR,
				RULE, &Rule{Key: "f", Cmp: CmpEq, Val: int64(6)},
				RPAREN,
				RPAREN,
				OR,
				LPAREN,
				RULE, &Rule{Key: "g", Cmp: CmpEq, Val: int64(7)},
				AND,
				RULE, &Rule{Key: "h", Cmp: CmpEq, Val: int64(8)},
				RPAREN,
				RPAREN,
				AND,
				LPAREN,
				RULE, &Rule{Key: "i", Cmp: CmpEq, Val: int64(9)},
				OR,
				RULE, &Rule{Key: "j", Cmp: CmpEq, Val: int64(0)},
				RPAREN,
			},
		},

		// Test that the tokenizer ignores whitespace correctly.
		{
			test: "PASS: ignore space between fql tokens but recognize space in text values",
			fql:  ` test_key   :    "  Hello "  `,
			want: list{
				RULE, &Rule{Key: "test_key", Cmp: CmpEq, Val: "  Hello "},
			},
		}, {
			test: "PASS: ignore space between fql tokens",
			fql:  ` test_key   :  >=  123  `,
			want: list{
				RULE, &Rule{Key: "test_key", Cmp: CmpGe, Val: int64(123)},
			},
		}, {
			test: "PASS: ignore space between fql tokens 2",
			fql:  ` test_key   :    false  `,
			want: list{
				RULE, &Rule{Key: "test_key", Cmp: CmpEq, Val: false},
			},
		}, {
			test: "PASS: ignore space between fql tokens 3",
			fql:  ` test_key   :    .007  `,
			want: list{
				RULE, &Rule{Key: "test_key", Cmp: CmpEq, Val: float64(0.007)},
			},
		}, {
			test: "PASS: ignore space between fql tokens 4",
			fql:  ` test_key   :    d-1234  `,
			want: list{
				RULE, &Rule{Key: "test_key", Cmp: CmpEq, Val: time.Unix(-1234, 0)},
			},
		}, {
			test: "PASS: ignore space between fql tokens and groups of tokens",
			fql:  "key_a: 123 ; (  key_b: 456 , key_c : 789 ) ",
			want: list{
				RULE, &Rule{Key: "key_a", Cmp: CmpEq, Val: int64(123)},
				AND,
				LPAREN,
				RULE, &Rule{Key: "key_b", Cmp: CmpEq, Val: int64(456)},
				OR,
				RULE, &Rule{Key: "key_c", Cmp: CmpEq, Val: int64(789)},
				RPAREN,
			},
		},

		// Test error handling.
		{
			test: "FAIL: space in key not allowed",
			fql:  `test key:"Hello"`, err: &Error{Key: "test key", Pos: 0, Code: ErrBadKey},
		}, {
			test: "FAIL: hyphen in key not allowed",
			fql:  `key_a:123;key-b:"???"`, err: &Error{Key: "key-b", Pos: 10, Code: ErrBadKey},
		}, {
			test: "FAIL: lone key",
			fql:  "test_key", err: &Error{Key: "test_key", Pos: 0, Code: ErrBadKey},
		}, {
			test: "FAIL: lone key 2",
			fql:  "   test_key   ", err: &Error{Key: "test_key", Pos: 3, Code: ErrBadKey},
		}, {
			test: "FAIL: bad key syntax",
			fql:  "key_a:123; - key_b:456", err: &Error{Key: "-", Pos: 11, Code: ErrBadKey},
		}, {
			test: "FAIL: bad token sequence",
			fql:  "key_a:123 key_b:456", err: &Error{Pos: 10, LastToken: RULE, CurrentToken: RULE, Code: ErrBadTokenSequence},
		}, {
			test: "FAIL: bad token sequence 2",
			fql:  "key_a:123,key_b:1234dff6544444;key_c:789", err: &Error{Pos: 20, LastToken: RULE, CurrentToken: RULE, Code: ErrBadTokenSequence},
		}, {
			test: "FAIL: extra closing parentheses",
			fql:  "key_a:123 ) key_b:456", err: &Error{Pos: 11, Code: ErrExtraClosingParen},
		}, {
			test: "FAIL: missing closing parentheses",
			fql:  "key_a:123;(key_b:456", err: &Error{Pos: 20, Code: ErrNoClosingParen},
		}, {
			test: "FAIL: missing rule value",
			fql:  "key_a:123,key_b:", err: &Error{Pos: 16, Key: "key_b", Code: ErrNoRuleValue},
		}, {
			test: "FAIL: missing rule value 2",
			fql:  "key_a:123,key_b:>=", err: &Error{Pos: 18, Key: "key_b", Code: ErrNoRuleValue},
		}, {
			test: "FAIL: missing rule value 3",
			fql:  "key_a:123,key_b:;key_c:789", err: &Error{Pos: 16, Key: "key_b", Code: ErrNoRuleValue},
		}, {
			test: "FAIL: bad boolean value",
			fql:  "key_a:123,key_b:truck;key_c:789", err: &Error{Pos: 16, Key: "key_b", Code: ErrBadBoolean},
		}, {
			test: "FAIL: bad boolean value 2",
			fql:  "key_a:123,key_b:t;key_c:789", err: &Error{Pos: 16, Key: "key_b", Code: ErrBadBoolean},
		}, {
			test: "FAIL: bad boolean value 3",
			fql:  "key_a:123,key_b:f;key_c:789", err: &Error{Pos: 16, Key: "key_b", Code: ErrBadBoolean},
		}, {
			test: "FAIL: bad boolean value 4",
			fql:  "key_a:123,key_b:falsy;key_c:789", err: &Error{Pos: 16, Key: "key_b", Code: ErrBadBoolean},
		}, {
			test: "FAIL: bad number value",
			fql:  "key_a:123,key_b:-;key_c:789", err: &Error{Pos: 17, Key: "key_b", Val: "-", Code: ErrBadNumber},
		}, {
			test: "FAIL: bad duration value",
			fql:  "key_a:123,key_b:d-;key_c:789", err: &Error{Pos: 18, Key: "key_b", Val: "-", Code: ErrBadDuration},
		}, {
			test: "FAIL: missing closing double quote",
			fql:  "key_a:123,key_b:\"Hello;key_c:789", err: &Error{Pos: 32, Key: "key_b", Val: "Hello;key_c:789", Code: ErrNoClosingDoubleQuote},
		}, {
			test: "FAIL: bad token sequence 3",
			fql:  `key_a:123,key_b:"They said \\"Hi!\"."`, err: &Error{Pos: 29, LastToken: RULE, CurrentToken: RULE, Code: ErrBadTokenSequence},
		}, {
			test: "FAIL: bad bool operator",
			fql:  `key_a:>=true`, err: &Error{Pos: 6, Key: "key_a", Cmp: CmpGe, Val: "true", Code: ErrBadBooleanOp},
		}, {
			test: "FAIL: bad null operator",
			fql:  `key_a:<null"`, err: &Error{Pos: 6, Key: "key_a", Cmp: CmpLt, Code: ErrBadNullOp},
		}, {
			test: "FAIL: bad key syntax",
			fql:  `:key_a:<null"`, err: &Error{Pos: 0, Key: ":", Code: ErrBadKey},
		}, {
			test: "FAIL: bad left-right parentheses sequence",
			fql:  `((()), abc)`, err: &Error{Pos: 4, Key: "", LastToken: LPAREN, CurrentToken: RPAREN, Code: ErrBadTokenSequence},
		}, {
			test: "FAIL: bad rule value",
			fql:  `key_a:(123;key_b:false`, err: &Error{Pos: 6, Key: "key_a", Code: ErrNoRuleValue},
		},
	}

	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			var got list
			var err error

			z := NewTokenizer(tt.fql)
			for {
				tok, e := z.Next()
				if e != nil {
					if e != EOF {
						err = e
						got = nil
					}
					break
				}
				got = append(got, tok)
				if tok == RULE {
					got = append(got, z.Rule())
				}
			}

			if e := compare.Compare(err, tt.err); e != nil {
				t.Error(e, "\n", err, "|", tt.err)
			}
			if err := compare.Compare(got, tt.want); err != nil {
				t.Error(err)
			}
		})
	}
}
