package fql

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// EOF is the error returned by the Tokenizer when no more input is available.
var EOF = io.EOF

type nullType int

type boolType bool

const (
	Null nullType = 1 // Null represents a "null" value in a filter rule.

	True  boolType = true
	False boolType = false

	// The eof rune signals the graceful end of a filter rule expression.
	eof = -1

	// These are the space characters defined by Go itself.
	spaceChars = " \t\r\n"
)

// OpType represents the type of a comparison operator.
type OpType uint32

const (
	// OpEq looks like ":" and indicates the equality operator.
	OpEq OpType = 1 + iota
	// OpNe looks like ":!" and indicates the inequality operator.
	OpNe
	// OpGt looks like ":>" and indicates the greater-than operator.
	OpGt
	// OpLt looks like ":<" and indicates the less-than operator.
	OpLt
	// OpGe looks like ":>=" and indicates the greater-than-or-equal operator.
	OpGe
	// OpLe looks like ":<=" and indicates the less-than-or-equal operator.
	OpLe
)

// String returns a string representation of the OpType.
func (t OpType) String() string {
	switch t {
	case OpEq:
		return "="
	case OpNe:
		return "!="
	case OpGt:
		return ">"
	case OpLt:
		return "<"
	case OpGe:
		return ">="
	case OpLe:
		return "<="
	}
	return "Invalid(" + strconv.Itoa(int(t)) + ")"
}

// length returns the rune count of the OpType.
func (t OpType) length() int {
	switch t {
	case OpEq:
		return 0
	case OpNe:
		return 1
	case OpGt:
		return 1
	case OpLt:
		return 1
	case OpGe:
		return 2
	case OpLe:
		return 2
	}
	return -1
}

// A TokenType is the type of a Token.
type TokenType uint32

const (
	// TokenGroupStart looks like "(" and indicates the start of a rule group.
	TokenGroupStart TokenType = 1 + iota
	// TokenGroupEnd looks like ")" and indicates the end of a rule group.
	TokenGroupEnd
	// TokenAND looks like ";" and indicates the logical AND operator.
	TokenAND
	// TokenOR looks like "," and indicates the logical OR operator.
	TokenOR
	// TokenRule looks like "<key>:[op]<value>" and indicates a filter rule.
	TokenRule
)

// String returns a string representation of the TokenType.
func (t TokenType) String() string {
	switch t {
	case TokenGroupStart:
		return "GroupStart"
	case TokenGroupEnd:
		return "GroupEnd"
	case TokenAND:
		return "AND"
	case TokenOR:
		return "OR"
	case TokenRule:
		return "Rule"
	}
	return "Invalid(" + strconv.Itoa(int(t)) + ")"
}

// A Token consists of a TokenType and a Rule if the TokenType is TokenRule,
// otherwise Rule will be nil.
type Token struct {
	Type TokenType
	Rule *Rule
}

// A Rule consists of a Key, an OpType, and a Value.
type Rule struct {
	Key string
	Op  OpType
	Val interface{}
}

// A Tokenizer returns a stream of FQL Tokens.
type Tokenizer struct {
	// The fql text to tokenize.
	input string
	// Current position of the scanner in the input.
	pos int
	// The width of the last rune read from input, allows "unreading" the
	// last read rune.
	width int
	// The current token's starting position in the input.
	start int
	// The current number of read open round brackets without a matching
	// closing bracket.
	group int
	// The type of the last token read.
	last TokenType
	// The current rule's op, used to check if it can be applied to the
	// subsequent value. For example if the rule's value is parsed as a
	// boolean we make sure that the op is not ( ">" | ">=" | "<" | "<=" ).
	op OpType
	// The current rule's key, used for error reporting.
	key string
	// The first error encountered during tokenization.
	err error
}

// NewTokenizer returns a new FQL Tokenizer for the given string.
func NewTokenizer(fqlString string) *Tokenizer {
	return &Tokenizer{input: fqlString}
}

// nextRune returns the next rune in the input.
func (t *Tokenizer) nextRune() (r rune) {
	if t.pos >= len(t.input) {
		t.width = 0
		return eof
	}
	r, t.width = utf8.DecodeRuneInString(t.input[t.pos:])
	t.pos += t.width
	return r
}

// backupRune steps back one rune. Can only be called once per call of next.
func (t *Tokenizer) backupRune() {
	t.pos -= t.width
}

// accept consumes the next rune if it is from the valid set.
func (t *Tokenizer) accept(valid string) bool {
	if strings.IndexRune(valid, t.nextRune()) >= 0 {
		return true
	}
	t.backupRune()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (t *Tokenizer) acceptRun(valid string) {
	for strings.ContainsRune(valid, t.nextRune()) {
	}
	t.backupRune()
}

// ignore skips over the pending input before this point.
func (t *Tokenizer) ignore() {
	t.start = t.pos
}

// ignoreSpace consumes and ignores a run of space runes.
func (t *Tokenizer) ignoreSpace() {
	t.acceptRun(spaceChars)
	t.ignore()
}

// cursegment returns the current segment in the input.
func (t *Tokenizer) cursegment() string {
	return t.input[t.start:t.pos]
}

// Next scans, parses, and returns the next token.
func (t *Tokenizer) Next() (*Token, error) {
	if t.err != nil {
		return nil, t.err
	}

	t.ignoreSpace()
	r := t.nextRune()
	switch r {
	case eof:
		if t.group > 0 {
			t.err = &Error{Pos: t.pos, Code: ErrNoClosingParen}
			return nil, t.err
		}
		t.err = EOF
		return nil, t.err
	case '(':
		t.group += 1
		t.last = TokenGroupStart
		return &Token{Type: TokenGroupStart}, nil
	case ')':
		if t.last == TokenGroupStart {
			t.err = &Error{Pos: t.pos, LastToken: t.last, CurrentToken: TokenGroupEnd, Code: ErrBadTokenSequence}
			return nil, t.err
		}
		if t.group < 1 {
			t.err = &Error{Pos: t.pos, Code: ErrExtraClosingParen}
			return nil, t.err
		}

		t.group -= 1
		t.last = TokenGroupEnd
		return &Token{Type: TokenGroupEnd}, nil
	case ';':
		if t.last == TokenAND || t.last == TokenOR || t.last == TokenGroupStart {
			t.err = &Error{Pos: t.pos, LastToken: t.last, CurrentToken: TokenAND, Code: ErrBadTokenSequence}
			return nil, t.err
		}
		t.last = TokenAND
		return &Token{Type: TokenAND}, nil
	case ',':
		if t.last == TokenAND || t.last == TokenOR || t.last == TokenGroupStart {
			t.err = &Error{Pos: t.pos, LastToken: t.last, CurrentToken: TokenOR, Code: ErrBadTokenSequence}
			return nil, t.err
		}
		t.last = TokenOR
		return &Token{Type: TokenOR}, nil
	}

	if !isFirstKeyRune(r) {
		badKey := t.cursegment()
		t.backupRune()
		t.err = &Error{Key: badKey, Pos: t.pos, Code: ErrBadKey}
		return nil, t.err
	}
	t.backupRune()

	if t.last == TokenGroupEnd || t.last == TokenRule {
		t.err = &Error{Pos: t.pos, LastToken: t.last, CurrentToken: TokenRule, Code: ErrBadTokenSequence}
		return nil, t.err
	}
	t.last = TokenRule

	tok, err := parseRuleToken(t)
	if err != nil {
		t.err = err
		return nil, t.err
	}
	return tok, nil
}

// parseRuleToken parses the next segment of the input as a rule. The next
// rune in the input is known to be a valid initial rule-key rune.
func parseRuleToken(t *Tokenizer) (tok *Token, err error) {
	tok = &Token{Type: TokenRule, Rule: new(Rule)}
	if tok.Rule.Key, err = parseRuleKey(t); err != nil {
		return nil, err
	}
	t.key = tok.Rule.Key // track key

	if tok.Rule.Op, err = parseRuleOp(t); err != nil {
		return nil, err
	}
	t.op = tok.Rule.Op // track operator

	if tok.Rule.Val, err = parseRuleValue(t); err != nil {
		return nil, err
	}

	return tok, nil
}

// parseRuleKey parses the next segment of the input as a rule-key. The next
// rune in the input is known to be a valid initial rule-key rune.
func parseRuleKey(t *Tokenizer) (string, error) {
	keyStart := t.pos
	for {
		r := t.nextRune()
		if isKeyRune(r) || (r == '.' && isFirstKeyRune(t.nextRune())) {
			continue
		}

		break
	}
	t.backupRune()

	key := t.input[keyStart:t.pos]

	t.acceptRun(spaceChars)
	if t.nextRune() != ':' {
		keyEnd := len(t.input)
		if i := strings.IndexByte(t.input[keyStart:], ':'); i > -1 {
			keyEnd = keyStart + i
		}

		badKey := strings.TrimSpace(t.input[keyStart:keyEnd])
		return "", &Error{Key: badKey, Pos: keyStart, Code: ErrBadKey}
	}
	return key, nil
}

// parseRuleOp parses the next segment of the input as a rule-op. The rule's
// key-value separator ":" is known to be present.
func parseRuleOp(t *Tokenizer) (OpType, error) {
	t.acceptRun(spaceChars)
	switch r := t.nextRune(); r {
	case eof:
		return 0, &Error{Pos: t.pos, Key: t.key, Code: ErrNoRuleValue}
	case '!':
		return OpNe, nil
	case '>':
		if t.nextRune() == '=' {
			return OpGe, nil
		}
		t.backupRune()
		return OpGt, nil
	case '<':
		if t.nextRune() == '=' {
			return OpLe, nil
		}
		t.backupRune()
		return OpLt, nil
	}

	t.backupRune()
	return OpEq, nil
}

// parseRuleValue parses the next segment of the input as a rule-value.
func parseRuleValue(t *Tokenizer) (v interface{}, err error) {
	t.acceptRun(spaceChars)

	r := t.nextRune()
	t.start = t.pos
	switch r {
	case eof:
		return 0, &Error{Pos: t.pos, Key: t.key, Code: ErrNoRuleValue}
	case '"':
		return parseText(t)
	case 'd':
		return parseTime(t)
	}

	t.backupRune()
	t.start = t.pos
	if r == 'n' {
		// get the operator position (for error reporting)
		opPos := t.pos - t.op.length()

		val, err := parseNull(t)
		if err != nil {
			return nil, err
		}

		// Assuming the value is a valid null check that the current
		// rule's operator is compatible with that type of value.
		if t.op != OpEq && t.op != OpNe {
			return nil, &Error{Pos: opPos, Key: t.key, Op: t.op, Code: ErrBadNullOp}
		}
		return val, nil
	}
	if r == 'f' || r == 't' {
		// get the operator position (for error reporting)
		opPos := t.pos - t.op.length()

		val, err := parseBool(t)
		if err != nil {
			return nil, err
		}

		// Assuming the value is a valid boolean check that the current
		// rule's operator is compatible with that type of value.
		if t.op != OpEq && t.op != OpNe {
			return nil, &Error{Pos: opPos, Key: t.key, Op: t.op, Val: strconv.FormatBool(bool(val)), Code: ErrBadBooleanOp}
		}
		return val, nil
	}
	if r == '+' || r == '-' || r == '.' || (r >= '0' && r <= '9') {
		return parseNumber(t)
	}

	return nil, &Error{Pos: t.pos, Key: t.key, Code: ErrNoRuleValue}
}

// parseText parses the next segment of the input as a string and returns it.
// The opening double quote '"' is known to be present. parseText scans and
// parses up until the next unescaped double quote.
func parseText(t *Tokenizer) (string, error) {
	var esc bool
Loop:
	for {
		switch r := t.nextRune(); {
		case r == eof:
			return "", &Error{Pos: t.pos, Key: t.key, Val: t.input[t.start:], Code: ErrNoClosingDoubleQuote}
		case r == '\\':
			esc = !esc
			if esc {
				t.backupRune()
				// Delete the escape char so it doesn't get escaped by Go.
				t.input = t.input[:t.pos] + t.input[t.pos+1:]
			}
		case r == '"' && !esc:
			break Loop
		default:
			esc = false
		}
	}

	text := t.input[t.start : t.pos-1]
	return text, nil
}

// parseTime parses the next segment of the input as time.Time and returns it.
// The value should be an integer denoted by a preceding 'd' which is known
// to be present, the integer should represent the number of seconds elapsed since
// January 1, 1970 UTC.
//
// NOTE(mkopriva): timezones are currently not supported.
func parseTime(t *Tokenizer) (time.Time, error) {
	// optional leading sign
	t.accept("+-")

	t.acceptRun("0123456789")
	d, err := strconv.ParseInt(t.cursegment(), 10, 64)
	if err != nil {
		return time.Time{}, &Error{Pos: t.pos, Key: t.key, Val: t.cursegment(), Code: ErrBadDuration}
	}
	return time.Unix(d, 0), nil
}

// parseNull parses the next segment of the input as the Null const and
// returns it. It is known that the next rune is 'n'.
func parseNull(t *Tokenizer) (nullType, error) {
	if strings.HasPrefix(t.input[t.pos:], "null") {
		t.pos += 4 // len("null")
		return Null, nil
	}
	return 0, errors.New("bad null")

}

// parseBool parses the next segment in the input as a bool and returns it.
// It is known that the next rune is either 't' or 'f'.
func parseBool(t *Tokenizer) (boolType, error) {
	if strings.HasPrefix(t.input[t.pos:], "true") {
		t.pos += 4 // len("true")
		return True, nil
	} else if strings.HasPrefix(t.input[t.pos:], "false") {
		t.pos += 5 // len("false")
		return False, nil
	}
	return false, &Error{Pos: t.pos, Key: t.key, Code: ErrBadBoolean}
}

// parseNumber parses the next segment in the input as either an int64 or a float64
// and returns it. The value should be an integer, a float, or a float with an exponent,
// it can also be preceded by a hyphen in which case it will be parsed as negative.
func parseNumber(t *Tokenizer) (v interface{}, err error) {
	const digits = "0123456789"
	var isfloat bool

	// optional leading sign
	t.accept("+-")

	t.acceptRun(digits)
	if t.accept(".") {
		isfloat = true
		t.acceptRun(digits)
	}
	if t.accept("eE") {
		t.accept("+-")
		t.acceptRun(digits)
	}

	if isfloat {
		v, err = strconv.ParseFloat(t.cursegment(), 64)
	} else {
		v, err = strconv.ParseInt(t.cursegment(), 10, 64)
	}
	if err != nil {
		return nil, &Error{Pos: t.pos, Key: t.key, Val: t.cursegment(), Code: ErrBadNumber}
	}
	return v, nil
}

// isKeyRune reports whether r is a valid "key" rune.
func isKeyRune(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

// isFirstKeyRune reports whether r is a valid first "key" rune.
func isFirstKeyRune(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

type ErrorCode uint32

const (
	ErrExtraClosingParen ErrorCode = 1 + iota
	ErrNoClosingParen
	ErrNoClosingDoubleQuote
	ErrNoRuleValue
	ErrBadBoolean
	ErrBadNumber
	ErrBadDuration
	ErrBadNullOp
	ErrBadBooleanOp
	ErrBadKey
	ErrBadTokenSequence
)

type Error struct {
	Code         ErrorCode
	Pos          int
	Key          string
	Op           OpType
	Val          string
	LastToken    TokenType
	CurrentToken TokenType
}

func (e *Error) Error() string {
	key := e.Key
	if len(key) > 10 {
		key = key[:10] + "..."
	}
	val := e.Val
	if len(val) > 10 {
		val = val[:10] + "..."
	}

	switch e.Code {
	case ErrExtraClosingParen:
		return fmt.Sprintf("pos(%d): Unexpected closing round bracket.", e.Pos)
	case ErrNoClosingParen:
		return fmt.Sprintf("pos(%d): Missing closing round bracket to match the open bracket.", e.Pos)
	case ErrNoClosingDoubleQuote:
		return fmt.Sprintf("pos(%d): The %q string %q is missing an end quote.", e.Pos, key, val)
	case ErrNoRuleValue:
		return fmt.Sprintf("pos(%d): Missing rule value for key %q.", e.Pos, key)
	case ErrBadBoolean:
		return fmt.Sprintf("pos(%d): Invalid boolean value for key %q. Expted the"+
			" values 'true' or 'false'.", e.Pos, key)
	case ErrBadNumber:
		return fmt.Sprintf("pos(%d): %q is not a valid number value for key %q.", e.Pos, val, key)
	case ErrBadDuration:
		return fmt.Sprintf("pos(%d): %q is not a valid duration value for key %q.", e.Pos, val, key)
	case ErrBadNullOp:
		return fmt.Sprintf("pos(%d): Invalid operator %q for %q's null value.", e.Pos, e.Op, key)
	case ErrBadBooleanOp:
		return fmt.Sprintf("pos(%d): Invalid operator %q for %q's boolean value %q.", e.Pos, e.Op, key, val)
	case ErrBadKey:
		return fmt.Sprintf("pos(%d): The key %q is not valid, keys can contain only alphanumeric"+
			" characters and the underscore character [_0-9a-Z].", e.Pos, key)
	case ErrBadTokenSequence:
		return fmt.Sprintf("pos(%d): Invalid token sequence. %q token cannot be followed by a %q token.", e.Pos, e.LastToken, e.CurrentToken)
	}
	return "bad syntax"
}
