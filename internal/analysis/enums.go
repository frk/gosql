package analysis

import (
	"go/types"
)

// TypeKind indicates the specific kind of a Go type.
type TypeKind uint

const (
	// basic
	TypeKindInvalid TypeKind = iota

	_basic_kind_start
	TypeKindBool
	TypeKindInt
	TypeKindInt8
	TypeKindInt16
	TypeKindInt32
	TypeKindInt64
	TypeKindUint
	TypeKindUint8
	TypeKindUint16
	TypeKindUint32
	TypeKindUint64
	TypeKindUintptr
	TypeKindFloat32
	TypeKindFloat64
	TypeKindComplex64
	TypeKindComplex128
	TypeKindString
	TypeKindUnsafePointer
	_basic_kind_end

	// non-basic
	TypeKindArray
	TypeKindInterface
	TypeKindMap
	TypeKindPtr
	TypeKindSlice
	TypeKindStruct
	TypeKindChan
	TypeKindFunc

	// alisases (basic)
	TypeKindByte = TypeKindUint8
	TypeKindRune = TypeKindInt32
)

func (k TypeKind) isBasic() bool { return _basic_kind_start < k && k < _basic_kind_end }

// BasicString returns a string representation of k.
func (k TypeKind) BasicString() string {
	if k.isBasic() {
		return basicTypeKinds[k]
	}
	return "<unknown>"
}

// Basic type kind string represenations indexed by typeKind.
var basicTypeKinds = [...]string{
	TypeKindBool:       "bool",
	TypeKindInt:        "int",
	TypeKindInt8:       "int8",
	TypeKindInt16:      "int16",
	TypeKindInt32:      "int32",
	TypeKindInt64:      "int64",
	TypeKindUint:       "uint",
	TypeKindUint8:      "uint8",
	TypeKindUint16:     "uint16",
	TypeKindUint32:     "uint32",
	TypeKindUint64:     "uint64",
	TypeKindUintptr:    "uintptr",
	TypeKindFloat32:    "float32",
	TypeKindFloat64:    "float64",
	TypeKindComplex64:  "complex64",
	TypeKindComplex128: "complex128",
	TypeKindString:     "string",
}

// typeKinds indexed by types.BasicKind.
var typesBasicKindToTypeKind = [...]TypeKind{
	types.Invalid:       TypeKindInvalid,
	types.Bool:          TypeKindBool,
	types.Int:           TypeKindInt,
	types.Int8:          TypeKindInt8,
	types.Int16:         TypeKindInt16,
	types.Int32:         TypeKindInt32,
	types.Int64:         TypeKindInt64,
	types.Uint:          TypeKindUint,
	types.Uint8:         TypeKindUint8,
	types.Uint16:        TypeKindUint16,
	types.Uint32:        TypeKindUint32,
	types.Uint64:        TypeKindUint64,
	types.Uintptr:       TypeKindUintptr,
	types.Float32:       TypeKindFloat32,
	types.Float64:       TypeKindFloat64,
	types.Complex64:     TypeKindComplex64,
	types.Complex128:    TypeKindComplex128,
	types.String:        TypeKindString,
	types.UnsafePointer: TypeKindUnsafePointer,
}

type QueryKind uint8

const (
	_ QueryKind = iota
	QueryKindInsert
	QueryKindUpdate
	QueryKindDelete

	_select_kind_start
	QueryKindSelect
	QueryKindSelectCount
	QueryKindSelectExists
	QueryKindSelectNotExists
	_select_kind_end
)

// isSelect reports whether or not the query kind is one of the select kinds.
func (k QueryKind) isSelect() bool { return _select_kind_start < k && k < _select_kind_end }

// isNonFromSelect reports whether or not the query kind is one of the non-from-select select kinds.
func (k QueryKind) isNonFromSelect() bool { return k.isSelect() && k != QueryKindSelect }

// String returns the string form of the QueryKind.
func (k QueryKind) String() string {
	switch k {
	case QueryKindInsert:
		return "Insert"
	case QueryKindUpdate:
		return "Update"
	case QueryKindDelete:
		return "Delete"
	case QueryKindSelect, QueryKindSelectCount, QueryKindSelectExists, QueryKindSelectNotExists:
		return "Select"
	}
	return "<Unknown QueryKind>"
}

// JoinType indicates the gosql.XxxJoin directive used in a query struct.
type JoinType uint8

const (
	_             JoinType = iota // no join
	JoinTypeCross                 // CROSS JOIN
	JoinTypeInner                 // INNER JOIN
	JoinTypeLeft                  // LEFT JOIN
	JoinTypeRight                 // RIGHT JOIN
	JoinTypeFull                  // FULL JOIN
)

// JoinTypes is an array of join directive names indexed by their respective JoinType values.
var JoinTypes = [...]string{
	JoinTypeCross: "CrossJoin",
	JoinTypeInner: "InnerJoin",
	JoinTypeLeft:  "LeftJoin",
	JoinTypeRight: "RightJoin",
	JoinTypeFull:  "FullJoin",
}

// stringToJoinType maps string literals to JoinType values. Used for parsing of directives.
var stringToJoinType map[string]JoinType

func init() {
	stringToJoinType = make(map[string]JoinType)
	for typ, str := range JoinTypes {
		if len(str) > 0 {
			stringToJoinType[str] = JoinType(typ)
		}
	}
}

// OrderDirection is used to specify the order direction in an ORDER BY clause.
type OrderDirection uint8

const (
	OrderAsc  OrderDirection = iota // ASC, default
	OrderDesc                       // DESC
)

// NullsPosition is used to specify the position of NULLs in an ORDER BY clause.
type NullsPosition uint8

const (
	_          NullsPosition = iota // none specified, i.e. default
	NullsFirst                      // NULLS FIRST
	NullsLast                       // NULLS LAST
)

// OverridingKind indicates the option used with the gosql.Override directive.
type OverridingKind uint8

const (
	_                OverridingKind = iota // no overriding
	OverridingSystem                       // OVERRIDING SYSTEM VALUE
	OverridingUser                         // OVERRIDING USER VALUE
)

// boolean operation
type Boolean uint8

const (
	_       Boolean = iota // no bool
	BoolAnd                // conjunction
	BoolOr                 // disjunction
	BoolNot                // negation
)

// Quantifier represents the type of a comparison predicate quantifier.
type Quantifier uint8

const (
	_         Quantifier = iota // no operator
	QuantAny                    // ANY
	QuantSome                   // SOME
	QuantAll                    // ALL
)

// stringToQuantifier is a map of string literals to supported quantifiers. Used for parsing of tags.
var stringToQuantifier = map[string]Quantifier{
	"any":  QuantAny,
	"some": QuantSome,
	"all":  QuantAll,
}

// Predicate represents the predicate type of a search condition.
type Predicate uint

const (
	_ Predicate = iota // no predicate

	_binary_pred_start
	IsEQ        // equals
	NotEQ       // not equals
	NotEQ2      // not equals
	IsLT        // less than
	IsGT        // greater than
	IsLTE       // less than or equal
	IsGTE       // greater than or equal
	IsDistinct  // IS DISTINCT FROM
	NotDistinct // IS NOT DISTINCT FROM
	_binary_pred_end

	_pattern_pred_start
	IsMatch    // match regular expression
	IsMatchi   // match regular expression (case insensitive)
	NotMatch   // not match regular expression
	NotMatchi  // not match regular expression (case insensitive)
	IsLike     // LIKE
	NotLike    // NOT LIKE
	IsILike    // ILIKE
	NotILike   // NOT ILIKE
	IsSimilar  // IS SIMILAR TO
	NotSimilar // IS NOT SIMILAR TO
	_pattern_pred_end

	_array_pred_start
	IsIn  // IN
	NotIn // NOT IN
	_array_pred_end

	_range_pred_start
	IsBetween      // BETWEEN x AND y
	NotBetween     // NOT BETWEEN x AND y
	IsBetweenSym   // BETWEEN SYMMETRIC x AND y
	NotBetweenSym  // NOT BETWEEN SYMMETRIC x AND y
	IsBetweenAsym  // BETWEEN ASYMMETRIC x AND y
	NotBetweenAsym // NOT BETWEEN ASYMMETRIC x AND y
	_range_pred_end

	_null_pred_start
	IsNull  // IS NULL
	NotNull // IS NOT NULL
	_null_pred_end

	_truth_pred_start
	IsTrue     // IS TRUE
	NotTrue    // IS NOT TRUE
	IsFalse    // IS FALSE
	NotFalse   // IS NOT FALSE
	IsUnknown  // IS UNKNOWN
	NotUnknown // IS NOT UNKNOWN
	_truth_pred_end
)

var predicates = [...]string{
	IsEQ:        "=",
	NotEQ:       "<>",
	NotEQ2:      "!=",
	IsLT:        "<",
	IsGT:        ">",
	IsLTE:       "<=",
	IsGTE:       ">=",
	IsDistinct:  "isdistinct",
	NotDistinct: "notdistinct",

	IsMatch:    "~",
	IsMatchi:   "~*",
	NotMatch:   "!~",
	NotMatchi:  "!~*",
	IsLike:     "islike",
	NotLike:    "notlike",
	IsILike:    "isilike",
	NotILike:   "notilike",
	IsSimilar:  "issimilar",
	NotSimilar: "notsimilar",

	IsIn:  "isin",
	NotIn: "notin",

	IsBetween:      "isbetween",
	NotBetween:     "notbetween",
	IsBetweenSym:   "isbetweensym",
	NotBetweenSym:  "notbetweensym",
	IsBetweenAsym:  "isbetweenasym",
	NotBetweenAsym: "notbetweenasym",

	IsNull:     "isnull",
	NotNull:    "notnull",
	IsTrue:     "istrue",
	NotTrue:    "nottrue",
	IsFalse:    "isfalse",
	NotFalse:   "notfalse",
	IsUnknown:  "isunknown",
	NotUnknown: "notunknown",
}

// predicateAdjectives is a whitelist of predicate adjectives and adverbs. Used for parsing tags.
var predicateAdjectives = []string{
	"between",
	"betweensym",
	"distinct",
	"false",
	"ilike",
	"in",
	"like",
	"null",
	"similar",
	"true",
	"unknown",
}

// IsBinary reports whether or not the predicate represents a binary comparison.
func (p Predicate) IsBinary() bool { return _binary_pred_start < p && p < _binary_pred_end }

// IsNull reports whether or not the predicate represents a NULL test.
func (p Predicate) IsNull() bool { return _null_pred_start < p && p < _null_pred_end }

// IsBoolean reports whether or not the predicate represents a boolean test.
func (p Predicate) IsBoolean() bool { return _truth_pred_start < p && p < _truth_pred_end }

// IsPatternMatch reports whether or not the predicate represents a pattern-match comparison.
func (p Predicate) IsPatternMatch() bool { return _pattern_pred_start < p && p < _pattern_pred_end }

// IsRange reports whether or not the predicate represents a range comparison.
func (p Predicate) IsRange() bool { return _range_pred_start < p && p < _range_pred_end }

// IsArray reports whether or not the predicate represents an array comparison.
func (p Predicate) IsArray() bool { return _array_pred_start < p && p < _array_pred_end }

// IsUnary reports whether or not the predicate represents a unary comparison.
func (p Predicate) IsUnary() bool { return p.IsNull() || p.IsBoolean() }

// CanQuantify reports whether or not the predicate can be used together with a quantifier.
func (p Predicate) CanQuantify() bool { return p.IsBinary() || p.IsPatternMatch() }

// IsOneOf reports whether or not p is one of the provided predicates.
func (p Predicate) IsOneOf(pp ...Predicate) bool {
	for i := 0; i < len(pp); i++ {
		if pp[i] == p {
			return true
		}
	}
	return false
}

// stringToPredicate is a map of string literals to supported predicates. Used for parsing tags.
var stringToPredicate map[string]Predicate

func init() {
	stringToPredicate = make(map[string]Predicate)
	for p, str := range predicates {
		if len(str) > 0 {
			stringToPredicate[str] = Predicate(p)
		}
	}
}

// LiteralType is a string representation of literal type or type name.
type LiteralType string

const (
	LiteralBool                     LiteralType = "bool"
	LiteralBoolSlice                LiteralType = "[]bool"
	LiteralBoolSliceSlice           LiteralType = "[][]bool"
	LiteralString                   LiteralType = "string"
	LiteralStringSlice              LiteralType = "[]string"
	LiteralStringSliceSlice         LiteralType = "[][]string"
	LiteralStringMap                LiteralType = "map[string]string"
	LiteralStringMapSlice           LiteralType = "[]map[string]string"
	LiteralByte                     LiteralType = "byte"
	LiteralByteSlice                LiteralType = "[]byte"
	LiteralByteSliceSlice           LiteralType = "[][]byte"
	LiteralByteSliceSliceSlice      LiteralType = "[][][]byte"
	LiteralByteArray16              LiteralType = "[16]byte"
	LiteralByteArray16Slice         LiteralType = "[][16]byte"
	LiteralRune                     LiteralType = "rune"
	LiteralRuneSlice                LiteralType = "[]rune"
	LiteralRuneSliceSlice           LiteralType = "[][]rune"
	LiteralInt                      LiteralType = "int"
	LiteralIntSlice                 LiteralType = "[]int"
	LiteralIntSliceSlice            LiteralType = "[][]int"
	LiteralIntArray2                LiteralType = "[2]int"
	LiteralIntArray2Slice           LiteralType = "[][2]int"
	LiteralInt8                     LiteralType = "int8"
	LiteralInt8Slice                LiteralType = "[]int8"
	LiteralInt8SliceSlice           LiteralType = "[][]int8"
	LiteralInt8Array2               LiteralType = "[2]int8"
	LiteralInt8Array2Slice          LiteralType = "[][2]int8"
	LiteralInt16                    LiteralType = "int16"
	LiteralInt16Slice               LiteralType = "[]int16"
	LiteralInt16SliceSlice          LiteralType = "[][]int16"
	LiteralInt16Array2              LiteralType = "[2]int16"
	LiteralInt16Array2Slice         LiteralType = "[][2]int16"
	LiteralInt32                    LiteralType = "int32"
	LiteralInt32Slice               LiteralType = "[]int32"
	LiteralInt32SliceSlice          LiteralType = "[][]int32"
	LiteralInt32Array2              LiteralType = "[2]int32"
	LiteralInt32Array2Slice         LiteralType = "[][2]int32"
	LiteralInt64                    LiteralType = "int64"
	LiteralInt64Slice               LiteralType = "[]int64"
	LiteralInt64SliceSlice          LiteralType = "[][]int64"
	LiteralInt64Array2              LiteralType = "[2]int64"
	LiteralInt64Array2Slice         LiteralType = "[][2]int64"
	LiteralUint                     LiteralType = "uint"
	LiteralUintSlice                LiteralType = "[]uint"
	LiteralUintSliceSlice           LiteralType = "[][]uint"
	LiteralUintArray2               LiteralType = "[2]uint"
	LiteralUintArray2Slice          LiteralType = "[][2]uint"
	LiteralUint8                    LiteralType = "uint8"
	LiteralUint8Slice               LiteralType = "[]uint8"
	LiteralUint8SliceSlice          LiteralType = "[][]uint8"
	LiteralUint8Array2              LiteralType = "[2]uint8"
	LiteralUint8Array2Slice         LiteralType = "[][2]uint8"
	LiteralUint16                   LiteralType = "uint16"
	LiteralUint16Slice              LiteralType = "[]uint16"
	LiteralUint16SliceSlice         LiteralType = "[][]uint16"
	LiteralUint16Array2             LiteralType = "[2]uint16"
	LiteralUint16Array2Slice        LiteralType = "[][2]uint16"
	LiteralUint32                   LiteralType = "uint32"
	LiteralUint32Slice              LiteralType = "[]uint32"
	LiteralUint32SliceSlice         LiteralType = "[][]uint32"
	LiteralUint32Array2             LiteralType = "[2]uint32"
	LiteralUint32Array2Slice        LiteralType = "[][2]uint32"
	LiteralUint64                   LiteralType = "uint64"
	LiteralUint64Slice              LiteralType = "[]uint64"
	LiteralUint64SliceSlice         LiteralType = "[][]uint64"
	LiteralUint64Array2             LiteralType = "[2]uint64"
	LiteralUint64Array2Slice        LiteralType = "[][2]uint64"
	LiteralFloat32                  LiteralType = "float32"
	LiteralFloat32Slice             LiteralType = "[]float32"
	LiteralFloat32SliceSlice        LiteralType = "[][]float32"
	LiteralFloat32Array2            LiteralType = "[2]float32"
	LiteralFloat32Array2Slice       LiteralType = "[][2]float32"
	LiteralFloat64                  LiteralType = "float64"
	LiteralFloat64Slice             LiteralType = "[]float64"
	LiteralFloat64SliceSlice        LiteralType = "[][]float64"
	LiteralFloat64Array2            LiteralType = "[2]float64"
	LiteralFloat64Array2Slice       LiteralType = "[][2]float64"
	LiteralFloat64Array2SliceSlice  LiteralType = "[][][2]float64"
	LiteralFloat64Array2Array2      LiteralType = "[2][2]float64"
	LiteralFloat64Array2Array2Slice LiteralType = "[][2][2]float64"
	LiteralFloat64Array3            LiteralType = "[3]float64"
	LiteralFloat64Array3Slice       LiteralType = "[][3]float64"
	LiteralIP                       LiteralType = "net.IP"
	LiteralIPSlice                  LiteralType = "[]net.IP"
	LiteralIPNet                    LiteralType = "net.IPNet"
	LiteralIPNetSlice               LiteralType = "[]net.IPNet"
	LiteralHardwareAddr             LiteralType = "net.HardwareAddr"
	LiteralHardwareAddrSlice        LiteralType = "[]net.HardwareAddr"
	LiteralTime                     LiteralType = "time.Time"
	LiteralTimeSlice                LiteralType = "[]time.Time"
	LiteralTimeArray2               LiteralType = "[2]time.Time"
	LiteralTimeArray2Slice          LiteralType = "[][2]time.Time"
	LiteralBigInt                   LiteralType = "big.Int"
	LiteralBigIntSlice              LiteralType = "[]big.Int"
	LiteralBigIntArray2             LiteralType = "[2]big.Int"
	LiteralBigIntArray2Slice        LiteralType = "[][2]big.Int"
	LiteralBigFloat                 LiteralType = "big.Float"
	LiteralBigFloatSlice            LiteralType = "[]big.Float"
	LiteralBigFloatArray2           LiteralType = "[2]big.Float"
	LiteralBigFloatArray2Slice      LiteralType = "[][2]big.Float"
	LiteralNullStringMap            LiteralType = "map[string]sql.NullString"
	LiteralNullStringMapSlice       LiteralType = "[]map[string]sql.NullString"
	LiteralStringPtrMap             LiteralType = "map[string]*string"
	LiteralStringPtrMapSlice        LiteralType = "[]map[string]*string"
	LiteralEmptyInterface           LiteralType = "interface{}"
)
