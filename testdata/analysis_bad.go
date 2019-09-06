package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

//BAD: missing relation field
type InsertTestBAD_NoRelfield struct {
	User *common.User
}

//BAD: invalid datatype kind
type InsertTestBAD3 struct {
	User string `rel:"users_table"`
}

//BAD: Delete with invalid relid
type DeleteTestBAD_BadRelId struct {
	Rel T `rel:"foo.123:bar"`
}

//BAD: Delete with illegal count field
type DeleteTestBAD_IllegalCountField struct {
	Count int `rel:"relation_a:a"`
}

//BAD: Update with illegal exists field
type UpdateTestBAD_IllegalExistsField struct {
	Exists bool `rel:"relation_a:a"`
}

//BAD: Insert with illegal notexists field
type InsertTestBAD_IllegalNotExistsField struct {
	NotExists bool `rel:"relation_a:a"`
}

//BAD: Select with illegal gosql.Relation directive
type SelectTestBAD_IllegalRelationDirective struct {
	_ gosql.Relation `rel:"relation_a:a"`
}

//BAD: Select with unnamed base struct type
type SelectTestBAD_UnnamedBaseStructType struct {
	Rel []*struct{} `rel:"relation_a:a"`
}

//BAD: Select with All directive
type SelectTestBAD_IllegalAllDirective struct {
	Rel []T `rel:"relation_a:a"`
	_   gosql.All
}

//BAD: Insert with All directive
type InsertTestBAD_IllegalAllDirective struct {
	Rel T `rel:"relation_a:a"`
	_   gosql.All
}

//BAD: Update with conflicting where producer
type UpdateTestBAD_ConflictWhereProducer struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		Id int `sql:"a.id"`
	}
	_ gosql.All
}

//BAD: Delete with illegal gosql.Default directive
type DeleteTestBAD_IllegalDefaultDirective struct {
	Rel T             `rel:"relation_a:a"`
	_   gosql.Default `sql:"*"`
}

//BAD: Update with empty gosql.Default directive collist
type UpdateTestBAD_EmptyDefaultDirectiveCollist struct {
	Rel T `rel:"relation_a:a"`
	_   gosql.Default
}

//BAD: Select with illegal gosql.Force directive
type SelectTestBAD_IllegalForceDirective struct {
	Rel T           `rel:"relation_a:a"`
	_   gosql.Force `sql:"*"`
}

//BAD: Update with bad gosql.Force directive colid
type UpdateTestBAD_BadForceDirectiveColId struct {
	Rel T           `rel:"relation_a:a"`
	_   gosql.Force `sql:"a.id,1234"`
}

//BAD: Filter with illegal gosql.Return directive
type FilterTestBAD_IllegalReturnDirective struct {
	Rel T            `rel:"relation_a:a"`
	_   gosql.Return `sql:"*"`
}

//BAD: Delete with conflicting result producer
type DeleteTestBAD_ConflictResultProducer struct {
	Rel T            `rel:"relation_a:a"`
	_   gosql.Return `sql:"*"`
	_   gosql.Return `sql:"a.id"`
}

//BAD: Update with empty gosql.Return directive collist
type UpdateTestBAD_EmptyReturnDirectiveCollist struct {
	Rel T `rel:"relation_a:a"`
	_   gosql.Return
}

//BAD: Insert with Limit field
type InsertTestBAD_IllegalLimitField struct {
	Rel T           `rel:"relation_a:a"`
	_   gosql.Limit `sql:"10"`
}

//BAD: Update with Offset field
type UpdateTestBAD_IllegalOffsetField struct {
	Rel T            `rel:"relation_a:a"`
	_   gosql.Offset `sql:"2"`
}

//BAD: Insert with illegal gosql.OrderBy directive
type InsertTestBAD_IllegalOrderByDirective struct {
	Rel T             `rel:"relation_a:a"`
	_   gosql.OrderBy `sql:"a.id"`
}

//BAD: Delete with illegal gosql.Override directive
type DeleteTestBAD_IllegalOverrideDirective struct {
	Rel T              `rel:"relation_a:a"`
	_   gosql.Override `sql:"user"`
}

//BAD: Select with illegal gosql.TextSearch directive
type SelectTestBAD_IllegalTextSearchDirective struct {
	Rel T                `rel:"relation_a:a"`
	_   gosql.TextSearch `sql:"a._document"`
}

//BAD: Select with illegal gosql.Column directive
type SelectTestBAD_IllegalColumnDirective struct {
	Rel T            `rel:"relation_a:a"`
	_   gosql.Column `sql:"a.id"`
}

//BAD: Insert with illegal Where block
type InsertTestBAD_IllegalWhereBlock struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		Id int `sql:"id"`
	}
}

//BAD: Update with illegal Join block
type UpdateTestBAD_IllegalJoinBlock struct {
	Rel  T `rel:"relation_a:a"`
	Join struct {
		_ gosql.Relation
	}
}

//BAD: Delete with illegal From block
type DeleteTestBAD_IllegalFromBlock struct {
	Rel  T `rel:"relation_a:a"`
	From struct {
		_ gosql.Relation
	}
}

//BAD: Select with illegal Using block
type SelectTestBAD_IllegalUsingBlock struct {
	Rel   T `rel:"relation_a:a"`
	Using struct {
		_ gosql.Relation
	}
}

//BAD: Update with illegal OnConflict block
type UpdateTestBAD_IllegalOnConflictBlock struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict struct {
		// ...
	}
}

//BAD: Select with illegal result field
type SelectTestBAD_IllegalResultField struct {
	Rel    T `rel:"relation_a:a"`
	Result T
}

//BAD: Select with conflicting limit producers
type SelectTestBAD_ConflictLimitProducer struct {
	Rel   []T         `rel:"relation_a:a"`
	_     gosql.Limit `sql:"10"`
	Limit int
}

//BAD: Select with conflicting offset producers
type SelectTestBAD_ConflictOffsetProducer struct {
	Rel    []T          `rel:"relation_a:a"`
	_      gosql.Offset `sql:"2"`
	Offset int
}

//BAD: Select with illegal rowsaffected field
type SelectTestBAD_IllegalRowsAffectedField struct {
	Rel          []T `rel:"relation_a:a"`
	RowsAffected int
}

//BAD: Insert with illegal filter field
type InsertTestBAD_IllegalFilterField struct {
	Rel []T `rel:"relation_a:a"`
	F   gosql.Filter
}

//BAD: Select with conflicting where producer
type SelectTestBAD_ConflictWhereProducer struct {
	Rel   []T `rel:"relation_a:a"`
	Where struct {
		Id int `sql:"id"`
	}
	F gosql.Filter
}

//BAD: Delete with conflicting error handlers
type DeleteTestBAD_ConflictWhereProducer struct {
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
type SelectTestBAD_IteratorWithTooManyMethods struct {
	Rel badIterator `rel:"relation_a:a"`
}

//BAD: Select with iterator with bad signature
type SelectTestBAD_IteratorWithBadSignature struct {
	Rel func(*common.User) int `rel:"relation_a:a"`
}

//BAD: Select with iterator with bad signature (interface)
type SelectTestBAD_IteratorWithBadSignatureIface struct {
	Rel badIterator2 `rel:"relation_a:a"`
}

//BAD: Select with imported iterator that has unexported method
type SelectTestBAD_IteratorWithUnexportedMethod struct {
	Rel common.BadIterator `rel:"relation_a:a"`
}

//BAD: Select with iterator with unnamed argument
type SelectTestBAD_IteratorWithUnnamedArgument struct {
	Rel func(*struct{}) error `rel:"relation_a:a"`
}

//BAD: Select with iterator with non-struct argument
type SelectTestBAD_IteratorWithNonStructArgument struct {
	Rel func(*notstruct) error `rel:"relation_a:a"`
}

type notstruct string

//BAD: Insert with bad struct base type
type InsertTestBAD_BadRelfiedlStructBaseType struct {
	Rel []*notstruct `rel:"relation_a:a"`
}

//BAD: Update with bad relfield type field's colid
type UpdateTestBAD_BadRelTypeFieldColId struct {
	Rel struct {
		Foo string `sql:"1234"`
	} `rel:"relation_a:a"`
}

//BAD: Update with conflicting where produceer
type UpdateTestBAD_ConflictWhereProducer2 struct {
	Rel   T `rel:"relation_a:a"`
	_     gosql.All
	Where struct {
		Id int `sql:"id"`
	}
}

//BAD: Delete with bad where block type
type DeleteTestBAD_BadWhereBlockType struct {
	Rel   T `rel:"relation_a:a"`
	Where []string
}

//BAD: Select with bad bool tag value
type SelectTestBAD_BadBoolTagValue struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		Id   int    `sql:"id"`
		Name string `sql:"name" bool:"abc"`
	}
}

//BAD: Select with bad nested where block type
type SelectTestBAD_BadNestedWhereBlockType struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		Id int       `sql:"id"`
		X  notstruct `sql:">"`
	}
}

//BAD: Select with bad gosql.Column expression LHS
type SelectTestBAD_BadColumnExpressionLHS struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		_ gosql.Column `sql:"123 = x"`
	}
}

//BAD: Select with bad gosql.Column cmpop combo
type SelectTestBAD_BadColumnCmpopCombo struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		_ gosql.Column `sql:"x isin any y"`
	}
}

//BAD: Delete with bad gosql.Column expression LHS
type DeleteTestBAD_BadColumnExpressionLHS struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		_ gosql.Column `sql:"123 isnull"`
	}
}

//BAD: Update with bad unary op
type UpdateTestBAD_BadUnaryOp struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		_ gosql.Column `sql:"x <="`
	}
}

//BAD: Update with extra scalar array op
type UpdateTestBAD_ExtraScalarrop struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		_ gosql.Column `sql:"x isnull any"`
	}
}

//BAD: Select with bad between field type
type SelectTestBAD_BadBetweenFieldType struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		between notstruct `sql:"a.foo isbetween"`
	}
}

//BAD: Select with bad number of fields in "between" struct
type SelectTestBAD_BadBetweenFieldType2 struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		between struct {
			x, y, z int
		} `sql:"a.foo isbetween"`
	}
}

//BAD: Select with bad colid in "between" struct field's tag
type SelectTestBAD_BadBetweenArgColId struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		between struct {
			_ gosql.Column `sql:"a.bar,x"`
			_ gosql.Column `sql:"123,y"`
		} `sql:"a.foo isbetween"`
	}
}

//BAD: Select with missing x / y in "between"
type SelectTestBAD_NoBetweenXYArg struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		between struct {
			_ gosql.Column `sql:"a.bar"`
			_ gosql.Column `sql:"a.baz,y"`
		} `sql:"a.foo isbetween"`
	}
}

//BAD: Select with bad "between" target colid
type SelectTestBAD_BadBetweenColId struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		between struct {
			_ gosql.Column `sql:"a.bar,x"`
			_ gosql.Column `sql:"a.baz,y"`
		} `sql:"123 isbetween"`
	}
}

//BAD: Delete with bad where field colid
type DeleteTestBAD_BadWhereFieldColId struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		Id int `sql:"123"`
	}
}

//BAD: Delete with bad where field cmpop combo
type DeleteTestBAD_BadWhereFieldCmpopCombo struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		Id int `sql:"a.id notin any"`
	}
}

//BAD: Delete with illegal where field unary comparison
type DeleteTestBAD_IllegalWhereFieldUnaryCmp struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		Id int `sql:"a.id istrue"`
	}
}

//BAD: Update with bad where field type for scalarrop
type UpdateTestBAD_BadWhereFieldTypeForScalarrop struct {
	Rel   T `rel:"relation_a:a"`
	Where struct {
		Id int `sql:"a.id = any"`
	}
}

//BAD: Select with bad join block type
type SelectTestBAD_BadJoinBlockType struct {
	Rel  T `rel:"relation_a:a"`
	Join notstruct
}

//BAD: Select with illegal join gosql.Relation directive
type SelectTestBAD_IllegalJoinBlockRelationDirective struct {
	Rel  T `rel:"relation_a:a"`
	Join struct {
		_ gosql.Relation `sql:"foobar"`
	}
}

//BAD: Delete with conflicting using block gosql.Relation directive
type DeleteTestBAD_ConflictRelationDirective struct {
	Rel   T `rel:"relation_a:a"`
	Using struct {
		_ gosql.Relation `sql:"foo"`
		_ gosql.Relation `sql:"bar"`
	}
}

//BAD: Update with bad from block gosql.Relation directive relid
type UpdateTestBAD_BadFromRelationRelId struct {
	Rel  T `rel:"relation_a:a"`
	From struct {
		_ gosql.Relation `sql:"123"`
	}
}

//BAD: Select with bad join block gosql.JoinXxx directive relid
type SelectTestBAD_BadJoinDirectiveRelId struct {
	Rel  T `rel:"relation_a:a"`
	Join struct {
		_ gosql.LeftJoin `sql:"123"`
	}
}

//BAD: Select with bad gosql.JoinXxx directive expression colid
type SelectTestBAD_BadJoinDirectiveExpressionColId struct {
	Rel  T `rel:"relation_a:a"`
	Join struct {
		_ gosql.LeftJoin `sql:"relation_b:b,123 = b.foo"`
	}
}

//BAD: Select with bad gosql.JoinXxx directive expression cmpop
type SelectTestBAD_BadJoinDirectiveExpressionCmpop struct {
	Rel  T `rel:"relation_a:a"`
	Join struct {
		_ gosql.LeftJoin `sql:"relation_b:b,b.foo ="`
	}
}

//BAD: Select with bad gosql.JoinXxx directive expression extra scalarrop
type SelectTestBAD_BadJoinDirectiveExpressionExtraScalarrop struct {
	Rel  T `rel:"relation_a:a"`
	Join struct {
		_ gosql.LeftJoin `sql:"relation_b:b,b.foo isnull any"`
	}
}

//BAD: Select with bad gosql.JoinXxx directive expression cmpop combo
type SelectTestBAD_BadJoinDirectiveExpressionCmpopCombo struct {
	Rel  T `rel:"relation_a:a"`
	Join struct {
		_ gosql.LeftJoin `sql:"relation_b:b,b.foo isin any a.bar"`
	}
}

//BAD: Delete with illegal join block directive
type DeleteTestBAD_IllegalJoinBlockDirective struct {
	Rel   T `rel:"relation_a:a"`
	Using struct {
		_ gosql.Column `sql:"a.foo"`
	}
}

//BAD: Insert with bad onconflict block type
type InsertTestBAD_BadOnConflictBlockType struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict notstruct
}

//BAD: Insert with conflicting onconflict block target
type InsertTestBAD_ConflictOnConflictBlockTargetProducer struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.Index  `sql:"some_index"`
		_ gosql.Column `sql:"a.id"`
	}
}

//BAD: Insert with conflicting onconflict block target
type InsertTestBAD_ConflictOnConflictBlockTargetProducer2 struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.Column `sql:"a.id"`
		_ gosql.Index  `sql:"some_index"`
	}
}

//BAD: Insert with conflicting onconflict block target
type InsertTestBAD_ConflictOnConflictBlockTargetProducer3 struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.Index      `sql:"some_index"`
		_ gosql.Constraint `sql:"some_constraint"`
	}
}

//BAD: Insert with conflicting onconflict block action
type InsertTestBAD_ConflictOnConflictBlockActionProducer struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.Column `sql:"a.id"`
		_ gosql.Update `sql:"a.foo"`
		_ gosql.Ignore
	}
}

//BAD: Insert with conflicting onconflict block action
type InsertTestBAD_ConflictOnConflictBlockActionProducer2 struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.Column `sql:"a.id"`
		_ gosql.Ignore
		_ gosql.Update `sql:"a.foo"`
	}
}

//BAD: Insert with bad onconflict column target value
type InsertTestBAD_BadOnConflictColumnTargetValue struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.Column `sql:"a.id,a.1234"`
	}
}

//BAD: Insert with bad onconflict index target identifier
type InsertTestBAD_BadOnConflictIndexTargetIdent struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.Index `sql:"1234"`
	}
}

//BAD: Insert with bad onconflict constraint target identifier
type InsertTestBAD_BadOnConflictConstraintTargetIdent struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.Constraint `sql:"1234"`
	}
}

//BAD: Insert with bad onconflict update action collist
type InsertTestBAD_BadOnConflictUpdateActionCollist struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.Column `sql:"a.id"`
		_ gosql.Update `sql:"a.id,a.1234"`
	}
}

//BAD: Insert with illegal onconflict directive
type InsertTestBAD_IllegalOnConflictDirective struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.LeftJoin `sql:"relation_b:b,b.id=a.id"`
	}
}

//BAD: Insert with illegal onconflict directive
type InsertTestBAD_NoOnConflictTarget struct {
	Rel        T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.Update `sql:"a.foo,a.bar"`
	}
}

//BAD: Select with bad limit field type
type SelectTestBAD_BadLimitFieldType struct {
	Rel   []T    `rel:"relation_a:a"`
	Limit string `sql:"123"`
}

//BAD: Select with no limit directive value
type SelectTestBAD_NoLimitDirectiveValue struct {
	Rel []T         `rel:"relation_a:a"`
	_   gosql.Limit `sql:""`
}

//BAD: Select with bad limit directive value
type SelectTestBAD_BadLimitDirectiveValue struct {
	Rel []T         `rel:"relation_a:a"`
	_   gosql.Limit `sql:"abc"`
}

//BAD: Select with bad offset field type
type SelectTestBAD_BadOffsetFieldType struct {
	Rel    []T    `rel:"relation_a:a"`
	Offset string `sql:"123"`
}

//BAD: Select with no offset directive value
type SelectTestBAD_NoOffsetDirectiveValue struct {
	Rel []T          `rel:"relation_a:a"`
	_   gosql.Offset `sql:""`
}

//BAD: Select with bad offset directive value
type SelectTestBAD_BadOffsetDirectiveValue struct {
	Rel []T          `rel:"relation_a:a"`
	_   gosql.Offset `sql:"abc"`
}

//BAD: Select with empty gosql.OrderBy directive list
type SelectTestBAD_EmptyOrderByDirectiveCollist struct {
	Rel []T `rel:"relation_a:a"`
	_   gosql.OrderBy
}

//BAD: Select with bad gosql.OrderBy directive nulls order option value
type SelectTestBAD_BadOrderByDirectiveNullsOrderValue struct {
	Rel []T           `rel:"relation_a:a"`
	_   gosql.OrderBy `sql:"a.id:nullsthird"`
}

//BAD: Select with bad gosql.OrderBy directive colid
type SelectTestBAD_BadOrderByDirectiveCollist struct {
	Rel []T           `rel:"relation_a:a"`
	_   gosql.OrderBy `sql:"-a.id:nullsfirst,a.1234"`
}

//BAD: Insert with bad gosql.Override directive kind value
type InsertTestBAD_BadOverrideDirectiveKindValue struct {
	Rel []T            `rel:"relation_a:a"`
	_   gosql.Override `sql:"foo"`
}

//BAD: Update with conflicting result producer
type UpdateTestBAD_ConflictResultProducer struct {
	Rel    T            `rel:"relation_a:a"`
	_      gosql.Return `sql:"*"`
	Result []T
}

//BAD: Update with bad result field type
type UpdateTestBAD_BadResultFieldType struct {
	Rel    T `rel:"relation_a:a"`
	Result []notstruct
}

//BAD: Delete with conflicting result producer
type DeleteTestBAD_ConflictResultProducer2 struct {
	Rel          T `rel:"relation_a:a"`
	Result       []T
	RowsAffected int
}

//BAD: Delete with bad rowsaffected field type
type DeleteTestBAD_BadRowsAffecteFieldType struct {
	Rel          T `rel:"relation_a:a"`
	RowsAffected string
}

//BAD: Filter with bad gosql.TextSearch directive colid
type FilterTestBAD_BadTextSearchDirectiveColId struct {
	Rel T                `rel:"relation_a:a"`
	_   gosql.TextSearch `sql:"123"`
}
