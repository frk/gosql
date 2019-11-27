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
//
// The io.EOF value is re-declared here only so that client code will not have
// to import the "io" package when using the "fql" package.
var EOF = io.EOF

const (
	// The eof rune signals the graceful end of a filter rule expression.
	eof rune = -1

	// These are the whitespace characters defined by Go itself.
	spacechars = " \t\r\n"
)

// The nullbit type is used to represent a value that is either "null" or not.
type nullbit bool

const (
	Null nullbit = true // Null represents a "null" value in a filter rule.
)

// CmpOp represents a comparison operator.
type CmpOp uint32

const (
	_     CmpOp = iota
	CmpEq       // ":" equality operator
	CmpNe       // ":!" inequality operator
	CmpGt       // ":>" greater-than operator
	CmpLt       // ":<" less-than operator
	CmpGe       // ":>=" greater-than-or-equal operator
	CmpLe       // ":<=" less-than-or-equal operator
)

// length returns the rune count of the CmpOp but minus the ":".
func (t CmpOp) length() int {
	switch t {
	case CmpEq:
		return 0
	case CmpNe:
		return 1
	case CmpGt:
		return 1
	case CmpLt:
		return 1
	case CmpGe:
		return 2
	case CmpLe:
		return 2
	}
	return -1
}

// A Token is the set of lexical tokens of FQL.
type Token uint32

const (
	_      Token = iota
	LPAREN       // (
	RPAREN       // )
	AND          // ;
	OR           // ,
	RULE         // <key>:[op]<value>
)

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
	// The current number of read open parentheses without a matching
	// closing parenthesis.
	group int
	// The type of the last token read.
	last Token
	// If the last Token is RULE this field will hold the parsed Rule,
	// otherwise it will be nil.
	rule *Rule
	// The first error encountered during tokenization.
	err error
}

// NewTokenizer returns a new FQL Tokenizer for the given string.
func NewTokenizer(fqlString string) *Tokenizer {
	return &Tokenizer{input: fqlString}
}

// Rule returns the current parsed Rule node, or nil if there isn't one.
func (t *Tokenizer) Rule() *Rule {
	return t.rule
}

// Next scans, parses, and returns the next token.
func (t *Tokenizer) Next() (Token, error) {
	if t.err != nil {
		return 0, t.err
	}
	t.rule = nil // reset

	t.eatspace()
	r := t.next()
	switch r {
	case eof:
		if t.group > 0 {
			t.err = &Error{Pos: t.pos, Code: ErrNoClosingParen}
			return 0, t.err
		}
		t.err = EOF
		return 0, t.err
	case '(':
		t.group += 1
		t.last = LPAREN
		return LPAREN, nil
	case ')':
		if t.last == LPAREN {
			t.err = &Error{Pos: t.pos, LastToken: t.last, CurrentToken: RPAREN, Code: ErrBadTokenSequence}
			return 0, t.err
		}
		if t.group < 1 {
			t.err = &Error{Pos: t.pos, Code: ErrExtraClosingParen}
			return 0, t.err
		}

		t.group -= 1
		t.last = RPAREN
		return RPAREN, nil
	case ';':
		if t.last == AND || t.last == OR || t.last == LPAREN {
			t.err = &Error{Pos: t.pos, LastToken: t.last, CurrentToken: AND, Code: ErrBadTokenSequence}
			return 0, t.err
		}
		t.last = AND
		return AND, nil
	case ',':
		if t.last == AND || t.last == OR || t.last == LPAREN {
			t.err = &Error{Pos: t.pos, LastToken: t.last, CurrentToken: OR, Code: ErrBadTokenSequence}
			return 0, t.err
		}
		t.last = OR
		return OR, nil
	}

	if !isfirstkeyrune(r) {
		badKey := t.current()
		t.backup()
		t.err = &Error{Key: badKey, Pos: t.pos, Code: ErrBadKey}
		return 0, t.err
	}
	t.backup()

	if t.last == RPAREN || t.last == RULE {
		t.err = &Error{Pos: t.pos, LastToken: t.last, CurrentToken: RULE, Code: ErrBadTokenSequence}
		return 0, t.err
	}

	rule := new(Rule)
	if err := rule.parse(t); err != nil {
		t.err = err
		return 0, t.err
	}

	t.last = RULE
	t.rule = rule
	return RULE, nil
}

// next returns the next rune in the input.
func (t *Tokenizer) next() (r rune) {
	if t.pos >= len(t.input) {
		t.width = 0
		return eof
	}
	r, t.width = utf8.DecodeRuneInString(t.input[t.pos:])
	t.pos += t.width
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (t *Tokenizer) backup() {
	t.pos -= t.width
}

// accept consumes the next rune if it is from the valid set.
func (t *Tokenizer) accept(valid string) bool {
	if strings.IndexRune(valid, t.next()) >= 0 {
		return true
	}
	t.backup()
	return false
}

// acceptrun consumes a run of runes from the valid set.
func (t *Tokenizer) acceptrun(valid string) {
	for strings.ContainsRune(valid, t.next()) {
	}
	t.backup()
}

// eatspace consumes and ignores a run of space runes.
func (t *Tokenizer) eatspace() {
	t.acceptrun(spacechars)
	t.ignore()
}

// ignore skips over the pending input before the current position.
func (t *Tokenizer) ignore() {
	t.start = t.pos
}

// current returns the current segment in the input.
func (t *Tokenizer) current() string {
	return t.input[t.start:t.pos]
}

// A Rule represents a filter rule and consists of a Key, a CmpOp, and a Value.
type Rule struct {
	Key string
	Cmp CmpOp
	Val interface{}
}

// parse parses the next segment of the input as a rule. The next
// rune in the input is known to be a valid initial rule-key rune.
func (rule *Rule) parse(t *Tokenizer) (err error) {
	if err = rule.parsekey(t); err != nil {
		return err
	}
	if err = rule.parsecmp(t); err != nil {
		return err
	}
	if err = rule.parsevalue(t); err != nil {
		return err
	}
	return nil
}

// parsekey parses the next segment of the input as a rule-key. The next
// rune in the input is known to be a valid initial rule-key rune.
func (rule *Rule) parsekey(t *Tokenizer) error {
	keystart := t.pos
	for {
		r := t.next()
		if iskeyrune(r) || (r == '.' && isfirstkeyrune(t.next())) {
			continue
		}

		break
	}
	t.backup()

	key := t.input[keystart:t.pos]

	t.acceptrun(spacechars)
	if t.next() != ':' {
		keyend := len(t.input)
		if i := strings.IndexByte(t.input[keystart:], ':'); i > -1 {
			keyend = keystart + i
		}

		badKey := strings.TrimSpace(t.input[keystart:keyend])
		return &Error{Key: badKey, Pos: keystart, Code: ErrBadKey}
	}

	rule.Key = key
	return nil
}

// parsecmp parses the next segment of the input as a rule-op. The rule's
// key-value separator ":" is known to be present.
func (rule *Rule) parsecmp(t *Tokenizer) error {
	t.acceptrun(spacechars)
	switch r := t.next(); r {
	case eof:
		return &Error{Pos: t.pos, Key: rule.Key, Code: ErrNoRuleValue}
	case '!':
		rule.Cmp = CmpNe
		return nil
	case '>':
		if t.next() == '=' {
			rule.Cmp = CmpGe
			return nil
		}
		t.backup()
		rule.Cmp = CmpGt
		return nil
	case '<':
		if t.next() == '=' {
			rule.Cmp = CmpLe
			return nil
		}
		t.backup()
		rule.Cmp = CmpLt
		return nil
	}

	t.backup()
	rule.Cmp = CmpEq
	return nil
}

// parsevalue parses the next segment of the input as a rule-value.
func (rule *Rule) parsevalue(t *Tokenizer) (err error) {
	t.acceptrun(spacechars)

	r := t.next()
	t.start = t.pos
	switch r {
	case eof:
		return &Error{Pos: t.pos, Key: rule.Key, Code: ErrNoRuleValue}
	case '"':
		return rule.parsetext(t)
	case 'd':
		return rule.parsetime(t)
	}

	t.backup()
	t.start = t.pos

	if r == 'n' {
		// get the position of the comparison operator  (for error reporting)
		cmppos := t.pos - rule.Cmp.length()

		if err := rule.parsenull(t); err != nil {
			return err
		}

		// Assuming the value is a valid null, then check that the current
		// rule's comparison operator is compatible with that type of value.
		if rule.Cmp != CmpEq && rule.Cmp != CmpNe {
			return &Error{Pos: cmppos, Key: rule.Key, Cmp: rule.Cmp, Code: ErrBadNullOp}
		}
		return nil
	}

	if r == 'f' || r == 't' {
		// get the position of the comparison operator  (for error reporting)
		cmppos := t.pos - rule.Cmp.length()

		if err := rule.parsebool(t); err != nil {
			return err
		}

		// Assuming the value is a valid boolean, then check that the current
		// rule's comparison operator is compatible with that type of value.
		if rule.Cmp != CmpEq && rule.Cmp != CmpNe {
			val, _ := rule.Val.(bool)
			return &Error{Pos: cmppos, Key: rule.Key, Cmp: rule.Cmp, Val: strconv.FormatBool(val), Code: ErrBadBooleanOp}
		}
		return nil
	}

	if r == '+' || r == '-' || r == '.' || (r >= '0' && r <= '9') {
		return rule.parsenumber(t)
	}

	return &Error{Pos: t.pos, Key: rule.Key, Code: ErrNoRuleValue}
}

// parsetext parses the next segment of the input as a string and sets it as the
// rule's value. The opening double quote `"` is known to be present. parsetext
// scans and parses up until the next unescaped double quote.
func (rule *Rule) parsetext(t *Tokenizer) error {
	var esc bool
Loop:
	for {
		switch r := t.next(); {
		case r == eof:
			return &Error{Pos: t.pos, Key: rule.Key, Val: t.input[t.start:], Code: ErrNoClosingDoubleQuote}
		case r == '\\':
			esc = !esc
			if esc {
				t.backup()
				// Delete the escape char so it doesn't get escaped by Go.
				t.input = t.input[:t.pos] + t.input[t.pos+1:]
			}
		case r == '"' && !esc:
			break Loop
		default:
			esc = false
		}
	}

	rule.Val = t.input[t.start : t.pos-1]
	return nil
}

// parsetime parses the next segment of the input as time.Time and sets it as the
// rule's value. The value should be an integer denoted by a preceding 'd' which
// is known to be present, the integer should represent the number of seconds
// elapsed since January 1, 1970 UTC.
//
// NOTE(mkopriva): timezones are currently not supported.
func (rule *Rule) parsetime(t *Tokenizer) error {
	// optional leading sign
	t.accept("+-")

	t.acceptrun("0123456789")
	d, err := strconv.ParseInt(t.current(), 10, 64)
	if err != nil {
		return &Error{Pos: t.pos, Key: rule.Key, Val: t.current(), Code: ErrBadDuration}
	}

	rule.Val = time.Unix(d, 0)
	return nil
}

// parsenull parses the next segment of the input as the Null const and sets it
// as the rule's value. It is known that the next rune is 'n'.
func (rule *Rule) parsenull(t *Tokenizer) error {
	if strings.HasPrefix(t.input[t.pos:], "null") {
		t.pos += 4 // len("null")
		rule.Val = Null
		return nil
	}
	return errors.New("bad null")

}

// parsebool parses the next segment in the input as a bool and sets it as
// the rule's value. It is known that the next rune is either 't' or 'f'.
func (rule *Rule) parsebool(t *Tokenizer) error {
	if strings.HasPrefix(t.input[t.pos:], "true") {
		t.pos += 4 // len("true")
		rule.Val = true
		return nil
	} else if strings.HasPrefix(t.input[t.pos:], "false") {
		t.pos += 5 // len("false")
		rule.Val = false
		return nil
	}
	return &Error{Pos: t.pos, Key: rule.Key, Code: ErrBadBoolean}
}

// parsenumber parses the next segment in the input as either an int64 or a float64
// and sets it as the rule's value. The value can be an integer, a float, or a float
// with an exponent, it can also be preceded by a hyphen in which case it will be
// parsed as negative.
func (rule *Rule) parsenumber(t *Tokenizer) (err error) {
	const digits = "0123456789"
	var isfloat bool

	// optional leading sign
	t.accept("+-")

	t.acceptrun(digits)
	if t.accept(".") {
		isfloat = true
		t.acceptrun(digits)
	}
	if t.accept("eE") {
		t.accept("+-")
		t.acceptrun(digits)
	}

	if isfloat {
		rule.Val, err = strconv.ParseFloat(t.current(), 64)
	} else {
		rule.Val, err = strconv.ParseInt(t.current(), 10, 64)
	}
	if err != nil {
		return &Error{Pos: t.pos, Key: rule.Key, Val: t.current(), Code: ErrBadNumber}
	}
	return nil
}

// iskeyrune reports whether r is a valid "key" rune.
func iskeyrune(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

// isfirstkeyrune reports whether r is a valid first "key" rune.
func isfirstkeyrune(r rune) bool {
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
	Cmp          CmpOp
	Val          string
	LastToken    Token
	CurrentToken Token
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
		return fmt.Sprintf("pos(%d): Unexpected closing parenthesis.", e.Pos)
	case ErrNoClosingParen:
		return fmt.Sprintf("pos(%d): Missing closing parenthesis to match the openning parenthesis.", e.Pos)
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
		return fmt.Sprintf("pos(%d): Invalid operator %q for %q's null value.", e.Pos, cmpop2string[e.Cmp], key)
	case ErrBadBooleanOp:
		return fmt.Sprintf("pos(%d): Invalid operator %q for %q's boolean value %q.", e.Pos, cmpop2string[e.Cmp], key, val)
	case ErrBadKey:
		return fmt.Sprintf("pos(%d): The key %q is not valid, keys can contain only alphanumeric"+
			" characters and the underscore character [_0-9a-Z].", e.Pos, key)
	case ErrBadTokenSequence:
		return fmt.Sprintf("pos(%d): Invalid token sequence. %q token cannot be followed by a %q token.", e.Pos, token2string[e.LastToken], token2string[e.CurrentToken])
	}
	return "bad syntax"
}

// TODO
type RuleError struct {
	Pos int
	Key string
	Cmp CmpOp
	Val string
}

// used or debugging
var cmpop2string = map[CmpOp]string{
	CmpEq: ":",
	CmpNe: ":!",
	CmpGt: ":>",
	CmpLt: ":<",
	CmpGe: ":>=",
	CmpLe: ":<=",
}

// used or debugging
var token2string = map[Token]string{
	LPAREN: "(",
	RPAREN: ")",
	AND:    ";",
	OR:     ",",
	RULE:   "<Rule>",
}
