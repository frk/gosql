package postgres

import (
	"github.com/frk/gosql/internal/analysis"
)

// relation kind
type RelKind string

const (
	RelKindOrdinaryTable    RelKind = "r"
	RelKindIndex            RelKind = "i"
	RelKindSequence         RelKind = "S"
	RelKindTOAST            RelKind = "t"
	RelKindView             RelKind = "v"
	RelKindMaterializedView RelKind = "m"
	RelKindCompositeType    RelKind = "c"
	RelKindForeignTable     RelKind = "f"
	RelKindPartitionedTable RelKind = "p"
	RelKindPartitionedIndex RelKind = "I"
)

// postgres type types
type TypeType string

const (
	TypeTypeBase      TypeType = "b"
	TypeTypeComposite TypeType = "c"
	TypeTypeDomain    TypeType = "d"
	TypeTypeEnum      TypeType = "e"
	TypeTypePseudo    TypeType = "p"
	TypeTypeRange     TypeType = "r"
)

// postgres type categories
type TypeCategory string

const (
	TypeCategoryArray       TypeCategory = "A"
	TypeCategoryBoolean     TypeCategory = "B"
	TypeCategoryComposite   TypeCategory = "C"
	TypeCategoryDatetime    TypeCategory = "D"
	TypeCategoryEnum        TypeCategory = "E"
	TypeCategoryGeometric   TypeCategory = "G"
	TypeCategoryNetaddress  TypeCategory = "I"
	TypeCategoryNumeric     TypeCategory = "N"
	TypeCategoryPseudo      TypeCategory = "P"
	TypeCategoryRange       TypeCategory = "R"
	TypeCategoryString      TypeCategory = "S"
	TypeCategoryTimespan    TypeCategory = "T"
	TypeCategoryUserdefined TypeCategory = "U"
	TypeCategoryBitstring   TypeCategory = "V"
	TypeCategoryUnknown     TypeCategory = "X"
)

// postgres constraint types
type ConstraintType string

const (
	ConstraintTypeCheck     ConstraintType = "c"
	ConstraintTypeFKey      ConstraintType = "f"
	ConstraintTypePKey      ConstraintType = "p"
	ConstraintTypeUnique    ConstraintType = "u"
	ConstraintTypeTrigger   ConstraintType = "t"
	ConstraintTypeExclusion ConstraintType = "x"
)

// postgres cast contexts
type CastContext string

const (
	CastContextExplicit   CastContext = "e"
	CastContextImplicit   CastContext = "i"
	CastContextAssignment CastContext = "a"
)

// Map of supported predicates to *equivalent* postgres comparison operators. For example
// the constructs IN and NOT IN are essentially the same as comparing the LHS to every
// element of the RHS with the operators "=" and "<>" respectively, and therefore the
// isin maps to "=" and notin maps to "<>".
var predicateToOprname = map[analysis.Predicate]string{
	analysis.IsEQ:        "=",
	analysis.NotEQ:       "<>",
	analysis.NotEQ2:      "<>",
	analysis.IsLT:        "<",
	analysis.IsGT:        ">",
	analysis.IsLTE:       "<=",
	analysis.IsGTE:       ">=",
	analysis.IsMatch:     "~",
	analysis.IsMatchi:    "~*",
	analysis.NotMatch:    "!~",
	analysis.NotMatchi:   "!~*",
	analysis.IsDistinct:  "<>",
	analysis.NotDistinct: "=",
	analysis.IsLike:      "~~",
	analysis.NotLike:     "!~~",
	analysis.IsILike:     "~~*",
	analysis.NotILike:    "!~~*",
	analysis.IsSimilar:   "~~",
	analysis.NotSimilar:  "!~~",
	analysis.IsIn:        "=",
	analysis.NotIn:       "<>",
}
