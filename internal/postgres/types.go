package postgres

import (
	"sync"

	"github.com/frk/gosql/internal/analysis"
	"github.com/frk/gosql/internal/postgres/oid"
)

////////////////////////////////////////////////////////////////////////////////
// Result Types
//
type (
	// FieldWrite holds the information needed by the generator to produce the
	// expression nodes that constitute a field-to-column write operation.
	FieldWrite struct {
		// Info on the field from which the column will be written.
		Field *analysis.FieldInfo
		// The column to which the data will be written.
		Column *Column
		// The column identifier.
		ColIdent analysis.ColIdent
		// The name of the valuer to be used for writing the column, or empty.
		Valuer string
	}

	// FieldRead holds the information needed by the generator to produce the
	// expression nodes that constitute a field-from-column read operation.
	FieldRead struct {
		// Info on the field into which the column will be read.
		Field *analysis.FieldInfo
		// The column from which the data will be read.
		Column *Column
		// The column identifier.
		ColIdent analysis.ColIdent
		// The name of the scanner to be used for reading the column, or empty.
		Scanner string
	}

	// FieldFilter ...
	FieldFilter struct {
		// Info on the field that holds the value to be used as the filter parameter.
		Field *analysis.FieldInfo
		// The column which to filter by.
		Column *Column
		// The column identifier.
		ColIdent analysis.ColIdent
		// The name of the valuer to be used for converting the field value.
		Valuer string
	}

	// Boolean
	Boolean struct {
		Value analysis.Boolean
	}

	// FieldConditional
	FieldConditional struct {
		// Name of the field that holds the value to be used in the conditional.
		FieldName string
		// Type info of the field that holds the value to be used in the conditional.
		FieldType analysis.TypeInfo
		// The identifier of the column to be used in the conditional.
		ColIdent analysis.ColIdent
		// The column to be used in the conditional.
		Column *Column
		// The type of the predicate, or 0.
		Predicate analysis.Predicate
		// The predicate quantifier, or 0.
		Quantifier analysis.Quantifier
		// Name of the modifier function, or empty.
		FuncName analysis.FuncName
		// Name of the valuer to be employed, or empty.
		Valuer string
	}

	// ColumnConditional holds the information needed by the generator
	// to produce a column-specific SQL boolean expression.
	ColumnConditional struct {
		// Left hand side column id.
		LHSColIdent analysis.ColIdent
		// Left hand side column.
		LHSColumn *Column
		// Right hand side column id, or empty.
		RHSColIdent analysis.ColIdent
		// Right hand side column, or nil.
		RHSColumn *Column
		// Right hand side literal expression, or empty.
		RHSLiteral string
		// Type of the right hand side column or literal expression.
		RHSType *Type
		// The type of the predicate, or 0.
		Predicate analysis.Predicate
		// The predicate quantifier, or 0.
		Quantifier analysis.Quantifier
	}

	// BetweenConditional
	BetweenConditional struct {
		// The name of the field containing the "between" info.
		FieldName string
		// The id of the predicand column.
		ColIdent analysis.ColIdent
		// The primary column predicand.
		Column *Column
		// The type of the between predicate.
		Predicate analysis.Predicate
		// The lower-bound range predicand.
		LowerBound RangeBound
		// The upper-bound range predicand.
		UpperBound RangeBound
	}

	// NestedConditional
	NestedConditional struct {
		FieldName    string
		Conditionals []WhereConditional
	}

	ConflictInfo struct {
		Target ConflictTarget
		Update []*Column
	}

	ConflictIndex struct {
		// The index predicate.
		Predicate string
		// The index expression.
		Expression string
	}

	ConflictConstraint struct {
		// The name of the constraint.
		Name string
	}

	WhereConditional interface{ whereConditional() }
	RangeBound       interface{ rangeBound() }
	ConflictTarget   interface{ conflictTarget() }

	TableJoinConditional interface {
		tableJoinConditional()
		whereConditional()
	}
)

////////////////////////////////////////////////////////////////////////////////
// PostgreSQL Catalog
//
type (
	// Catalog holds information on various objects of the database.
	Catalog struct {
		// populated by loadCatalog
		Types     map[oid.OID]*Type
		Operators map[OpKey]*Operator
		Casts     map[CastKey]*Cast
		Procs     map[string][]*Proc

		// populated by loadRelation
		Relations map[analysis.RelIdent]*Relation
		// sync.RWMutex is necessary to be used only for the Relations map,
		// the rest is read-only once initialized.
		sync.RWMutex
	}

	// helper type used to map an Operator value
	OpKey struct {
		Name  string
		Left  oid.OID
		Right oid.OID
	}

	// helper type used to map a Cast value
	CastKey struct {
		Target oid.OID
		Source oid.OID
	}
)

////////////////////////////////////////////////////////////////////////////////
// PostgreSQL Catalog Objects
//
type (
	// Relation holds the info of a "pg_class" entry that represents
	// a table, view, or materialized view.
	Relation struct {
		// The object identifier of the relation.
		OID oid.OID
		// The name of the relation.
		Name string
		// The name of the schema to which the relation belongs.
		Schema string
		// The relation's kind, we're only interested in r, v, and m.
		RelKind RelKind
		// List of columns associated with the relation.
		Columns []*Column
		// List of constraints applied to the relation.
		Constraints []*Constraint
		// List of indexes applied to the relation.
		Indexes []*Index
	}

	// Column holds the info of a "pg_attribute" entry that represents
	// a column of a relation.
	Column struct {
		// The number of the column. Ordinary columns are numbered from 1 up.
		Num int16
		// The name of the member's column.
		Name string
		// Records type-specific data supplied at table creation time (for example,
		// the maximum length of a varchar column). It is passed to type-specific
		// input functions and length coercion functions. The value will generally
		// be -1 for types that do not need.
		//
		// NOTE(mkopriva): to get the actual value subtract 4.
		// NOTE(mkopriva): in the case of NUMERIC(precision, scale) types, to
		// calculate the precision use ((typmod - 4) >> 16) & 65535 and to
		// calculate the scale use (typmod - 4) && 65535
		TypeMod int
		// Indicates whether or not the column has a NOT NULL constraint.
		// NOTE(mkopriva): this is ambiguous if the column's relation is a view.
		HasNotNull bool
		// Indicates whether or not the column has a DEFAULT value.
		HasDefault bool
		// Reports whether or not the column is a primary key.
		IsPrimary bool
		// The number of dimensions if the column is an array type, otherwise 0.
		NumDims int
		// The OID of the column's type.
		TypeOID oid.OID
		// Info about the column's type.
		Type *Type
		// The Relation to which the Column belongs.
		Relation *Relation
	}

	// Type holds the info of a "pg_type" entry that represents a column's data type.
	Type struct {
		// The object identifier of the type.
		OID oid.OID
		// The name of the type.
		Name string
		// The formatted name of the type.
		NameFmt string
		// The number of bytes for fixed-size types, negative for variable length types.
		Length int
		// The type's type.
		Type TypeType
		// An arbitrary classification of data types that is used by the parser
		// to determine which implicit casts should be "preferred".
		Category TypeCategory
		// True if the type is a preferred cast target within its category.
		IsPreferred bool
		// If this is an array type then elem identifies the element type
		// of that array type.
		Elem oid.OID
	}

	// Index holds the info of a "pg_index" entry that represents a table's index.
	Index struct {
		// The object identifier of the index.
		OID oid.OID
		// The name of the index.
		Name string
		// The total number of columns in the index; this number includes
		// both key and included attributes.
		NumAtts int
		// If true, this is a unique index.
		IsUnique bool
		// If true, this index represents the primary key of the table.
		IsPrimary bool
		// If true, this index supports an exclusion constraint.
		IsExclusion bool
		// If true, the uniqueness check is enforced immediately on insertion.
		IsImmediate bool
		// If true, the index is currently ready for inserts. False means the
		// index must be ignored by INSERT/UPDATE operations.
		IsReady bool
		// This is an array of values that indicate which table columns this index
		// indexes. For example a value of 1 3 would mean that the first
		// and the third table columns make up the index entries. Key columns come
		// before non-key (included) columns. A zero in this array indicates that
		// the corresponding index attribute is an expression over the table columns,
		// rather than a simple column reference.
		Key []int16
		// The index definition.
		Definition string
		// The index predicate (optional).
		Predicate string
		// Parsed index expression.
		Expression string
	}

	// Constraint holds the info of a "pg_constraint" entry that represents
	// a constraint's on a table.
	Constraint struct {
		// The object identifier of the constraint.
		OID oid.OID
		// Constraint name (not necessarily unique!)
		Name string
		// The type of the constraint
		Type ConstraintType
		// Indicates whether or not the constraint is deferrable
		IsDeferrable bool
		// Indicates whether or not the constraint is deferred by default
		IsDeferred bool
		// If a table constraint (including foreign keys, but not constraint triggers),
		// list of the constrained columns
		Key []int64
		// If a foreign key, list of the referenced columns
		FKey []int64
	}

	// Operator holds info on a "pg_operator" entry.
	Operator struct {
		// The object identifier of the operator.
		OID oid.OID
		// The name of the operator.
		Name string
		// The kind (infix, prefix, or postfix) of the operator.
		Kind string
		// The type oid of the left operand.
		Left oid.OID
		// The type oid of the right operand.
		Right oid.OID
		// The type oid of the result.
		Result oid.OID
	}

	// Cast holds info on a "pg_cast" entry.
	Cast struct {
		// The object identifier of the cast.
		OID oid.OID
		// The oid of the source data type.
		Source oid.OID
		// The oid of the target data type.
		Target oid.OID
		// The context in which the cast can be invoked.
		Context CastContext
	}

	// Proc holds info on a "pg_proc" entry.
	// Current support is limited to functions with 1 input argument
	// and 1 return value, hence ArgType & RetType are single OIDs.
	Proc struct {
		// The object identifier of the procedure.
		OID oid.OID
		// The name of the function.
		Name string
		// The type oid of the function's input argument.
		ArgType oid.OID
		// The type oid of the function's return value.
		RetType oid.OID
		// Indicates whether or not the function is an aggregate function.
		IsAgg bool
	}
)

func (w FieldWrite) NeedsNULLIF() bool {
	return w.Column.IsNULLable() && w.Field.Type.Kind != analysis.TypeKindPtr &&
		!w.Field.Type.ImplementsValuer() && w.Valuer == ""
}

func (r FieldRead) NeedsCOALESCE() bool {
	if r.Field.UseCoalesce {
		return true
	}

	return r.Column.IsNULLable() && r.Field.Type.Kind != analysis.TypeKindPtr &&
		!r.Field.Type.ImplementsScanner() && r.Scanner == ""
}

func (t Type) is(oids ...oid.OID) bool {
	for _, id := range oids {
		if t.OID == id {
			return true
		}
	}
	return false
}

func (t Type) ZeroValueLiteral() (lit string, ok bool) {
	lit, ok = oid.TypeToZeroValue[t.OID]
	return lit, ok
}

func (t *Type) GetNameFmt() string {
	if t == nil {
		return "<nil>"
	}
	return t.NameFmt
}

func (c *Column) IsNULLable() bool {
	return c.Relation != nil && c.Relation.RelKind == RelKindOrdinaryTable &&
		c.HasNotNull == false
}

// TableJoinConditional implementations
func (*Boolean) tableJoinConditional()           {}
func (*ColumnConditional) tableJoinConditional() {}

// WhereConditional implementations
func (*Boolean) whereConditional()            {}
func (*FieldConditional) whereConditional()   {}
func (*ColumnConditional) whereConditional()  {}
func (*BetweenConditional) whereConditional() {}
func (*NestedConditional) whereConditional()  {}

// RangeBound implementations
func (*FieldConditional) rangeBound()  {}
func (*ColumnConditional) rangeBound() {}

// ConflictTarget implementations
func (*ConflictIndex) conflictTarget()      {}
func (*ConflictConstraint) conflictTarget() {}
