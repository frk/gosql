package postgres

import (
	"fmt"

	"github.com/frk/gosql/internal/analysis"
)

/*
postgres.ErrorObject{
	Code:0xd,
	PkgPath:"../testdata/generator/pgsql",
	DBName:"gosql_test_db",
	RelName:"pgsql_test",
	RelSchema:"public",
	ColName:"col_bpchar",
	ColQualifier:"",
	ColType:"character",
	Col2Name:"",
	Col2Qualifier:"",
	Col2Type:"",
	LitExpr:"",
	Predicate:0x0,
	Quantifier:0x0,
	TargetName:"insertbasicquery",
	FieldType:"byte",
	FieldName:"bpchar",
	FieldTag:"sql:\"col_bpchar\"",
	FileName:"../testdata/generator/pgsql/datatype.go",
	FileLine:10}
*/
type ErrorCode uint

const (
	_ ErrorCode = iota
	ErrNoRelation
	ErrNoColumn
	ErrNoColumnType
	ErrNoColumnDefault
	ErrNoColumnField
	ErrNoProc
	ErrNoType
	ErrNoUniqueIndex
	ErrNoUniqueConstraint
	ErrBadAlias
	ErrBadColumnType
	ErrBadColumnQualifier
	ErrBadFieldWriteType
	ErrBadFieldReadType
	ErrBadProcType
	ErrBadUnaryPredicateType
	ErrBadLiteralExpression
	ErrBadPredicateQuantifierType
	ErrBadComparisonOperation
)

type Error struct {
	Code          ErrorCode
	PkgPath       string
	DBName        string
	RelName       string
	RelSchema     string
	ColName       string
	ColQualifier  string
	ColType       string
	Col2Name      string
	Col2Qualifier string
	Col2Type      string
	LitExpr       string
	Predicate     analysis.Predicate
	Quantifier    analysis.Quantifier
	TargetName    string
	FieldType     string
	FieldName     string
	FieldTag      string
	FileName      string
	FileLine      int
}

func (e Error) Error() string {
	type ErrorObject Error
	return fmt.Sprintf("[ TODO ERROR MESG ] %#v\n\n", ErrorObject(e))
}

type eopt struct {
	c    ErrorCode
	rid  analysis.RelIdent
	cid  analysis.ColIdent
	cid2 analysis.ColIdent
	ptr  analysis.FieldPtr
	lit  string
	rel  *Relation
	col  *Column
	col2 *Column
	typ2 *Type
	pre  analysis.Predicate
	qua  analysis.Quantifier
}

// error constructs and returns a new Error value.
func (c *checker) error2(opt eopt) error {
	e := Error{Code: opt.c}
	e.PkgPath = c.info.PkgPath
	e.TargetName = c.info.TypeName
	e.DBName = c.db.name
	e.LitExpr = opt.lit
	e.Predicate = opt.pre
	e.Quantifier = opt.qua

	// get rel info from *Relation or RelIdent
	if opt.rel != nil {
		e.RelName = opt.rel.Name
		e.RelSchema = opt.rel.Schema
	} else if opt.rid.Name != "" {
		e.RelName = opt.rid.Name
		e.RelSchema = opt.rid.Qualifier
	}

	// get col info from *Column or ColIdent
	if opt.col != nil {
		e.ColName = opt.col.Name
		e.ColType = opt.col.Type.NameFmt
		e.ColQualifier = opt.cid.Qualifier

		// If no rel info exists try the col.Relation
		if e.RelName == "" && opt.col.Relation != nil {
			e.RelName = opt.col.Relation.Name
			e.RelSchema = opt.col.Relation.Schema
		}
	} else if opt.cid.Name != "" {
		e.ColName = opt.cid.Name
		e.ColQualifier = opt.cid.Qualifier

		// If no rel info exists try the relMap
		if e.RelName == "" {
			if rel, ok := c.relMap[e.ColQualifier]; ok {
				e.RelName = rel.Name
				e.RelSchema = rel.Schema
			}
		}
	}

	// get col2 info from *Column or ColIdent
	if opt.col2 != nil {
		e.Col2Name = opt.col2.Name
		e.Col2Type = opt.col2.Type.NameFmt
		e.Col2Qualifier = opt.cid2.Qualifier
	} else if opt.cid2.Name != "" {
		e.ColName = opt.cid2.Name
		e.ColQualifier = opt.cid2.Qualifier
	}

	if opt.typ2 != nil {
		e.Col2Type = opt.typ2.NameFmt
	}

	// get field info from FieldPtr
	if f, ok := c.info.FieldMap[opt.ptr]; ok {
		p := c.info.FileSet.Position(f.Var.Pos())
		e.FieldType = f.Var.Type().String()
		e.FieldName = f.Var.Name()
		e.FieldTag = f.Tag
		e.FileName = p.Filename
		e.FileLine = p.Line
	}
	return e
}
