package fql

import (
	"testing"
	"time"

	"github.com/frk/compare"
)

func TestTokenizer(t *testing.T) {
	tests := []struct {
		test string
		fql  string
		want []*Token
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
			want: []*Token{{Type: TokenRule, Rule: &Rule{Key: "test_key", Op: OpEq, Val: True}}},
		}, {
			test: "PASS: parse 'false' value",
			fql:  "test_key:false",
			want: []*Token{{Type: TokenRule, Rule: &Rule{Key: "test_key", Op: OpEq, Val: False}}},
		}, {
			test: "PASS: parse 'false' value",
			fql:  "test_key:>18",
			want: []*Token{{Type: TokenRule, Rule: &Rule{Key: "test_key", Op: OpGt, Val: int64(18)}}},
		}, {
			test: "PASS: parse integer value",
			fql:  "test_key:>=90",
			want: []*Token{{Type: TokenRule, Rule: &Rule{Key: "test_key", Op: OpGe, Val: int64(90)}}},
		}, {
			test: "PASS: parse negative float value",
			fql:  "test_key:<-45.004532",
			want: []*Token{{Type: TokenRule, Rule: &Rule{Key: "test_key", Op: OpLt, Val: float64(-45.004532)}}},
		}, {
			test: "PASS: parse duration value",
			fql:  "test_key:<=d1257894000",
			want: []*Token{{Type: TokenRule, Rule: &Rule{Key: "test_key", Op: OpLe, Val: time.Unix(1257894000, 0)}}},
		}, {
			test: "PASS: parse text value",
			fql:  `test_key:!"John Doe"`,
			want: []*Token{{Type: TokenRule, Rule: &Rule{Key: "test_key", Op: OpNe, Val: "John Doe"}}},
		}, {
			test: "PASS: parse text value with escaped quotes",
			fql:  `test_key:"They said \"Hi!\"."`,
			want: []*Token{{Type: TokenRule, Rule: &Rule{Key: "test_key", Op: OpEq, Val: "They said \"Hi!\"."}}},
		}, {
			test: "PASS: parse text value with escaped escapes",
			fql:  `test_key:"They said \\\\\"Hi!\"."`,
			want: []*Token{{Type: TokenRule, Rule: &Rule{Key: "test_key", Op: OpEq, Val: "They said \\\\\"Hi!\"."}}},
		}, {
			test: "PASS: parse null",
			fql:  `test_key:null`,
			want: []*Token{{Type: TokenRule, Rule: &Rule{Key: "test_key", Op: OpEq, Val: Null}}},
		}, {
			test: "PASS: parse null 2",
			fql:  `test_key:!null`,
			want: []*Token{{Type: TokenRule, Rule: &Rule{Key: "test_key", Op: OpNe, Val: Null}}},
		}, {
			test: "PASS: parse integer 0",
			fql:  `test_key:!0`,
			want: []*Token{{Type: TokenRule, Rule: &Rule{Key: "test_key", Op: OpNe, Val: int64(0)}}},
		},

		// Test that logical operators are tokenized accurately.
		{
			test: "PASS: tokenize logical AND",
			fql:  "key_a:true;key_b:false",
			want: []*Token{
				{Type: TokenRule, Rule: &Rule{Key: "key_a", Op: OpEq, Val: True}},
				{Type: TokenAND},
				{Type: TokenRule, Rule: &Rule{Key: "key_b", Op: OpEq, Val: False}},
			},
		}, {
			test: "PASS: tokenize logical OR and AND",
			fql:  "key_a:true,key_b:false;key_c:>d1257894000",
			want: []*Token{
				{Type: TokenRule, Rule: &Rule{Key: "key_a", Op: OpEq, Val: True}},
				{Type: TokenOR},
				{Type: TokenRule, Rule: &Rule{Key: "key_b", Op: OpEq, Val: False}},
				{Type: TokenAND},
				{Type: TokenRule, Rule: &Rule{Key: "key_c", Op: OpGt, Val: time.Unix(1257894000, 0)}},
			},
		},

		// Test tokenization of grouped rules.
		{
			test: "PASS: tokenize simple group",
			fql:  "(test_key:123)",
			want: []*Token{
				{Type: TokenGroupStart},
				{Type: TokenRule, Rule: &Rule{Key: "test_key", Op: OpEq, Val: int64(123)}},
				{Type: TokenGroupEnd},
			},
		}, {
			test: "PASS: tokenize simple group together with non-group",
			fql:  "(key_a:123),key_b:456",
			want: []*Token{
				{Type: TokenGroupStart},
				{Type: TokenRule, Rule: &Rule{Key: "key_a", Op: OpEq, Val: int64(123)}},
				{Type: TokenGroupEnd},
				{Type: TokenOR},
				{Type: TokenRule, Rule: &Rule{Key: "key_b", Op: OpEq, Val: int64(456)}},
			},
		}, {
			test: "PASS: tokenize group with multiple rules",
			fql:  "key_a:123;(key_b:456,key_c:789)",
			want: []*Token{
				{Type: TokenRule, Rule: &Rule{Key: "key_a", Op: OpEq, Val: int64(123)}},
				{Type: TokenAND},
				{Type: TokenGroupStart},
				{Type: TokenRule, Rule: &Rule{Key: "key_b", Op: OpEq, Val: int64(456)}},
				{Type: TokenOR},
				{Type: TokenRule, Rule: &Rule{Key: "key_c", Op: OpEq, Val: int64(789)}},
				{Type: TokenGroupEnd},
			},
		}, {
			test: "PASS: tokenize nested groups",
			fql:  "a:1;(b:2,(c:3;d:4;(e:5,f:6)),(g:7;h:8));(i:9,j:0)",
			want: []*Token{
				{Type: TokenRule, Rule: &Rule{Key: "a", Op: OpEq, Val: int64(1)}},
				{Type: TokenAND},
				{Type: TokenGroupStart},
				{Type: TokenRule, Rule: &Rule{Key: "b", Op: OpEq, Val: int64(2)}},
				{Type: TokenOR},
				{Type: TokenGroupStart},
				{Type: TokenRule, Rule: &Rule{Key: "c", Op: OpEq, Val: int64(3)}},
				{Type: TokenAND},
				{Type: TokenRule, Rule: &Rule{Key: "d", Op: OpEq, Val: int64(4)}},
				{Type: TokenAND},
				{Type: TokenGroupStart},
				{Type: TokenRule, Rule: &Rule{Key: "e", Op: OpEq, Val: int64(5)}},
				{Type: TokenOR},
				{Type: TokenRule, Rule: &Rule{Key: "f", Op: OpEq, Val: int64(6)}},
				{Type: TokenGroupEnd},
				{Type: TokenGroupEnd},
				{Type: TokenOR},
				{Type: TokenGroupStart},
				{Type: TokenRule, Rule: &Rule{Key: "g", Op: OpEq, Val: int64(7)}},
				{Type: TokenAND},
				{Type: TokenRule, Rule: &Rule{Key: "h", Op: OpEq, Val: int64(8)}},
				{Type: TokenGroupEnd},
				{Type: TokenGroupEnd},
				{Type: TokenAND},
				{Type: TokenGroupStart},
				{Type: TokenRule, Rule: &Rule{Key: "i", Op: OpEq, Val: int64(9)}},
				{Type: TokenOR},
				{Type: TokenRule, Rule: &Rule{Key: "j", Op: OpEq, Val: int64(0)}},
				{Type: TokenGroupEnd},
			},
		},

		// Test that the tokenizer ignores whitespace correctly.
		{
			test: "PASS: ignore space between fql tokens but recognize space in text values",
			fql:  ` test_key   :    "  Hello "  `,
			want: []*Token{
				{Type: TokenRule, Rule: &Rule{Key: "test_key", Op: OpEq, Val: "  Hello "}},
			},
		}, {
			test: "PASS: ignore space between fql tokens",
			fql:  ` test_key   :  >=  123  `,
			want: []*Token{
				{Type: TokenRule, Rule: &Rule{Key: "test_key", Op: OpGe, Val: int64(123)}},
			},
		}, {
			test: "PASS: ignore space between fql tokens 2",
			fql:  ` test_key   :    false  `,
			want: []*Token{
				{Type: TokenRule, Rule: &Rule{Key: "test_key", Op: OpEq, Val: False}},
			},
		}, {
			test: "PASS: ignore space between fql tokens 3",
			fql:  ` test_key   :    .007  `,
			want: []*Token{
				{Type: TokenRule, Rule: &Rule{Key: "test_key", Op: OpEq, Val: float64(0.007)}},
			},
		}, {
			test: "PASS: ignore space between fql tokens 4",
			fql:  ` test_key   :    d-1234  `,
			want: []*Token{
				{Type: TokenRule, Rule: &Rule{Key: "test_key", Op: OpEq, Val: time.Unix(-1234, 0)}},
			},
		}, {
			test: "PASS: ignore space between fql tokens and groups of tokens",
			fql:  "key_a: 123 ; (  key_b: 456 , key_c : 789 ) ",
			want: []*Token{
				{Type: TokenRule, Rule: &Rule{Key: "key_a", Op: OpEq, Val: int64(123)}},
				{Type: TokenAND},
				{Type: TokenGroupStart},
				{Type: TokenRule, Rule: &Rule{Key: "key_b", Op: OpEq, Val: int64(456)}},
				{Type: TokenOR},
				{Type: TokenRule, Rule: &Rule{Key: "key_c", Op: OpEq, Val: int64(789)}},
				{Type: TokenGroupEnd},
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
			fql:  "key_a:123 key_b:456", err: &Error{Pos: 10, LastToken: TokenRule, CurrentToken: TokenRule, Code: ErrBadTokenSequence},
		}, {
			test: "FAIL: bad token sequence 2",
			fql:  "key_a:123,key_b:1234dff6544444;key_c:789", err: &Error{Pos: 20, LastToken: TokenRule, CurrentToken: TokenRule, Code: ErrBadTokenSequence},
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
			fql:  `key_a:123,key_b:"They said \\"Hi!\"."`, err: &Error{Pos: 29, LastToken: TokenRule, CurrentToken: TokenRule, Code: ErrBadTokenSequence},
		}, {
			test: "FAIL: bad bool operator",
			fql:  `key_a:>=true`, err: &Error{Pos: 6, Key: "key_a", Op: OpGe, Val: "true", Code: ErrBadBooleanOp},
		}, {
			test: "FAIL: bad null operator",
			fql:  `key_a:<null"`, err: &Error{Pos: 6, Key: "key_a", Op: OpLt, Code: ErrBadNullOp},
		}, {
			test: "FAIL: bad key syntax",
			fql:  `:key_a:<null"`, err: &Error{Pos: 0, Key: ":", Code: ErrBadKey},
		},
	}

	// // TODO: test edge cases
	// // - ((()), foo)
	// // - key_a:(123;key_b: ...

	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			var got []*Token
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
			}

			if e := compare.Compare(err, tt.err); e != nil {
				t.Error(e, "|", err, "|", tt.err)
			}
			if err := compare.Compare(got, tt.want); err != nil {
				t.Error(err)
			}
		})
	}
}
