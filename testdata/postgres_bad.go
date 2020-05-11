// WARNING(mkopriva): IF NOT NECESSARY TRY NOT TO DO ANYTHING THAT WILL REORDER
// LINES OR ADD NEW LINES IN THE MIDDLE OF THE FILE, ONLY ADD NEW CODE AT THE TAIL END.
package testdata

import (
	"time"

	"github.com/frk/gosql"
)

// column_tests_1 record
type CT1 struct {
	A int       `sql:"col_a"`
	B string    `sql:"col_b"`
	C bool      `sql:"col_c"`
	D float64   `sql:"col_d"`
	E time.Time `sql:"col_e"`
}

//BAD: relation does not exist
type SelectPostgresTestBAD_NoRelation struct {
	Columns CT1 `rel:"norel"`
}

//BAD: Join relation does not exist
type DeletePostgresTestBAD_JoinNoRelation struct {
	Rel   CT1 `rel:"column_tests_1:a"`
	Using struct {
		_ gosql.Relation `sql:"norel:b"`
	}
	Where struct {
		_ gosql.Column `sql:"a.col_a = b.col_a"`
	}
}

//BAD: Join relation does not exist
type DeletePostgresTestBAD_JoinNoRelation2 struct {
	Rel   CT1 `rel:"column_tests_1:a"`
	Using struct {
		_ gosql.Relation `sql:"column_tests_2:b"`
		_ gosql.LeftJoin `sql:"norel:c,c.b_id = b.id"`
	}
	Where struct {
		_ gosql.Column `sql:"a.col_a = b.col_foo"`
	}
}

//BAD: Alias referenced in join doesn't match the joined table's alias
type SelectPostgresTestBAD_JoinNoAliasRelation struct {
	Columns CT1 `rel:"column_tests_1:a"`
	Join    struct {
		_ gosql.LeftJoin `sql:"column_tests_2:b,x.col_foo = a.col_a"`
	}
}

//BAD: Alias referenced in join doesn't match any known relation
type SelectPostgresTestBAD_JoinNoAliasRelation2 struct {
	Columns CT1 `rel:"column_tests_1:a"`
	Join    struct {
		_ gosql.LeftJoin `sql:"column_tests_2:b,b.col_foo = x.col_a"`
	}
}

//BAD: Join column does not exist
type SelectPostgresTestBAD_JoinNoColumn struct {
	Columns CT1 `rel:"column_tests_1:a"`
	Join    struct {
		_ gosql.LeftJoin `sql:"column_tests_2:b,b.nocol = a.nocol"`
	}
}

//BAD: Join column does not exist
type SelectPostgresTestBAD_JoinNoColumn2 struct {
	Columns CT1 `rel:"column_tests_1:a"`
	Join    struct {
		_ gosql.LeftJoin `sql:"column_tests_2:b,b.col_foo = a.nocol"`
	}
}

//BAD: Join column type cannot be used in bool predicate
type SelectPostgresTestBAD_JoinBadUnaryBoolColumn struct {
	Columns CT1 `rel:"column_tests_1:a"`
	Join    struct {
		_ gosql.LeftJoin `sql:"column_tests_2:b,b.col_foo istrue"`
	}
}

//BAD: Join column with NOT NULL cannot be used in null predicate
type SelectPostgresTestBAD_JoinBadUnaryNullColumn struct {
	Columns CT1 `rel:"column_tests_1:a"`
	Join    struct {
		_ gosql.LeftJoin `sql:"column_tests_2:b,b.col_baz isnull"`
	}
}

//BAD: Join column with bad literal expression
type SelectPostgresTestBAD_JoinBadLiteralExpression struct {
	Columns CT1 `rel:"column_tests_1:a"`
	Join    struct {
		_ gosql.LeftJoin `sql:"column_tests_2:b,b.col_baz = 'foo'bar "`
	}
}

//BAD: Join column with bad quantifier colum type
type SelectPostgresTestBAD_JoinBadQuantifierColumnType struct {
	Columns CT1 `rel:"column_tests_1:a"`
	Join    struct {
		_ gosql.LeftJoin `sql:"column_tests_2:b,b.col_foo >any a.col_a"`
	}
}

//BAD: Join column with bad comparison operands' types
type SelectPostgresTestBAD_JoinBadComparisonOperandType struct {
	Columns CT1 `rel:"column_tests_1:a"`
	Join    struct {
		_ gosql.LeftJoin `sql:"column_tests_2:b,b.col_baz < 'baz'"`
	}
}

//BAD: onconflict block column target not found
type InsertPostgresTestBAD_OnConflictNoColumn struct {
	Rel        CT1 `rel:"column_tests_1:c"`
	OnConflict struct {
		_ gosql.Column `sql:"c.col_xyz"`
	}
}

//BAD: onconflict block column don't match any unique index
type InsertPostgresTestBAD_OnConflictColumnNoIndexMatch struct {
	Rel        CT1 `rel:"column_tests_1:c"`
	OnConflict struct {
		_ gosql.Column `sql:"c.col_a,c.col_b"`
	}
}

//BAD: onconflict block index not found
type InsertPostgresTestBAD_OnConflictNoIndex struct {
	Rel        CT1 `rel:"column_tests_1:c"`
	OnConflict struct {
		_ gosql.Index `sql:"some_index"`
	}
}

//BAD: onconflict block unique index not found
type InsertPostgresTestBAD_OnConflictNoUniqueIndex struct {
	Rel        CT1 `rel:"column_tests_2:c"`
	OnConflict struct {
		_ gosql.Index `sql:"column_tests_2_nonunique_index"`
	}
}

//BAD: onconflict block constraint not found
type InsertPostgresTestBAD_OnConflictNoConstraint struct {
	Rel        CT1 `rel:"column_tests_1:c"`
	OnConflict struct {
		_ gosql.Constraint `sql:"some_constraint"`
	}
}

//BAD: onconflict block constraint not found
type InsertPostgresTestBAD_OnConflictNoUniqueConstraint struct {
	Rel        CT1 `rel:"column_tests_2:c"`
	OnConflict struct {
		_ gosql.Constraint `sql:"column_tests_2_nonunique_constraint"`
	}
}

//BAD: onconflict block update column not found
type InsertPostgresTestBAD_OnConflictUpdateColumnNotFound struct {
	Rel        CT1 `rel:"column_tests_2:c"`
	OnConflict struct {
		_ gosql.Constraint `sql:"column_tests_2_unique_constraint"`
		_ gosql.Update     `sql:"c.col_a,c.col_b,c.col_xyz"`
	}
}

//BAD: whereblock field not found
type SelectPostgresTestBAD_WhereFieldNotFound struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		Id int `sql:"c.id"`
	}
}

//BAD: whereblock aliased relation not found
type SelectPostgresTestBAD_WhereAliasNotFound struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		Id int `sql:"x.id"`
	}
}

//BAD: whereblock aliased relation not found
type SelectPostgresTestBAD_WherePointerFieldForNonNullColumn struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		B *string `sql:"c.col_b"`
	}
}

//BAD: whereblock wrong field type for quantifier
type SelectPostgresTestBAD_WhereBadFieldTypeForQuantifier struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		B string `sql:"c.col_b isin"`
	}
}

//BAD: whereblock cannot compare types
type SelectPostgresTestBAD_WhereCannotCompareTypes struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		D float64 `sql:"c.col_e"`
	}
}

//BAD: whereblock with bad argument type for funcname
type SelectPostgresTestBAD_WhereColumnTypeForFuncname struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		D float64 `sql:"c.col_d,@lower"`
	}
}

//BAD: whereblock with column not found
type SelectPostgresTestBAD_WhereColumnNotFound struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		_ gosql.Column `sql:"c.col_xyz istrue"`
	}
}

//BAD: whereblock with column not found 2 (bad alias)
type SelectPostgresTestBAD_WhereColumnNotFoundBadAlias struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		_ gosql.Column `sql:"x.col_a = 123"`
	}
}

//BAD: whereblock with column bad bool operation
type SelectPostgresTestBAD_WhereColumnBadBoolOp struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		_ gosql.Column `sql:"c.col_a istrue"`
	}
}

//BAD: whereblock with column bad NULL operation
type SelectPostgresTestBAD_WhereColumnBadNULLOp struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		_ gosql.Column `sql:"c.col_b isnull"`
	}
}

//BAD: whereblock with RHS column not found
type SelectPostgresTestBAD_WhereColumnNotFoundRHS struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		_ gosql.Column `sql:"c.col_a = c.col_xyz"`
	}
}

//BAD: whereblock with RHS column not found 2 (bad alias)
type SelectPostgresTestBAD_WhereColumnNotFoundRHSBadAlias struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		_ gosql.Column `sql:"c.col_a = x.col_a"`
	}
}

//BAD: whereblock with column bad literal expression
type SelectPostgresTestBAD_WhereColumnBadLiteralExpression struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		_ gosql.Column `sql:"c.col_a = 123abc"`
	}
}

//BAD: whereblock wrong column type for quantifier
type SelectPostgresTestBAD_WhereColumnBadTypeForQuantifier struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		_ gosql.Column `sql:"c.col_a isin c.col_b"`
	}
}

//BAD: whereblock wrong column type for comparison
type SelectPostgresTestBAD_WhereColumnBadTypeComparison struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		_ gosql.Column `sql:"c.col_a = c.col_b"`
	}
}

//BAD: whereblock between column not found
type SelectPostgresTestBAD_WhereBetweenColumnNotFound struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		a struct {
			_ gosql.Column `sql:"c.col_a,x"`
			_ gosql.Column `sql:"c.col_c,y"`
		} `sql:"c.col_xyz isbetween"`
	}
}

//BAD: whereblock between relation not found (bad alias)
type SelectPostgresTestBAD_WhereBetweenRelationNotFound struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		a struct {
			_ gosql.Column `sql:"c.col_b,x"`
			_ gosql.Column `sql:"c.col_c,y"`
		} `sql:"x.col_a isbetween"`
	}
}

//BAD: whereblock between column not found
type SelectPostgresTestBAD_WhereBetweenArgColumnNotFound struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		a struct {
			_ gosql.Column `sql:"c.col_xyz,x"`
			_ gosql.Column `sql:"c.col_c,y"`
		} `sql:"c.col_a isbetween"`
	}
}

//BAD: whereblock between relation not found (bad alias)
type SelectPostgresTestBAD_WhereBetweenArgRelationNotFound struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		a struct {
			_ gosql.Column `sql:"x.col_b,x"`
			_ gosql.Column `sql:"c.col_c,y"`
		} `sql:"c.col_a isbetween"`
	}
}

//BAD: whereblock between comparison bad arg type
type SelectPostgresTestBAD_WhereBetweenComparisonBadArgType struct {
	Rel   CT1 `rel:"column_tests_1:c"`
	Where struct {
		a struct {
			x int    `sql:"x"`
			y string `sql:"y"`
		} `sql:"c.col_a isbetween"`
	}
}

//BAD: orderby column not found
type SelectPostgresTestBAD_OrderByColumnNotFound struct {
	Rel CT1           `rel:"column_tests_1:c"`
	_   gosql.OrderBy `sql:"c.col_a,c.col_xyz"`
}

//BAD: orderby relation not found
type SelectPostgresTestBAD_OrderByRelationNotFound struct {
	Rel CT1           `rel:"column_tests_1:c"`
	_   gosql.OrderBy `sql:"x.col_a"`
}

//BAD: default bad relation alias
type InsertPostgresTestBAD_DefaultBadRelationAlias struct {
	Rel CT1           `rel:"column_tests_1:c"`
	_   gosql.Default `sql:"x.col_b"`
}

//BAD: default column not found
type InsertPostgresTestBAD_DefaultColumnNotFound struct {
	Rel CT1           `rel:"column_tests_1:c"`
	_   gosql.Default `sql:"c.col_xyz"`
}

//BAD: default not set
type InsertPostgresTestBAD_DefaultNotSet struct {
	Rel CT1           `rel:"column_tests_1:c"`
	_   gosql.Default `sql:"c.col_b"`
}

//BAD: force column not found
type InsertPostgresTestBAD_ForceColumnNotFound struct {
	Rel CT1         `rel:"column_tests_1:c"`
	_   gosql.Force `sql:"c.col_xyz"`
}

//BAD: force relation not found
type InsertPostgresTestBAD_ForceRelationNotFound struct {
	Rel CT1         `rel:"column_tests_1:c"`
	_   gosql.Force `sql:"x.col_a"`
}

//BAD: returning column not found
type UpdatePostgresTestBAD_ReturnColumnNotFound struct {
	Rel CT1          `rel:"column_tests_1:c"`
	_   gosql.Return `sql:"c.col_xyz"`
}

//BAD: returning relation not found
type UpdatePostgresTestBAD_ReturnRelationNotFound struct {
	Rel CT1          `rel:"column_tests_1:c"`
	_   gosql.Return `sql:"x.col_a"`
}

//BAD: textsearch column not found
type FilterPostgresTestBAD_TextSearchColumnNotFound struct {
	_ CT1              `rel:"column_tests_1:c"`
	_ gosql.TextSearch `sql:"c.col_xyz"`
}

//BAD: textsearch relation not found (bad alias)
type FilterPostgresTestBAD_TextSearchRelationNotFound struct {
	_ CT1              `rel:"column_tests_1:c"`
	_ gosql.TextSearch `sql:"x.col_b"`
}

//BAD: textsearch bad column type
type FilterPostgresTestBAD_TextSearchBadColumnType struct {
	_ CT1              `rel:"column_tests_1:c"`
	_ gosql.TextSearch `sql:"c.col_b"`
}

//BAD: target relation column not found
type SelectPostgresTestBAD_RelationColumnNotFound struct {
	Rel struct {
		Xyz string `sql:"col_xyz"`
	} `rel:"column_tests_1:c"`
}

//BAD: target relation column not found (bad alias)
type SelectPostgresTestBAD_RelationColumnAliasNotFound struct {
	Rel struct {
		Xyz string `sql:"x.col_a"`
	} `rel:"column_tests_1:c"`
}

//BAD: target relation column not found
type InsertPostgresTestBAD_RelationColumnNotFound struct {
	Rel struct {
		Xyz string `sql:"col_xyz"`
	} `rel:"column_tests_1:c"`
}

//BAD: JSON option on non json column
type InsertPostgresTestBAD_BadJSONOption struct {
	Rel struct {
		B string `sql:"col_b,json"`
	} `rel:"column_tests_1:c"`
}

//BAD: XML option on non xml column
type InsertPostgresTestBAD_BadXMLOption struct {
	Rel struct {
		B string `sql:"col_b,xml"`
	} `rel:"column_tests_1:c"`
}

//BAD: field type incompatible with column type
type InsertPostgresTestBAD_BadFieldToColumnType struct {
	Rel struct {
		B int `sql:"col_c"`
	} `rel:"column_tests_1:c"`
}

//BAD: result column not found
type InsertPostgresTestBAD_ResultColumnNotFound struct {
	Rel    CT1 `rel:"column_tests_1:c"`
	Result struct {
		A int `sql:"col_xyz"`
	}
}
