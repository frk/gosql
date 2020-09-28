// WARNING(mkopriva): IF NOT NECESSARY TRY NOT TO DO ANYTHING THAT WILL REORDER
// LINES OR ADD NEW LINES IN THE MIDDLE OF THE FILE, ONLY ADD NEW CODE AT THE TAIL END.
package testdata

import (
	"context"

	"github.com/frk/gosql"
	"github.com/frk/gosql/internal/testdata/common"
)

//BAD: missing relation field
type InsertAnalysisTestBAD_NoDataField struct {
	User *common.User
}

//BAD: invalid datatype kind
type InsertAnalysisTestBAD3 struct {
	User string `rel:"users_table"`
}

//BAD: Delete with invalid relid
type DeleteAnalysisTestBAD_BadRelId struct {
	Rel T `rel:"foo.123:bar"`
}

//BAD: Select with multiple rel tags
type SelectAnalysisTestBAD_MultipleRelTags struct {
	Rel1 T `rel:"relation_a:a"`
	Rel2 T `rel:""`
}

//BAD: Delete with illegal count field
type DeleteAnalysisTestBAD_IllegalCountField struct {
	Count int `rel:"relation_a:a"`
}

//BAD: Update with illegal exists field
type UpdateAnalysisTestBAD_IllegalExistsField struct {
	Exists bool `rel:"relation_a:a"`
}

//BAD: Insert with illegal notexists field
type InsertAnalysisTestBAD_IllegalNotExistsField struct {
	NotExists bool `rel:"relation_a:a"`
}

//BAD: Select with illegal gosql.Relation directive
type SelectAnalysisTestBAD_IllegalRelationDirective struct {
	_ gosql.Relation `rel:"relation_a:a"`
}

//BAD: Select with unnamed base struct type
type SelectAnalysisTestBAD_UnnamedBaseStructType struct {
	Rel []*struct{} `rel:"relation_a:a"`
}

//BAD: Select with All directive
type SelectAnalysisTestBAD_IllegalAllDirective struct {
	Rel []T `rel:"relation_a:a"`
	_   gosql.All
}

//BAD: Insert with All directive
type InsertAnalysisTestBAD_IllegalAllDirective struct {
	Rel T `rel:"relation_a:a"`
	_   gosql.All
}

//BAD: Update with conflicting where producer
type UpdateAnalysisTestBAD_ConflictWhereProducer struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		Id int `sql:"a.id"`
	}
	_ gosql.All
}

//BAD: Delete with illegal gosql.Default directive
type DeleteAnalysisTestBAD_IllegalDefaultDirective struct {
	Rel T             `rel:"relation_a:a"`
	_   gosql.Default `sql:"*"`
}

//BAD: Update with empty gosql.Default directive collist
type UpdateAnalysisTestBAD_EmptyDefaultDirectiveCollist struct {
	Rel T `rel:"relation_a:a"`
	_   gosql.Default
}

//BAD: Select with illegal gosql.Force directive
type SelectAnalysisTestBAD_IllegalForceDirective struct {
	Rel T           `rel:"relation_a:a"`
	_   gosql.Force `sql:"*"`
}

//BAD: Update with bad gosql.Force directive colid
type UpdateAnalysisTestBAD_BadForceDirectiveColId struct {
	Rel T           `rel:"relation_a:a"`
	_   gosql.Force `sql:"a.id,1234"`
}

//BAD: Filter with illegal gosql.Return directive
type FilterAnalysisTestBAD_IllegalReturnDirective struct {
	Rel T            `rel:"relation_a:a"`
	_   gosql.Return `sql:"*"`
}

//BAD: Delete with conflicting result producer
type DeleteAnalysisTestBAD_ConflictResultProducer struct {
	Rel T            `rel:"relation_a:a"`
	_   gosql.Return `sql:"*"`
	_   gosql.Return `sql:"a.id"`
}

//BAD: Update with empty gosql.Return directive collist
type UpdateAnalysisTestBAD_EmptyReturnDirectiveCollist struct {
	Rel T `rel:"relation_a:a"`
	_   gosql.Return
}

//BAD: Insert with Limit field
type InsertAnalysisTestBAD_IllegalLimitField struct {
	Rel T           `rel:"relation_a:a"`
	_   gosql.Limit `sql:"10"`
}

//BAD: Update with Offset field
type UpdateAnalysisTestBAD_IllegalOffsetField struct {
	Rel T            `rel:"relation_a:a"`
	_   gosql.Offset `sql:"2"`
}

//BAD: Insert with illegal gosql.OrderBy directive
type InsertAnalysisTestBAD_IllegalOrderByDirective struct {
	Rel T             `rel:"relation_a:a"`
	_   gosql.OrderBy `sql:"a.id"`
}

//BAD: Delete with illegal gosql.Override directive
type DeleteAnalysisTestBAD_IllegalOverrideDirective struct {
	Rel T              `rel:"relation_a:a"`
	_   gosql.Override `sql:"user"`
}

//BAD: Select with illegal gosql.TextSearch directive
type SelectAnalysisTestBAD_IllegalTextSearchDirective struct {
	Rel T                `rel:"relation_a:a"`
	_   gosql.TextSearch `sql:"a._document"`
}

//BAD: Select with illegal gosql.Column directive
type SelectAnalysisTestBAD_IllegalColumnDirective struct {
	Rel T            `rel:"relation_a:a"`
	_   gosql.Column `sql:"a.id"`
}

//BAD: Insert with illegal Where block
type InsertAnalysisTestBAD_IllegalWhereBlock struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		Id int `sql:"id"`
	}
}

//BAD: Update with illegal Join block
type UpdateAnalysisTestBAD_IllegalJoinBlock struct {
	Rel  T `rel:"relation_a:a"`
	Join struct {
		_ gosql.Relation
	}
}

//BAD: Delete with illegal From block
type DeleteAnalysisTestBAD_IllegalFromBlock struct {
	Rel  T `rel:"relation_a:a"`
	From struct {
		_ gosql.Relation
	}
}

//BAD: Select with illegal Using block
type SelectAnalysisTestBAD_IllegalUsingBlock struct {
	Rel   T `rel:"relation_a:a"`
	Using struct {
		_ gosql.Relation
	}
}

//BAD: Update with illegal OnConflict block
type UpdateAnalysisTestBAD_IllegalOnConflictBlock struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict struct {
		// ...
	}
}

//BAD: Select with illegal result field
type SelectAnalysisTestBAD_IllegalResultField struct {
	Rel    T `rel:"relation_a:a"`
	Result T
}

//BAD: Select with conflicting limit producers
type SelectAnalysisTestBAD_ConflictLimitProducer struct {
	Rel   []T         `rel:"relation_a:a"`
	_     gosql.Limit `sql:"10"`
	Limit int
}

//BAD: Select with conflicting offset producers
type SelectAnalysisTestBAD_ConflictOffsetProducer struct {
	Rel    []T `rel:"relation_a:a"`
	Offset int
	_      gosql.Offset `sql:"2"`
}

//BAD: Select with illegal rowsaffected field
type SelectAnalysisTestBAD_IllegalRowsAffectedField struct {
	Rel          []T `rel:"relation_a:a"`
	RowsAffected int
}

//BAD: Insert with illegal filter field
type InsertAnalysisTestBAD_IllegalFilterField struct {
	Rel []T `rel:"relation_a:a"`
	F   gosql.Filter
}

//BAD: Select with conflicting where producer
type SelectAnalysisTestBAD_ConflictWhereProducer struct {
	Rel   []T `rel:"relation_a:a"`
	Where struct {
		Id int `sql:"id"`
	}
	F gosql.Filter
}

//BAD: Delete with conflicting error handlers
type DeleteAnalysisTestBAD_ConflictErrorHandler struct {
	Rel          T `rel:"relation_a:a"`
	ErrorHandler myerrorhandler
	erh          myerrorhandler
}

type badIterator interface { // too many methods
	Fn1(*common.User) error
	Fn2(*common.User) error
}

type badIterator2 interface { // bad signature
	Fn(*common.User) int
}

//BAD: Select with iterator with too many methods
type SelectAnalysisTestBAD_IteratorWithTooManyMethods struct {
	Rel badIterator `rel:"relation_a:a"`
}

//BAD: Select with iterator with bad signature
type SelectAnalysisTestBAD_IteratorWithBadSignature struct {
	Rel func(*common.User) int `rel:"relation_a:a"`
}

//BAD: Select with iterator with bad signature (interface)
type SelectAnalysisTestBAD_IteratorWithBadSignatureIface struct {
	Rel badIterator2 `rel:"relation_a:a"`
}

//BAD: Select with imported iterator that has unexported method
type SelectAnalysisTestBAD_IteratorWithUnexportedMethod struct {
	Rel common.BadIterator `rel:"relation_a:a"`
}

//BAD: Select with iterator with unnamed argument
type SelectAnalysisTestBAD_IteratorWithUnnamedArgument struct {
	Rel func(*struct{}) error `rel:"relation_a:a"`
}

//BAD: Select with iterator with non-struct argument
type SelectAnalysisTestBAD_IteratorWithNonStructArgument struct {
	Rel func(*notstruct) error `rel:"relation_a:a"`
}

type notstruct string

//BAD: Insert with bad struct base type
type InsertAnalysisTestBAD_BadRelfiedlStructBaseType struct {
	Rel []*notstruct `rel:"relation_a:a"`
}

//BAD: Update with bad dataType field's colid
type UpdateAnalysisTestBAD_BadRelTypeFieldColId struct {
	Rel struct {
		Foo string `sql:"1234"`
	} `rel:"relation_a:a"`
}

//BAD: Update with conflicting where produceer
type UpdateAnalysisTestBAD_ConflictWhereProducer2 struct {
	Rel   T `rel:"relation_a:a"`
	_     gosql.All
	Where struct {
		Id int `sql:"id"`
	}
}

//BAD: Delete with bad where block type
type DeleteAnalysisTestBAD_BadWhereBlockType struct {
	Rel   T `rel:"relation_a:a"`
	Where []string
}

//BAD: Select with bad bool tag value
type SelectAnalysisTestBAD_BadBoolTagValue struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		Id   int    `sql:"id"`
		Name string `sql:"name" bool:"abc"`
	}
}

//BAD: Select with bad nested where block type
type SelectAnalysisTestBAD_BadNestedWhereBlockType struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		Id int       `sql:"id"`
		X  notstruct `sql:">"`
	}
}

//BAD: Select with bad gosql.Column expression LHS
type SelectAnalysisTestBAD_BadColumnExpressionLHS struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		_ gosql.Column `sql:"123 = x"`
	}
}

//BAD: Select with bad gosql.Column predicate combo
type SelectAnalysisTestBAD_BadColumnPredicateCombo struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		_ gosql.Column `sql:"x isin any y"`
	}
}

//BAD: Delete with bad gosql.Column expression LHS
type DeleteAnalysisTestBAD_BadColumnExpressionLHS struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		_ gosql.Column `sql:"123 isnull"`
	}
}

//BAD: Update with bad unary op
type UpdateAnalysisTestBAD_BadUnaryOp struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		_ gosql.Column `sql:"x <="`
	}
}

//BAD: Update with extra quantifier
type UpdateAnalysisTestBAD_ExtraQuantifier struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		_ gosql.Column `sql:"x isnull any"`
	}
}

//BAD: Select with bad between field type
type SelectAnalysisTestBAD_BadBetweenFieldType struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		between notstruct `sql:"a.foo isbetween"`
	}
}

//BAD: Select with bad number of fields in "between" struct
type SelectAnalysisTestBAD_BadBetweenFieldType2 struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		between struct {
			x, y, z int
		} `sql:"a.foo isbetween"`
	}
}

//BAD: Select with bad colid in "between" struct field's tag
type SelectAnalysisTestBAD_BadBetweenArgColId struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		between struct {
			_ gosql.Column `sql:"a.bar,x"`
			_ gosql.Column `sql:"123,y"`
		} `sql:"a.foo isbetween"`
	}
}

//BAD: Select with missing x / y in "between"
type SelectAnalysisTestBAD_NoBetweenXYArg struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		between struct {
			_ gosql.Column `sql:"a.bar"`
			_ gosql.Column `sql:"a.baz,y"`
		} `sql:"a.foo isbetween"`
	}
}

//BAD: Select with bad "between" target colid
type SelectAnalysisTestBAD_BadBetweenColId struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		between struct {
			_ gosql.Column `sql:"a.bar,x"`
			_ gosql.Column `sql:"a.baz,y"`
		} `sql:"123 isbetween"`
	}
}

//BAD: Delete with bad where field colid
type DeleteAnalysisTestBAD_BadWhereFieldColId struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		Id int `sql:"123"`
	}
}

//BAD: Delete with bad where field predicate combo
type DeleteAnalysisTestBAD_BadWhereFieldPredicateCombo struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		Id []int `sql:"a.id notin any"`
	}
}

//BAD: Delete with illegal where field unary comparison
type DeleteAnalysisTestBAD_IllegalWhereFieldUnaryPredicate struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		Id int `sql:"a.id istrue"`
	}
}

//BAD: Update with bad where field type for quantifier
type UpdateAnalysisTestBAD_BadWhereFieldTypeForQuantifier struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		Id int `sql:"a.id = any"`
	}
}

//BAD: Select with bad join block type
type SelectAnalysisTestBAD_BadJoinBlockType struct {
	Rel  T `rel:"relation_a:a"`
	Join notstruct
}

//BAD: Select with illegal join gosql.Relation directive
type SelectAnalysisTestBAD_IllegalJoinBlockRelationDirective struct {
	Rel  T `rel:"relation_a:a"`
	Join struct {
		_ gosql.Relation `sql:"foobar"`
	}
}

//BAD: Delete with conflicting using block gosql.Relation directive
type DeleteAnalysisTestBAD_ConflictRelationDirective struct {
	Rel   T `rel:"relation_a:a"`
	Using struct {
		_ gosql.Relation `sql:"foo"`
		_ gosql.Relation `sql:"bar"`
	}
}

//BAD: Update with bad from block gosql.Relation directive relid
type UpdateAnalysisTestBAD_BadFromRelationRelId struct {
	Rel  T `rel:"relation_a:a"`
	From struct {
		_ gosql.Relation `sql:"123"`
	}
}

//BAD: Select with bad join block gosql.JoinXxx directive relid
type SelectAnalysisTestBAD_BadJoinDirectiveRelId struct {
	Rel  T `rel:"relation_a:a"`
	Join struct {
		_ gosql.LeftJoin `sql:"123"`
	}
}

//BAD: Select with bad gosql.JoinXxx directive expression colid
type SelectAnalysisTestBAD_BadJoinDirectiveExpressionColId struct {
	Rel  T `rel:"relation_a:a"`
	Join struct {
		_ gosql.LeftJoin `sql:"relation_b:b,123 = b.foo"`
	}
}

//BAD: Select with bad gosql.JoinXxx directive expression predicate
type SelectAnalysisTestBAD_BadJoinDirectiveExpressionPredicate struct {
	Rel  T `rel:"relation_a:a"`
	Join struct {
		_ gosql.LeftJoin `sql:"relation_b:b,b.foo ="`
	}
}

//BAD: Select with bad gosql.JoinXxx directive expression extra quantifier
type SelectAnalysisTestBAD_BadJoinDirectiveExpressionExtraQuantifier struct {
	Rel  T `rel:"relation_a:a"`
	Join struct {
		_ gosql.LeftJoin `sql:"relation_b:b,b.foo isnull any"`
	}
}

//BAD: Select with bad gosql.JoinXxx directive expression predicate combo
type SelectAnalysisTestBAD_BadJoinDirectiveExpressionPredicateCombo struct {
	Rel  T `rel:"relation_a:a"`
	Join struct {
		_ gosql.LeftJoin `sql:"relation_b:b,b.foo isin any a.bar"`
	}
}

//BAD: Delete with illegal join block directive
type DeleteAnalysisTestBAD_IllegalJoinBlockDirective struct {
	Rel   T `rel:"relation_a:a"`
	Using struct {
		_ gosql.Column `sql:"a.foo"`
	}
}

//BAD: Insert with bad onconflict block type
type InsertAnalysisTestBAD_BadOnConflictBlockType struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict notstruct
}

//BAD: Insert with conflicting onconflict block target
type InsertAnalysisTestBAD_ConflictOnConflictBlockTargetProducer struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.Index  `sql:"some_index"`
		_ gosql.Column `sql:"a.id"`
	}
}

//BAD: Insert with conflicting onconflict block target
type InsertAnalysisTestBAD_ConflictOnConflictBlockTargetProducer2 struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.Column `sql:"a.id"`
		_ gosql.Index  `sql:"some_index"`
	}
}

//BAD: Insert with conflicting onconflict block target
type InsertAnalysisTestBAD_ConflictOnConflictBlockTargetProducer3 struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.Index      `sql:"some_index"`
		_ gosql.Constraint `sql:"some_constraint"`
	}
}

//BAD: Insert with conflicting onconflict block action
type InsertAnalysisTestBAD_ConflictOnConflictBlockActionProducer struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.Column `sql:"a.id"`
		_ gosql.Update `sql:"a.foo"`
		_ gosql.Ignore
	}
}

//BAD: Insert with conflicting onconflict block action
type InsertAnalysisTestBAD_ConflictOnConflictBlockActionProducer2 struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.Column `sql:"a.id"`
		_ gosql.Ignore
		_ gosql.Update `sql:"a.foo"`
	}
}

//BAD: Insert with bad onconflict column target value
type InsertAnalysisTestBAD_BadOnConflictColumnTargetValue struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.Column `sql:"a.id,a.1234"`
	}
}

//BAD: Insert with bad onconflict index target identifier
type InsertAnalysisTestBAD_BadOnConflictIndexTargetIdent struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.Index `sql:"1234"`
	}
}

//BAD: Insert with bad onconflict constraint target identifier
type InsertAnalysisTestBAD_BadOnConflictConstraintTargetIdent struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.Constraint `sql:"1234"`
	}
}

//BAD: Insert with bad onconflict update action collist
type InsertAnalysisTestBAD_BadOnConflictUpdateActionCollist struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.Column `sql:"a.id"`
		_ gosql.Update `sql:"a.id,a.1234"`
	}
}

//BAD: Insert with illegal onconflict directive
type InsertAnalysisTestBAD_IllegalOnConflictDirective struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.LeftJoin `sql:"relation_b:b,b.id=a.id"`
	}
}

//BAD: Insert with illegal onconflict directive
type InsertAnalysisTestBAD_NoOnConflictTarget struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.Update `sql:"a.foo,a.bar"`
	}
}

//BAD: Select with bad limit field type
type SelectAnalysisTestBAD_BadLimitFieldType struct {
	Rel   []T    `rel:"relation_a:a"`
	Limit string `sql:"123"`
}

//BAD: Select with no limit directive value
type SelectAnalysisTestBAD_NoLimitDirectiveValue struct {
	Rel []T         `rel:"relation_a:a"`
	_   gosql.Limit `sql:""`
}

//BAD: Select with bad limit directive value
type SelectAnalysisTestBAD_BadLimitDirectiveValue struct {
	Rel []T         `rel:"relation_a:a"`
	_   gosql.Limit `sql:"abc"`
}

//BAD: Select with bad offset field type
type SelectAnalysisTestBAD_BadOffsetFieldType struct {
	Rel    []T    `rel:"relation_a:a"`
	Offset string `sql:"123"`
}

//BAD: Select with no offset directive value
type SelectAnalysisTestBAD_NoOffsetDirectiveValue struct {
	Rel []T          `rel:"relation_a:a"`
	_   gosql.Offset `sql:""`
}

//BAD: Select with bad offset directive value
type SelectAnalysisTestBAD_BadOffsetDirectiveValue struct {
	Rel []T          `rel:"relation_a:a"`
	_   gosql.Offset `sql:"abc"`
}

//BAD: Select with empty gosql.OrderBy directive list
type SelectAnalysisTestBAD_EmptyOrderByDirectiveCollist struct {
	Rel []T `rel:"relation_a:a"`
	_   gosql.OrderBy
}

//BAD: Select with bad gosql.OrderBy directive nulls order option value
type SelectAnalysisTestBAD_BadOrderByDirectiveNullsOrderValue struct {
	Rel []T           `rel:"relation_a:a"`
	_   gosql.OrderBy `sql:"a.id:nullsthird"`
}

//BAD: Select with bad gosql.OrderBy directive colid
type SelectAnalysisTestBAD_BadOrderByDirectiveCollist struct {
	Rel []T           `rel:"relation_a:a"`
	_   gosql.OrderBy `sql:"-a.id:nullsfirst,a.1234"`
}

//BAD: Insert with bad gosql.Override directive kind value
type InsertAnalysisTestBAD_BadOverrideDirectiveKindValue struct {
	Rel []T            `rel:"relation_a:a"`
	_   gosql.Override `sql:"foo"`
}

//BAD: Update with conflicting result producer
type UpdateAnalysisTestBAD_ConflictResultProducer struct {
	Rel    T            `rel:"relation_a:a"`
	_      gosql.Return `sql:"*"`
	Result []T
}

//BAD: Update with bad result field type
type UpdateAnalysisTestBAD_BadResultFieldType struct {
	Rel    T `rel:"relation_a:a"`
	Result []notstruct
}

//BAD: Delete with conflicting result producer
type DeleteAnalysisTestBAD_ConflictResultProducer2 struct {
	Rel          T `rel:"relation_a:a"`
	Result       []T
	RowsAffected int
}

//BAD: Delete with bad rowsaffected field type
type DeleteAnalysisTestBAD_BadRowsAffecteFieldType struct {
	Rel          T `rel:"relation_a:a"`
	RowsAffected string
}

//BAD: Filter with bad gosql.TextSearch directive colid
type FilterAnalysisTestBAD_BadTextSearchDirectiveColId struct {
	Rel T                `rel:"relation_a:a"`
	_   gosql.TextSearch `sql:"123"`
}

//BAD: Update slice with All directive
type UpdateAnalysisTestBAD_IllegalAllDirective struct {
	Rel []T `rel:"relation_a:a"`
	_   gosql.All
}

//BAD: Update slice with Where struct
type UpdateAnalysisTestBAD_IllegalWhereStruct struct {
	Rel   []T `rel:"relation_a:a"`
	Where struct {
		Name string `sql:"name"`
	}
}

//BAD: Update slice with Filter field
type UpdateAnalysisTestBAD_IllegalFilterField struct {
	Rel []T `rel:"relation_a:a"`
	F   gosql.Filter
}

//BAD: Delete with illegal unary predicate in expression
type DeleteAnalysisTestBAD_IllegalUnaryPredicateInExpression struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		_ gosql.Column `sql:"a.id isfalse a.foo"`
	}
}

//BAD: Select with illegal unary predicate in gosql.JoinXxx directive expression
type SelectAnalysisTestBAD_IllegalUnaryPredicateInJoinDirectiveExpression struct {
	Rel  T `rel:"relation_a:a"`
	Join struct {
		_ gosql.LeftJoin `sql:"relation_b:b,b.foo istrue a.bar"`
	}
}

//BAD: where block with list predicate on field of non-sequence type
type DeleteAnalysisTestBAD_ListPredicate struct {
	Rel   T `rel:"relation_a"`
	Where struct {
		a int `sql:"column_a isin"`
	}
}

//BAD: DELETE with conflicting relation name
type DeleteAnalysisTestBAD_ConflictingRelationName struct {
	Rel   T `rel:"relation_a"`
	Using struct {
		_ gosql.Relation `sql:"relation_d:d"`
		_ gosql.LeftJoin `sql:"relation_a,id = d.a_id"`
	}
	Where struct {
		_ gosql.Column `sql:"a.id = d.a_id"`
	}
}

//BAD: DELETE with conflicting relation alias
type DeleteAnalysisTestBAD_ConflictingRelationAlias struct {
	Rel   T `rel:"relation_a:a"`
	Using struct {
		_ gosql.Relation `sql:"relation_b:b"`
		_ gosql.LeftJoin `sql:"relation_c:a,a.id = b.c_id"`
	}
	Where struct {
		_ gosql.Column `sql:"a.id = b.a_id"`
	}
}

//BAD: Filter with conflicting rel tag
type FilterAnalysisTestBAD_ConflictingRelTag struct {
	_ T                `rel:"relation_a:a"`
	_ gosql.TextSearch `rel:"a.ts_document"`
}

//BAD: Filter with illegal iterator type for rel field
type FilterAnalysisTestBAD_IllegalIteratorType struct {
	_ func(*T) error   `rel:"relation_a:a"`
	_ gosql.TextSearch `sql:"a.ts_document"`
}

//BAD: Select with conflicting relation name
type SelectAnalysisTestBAD_ConflictingRelName struct {
	Join struct {
		_ gosql.LeftJoin `sql:"relation_a,relation_a.foo istrue"`
	}
	Rel T `rel:"relation_a"`
}

//BAD: Select with conflicting relation alias
type SelectAnalysisTestBAD_ConflictingRelAlias struct {
	Join struct {
		_ gosql.LeftJoin `sql:"relation_b:a,a.foo istrue"`
	}
	Rel T `rel:"relation_a:a"`
}

//BAD: Update with conflicting relation name
type UpdateAnalysisTestBAD_ConflictingRelationName struct {
	Rel  T `rel:"relation_a"`
	From struct {
		_ gosql.Relation `sql:"relation_a"`
		_ gosql.LeftJoin `sql:"relation_d:d,d.id = relation_a.d_id"`
	}
	Where struct {
		Value string `sql:"d.col_value"`
	}
}

//BAD: Update with conflicting relation alias
type UpdateAnalysisTestBAD_ConflictingRelationAlias struct {
	Rel  T `rel:"relation_a:a"`
	From struct {
		_ gosql.Relation `sql:"relation_b:a"`
		_ gosql.LeftJoin `sql:"relation_c:c,c.id = a.c_id"`
	}
	Where struct {
		Value string `sql:"c.col_value"`
	}
}

//BAD: Select with join condition containing unknown qualifier
type SelectAnalysisTestBAD_UnknownColumnQualifier struct {
	Rel  T `rel:"relation_a:a"`
	Join struct {
		_ gosql.LeftJoin `sql:"relation_b:b,b.id = c.b_id"`
	}
}

//BAD: Select with join condition containing unknown qualifier
type SelectAnalysisTestBAD_UnknownColumnQualifier2 struct {
	Rel  T `rel:"relation_a"`
	Join struct {
		_ gosql.LeftJoin `sql:"relation_b,relation_b.id = relation_c.b_id"`
	}
}

//BAD: textsearch relation not found (bad alias)
type FilterAnalysisTestBAD_UnknownColumnQualifierInTextSearch struct {
	_ CT1              `rel:"column_tests_1:c"`
	_ gosql.TextSearch `sql:"x.col_b"`
}

//BAD: Unknown column qualifier in gosql.Return directive
type UpdateAnalysisTestBAD_UnknownColumnQualifierInReturn struct {
	Rel CT1          `rel:"column_tests_1:c"`
	_   gosql.Return `sql:"x.col_a"`
}

//BAD: Unknown column qualifier in gosql.LeftJoin directive
type SelectAnalysisTestBAD_UnknownColumnQualifierInJoin struct {
	Columns CT1 `rel:"column_tests_1:a"`
	Join    struct {
		_ gosql.LeftJoin `sql:"column_tests_2:b,x.col_foo = a.col_a"`
	}
}

//BAD: Unknown column qualifier in gosql.LeftJoin directive (2)
type SelectAnalysisTestBAD_UnknownColumnQualifierInJoin2 struct {
	Columns CT1 `rel:"column_tests_1:a"`
	Join    struct {
		_ gosql.LeftJoin `sql:"column_tests_2:b,b.col_foo = x.col_a"`
	}
}

//BAD: Unknown column qualifier in gosql.Force directive
type InsertAnalysisTestBAD_UnknownColumnQualifierInForce struct {
	Rel CT1         `rel:"column_tests_1:c"`
	_   gosql.Force `sql:"x.col_a"`
}

//BAD: Unknown column qualifier in Where field
type SelectAnalysisTestBAD_UnknownColumnQualifierInWhereField struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		Id int `sql:"x.id"`
	}
}

//BAD: Unknown column qualifier in Where gosql.Column directive
type SelectAnalysisTestBAD_UnknownColumnQualifierInWhereColumn struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		_ gosql.Column `sql:"x.col_a = 123"`
	}
}

//BAD: Unknown column qualifier in Where gosql.Column directive (2)
type SelectAnalysisTestBAD_UnknownColumnQualifierInWhereColumn2 struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		_ gosql.Column `sql:"c.col_a = x.col_a"`
	}
}

//BAD: Unknown column qualifier in gosql.OrderBy directive
type SelectAnalysisTestBAD_UnknownColumnQualifierInOrderBy struct {
	Rel CT1           `rel:"column_tests_1:c"`
	_   gosql.OrderBy `sql:"x.col_a"`
}

//BAD: Unknown column qualifier in Between struct tag
type SelectAnalysisTestBAD_UnknownColumnQualifierInBetween struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		a struct {
			_ gosql.Column `sql:"c.col_b,x"`
			_ gosql.Column `sql:"c.col_c,y"`
		} `sql:"x.col_a isbetween"`
	}
}

//BAD: Unknown column qualifier in Between struct gosql.Column
type SelectAnalysisTestBAD_UnknownColumnQualifierInBetweenColumn struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		a struct {
			_ gosql.Column `sql:"x.col_b,x"`
			_ gosql.Column `sql:"c.col_c,y"`
		} `sql:"c.col_a isbetween"`
	}
}

//BAD: Unknown column qualifier in gosql.Default alias
type InsertAnalysisTestBAD_UnknownColumnQualifierInDefault struct {
	Rel CT1           `rel:"column_tests_1:c"`
	_   gosql.Default `sql:"x.col_b"`
}

//BAD: Join conditional operand on the wrong side
type SelectAnalysisTestBAD_JoinConditionalLHSOperand struct {
	Columns CT1 `rel:"column_tests_1:a"`
	Join    struct {
		_ gosql.LeftJoin `sql:"column_tests_2:b,a.col_b = b.col_bar"`
	}
}

//BAD: Insert with "field-less" column in Return directive
type InsertAnalysisTestBAD_ReturnColumnNoField struct {
	Rel T2           `rel:"relation_a:a"`
	_   gosql.Return `sql:"a.foo,a.bar,a.baz,a.quux"`
}

//BAD: Update with "field-less" column in Force directive
type UpdateAnalysisTestBAD_ForceColumnNoField struct {
	Rel T2          `rel:"relation_a:a"`
	_   gosql.Force `sql:"a.foo,a.bar,a.baz,a.quux"`
}

//BAD: conflicting context field
type InsertAnalysisTestBAD_WithContextConflict struct {
	context.Context
	Rel []T `rel:"relation_a:a"`
	ctx context.Context
}

//OK: Filter with missing constructor
type FilterAnalysisTestBAD_NoFilterConstructor struct {
	_ T                `rel:"relation_a:a"`
	_ gosql.TextSearch `sql:"a.ts_document"`
}

//OK: Filter with conflicting constructors
type FilterAnalysisTestBAD_ConflictingFilterConstructor struct {
	_ T                `rel:"relation_a:a"`
	_ gosql.TextSearch `sql:"a.ts_document"`
	common.FilterMaker
	maker common.FilterMaker
}
