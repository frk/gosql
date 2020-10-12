package analysis

import (
	"go/types"

	"github.com/frk/tagutil"
)

////////////////////////////////////////////////////////////////////////////////
// Target Types
//
type (
	// QueryStruct represents the result of the analysis of a target "query" struct type.
	QueryStruct struct {
		// Name of the query struct type.
		TypeName string
		// The queryKind analyzed from query struct's type name.
		Kind QueryKind
		// The primary, `rel`-tagged field of the query struct.
		Rel *RelField
		// Info on the "result" field of the query struct type, or nil.
		Result *ResultField
		// Info on the "where" struct field of the query struct type, or nil.
		Where *WhereStruct
		// Info on the "join", "using", or "from" struct field of the query struct type, or nil.
		Join *JoinStruct
		// Info on the "onConflict" struct field of the query struct type, or nil.
		OnConflict *OnConflictStruct
		// Info on the gosql.OrderBy directive field of the query struct type, or nil.
		OrderBy *OrderByDirective
		// Info on the "limit" field or the gosql.Limit directive of the query struct type, or nil.
		Limit *LimitField
		// Info on the "offset" field or the gosql.Offset directive of the query struct type, or nil.
		Offset *OffsetField
		// Info on the "rowsaffected" field of the query struct type, or nil.
		RowsAffected *RowsAffectedField
		// Info on the gosql.Force directive field of the query struct type, or nil.
		Force *ForceDirective
		// Info on the gosql.Default directive field of the query struct type, or nil.
		Default *DefaultDirective
		// Info on the gosql.Return directive field of the query struct type, or nil.
		Return *ReturnDirective
		// Info on the gosql.Override directive field of the query struct type, or nil.
		Override *OverrideDirective
		// Info on the gosql.ErrorHandler or gosql.ErrorInfoHandler field of the query struct type, or nil.
		ErrorHandler *ErrorHandlerField
		// Info on the gosql.Filter field of the query struct type, or nil.
		Filter *FilterField
		// Info on the context.Context field of the query struct type, or nil.
		Context *ContextField
		// Info on the gosql.All directive field of the query struct type, or nil.
		All *AllDirective
	}

	// FilterStruct represents the result of the analysis of a target "filter" struct type.
	FilterStruct struct {
		// Name of the filter struct type.
		TypeName string
		// The primary, `rel`-tagged field of the filter struct type.
		Rel *RelField
		// Info on the gosql.TextSearch directive field of the filter struct type, or nil.
		TextSearch *TextSearchDirective
		// Info on the field that implements the gosql.FilterConstructor interface.
		FilterConstructor *FilterConstructorField
	}

	// The TargetStruct interface is implemented by the QueryStruct and FilterStruct types.
	TargetStruct interface {
		// GetRelField should return the RelField of the target relation.
		GetRelField() *RelField
		// GetRelIdent should return the RelIdent of the target relation.
		GetRelIdent() RelIdent
		targetStruct()
	}
)

////////////////////////////////////////////////////////////////////////////////
// Rel Fields
//
type (
	// RelField is the primary, `rel`-tagged field of a "query" or "filter" struct type.
	RelField struct {
		// Name of the field.
		FieldName string
		// The relation identifier parsed from the field's `rel` tag.
		Id RelIdent
		// The type information of the field.
		Type RelType
		// Indicates whether or not the gosql.Relation directive was used.
		IsDirective bool
	}

	// ResultField is the result of analyzing a query struct's field named "result" (case insensitive).
	ResultField struct {
		// Name of the field (case preserved).
		FieldName string
		// The type information of the field.
		Type RelType
	}
)

////////////////////////////////////////////////////////////////////////////////
// Rel Type
//
type (
	// RelType holds the type information of the "rel" field of a "query" or "filter" struct type.
	RelType struct {
		// Information on the "rel" field's base type.
		Base TypeInfo
		// Indicates whether or not the base type's a pointer type.
		IsPointer bool
		// Indicates whether or not the base type's a slice type.
		IsSlice bool
		// Indicates whether or not the base type's an array type.
		IsArray bool
		// If the base type's an array type, this field will hold the array's length.
		ArrayLen int64
		// If set, indicates that the type is handled by an iterator.
		IsIter bool
		// If set the value will hold the method name of the iterator interface.
		IterMethod string
		// Indicates whether or not the type implements the gosql.AfterScanner interface.
		IsAfterScanner bool
		// fields holds the information on the type's struct fields.
		Fields []*FieldInfo

		// A map of FieldInfo pointers to source code information about the field.
		// Used only for error reporting, after analysis of the RelType instance
		// this map is merged into the fieldMap of the targetInfo instance.
		//
		// Use `cmp:"-"` to ignore the map's contents during testing since
		// the map as an intermediary container is of secondary importance,
		// plus initializing it for the comparison to succeed would be too
		// much work with little to no return.
		FieldMap map[FieldPtr]FieldVar `cmp:"-"`
	}

	// TypeInfo holds detailed information about a Go type.
	TypeInfo struct {
		// The name of a named type or empty string for unnamed types
		Name string
		// The kind of the go type.
		Kind TypeKind
		// The package import path.
		PkgPath string
		// The package's name.
		PkgName string
		// The local package name (including ".").
		PkgLocal string
		// Indicates whether or not the package is imported.
		IsImported bool
		// Indicates whether or not the type implements the sql.Scanner interface.
		IsScanner bool
		// Indicates whether or not the type implements the driver.Valuer interface.
		IsValuer bool
		// Indicates whether or not the type implements the json.Marshaler interface.
		IsJSONMarshaler bool
		// Indicates whether or not the type implements the json.Unmarshaler interface.
		IsJSONUnmarshaler bool
		// Indicates whether or not the type implements the xml.Marshaler interface.
		IsXMLMarshaler bool
		// Indicates whether or not the type implements the xml.Unmarshaler interface.
		IsXMLUnmarshaler bool
		// Indicates whether or not the type is an empty interface type.
		IsEmptyInterface bool
		// Indicates whether or not the type is the "byte" alias type.
		IsByte bool
		// Indicates whether or not the type is the "rune" alias type.
		IsRune bool
		// If kind is array, ArrayLen will hold the array's length.
		ArrayLen int64
		// If kind is map, key will hold the info on the map's key type.
		Key *TypeInfo
		// If kind is map, elem will hold the info on the map's value type.
		// If kind is ptr, elem will hold the info on pointed-to type.
		// If kind is slice/array, elem will hold the info on slice/array element type.
		Elem *TypeInfo
	}

	// FieldInfo is the result of analyzing a RelType's field and its associated tag.
	FieldInfo struct {
		// Name of the struct field.
		Name string
		// Info about the field's type.
		Type TypeInfo
		// If the field is nested, selector will hold the parent fields' information.
		Selector []*FieldSelectorNode
		// Indicates whether or not the field is embedded.
		IsEmbedded bool
		// Indicates whether or not the field is exported.
		IsExported bool
		// The field's parsed tag.
		Tag tagutil.Tag
		// The column identifier parsed from the field's `sql` tag.
		ColIdent ColIdent
		// If set, indicates the the "nullempty" option was used in the field's `sql` tag.
		NullEmpty bool
		// If set, indicates the the "ro" option was used in the field's `sql` tag.
		ReadOnly bool
		// If set, indicates the the "wo" option was used in the field's `sql` tag.
		WriteOnly bool
		// If set, indicates the the "default" option was used in the field's `sql` tag.
		UseDefault bool
		// If set, indicates the the "add" option was used in the field's `sql` tag.
		UseAdd bool
		// If set, indicates the the "coalesce" option was used in the field's `sql` tag.
		UseCoalesce bool
		// Will hold the "alternative" value as parsed from the "coalesce" option of the field's `sql` tag.
		CoalesceValue string
	}

	// FieldSelectorNode represents a single node in a nested field's "selector".
	// This is a stripped-down version of FieldInfo that holds only that information
	// that is needed by the generator to produce correct Go field selector expressions.
	FieldSelectorNode struct {
		// The name of the field.
		Name string
		// The tag of the field.
		Tag tagutil.Tag
		// The name of the field's type. Empty if the type is unnamed.
		TypeName string
		// The package import path for the field's type. Empty if the type is unnamed.
		TypePkgPath string
		// The name of the package of the field's type. Empty if the type is unnamed.
		TypePkgName string
		// The local name of the imported package of the field's type (including ".").
		// Empty if the type is unnamed.
		TypePkgLocal string
		// Indicates whether or not the type is imported.
		IsImported bool
		// Indicates whether or not the field is embedded.
		IsEmbedded bool
		// Indicates whether or not the field is exported.
		IsExported bool
		// Indicates whether or not the field type is a pointer type.
		IsPointer bool
	}
)

////////////////////////////////////////////////////////////////////////////////
// Where Struct
//
type (
	// WhereStruct represents a struct analyzed from a QueryStruct's
	// field named "where" (case insensitive).
	WhereStruct struct {
		// Name of the field (case preserved).
		FieldName string
		// List of different items specified by the WhereStruct.
		Items []WhereItem
	}

	// WhereBoolTag is the result of analyzing a `bool` tag inside a WhereStruct.
	// The absence of a `bool` key in a WhereStruct's field tag produces an instance
	// of WhereBoolTag with the default value.
	WhereBoolTag struct {
		// The boolean value.
		Value Boolean
	}

	// WhereStructField is the result of analyzing a plain WhereStruct
	// field and its associated `sql` tag.
	WhereStructField struct {
		// The name of the field.
		Name string
		// The field's type information.
		Type TypeInfo
		// The column identifier parsed from the `sql` tag.
		ColIdent ColIdent
		// The predicate type parsed from the `sql` tag, or 0.
		Predicate Predicate
		// The quantifier parsed from the `sql` tag, or 0.
		Quantifier Quantifier
		// The name of the function parsed from the `sql` tag, or empty.
		FuncName FuncName
	}

	// WhereColumnDirective is the result of analyzing a WhereStruct gosql.Column
	// directive and its associated `sql` tag. The content of the `sql` tag is
	// expected to be a predicate expression.
	WhereColumnDirective struct {
		// The LHS column identifier parsed from the `sql` tag.
		LHSColIdent ColIdent
		// The RHS column identifier parsed from the `sql` tag, or empty.
		RHSColIdent ColIdent
		// The RHS literal value parsed from the `sql` tag, or empty.
		RHSLiteral string
		// The predicate type parsed from the `sql` tag, or 0.
		Predicate Predicate
		// The quantifier parsed from the `sql` tag, or 0.
		Quantifier Quantifier
	}

	// WhereBetweenStruct is the result of analyzing a field whose `sql` tag
	// contains one of the "between" predicates.
	WhereBetweenStruct struct {
		// The name of the field with the "between" tag.
		FieldName string
		// The column identifier as extracted from the "between" field. The ColIdent
		// identifies the column to be used as the primary predicand of the BETWEEN predicate.
		ColIdent ColIdent
		// The type of the "between" predicate extracted from the field's tag.
		Predicate Predicate
		// The lower-bound range predicand of the "between" predicate as extracted
		// from the BetweenStruct's sub-field tagged with "lower", or "x".
		LowerBound RangeBound
		// The upper-bound range predicand of the "between" predicate as extracted
		// from the BetweenStruct's sub-field tagged with "upper", or "y".
		UpperBound RangeBound
	}

	// BetweenStructField is the result of analyzing a plain BetweenStruct field.
	BetweenStructField struct {
		// The name of the field.
		Name string
		// The information on the type of the field.
		Type TypeInfo
	}

	// BetweenColumnDirective is the result of analyzing a "_ gosql.Column"
	// directive field in a BetweenStruct.
	BetweenColumnDirective struct {
		// The column identifier parsed from the field's `sql` tag.
		ColIdent ColIdent
	}

	// The WhereItem interface is implemented by the WhereStruct, WhereBoolTag,
	// WhereStructField, WhereColumnDirective, and WhereBetweenStruct types.
	WhereItem interface{ whereItem() }

	// The RangeBound interface is implemented by the BetweenStructField
	// and BetweenColumnDirective types.
	RangeBound interface{ rangeBound() }
)

////////////////////////////////////////////////////////////////////////////////
// Join/From/Using Struct
//
type (
	// JoinStruct represents a struct analyzed from a QueryStruct's field
	// named "join", "from", or "using" (case insensitive).
	JoinStruct struct {
		// Name of the field (case preserved).
		FieldName string
		// Info on the gosql.Relation directive of the "from" and "using"
		// JoinStruct variation, nil for the "join" variation.
		Relation *RelationDirective
		// The list of gosql.JoinXxx directives declared in the JoinStruct.
		Directives []*JoinDirective
	}

	// JoinDirective is the result of analyzing one of the "_ gosql.JoinXxx" directives.
	JoinDirective struct {
		// The type of the join.
		JoinType JoinType
		// The relation identifier parsed from the `sql` tag.
		RelIdent RelIdent
		// List of items parsed from the `sql` tag.
		TagItems []JoinTagItem
	}

	// JoinBoolTagItem is the boolean operator parsed from a join directive's `sql` tag.
	JoinBoolTagItem WhereBoolTag

	// JoinConditionTagItem is the conditional expression parsed from a join directive's `sql` tag.
	JoinConditionTagItem WhereColumnDirective

	// The JoinTagItem interface is implemented by the JoinBoolTagItem
	// and JoinConditionTagItem types.
	JoinTagItem interface{ joinTagItem() }
)

////////////////////////////////////////////////////////////////////////////////
// On Conflict
//
type (
	// OnConflictStruct represents a struct analyzed from a QueryStruct's
	// field named "onconflict" (case insensitive).
	OnConflictStruct struct {
		// Name of the field (case preserved).
		FieldName string
		// If set, indicates that the gosql.Column "conflict_target" directive was used.
		Column *ColumnDirective
		// If set, indicates that the gosql.Index "conflict_target" directive was used.
		Index *IndexDirective
		// If set, indicates that the gosql.Constraint "conflict_target" directive was used.
		Constraint *ConstraintDirective
		// If set, indicates that the gosql.Ignore "conflict_action" directive was used.
		Ignore *IgnoreDirective
		// If set, indicates that the gosql.Update "conflict_action" directive was used.
		Update *UpdateDirective
	}
)

////////////////////////////////////////////////////////////////////////////////
// Fields
//
type (
	// The LimitField is the result of analyzing a query struct's "limit"
	// (case insensitive) field or the "_ gosql.Limit" directive.
	LimitField struct {
		// The name of the field (case preserved), or empty if it's a gosql.Limit directive.
		Name string
		// The value provided in the limit field's / directive's `sql` tag.
		// If the LimitField was produced from a directive the value will be
		// used as a constant.
		//
		// If the LimitField was produced from a normal field the value should *only*
		// be used if the field's actual value is empty, at runtime during the query's
		// execution, essentially acting as a default fallback.
		Value uint64
	}

	// The OffsetField is the result of analyzing a query struct's "offset"
	// (case insensitive) field or the "_ gosql.Offset" directive.
	OffsetField struct {
		// The name of the field (case preserved), or empty if it's a gosql.Offset directive.
		Name string
		// The value provided in the offset field's / directive's `sql` tag.
		// If the OffsetField was produced from a directive the value will be
		// used as a constant.
		//
		// If the OffsetField was produced from a normal field the value should *only*
		// be used if the field's actual value is empty, at runtime during the query's
		// execution, essentially acting as a default fallback.
		Value uint64
	}

	// RowsAffectedField is the result of the analyzing a query struct's
	// "rowsaffected" (case insensitive) field.
	RowsAffectedField struct {
		// Name of the field (case preserved).
		Name string
		// The kind of the field's type.
		TypeKind TypeKind
	}

	// ErrorHandlerField is the result of analyzing a query struct's field whose
	// type implements the gosql.ErrorHandler or gosql.ErrorInfoHandler interface.
	ErrorHandlerField struct {
		// Name of the field (case preserved).
		Name string
		// Indicates whether or not the field's type implements
		// the gosql.ErrorInfoHandler interface.
		IsInfo bool
	}

	// FilterField is the result of analyzing a query struct's gosql.Filter type field.
	FilterField struct {
		// Name of the field (case preserved).
		Name string
	}

	// ContextField is the result of analyzing a query struct's context.Context type field.
	ContextField struct {
		// Name of the field (case preserved).
		Name string
	}

	// FilterConstructorField is the result of analyzing a filter struct's
	// field that implements the gosql.FilterConstructor interface.
	FilterConstructorField struct {
		// Name of the field (case preserved).
		Name string
	}
)

////////////////////////////////////////////////////////////////////////////////
// Directives
//
type (
	// IgnoreDirective is the result of analyzing the "_ gosql.Ignore" directive.
	IgnoreDirective struct {
		// empty
	}

	// AllDirective is the result of analyzing the "_ gosql.All" directive.
	AllDirective struct {
		// empty
	}

	// IndexDirective is the result of analyzing the "_ gosql.Index" directive.
	IndexDirective struct {
		// The name of the index as parsed from the `sql` tag of the directive.
		Name string
	}

	// ConstraintDirective is the result of analyzing the "_ gosql.Constraint" directive.
	ConstraintDirective struct {
		// The name of the constraint as parsed from the `sql` tag of the directive.
		Name string
	}

	// TextSearchDirective is the result of analyzing the "_ gosql.TextSearch" directive.
	TextSearchDirective struct {
		// The column identifier as parsed from the `sql` tag of the directive.
		ColIdent
	}

	// RelationDirective is the result of analyzing the "_ gosql.Relation" directive.
	RelationDirective struct {
		// The relation identifier as parsed from the `rel` tag of the directive.
		RelIdent RelIdent
	}

	// ColumnDirective is the result of analyzing the "_ gosql.Column" directive.
	ColumnDirective struct {
		// The list of column identifiers as parsed from the `sql` tag of the directive.
		ColIdents []ColIdent
	}

	// ForceDirective is the result of analyzing the "_ gosql.Force" directive.
	ForceDirective struct {
		// The list of column identifiers as parsed from the `sql` tag of the directive.
		ColIdentList
	}

	// ReturnDirective is the result of analyzing the "_ gosql.Return" directive.
	ReturnDirective struct {
		// The list of column identifiers as parsed from the `sql` tag of the directive.
		ColIdentList
	}

	// UpdateDirective is the result of analyzing the "_ gosql.Update" directive.
	UpdateDirective struct {
		// The list of column identifiers as parsed from the `sql` tag of the directive.
		ColIdentList
	}

	// DefaultDirective is the result of analyzing the "_ gosql.Default" directive.
	DefaultDirective struct {
		// The list of column identifiers as parsed from the `sql` tag of the directive.
		ColIdentList
	}

	// OverrideDirective is the result of analyzing the "_ gosql.Override" directive.
	OverrideDirective struct {
		// The OVERRIDING value as parsed from the `sql` tag of the directive.
		Kind OverridingKind
	}

	// OrderByDirective is the result of analyzing the "_ gosql.OrderBy" directive.
	OrderByDirective struct {
		// The list of ORDER BY items as parsed from the `sql` tag of the directive.
		Items []OrderByTagItem
	}

	// OrderByTagItem represents a single item parsed from the tag of a "_ gosql.OrderBy" directive.
	OrderByTagItem struct {
		// The column identifier parsed from the `sql` tag.
		ColIdent ColIdent
		// The ordering direction parsed from the `sql` tag.
		Direction OrderDirection
		// The NULL position parsed from the `sql` tag.
		Nulls NullsPosition
	}
)

////////////////////////////////////////////////////////////////////////////////
// Identifiers
//
type (
	// RelIdent represents a relation identifier as parsed from a `rel` or `sql` tag.
	RelIdent struct {
		// The name of the relation.
		Name string
		// The alias of the relation, or empty.
		Alias string
		// The schema qualifier of the relation, or empty.
		Qualifier string
	}

	// ColIdent represents a column identifier as parsed from a field's `sql` tag.
	ColIdent struct {
		// The name of the column.
		Name string
		// The table or table alias qualifier of the column, or empty.
		Qualifier string
	}

	// ColIdentList represents a list of column identifiers as parsed from a field's `sql` tag.
	ColIdentList struct {
		// The individual column identifiers parsed from the `sql` tag.
		Items []ColIdent
		// If set to `true`, indicates that "*" was used in the `sql` tag.
		All bool
	}
)

////////////////////////////////////////////////////////////////////////////////
// Miscellaneous
//
type (
	// FuncName is the name of a database function that can either be used to modify
	// a value, like lower, upper, etc. or a function that can be used as an aggregate.
	FuncName string

	// FieldPtr represents a pointer to the result of an analyzed struct field.
	FieldPtr interface{}

	// FieldVar holds the types.Var represenation and the raw tag of a struct field.
	FieldVar struct {
		// types.Var representation of the struct field.
		Var *types.Var
		// The raw string value of the field's tag.
		Tag string
	}
)

func (id ColIdent) IsEmpty() bool { return id == ColIdent{} }

func (id ColIdent) String() string {
	if len(id.Qualifier) > 0 {
		return id.Qualifier + "." + id.Name
	}
	return id.Name
}

func (s *QueryStruct) IsSelectCountOrExists() bool {
	return s.Kind.isSelect() && s.Kind != QueryKindSelect
}

func (s *QueryStruct) IsUpdateSlice() bool {
	return s.Kind == QueryKindUpdate && s.Rel.Type.IsSlice
}

func (s *QueryStruct) IsUpdateWithPKeys() bool {
	return s.Kind == QueryKindUpdate && s.HasNoQualifier()
}

func (s *QueryStruct) IsUpdateWithoutPKeys() bool {
	return s.Kind == QueryKindUpdate &&
		(s.Rel.Type.IsSingle()) &&
		(s.Where != nil || s.Filter != nil || s.All != nil)
}

func (s *QueryStruct) IsInsertOrUpdate() bool {
	return s.Kind == QueryKindInsert || s.Kind == QueryKindUpdate
}

func (s *QueryStruct) IsInsertOrUpdateSlice() bool {
	return s.IsInsertOrUpdate() && s.Rel.Type.IsSlice
}

func (s *QueryStruct) IsWithoutOutput() bool {
	return s.Result == nil && s.Return == nil && !s.Kind.isSelect()
}

func (s *QueryStruct) IsSingleOutput() bool {
	return (s.Kind.isSelect() && s.Rel.Type.IsSingle()) ||
		(s.Result != nil && s.Result.Type.IsSingle()) ||
		(s.Return != nil && s.Rel.Type.IsSingle())
}

func (s *QueryStruct) HasNoErrorInfoHandler() bool {
	return s.ErrorHandler == nil || s.ErrorHandler.IsInfo == false
}

func (s *QueryStruct) HasNoQualifier() bool {
	return s.Where == nil && s.Filter == nil && s.All == nil
}

func (s *QueryStruct) HasQualifier() bool {
	return s.Where != nil || s.Filter != nil || s.All != nil
}

func (s *QueryStruct) OutputRelFieldName() string {
	if s.Result != nil {
		return s.Result.FieldName
	}
	return s.Rel.FieldName
}

// OutputRelType returns the RelType representation of the QueryStruct's output.
func (s *QueryStruct) OutputRelType() RelType {
	if s.Result != nil {
		return s.Result.Type
	}
	return s.Rel.Type
}

// OutputIsAfterScanner reports whether or not the QueryStruct's
// output type implements the AfterScanner interface.
func (s *QueryStruct) OutputIsAfterScanner() bool {
	return s.OutputRelType().IsAfterScanner
}

// IsSingle reports whether or not the RelType represents a non-collection type.
func (t *RelType) IsSingle() bool {
	return !t.IsSlice && !t.IsIter && !t.IsArray
}

// HasFieldWithColumn reports whether or not the RelType has a field with a
// column identifier whose name matches the given colName.
func (t *RelType) HasFieldWithColumn(colName string) bool {
	for _, f := range t.Fields {
		if f.ColIdent.Name == colName {
			return true
		}
	}
	return false
}

func (id *RelIdent) QualifiedName() string {
	if len(id.Qualifier) > 0 {
		return id.Qualifier + "." + id.Name
	}
	return id.Name
}

// GetRelField (implements TargetStruct) returns the RelField of the QueryStruct.
func (s *QueryStruct) GetRelField() *RelField { return s.Rel }

// GetRelField (implements TargetStruct) returns the RelField of the FilterStruct.
func (s *FilterStruct) GetRelField() *RelField { return s.Rel }

// GetRelIdent (implements TargetStruct) returns the RelIdent of the QueryStruct.
func (s *QueryStruct) GetRelIdent() RelIdent { return s.Rel.Id }

// GetRelIdent (implements TargetStruct) returns the RelIdent of the FilterStruct.
func (s *FilterStruct) GetRelIdent() RelIdent { return s.Rel.Id }

// TargetStruct implementations
func (*QueryStruct) targetStruct()  {}
func (*FilterStruct) targetStruct() {}

// whereItem implementations
func (*WhereStruct) whereItem()          {}
func (*WhereBoolTag) whereItem()         {}
func (*WhereStructField) whereItem()     {}
func (*WhereColumnDirective) whereItem() {}
func (*WhereBetweenStruct) whereItem()   {}

// rangeBound implementations
func (*BetweenStructField) rangeBound()     {}
func (*BetweenColumnDirective) rangeBound() {}

// joinTagItem implementations
func (*JoinBoolTagItem) joinTagItem()      {}
func (*JoinConditionTagItem) joinTagItem() {}
