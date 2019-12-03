package sqlang

import (
	"github.com/frk/gosql/internal/writer"
)

type Predicate interface {
	Node
	predicateNode()
}

// TRUTH type is used to represet a truth value.
type TRUTH uint8

const (
	UNKNOWN TRUTH = iota
	TRUE
	FALSE
)

func (t TRUTH) Walk(w *writer.Writer) {
	switch t {
	case UNKNOWN:
		w.Write("UNKNOWN")
	case TRUE:
		w.Write("TRUE")
	case FALSE:
		w.Write("FALSE")
	}
}

// TruthPredicate produces an SQL predicate for comparing a row value against a truth value.
type TruthPredicate struct {
	Not       bool
	Truth     TRUTH
	Predicand ValueExpr
}

func (p TruthPredicate) Walk(w *writer.Writer) {
	p.Predicand.Walk(w)
	if p.Not {
		w.Write(" IS NOT ")
	} else {
		w.Write(" IS ")
	}
	p.Truth.Walk(w)
}

// CMPOP is used to specify the comparison operator of the ComparisonPredicate.
type CMPOP uint8

const (
	EQUAL CMPOP = 1 + iota
	NOT_EQUAL
	NOT_EQUAL2 // support for !=
	GREATER_THAN
	GREATER_THAN_EQUAL
	LESS_THAN
	LESS_THAN_EQUAL
)

func (op CMPOP) Walk(w *writer.Writer) {
	switch op {
	case EQUAL:
		w.Write(" = ")
	case NOT_EQUAL:
		w.Write(" <> ")
	case NOT_EQUAL2:
		w.Write(" != ")
	case GREATER_THAN:
		w.Write(" > ")
	case GREATER_THAN_EQUAL:
		w.Write(" >= ")
	case LESS_THAN:
		w.Write(" < ")
	case LESS_THAN_EQUAL:
		w.Write(" <= ")
	}
}

// ComparisonPredicate produces an SQL predicate for comparing two row values.
type ComparisonPredicate struct {
	Cmp        CMPOP
	LPredicand ValueExpr
	RPredicand ValueExpr
}

func (p ComparisonPredicate) Walk(w *writer.Writer) {
	p.LPredicand.Walk(w)
	p.Cmp.Walk(w)
	p.RPredicand.Walk(w)
}

// SYMMETRY is used to specify the symmetry of the BetweenPredicate.
type SYMMETRY uint8

const (
	ASYMMETRIC SYMMETRY = 1 + iota
	SYMMETRIC
)

func (s SYMMETRY) Walk(w *writer.Writer) {
	switch s {
	case ASYMMETRIC:
		w.Write("ASYMMETRIC ")
	case SYMMETRIC:
		w.Write("SYMMETRIC ")
	}
}

// BetweenPredicate produces the "[ NOT ] BETWEEN [ ASYMMETRIC | SYMMETRIC ]"
// SQL predicate for range comparison.
type BetweenPredicate struct {
	Not       bool
	Sym       SYMMETRY
	Predicand ValueExpr
	LowEnd    ValueExpr // the low end of the range
	HighEnd   ValueExpr // the hight end of the range
}

func (p BetweenPredicate) Walk(w *writer.Writer) {
	p.Predicand.Walk(w)

	if p.Not {
		w.Write(" NOT BETWEEN ")
	} else {
		w.Write(" BETWEEN ")
	}
	p.Sym.Walk(w)

	p.LowEnd.Walk(w)
	w.Write(" AND ")
	p.HighEnd.Walk(w)
}

// InPredicate produces the "[ NOT ] IN" SQL predicate for quantified comparison.
type InPredicate struct {
	Not       bool
	Predicand ValueExpr
	ValueList ValueExpr
}

func (p InPredicate) Walk(w *writer.Writer) {
	p.Predicand.Walk(w)

	if p.Not {
		w.Write(" NOT IN ")
	} else {
		w.Write(" IN ")
	}

	w.Write("(")
	p.ValueList.Walk(w)
	w.Write(")")
}

// LikePredicate produces the "[ NOT ] LIKE" SQL predicate for pattern-match comparison.
type LikePredicate struct {
	Not       bool
	Predicand ValueExpr
	Pattern   ValueExpr
	Escape    EscapeCharacter
}

func (p LikePredicate) Walk(w *writer.Writer) {
	p.Predicand.Walk(w)

	if p.Not {
		w.Write(" NOT LIKE ")
	} else {
		w.Write(" LIKE ")
	}

	p.Pattern.Walk(w)

	if p.Escape != nil {
		w.Write(" ESCAPE ")
		p.Escape.Walk(w)
	}
}

// ILikePredicate produces the "[ NOT ] ILIKE" SQL predicate for case-insensitive
// pattern-match comparison. Note that the ILIKE predicate is postgres specific.
type ILikePredicate struct {
	Not       bool
	Predicand ValueExpr
	Pattern   ValueExpr
	Escape    EscapeCharacter
}

func (p ILikePredicate) Walk(w *writer.Writer) {
	p.Predicand.Walk(w)

	if p.Not {
		w.Write(" NOT ILIKE ")
	} else {
		w.Write(" ILIKE ")
	}

	p.Pattern.Walk(w)

	if p.Escape != nil {
		w.Write(" ESCAPE ")
		p.Escape.Walk(w)
	}
}

// SimilarPredicate produces the "[ NOT ] SIMILAR TO" SQL predicate for
// pattern-match comparison by means of regular expressions.
type SimilarPredicate struct {
	Not       bool
	Predicand ValueExpr
	Pattern   ValueExpr
	Escape    EscapeCharacter
}

func (p SimilarPredicate) Walk(w *writer.Writer) {
	p.Predicand.Walk(w)

	if p.Not {
		w.Write(" NOT SIMILAR TO ")
	} else {
		w.Write(" SIMILAR TO ")
	}

	p.Pattern.Walk(w)

	if p.Escape != nil {
		w.Write(" ESCAPE ")
		p.Escape.Walk(w)
	}
}

// REGEXOP is used to specify the operator of the RegexPredicate.
type REGEXOP uint8

const (
	MATCH REGEXOP = iota
	MATCH_CI
	NOT_MATCH
	NOT_MATCH_CI
)

func (op REGEXOP) Walk(w *writer.Writer) {
	switch op {
	case MATCH:
		w.Write(" ~ ")
	case MATCH_CI:
		w.Write(" ~* ")
	case NOT_MATCH:
		w.Write(" !~ ")
	case NOT_MATCH_CI:
		w.Write(" !~* ")
	}
}

// RegexPredicate produces a PostgreSQL predicate for regular expression pattern-match comparison.
type RegexPredicate struct {
	Op        REGEXOP
	Predicand ValueExpr
	Pattern   ValueExpr
}

func (p RegexPredicate) Walk(w *writer.Writer) {
	p.Predicand.Walk(w)
	p.Op.Walk(w)
	p.Pattern.Walk(w)
}

// NullPredicate produces the "IS [NOT] NULL" SQL predicate for testing against null values.
type NullPredicate struct {
	Not       bool
	Predicand ValueExpr
}

func (p NullPredicate) Walk(w *writer.Writer) {
	p.Predicand.Walk(w)
	if p.Not {
		w.Write(" IS NOT NULL ")
	} else {
		w.Write(" IS NULL ")
	}
}

// QUANTIFIER is used to specify the quantifier of the ArrayComparisonPredicate.
type QUANTIFIER uint8

const (
	ALL QUANTIFIER = iota
	ANY
	SOME
)

func (q QUANTIFIER) Walk(w *writer.Writer) {
	switch q {
	case ALL:
		w.Write("ALL ")
	case ANY:
		w.Write("ANY ")
	case SOME:
		w.Write("SOME ")
	}
}

// ArrayComparisonPredicate produces an SQL predicate for comparison against an array of values.
type ArrayComparisonPredicate struct {
	Predicand ValueExpr
	Cmp       CMPOP
	Qua       QUANTIFIER
	Array     ValueExpr
}

func (p ArrayComparisonPredicate) Walk(w *writer.Writer) {
	p.Predicand.Walk(w)
	p.Cmp.Walk(w)
	p.Qua.Walk(w)
	w.Write("(")
	p.Array.Walk(w)
	w.Write(")")
}

// DistinctPredicate produces the "IS [NOT] DISTINCT FROM" SQL predicate for
// equality comparison, treating nulls as ordinary values.
type DistinctPredicate struct {
	Not        bool
	LPredicate ValueExpr
	RPredicate ValueExpr
}

func (p DistinctPredicate) Walk(w *writer.Writer) {
	p.LPredicate.Walk(w)
	if p.Not {
		w.Write(" IS NOT DISTINCT FROM ")
	} else {
		w.Write(" IS DISTINCT FROM ")
	}
	p.RPredicate.Walk(w)
}

func (TruthPredicate) predicateNode()           {}
func (ComparisonPredicate) predicateNode()      {}
func (BetweenPredicate) predicateNode()         {}
func (InPredicate) predicateNode()              {}
func (LikePredicate) predicateNode()            {}
func (ILikePredicate) predicateNode()           {}
func (SimilarPredicate) predicateNode()         {}
func (RegexPredicate) predicateNode()           {}
func (NullPredicate) predicateNode()            {}
func (ArrayComparisonPredicate) predicateNode() {}
func (DistinctPredicate) predicateNode()        {}

func (TruthPredicate) boolValueExpr()           {}
func (ComparisonPredicate) boolValueExpr()      {}
func (BetweenPredicate) boolValueExpr()         {}
func (InPredicate) boolValueExpr()              {}
func (LikePredicate) boolValueExpr()            {}
func (ILikePredicate) boolValueExpr()           {}
func (SimilarPredicate) boolValueExpr()         {}
func (RegexPredicate) boolValueExpr()           {}
func (NullPredicate) boolValueExpr()            {}
func (ArrayComparisonPredicate) boolValueExpr() {}
func (DistinctPredicate) boolValueExpr()        {}

// TODO(mkopriva): these need SUBQUERY expression...
// type QuantifiedComparisonPredicate struct{  }
// func (QuantifiedComparisonPredicate) predicateNode() {}
// type ExistsPredicate struct{}
// func (ExistsPredicate) predicateNode()               {}

////////////////////////////////////////////////////////////////////////////////

type EscapeCharacter interface {
	Node
	escapeCharacterNode()
}

type EmptyEscapeCharacter struct{}

func (EmptyEscapeCharacter) Walk(w *writer.Writer) {
	w.Write("''")
}

type CharacterValueExpr byte

func (x CharacterValueExpr) Walk(w *writer.Writer) {
	w.Write("'" + string(x) + "'")
}

func (EmptyEscapeCharacter) escapeCharacterNode() {}
func (CharacterValueExpr) escapeCharacterNode()   {}
